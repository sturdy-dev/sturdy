import { gql, useMutation } from '@urql/vue'
import { UpdateInstallationInput } from '../__generated__/types'
import { DeepMaybeRef } from '@vueuse/core'
import { UpdateResolver } from '@urql/exchange-graphcache'
import {
  CreateOrganizationMutation,
  CreateOrganizationMutationVariables,
} from './__generated__/useCreateOrganization'
import {
  UpdateInstallationMutation,
  UpdateInstallationMutationVariables,
} from './__generated__/useUpdateInstallation'

const UPDATE_INSTALLATION = gql`
  mutation UpdateInstallation($input: UpdateInstallationInput!) {
    updateInstallation(input: $input) {
      id
      license {
        id
      }
    }
  }
`

export function useUpdateInstallation(): (
  input: DeepMaybeRef<UpdateInstallationInput>
) => Promise<UpdateInstallationMutation> {
  const { executeMutation } = useMutation<
    UpdateInstallationMutation,
    DeepMaybeRef<UpdateInstallationMutationVariables>
  >(UPDATE_INSTALLATION)
  return async (input): Promise<UpdateInstallationMutation> => {
    const result = await executeMutation({ input })
    if (result.error) throw result.error
    if (!result.data) throw new Error('No data returned')
    return result.data
  }
}

export const updateInstallationUpdateResolver: UpdateResolver<
  CreateOrganizationMutation,
  CreateOrganizationMutationVariables
> = (result, args, cache, info) => {
  // nothing
}
