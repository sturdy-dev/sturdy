import { gql, useMutation } from '@urql/vue'
import type { UpdateResolver } from '@urql/exchange-graphcache'
import type { DeepMaybeRef } from '@vueuse/core'
import type { CreateOrUpdateBuildkiteIntegrationInput } from '../__generated__/types'
import type {
  CreateOrUpdateBuildkiteIntegrationMutation,
  CreateOrUpdateBuildkiteIntegrationMutationVariables,
} from './__generated__/useCreateOrUpdateBuildkiteIntegration'

const CREATE_BUILDKITE_INTEGRATION = gql`
  mutation CreateOrUpdateBuildkiteIntegration($input: CreateOrUpdateBuildkiteIntegrationInput!) {
    createOrUpdateBuildkiteIntegration(input: $input) {
      id
      ... on BuildkiteIntegration {
        id
        configuration {
          id
          organizationName
          pipelineName
          apiToken
          webhookSecret
        }
      }
    }
  }
`

export function useCreateOrUpdateBuildkiteIntegration(): (
  input: DeepMaybeRef<CreateOrUpdateBuildkiteIntegrationInput>
) => Promise<CreateOrUpdateBuildkiteIntegrationMutation> {
  const { executeMutation } = useMutation<
    CreateOrUpdateBuildkiteIntegrationMutation,
    DeepMaybeRef<CreateOrUpdateBuildkiteIntegrationMutationVariables>
  >(CREATE_BUILDKITE_INTEGRATION)
  return async (input) => {
    const result = await executeMutation({ input })
    if (result.error) throw result.error
    if (result.data) {
      return result.data
    }
    throw new Error('unexpected result')
  }
}

export const createOrUpdateBuildkiteIntegrationUpdateResolver: UpdateResolver<
  CreateOrUpdateBuildkiteIntegrationMutation,
  CreateOrUpdateBuildkiteIntegrationMutationVariables
> = (result, args, cache, info) => {
  // update cache manually if needed
}
