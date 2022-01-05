import { gql, useSubscription } from '@urql/vue'
import { DeepMaybeRef, MaybeRef } from '@vueuse/core'
import {
  ChangeQuery,
  ChangeQueryVariables,
  UpdatedChangesStatusesSubscription,
  UpdatedChangesStatusesSubscriptionVariables,
} from './__generated__/useUpdatedChangesStatuses'
import { UpdateResolver } from '@urql/exchange-graphcache'

const UPDATED_CHANGES_STATUSES = gql`
  subscription UpdatedChangesStatuses($changeIDs: [ID!]!) {
    updatedChangesStatuses(changeIDs: $changeIDs) {
      id
      type
      title
      description
      timestamp
      detailsUrl
      change {
        id
      }
    }
  }
`

const CHANGE_STATUSES = gql`
  query change($id: ID!) {
    change(id: $id) {
      id
      statuses {
        id
        title
      }
    }
  }
`

export function useUpdatedChangesStatuses(changeIDs: MaybeRef<string[]>) {
  useSubscription<
    UpdatedChangesStatusesSubscription,
    DeepMaybeRef<UpdatedChangesStatusesSubscriptionVariables>
  >({
    query: UPDATED_CHANGES_STATUSES,
    variables: { changeIDs: changeIDs },
  })
}

export const updatedChangesStatusesResolver: UpdateResolver<
  UpdatedChangesStatusesSubscription,
  UpdatedChangesStatusesSubscriptionVariables
> = (result, args, cache, info) => {
  const newStatus = result.updatedChangesStatuses
  if (!newStatus.change) return
  cache.updateQuery<ChangeQuery, ChangeQueryVariables>(
    {
      query: CHANGE_STATUSES,
      variables: { id: newStatus.change.id },
    },
    (data) => {
      if (!data) return data
      if (!data.change) return data

      data.change.statuses = [
        ...data.change.statuses.filter((s) => s.title != newStatus.title),
        newStatus,
      ]
      return data
    }
  )
}
