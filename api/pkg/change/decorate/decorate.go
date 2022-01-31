package decorate

import (
	"database/sql"
	"errors"
	"fmt"
	"getsturdy.com/api/pkg/author"
	"getsturdy.com/api/pkg/change"
	db_change "getsturdy.com/api/pkg/change/db"
	db_codebase "getsturdy.com/api/pkg/codebase/db"
	"getsturdy.com/api/pkg/jsontime"
	db_user "getsturdy.com/api/pkg/users/db"
	"getsturdy.com/api/vcs"
	"strings"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type DecoratedChange struct {
	ChangeID change.ID `json:"change_id"`

	CreatedAt jsontime.Time `json:"created_at"`
	CommitID  string        `json:"commit_id"`
	IsLanded  bool          `json:"is_landed"`

	Title       string `json:"title"`
	Description string `json:"description"` // can be HTML

	Meta   ChangeMetadata `json:"meta"` // Data from the commit message
	Author author.Author  `json:"author"`
}

func DecorateChanges(changes []*vcs.LogEntry, userRepo db_user.Repository, logger *zap.Logger, changeRepo db_change.Repository, changeCommitRepo db_change.CommitRepository, codebaseUserRepo db_codebase.CodebaseUserRepository, codebaseID string, isTrunk bool) ([]DecoratedChange, error) {
	res := make([]DecoratedChange, len(changes))
	var err error
	for k, entry := range changes {
		res[k], err = DecorateChange(entry, userRepo, logger, changeRepo, changeCommitRepo, codebaseUserRepo, codebaseID, isTrunk)
		if err != nil {
			return nil, err
		}
	}
	return res, nil
}

func DecorateChange(entry *vcs.LogEntry, userRepo db_user.Repository, logger *zap.Logger, changeRepo db_change.Repository, changeCommitRepo db_change.CommitRepository, codebaseUserRepo db_codebase.CodebaseUserRepository, codebaseID string, isTrunk bool) (DecoratedChange, error) {
	meta := ParseCommitMessage(entry.RawCommitMessage)

	var ch change.Change
	var chCommit change.ChangeCommit
	var err error

	chCommit, err = changeCommitRepo.GetByCommitID(entry.ID, codebaseID)
	if errors.Is(err, sql.ErrNoRows) {
		// Create both a change and a change
		ch = change.Change{
			ID:              change.ID(uuid.New().String()),
			CodebaseID:      codebaseID,
			GitCreatedAt:    &entry.Time,
			GitCreatorName:  &entry.Name,
			GitCreatorEmail: &entry.Email,
		}

		// TODO: Remove this! This import can not be trusted, and the association between commits and users should only be done on GitHub push events
		if meta.UserID != "" {
			if _, err := codebaseUserRepo.GetByUserAndCodebase(meta.UserID, codebaseID); err == nil {
				ch.UserID = &meta.UserID
			}
		}

		err = changeRepo.Insert(ch)
		if err != nil {
			return DecoratedChange{}, fmt.Errorf("failed to create change: %w", err)
		}

		chCommit = change.ChangeCommit{
			ChangeID:   ch.ID,
			CommitID:   entry.ID,
			CodebaseID: codebaseID,
			Trunk:      isTrunk,
		}

		err = changeCommitRepo.Insert(chCommit)
		if err != nil {
			return DecoratedChange{}, fmt.Errorf("failed to create change_commit: %w", err)
		}
	} else if err != nil {
		return DecoratedChange{}, fmt.Errorf("failed to get change from commit: %w", err)
	} else {
		// Load change
		ch, err = changeRepo.Get(chCommit.ChangeID)
		if err != nil {
			return DecoratedChange{}, fmt.Errorf("failed to get change: %w", err)
		}
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

	if !chCommit.Trunk && isTrunk {
		chCommit.Trunk = true
		if err := changeCommitRepo.Update(chCommit); err != nil {
			return DecoratedChange{}, fmt.Errorf("failed to update change commit: %w", err)
		}
	}

	// TODO: Remove this! This import can not be trusted, and the association between commits and users should only be done on GitHub push events
	if ch.UserID == nil && len(meta.UserID) > 0 {
		if _, err := codebaseUserRepo.GetByUserAndCodebase(meta.UserID, codebaseID); err == nil {
			ch.UserID = &meta.UserID
			updatedChange = true
		}
	}

	// Save changes to the db
	if updatedChange {
		if err := changeRepo.Update(ch); err != nil {
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
		authorObj, err := author.GetAuthor(meta.UserID, userRepo)
		if err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				logger.Error("failed to decorate change", zap.Error(err))
			}
			res.Author = author.Author{
				Name:           entry.Name,
				IsExternalUser: true,
			}
		} else {
			res.Author = authorObj
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
