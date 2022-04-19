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
import type { DeepMaybeRef } from '@vueuse/core'
import { defineComponent, inject, ref, computed } from 'vue'
import type { PropType, Ref } from 'vue'
import { BellIcon as BellIconSolid } from '@heroicons/vue/solid'
import NotificationOverlay from './Overlay.vue'
import { NOTIFICATION_FRAGMENT as NOTIFICATION_DATA_FRAGMENT } from './Feed.vue'
import { gql, useMutation, useQuery } from '@urql/vue'
import { useUpdatedNotifications } from '../../subscriptions/useUpdatedNotifications'
import {
  type User,
  NotificationType,
  NotificationChannel,
  type NotificationPreference,
} from '../../__generated__/types'
import mentionify from '../../atoms/mentionify'
import type {
  NotificationFragment,
  NotificationIconQuery,
  NotificationIconQueryVariables,
} from './__generated__/Icon'
import { Slug } from '../../slug'
import { Feature } from '../../__generated__/types'

const NOTIFICATION_FRAGMENT = gql`
  fragment Notification on Notification {
    id
    archivedAt
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
    case 'InvitedToCodebaseNotification':
      return `You have joined ${data.codebase.name}`
    case 'InvitedToOrganizationNotification':
      return `You have joined ${data.organization.name}`
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
    case 'InvitedToCodebaseNotification':
      return `You have been invited to ${data.codebase.name}`
    case 'InvitedToOrganizationNotification':
      return `You have been invited to ${data.organization.name}`
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
    case 'InvitedToCodebaseNotification':
      return defaultIcon
    case 'InvitedToOrganizationNotification':
      return defaultIcon
    default:
      return defaultIcon
  }
}

const notificationCodebaseMembers = (data: NotificationFragment): User[] => {
  switch (data.__typename) {
    case 'CommentNotification':
      return data.comment.codebase.members as User[]
    case 'GitHubRepositoryImported':
      return data.repository.codebase.members as User[]
    case 'NewSuggestionNotification':
      return data.suggestion.for.codebase.members as User[]
    case 'RequestedReviewNotification':
      return data.review.workspace.codebase.members as User[]
    case 'ReviewNotification':
      return data.review.workspace.codebase.members as User[]
    case 'InvitedToCodebaseNotification':
      return data.codebase.members as User[]
    case 'InvitedToOrganizationNotification':
      return data.organization.members as User[]
    default:
      return []
  }
}

const notificationUrl = (data: NotificationFragment): string => {
  switch (data.__typename) {
    case 'CommentNotification':
      switch (data.comment.__typename) {
        case 'ReplyComment':
          if (data.comment.parent.workspace) {
            return `/${Slug(data.comment.codebase.name, data.comment.codebase.shortID)}/${
              data.comment.parent.workspace.id
            }#${data.comment.id}`
          } else if (data.comment.parent.change) {
            return `/${Slug(data.comment.codebase.name, data.comment.codebase.shortID)}/${
              data.comment.parent.change.id
            }#${data.comment.id}`
          } else {
            return '/'
          }
        case 'TopComment':
          if (data.comment.workspace) {
            return `/${Slug(data.comment.codebase.name, data.comment.codebase.shortID)}/${
              data.comment.workspace.id
            }#${data.comment.id}`
          } else if (data.comment.change) {
            return `/${Slug(data.comment.codebase.name, data.comment.codebase.shortID)}/${
              data.comment.change.id
            }#${data.comment.id}`
          } else {
            return '/'
          }
        default:
          return '/'
      }
    case 'GitHubRepositoryImported':
      return `/${Slug(data.repository.codebase.name, data.repository.codebase.shortID)}`
    case 'NewSuggestionNotification':
      return `/${Slug(data.suggestion.for.codebase.name, data.suggestion.for.codebase.shortID)}/${
        data.suggestion.for.id
      }`
    case 'RequestedReviewNotification':
      return `/${Slug(
        data.review.workspace.codebase.name,
        data.review.workspace.codebase.shortID
      )}/${data.review.workspace.id}`
    case 'ReviewNotification':
      return `/${Slug(
        data.review.workspace.codebase.name,
        data.review.workspace.codebase.shortID
      )}/${data.review.workspace.id}`
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
    const features = inject<Ref<Array<Feature>>>('features', ref([]))
    const isGitHubEnabled = computed(() => features?.value?.includes(Feature.GitHub))

    let { data, fetching, error, executeQuery } = useQuery<
      NotificationIconQuery,
      DeepMaybeRef<NotificationIconQueryVariables>
    >({
      query: gql`
        query NotificationIcon($isGitHubEnabled: Boolean!) {
          notifications {
            ...Notification
          }
        }
        ${NOTIFICATION_FRAGMENT}
      `,
      requestPolicy: 'cache-and-network',
      variables: {
        isGitHubEnabled,
      },
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
    nonArchivedNotifications: function (unreadNotifications) {
      if (window?.ipc?.setBadgeCount) window.ipc.setBadgeCount(unreadNotifications?.length)
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
      const members = notificationCodebaseMembers(notification)
      const body = dirtyNotificationBody
        ? mentionify(dirtyNotificationBody, '@', [this.user].concat(members))
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
