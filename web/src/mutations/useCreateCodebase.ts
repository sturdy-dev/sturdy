import { gql, useMutation } from '@urql/vue'
import { Codebase, CreateCodebaseInput } from '../__generated__/types'
import { DeepMaybeRef } from '@vueuse/core'
import {
  CreateCodebaseMutation,
  CreateCodebaseMutationVariables,
  CreateCodebaseOrganizationCodebasesQuery,
  CreateCodebaseOrganizationCodebasesQueryVariables,
} from './__generated__/useCreateCodebase'
import { UpdateResolver } from '@urql/exchange-graphcache'

const CREATE_CODEBASE = gql`
  mutation CreateCodebase($input: CreateCodebaseInput!) {
    createCodebase(input: $input) {
      id
      shortID
      name
      organization {
        id
      }
    }
  }
`

const ORGANIZATION_CODEBASES = gql`
  query CreateCodebaseOrganizationCodebases($organizationID: ID!) {
    organization(id: $organizationID) {
      id
      codebases {
        id
      }
    }
  }
`

export function useCreateCodebase(): (
  input: DeepMaybeRef<CreateCodebaseInput>
) => Promise<CreateCodebaseMutation> {
  const { executeMutation } = useMutation<
    CreateCodebaseMutation,
    DeepMaybeRef<CreateCodebaseMutationVariables>
  >(CREATE_CODEBASE)

  return async (input): Promise<CreateCodebaseMutation> => {
    const result = await executeMutation({ input })
    if (result.error) throw result.error
    if (!result.data) throw new Error('No data returned')
    return result.data
  }
}

export const createCodebaseUpdateResolver: UpdateResolver<
  CreateCodebaseMutation,
  CreateCodebaseMutationVariables
> = (result, args, cache, info) => {
  // Add codebase to list of codebases in organization
  if (args.input.organizationID) {
    cache.updateQuery<
      CreateCodebaseOrganizationCodebasesQuery,
      CreateCodebaseOrganizationCodebasesQueryVariables
    >(
      {
        query: ORGANIZATION_CODEBASES,
        variables: { organizationID: args.input.organizationID },
      },
      (data) => {
        if (data && !data.organization.codebases.some((c) => c.id === result.createCodebase.id)) {
          data.organization.codebases.push(result.createCodebase)
        }
        return data
      }
    )
  }
}
