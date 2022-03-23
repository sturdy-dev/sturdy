import type { UpdateResolver } from '@urql/exchange-graphcache'
import { updatedCommentUpdateResolver } from './useUpdatedComment'
import { updatedCodebaseUpdateResolver } from './useUpdatedCodebase'
import { updatedNotificationsUpdateResolver } from './useUpdatedNotifications'
import { updatedWorkspaceActivityUpdateResolver } from './useUpdatedWorkspaceActivity'
import { updatedWorkspaceUpdateResolver } from './useUpdatedWorkspace'
import { updatedGitHubPullRequestUpdateResolver } from './useUpdatedGitHubPullRequest'
import { updatedWorkspacePresenceResolver } from './useUpdatedWorkspacePresence'
import { updatedSuggestionResolver } from './useUpdatedSuggestion'
import { updatedChangesStatusesResolver } from './useUpdatedChangesStatuses'
import { updatedGitHubPullRequestStatusesResolver } from './useUpdatedGitHubPullRequestStatuses'
import { updatedWorkspaceWatchersResolver } from './useUpdatedWorkspaceWathcers'
import { updatedReviewsResolver } from './useUpdatedReviews'
import { updatedViewsUpdateResolver } from './useUpdatedViews'
import { updatedOrganizationUpdateResolver } from './useUpdatedOrganization'
import { updateWorkspaceDiffsResolver } from './useUpdatedWorkspaceDiffs'

export const subscriptionUpdateResolvers: Record<string, UpdateResolver> = {
  updatedComment: updatedCommentUpdateResolver,
  updatedCodebase: updatedCodebaseUpdateResolver,
  updatedNotifications: updatedNotificationsUpdateResolver,
  updatedWorkspaceActivity: updatedWorkspaceActivityUpdateResolver,
  updatedWorkspace: updatedWorkspaceUpdateResolver,
  updatedGitHubPullRequest: updatedGitHubPullRequestUpdateResolver,
  updatedViews: updatedViewsUpdateResolver,
  updatedWorkspacePresence: updatedWorkspacePresenceResolver,
  updatedSuggestion: updatedSuggestionResolver,
  updatedChangesStatuses: updatedChangesStatusesResolver,
  updatedGitHubPullRequestStatuses: updatedGitHubPullRequestStatusesResolver,
  updatedWorkspaceWatchers: updatedWorkspaceWatchersResolver,
  updatedReviews: updatedReviewsResolver,
  updatedOrganization: updatedOrganizationUpdateResolver,
  updatedWorkspaceDiffs: updateWorkspaceDiffsResolver,
}
