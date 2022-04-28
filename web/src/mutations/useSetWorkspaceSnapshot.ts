import { gql, useMutation } from '@urql/vue'
import type {
  SetWorkspaceSnapshotMutation,
  SetWorkspaceSnapshotMutationVariables,
} from './__generated__/useSetWorkspaceSnapshot'
import type { DeepMaybeRef } from '@vueuse/core'
import type { UpdateResolver } from '@urql/exchange-graphcache'
import type { SetWorkspaceSnapshotInput } from '../__generated__/types'

const SET_WORKSPACE_SNAPSHOT = gql`
  mutation SetWorkspaceSnapshot($input: SetWorkspaceSnapshotInput!) {
    setWorkspaceSnapshot(input: $input) {
      id
      snapshot {
        id
      }
    }
  }
`

export function useSetWorkspaceSnapshot(): (
  input: DeepMaybeRef<SetWorkspaceSnapshotInput>
) => Promise<void> {
  const { executeMutation } = useMutation<
    SetWorkspaceSnapshotMutation,
    DeepMaybeRef<SetWorkspaceSnapshotMutationVariables>
  >(SET_WORKSPACE_SNAPSHOT)
  return async (input): Promise<void> => {
    const result = await executeMutation({ input })
    if (result.error) throw result.error
  }
}

export const setWorkspaceSnapshot: UpdateResolver<
  SetWorkspaceSnapshotMutation,
  SetWorkspaceSnapshotMutationVariables
> = (result, args, cache, info) => {
  // update cache manually if needed
}
