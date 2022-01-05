import { UpdateResolver } from '@urql/exchange-graphcache'
import { DeepMaybeRef } from '@vueuse/core'
import gql from 'graphql-tag'
import { useMutation } from '@urql/vue'
import { UnwatchWorkspaceInput } from '../__generated__/types'
import {
  UnwatchWorkspaceMutation,
  UnwatchWorkspaceMutationVariables,
  UnwatchWorkspaceWorkspaceWatchersQuery,
  UnwatchWorkspaceWorkspaceWatchersQueryVariables,
} from './__generated__/useUnwatchWorkspace'

const WATCH_WORKSPACE = gql`
  mutation UnwatchWorkspace($input: UnwatchWorkspaceInput!) {
    unwatchWorkspace(input: $input) {
      user {
        id
      }
      status
    }
  }
`

const WORKSPACE_WATCHERS = gql`
  query UnwatchWorkspaceWorkspaceWatchers($workspaceID: ID!) {
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

export function useUnwatchWorkspace(): (
  input: DeepMaybeRef<UnwatchWorkspaceInput>
) => Promise<void> {
  const { executeMutation } = useMutation<
    UnwatchWorkspaceMutation,
    DeepMaybeRef<UnwatchWorkspaceMutationVariables>
  >(WATCH_WORKSPACE)
  return async (input) => {
    const result = await executeMutation({ input })
    if (result.error) throw result.error
  }
}

export const unwatchWorkspaceResolver: UpdateResolver<
  UnwatchWorkspaceMutation,
  UnwatchWorkspaceMutationVariables
> = (result, args, cache, info) => {
  // remove the user from the workspace's watchers
  cache.updateQuery<
    UnwatchWorkspaceWorkspaceWatchersQuery,
    UnwatchWorkspaceWorkspaceWatchersQueryVariables
  >(
    {
      query: WORKSPACE_WATCHERS,
      variables: { workspaceID: args.input.workspaceID },
    },
    (data) => {
      if (!data) return data
      data.workspace.watchers = data.workspace.watchers.filter(
        (w) => w.user.id !== result.unwatchWorkspace.user.id
      )
      return data
    }
  )
}
