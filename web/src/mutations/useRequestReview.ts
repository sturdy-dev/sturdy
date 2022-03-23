import type { RequestReviewInput } from '../__generated__/types'
import { gql, useMutation } from '@urql/vue'
import type {
  RequestReviewMutation,
  RequestReviewMutationVariables,
  RequestReviewWorkspaceReviewsQuery,
  RequestReviewWorkspaceReviewsQueryVariables,
} from './__generated__/useRequestReview'
import type { UpdateResolver } from '@urql/exchange-graphcache'
import type { DeepMaybeRef } from '@vueuse/core'

const REQUEST_REVIEW = gql`
  mutation RequestReview($input: RequestReviewInput!) {
    requestReview(input: $input) {
      id
      grade
      author {
        id
        name
        avatarUrl
      }
    }
  }
`

const REQUEST_REVIEW_WORKSPACE_REVIEWS = gql`
  query RequestReviewWorkspaceReviews($workspaceID: ID!) {
    workspace(id: $workspaceID) {
      id
      reviews {
        id
      }
    }
  }
`

export function useRequestReview(): (input: DeepMaybeRef<RequestReviewInput>) => Promise<void> {
  const { executeMutation } = useMutation<
    RequestReviewMutation,
    DeepMaybeRef<RequestReviewMutationVariables>
  >(REQUEST_REVIEW)

  return async (input) => {
    const result = await executeMutation({ input })

    if (result.error) {
      throw result.error
    }

    console.log('requestReviewResult', result)
  }
}

export const requestReviewUpdateResolver: UpdateResolver<
  RequestReviewMutation,
  RequestReviewMutationVariables
> = (result, args, cache, info) => {
  if (args.input.workspaceID) {
    cache.updateQuery<
      RequestReviewWorkspaceReviewsQuery,
      RequestReviewWorkspaceReviewsQueryVariables
    >(
      {
        query: REQUEST_REVIEW_WORKSPACE_REVIEWS,
        variables: { workspaceID: args.input.workspaceID },
      },
      (data) => {
        // Add review if not exists
        if (data && !data.workspace.reviews.some((c) => c.id === result.requestReview.id)) {
          data.workspace.reviews.push(result.requestReview)
        }
        return data
      }
    )
  }
}
