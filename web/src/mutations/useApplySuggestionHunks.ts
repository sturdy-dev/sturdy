import { DeepMaybeRef } from '@vueuse/core'
import { gql, useMutation } from '@urql/vue'
import { ApplySuggestionHunksInput } from '../__generated__/types'
import {
  ApplySuggestionHunksMutation,
  ApplySuggestionHunksMutationVariables,
} from './__generated__/useApplySuggestionHunks'

const APPLY_SUGGESTION = gql`
  mutation ApplySuggestionHunks($input: ApplySuggestionHunksInput!) {
    applySuggestionHunks(input: $input) {
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
export function useApplySuggestionHunks(): (
  input: DeepMaybeRef<ApplySuggestionHunksInput>
) => Promise<void> {
  const { executeMutation } = useMutation<
    ApplySuggestionHunksMutation,
    DeepMaybeRef<ApplySuggestionHunksMutationVariables>
  >(APPLY_SUGGESTION)
  return async (input) => {
    const result = await executeMutation({ input })
    if (result.error) throw result.error
  }
}
