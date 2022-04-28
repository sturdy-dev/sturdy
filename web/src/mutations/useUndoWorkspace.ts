import type { UpdateResolver } from '@urql/exchange-graphcache'
import type { DeepMaybeRef } from '@vueuse/core'
import gql from 'graphql-tag'
import { useMutation } from '@urql/vue'
import type {
  UndoWorkspaceMutation,
  UndoWorkspaceMutationVariables,
} from './__generated__/useUndoWorkspace'

const UNDO_WORKSPACE = gql`
  mutation UndoWorkspace($id: ID!) {
    undoWorkspace(id: $id) {
      id
    }
  }
`

export function useUndoWorkspace(): (id: DeepMaybeRef<string>) => Promise<void> {
  const { executeMutation } = useMutation<
    UndoWorkspaceMutation,
    DeepMaybeRef<UndoWorkspaceMutationVariables>
  >(UNDO_WORKSPACE)
  return async (id) => {
    const result = await executeMutation({ id })
    if (result.error) throw result.error
  }
}

export const undoWorkspaceResolver: UpdateResolver<
  UndoWorkspaceMutation,
  UndoWorkspaceMutationVariables
> = (result, args, cache, info) => {}
