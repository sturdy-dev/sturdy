import { gql, useSubscription } from '@urql/vue'
import { UpdateResolver } from '@urql/exchange-graphcache'
import {
  UpdatedWorkspaceActivityLatestActivityQuery,
  UpdatedWorkspaceActivityLatestActivityQueryVariables,
  UpdatedWorkspaceActivitySubscription,
  UpdatedWorkspaceActivitySubscriptionVariables,
  UpdatedWorkspaceActivityWorkspaceAllActivityQuery,
  UpdatedWorkspaceActivityWorkspaceAllActivityQueryVariables,
} from './__generated__/useUpdatedWorkspaceActivity'
import { DeepMaybeRef } from '@vueuse/core'

const UPDATED_WORKSPACE_ACTIVITY = gql`
  subscription UpdatedWorkspaceActivity {
    updatedWorkspaceActivity {
      id

      isRead
      workspace {
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

const UPDATED_WORKSPACE_ACTIVITY_LATEST_ACTIVITY = gql`
  query UpdatedWorkspaceActivityLatestActivity {
    codebases {
      id
      workspaces {
        id
        activity(input: { unreadOnly: true, limit: 1 }) {
          id
          isRead
        }
      }
    }
  }
`

const UPDATED_WORKSPACE_ACTIVITY_WORKSPACE_ALL_ACTIVITY = gql`
  query UpdatedWorkspaceActivityWorkspaceAllActivity($id: ID!) {
    workspace(id: $id) {
      id
      activity {
        id
        isRead
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
  const updatedActivity = result.updatedWorkspaceActivity

  // Add as latest activity
  cache.updateQuery<
    UpdatedWorkspaceActivityLatestActivityQuery,
    UpdatedWorkspaceActivityLatestActivityQueryVariables
  >(
    {
      query: UPDATED_WORKSPACE_ACTIVITY_LATEST_ACTIVITY,
    },
    (data) => {
      if (!data) return data
      for (const cb of data.codebases) {
        for (const ws of cb.workspaces) {
          if (ws.id === result.updatedWorkspaceActivity.workspace.id) {
            ws.activity = [updatedActivity]
          }
        }
      }
      return data
    }
  )

  // Add to workspace list of activity
  cache.updateQuery<
    UpdatedWorkspaceActivityWorkspaceAllActivityQuery,
    UpdatedWorkspaceActivityWorkspaceAllActivityQueryVariables
  >(
    {
      query: UPDATED_WORKSPACE_ACTIVITY_WORKSPACE_ALL_ACTIVITY,
      variables: { id: updatedActivity.workspace.id },
    },
    (data) => {
      if (!data) {
        return {
          workspace: {
            __typename: 'Workspace',
            id: updatedActivity.workspace.id,
            activity: [updatedActivity],
          },
        }
      }
      if (!data) return data
      data.workspace.activity = [
        updatedActivity,
        ...data.workspace.activity.filter((activity) => activity.id != updatedActivity.id),
      ]
      return data
    }
  )
}
