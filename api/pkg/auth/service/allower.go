package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"getsturdy.com/api/pkg/auth"
	"getsturdy.com/api/pkg/changes"
	"getsturdy.com/api/pkg/codebase"
	"getsturdy.com/api/pkg/codebase/acl"
	"getsturdy.com/api/pkg/suggestions"
	"getsturdy.com/api/pkg/unidiff"
	"getsturdy.com/api/pkg/users"
	"getsturdy.com/api/pkg/workspaces"
)

var (
	allAllowed, _  = unidiff.NewAllower("*")
	noneAllowed, _ = unidiff.NewAllower()
)

func (s *Service) GetAllower(ctx context.Context, obj any) (*unidiff.Allower, error) {
	if obj == nil {
		return noneAllowed, nil
	}

	subject, found := auth.FromContext(ctx)
	if !found {
		return noneAllowed, nil
	}

	switch subject.Type {
	case auth.SubjectMutagen:
		subjectID := users.ID(subject.ID)
		// TODO: mutagen request should be authenticated
		switch object := obj.(type) {
		case *codebase.Codebase:
			return s.getUserCodebaseAllower(ctx, subjectID, object)
		case codebase.Codebase:
			return s.getUserCodebaseAllower(ctx, subjectID, &object)
		}

	case auth.SubjectUser:
		subjectID := users.ID(subject.ID)
		switch object := obj.(type) {
		case *codebase.Codebase:
			return s.getUserCodebaseAllower(ctx, subjectID, object)
		case codebase.Codebase:
			return s.getUserCodebaseAllower(ctx, subjectID, &object)
		case changes.Change:
			return s.getUserChangeAllower(ctx, subjectID, &object)
		case *changes.Change:
			return s.getUserChangeAllower(ctx, subjectID, object)
		case workspaces.Workspace:
			return s.getUserWorkspaceAllower(ctx, subjectID, &object)
		case *workspaces.Workspace:
			return s.getUserWorkspaceAllower(ctx, subjectID, object)
		case suggestions.Suggestion:
			return s.getUserSuggestionAllower(ctx, subjectID, &object)
		case *suggestions.Suggestion:
			return s.getUserSuggestionAllower(ctx, subjectID, object)
		}

	case auth.SubjectCI:
		switch object := obj.(type) {
		case *changes.Change:
			return s.getCIChangeAllower(ctx, subject.ID, object)
		case changes.Change:
			return s.getCIChangeAllower(ctx, subject.ID, &object)
		}

	case auth.SubjectAnonymous:
		switch object := obj.(type) {
		case *changes.Change:
			return s.getAnonymousChangeAllower(ctx, object)
		case changes.Change:
			return s.getAnonymousChangeAllower(ctx, &object)
		case workspaces.Workspace:
			return s.getAnonymousWorkspaceAllower(ctx, &object)
		case *workspaces.Workspace:
			return s.getAnonymousWorkspaceAllower(ctx, object)
		case *codebase.Codebase:
			return s.getAnonymousCodebaseAllower(ctx, object)
		case codebase.Codebase:
			return s.getAnonymousCodebaseAllower(ctx, &object)
		}
	}

	return noneAllowed, nil
}

func (s *Service) getUserChangeAllower(ctx context.Context, userID users.ID, change *changes.Change) (*unidiff.Allower, error) {
	cb, err := s.codebaseService.GetByID(ctx, change.CodebaseID)
	if err != nil {
		return nil, fmt.Errorf("failed to get codebase: %w", err)
	}
	return s.getUserCodebaseAllower(ctx, userID, cb)
}

func (s *Service) getUserWorkspaceAllower(ctx context.Context, userID users.ID, workspace *workspaces.Workspace) (*unidiff.Allower, error) {
	cb, err := s.codebaseService.GetByID(ctx, workspace.CodebaseID)
	if err != nil {
		return nil, fmt.Errorf("failed to get codebase: %w", err)
	}
	return s.getUserCodebaseAllower(ctx, userID, cb)
}

func (s *Service) getUserSuggestionAllower(ctx context.Context, userID users.ID, suggestion *suggestions.Suggestion) (*unidiff.Allower, error) {
	cb, err := s.codebaseService.GetByID(ctx, suggestion.CodebaseID)
	if err != nil {
		return nil, fmt.Errorf("failed to get codebase: %w", err)
	}
	return s.getUserCodebaseAllower(ctx, userID, cb)
}

func (s *Service) getUserCodebaseAllower(ctx context.Context, userID users.ID, codebase *codebase.Codebase) (*unidiff.Allower, error) {
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
		acl.Identity{Type: acl.Users, ID: user.ID.String()},
		acl.ActionWrite,
		acl.Files,
	)

	return unidiff.NewAllower(append(allowedByEmail, allowedByID...)...)
}

func (s *Service) getCIChangeAllower(ctx context.Context, changeID string, change *changes.Change) (*unidiff.Allower, error) {
	if changeID != string(change.ID) {
		return noneAllowed, nil
	}
	return allAllowed, nil
}

func (s *Service) getAnonymousWorkspaceAllower(ctx context.Context, workspace *workspaces.Workspace) (*unidiff.Allower, error) {
	cb, err := s.codebaseService.GetByID(ctx, workspace.CodebaseID)
	if err != nil {
		return nil, fmt.Errorf("failed to get codebase: %w", err)
	}
	return s.getAnonymousCodebaseAllower(ctx, cb)
}
func (s *Service) getAnonymousChangeAllower(ctx context.Context, change *changes.Change) (*unidiff.Allower, error) {
	cb, err := s.codebaseService.GetByID(ctx, change.CodebaseID)
	if err != nil {
		return nil, fmt.Errorf("failed to get codebase: %w", err)
	}
	return s.getAnonymousCodebaseAllower(ctx, cb)
}

func (s *Service) getAnonymousCodebaseAllower(ctx context.Context, cb *codebase.Codebase) (*unidiff.Allower, error) {
	if !cb.IsPublic {
		// if codebase is not public, then anonymous users can't see any files.
		return noneAllowed, nil
	}

	// for public codebases, use acls to determine what files are visible.

	aclPolicy, err := s.aclProvider.GetByCodebaseID(ctx, cb.ID)
	if errors.Is(err, sql.ErrNoRows) {
		return noneAllowed, nil
	} else if err != nil {
		return nil, fmt.Errorf("failed to get acl policy: %w", err)
	}

	allowedByID := aclPolicy.Policy.List(
		acl.Identity{Type: acl.Users, ID: "anonymous"},
		acl.ActionWrite,
		acl.Files,
	)

	return unidiff.NewAllower(allowedByID...)
}
