import type { UpdateResolver } from '@urql/exchange-graphcache'
import gql from 'graphql-tag'
import { useMutation } from '@urql/vue'
import type {
  DeleteCommentMutation,
  DeleteCommentMutationVariables,
} from './__generated__/useDeleteComment'
import type { DeepMaybeRef, MaybeRef } from '@vueuse/core'

const DELETE_COMMENT = gql`
  mutation DeleteComment($commentID: ID!) {
    deleteComment(id: $commentID) {
      id
      deletedAt
    }
  }
`

export function useDeleteComment(): (commentID: MaybeRef<string>) => Promise<void> {
  const { executeMutation } = useMutation<
    DeleteCommentMutation,
    DeepMaybeRef<DeleteCommentMutationVariables>
  >(DELETE_COMMENT)

  return async (commentID) => {
    const result = await executeMutation({ commentID })

    if (result.error) {
      throw result.error
    }

    console.log('delete comment', result)
  }
}

export const deleteCommentUpdateResolver: UpdateResolver = (parent, args, cache, info) => {
  // Update cache manually if needed
}
