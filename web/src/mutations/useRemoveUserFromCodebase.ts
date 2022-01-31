import { UpdateResolver } from '@urql/exchange-graphcache'
import { gql, useMutation } from '@urql/vue'
import { RemoveUserFromCodebaseInput } from '../__generated__/types'
import { DeepMaybeRef } from '@vueuse/core'
import {
  RemoveUserFromCodebaseMutation,
  RemoveUserFromCodebaseMutationVariables,
} from './__generated__/useRemoveUserFromCodebase'

const REMOVE_USER_FROM_CODEBASE = gql`
  mutation RemoveUserFromCodebase($input: RemoveUserFromCodebaseInput!) {
    removeUserFromCodebase(input: $input) {
      id
      members {
        id
      }
    }
  }
`

export function useRemoveUserFromCodebase(): (
  input: DeepMaybeRef<RemoveUserFromCodebaseInput>
) => Promise<void> {
  const { executeMutation } = useMutation<
    RemoveUserFromCodebaseMutation,
    DeepMaybeRef<RemoveUserFromCodebaseMutationVariables>
  >(REMOVE_USER_FROM_CODEBASE)

  return async (input) => {
    const result = await executeMutation({ input })
    if (result.error) {
      throw result.error
    }
  }
}

export const removeUserFromCodebaseUpdateResolver: UpdateResolver<
  RemoveUserFromCodebaseMutation,
  RemoveUserFromCodebaseMutationVariables
> = (result, args, cache, info) => {
  // not doing anything
}
