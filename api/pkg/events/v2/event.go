package events

import (
	"getsturdy.com/api/pkg/codebases"
	"getsturdy.com/api/pkg/github"
	"getsturdy.com/api/pkg/notification"
	"getsturdy.com/api/pkg/onboarding"
	"getsturdy.com/api/pkg/organization"
	"getsturdy.com/api/pkg/review"
	"getsturdy.com/api/pkg/statuses"
	"getsturdy.com/api/pkg/view"
	"getsturdy.com/api/pkg/workspaces"
	"getsturdy.com/api/pkg/workspaces/watchers"
)

type Type uint

const (
	TypeUndefined Type = iota

	CodebaseEvent
	CodebaseUpdated

	ViewUpdated
	ViewStatusUpdated

	WorkspaceUpdated
	WorkspaceUpdatedComments
	WorkspaceUpdatedReviews
	WorkspaceUpdatedActivity
	WorkspaceUpdatedSnapshot
	WorkspaceUpdatedPresence
	WorkspaceUpdatedSuggestion
	WorkspaceWatchingStatusUpdated

	ReviewUpdated
	GitHubPRUpdated
	NotificationEvent
	StatusUpdated
	CompletedOnboardingStep
	OrganizationUpdated
)

func (t Type) String() string {
	switch t {
	case TypeUndefined:
		return "TypeUndefined"
	case CodebaseEvent:
		return "CodebaseEvent"
	case CodebaseUpdated:
		return "CodebaseUpdated"
	case ViewUpdated:
		return "ViewUpdated"
	case ViewStatusUpdated:
		return "ViewStatusUpdated"
	case WorkspaceUpdated:
		return "WorkspaceUpdated"
	case WorkspaceUpdatedComments:
		return "WorkspaceUpdatedComments"
	case WorkspaceUpdatedReviews:
		return "WorkspaceUpdatedReviews"
	case ReviewUpdated:
		return "ReviewUpdated"
	case WorkspaceUpdatedActivity:
		return "WorkspaceUpdatedActivity"
	case WorkspaceUpdatedSnapshot:
		return "WorkspaceUpdatedSnapshot"
	case WorkspaceUpdatedPresence:
		return "WorkspaceUpdatedPresence"
	case WorkspaceUpdatedSuggestion:
		return "WorkspaceUpdatedSuggestion"
	case GitHubPRUpdated:
		return "GitHubPRUpdated"
	case NotificationEvent:
		return "NotificationEvent"
	case StatusUpdated:
		return "StatusUpdated"
	case CompletedOnboardingStep:
		return "CompletedOnboardingStep"
	case WorkspaceWatchingStatusUpdated:
		return "WorkspaceWatchingStatusUpdated"
	case OrganizationUpdated:
		return "OrganizationUpdated"
	default:
		return "Unknown"
	}
}

type event struct {
	Type Type

	Codebase          *codebases.Codebase
	View              *view.View
	Workspace         *workspaces.Workspace
	Review            *review.Review
	GitHubPullRequest *github.PullRequest
	Notification      *notification.Notification
	Status            *statuses.Status
	OnboardingStep    *onboarding.Step
	WorkspaceWatcher  *watchers.Watcher
	Organization      *organization.Organization
}
