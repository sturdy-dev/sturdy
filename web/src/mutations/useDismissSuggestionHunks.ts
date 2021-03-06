import type { DeepMaybeRef } from '@vueuse/core'
import { gql, useMutation } from '@urql/vue'
import type { DismissSuggestionHunksInput } from '../__generated__/types'
import type {
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
          _id
          hunkID
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
