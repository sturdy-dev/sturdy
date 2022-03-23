import type { UpdateResolver } from '@urql/exchange-graphcache'
import type { DeepMaybeRef } from '@vueuse/core'
import gql from 'graphql-tag'
import { useMutation } from '@urql/vue'
import type { WatchWorkspaceInput } from '../__generated__/types'
import type {
  WatchWorkspaceMutation,
  WatchWorkspaceMutationVariables,
  WatchWorkspaceWorkspaceWatchersQuery,
  WatchWorkspaceWorkspaceWatchersQueryVariables,
} from './__generated__/useWatchWorkspace'

const WATCH_WORKSPACE = gql`
  mutation WatchWorkspace($input: WatchWorkspaceInput!) {
    watchWorkspace(input: $input) {
      user {
        id
      }
    }
  }
`

const WORKSPACE_WATCHERS = gql`
  query WatchWorkspaceWorkspaceWatchers($workspaceID: ID!) {
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

export function useWatchWorkspace(): (input: DeepMaybeRef<WatchWorkspaceInput>) => Promise<void> {
  const { executeMutation } = useMutation<
    WatchWorkspaceMutation,
    DeepMaybeRef<WatchWorkspaceMutationVariables>
  >(WATCH_WORKSPACE)
  return async (input) => {
    const result = await executeMutation({ input })
    if (result.error) throw result.error
  }
}

export const watchWorkspaceResolver: UpdateResolver<
  WatchWorkspaceMutation,
  WatchWorkspaceMutationVariables
> = (result, args, cache, info) => {
  // add the user from the workspace's watchers
  cache.updateQuery<
    WatchWorkspaceWorkspaceWatchersQuery,
    WatchWorkspaceWorkspaceWatchersQueryVariables
  >(
    {
      query: WORKSPACE_WATCHERS,
      variables: { workspaceID: args.input.workspaceID },
    },
    (data) => {
      if (!data) return data
      const watcherExists = data.workspace.watchers.some(
        (w) => w.user.id === result.watchWorkspace.user.id
      )
      if (watcherExists) return data
      data.workspace.watchers.push(result.watchWorkspace)
      return data
    }
  )
}
