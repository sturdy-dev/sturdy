package service

import (
	"context"
	"fmt"

	"getsturdy.com/api/pkg/auth"
	"getsturdy.com/api/pkg/change"
	"getsturdy.com/api/pkg/codebase"
	provider_acl "getsturdy.com/api/pkg/codebase/acl/provider"
	service_codebase "getsturdy.com/api/pkg/codebase/service"
	"getsturdy.com/api/pkg/comments"
	"getsturdy.com/api/pkg/github"
	"getsturdy.com/api/pkg/organization"
	service_organization "getsturdy.com/api/pkg/organization/service"
	"getsturdy.com/api/pkg/review"
	"getsturdy.com/api/pkg/suggestions"
	service_user "getsturdy.com/api/pkg/users/service"
	"getsturdy.com/api/pkg/view"
	"getsturdy.com/api/pkg/workspace"
	"getsturdy.com/api/pkg/workspace/activity"
	service_workspace "getsturdy.com/api/pkg/workspace/service"
)

type Service struct {
	codebaseService     *service_codebase.Service
	userService         service_user.Service
	workspaceService    service_workspace.Service
	aclProvider         *provider_acl.Provider
	organizationService *service_organization.Service
}

func New(
	codebaseService *service_codebase.Service,
	userService service_user.Service,
	workspaceService service_workspace.Service,
	aclProvider *provider_acl.Provider,
	organizationService *service_organization.Service,
) *Service {
	return &Service{
		codebaseService:     codebaseService,
		userService:         userService,
		workspaceService:    workspaceService,
		aclProvider:         aclProvider,
		organizationService: organizationService,
	}
}

type accessType uint

const (
	accessTypeUnknown accessType = iota
	accessTypeRead
	accessTypeWrite
)

// CanRead checks if the user has the read permission on the given object.
func (s *Service) CanRead(ctx context.Context, obj interface{}) error {
	return s.hasAccess(ctx, accessTypeRead, obj)
}

// CanWrite checks if the user has the write permission on the given object.
func (s *Service) CanWrite(ctx context.Context, obj interface{}) error {
	return s.hasAccess(ctx, accessTypeWrite, obj)
}

// hasAccess checks if the user has the given permission on the given object.
func (s *Service) hasAccess(ctx context.Context, at accessType, obj interface{}) error {
	subject, found := auth.FromContext(ctx)
	if !found {
		return fmt.Errorf("subject is not found in the context: %w", auth.ErrUnauthenticated)
	}

	if obj == nil {
		return fmt.Errorf("nil object provided: %w", auth.ErrForbidden)
	}

	switch subject.Type {
	case auth.SubjectUser:
		switch object := obj.(type) {
		case review.Review:
			return s.canUserAccessReview(ctx, subject.ID, at, &object)
		case *review.Review:
			return s.canUserAccessReview(ctx, subject.ID, at, object)
		case view.View:
			return s.canUserAccessView(ctx, subject.ID, at, &object)
		case *view.View:
			return s.canUserAccessView(ctx, subject.ID, at, object)
		case github.GitHubRepository:
			return s.canUserAccessGitHubRepo(ctx, subject.ID, at, &object)
		case *github.GitHubRepository:
			return s.canUserAccessGitHubRepo(ctx, subject.ID, at, object)
		case github.GitHubPullRequest:
			return s.canUserAccessGitHubPullRequest(ctx, subject.ID, at, &object)
		case *github.GitHubPullRequest:
			return s.canUserAccessGitHubPullRequest(ctx, subject.ID, at, object)
		case comments.Comment:
			return s.canUserAccessComment(ctx, subject.ID, at, &object)
		case *comments.Comment:
			return s.canUserAccessComment(ctx, subject.ID, at, object)
		case codebase.Codebase:
			return s.canUserAccessCodebase(ctx, subject.ID, at, &object)
		case *codebase.Codebase:
			return s.canUserAccessCodebase(ctx, subject.ID, at, object)
		case change.Change:
			return s.canUserAccessChange(ctx, subject.ID, at, &object)
		case *change.Change:
			return s.canUserAccessChange(ctx, subject.ID, at, object)
		case change.ChangeCommit:
			return s.canUserAccessChangeCommit(ctx, subject.ID, at, &object)
		case *change.ChangeCommit:
			return s.canUserAccessChangeCommit(ctx, subject.ID, at, object)
		case activity.WorkspaceActivity:
			return s.canUserAccessWorkspaceActivity(ctx, subject.ID, at, &object)
		case *activity.WorkspaceActivity:
			return s.canUserAccessWorkspaceActivity(ctx, subject.ID, at, object)
		case workspace.Workspace:
			return s.canUserAccessWorkspace(ctx, subject.ID, at, &object)
		case *workspace.Workspace:
			return s.canUserAccessWorkspace(ctx, subject.ID, at, object)
		case suggestions.Suggestion:
			return s.canUserAccessSuggestion(ctx, subject.ID, at, &object)
		case *suggestions.Suggestion:
			return s.canUserAccessSuggestion(ctx, subject.ID, at, object)
		case organization.Organization:
			return s.canUserAccessOrganization(ctx, subject.ID, at, &object)
		case *organization.Organization:
			return s.canUserAccessOrganization(ctx, subject.ID, at, object)
		default:
			return fmt.Errorf("unsupported object type '%T' for user: %w", obj, auth.ErrForbidden)
		}
	case auth.SubjectCI:
		switch object := obj.(type) {
		case change.Change:
			return s.canCIAccessChange(ctx, subject.ID, &object)
		case *change.Change:
			return s.canCIAccessChange(ctx, subject.ID, object)
		default:
			return fmt.Errorf("unsupported object type '%T' for ci: %w", obj, auth.ErrForbidden)
		}
	case auth.SubjectAnonymous:
		switch object := obj.(type) {
		case review.Review:
			return s.canAnonymousAccessReview(ctx, at, &object)
		case *review.Review:
			return s.canAnonymousAccessReview(ctx, at, object)
		case github.GitHubRepository:
			return s.canAnonymousAccessGitHubRepo(ctx, at, &object)
		case *github.GitHubRepository:
			return s.canAnonymousAccessGitHubRepo(ctx, at, object)
		case comments.Comment:
			return s.canAnonymousAccessComment(ctx, at, &object)
		case *comments.Comment:
			return s.canAnonymousAccessComment(ctx, at, object)
		case change.Change:
			return s.canAnonymousAccessChange(ctx, at, &object)
		case *change.Change:
			return s.canAnonymousAccessChange(ctx, at, object)
		case change.ChangeCommit:
			return s.canAnonymousAccessChangeCommit(ctx, at, &object)
		case *change.ChangeCommit:
			return s.canAnonymousAccessChangeCommit(ctx, at, object)
		case codebase.Codebase:
			return s.canAnonymousAccessCodebase(ctx, at, &object)
		case *codebase.Codebase:
			return s.canAnonymousAccessCodebase(ctx, at, object)
		case workspace.Workspace:
			return s.canAnonymousAccessWorkspace(ctx, at, &object)
		case *workspace.Workspace:
			return s.canAnonymousAccessWorkspace(ctx, at, object)
		case activity.WorkspaceActivity:
			return s.canAnonymousAccessActivity(ctx, at, &object)
		case *activity.WorkspaceActivity:
			return s.canAnonymousAccessActivity(ctx, at, object)
		case view.View:
			return s.canAnonymousAccessView(ctx, at, &object)
		case *view.View:
			return s.canAnonymousAccessView(ctx, at, object)
		case suggestions.Suggestion:
			return s.canAnonymousAccessSuggestion(ctx, at, &object)
		case *suggestions.Suggestion:
			return s.canAnonymousAccessSuggestion(ctx, at, object)
		case organization.Organization:
			return s.canAnonymousAccessOrganization(ctx, at, &object)
		case *organization.Organization:
			return s.canAnonymousAccessOrganization(ctx, at, object)
		default:
			return fmt.Errorf("unsupported object type '%T' for anonymous: %w", obj, auth.ErrForbidden)
		}
	default:
		return fmt.Errorf("unsupported subject type '%s': %w", subject.Type, auth.ErrForbidden)
	}
}

func (s *Service) canCIAccessChange(ctx context.Context, changeID string, change *change.Change) error {
	if changeID != string(change.ID) {
		return fmt.Errorf("ci doesn't have access to the change: %w", auth.ErrForbidden)
	}
	return nil
}

func (s *Service) canUserAccessChange(ctx context.Context, userID string, at accessType, change *change.Change) error {
	cb, err := s.codebaseService.GetByID(ctx, change.CodebaseID)
	if err != nil {
		return fmt.Errorf("failed to get codebase: %w", err)
	}
	return s.canUserAccessCodebase(ctx, userID, at, cb)
}

func (s *Service) canUserAccessChangeCommit(ctx context.Context, userID string, at accessType, changeCommit *change.ChangeCommit) error {
	cb, err := s.codebaseService.GetByID(ctx, changeCommit.CodebaseID)
	if err != nil {
		return fmt.Errorf("failed to get codebase: %w", err)
	}
	return s.canUserAccessCodebase(ctx, userID, at, cb)
}

func (s *Service) canUserAccessComment(ctx context.Context, userID string, at accessType, comment *comments.Comment) error {
	// user can access the comment if he is the author or if he has the read permission on the codebase
	if comment.UserID == userID {
		return nil
	}

	if at == accessTypeWrite && comment.UserID != userID {
		return fmt.Errorf("only owners can update comments: %w", auth.ErrForbidden)
	}

	cb, err := s.codebaseService.GetByID(ctx, comment.CodebaseID)
	if err != nil {
		return fmt.Errorf("failed to get codebase by id: %w", err)
	}
	return s.canUserAccessCodebase(ctx, userID, at, cb)
}

func (s *Service) canUserAccessCodebase(ctx context.Context, userID string, at accessType, codebase *codebase.Codebase) error {
	// Everyone can read public codebases
	if at == accessTypeRead && codebase.IsPublic {
		return nil
	}

	accessAllowed, err := s.codebaseService.CanAccess(ctx, userID, codebase.ID)
	if err != nil {
		return fmt.Errorf("failed to check if user can access codebase: %w", err)
	}

	if accessAllowed {
		return nil
	}

	if codebase.OrganizationID != nil {
		org, err := s.organizationService.GetByID(ctx, *codebase.OrganizationID)
		if err != nil {
			return fmt.Errorf("failed to check if user can access codebase: %w", err)
		}

		if err := s.canUserAccessOrganization(ctx, userID, accessTypeWrite, org); err == nil {
			return nil
		}
	}

	return fmt.Errorf("user doesn't have acces to the codebase: %w", auth.ErrForbidden)
}

func (s *Service) canAnonymousAccessCodebase(ctx context.Context, at accessType, codebase *codebase.Codebase) error {
	if at == accessTypeRead && codebase.IsPublic {
		return nil
	}
	return fmt.Errorf("anonymous users can only read public codebaes: %w", auth.ErrForbidden)
}

func (s *Service) canUserAccessSuggestion(ctx context.Context, userID string, at accessType, suggestion *suggestions.Suggestion) error {
	cb, err := s.codebaseService.GetByID(ctx, suggestion.CodebaseID)
	if err != nil {
		return fmt.Errorf("failed to get codebase: %w", err)
	}
	return s.canUserAccessCodebase(ctx, userID, at, cb)
}

func (s *Service) canUserAccessWorkspace(ctx context.Context, userID string, at accessType, workspace *workspace.Workspace) error {
	// user can access a workspace if they can access the codebase it's in
	cb, err := s.codebaseService.GetByID(ctx, workspace.CodebaseID)
	if err != nil {
		return fmt.Errorf("failed to get codebase: %w", err)
	}

	return s.canUserAccessCodebase(ctx, userID, at, cb)
}

func (s *Service) canAnonymousAccessSuggestion(ctx context.Context, at accessType, suggestion *suggestions.Suggestion) error {
	if at != accessTypeRead {
		return fmt.Errorf("anonymous users can only read suggestion: %w", auth.ErrForbidden)
	}

	// user can access a workspace if they can access the codebase it's in
	cb, err := s.codebaseService.GetByID(ctx, suggestion.CodebaseID)
	if err != nil {
		return fmt.Errorf("failed to get codebase: %w", err)
	}

	return s.canAnonymousAccessCodebase(ctx, at, cb)
}

func (s *Service) canAnonymousAccessWorkspace(ctx context.Context, at accessType, workspace *workspace.Workspace) error {
	if at != accessTypeRead {
		return fmt.Errorf("anonymous users can only read workspaces: %w", auth.ErrForbidden)
	}

	// user can access a workspace if they can access the codebase it's in
	cb, err := s.codebaseService.GetByID(ctx, workspace.CodebaseID)
	if err != nil {
		return fmt.Errorf("failed to get codebase: %w", err)
	}

	return s.canAnonymousAccessCodebase(ctx, at, cb)
}

func (s *Service) canAnonymousAccessActivity(ctx context.Context, at accessType, activity *activity.WorkspaceActivity) error {
	if at != accessTypeRead {
		return fmt.Errorf("anonymous users can only read activities: %w", auth.ErrForbidden)
	}

	// user can access an activity if they can access the workspace it's from
	ws, err := s.workspaceService.GetByID(ctx, activity.WorkspaceID)
	if err != nil {
		return fmt.Errorf("failed to get workspace: %w", err)
	}

	return s.canAnonymousAccessWorkspace(ctx, at, ws)
}

func (s *Service) canAnonymousAccessComment(ctx context.Context, at accessType, comment *comments.Comment) error {
	if at != accessTypeRead {
		return fmt.Errorf("anonymous users can only read comments: %w", auth.ErrForbidden)
	}

	// user can access a comment if they can access the codebase it's in
	cb, err := s.codebaseService.GetByID(ctx, comment.CodebaseID)
	if err != nil {
		return fmt.Errorf("failed to get codebase: %w", err)
	}

	return s.canAnonymousAccessCodebase(ctx, at, cb)
}

func (s *Service) canAnonymousAccessChange(ctx context.Context, at accessType, change *change.Change) error {
	if at != accessTypeRead {
		return fmt.Errorf("anonymous users can only read changes: %w", auth.ErrForbidden)
	}

	// user can access a change if they can access the codebase it's in
	cb, err := s.codebaseService.GetByID(ctx, change.CodebaseID)
	if err != nil {
		return fmt.Errorf("failed to get codebase: %w", err)
	}

	return s.canAnonymousAccessCodebase(ctx, at, cb)
}

func (s *Service) canAnonymousAccessChangeCommit(ctx context.Context, at accessType, changeCommit *change.ChangeCommit) error {
	if at != accessTypeRead {
		return fmt.Errorf("anonymous users can only read changes: %w", auth.ErrForbidden)
	}

	// user can access a change if they can access the codebase it's in
	cb, err := s.codebaseService.GetByID(ctx, changeCommit.CodebaseID)
	if err != nil {
		return fmt.Errorf("failed to get codebase: %w", err)
	}

	return s.canAnonymousAccessCodebase(ctx, at, cb)
}

func (s *Service) canUserAccessWorkspaceActivity(ctx context.Context, userID string, at accessType, a *activity.WorkspaceActivity) error {
	// user can access a workspace activity if they can access the workspace it's in
	ws, err := s.workspaceService.GetByID(ctx, a.WorkspaceID)
	if err != nil {
		return fmt.Errorf("failed to get workspace: %w", err)
	}
	return s.canUserAccessWorkspace(ctx, userID, at, ws)
}

func (s *Service) canUserAccessGitHubRepo(ctx context.Context, userID string, at accessType, repo *github.GitHubRepository) error {
	// user can access a github repository if they can access the codebase it's in
	cb, err := s.codebaseService.GetByID(ctx, repo.CodebaseID)
	if err != nil {
		return fmt.Errorf("failed to get codebase: %w", err)
	}
	return s.canUserAccessCodebase(ctx, userID, at, cb)
}

func (s *Service) canUserAccessGitHubPullRequest(ctx context.Context, userID string, at accessType, pr *github.GitHubPullRequest) error {
	// user can access a github pr if they can access the codebase it's in
	cb, err := s.codebaseService.GetByID(ctx, pr.CodebaseID)
	if err != nil {
		return fmt.Errorf("failed to get codebase: %w", err)
	}
	return s.canUserAccessCodebase(ctx, userID, at, cb)
}

func (s *Service) canAnonymousAccessGitHubRepo(ctx context.Context, at accessType, repo *github.GitHubRepository) error {
	if at != accessTypeRead {
		return fmt.Errorf("anonymous users can only read github repositories: %w", auth.ErrForbidden)
	}
	// user can access a github repository if they can access the codebase it's in
	cb, err := s.codebaseService.GetByID(ctx, repo.CodebaseID)
	if err != nil {
		return fmt.Errorf("failed to get codebase: %w", err)
	}
	return s.canAnonymousAccessCodebase(ctx, at, cb)
}

func (s *Service) canAnonymousAccessView(ctx context.Context, at accessType, v *view.View) error {
	if at == accessTypeWrite {
		return fmt.Errorf("anonymous users can only read views: %w", auth.ErrForbidden)
	}
	// user can access a view if they can access the codebase it's in
	cb, err := s.codebaseService.GetByID(ctx, v.CodebaseID)
	if err != nil {
		return fmt.Errorf("failed to get codebase: %w", err)
	}
	return s.canAnonymousAccessCodebase(ctx, at, cb)
}

func (s *Service) canUserAccessView(ctx context.Context, userID string, at accessType, v *view.View) error {
	if at == accessTypeWrite && v.UserID != userID {
		return fmt.Errorf("only owner can write to a view: %w", auth.ErrForbidden)
	}
	// user can access a view if they can access the codebase it's in
	cb, err := s.codebaseService.GetByID(ctx, v.CodebaseID)
	if err != nil {
		return fmt.Errorf("failed to get codebase: %w", err)
	}
	return s.canUserAccessCodebase(ctx, userID, at, cb)
}

func (s *Service) canUserAccessReview(ctx context.Context, userID string, at accessType, r *review.Review) error {
	// user can access a review if they can access the codebase it's in
	cb, err := s.codebaseService.GetByID(ctx, r.CodebaseID)
	if err != nil {
		return fmt.Errorf("failed to get codebase: %w", err)
	}
	return s.canUserAccessCodebase(ctx, userID, at, cb)
}

func (s *Service) canAnonymousAccessReview(ctx context.Context, at accessType, r *review.Review) error {
	// anonymous users can only read reviews.
	if at != accessTypeRead {
		return fmt.Errorf("anonymous users can only read reviews: %w", auth.ErrForbidden)
	}
	// user can access a review if they can access the codebase it's in
	cb, err := s.codebaseService.GetByID(ctx, r.CodebaseID)
	if err != nil {
		return fmt.Errorf("failed to get codebase: %w", err)
	}
	return s.canAnonymousAccessCodebase(ctx, at, cb)
}

func (s *Service) canUserAccessOrganization(ctx context.Context, userID string, at accessType, org *organization.Organization) error {
	// user can access a organization if they are a member of it
	_, err := s.organizationService.GetMemberByUserIDAndOrganizationID(ctx, userID, org.ID)
	if err == nil {
		return nil
	}

	// user can read (but not write) a organization if they are a member of any of it's codebases
	if at == accessTypeRead {
		ok, err := s.codebaseService.UserIsMemberOfCodebaseInOrganization(ctx, userID, org.ID)
		if err != nil {
			return fmt.Errorf("user does not have access to organization: %w", auth.ErrForbidden)
		}
		if ok {
			return nil
		}
	}

	return fmt.Errorf("user does not have access to organization: %w", auth.ErrForbidden)
}

func (s *Service) canAnonymousAccessOrganization(ctx context.Context, at accessType, org *organization.Organization) error {
	return fmt.Errorf("anonymous users can't access organizations: %w", auth.ErrForbidden)
}
