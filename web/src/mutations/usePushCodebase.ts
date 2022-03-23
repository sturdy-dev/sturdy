import { gql, useMutation } from '@urql/vue'
import type { Ref } from 'vue'
import type { PushCodebaseInput } from '../__generated__/types'
import type { DeepMaybeRef } from '@vueuse/core'
import type {
  PushCodebaseMutation,
  PushCodebaseMutationVariables,
} from './__generated__/usePushCodebase'
import type { UpdateResolver } from '@urql/exchange-graphcache'

const PUSH_CODEBASE = gql<PushCodebaseMutation, DeepMaybeRef<PushCodebaseMutationVariables>>`
  mutation PushCodebase($input: PushCodebaseInput!) {
    pushCodebase(input: $input) {
      id
    }
  }
`

export function usePushCodebase(): {
  mutating: Ref<boolean>
  pushCodebase(input: DeepMaybeRef<PushCodebaseInput>): Promise<void>
} {
  const { executeMutation, fetching: mutating } = useMutation(PUSH_CODEBASE)

  return {
    mutating,
    async pushCodebase(input: DeepMaybeRef<PushCodebaseInput>) {
      const result = await executeMutation({ input })
      if (result.error) {
        throw result.error
      }
    },
  }
}

export const pushCodebaseUpdateResolver: UpdateResolver<
  PushCodebaseMutation,
  PushCodebaseMutationVariables
> = (result, args, cache, info) => {
  // add if needed
}
