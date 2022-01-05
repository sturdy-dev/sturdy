import { DeepMaybeRef, MaybeRef } from '@vueuse/core'
import { gql, useSubscription } from '@urql/vue'
import {
  UpdatedWorkspaceWatchersSubscription,
  UpdatedWorkspaceWatchersSubscriptionVariables,
  WorkspaceWatchersQuery,
  WorkspaceWatchersQueryVariables,
} from './__generated__/useUpdatedWorkspaceWathcers'
import { UpdateResolver } from '@urql/exchange-graphcache'

const UPDATED_WORKSPACE_WATCHERS = gql`
  subscription UpdatedWorkspaceWatchers($workspaceID: ID!) {
    updatedWorkspaceWatchers(workspaceID: $workspaceID) {
      status
      user {
        id
      }
    }
  }
`

const WORKSPACE_WATCHERS = gql`
  query WorkspaceWatchers($workspaceID: ID!) {
    workspace(id: $workspaceID) {
      id
      watchers {
        user {
          id
        }
      }
    }
  }
`

export function useUpdatedWorkspaceWatchers(workspaceID: MaybeRef<string>) {
  useSubscription<
    UpdatedWorkspaceWatchersSubscription,
    DeepMaybeRef<UpdatedWorkspaceWatchersSubscriptionVariables>
  >({
    query: UPDATED_WORKSPACE_WATCHERS,
    variables: { workspaceID },
  })
}

export const updatedWorkspaceWatchersResolver: UpdateResolver<
  UpdatedWorkspaceWatchersSubscription,
  UpdatedWorkspaceWatchersSubscriptionVariables
> = (result, args, cache, info) => {
  cache.updateQuery<WorkspaceWatchersQuery, WorkspaceWatchersQueryVariables>(
    {
      query: WORKSPACE_WATCHERS,
      variables: {
        workspaceID: args.workspaceID,
      },
    },
    (data) => {
      if (!data) return data

      // replace exsisting user wathcer
      data.workspace.watchers = data.workspace.watchers
        .filter((watcher) => watcher.user.id !== result.updatedWorkspaceWatchers.user.id)
        .concat(result.updatedWorkspaceWatchers)

      return data
    }
  )
}
