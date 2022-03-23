import { gql, useSubscription } from '@urql/vue'
import type { DeepMaybeRef } from '@vueuse/core'
import type { UpdateResolver } from '@urql/exchange-graphcache'
import type {
  UpdatedReviewsSubscription,
  UpdatedReviewsSubscriptionVariables,
  WorkspaceReviewsQuery,
  WorkspaceReviewsQueryVariables,
} from './__generated__/useUpdatedReviews'

const UPDATED_REVIEWS = gql`
  subscription UpdatedReviews {
    updatedReviews {
      id
      author {
        id
        name
        avatarUrl
      }
      grade
      createdAt
      dismissedAt
      isReplaced
      requestedBy {
        id
        name
        avatarUrl
      }
      workspace {
        id
      }
    }
  }
`

const WORKSPACE_REVIEWS = gql`
  query WorkspaceReviews($workspaceId: ID!) {
    workspace(id: $workspaceId) {
      id
      reviews {
        id
        grade
        author {
          id
        }
      }
    }
  }
`

export function useUpdatedReviews() {
  useSubscription<UpdatedReviewsSubscription, DeepMaybeRef<UpdatedReviewsSubscriptionVariables>>({
    query: UPDATED_REVIEWS,
  })
}

export const updatedReviewsResolver: UpdateResolver<
  UpdatedReviewsSubscription,
  UpdatedReviewsSubscriptionVariables
> = (result, args, cache, info) => {
  const updatedReview = result.updatedReviews
  cache.updateQuery<WorkspaceReviewsQuery, WorkspaceReviewsQueryVariables>(
    {
      query: WORKSPACE_REVIEWS,
      variables: { workspaceId: updatedReview.workspace.id },
    },
    (data) => {
      if (!data) {
        return {
          workspace: {
            __typename: 'Workspace',
            id: updatedReview.workspace.id,
            reviews: [updatedReview],
          },
        }
      }
      data.workspace.reviews = [
        updatedReview,
        ...data.workspace.reviews.filter((review) => review.id !== updatedReview.id),
      ]
      return data
    }
  )
}
