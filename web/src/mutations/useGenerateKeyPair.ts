import { gql, useMutation } from '@urql/vue'
import type { DeepMaybeRef } from '@vueuse/core'
import type { GenerateKeyPairInput } from '../__generated__/types'
import type {
  GenerateKeyPairMutation,
  GenerateKeyPairMutationVariables,
} from './__generated__/useGenerateKeyPair'
import type { Ref } from 'vue'

const GENERATE_KEYPAIR = gql`
  mutation GenerateKeyPair($input: GenerateKeyPairInput!) {
    generateKeyPair(input: $input) {
      id
      publicKey
    }
  }
`

export function useGenerateKeyPair(): {
  mutating: Ref<boolean>
  generateKeyPair(input: DeepMaybeRef<GenerateKeyPairInput>): Promise<GenerateKeyPairMutation>
} {
  const { executeMutation, fetching: mutating } = useMutation<
    GenerateKeyPairMutation,
    DeepMaybeRef<GenerateKeyPairMutationVariables>
  >(GENERATE_KEYPAIR)

  return {
    mutating,
    async generateKeyPair(
      input: DeepMaybeRef<GenerateKeyPairInput>
    ): Promise<GenerateKeyPairMutation> {
      const result = await executeMutation({ input })
      if (result.error) {
        throw result.error
      }
      if (result.data) {
        return result.data
      }
      throw new Error('unexpected result')
    },
  }
}
