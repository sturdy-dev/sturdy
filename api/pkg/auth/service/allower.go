package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"getsturdy.com/api/pkg/auth"
	"getsturdy.com/api/pkg/change"
	"getsturdy.com/api/pkg/codebase"
	"getsturdy.com/api/pkg/codebase/acl"
	"getsturdy.com/api/pkg/suggestions"
	"getsturdy.com/api/pkg/unidiff"
	"getsturdy.com/api/pkg/workspace"
)

var (
	allAllowed, _  = unidiff.NewAllower("*")
	noneAllowed, _ = unidiff.NewAllower()
)

func (s *Service) GetAllower(ctx context.Context, obj interface{}) (*unidiff.Allower, error) {
	if obj == nil {
		return noneAllowed, nil
	}

	subject, found := auth.FromContext(ctx)
	if !found {
		return noneAllowed, nil
	}

	switch subject.Type {
	case auth.SubjectMutagen:
		// TODO: mutagen request should be authenticated
		switch object := obj.(type) {
		case *codebase.Codebase:
			return s.getUserCodebaseAllower(ctx, subject.ID, object)
		case codebase.Codebase:
			return s.getUserCodebaseAllower(ctx, subject.ID, &object)
		}

	case auth.SubjectUser:
		switch object := obj.(type) {
		case *codebase.Codebase:
			return s.getUserCodebaseAllower(ctx, subject.ID, object)
		case codebase.Codebase:
			return s.getUserCodebaseAllower(ctx, subject.ID, &object)
		case change.ChangeCommit:
			return s.getUserChangeCommitAllower(ctx, subject.ID, &object)
		case *change.ChangeCommit:
			return s.getUserChangeCommitAllower(ctx, subject.ID, object)
		case change.Change:
			return s.getUserChangeAllower(ctx, subject.ID, &object)
		case *change.Change:
			return s.getUserChangeAllower(ctx, subject.ID, object)
		case workspace.Workspace:
			return s.getUserWorkspaceAllower(ctx, subject.ID, &object)
		case *workspace.Workspace:
			return s.getUserWorkspaceAllower(ctx, subject.ID, object)
		case suggestions.Suggestion:
			return s.getUserSuggestionAllower(ctx, subject.ID, &object)
		case *suggestions.Suggestion:
			return s.getUserSuggestionAllower(ctx, subject.ID, object)
		}

	case auth.SubjectCI:
		switch object := obj.(type) {
		case *change.ChangeCommit:
			return s.getCIChangeCommit(ctx, subject.ID, object)
		case change.ChangeCommit:
			return s.getCIChangeCommit(ctx, subject.ID, &object)
		case *change.Change:
			return s.getCIChangeAllower(ctx, subject.ID, object)
		case change.Change:
			return s.getCIChangeAllower(ctx, subject.ID, &object)
		}

	case auth.SubjectAnonymous:
		switch object := obj.(type) {
		case *change.ChangeCommit:
			return s.getAnonymousChangeCommitAllower(ctx, object)
		case change.ChangeCommit:
			return s.getAnonymousChangeCommitAllower(ctx, &object)
		case *change.Change:
			return s.getAnonymousChangeAllower(ctx, object)
		case change.Change:
			return s.getAnonymousChangeAllower(ctx, &object)
		case workspace.Workspace:
			return s.getAnonymousWorkspaceAllower(ctx, &object)
		case *workspace.Workspace:
			return s.getAnonymousWorkspaceAllower(ctx, object)
		case *codebase.Codebase:
			return s.getAnonymousCodebaseAllower(ctx, object)
		case codebase.Codebase:
			return s.getAnonymousCodebaseAllower(ctx, &object)
		}
	}

	return noneAllowed, nil
}

func (s *Service) getUserChangeCommitAllower(ctx context.Context, userID string, changeCommit *change.ChangeCommit) (*unidiff.Allower, error) {
	cb, err := s.codebaseService.GetByID(ctx, changeCommit.CodebaseID)
	if err != nil {
		return nil, fmt.Errorf("failed to get codebase: %w", err)
	}
	return s.getUserCodebaseAllower(ctx, userID, cb)
}

func (s *Service) getUserChangeAllower(ctx context.Context, userID string, change *change.Change) (*unidiff.Allower, error) {
	cb, err := s.codebaseService.GetByID(ctx, change.CodebaseID)
	if err != nil {
		return nil, fmt.Errorf("failed to get codebase: %w", err)
	}
	return s.getUserCodebaseAllower(ctx, userID, cb)
}

func (s *Service) getUserWorkspaceAllower(ctx context.Context, userID string, workspace *workspace.Workspace) (*unidiff.Allower, error) {
	cb, err := s.codebaseService.GetByID(ctx, workspace.CodebaseID)
	if err != nil {
		return nil, fmt.Errorf("failed to get codebase: %w", err)
	}
	return s.getUserCodebaseAllower(ctx, userID, cb)
}

func (s *Service) getUserSuggestionAllower(ctx context.Context, userID string, suggestion *suggestions.Suggestion) (*unidiff.Allower, error) {
	cb, err := s.codebaseService.GetByID(ctx, suggestion.CodebaseID)
	if err != nil {
		return nil, fmt.Errorf("failed to get codebase: %w", err)
	}
	return s.getUserCodebaseAllower(ctx, userID, cb)
}

func (s *Service) getUserCodebaseAllower(ctx context.Context, userID string, codebase *codebase.Codebase) (*unidiff.Allower, error) {
	aclPolicy, err := s.aclProvider.GetByCodebaseID(ctx, codebase.ID)
	if errors.Is(err, sql.ErrNoRows) {
		return noneAllowed, nil
	} else if err != nil {
		return nil, fmt.Errorf("failed to get acl policy: %w", err)
	}

	user, err := s.userService.GetByID(ctx, userID)
	if errors.Is(err, sql.ErrNoRows) {
		return noneAllowed, nil
	} else if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	allowedByEmail := aclPolicy.Policy.List(
		acl.Identity{Type: acl.Users, ID: user.Email},
		acl.ActionWrite,
		acl.Files,
	)

	allowedByID := aclPolicy.Policy.List(
		acl.Identity{Type: acl.Users, ID: user.ID},
		acl.ActionWrite,
		acl.Files,
	)

	return unidiff.NewAllower(append(allowedByEmail, allowedByID...)...)
}

func (s *Service) getCIChangeCommit(ctx context.Context, changeID string, changeCommit *change.ChangeCommit) (*unidiff.Allower, error) {
	if changeID != string(changeCommit.ChangeID) {
		return noneAllowed, nil
	}
	return allAllowed, nil
}

func (s *Service) getCIChangeAllower(ctx context.Context, changeID string, change *change.Change) (*unidiff.Allower, error) {
	if changeID != string(change.ID) {
		return noneAllowed, nil
	}
	return allAllowed, nil
}

func (s *Service) getAnonymousWorkspaceAllower(ctx context.Context, workspace *workspace.Workspace) (*unidiff.Allower, error) {
	cb, err := s.codebaseService.GetByID(ctx, workspace.CodebaseID)
	if err != nil {
		return nil, fmt.Errorf("failed to get codebase: %w", err)
	}
	return s.getAnonymousCodebaseAllower(ctx, cb)
}

func (s *Service) getAnonymousChangeCommitAllower(ctx context.Context, changeCommit *change.ChangeCommit) (*unidiff.Allower, error) {
	cb, err := s.codebaseService.GetByID(ctx, changeCommit.CodebaseID)
	if err != nil {
		return nil, fmt.Errorf("failed to get codebase: %w", err)
	}
	return s.getAnonymousCodebaseAllower(ctx, cb)
}

func (s *Service) getAnonymousChangeAllower(ctx context.Context, change *change.Change) (*unidiff.Allower, error) {
	cb, err := s.codebaseService.GetByID(ctx, change.CodebaseID)
	if err != nil {
		return nil, fmt.Errorf("failed to get codebase: %w", err)
	}
	return s.getAnonymousCodebaseAllower(ctx, cb)
}

func (s *Service) getAnonymousCodebaseAllower(ctx context.Context, cb *codebase.Codebase) (*unidiff.Allower, error) {
	if cb.IsPublic {
		return allAllowed, nil
	}

	return noneAllowed, nil
}
