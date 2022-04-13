import { gql, useSubscription } from '@urql/vue'
import type {
  UpdatedCommentSubscription,
  UpdatedCommentSubscriptionVariables,
  UpdatedCommentWorkspaceCommentsQuery,
  UpdatedCommentWorkspaceCommentsQueryVariables,
} from './__generated__/useUpdatedComment'
import type { DeepMaybeRef, MaybeRef } from '@vueuse/core'
import type { UpdateResolver } from '@urql/exchange-graphcache'

const UPDATED_COMMENT = gql`
  subscription UpdatedComment($workspaceID: ID!, $viewID: ID) {
    updatedComment(workspaceID: $workspaceID, viewID: $viewID) {
      id
      message
      deletedAt
      createdAt

      author {
        id
        name
        email
      }

      ... on TopComment {
        codeContext {
          id
          lineStart
          lineEnd
          lineIsNew
          context
          contextStartsAtLine
          path
        }
        workspace {
          id
        }
        replies {
          id
          message
        }
        resolved
      }

      ... on ReplyComment {
        parent {
          id
        }
      }
    }
  }
`

const WORKSPACE_COMMENTS = gql`
  query UpdatedCommentWorkspaceComments($workspaceID: ID!) {
    workspace(id: $workspaceID) {
      id
      comments {
        id
      }
    }
  }
`

export function useUpdatedComment(workspaceID: MaybeRef<string>, viewID?: MaybeRef<string>) {
  useSubscription<
    UpdatedCommentSubscription,
    UpdatedCommentSubscription,
    DeepMaybeRef<UpdatedCommentSubscriptionVariables>
  >({
    query: UPDATED_COMMENT,
    variables: { workspaceID, viewID },
  })
}

export const updatedCommentUpdateResolver: UpdateResolver<
  UpdatedCommentSubscription,
  UpdatedCommentSubscriptionVariables
> = (result, args, cache, info) => {
  if (
    result &&
    result.updatedComment.__typename === 'TopComment' &&
    result.updatedComment.workspace
  ) {
    cache.updateQuery<
      UpdatedCommentWorkspaceCommentsQuery,
      UpdatedCommentWorkspaceCommentsQueryVariables
    >(
      {
        query: WORKSPACE_COMMENTS,
        variables: { workspaceID: result.updatedComment.workspace.id },
      },
      (data) => {
        // Add comment if not exists
        if (
          data &&
          !data.workspace.comments.some((c) => c.id === result.updatedComment.id) &&
          result.updatedComment.__typename === 'TopComment'
        ) {
          data.workspace.comments.push(result.updatedComment)
        }
        return data
      }
    )
  }

  // Add replies to top comments
  if (
    result &&
    result.updatedComment.__typename === 'ReplyComment' &&
    result.updatedComment.parent
  ) {
    const repliesList = cache.resolve(
      { __typename: 'TopComment', id: result.updatedComment.parent.id },
      'replies'
    ) as Array<string>

    const selfKey = cache.keyOfEntity(result.updatedComment)

    if (repliesList && selfKey) {
      if (!repliesList.includes(selfKey)) {
        repliesList.push(selfKey)
        cache.link({ __typename: 'TopComment', id: args.input.inReplyTo }, 'replies', repliesList)
      }
    }
  }
}
