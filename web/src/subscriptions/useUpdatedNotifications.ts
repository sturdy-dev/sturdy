import { gql, useSubscription } from '@urql/vue'
import {
  UpdatedNotificationsNotificationsQuery,
  UpdatedNotificationsNotificationsQueryVariables,
  UpdatedNotificationsSubscription,
  UpdatedNotificationsSubscriptionVariables,
} from './__generated__/useUpdatedNotifications'
import { UpdateResolver } from '@urql/exchange-graphcache'
import { DeepMaybeRef } from '@vueuse/core'

const UPDATED_NOTIFICATIONS = gql`
  subscription UpdatedNotifications {
    updatedNotifications {
      id
      archivedAt
    }
  }
`

export function useUpdatedNotifications() {
  useSubscription<
    UpdatedNotificationsSubscription,
    UpdatedNotificationsSubscription,
    DeepMaybeRef<UpdatedNotificationsSubscriptionVariables>
  >({
    query: UPDATED_NOTIFICATIONS,
  })
}

const UPDATED_NOTIFICATIONS_NOTIFICATIONS = gql`
  query UpdatedNotificationsNotifications {
    notifications {
      id
    }
  }
`

export const updatedNotificationsUpdateResolver: UpdateResolver<
  UpdatedNotificationsSubscription,
  UpdatedNotificationsSubscriptionVariables
> = (result, args, cache, info) => {
  cache.updateQuery<
    UpdatedNotificationsNotificationsQuery,
    UpdatedNotificationsNotificationsQueryVariables
  >(
    {
      query: UPDATED_NOTIFICATIONS_NOTIFICATIONS,
    },
    (data) => {
      // Add codebase if not exists
      if (
        data?.notifications &&
        !data.notifications.some((c) => c.id === result.updatedNotifications.id)
      ) {
        data.notifications.push(result.updatedNotifications)
      }
      return data
    }
  )
}
