import { UpdateResolver } from '@urql/exchange-graphcache'
import { gql, useMutation } from '@urql/vue'
import { RemoveUserFromOrganizationInput } from '../__generated__/types'
import { DeepMaybeRef } from '@vueuse/core'
import {
  RemoveUserFromOrganizationMutation,
  RemoveUserFromOrganizationMutationVariables,
} from './__generated__/useRemoveUserFromOrganization'

const REMOVE_USER_FROM_ORGANIZATION = gql`
  mutation RemoveUserFromOrganization($input: RemoveUserFromOrganizationInput!) {
    removeUserFromOrganization(input: $input) {
      id
      members {
        id
      }
    }
  }
`

export function useRemoveUserFromOrganization(): (
  input: DeepMaybeRef<RemoveUserFromOrganizationInput>
) => Promise<void> {
  const { executeMutation } = useMutation<
    RemoveUserFromOrganizationMutation,
    DeepMaybeRef<RemoveUserFromOrganizationMutationVariables>
  >(REMOVE_USER_FROM_ORGANIZATION)

  return async (input) => {
    const result = await executeMutation({ input })
    if (result.error) {
      throw result.error
    }
  }
}

export const removeUserFromOrganizationUpdateResolver: UpdateResolver<
  RemoveUserFromOrganizationMutation,
  RemoveUserFromOrganizationMutationVariables
> = (result, args, cache, info) => {
  // not doing anything
}
