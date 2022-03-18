import { gql, useSubscription } from '@urql/vue'
import { DeepMaybeRef, MaybeRef } from '@vueuse/core'
import {
  UpdatedWorkspaceAuthorViewQuery,
  UpdatedWorkspaceAuthorViewQueryVariables,
  UpdatedWorkspaceSubscription,
  UpdatedWorkspaceSubscriptionVariables,
} from './__generated__/useUpdatedWorkspace'
import { UpdateResolver } from '@urql/exchange-graphcache'

const UPDATED_WORKSPACE = gql`
  subscription UpdatedWorkspace($shortCodebaseID: ID, $workspaceID: ID) {
    updatedWorkspace(shortCodebaseID: $shortCodebaseID, workspaceID: $workspaceID) {
      id
      name
      lastLandedAt
      conflicts
      updatedAt
      upToDateWithTrunk
      lastActivityAt
      draftDescription

      diffsCount
      commentsCount

      archivedAt
      unarchivedAt

      author {
        id
      }

      view {
        id
        mountPath
        shortMountPath
        mountHostname
        lastUsedAt
      }

      # Necessary to add new workspaces to the list of workspaces in a codebase
      codebase {
        id
        workspaces {
          id
        }
      }

      gitHubPullRequest {
        id
        pullRequestNumber
        open
        merged
        state
      }

      reviews {
        id
        grade
        createdAt
        dismissedAt
        author {
          id
          name
          avatarUrl
        }
      }

      headChange {
        id
        title
        trunkCommitID
        createdAt
        author {
          id
          name
          avatarUrl
        }
      }
    }
  }
`

export function useUpdatedWorkspaceByCodebase(
  shortCodebaseID: MaybeRef<string>,
  opts?: { pause?: MaybeRef<boolean> }
) {
  return useUpdatedWorkspaceSubscription({ shortCodebaseID: shortCodebaseID }, opts)
}

export function useUpdatedWorkspace(
  workspaceID: MaybeRef<string>,
  opts?: { pause?: MaybeRef<boolean> }
) {
  return useUpdatedWorkspaceSubscription({ workspaceID }, opts)
}

function useUpdatedWorkspaceSubscription(
  variables: DeepMaybeRef<UpdatedWorkspaceSubscriptionVariables>,
  { pause = false }: { pause?: MaybeRef<boolean> } = {}
) {
  useSubscription<
    UpdatedWorkspaceSubscription,
    UpdatedWorkspaceSubscription,
    DeepMaybeRef<UpdatedWorkspaceSubscriptionVariables>
  >({
    query: UPDATED_WORKSPACE,
    variables,
    pause,
  })
}

const UPDATED_WORKSPACE_AUTHOR_VIEW = gql`
  query UpdatedWorkspaceAuthorView {
    user {
      id
      views {
        id
        workspace {
          id
        }
      }
    }
  }
`

export const updatedWorkspaceUpdateResolver: UpdateResolver<
  UpdatedWorkspaceSubscription,
  UpdatedWorkspaceSubscriptionVariables
> = (parent, args, cache, info) => {
  // When a workspace is updated to have a view, add it to the author's
  // list of views (if they're the current user)
  cache.updateQuery<UpdatedWorkspaceAuthorViewQuery, UpdatedWorkspaceAuthorViewQueryVariables>(
    {
      query: UPDATED_WORKSPACE_AUTHOR_VIEW,
    },
    (data) => {
      if (
        data?.user.id === parent.updatedWorkspace.author.id &&
        parent.updatedWorkspace.view &&
        !data.user.views.some((view) => view.workspace?.id === parent.updatedWorkspace.id)
      ) {
        data.user.views.push(parent.updatedWorkspace.view)
      }
      return data
    }
  )
}
