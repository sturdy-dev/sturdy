import { gql, useMutation } from '@urql/vue'
import type { Ref } from 'vue'
import type { PullCodebaseInput } from '../__generated__/types'
import type { DeepMaybeRef } from '@vueuse/core'
import type {
  PullCodebaseMutation,
  PullCodebaseMutationVariables,
} from './__generated__/usePullCodebase'
import type { UpdateResolver } from '@urql/exchange-graphcache'

const PULL_CODEBASE = gql<PullCodebaseMutation, DeepMaybeRef<PullCodebaseMutationVariables>>`
  mutation PullCodebase($input: PullCodebaseInput!) {
    pullCodebase(input: $input) {
      id
      changes(input: { limit: 1 }) {
        id
      }
    }
  }
`

export function usePullCodebase(): {
  mutating: Ref<boolean>
  pullCodebase(input: DeepMaybeRef<PullCodebaseInput>): Promise<void>
} {
  const { executeMutation, fetching: mutating } = useMutation(PULL_CODEBASE)

  return {
    mutating,
    async pullCodebase(input: DeepMaybeRef<PullCodebaseInput>) {
      const result = await executeMutation({ input })
      if (result.error) {
        throw result.error
      }
    },
  }
}

export const pullCodebaseUpdateResolver: UpdateResolver<
  PullCodebaseMutation,
  PullCodebaseMutationVariables
> = (result, args, cache, info) => {
  if (!result) {
    return
  }
  if (result.pullCodebase && result.pullCodebase.__typename) {
    const key = cache.keyOfEntity({
      __typename: result.pullCodebase.__typename,
      id: result.pullCodebase.id,
    })
    cache
      .inspectFields(key)
      .filter(({ fieldName }) => fieldName === 'changes')
      .forEach((f) => {
        cache.invalidate(key, f.fieldKey)
      })
  }
}
