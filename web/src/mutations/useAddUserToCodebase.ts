import { UpdateResolver } from '@urql/exchange-graphcache'
import { gql, useMutation } from '@urql/vue'
import { AddUserToCodebaseInput } from '../__generated__/types'
import { DeepMaybeRef } from '@vueuse/core'
import {
  AddUserToCodebaseMutation,
  AddUserToCodebaseMutationVariables,
} from './__generated__/useAddUserToCodebase'

const ADD_USER_TO_CODEBASE = gql`
  mutation AddUserToCodebase($input: AddUserToCodebaseInput!) {
    addUserToCodebase(input: $input) {
      id
      members {
        id
      }
    }
  }
`

export function useAddUserToCodebase(): (
  input: DeepMaybeRef<AddUserToCodebaseInput>
) => Promise<void> {
  const { executeMutation } = useMutation<
    AddUserToCodebaseMutation,
    DeepMaybeRef<AddUserToCodebaseMutationVariables>
  >(ADD_USER_TO_CODEBASE)

  return async (input) => {
    const result = await executeMutation({ input })
    if (result.error) {
      throw result.error
    }
  }
}

export const addUserToCodebaseUpdateResolver: UpdateResolver<
  AddUserToCodebaseMutation,
  AddUserToCodebaseMutationVariables
> = (result, args, cache, info) => {
  // not doing anything
}
