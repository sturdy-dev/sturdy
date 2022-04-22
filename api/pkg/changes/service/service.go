package service

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"sync"
	"time"

	git "github.com/libgit2/git2go/v33"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"getsturdy.com/api/pkg/changes"
	db_change "getsturdy.com/api/pkg/changes/db"
	"getsturdy.com/api/pkg/changes/message"
	"getsturdy.com/api/pkg/codebases"
	db_codebases "getsturdy.com/api/pkg/codebases/db"
	"getsturdy.com/api/pkg/snapshots/snapshotter"
	"getsturdy.com/api/pkg/unidiff"
	"getsturdy.com/api/pkg/workspaces"
	"getsturdy.com/api/vcs"
	"getsturdy.com/api/vcs/executor"

	"github.com/google/uuid"
)

type Service struct {
	changeRepo       db_change.Repository
	codebaseRepo     db_codebases.CodebaseRepository
	logger           *zap.Logger
	executorProvider executor.Provider
	snap             snapshotter.Snapshotter

	// TODO: remove once we've added unique index on changes(codebase_id, commit_id)
	createWithChangeAsParentMutex sync.Mutex
}

func New(
	changeRepo db_change.Repository,
	codebaseRepo db_codebases.CodebaseRepository,
	logger *zap.Logger,
	executorProvider executor.Provider,
	snap snapshotter.Snapshotter,
) *Service {
	return &Service{
		changeRepo:       changeRepo,
		codebaseRepo:     codebaseRepo,
		logger:           logger.Named("changeService"),
		executorProvider: executorProvider,
		snap:             snap,
	}
}

func (svc *Service) ListChanges(ctx context.Context, ids ...changes.ID) ([]*changes.Change, error) {
	return svc.changeRepo.ListByIDs(ctx, ids...)
}

func (svc *Service) GetChangeByID(ctx context.Context, id changes.ID) (*changes.Change, error) {
	ch, err := svc.changeRepo.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return ch, nil
}

func (svc *Service) GetByCommitAndCodebase(ctx context.Context, commitID string, codebaseID codebases.ID) (*changes.Change, error) {
	ch, err := svc.changeRepo.GetByCommitID(ctx, commitID, codebaseID)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		// import
		return svc.importCommitToChange(ctx, codebaseID, commitID)
	case err != nil:
		return nil, err
	default:
		return ch, nil
	}
}

func (svc *Service) CreateOnTop(ctx context.Context, ws *workspaces.Workspace, commitID string) (*changes.Change, error) {
	headChange, err := svc.head(ctx, ws.CodebaseID)
	switch {
	case errors.Is(err, ErrNotFound):
		return svc.CreateWithChangeAsParent(ctx, ws, commitID, nil)
	case err != nil:
		return nil, fmt.Errorf("could not get head change: %w", err)
	default:
		return svc.CreateWithChangeAsParent(ctx, ws, commitID, &headChange.ID)
	}
}

func (svc *Service) CreateWithCommitAsParent(ctx context.Context, ws *workspaces.Workspace, commitID, parentCommitID string) (*changes.Change, error) {
	var parentChangeID *changes.ID

	parent, err := svc.getChangeFromCommit(ctx, ws.CodebaseID, parentCommitID)
	switch {
	case err == nil:
		parentChangeID = &parent.ID
	case errors.Is(err, ErrNotFound):
		// nothing
	default:
		return nil, fmt.Errorf("failed to get change from parent commit: %w", err)
	}

	return svc.CreateWithChangeAsParent(ctx, ws, commitID, parentChangeID)
}

var ErrAlreadyExists = fmt.Errorf("change already exists")

func (svc *Service) CreateWithChangeAsParent(ctx context.Context, ws *workspaces.Workspace, commitID string, parentChangeID *changes.ID) (*changes.Change, error) {
	svc.createWithChangeAsParentMutex.Lock()
	defer svc.createWithChangeAsParentMutex.Unlock()

	if _, err := svc.changeRepo.GetByCommitID(ctx, commitID, ws.CodebaseID); errors.Is(err, sql.ErrNoRows) {
		// ok
	} else if err != nil {
		return nil, fmt.Errorf("could not get change by commit id: %w", err)
	} else {
		return nil, fmt.Errorf("change with commit %s already exists: %w", commitID, ErrAlreadyExists)
	}

	changeID := changes.ID(uuid.NewString())
	t := time.Now()

	cleanCommitMessage := message.CommitMessage(ws.DraftDescription)
	title := message.Title(cleanCommitMessage)

	changeChange := changes.Change{
		ID:                 changeID,
		CodebaseID:         ws.CodebaseID,
		Title:              &title,
		UpdatedDescription: ws.DraftDescription,
		UserID:             &ws.UserID,
		CreatedAt:          &t,
		CommitID:           &commitID,
		ParentChangeID:     parentChangeID,
		WorkspaceID:        &ws.ID,
	}

	if err := svc.changeRepo.Insert(ctx, changeChange); err != nil {
		return nil, fmt.Errorf("failed to insert change: %w", err)
	}

	return &changeChange, nil
}

func (svc *Service) head(ctx context.Context, codebaseID codebases.ID) (*changes.Change, error) {
	// To find the root commit, peek into git
	var headCommitID string

	getHeadCommit := func(repo vcs.RepoGitReader) error {
		headCommit, err := repo.HeadCommit()
		if err != nil {
			return fmt.Errorf("could not find head commit: %w", err)
		}
		headCommitID = headCommit.Id().String()
		return nil
	}

	err := svc.executorProvider.New().GitRead(getHeadCommit).ExecTrunk(codebaseID, "changeServiceChangelog")
	switch {
	case errors.Is(err, vcs.ErrNotFound):
		return nil, ErrNotFound
	case err != nil:
		return nil, fmt.Errorf("could not get head commit: %w", err)
	default:
		return svc.getChangeFromCommit(ctx, codebaseID, headCommitID)
	}
}

// Changelog returns a list of changes for the given codebaesID in the descending order.
//
// limit - the maximum number of changes to return
// before - if set, used as a change id to start the list from
//          if not set, list will start from the head
func (svc *Service) Changelog(ctx context.Context, codebaseID codebases.ID, limit int, before *changes.ID) ([]*changes.Change, error) {
	var (
		startFrom *changes.Change
		err       error
		res       []*changes.Change
	)

	if before == nil {
		startFrom, err = svc.head(ctx, codebaseID)
		res = append(res, startFrom)
	} else {
		startFrom, err = svc.changeRepo.Get(ctx, *before)
	}

	switch {
	case errors.Is(err, ErrNotFound):
		return nil, nil
	case err != nil:
		return nil, fmt.Errorf("could not get head change: %w", err)
	}

	nextChange := startFrom
	for len(res) < limit {
		next, err := svc.ParentChange(ctx, nextChange)
		switch {
		case errors.Is(err, ErrNotFound):
			return res, nil
		case err != nil:
			return nil, err
		case err == nil:
			res = append(res, next)
			nextChange = next
		}
	}

	return res, nil
}

var ErrNotFound = errors.New("not found")

// ChildChange return the first child change of the given change.
func (svc *Service) ChildChange(ctx context.Context, ch *changes.Change) (*changes.Change, error) {
	if child, err := svc.changeRepo.GetByParentChangeID(ctx, ch.ID); errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, fmt.Errorf("could not get child change: %w", err)
	} else {
		return child, nil
	}
}

// ParentChange returns the parent change of the given change.
func (svc *Service) ParentChange(ctx context.Context, ch *changes.Change) (*changes.Change, error) {
	if ch.ParentChangeID != nil {
		// get from db
		next, err := svc.changeRepo.Get(ctx, *ch.ParentChangeID)
		switch {
		case err == nil:
			return next, nil
		case errors.Is(err, sql.ErrNoRows):
			// parent not found in db? re-create and import from git
		case err != nil:
			return nil, fmt.Errorf("could not get parent change from db id=%s: %w", *ch.ParentChangeID, err)
		}
	}

	// get parents from git
	var parents []string
	getCurrentFromGit := func(repo vcs.RepoGitReader) error {
		details, err := repo.GetCommitDetails(*ch.CommitID)
		if err != nil {
			return fmt.Errorf("could not get commit details from repo: %w", err)
		}
		parents = details.Parents
		return nil
	}
	if err := svc.executorProvider.New().GitRead(getCurrentFromGit).ExecTrunk(ch.CodebaseID, "changeService.parentChange"); err != nil {
		return nil, fmt.Errorf("could not get from git: %w", err)
	}

	// this commit is a root commit
	if len(parents) == 0 {
		return nil, ErrNotFound
	}

	// the first parent (usually) refers to the state of the branch that that the branch was merged _into_, prior to the merge.
	parent, err := svc.importCommitToChange(ctx, ch.CodebaseID, parents[0])
	if err != nil {
		return nil, fmt.Errorf("failed to import parent: %w", err)
	}

	// update the "current" commit and mark the new commit as it's parent
	ch.ParentChangeID = &parent.ID
	err = svc.changeRepo.Update(ctx, *ch)
	if err != nil {
		return nil, fmt.Errorf("failed to update change parent: %w", err)
	}

	return parent, nil
}

func (svc *Service) getChangeFromCommit(ctx context.Context, codebaseID codebases.ID, commitID string) (*changes.Change, error) {
	ch, err := svc.changeRepo.GetByCommitID(ctx, commitID, codebaseID)
	switch {
	case err == nil:
		return ch, nil
	case errors.Is(err, sql.ErrNoRows):
		return svc.importCommitToChange(ctx, codebaseID, commitID)
	default:
		return nil, fmt.Errorf("failed to get change from db: %w", err)
	}
}

func (svc *Service) importCommitToChange(ctx context.Context, codebaseID codebases.ID, commitID string) (*changes.Change, error) {
	// if the change exists in the db, use it!
	{
		fromDb, err := svc.changeRepo.GetByCommitID(ctx, commitID, codebaseID)
		switch {
		case err == nil:
			return fromDb, nil
		case errors.Is(err, sql.ErrNoRows):
		case err != nil:
			return nil, fmt.Errorf("could not lookup change by commit: %w", err)
		}
	}

	var details *vcs.CommitDetails
	var err error

	getCommit := func(repo vcs.RepoGitReader) error {
		details, err = repo.GetCommitDetails(commitID)
		if err != nil {
			return fmt.Errorf("could not get commit details: %w", err)
		}
		return nil
	}
	if err := svc.executorProvider.New().GitRead(getCommit).ExecTrunk(codebaseID, "changeServiceChangelog"); err != nil {
		return nil, err
	}

	// don't import Sturdy-style root commits
	if len(details.Parents) == 0 && details.Message == "Root Commit" {
		return nil, ErrNotFound
	}

	description := strings.TrimSpace(details.Message)
	title := message.Title(description)
	description = strings.NewReplacer("\r\n", "<br>", "\n", "<br>").Replace(description)

	// CreateWithCommitAsParent change!
	ch := changes.Change{
		ID:                 changes.ID(uuid.NewString()),
		CodebaseID:         codebaseID,
		Title:              &title,
		UpdatedDescription: description,
		UserID:             nil,
		CreatedAt:          nil, // Set?
		GitCreatedAt:       &details.Author.When,
		GitCreatorName:     &details.Author.Name,
		GitCreatorEmail:    &details.Author.Email,
		CommitID:           &commitID,
		ParentChangeID:     nil, // Parent is starts out as nil. If/when the parent commit is imported, this value will be set.
		WorkspaceID:        nil, // WorkspaceID is nil, since this is an imported change.
	}

	if err := svc.changeRepo.Insert(ctx, ch); err != nil {
		return nil, fmt.Errorf("could not write new change to db: %w", err)
	}

	return &ch, nil
}

func (svc *Service) Diffs(ctx context.Context, ch *changes.Change, allower *unidiff.Allower) ([]unidiff.FileDiff, error) {
	parent, err := svc.ParentChange(ctx, ch)
	switch {
	case errors.Is(err, ErrNotFound):
		// use diffToRoot
	case err != nil:
		return nil, fmt.Errorf("could not get change parent: %w", err)
	}

	var diff *git.Diff
	diffBetweenCommits := func(repo vcs.RepoGitReader) error {
		diff, err = repo.DiffCommits(*parent.CommitID, *ch.CommitID)
		if err != nil {
			return fmt.Errorf("could not get diffs: %w", err)
		}
		return nil
	}

	diffToRoot := func(repo vcs.RepoGitReader) error {
		diff, err = repo.DiffCommitToRoot(*ch.CommitID)
		if err != nil {
			return fmt.Errorf("could not get diff to root: %w", err)
		}
		return nil
	}

	var fn func(repo vcs.RepoGitReader) error
	if parent != nil {
		fn = diffBetweenCommits
	} else {
		fn = diffToRoot
	}

	err = svc.executorProvider.New().GitRead(fn).ExecTrunk(ch.CodebaseID, "changeService.Diffs")
	if err != nil {
		return nil, err
	}

	decoratedDiff, err := unidiff.NewUnidiff(
		unidiff.NewGitPatchReader(diff),
		svc.logger,
	).WithAllower(allower).Decorate()
	if err != nil {
		return nil, fmt.Errorf("failed to generate unidiff for diff: %w", err)
	}

	return decoratedDiff, nil
}

func (svc *Service) HeadChange(ctx context.Context, cb *codebases.Codebase) (*changes.Change, error) {
	if cb.CalculatedHeadChangeID && cb.CachedHeadChangeID == nil {
		return nil, ErrNotFound
	}

	if cb.CalculatedHeadChangeID && cb.CachedHeadChangeID != nil {
		return svc.GetChangeByID(ctx, changes.ID(*cb.CachedHeadChangeID))
	}

	headChange, err := svc.head(ctx, cb.ID)
	switch {
	case errors.Is(err, ErrNotFound):
		cb.CalculatedHeadChangeID = true
		cb.CachedHeadChangeID = nil
	case err != nil:
		return nil, fmt.Errorf("failed to get changelog: %w", err)
	default:
		cb.CalculatedHeadChangeID = true
		cb.CachedHeadChangeID = (*string)(&headChange.ID)
	}

	if err := svc.codebaseRepo.Update(cb); err != nil {
		return nil, fmt.Errorf("failed to update codebase head change: %w", err)
	}

	if headChange == nil {
		return nil, ErrNotFound
	}

	return headChange, nil
}

func (svc *Service) UnsetHeadChangeCache(codebaseID codebases.ID) error {
	cb, err := svc.codebaseRepo.Get(codebaseID)
	if err != nil {
		return fmt.Errorf("failed to get codebase: %w", err)
	}

	cb.CalculatedHeadChangeID = false
	cb.CachedHeadChangeID = nil

	if err := svc.codebaseRepo.Update(cb); err != nil {
		return fmt.Errorf("failed to set codebase head change: %w", err)
	}

	return nil
}

func (svc *Service) SetAsHeadChange(ch *changes.Change) error {
	cb, err := svc.codebaseRepo.Get(ch.CodebaseID)
	if err != nil {
		return fmt.Errorf("failed to get codebase: %w", err)
	}

	cb.CalculatedHeadChangeID = true
	cb.CachedHeadChangeID = (*string)(&ch.ID)

	if err := svc.codebaseRepo.Update(cb); err != nil {
		return fmt.Errorf("failed to set codebase head change: %w", err)
	}

	return nil
}
