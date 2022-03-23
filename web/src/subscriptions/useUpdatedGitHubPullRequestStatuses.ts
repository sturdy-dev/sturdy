import { gql, useSubscription } from '@urql/vue'
import type { DeepMaybeRef, MaybeRef } from '@vueuse/core'
import type {
  WorkspaceQuery,
  WorkspaceQueryVariables,
  UpdatedGitHubPullRequestStatusesSubscription,
  UpdatedGitHubPullRequestStatusesSubscriptionVariables,
} from './__generated__/useUpdatedGitHubPullRequestStatuses'
import type { UpdateResolver } from '@urql/exchange-graphcache'

const UPDATED_GITHUB_PULL_REQUEST_STATUSES = gql`
  subscription UpdatedGitHubPullRequestStatuses($id: ID!) {
    updatedGitHubPullRequestStatuses(id: $id) {
      id
      type
      title
      description
      timestamp
      detailsUrl
      gitHubPullRequest {
        id
        workspace {
          id
        }
      }
    }
  }
`

const WORKSPACE_PULL_REQUEST_STATUSES = gql`
  query workspace($id: ID!) {
    workspace(id: $id) {
      id
      gitHubPullRequest {
        statuses {
          id
          title
        }
      }
    }
  }
`

export function useUpdatedGitHubPullRequestStatuses(id: MaybeRef<string>) {
  useSubscription<
    UpdatedGitHubPullRequestStatusesSubscription,
    DeepMaybeRef<UpdatedGitHubPullRequestStatusesSubscriptionVariables>
  >({
    query: UPDATED_GITHUB_PULL_REQUEST_STATUSES,
    variables: { id: id },
  })
}

export const updatedGitHubPullRequestStatusesResolver: UpdateResolver<
  UpdatedGitHubPullRequestStatusesSubscription,
  UpdatedGitHubPullRequestStatusesSubscriptionVariables
> = (result, args, cache, info) => {
  const newStatus = result.updatedGitHubPullRequestStatuses
  if (!newStatus.gitHubPullRequest) return
  cache.updateQuery<WorkspaceQuery, WorkspaceQueryVariables>(
    {
      query: WORKSPACE_PULL_REQUEST_STATUSES,
      variables: { id: newStatus.gitHubPullRequest.workspace.id },
    },
    (data) => {
      if (!data) return data
      if (!data.workspace.gitHubPullRequest) return data

      data.workspace.gitHubPullRequest.statuses = [
        ...data.workspace.gitHubPullRequest.statuses.filter((s) => s.title != newStatus.title),
        newStatus,
      ]
      return data
    }
  )
}
