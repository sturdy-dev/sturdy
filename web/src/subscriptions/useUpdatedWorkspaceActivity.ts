import { gql, useSubscription } from '@urql/vue'
import type { Entity, UpdateResolver } from '@urql/exchange-graphcache'
import type {
  UpdatedWorkspaceActivitySubscription,
  UpdatedWorkspaceActivitySubscriptionVariables,
} from './__generated__/useUpdatedWorkspaceActivity'
import type { DeepMaybeRef } from '@vueuse/core'

const UPDATED_WORKSPACE_ACTIVITY = gql`
  subscription UpdatedWorkspaceActivity {
    updatedWorkspaceActivity {
      id

      isRead
      workspace {
        id
      }

      change {
        id
      }

      createdAt
      author {
        id
        name
        avatarUrl
      }
      ... on WorkspaceCommentActivity {
        comment {
          id
          message
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
          }
          ... on ReplyComment {
            parent {
              id
              message
              author {
                id
                name
              }
            }
          }
        }
      }
      ... on WorkspaceCreatedChangeActivity {
        change {
          id
          title
          trunkCommitID
        }
      }
      ... on WorkspaceRequestedReviewActivity {
        review {
          id
          grade
          createdAt
          dismissedAt
          isReplaced
          author {
            id
            name
            avatarUrl
          }
        }
      }
      ... on WorkspaceReviewedActivity {
        review {
          id
          grade
          createdAt
          isReplaced
          dismissedAt
          author {
            id
            name
            avatarUrl
          }
        }
      }
    }
  }
`

export function useUpdatedWorkspaceActivity() {
  useSubscription<
    UpdatedWorkspaceActivitySubscription,
    UpdatedWorkspaceActivitySubscription,
    DeepMaybeRef<UpdatedWorkspaceActivitySubscriptionVariables>
  >({
    query: UPDATED_WORKSPACE_ACTIVITY,
  })
}

export const updatedWorkspaceActivityUpdateResolver: UpdateResolver<
  UpdatedWorkspaceActivitySubscription,
  UpdatedWorkspaceActivitySubscriptionVariables
> = (result, args, cache, info) => {
  if (result && result.updatedWorkspaceActivity.workspace) {
    const activityList = cache.resolve(
      { __typename: 'Workspace', id: result.updatedWorkspaceActivity.workspace.id },
      'activity'
    ) as Array<string>
    const selfKey = cache.keyOfEntity(result.updatedWorkspaceActivity as Entity)

    if (activityList && selfKey) {
      if (!activityList.includes(selfKey)) {
        activityList.push(selfKey)
        cache.link(
          { __typename: 'Workspace', id: result.updatedWorkspaceActivity.workspace.id },
          'activity',
          activityList
        )
      }
    }
  }

  if (result && result.updatedWorkspaceActivity.change) {
    const activityList = cache.resolve(
      { __typename: 'Change', id: result.updatedWorkspaceActivity.change.id },
      'activity'
    ) as Array<string>
    const selfKey = cache.keyOfEntity(result.updatedWorkspaceActivity as Entity)

    if (activityList && selfKey) {
      if (!activityList.includes(selfKey)) {
        activityList.push(selfKey)
        cache.link(
          { __typename: 'Change', id: result.updatedWorkspaceActivity.change.id },
          'activity',
          activityList
        )
      }
    }
  }
}
