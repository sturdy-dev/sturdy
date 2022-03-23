import type { DeepMaybeRef, MaybeRef } from '@vueuse/core'
import { gql, useSubscription } from '@urql/vue'
import type { UpdateResolver } from '@urql/exchange-graphcache'
import type {
  CreatePresenceCodebaseWorkspacePresenceQuery,
  CreatePresenceCodebaseWorkspacePresenceQueryVariables,
  CreatePresenceWorkspacePresenceQuery,
  CreatePresenceWorkspacePresenceQueryVariables,
  UpdatedWorkspacePresenceSubscription,
  UpdatedWorkspacePresenceSubscriptionVariables,
} from './__generated__/useUpdatedWorkspacePresence'

const UPDATED_WORKSPACE_PRESENCE = gql`
  subscription UpdatedWorkspacePresence($workspaceID: ID) {
    updatedWorkspacePresence(workspaceID: $workspaceID) {
      id
      state
      lastActiveAt
      workspace {
        id
        codebase {
          id
        }
      }
    }
  }
`

const WORKSPACE_PRESENCE = gql`
  query CreatePresenceWorkspacePresence($workspaceID: ID!) {
    workspace(id: $workspaceID) {
      id
      presence {
        id
      }
    }
  }
`

const ALL_CODEBASES_ALL_WORKSPACES_PRESENCE = gql`
  query CreatePresenceCodebaseWorkspacePresence {
    codebases {
      id
      workspaces {
        id
        presence {
          id
        }
      }
    }
  }
`

export function useUpdatedWorkspacePresence(
  workspaceID: MaybeRef<string> | MaybeRef<undefined> | MaybeRef<null>,
  { pause = false }: { pause?: MaybeRef<boolean> } = {}
) {
  useSubscription<
    UpdatedWorkspacePresenceSubscription,
    UpdatedWorkspacePresenceSubscription,
    DeepMaybeRef<UpdatedWorkspacePresenceSubscriptionVariables>
  >({
    query: UPDATED_WORKSPACE_PRESENCE,
    variables: { workspaceID },
    pause,
  })
}

export const updatedWorkspacePresenceResolver: UpdateResolver<
  UpdatedWorkspacePresenceSubscription,
  UpdatedWorkspacePresenceSubscriptionVariables
> = (result, args, cache, info) => {
  if (args.workspaceID) {
    cache.updateQuery<
      CreatePresenceWorkspacePresenceQuery,
      CreatePresenceWorkspacePresenceQueryVariables
    >(
      {
        query: WORKSPACE_PRESENCE,
        variables: { workspaceID: args.workspaceID },
      },
      (data) => {
        // Add presence if not exists
        if (
          data &&
          !data.workspace.presence.some((c) => c?.id === result.updatedWorkspacePresence.id)
        ) {
          data.workspace.presence.push(result.updatedWorkspacePresence)
        }
        return data
      }
    )
  }

  cache.updateQuery<
    CreatePresenceCodebaseWorkspacePresenceQuery,
    CreatePresenceCodebaseWorkspacePresenceQueryVariables
  >(
    {
      query: ALL_CODEBASES_ALL_WORKSPACES_PRESENCE,
    },
    (data) => {
      if (data?.codebases) {
        // Add presence if not exists
        for (const cb of data.codebases) {
          if (cb.id !== result.updatedWorkspacePresence.workspace.codebase.id) {
            continue
          }

          for (const ws of cb.workspaces) {
            if (ws.id !== result.updatedWorkspacePresence.workspace.id) {
              continue
            }
            // Add to presence if not exists
            if (!ws.presence.some((c) => c?.id === result.updatedWorkspacePresence.id)) {
              ws.presence.push(result.updatedWorkspacePresence)
            }
          }
        }
      }

      return data
    }
  )
}
