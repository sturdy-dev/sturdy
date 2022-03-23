import { gql, useMutation } from '@urql/vue'
import type { UpdateResolver } from '@urql/exchange-graphcache'
import type {
  CreateServiceTokenMutation,
  CreateServiceTokenMutationVariables,
} from './__generated__/useCreateServiceToken'
import type { CreateServiceTokenInput, ServiceToken } from '../__generated__/types'
import type { DeepMaybeRef } from '@vueuse/core'

const CREATE_SERVICE_TOKEN = gql`
  mutation CreateServiceToken($input: CreateServiceTokenInput!) {
    createServiceToken(input: $input) {
      id
      name
      createdAt
      lastUsedAt
      token
    }
  }
`

export function useCreateServiceToken(): (
  input: DeepMaybeRef<CreateServiceTokenInput>
) => Promise<ServiceToken> {
  const { executeMutation } = useMutation<
    CreateServiceTokenMutation,
    DeepMaybeRef<CreateServiceTokenMutationVariables>
  >(CREATE_SERVICE_TOKEN)
  return async (input): Promise<ServiceToken> => {
    const result = await executeMutation({ input })
    if (result.error) throw result.error
    if (!result.data) throw new Error('No data returned')
    return result.data.createServiceToken
  }
}

export const createServiceTokenResolver: UpdateResolver<
  CreateServiceTokenMutation,
  CreateServiceTokenMutationVariables
> = (result, args, cache, info) => {
  // update cache manually if needed
}
