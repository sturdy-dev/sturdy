import { DeepMaybeRef } from '@vueuse/core'
import gql from 'graphql-tag'
import { useMutation } from '@urql/vue'
import { DeleteIntegrationInput } from '../__generated__/types'
import {
  DeleteIntegrationMutation,
  DeleteIntegrationMutationVariables,
} from './__generated__/useDeleteIntegration'

const DELETE_INTEGRATION = gql`
  mutation DeleteIntegration($input: DeleteIntegrationInput!) {
    deleteIntegration(input: $input) {
      id
      deletedAt
    }
  }
`

export function useDeleteIntegration(): (
  input: DeepMaybeRef<DeleteIntegrationInput>
) => Promise<void> {
  const { executeMutation } = useMutation<
    DeleteIntegrationMutation,
    DeepMaybeRef<DeleteIntegrationMutationVariables>
  >(DELETE_INTEGRATION)
  return async (input) => {
    const result = await executeMutation({ input })
    if (result.error) {
      throw result.error
    }
  }
}
