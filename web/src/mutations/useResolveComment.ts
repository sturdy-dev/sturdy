import type { UpdateResolver } from '@urql/exchange-graphcache'
import gql from 'graphql-tag'
import { useMutation } from '@urql/vue'
import type {
  ResolveCommentMutation,
  ResolveCommentMutationVariables,
} from './__generated__/useResolveComment'
import type { DeepMaybeRef, MaybeRef } from '@vueuse/core'

const RESOLVE_COMMENT = gql`
  mutation ResolveComment($commentID: ID!) {
    resolveComment(id: $commentID) {
      id
      ... on TopComment {
        resolved
        resolvedBy {
          id
          name
        }
      }
    }
  }
`

export function useResolveComment(): (commentID: MaybeRef<string>) => Promise<void> {
  const { executeMutation } = useMutation<
    ResolveCommentMutation,
    DeepMaybeRef<ResolveCommentMutationVariables>
  >(RESOLVE_COMMENT)

  return async (commentID) => {
    const result = await executeMutation({ commentID })

    if (result.error) {
      throw result.error
    }

    console.log('resolve comment', result)
  }
}

export const resolveCommentUpdateResolver: UpdateResolver = (parent, args, cache, info) => {
  // Update cache manually if needed
}
