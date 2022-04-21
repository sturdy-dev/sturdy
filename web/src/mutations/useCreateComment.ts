import type { Entity, UpdateResolver } from '@urql/exchange-graphcache'
import type {
  CreateCommentChangeCommentsQuery,
  CreateCommentChangeCommentsQueryVariables,
  CreateCommentMutation,
  CreateCommentMutationVariables,
  CreateCommentWorkspaceCommentsQuery,
  CreateCommentWorkspaceCommentsQueryVariables,
} from './__generated__/useCreateComment'
import { gql, useMutation } from '@urql/vue'
import type { CreateCommentInput } from '../__generated__/types'
import type { DeepMaybeRef } from '@vueuse/core'

const CREATE_COMMENT = gql`
  mutation CreateComment($input: CreateCommentInput!) {
    createComment(input: $input) {
      id
      message
      createdAt
      deletedAt
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

        resolved

        replies {
          id
        }
      }
    }
  }
`

const WORKSPACE_COMMENTS = gql`
  query CreateCommentWorkspaceComments($workspaceID: ID!) {
    workspace(id: $workspaceID) {
      id
      comments {
        id
      }
    }
  }
`

const CHANGE_COMMENTS = gql`
  query CreateCommentChangeComments($changeID: ID!) {
    change(id: $changeID) {
      id
      comments {
        id
      }
    }
  }
`

export function useCreateComment(): (input: DeepMaybeRef<CreateCommentInput>) => Promise<void> {
  const { executeMutation } = useMutation<
    CreateCommentMutation,
    DeepMaybeRef<CreateCommentMutationVariables>
  >(CREATE_COMMENT)

  return async (input) => {
    const result = await executeMutation({ input })
    if (result.error) {
      throw result.error
    }
  }
}

export const createCommentUpdateResolver: UpdateResolver<
  CreateCommentMutation,
  CreateCommentMutationVariables
> = (result, args, cache, info) => {
  // Comments on a workspace
  if (args.input.workspaceID) {
    cache.updateQuery<
      CreateCommentWorkspaceCommentsQuery,
      CreateCommentWorkspaceCommentsQueryVariables
    >(
      {
        query: WORKSPACE_COMMENTS,
        variables: { workspaceID: args.input.workspaceID },
      },
      (data) => {
        // Add comment if not exists
        if (
          data &&
          result.createComment.__typename === 'TopComment' &&
          !data.workspace.comments.some((c) => c.id === result.createComment.id)
        ) {
          console.log('push comment to workspace')
          data.workspace.comments.push(result.createComment)
        }
        return data
      }
    )
  }

  // Comment on a change
  if (args.input.changeID) {
    cache.updateQuery<CreateCommentChangeCommentsQuery, CreateCommentChangeCommentsQueryVariables>(
      {
        query: CHANGE_COMMENTS,
        variables: { changeID: args.input.changeID },
      },
      (data) => {
        // Add comment if not exists
        if (
          data?.change &&
          result.createComment.__typename === 'TopComment' &&
          !data.change.comments.some((c) => c.id === result.createComment.id)
        ) {
          console.log('push comment to change')
          data.change.comments.push(result.createComment)
        }
        return data
      }
    )
  }

  // Add replies to top comments
  if (args.input.inReplyTo) {
    const repliesList = cache.resolve(
      { __typename: 'TopComment', id: args.input.inReplyTo },
      'replies'
    ) as Array<string>

    const selfKey = cache.keyOfEntity(result.createComment as Entity)

    if (repliesList && selfKey) {
      if (!repliesList.includes(selfKey)) {
        console.log("push comment's reply to top comment")
        repliesList.push(selfKey)
        cache.link({ __typename: 'TopComment', id: args.input.inReplyTo }, 'replies', repliesList)
      }
    }
  }
}
