import type { UpdateResolver } from '@urql/exchange-graphcache'
import type { DeepMaybeRef } from '@vueuse/core'
import gql from 'graphql-tag'
import { useMutation } from '@urql/vue'
import type { VerifyEmailInput } from '../__generated__/types'
import type {
  VerifyEmailMutation,
  VerifyEmailMutationVariables,
} from './__generated__/useVerifyEmail'

const VERIFY_EMAIL = gql`
  mutation VerifyEmail($input: VerifyEmailInput!) {
    verifyEmail(input: $input) {
      id
      emailVerified
    }
  }
`

export function useVerifyEmail(): (input: DeepMaybeRef<VerifyEmailInput>) => Promise<void> {
  const { executeMutation } = useMutation<
    VerifyEmailMutation,
    DeepMaybeRef<VerifyEmailMutationVariables>
  >(VERIFY_EMAIL)
  return async (input) => {
    const result = await executeMutation({ input })
    if (result.error) {
      throw result.error
    }
  }
}

export const verifyEmailResolver: UpdateResolver = (parent, args, cache, info) => {
  // update cache manually if needed
}
