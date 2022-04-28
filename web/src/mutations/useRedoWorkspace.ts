import type { UpdateResolver } from '@urql/exchange-graphcache'
import type { DeepMaybeRef } from '@vueuse/core'
import gql from 'graphql-tag'
import { useMutation } from '@urql/vue'
import type {
  RedoWorkspaceMutation,
  RedoWorkspaceMutationVariables,
} from './__generated__/useRedoWorkspace'

const UNDO_WORKSPACE = gql`
  mutation RedoWorkspace($id: ID!) {
    redoWorkspace(id: $id) {
      id
    }
  }
`

export function useUndoWorkspace(): (id: DeepMaybeRef<string>) => Promise<void> {
  const { executeMutation } = useMutation<
    RedoWorkspaceMutation,
    DeepMaybeRef<RedoWorkspaceMutationVariables>
  >(UNDO_WORKSPACE)
  return async (id) => {
    const result = await executeMutation({ id })
    if (result.error) throw result.error
  }
}

export const redoWorkspaceResolver: UpdateResolver<
  RedoWorkspaceMutation,
  RedoWorkspaceMutationVariables
> = (result, args, cache, info) => {}
