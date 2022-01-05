import { DeepMaybeRef, MaybeRef } from '@vueuse/core'
import { gql, useSubscription } from '@urql/vue'
import {
  UpdatedGitHubPullRequestSubscription,
  UpdatedGitHubPullRequestSubscriptionVariables,
} from './__generated__/useUpdatedGitHubPullRequest'
import { UpdateResolver } from '@urql/exchange-graphcache'

const UPDATED_GIT_HUB_PULL_REQUEST = gql`
  subscription UpdatedGitHubPullRequest($workspaceID: ID!) {
    updatedGitHubPullRequest(workspaceID: $workspaceID) {
      id
      pullRequestNumber
      open
      merged
      mergedAt
      workspace {
        id
        upToDateWithTrunk
      }
    }
  }
`

export function useUpdatedGitHubPullRequest(workspaceID: MaybeRef<string>) {
  useSubscription<
    UpdatedGitHubPullRequestSubscription,
    UpdatedGitHubPullRequestSubscription,
    DeepMaybeRef<UpdatedGitHubPullRequestSubscriptionVariables>
  >({
    query: UPDATED_GIT_HUB_PULL_REQUEST,
    variables: { workspaceID: workspaceID },
  })
}

export const updatedGitHubPullRequestUpdateResolver: UpdateResolver<
  UpdatedGitHubPullRequestSubscription,
  UpdatedGitHubPullRequestSubscriptionVariables
> = (parent, args, cache, info) => {
  // Update cache manually if needed
}
