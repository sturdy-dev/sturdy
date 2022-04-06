import { gql, useSubscription } from '@urql/vue'
import type { DeepMaybeRef, MaybeRef } from '@vueuse/core'
import type {
  WorkspaceStatusesQuery,
  WorkspaceStatusesQueryVariables,
  UpdatedWorkspacesStatusesSubscription,
  UpdatedWorkspacesStatusesSubscriptionVariables,
} from './__generated__/useUpdatedWorkspacesStatuses'
import type { UpdateResolver } from '@urql/exchange-graphcache'

const UPDATED_WORKSPACES_STATUSES = gql`
  subscription UpdatedWorkspacesStatuses($workspaceIds: [ID!]!) {
    updatedWorkspacesStatuses(workspaceIds: $workspaceIds) {
      id
      type
      title
      description
      timestamp
      detailsUrl
      workspace {
        id
      }
    }
  }
`

const WORKSPACE_STATUSES = gql`
  query WorkspaceStatuses($id: ID!) {
    workspace(id: $id) {
      id
      statuses {
        id
        title
      }
    }
  }
`

export function useUpdatedWorkspacesStatuses(workspacesIds: MaybeRef<string[]>) {
  useSubscription<
    UpdatedWorkspacesStatusesSubscription,
    DeepMaybeRef<UpdatedWorkspacesStatusesSubscriptionVariables>
  >({
    query: UPDATED_WORKSPACES_STATUSES,
    variables: { workspaceIds: workspacesIds },
  })
}

export const updatedWorkspacesStatusesResolver: UpdateResolver<
  UpdatedWorkspacesStatusesSubscription,
  UpdatedWorkspacesStatusesSubscriptionVariables
> = (result, args, cache, info) => {
  const newStatus = result.updatedWorkspacesStatuses
  if (!newStatus.workspace) return
  cache.updateQuery<WorkspaceStatusesQuery, WorkspaceStatusesQueryVariables>(
    {
      query: WORKSPACE_STATUSES,
      variables: { id: newStatus.workspace.id },
    },
    (data) => {
      if (!data) return data
      if (!data.workspace) return data

      data.workspace.statuses = [
        ...data.workspace.statuses.filter((s) => s.title != newStatus.title),
        newStatus,
      ]
      return data
    }
  )
}
