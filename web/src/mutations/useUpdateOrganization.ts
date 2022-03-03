import { UpdateResolver } from '@urql/exchange-graphcache'
import { gql, useMutation } from '@urql/vue'
import { UpdateOrganizationInput } from '../__generated__/types'
import { DeepMaybeRef } from '@vueuse/core'
import {
  UpdateOrganizationMutation,
  UpdateOrganizationMutationVariables,
} from './__generated__/useUpdateOrganization'

const UPDATE_ORGANIZATION = gql`
  mutation UpdateOrganization($input: UpdateOrganizationInput!) {
    updateOrganization(input: $input) {
      id
      name
    }
  }
`

export function useUpdateOrganization(): (
  input: DeepMaybeRef<UpdateOrganizationInput>
) => Promise<void> {
  const { executeMutation } = useMutation<
    UpdateOrganizationMutation,
    DeepMaybeRef<UpdateOrganizationMutationVariables>
  >(UPDATE_ORGANIZATION)

  return async (input) => {
    const result = await executeMutation({ input })
    if (result.error) {
      throw result.error
    }
  }
}

export const updateOrganizationResolver: UpdateResolver<
  UpdateOrganizationMutation,
  UpdateOrganizationMutationVariables
> = (result, args, cache, info) => {
  // not doing anything
}
