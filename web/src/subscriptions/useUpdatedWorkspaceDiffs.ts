import { gql, useSubscription } from '@urql/vue'
import { Entity, UpdateResolver } from '@urql/exchange-graphcache'
import { DeepMaybeRef } from '@vueuse/core'
import {
  UpdatedWorkspaceDiffsSubscription,
  UpdatedWorkspaceDiffsSubscriptionVariables,
} from './__generated__/useUpdatedWorkspaceDiffs'
import { Ref } from 'vue'

const UPDATED_WORKSPACE_DIFFS = gql`
  subscription UpdatedWorkspaceDiffs($workspaceID: ID!) {
    updatedWorkspaceDiffs(workspaceID: $workspaceID) {
      id
      preferredName
      hunks {
        id
        patch
      }
    }
  }
`

export function useUpdatedWorkspaceDiffs(
  workspaceID: string,
  pause: Ref<boolean> | boolean | undefined
) {
  useSubscription<
    UpdatedWorkspaceDiffsSubscription,
    // UpdatedWorkspaceDiffsSubscription,
    DeepMaybeRef<UpdatedWorkspaceDiffsSubscriptionVariables>
  >({
    query: UPDATED_WORKSPACE_DIFFS,
    pause,
    variables: { workspaceID },
  })
}

export const updateWorkspaceDiffsResolver: UpdateResolver<
  UpdatedWorkspaceDiffsSubscription,
  UpdatedWorkspaceDiffsSubscriptionVariables
> = (result, args, cache, info) => {
  // Replace the existing list of diffs in the workspace with the new list
  if (result && result.updatedWorkspaceDiffs) {
    const newKeys = Array<Entity>()
    for (const diff of result.updatedWorkspaceDiffs) {
      if (diff.__typename) {
        newKeys.push(cache.keyOfEntity({ __typename: diff.__typename, id: diff.id }))
      }
    }
    cache.link({ __typename: 'Workspace', id: args.workspaceID }, 'diffs', newKeys)
  }
}
