package decorate

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"getsturdy.com/api/pkg/author"
	"getsturdy.com/api/pkg/change"
	db_change "getsturdy.com/api/pkg/change/db"
	service_codebase "getsturdy.com/api/pkg/codebase/service"
	vcs_codebase "getsturdy.com/api/pkg/codebase/vcs"
	"getsturdy.com/api/pkg/jsontime"
	service_user "getsturdy.com/api/pkg/users/service"
	"getsturdy.com/api/vcs"
	"getsturdy.com/api/vcs/executor"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Decorator struct {
	changeRepo      db_change.Repository
	userService     service_user.Service
	codebaseService *service_codebase.Service

	executorProvider executor.Provider

	logger *zap.Logger
}

func New(
	changeRepo db_change.Repository,
	userService service_user.Service,
	codebaseService *service_codebase.Service,

	executorProvider executor.Provider,

	logger *zap.Logger,
) *Decorator {
	return &Decorator{
		changeRepo:      changeRepo,
		userService:     userService,
		codebaseService: codebaseService,

		executorProvider: executorProvider,

		logger: logger.Named("changeDecorator"),
	}
}

type DecoratedChange struct {
	ChangeID change.ID `json:"change_id"`

	CreatedAt jsontime.Time `json:"created_at"`
	CommitID  string        `json:"commit_id"`
	IsLanded  bool          `json:"is_landed"`

	Title       string `json:"title"`
	Description string `json:"description"` // can be HTML

	Meta   change.ChangeMetadata `json:"meta"` // Data from the commit message
	Author author.Author         `json:"author"`
}

func (d *Decorator) DecorateChanges(ctx context.Context, changes []*vcs.LogEntry, codebaseID string) ([]DecoratedChange, error) {
	res := make([]DecoratedChange, len(changes))
	var err error
	for k, entry := range changes {
		res[k], err = d.DecorateChange(ctx, entry, codebaseID)
		if err != nil {
			return nil, err
		}
	}
	return res, nil
}

func (d *Decorator) DecorateChange(ctx context.Context, entry *vcs.LogEntry, codebaseID string) (DecoratedChange, error) {
	meta := change.ParseCommitMessage(entry.RawCommitMessage)

	var ch *change.Change
	var err error

	ch, err = d.changeRepo.GetByCommitID(context.Background(), entry.ID, codebaseID)
	if errors.Is(err, sql.ErrNoRows) {
		// Create a change
		ch = &change.Change{
			ID:              change.ID(uuid.New().String()),
			CodebaseID:      codebaseID,
			GitCreatedAt:    &entry.Time,
			GitCreatorName:  &entry.Name,
			GitCreatorEmail: &entry.Email,
			CommitID:        &entry.ID,
		}

		// TODO: Remove this! This import can not be trusted, and the association between commits and users should only be done on GitHub push events
		if meta.UserID != "" {
			if ok, err := d.codebaseService.CanAccess(ctx, meta.UserID, codebaseID); err == nil && ok {
				ch.UserID = &meta.UserID
			}
		}

		err = d.changeRepo.Insert(*ch)
		if err != nil {
			return DecoratedChange{}, fmt.Errorf("failed to create change: %w", err)
		}
	} else if err != nil {
		return DecoratedChange{}, fmt.Errorf("failed to get change from commit: %w", err)
	}

	// If git metadata is missing (changes imported before 2021-09-21), update the entry
	var updatedChange bool
	if ch.GitCreatorEmail == nil || ch.GitCreatedAt == nil || ch.GitCreatorName == nil {
		ch.GitCreatorEmail = &entry.Email
		ch.GitCreatorName = &entry.Name
		ch.GitCreatedAt = &entry.Time
		updatedChange = true
	}

	if ch.Title == nil || len(ch.UpdatedDescription) == 0 {
		// Import title and description from the commit
		ch.UpdatedDescription = meta.Description
		title := firstLine(meta.Description)
		ch.Title = &title
		updatedChange = true
	}

	// TODO: Remove this! This import can not be trusted, and the association between commits and users should only be done on GitHub push events
	if ch.UserID == nil && len(meta.UserID) > 0 {
		if ok, err := d.codebaseService.CanAccess(ctx, meta.UserID, codebaseID); err == nil && ok {
			ch.UserID = &meta.UserID
			updatedChange = true
		}
	}

	// Save changes to the db
	if updatedChange {
		if err := d.changeRepo.Update(*ch); err != nil {
			return DecoratedChange{}, fmt.Errorf("failed to update change: %w", err)
		}
	}

	res := DecoratedChange{
		ChangeID: ch.ID,

		CreatedAt: jsontime.Time(entry.Time),
		CommitID:  entry.ID,
		IsLanded:  entry.IsLanded,

		Title:       firstString(ch.Title),
		Description: ch.UpdatedDescription,

		Meta: meta,
	}

	if len(meta.UserID) > 0 {
		authorObj, err := d.userService.GetAsAuthor(ctx, meta.UserID)
		if err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				d.logger.Error("failed to decorate change", zap.Error(err))
			}
			res.Author = author.Author{
				Name:           entry.Name,
				IsExternalUser: true,
			}
		} else {
			res.Author = *authorObj
		}
	} else {
		res.Author = author.Author{
			Name:           entry.Name,
			IsExternalUser: true,
		}
	}

	return res, nil
}

func firstString(in ...*string) string {
	for _, v := range in {
		if v != nil && len(*v) > 0 {
			return *v
		}
	}
	return ""
}

func firstLine(in string) string {
	idx := strings.IndexByte(in, '\n')
	if idx < 0 {
		return in
	}
	return in[0:idx]
}

func (d *Decorator) List(ctx context.Context, codebaseID string, limit int) ([]DecoratedChange, error) {
	// vcs.ListChanges and decorate.DecorateChanges will import all commits to Sturdy.
	// This is not ideal. If we could make sure that the database is already is up to date with the Git state,
	// we would not have to read from disk here.
	var gitChangeLog []*vcs.LogEntry
	if err := d.executorProvider.New().GitRead(func(repo vcs.RepoGitReader) error {
		var err error
		gitChangeLog, err = vcs_codebase.ListChanges(repo, limit)
		if err != nil {
			return fmt.Errorf("failed to list changes: %w", err)
		}
		return nil
	}).ExecTrunk(codebaseID, "codebase.Changes"); err != nil {
		return nil, err
	}

	decoratedLog, err := d.DecorateChanges(ctx, gitChangeLog, codebaseID)
	if err != nil {
		return nil, err
	}

	return decoratedLog, nil
}
