package executor

import (
	"fmt"
	"runtime/debug"
	"sync"
	"time"

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
	Read(ReadFunc) Executor
	// Write schedules a function that can write files to the filesystem.
	Write(WriteFunc) Executor
	// Git schedules a function that can access git repository, but not the files on the filesystem.
	Git(GitFunc) Executor

	// Schedule schedules a function that can access repo provider. It is deprecated, prefer to use
	// Read, Write or Git functions instead to get more granular locking.
	Schedule(fn ExecuteFunc) Executor

	// AllowRebasingState allows the executor to open a repo in a rebasing state. By default this is not allowed.
	AllowRebasingState() Executor
	// AssertBranchName asserts that the repo is in the expected branch.
	AssertBranchName(string) Executor

	// ExecView executes all of the scheduled functions for the given view repository.
	ExecView(codebaseID, viewID, actionName string) error
	// ExecTrunk executes all of the scheduled functions for the given trunk repository.
	ExecTrunk(codebaseID, actionName string) error
}

type executeFunc struct {
	fun ExecuteFunc

	repoFun GitFunc

	writeFun WriteFunc

	readFun ReadFunc
}

func (e *executeFunc) Exec(repo *onceRepo) error {
	if e.fun != nil {
		return e.fun(repo.repoProvider)
	}

	r, err := repo.Get()
	if err != nil {
		return err
	}

	if e.repoFun != nil {
		return e.repoFun(r)
	}

	if e.writeFun != nil {
		return e.writeFun(r)
	}

	if e.readFun != nil {
		return e.readFun(r)
	}

	return nil
}

func (e *executeFunc) Write() bool {
	return e.fun != nil && e.writeFun != nil
}

func (e *executeFunc) Read() bool {
	return !e.Write() && e.readFun != nil
}

type executor struct {
	funs []*executeFunc

	writeLock bool
	readLock  bool

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

type WriteFunc func(vcs.RepoWriter) error

func (e *executor) Write(fn WriteFunc) Executor {
	e.writeLock = true
	e.funs = append(e.funs, &executeFunc{
		writeFun: fn,
	})
	return e
}

type ReadFunc func(vcs.RepoReader) error

func (e *executor) Read(fn ReadFunc) Executor {
	e.readLock = true
	e.funs = append(e.funs, &executeFunc{
		readFun: fn,
	})
	return e
}

type GitFunc func(repo vcs.Repo) error

func (e *executor) Git(fn GitFunc) Executor {
	e.funs = append(e.funs, &executeFunc{
		repoFun: fn,
	})
	return e
}

func (e *executor) gitFirst(fn GitFunc) Executor {
	e.funs = append([]*executeFunc{{repoFun: fn}}, e.funs...)
	return e
}

func (e *executor) AllowRebasingState() Executor {
	e.allowRebasing = true
	return e
}

func (e *executor) AssertBranchName(name string) Executor {
	return e.Git(func(repo vcs.Repo) error {
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

func (e *executor) ExecView(codebaseID, viewID, actionName string) (err error) {
	return e.exec(codebaseID, &viewID, actionName)
}

func (e *executor) ExecTrunk(codebaseID, actionName string) (err error) {
	return e.exec(codebaseID, nil, actionName)
}

func (e *executor) exec(codebaseID string, viewID *string, actionName string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("recovered call in vcs executor: %v\nStacktrace: %s", r, string(debug.Stack()))
		}
	}()

	logger := e.logger.With(zap.String("codebase_id", codebaseID), zap.String("action_name", actionName))
	if viewID != nil {
		logger = logger.With(zap.String("view_id", *viewID))
	}

	t0 := time.Now()
	defer logger.Info("git executor completed", zap.Duration("duration", time.Since(t0)))

	if !e.allowRebasing {
		e.gitFirst(func(repo vcs.Repo) error {
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

	lock := e.locks.Get(codebaseID, viewID)
	if e.writeLock {
		if err := lock.Lock(); err != nil {
			return fmt.Errorf("failed to acquire write lock: %w", err)
		}
		defer func() {
			if unlockErr := lock.Unlock(); unlockErr != nil {
				err = fmt.Errorf("failed to release write lock: %w", unlockErr)
			}
		}()
	} else if e.readLock {
		if err := lock.RLock(); err != nil {
			return fmt.Errorf("failed to acquire read lock: %w", err)
		}
		defer func() {
			if unlockErr := lock.RUnlock(); unlockErr != nil {
				err = fmt.Errorf("failed to release read lock: %w", unlockErr)
			}
		}()
	}

	onceRepo := openOnce(e.repoProvider, codebaseID, viewID)

	for _, fn := range e.funs {
		if err := fn.Exec(onceRepo); err != nil {
			return err
		}
	}

	return nil
}

type onceRepo struct {
	codebaseID   string
	viewID       *string
	repoProvider provider.RepoProvider

	repo vcs.RepoWriter
	err  error

	once *sync.Once
}

func openOnce(repoProvider provider.RepoProvider, codebaseID string, viewID *string) *onceRepo {
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
