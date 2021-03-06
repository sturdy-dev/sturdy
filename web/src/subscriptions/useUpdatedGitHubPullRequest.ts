import type { DeepMaybeRef } from '@vueuse/core'
import { gql, useSubscription } from '@urql/vue'
import type {
  UpdatedGitHubPullRequestSubscription,
  UpdatedGitHubPullRequestSubscriptionVariables,
} from './__generated__/useUpdatedGitHubPullRequest'
import type { UpdateResolver } from '@urql/exchange-graphcache'

const UPDATED_GIT_HUB_PULL_REQUEST = gql`
  subscription UpdatedGitHubPullRequest($workspaceID: ID!) {
    updatedGitHubPullRequest(workspaceID: $workspaceID) {
      id
      pullRequestNumber
      open
      merged
      state
      mergedAt
      workspace {
        id
        upToDateWithTrunk
        change {
          id
        }
      }
    }
  }
`

export function useUpdatedGitHubPullRequest(workspaceID: DeepMaybeRef<string>) {
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
