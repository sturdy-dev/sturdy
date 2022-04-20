package executor

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"path"
	"runtime/debug"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"getsturdy.com/api/pkg/codebases"
	"getsturdy.com/api/vcs"
	"getsturdy.com/api/vcs/provider"

	"go.uber.org/zap"
)

// Executor proves an interface to schedule git operations, to be executed synchronously, in an expected state,
// and exclusively.
//
// This executor replaces earlier direct access to provider.RepoProvider, which when using the executor will be injected instead.
type Executor interface {
	// Read schedules a function that can read files from the filesystem.
	Read(FileReadFunc) Executor
	// Write schedules a function that can write files to the filesystem.
	Write(FileWriteFunc) Executor
	// GitRead schedules a function that can read .git repository, but not the files on the filesystem.
	GitRead(fn GitReadFunc) Executor
	// GitWrite schedules a function that can read and write .git repository, but not the files on the filesystem.
	GitWrite(fn GitWriteFunc) Executor
	// FileReadGitWrite schedules a function that can read files from the filesystem, as well as read/write files in .git
	FileReadGitWrite(fn FileReadGitWriteFunc) Executor

	// Schedule schedules a function that can access repo provider. It is deprecated, prefer to use
	// Read, Write or Git functions instead to get more granular locking.
	Schedule(fn ExecuteFunc) Executor

	// AllowRebasingState allows the executor to open a repo in a rebasing state. By default this is not allowed.
	AllowRebasingState() Executor
	// AssertBranchName asserts that the repo is in the expected branch.
	AssertBranchName(string) Executor

	// ExecView executes all of the scheduled functions for the given view repository.
	ExecView(codebaseID codebases.ID, viewID, actionName string) error
	// ExecTrunk executes all of the scheduled functions for the given trunk repository.
	ExecTrunk(codebaseID codebases.ID, actionName string) error
	// ExecTemporaryView creates a view for the given codebase, cloning it from the trunk,
	// executes all the scheduled functions, and then deletes the view.
	ExecTemporaryView(codebaseID codebases.ID, actionName string) error
}

type executeFunc struct {
	fun ExecuteFunc

	repoReaderFun GitReadFunc  // read .git
	repoWriterFun GitWriteFunc // write .git

	fileReadFun  FileReadFunc  // read fs
	fileWriteFun FileWriteFunc // write fs

	fileReadGitWriteFun FileReadGitWriteFunc // read fs and write .git
}

func (e *executeFunc) Exec(repo *onceRepo) error {
	if e.fun != nil {
		return e.fun(repo.repoProvider)
	}

	r, err := repo.Get()
	if err != nil {
		return err
	}

	if e.repoReaderFun != nil {
		return e.repoReaderFun(r)
	}

	if e.repoWriterFun != nil {
		return e.repoWriterFun(r)
	}

	if e.fileReadGitWriteFun != nil {
		return e.fileReadGitWriteFun(r)
	}

	if e.fileWriteFun != nil {
		return e.fileWriteFun(r)
	}

	if e.fileReadFun != nil {
		return e.fileReadFun(r)
	}

	return nil
}

func (e *executeFunc) Write() bool {
	return e.fun != nil && e.fileWriteFun != nil
}

func (e *executeFunc) Read() bool {
	return !e.Write() && e.fileReadFun != nil
}

type executor struct {
	funs []*executeFunc

	// file system locks
	writeLock bool
	readLock  bool

	// in memory locks
	inMemoryReadLock  bool
	inMemoryWriteLock bool

	allowRebasing bool

	logger       *zap.Logger
	repoProvider provider.RepoProvider
	locks        *locker
}

func newExecutor(logger *zap.Logger, repoProvider provider.RepoProvider, locks *locker) *executor {
	return &executor{
		logger:       logger,
		repoProvider: repoProvider,
		locks:        locks,
	}
}

// ExecuteFunc is deprecated. Use ExecuteRepoFunc instead
type ExecuteFunc func(repoProvider provider.RepoProvider) error

// Schedule is deprecated. Use ScheduleRepo instead
func (e *executor) Schedule(fn ExecuteFunc) Executor {
	e.writeLock = true
	e.funs = append(e.funs, &executeFunc{
		fun: fn,
	})
	return e
}

type FileWriteFunc func(vcs.RepoWriter) error

func (e *executor) Write(fn FileWriteFunc) Executor {
	e.writeLock = true
	e.funs = append(e.funs, &executeFunc{
		fileWriteFun: fn,
	})
	return e
}

type FileReadFunc func(vcs.RepoReader) error

func (e *executor) Read(fn FileReadFunc) Executor {
	e.readLock = true
	e.funs = append(e.funs, &executeFunc{
		fileReadFun: fn,
	})
	return e
}

type GitReadFunc func(repo vcs.RepoGitReader) error
type GitWriteFunc func(repo vcs.RepoGitWriter) error
type FileReadGitWriteFunc func(repo vcs.RepoReaderGitWriter) error

func (e *executor) GitRead(fn GitReadFunc) Executor {
	e.inMemoryReadLock = true
	e.funs = append(e.funs, &executeFunc{
		repoReaderFun: fn,
	})
	return e
}

func (e *executor) GitWrite(fn GitWriteFunc) Executor {
	e.inMemoryWriteLock = true
	e.funs = append(e.funs, &executeFunc{
		repoWriterFun: fn,
	})
	return e
}

func (e *executor) FileReadGitWrite(fn FileReadGitWriteFunc) Executor {
	e.inMemoryWriteLock = true
	e.readLock = true

	e.funs = append(e.funs, &executeFunc{
		fileReadGitWriteFun: fn,
	})
	return e
}

func (e *executor) gitFirst(fn GitReadFunc) Executor {
	return e.prepend(&executeFunc{repoReaderFun: fn})
}

func (e *executor) prepend(f *executeFunc) Executor {
	e.funs = append([]*executeFunc{f}, e.funs...)
	return e
}

func (e *executor) AllowRebasingState() Executor {
	e.allowRebasing = true
	return e
}

func (e *executor) AssertBranchName(name string) Executor {
	return e.GitRead(func(repo vcs.RepoGitReader) error {
		if repo.IsRebasing() {
			return nil
		}

		headBranch, err := repo.HeadBranch()
		if err != nil && !e.allowRebasing {
			// If the repo is currently rebasing, there won't be any head branch
			// don't error out if rebases are allowed
			return fmt.Errorf("no head branch could be found (and rebases are not allowed): %w", err)

		}

		if err == nil && headBranch != name {
			return fmt.Errorf("expected %s got %s: %w", name, headBranch, ErrUnexpectedBranch)
		}

		return nil
	})
}

var ErrIsRebasing = fmt.Errorf("unexpected git executor state, is rebasing")
var ErrUnexpectedBranch = fmt.Errorf("unexpected git executor state, on unexpected branch")

const (
	inUsePrefix = "using-"
	tmpPrefix   = "tmp-"
)

func filter[T any](slice []T, f func(T) bool) []T {
	filtered := make([]T, 0, len(slice))
	for _, s := range slice {
		if f(s) {
			filtered = append(filtered, s)
		}
	}
	return filtered
}

func isTemporaryView(name string) bool {
	return strings.HasPrefix(name, tmpPrefix)
}

func random[T any](slice []T) T {
	randomIndex := rand.Intn(len(slice))
	return slice[randomIndex]
}

func (e *executor) getTemporaryViewID(codebaseID codebases.ID) (string, error) {
	codebasePath := path.Dir(e.repoProvider.TrunkPath(codebaseID))
	codebaseDir, err := os.Open(codebasePath)
	if err != nil {
		return "", fmt.Errorf("failed to open codebase path: %w", err)
	}

	// get all existing paths
	viewPaths, err := codebaseDir.Readdirnames(-1)
	if err != nil {
		return "", fmt.Errorf("failed to list views: %w", err)
	}

	// find temporary views among them
	temporaryViews := filter(viewPaths, isTemporaryView)
	if len(temporaryViews) == 0 {
		// add a new tmp view to the pool
		return fmt.Sprintf("%s%s%s", inUsePrefix, tmpPrefix, uuid.NewString()), nil
	}

	// re-use existing view from the pool
	viewID := random(temporaryViews)
	return e.temporaryViewMarkUsing(codebaseID, viewID)
}

func (e *executor) temporaryViewMarkUsing(codebaseID codebases.ID, viewID string) (string, error) {
	if strings.HasPrefix(viewID, inUsePrefix) {
		return "", fmt.Errorf("view already in use")
	}

	codebasePath := path.Dir(e.repoProvider.TrunkPath(codebaseID))
	inUseID := fmt.Sprintf("%s%s", inUsePrefix, viewID)
	if err := os.Rename(
		path.Join(codebasePath, viewID),
		path.Join(codebasePath, inUseID),
	); err != nil {
		return "", fmt.Errorf("failed to mark view as in use: %w", err)
	}
	return inUseID, nil
}

func (e *executor) temporaryViewMarkNotUsing(codebaseID codebases.ID, viewID string) (string, error) {
	if !strings.HasPrefix(viewID, inUsePrefix) {
		return viewID, nil
	}

	codebasesDir := path.Dir(e.repoProvider.TrunkPath(codebaseID))
	notInUse := strings.TrimLeft(viewID, inUsePrefix)
	if err := os.Rename(
		path.Join(codebasesDir, viewID),
		path.Join(codebasesDir, notInUse),
	); err != nil {
		return "", fmt.Errorf("failed to mark view as not using: %w", err)
	}
	return notInUse, nil
}

func (e *executor) ExecTemporaryView(codebaseID codebases.ID, actionName string) error {
	e.allowRebasing = true

	viewID, err := e.getTemporaryViewID(codebaseID)
	if err != nil {
		return err
	}

	codebasePath := e.repoProvider.TrunkPath(codebaseID)
	defer func() {
		if _, err := e.temporaryViewMarkNotUsing(codebaseID, viewID); err != nil {
			e.logger.Warn("failed to return tmp view to the pool: %w", zap.Error(err))
		}
	}()
	return e.prepend(&executeFunc{
		fun: func(repoProvider provider.RepoProvider) error {
			viewPath := repoProvider.ViewPath(codebaseID, viewID)
			if _, err := os.Stat(viewPath); errors.Is(err, os.ErrNotExist) {
				if _, err := vcs.CloneRepo(codebasePath, viewPath); err != nil {
					return fmt.Errorf("failed to clone repo: %w", err)
				}
				return nil
			} else if err != nil {
				return fmt.Errorf("failed to stat view: %w", err)
			} else {
				// checkout sturdy trunk as if it is a new view
				repo, err := repoProvider.ViewRepo(codebaseID, viewID)
				if err != nil {
					return fmt.Errorf("failed to open view repo: %w", err)
				}
				if err := repo.FetchBranch("sturdytrunk"); err != nil {
					return fmt.Errorf("failed to fetch branch: %w", err)
				}
				if err := repo.CheckoutBranchWithForce("sturdytrunk"); err != nil {
					return fmt.Errorf("failed to checkout branch: %w", err)
				}
				return nil
			}
		},
	}).ExecView(codebaseID, viewID, actionName)
}

func (e *executor) ExecView(codebaseID codebases.ID, viewID, actionName string) error {
	return e.exec(codebaseID, &viewID, actionName)
}

func (e *executor) ExecTrunk(codebaseID codebases.ID, actionName string) error {
	return e.exec(codebaseID, nil, actionName)
}

func (e *executor) exec(codebaseID codebases.ID, viewID *string, actionName string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("recovered call in vcs executor: %v\nStacktrace: %s", r, string(debug.Stack()))
		}
	}()

	logger := e.logger.With(zap.Stringer("codebase_id", codebaseID), zap.String("action_name", actionName))
	if viewID != nil {
		logger = logger.With(zap.String("view_id", *viewID))
	}

	t0 := time.Now()
	execT0 := time.Now() // is overwritten when the exec starts
	defer func() {
		logger.Info("git executor completed",
			zap.Duration("duration", time.Since(t0)),
			zap.Duration("exec_duration", time.Since(execT0)))
	}()

	if !e.allowRebasing {
		e.gitFirst(func(repo vcs.RepoGitReader) error {
			if repo.IsRebasing() {
				return ErrIsRebasing
			}
			return nil
		})
	}

	shouldRun := len(e.funs) > 0
	if !shouldRun {
		return nil
	}

	if e.writeLock {
		lock := e.locks.Get(codebaseID, viewID)
		if err := lock.Lock(); err != nil {
			return fmt.Errorf("failed to acquire write lock: %w", err)
		}
		defer func() {
			if unlockErr := lock.Unlock(); unlockErr != nil {
				err = fmt.Errorf("failed to release write lock: %w", unlockErr)
			}
		}()
	} else if e.readLock {
		lock := e.locks.Get(codebaseID, viewID)
		if err := lock.RLock(); err != nil {
			return fmt.Errorf("failed to acquire read lock: %w", err)
		}
		defer func() {
			if unlockErr := lock.RUnlock(); unlockErr != nil {
				err = fmt.Errorf("failed to release read lock: %w", unlockErr)
			}
		}()
	}

	if e.inMemoryWriteLock {
		lock := e.locks.GetInMemory(codebaseID, viewID)
		if err := lock.Lock(); err != nil {
			return fmt.Errorf("failed to acquire in-memory write lock: %w", err)
		}
		defer func() {
			if unlockErr := lock.Unlock(); unlockErr != nil {
				err = fmt.Errorf("failed to release in-memory write lock: %w", unlockErr)
			}
		}()
	} else if e.inMemoryReadLock {
		lock := e.locks.GetInMemory(codebaseID, viewID)
		if err := lock.RLock(); err != nil {
			return fmt.Errorf("failed to acquire in-memory read lock: %w", err)
		}
		defer func() {
			if unlockErr := lock.RUnlock(); unlockErr != nil {
				err = fmt.Errorf("failed to release in-memory read lock: %w", unlockErr)
			}
		}()
	}

	defer getMeterFunc(actionName)()

	execT0 = time.Now()

	onceRepo := openOnce(e.repoProvider, codebaseID, viewID)
	for _, fn := range e.funs {
		if err := fn.Exec(onceRepo); err != nil {
			return err
		}
	}

	return nil
}

type onceRepo struct {
	codebaseID   codebases.ID
	viewID       *string
	repoProvider provider.RepoProvider

	repo vcs.RepoWriter
	err  error

	once *sync.Once
}

func openOnce(repoProvider provider.RepoProvider, codebaseID codebases.ID, viewID *string) *onceRepo {
	return &onceRepo{
		codebaseID:   codebaseID,
		viewID:       viewID,
		repoProvider: repoProvider,
		once:         &sync.Once{},
	}
}

func (or *onceRepo) Get() (vcs.RepoWriter, error) {
	or.once.Do(func() {
		if or.viewID != nil {
			or.repo, or.err = or.repoProvider.ViewRepo(or.codebaseID, *or.viewID)
		} else {
			or.repo, or.err = or.repoProvider.TrunkRepo(or.codebaseID)
		}
	})
	return or.repo, or.err
}

var (
	meteredMethod = promauto.NewHistogramVec(
		prometheus.HistogramOpts{Name: "sturdy_executor_call_millis",
			Buckets: []float64{
				.01, .025, .05,
				.1, .25, .5,
				1, 2.5, 5,
				10, 25, 50,
				100, 250, 500,
				1000, 2500, 5000,
				10000, 25000, 50000,
				100000, 250000, 500000,
			},
		}, []string{"action"})
)

func getMeterFunc(action string) func() {
	t0 := time.Now()
	return func() {
		meteredMethod.With(prometheus.Labels{"action": action}).Observe(float64(time.Since(t0).Milliseconds()))
	}
}
