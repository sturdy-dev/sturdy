import { gql, useMutation } from '@urql/vue'
import { DeepMaybeRef } from '@vueuse/core'
import { RemovePatchesInput } from '../__generated__/types'
import {
  RemovePatchesMutation,
  RemovePatchesMutationVariables,
} from './__generated__/useRemovePatches'

const REMOVE_PATCHES = gql`
  mutation RemovePatches($input: RemovePatchesInput!) {
    removePatches(input: $input) {
      id
    }
  }
`

export function useRemovePatches(): (input: DeepMaybeRef<RemovePatchesInput>) => Promise<void> {
  const { executeMutation } = useMutation<
    RemovePatchesMutation,
    DeepMaybeRef<RemovePatchesMutationVariables>
  >(REMOVE_PATCHES)

  return async (input) => {
    const result = await executeMutation({ input })
    if (result.error) {
      throw result.error
    }
  }
}
