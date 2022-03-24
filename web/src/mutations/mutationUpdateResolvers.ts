import type { OptimisticMutationResolver, UpdateResolver } from '@urql/exchange-graphcache'
import { createCommentUpdateResolver } from './useCreateComment'
import { createOrUpdateReviewUpdateResolver } from './useCreateOrUpdateReview'
import { requestReviewUpdateResolver } from './useRequestReview'
import { createWorkspaceUpdateResolver } from './useCreateWorkspace'
import { updateCommentUpdateResolver } from './useUpdateComment'
import { deleteCommentUpdateResolver } from './useDeleteComment'
import {
  openWorkspaceOnViewOptimisticMutationResolver,
  openWorkspaceOnViewUpdateResolver,
} from './useOpenWorkspaceOnView'
import { updateNotificationPreferenceResolver } from './useUpdateNotificationPreference'
import { verifyEmailResolver } from './useVerifyEmail'
import { watchWorkspaceResolver } from './useWatchWorkspace'
import { unwatchWorkspaceResolver } from './useUnwatchWorkspace'
import { createServiceTokenResolver } from './useCreateServiceToken'
import { createOrUpdateBuildkiteIntegrationUpdateResolver } from './useCreateOrUpdateBuildkiteIntegration'
import { triggerInstantIntegrationUpdateResolver } from './useTriggerInstantIntegration'
import { createSuggestionUpdateResolver } from './useCreateSuggestion'
import { landWorkspaceChangeUpdateResolver } from './useLandWorkspaceChange'
import { setupGitHubUpdateResolver } from './useSetupGitHubRepository'
import { createCodebaseUpdateResolver } from './useCreateCodebase'
import { createOrganizationUpdateResolver } from './useCreateOrganization'
import { removeUserFromOrganizationUpdateResolver } from './useRemoveUserFromOrganization'
import { addUserToOrganizationUpdateResolver } from './useAddUserToOrganization'
import { addUserToCodebaseUpdateResolver } from './useAddUserToCodebase'
import { removeUserFromCodebaseUpdateResolver } from './useRemoveUserFromCodebase'
import { updateInstallationUpdateResolver } from './useUpdateInstallation'
import { updateOrganizationResolver } from './useUpdateOrganization'
import { createOrUpdateCodebaseRemoteUpdateResolver } from './useCreateOrUpdateGitRemote'
import { updateWorkspaceUpdateResolver } from './useUpdateWorkspace'
import { pullCodebaseUpdateResolver } from './usePullCodebase'
import { pushCodebaseUpdateResolver } from './usePushCodebase'
import { archiveWorkspaceUpdateResolver } from './useArchiveWorkspace'

export const mutationUpdateResolvers: Record<string, UpdateResolver> = {
  createComment: createCommentUpdateResolver,
  updateComment: updateCommentUpdateResolver,
  deleteComment: deleteCommentUpdateResolver,
  createOrUpdateReview: createOrUpdateReviewUpdateResolver,
  requestReview: requestReviewUpdateResolver,
  createWorkspace: createWorkspaceUpdateResolver,
  openWorkspaceOnView: openWorkspaceOnViewUpdateResolver,
  updateNotificationPreference: updateNotificationPreferenceResolver,
  verifyEmail: verifyEmailResolver,
  watchWorkspace: watchWorkspaceResolver,
  unwatchWorkspace: unwatchWorkspaceResolver,
  createServiceToken: createServiceTokenResolver,
  createOrUpdateBuildkiteIntegration: createOrUpdateBuildkiteIntegrationUpdateResolver,
  triggerInstantIntegration: triggerInstantIntegrationUpdateResolver,
  createSuggestion: createSuggestionUpdateResolver,
  landWorkspaceChange: landWorkspaceChangeUpdateResolver,
  createCodebase: createCodebaseUpdateResolver,
  setupGitHubRepository: setupGitHubUpdateResolver,
  createOrganization: createOrganizationUpdateResolver,
  removeUserFromOrganization: removeUserFromOrganizationUpdateResolver,
  addUserToOrganization: addUserToOrganizationUpdateResolver,
  addUserToCodebase: addUserToCodebaseUpdateResolver,
  removeUserFromCodebase: removeUserFromCodebaseUpdateResolver,
  updateInstallation: updateInstallationUpdateResolver,
  updateOrganization: updateOrganizationResolver,
  createOrUpdateCodebaseRemote: createOrUpdateCodebaseRemoteUpdateResolver,
  updateWorkspace: updateWorkspaceUpdateResolver,
  pullCodebase: pullCodebaseUpdateResolver,
  pushCodebase: pushCodebaseUpdateResolver,
  archiveWorkspace: archiveWorkspaceUpdateResolver,
}

export const optimisticMutationResolvers: Record<string, OptimisticMutationResolver> = {
  openWorkspaceOnView: openWorkspaceOnViewOptimisticMutationResolver,
}
