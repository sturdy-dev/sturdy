import { UpdateResolver } from '@urql/exchange-graphcache'
import { useMutation } from '@urql/vue'
import gql from 'graphql-tag'
import { UpdateCommentInput } from '../__generated__/types'
import { DeepMaybeRef } from '@vueuse/core'
import {
  UpdateCommentMutation,
  UpdateCommentMutationVariables,
} from './__generated__/useUpdateComment'

const UPDATE_COMMENT = gql`
  mutation UpdateComment($input: UpdateCommentInput!) {
    updateComment(input: $input) {
      id
      message
    }
  }
`

export function useUpdateComment(): (input: DeepMaybeRef<UpdateCommentInput>) => Promise<void> {
  const { executeMutation } = useMutation<
    UpdateCommentMutation,
    DeepMaybeRef<UpdateCommentMutationVariables>
  >(UPDATE_COMMENT)
  return async (input) => {
    const result = await executeMutation({ input })

    if (result.error) {
      throw result.error
    }

    console.log('update comment', result)
  }
}

export const updateCommentUpdateResolver: UpdateResolver = (parent, args, cache, info) => {
  // Update cache manually if needed
}
