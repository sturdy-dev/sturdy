import { gql, useMutation } from '@urql/vue'
import type { CreateOrganizationInput } from '../__generated__/types'
import type { DeepMaybeRef } from '@vueuse/core'
import type { UpdateResolver } from '@urql/exchange-graphcache'
import type {
  AllOrganizationsQuery,
  AllOrganizationsQueryVariables,
  CreateOrganizationMutation,
  CreateOrganizationMutationVariables,
} from './__generated__/useCreateOrganization'

const CREATE_ORGANIZATION = gql`
  mutation CreateOrganization($input: CreateOrganizationInput!) {
    createOrganization(input: $input) {
      id
      shortID
      name
    }
  }
`

const ALL_ORGANIZATIONS = gql`
  query AllOrganizations {
    organizations {
      id
    }
  }
`

export function useCreateOrganization(): (
  input: DeepMaybeRef<CreateOrganizationInput>
) => Promise<CreateOrganizationMutation> {
  const { executeMutation } = useMutation<
    CreateOrganizationMutation,
    DeepMaybeRef<CreateOrganizationMutationVariables>
  >(CREATE_ORGANIZATION)
  return async (input): Promise<CreateOrganizationMutation> => {
    const result = await executeMutation({ input })
    if (result.error) throw result.error
    if (!result.data) throw new Error('No data returned')
    return result.data
  }
}

export const createOrganizationUpdateResolver: UpdateResolver<
  CreateOrganizationMutation,
  CreateOrganizationMutationVariables
> = (result, args, cache, info) => {
  // Add codebase to list of codebases in organization

  cache.updateQuery<AllOrganizationsQuery, AllOrganizationsQueryVariables>(
    {
      query: ALL_ORGANIZATIONS,
    },
    (data) => {
      if (data && !data.organizations.some((o) => o.id === result.createOrganization.id)) {
        data.organizations.push(result.createOrganization)
      }
      return data
    }
  )
}
