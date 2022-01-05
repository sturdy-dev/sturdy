import { DeepMaybeRef } from '@vueuse/core'
import { gql, useMutation } from '@urql/vue'
import { DismissSuggestionInput } from '../__generated__/types'
import {
  DismissSuggestionMutation,
  DismissSuggestionMutationVariables,
} from './__generated__/useDismissSuggestion'

const DISMISS_SUGGESTION_HUNKS = gql`
  mutation DismissSuggestion($input: DismissSuggestionInput!) {
    dismissSuggestion(input: $input) {
      id
      dismissedAt
    }
  }
`
export function useDismissSuggestion(): (
  input: DeepMaybeRef<DismissSuggestionInput>
) => Promise<void> {
  const { executeMutation } = useMutation<
    DismissSuggestionMutation,
    DeepMaybeRef<DismissSuggestionMutationVariables>
  >(DISMISS_SUGGESTION_HUNKS)
  return async (input) => {
    const result = await executeMutation({ input })
    if (result.error) throw result.error
  }
}
