<template>
  <div>
    <button
      class="flex items-center justify-center cursor-pointer p-2 hover:bg-warmgray-300 text-gray-400 hover:text-gray-700 transition rounded-md"
      @click="clickIcon"
    >
      <span class="sr-only">Open options</span>

      <div class="relative">
        <BellIconSolid class="h-5 w-5" aria-hidden="true" />
        <span
          v-if="nonArchivedNotifications?.length"
          class="flex absolute h-5 w-5 top-0 right-0 -mt-3 -mr-3"
        >
          <span
            class="relative inline-flex items-center justify-center px-2 py-2 text-xs font-bold leading-none text-white bg-red-500 rounded-full"
          >
            {{ nonArchivedNotifications?.length }}
          </span>
        </span>
      </div>
    </button>

    <NotificationOverlay
      v-if="data"
      :open="open"
      :notifications="data.notifications"
      :user="user"
      @close="open = false"
    />
  </div>
</template>

<script lang="ts">
import { defineComponent, PropType } from 'vue'
import { BellIcon as BellIconSolid } from '@heroicons/vue/solid'
import NotificationOverlay from './Overlay.vue'
import { NOTIFICATION_FRAGMENT as NOTIFICATION_DATA_FRAGMENT } from './Feed.vue'
import { gql, useMutation, useQuery } from '@urql/vue'
import { useUpdatedNotifications } from '../../subscriptions/useUpdatedNotifications'
import {
  User,
  NotificationType,
  NotificationChannel,
  NotificationPreference,
} from '../../__generated__/types'
import mentionify from '../shared/mentionify'
import {
  NotificationFragment,
  NotificationIconQuery,
  NotificationIconQueryVariables,
} from './__generated__/Icon'
import { Slug } from '../../slug'

const NOTIFICATION_FRAGMENT = gql`
  fragment Notification on Notification {
    id
    archivedAt
    codebase {
      id
      members {
        id
        name
      }
    }
    ...NotificationData
  }
  ${NOTIFICATION_DATA_FRAGMENT}
`

const notificationTitle = (data: NotificationFragment): string => {
  switch (data.__typename) {
    case 'ReviewNotification':
      return data.review?.author?.name + ' reviewed ' + data.review?.workspace?.name
    case 'NewSuggestionNotification':
      return `${data.suggestion.author.name} proposed a suggestion`
    case 'GitHubRepositoryImported':
      return `Repository imported`
    case 'RequestedReviewNotification':
      if (data.review.requestedBy) {
        return data.review.requestedBy.name + ' requested review on ' + data.review.workspace.name
      }
      return 'Review requested on ' + data.review.workspace.name
    case 'CommentNotification':
      switch (data.comment.__typename) {
        case 'ReplyComment':
          if (data.comment?.parent?.workspace) {
            return data.comment.author.name + ' replied in ' + data.comment.parent.workspace.name
          }
          return data.comment.author.name + ' replied'
        case 'TopComment':
          if (data.comment?.workspace) {
            return data.comment.author.name + ' commented on ' + data.comment.workspace.name
          }
          return data.comment.author.name + ' commented'
        default:
          return 'New comment'
      }
    default:
      return 'New notification in Sturdy'
  }
}
const notificationBody = (data: NotificationFragment): string | undefined => {
  switch (data.__typename) {
    case 'CommentNotification':
      return data.comment.message
    case 'GitHubRepositoryImported':
      return `${data.repository.name} is now available`
    case 'NewSuggestionNotification':
      return `${data.suggestion.author.name} proposed a suggestion on ${data.suggestion.for.name}`
    case 'RequestedReviewNotification':
      return undefined
    case 'ReviewNotification':
      if (data.review?.grade === 'Approve') {
        return 'Looks good to me!'
      }
      if (data.review?.grade === 'Reject') {
        return 'I have some feedback'
      }
      return undefined
    default:
      return undefined
  }
}

const notificationIcon = (data: NotificationFragment): string => {
  const defaultIcon = 'https://getsturdy.com/favicon.ico'
  switch (data.__typename) {
    case 'CommentNotification':
      return data.comment.author.avatarUrl ? data.comment.author.avatarUrl : defaultIcon
    case 'GitHubRepositoryImported':
      return defaultIcon
    case 'NewSuggestionNotification':
      return data.suggestion.author.avatarUrl ? data.suggestion.author.avatarUrl : defaultIcon
    case 'RequestedReviewNotification':
      return defaultIcon
    case 'ReviewNotification':
      return data.review.author.avatarUrl ? data.review.author.avatarUrl : defaultIcon
    default:
      return defaultIcon
  }
}

const notificationUrl = (data: NotificationFragment): string => {
  switch (data.__typename) {
    case 'CommentNotification':
      switch (data.comment.__typename) {
        case 'ReplyComment':
          if (data.comment.parent.workspace) {
            return `/${Slug(data.codebase.name, data.codebase.shortID)}/${
              data.comment.parent.workspace.id
            }#${data.comment.id}`
          } else if (data.comment.parent.change) {
            return `/${Slug(data.codebase.name, data.codebase.shortID)}/${
              data.comment.parent.change.id
            }#${data.comment.id}`
          } else {
            return '/'
          }
        case 'TopComment':
          if (data.comment.workspace) {
            return `/${Slug(data.codebase.name, data.codebase.shortID)}/${
              data.comment.workspace.id
            }#${data.comment.id}`
          } else if (data.comment.change) {
            return `/${Slug(data.codebase.name, data.codebase.shortID)}/${data.comment.change.id}#${
              data.comment.id
            }`
          } else {
            return '/'
          }
        default:
          return '/'
      }
    case 'GitHubRepositoryImported':
      return `/${Slug(data.codebase.name, data.codebase.shortID)}`
    case 'NewSuggestionNotification':
      return `/${Slug(data.codebase.name, data.codebase.shortID)}/${data.suggestion.for.id}`
    case 'RequestedReviewNotification':
      return `/${Slug(data.codebase.name, data.codebase.shortID)}/${data.review.workspace.id}`
    case 'ReviewNotification':
      return `/${Slug(data.codebase.name, data.codebase.shortID)}/${data.review.workspace.id}`
    default:
      return '/'
  }
}

export default defineComponent({
  components: {
    NotificationOverlay,
    BellIconSolid,
  },
  props: {
    user: { type: Object as PropType<User>, required: true },
  },
  setup() {
    let { data, fetching, error, executeQuery } = useQuery<
      NotificationIconQuery,
      NotificationIconQueryVariables
    >({
      query: gql`
        query NotificationIcon {
          notifications {
            ...Notification
          }
        }
        ${NOTIFICATION_FRAGMENT}
      `,
      requestPolicy: 'cache-and-network',
    })

    let { executeMutation: archiveNotificationsResult } = useMutation(
      gql`
        mutation NotificationIconArchive($ids: [ID!]!) {
          archiveNotifications(input: { ids: $ids }) {
            id
            archivedAt
          }
        }
      `
    )

    useUpdatedNotifications()

    return {
      data,
      fetching,
      error,
      refresh: async () => {
        await executeQuery({
          requestPolicy: 'network-only',
        })
      },

      async archiveNotifications(ids: string[]) {
        const variables = { ids }
        await archiveNotificationsResult(variables).then((result) => {
          console.log('archive notifications', result)
        })
      },
    }
  },
  data() {
    return {
      open: false,
      mostRecentNotification: null as NotificationFragment | null,
      initedNotifications: false,
    }
  },
  computed: {
    nonArchivedNotifications() {
      return this.data?.notifications?.filter((n: NotificationFragment) => !n.archivedAt)
    },
    notifyOn(): NotificationType[] {
      const filterWebNotifications = (p: NotificationPreference): boolean => {
        const isWeb = p.channel === NotificationChannel.Web
        const isEnabled = p.enabled
        return isWeb && isEnabled
      }
      return this.user.notificationPreferences.filter(filterWebNotifications).map((p) => p.type)
    },
  },
  watch: {
    'data.notifications': function (notifications) {
      if (notifications) {
        this.sendNotifications(notifications)
      }
    },
  },
  methods: {
    async clickIcon() {
      this.open = !this.open
      if (this.open) {
        await this.refresh()
        this.archiveAll()
        this.requestPermission()
      }
    },
    archiveAll() {
      if (!this.data) return
      const ids = this.data.notifications.filter((n) => !n.archivedAt).map((n) => n.id)
      if (ids.length === 0) return
      this.archiveNotifications(ids)
    },
    supportsBrowserNotifications() {
      return 'Notification' in window
    },
    sendNotification(notification: NotificationFragment) {
      if (!this.supportsBrowserNotifications()) return
      const title = notificationTitle(notification)
      const dirtyNotificationBody = notificationBody(notification)
      const body = dirtyNotificationBody
        ? mentionify(
            dirtyNotificationBody,
            '@',
            [this.user].concat(notification.codebase.members as User[])
          )
        : dirtyNotificationBody
      const icon = notificationIcon(notification)
      const url = notificationUrl(notification)

      try {
        const n = new Notification(title, {
          body: body,
          icon: icon,
          data: {
            id: notification.id,
            url: url,
          },
        })
        n.onclick = (e) => {
          const notification = e.target as Notification
          this.$router.push(notification.data.url)
          this.archiveNotifications([notification.data.id])
        }
      } catch (err) {
        console.error('failed to send notification', err)
      }
    },
    requestPermission() {
      Notification.requestPermission(function (result) {
        console.log('requestPermission', result)
      })
    },
    sendNotifications(notifications: NotificationFragment[]) {
      // Init with the most recent notification
      if (!this.initedNotifications) {
        if (notifications.length === 0) {
          this.mostRecentNotification = null
          this.initedNotifications = true
          return
        }

        this.mostRecentNotification = notifications[0]
        this.initedNotifications = true
        return
      }

      // Send for all newer than mostRecentNotification
      for (const notif of notifications) {
        if (notif.id === this.mostRecentNotification?.id) {
          break
        }

        if (!this.notifyOn.includes(notif.type)) {
          continue
        }

        this.sendNotification(notif)
      }

      if (notifications.length > 0) {
        this.mostRecentNotification = notifications[0]
      }
    },
  },
})
</script>
