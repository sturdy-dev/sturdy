import { DeepMaybeRef } from '@vueuse/core'
import { gql, useMutation } from '@urql/vue'
import { DismissSuggestionHunksInput } from '../__generated__/types'
import {
  DismissSuggestionHunksMutation,
  DismissSuggestionHunksMutationVariables,
} from './__generated__/useDismissSuggestionHunks'

const DISMISS_SUGGESTION_HUNKS = gql`
  mutation DismissSuggestionHunks($input: DismissSuggestionHunksInput!) {
    dismissSuggestionHunks(input: $input) {
      id

      diffs {
        id

        origName
        newName
        preferredName

        isDeleted
        isNew
        isMoved

        hunks {
          id
          patch

          isOutdated
          isApplied
          isDismissed
        }
      }
    }
  }
`
export function useDismissSuggestionHunks(): (
  input: DeepMaybeRef<DismissSuggestionHunksInput>
) => Promise<void> {
  const { executeMutation } = useMutation<
    DismissSuggestionHunksMutation,
    DeepMaybeRef<DismissSuggestionHunksMutationVariables>
  >(DISMISS_SUGGESTION_HUNKS)
  return async (input) => {
    const result = await executeMutation({ input })
    if (result.error) throw result.error
  }
}
