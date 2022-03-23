import { gql, useMutation } from '@urql/vue'
import type { UpdateResolver } from '@urql/exchange-graphcache'
import type { CreateReviewInput } from '../__generated__/types'
import type {
  CreateOrUpdateReviewMutation,
  CreateOrUpdateReviewMutationVariables,
  CreateOrUpdateReviewWorkspaceReviewsQuery,
  CreateOrUpdateReviewWorkspaceReviewsQueryVariables,
} from './__generated__/useCreateOrUpdateReview'
import type { DeepMaybeRef } from '@vueuse/core'

const CREATE_OR_UPDATE_REVIEW = gql`
  mutation CreateOrUpdateReview($input: CreateReviewInput!) {
    createOrUpdateReview(input: $input) {
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

const CREATE_OR_UPDATE_REVIEW_WORKSPACE_REVIEWS = gql`
  query CreateOrUpdateReviewWorkspaceReviews($workspaceID: ID!) {
    workspace(id: $workspaceID) {
      id
      reviews {
        id
      }
    }
  }
`

export function useCreateOrUpdateReview(): (
  input: DeepMaybeRef<CreateReviewInput>
) => Promise<void> {
  const { executeMutation } = useMutation<
    CreateOrUpdateReviewMutation,
    DeepMaybeRef<CreateOrUpdateReviewMutationVariables>
  >(CREATE_OR_UPDATE_REVIEW)

  return async (input) => {
    const result = await executeMutation({ input })

    if (result.error) {
      throw result.error
    }

    console.log('createOrUpdateReview', result)
  }
}

export const createOrUpdateReviewUpdateResolver: UpdateResolver<
  CreateOrUpdateReviewMutation,
  CreateOrUpdateReviewMutationVariables
> = (result, args, cache, info) => {
  if (args.input.workspaceID) {
    cache.updateQuery<
      CreateOrUpdateReviewWorkspaceReviewsQuery,
      CreateOrUpdateReviewWorkspaceReviewsQueryVariables
    >(
      {
        query: CREATE_OR_UPDATE_REVIEW_WORKSPACE_REVIEWS,
        variables: { workspaceID: args.input.workspaceID },
      },
      (data) => {
        // Add review if not exists
        if (data && !data.workspace.reviews.some((c) => c.id === result.createOrUpdateReview.id)) {
          data.workspace.reviews.push(result.createOrUpdateReview)
        }
        return data
      }
    )
  }
}
