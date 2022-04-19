<template>
  <div class="flow-root">
    <ul class="-mb-8">
      <li v-for="(notification, idx) in showNotifications" :key="notification.id">
        <div class="relative pb-8">
          <span
            v-if="idx !== notifications.length - 1"
            class="absolute top-5 left-5 -ml-px h-full w-0.5 bg-gray-200"
            aria-hidden="true"
          />
          <div class="relative flex items-start space-x-3 flex-nowrap">
            <CommentNotification
              v-if="notification.comment"
              :data="notification"
              :user="user"
              @close="$emit('close')"
            />
            <InvitedToCodebase
              v-else-if="notification.__typename === 'InvitedToCodebaseNotification'"
              :data="notification"
              @close="$emit('close')"
            />
            <InvitedToOrganization
              v-else-if="notification.__typename === 'InvitedToOrganizationNotification'"
              :data="notification"
              @close="$emit('close')"
            />
            <RequestedReviewNotification
              v-else-if="notification.__typename === 'RequestedReviewNotification'"
              :data="notification"
              @close="$emit('close')"
            />
            <ReviewNotification
              v-else-if="notification.__typename === 'ReviewNotification'"
              :data="notification"
              @close="$emit('close')"
            />
            <NewSuggestionNotification
              v-else-if="notification.__typename === 'NewSuggestionNotification'"
              :data="notification"
              @close="$emit('close')"
            />
            <GitHubRepositoryImportedNotification
              v-else-if="notification.__typename === 'GitHubRepositoryImported'"
              :data="notification"
              @close="$emit('close')"
            />
          </div>
        </div>
      </li>
    </ul>
  </div>
</template>

<script>
import CommentNotification, { NOTIFICATION_COMMENT_FRAGMENT } from './notifications/Comment.vue'
import { onUnmounted, ref } from 'vue'
import RequestedReviewNotification, {
  REQUESTED_REVIEW_NOTIFICATION_FRAGMENT,
} from './notifications/RequestedReview.vue'
import ReviewNotification, { REVIEW_NOTIFICATION_FRAGMENT } from './notifications/Review.vue'
import NewSuggestionNotification, {
  NEW_SUGGESTION_NOTIFICATION_FRAGMENT,
} from './notifications/NewSuggestion.vue'
import GitHubRepositoryImportedNotification, {
  GITHUB_REPOSITORY_IMPORTED_NOTIFICATION_FRAGMENT,
} from './notifications/GitHubRepositoryImoprted.vue'
import { gql } from '@urql/vue'
import InvitedToCodebase, {
  INVITED_TO_CODEBASE_NOTIFICATION_FRAGMENT,
} from './notifications/InvitedToCodebase.vue'
import InvitedToOrganization, {
  INVITED_TO_ORGANIZATION_NOTIFICATION_FRAGMENT,
} from './notifications/InvitedToOrganization.vue'

export const NOTIFICATION_FRAGMENT = gql`
  fragment NotificationData on Notification {
    id
    ... on CommentNotification {
      ...NotificationComment
    }
    ... on RequestedReviewNotification {
      ...RequestedReviewNotification
    }
    ... on ReviewNotification {
      ...ReviewNotification
    }
    ... on NewSuggestionNotification {
      ...NewSuggestionNotification
    }
    ... on InvitedToCodebaseNotification {
      ...InvitedToCodebaseNotification
    }
    ... on InvitedToOrganizationNotification {
      ...InvitedToOrganizationNotification
    }
    ... on GitHubRepositoryImported @include(if: $isGitHubEnabled) {
      ...GitHubRepositoryImported
    }
  }
  ${INVITED_TO_ORGANIZATION_NOTIFICATION_FRAGMENT}
  ${INVITED_TO_CODEBASE_NOTIFICATION_FRAGMENT}
  ${NOTIFICATION_COMMENT_FRAGMENT}
  ${NEW_SUGGESTION_NOTIFICATION_FRAGMENT}
  ${REQUESTED_REVIEW_NOTIFICATION_FRAGMENT}
  ${REVIEW_NOTIFICATION_FRAGMENT}
  ${GITHUB_REPOSITORY_IMPORTED_NOTIFICATION_FRAGMENT}
`

const supportedTypes = {
  CommentNotification: true,
  RequestedReviewNotification: true,
  ReviewNotification: true,
  NewSuggestionNotification: true,
  GitHubRepositoryImported: true,
  InvitedToCodebaseNotification: true,
  InvitedToOrganizationNotification: true,
}

export default {
  components: {
    ReviewNotification,
    RequestedReviewNotification,
    CommentNotification,
    NewSuggestionNotification,
    GitHubRepositoryImportedNotification,
    InvitedToCodebase,
    InvitedToOrganization,
  },
  props: {
    notifications: {
      type: Array,
      required: true,
    },
    user: {
      type: Object,
      required: true,
    },
  },
  emits: ['close'],
  setup() {
    const now = ref(new Date())
    const updateNow = setInterval(() => {
      now.value = new Date()
    }, 1000)
    onUnmounted(() => {
      clearInterval(updateNow)
    })

    return {
      now,
    }
  },
  computed: {
    showNotifications() {
      const showNotification = (n) => {
        switch (true) {
          case n.__typename === 'NewSuggestionNotification':
            return !!n.suggestion.for
          default:
            return supportedTypes[n.__typename]
        }
      }
      return this.notifications.filter(showNotification)
    },
  },
}
</script>
