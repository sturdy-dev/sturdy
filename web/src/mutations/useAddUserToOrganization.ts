import { UpdateResolver } from '@urql/exchange-graphcache'
import { gql, useMutation } from '@urql/vue'
import { AddUserToOrganizationInput } from '../__generated__/types'
import { DeepMaybeRef } from '@vueuse/core'
import {
  AddUserToOrganizationMutation,
  AddUserToOrganizationMutationVariables,
} from './__generated__/useAddUserToOrganization'

const ADD_USER_TO_ORGANIZATION = gql`
  mutation AddUserToOrganization($input: AddUserToOrganizationInput!) {
    addUserToOrganization(input: $input) {
      id
      members {
        id
      }
    }
  }
`

export function useAddUserToOrganization(): (
  input: DeepMaybeRef<AddUserToOrganizationInput>
) => Promise<void> {
  const { executeMutation } = useMutation<
    AddUserToOrganizationMutation,
    DeepMaybeRef<AddUserToOrganizationMutationVariables>
  >(ADD_USER_TO_ORGANIZATION)

  return async (input) => {
    const result = await executeMutation({ input })
    if (result.error) {
      throw result.error
    }
  }
}

export const addUserToOrganizationUpdateResolver: UpdateResolver<
  AddUserToOrganizationMutation,
  AddUserToOrganizationMutationVariables
> = (result, args, cache, info) => {
  // not doing anything
}
