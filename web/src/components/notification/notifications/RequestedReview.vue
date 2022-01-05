<template>
  <div class="relative">
    <Avatar
      class="rounded-full bg-gray-400 flex items-center justify-center ring-8 ring-white"
      size="10"
      :author="data.review.requestedBy"
    />

    <span class="absolute -bottom-0.5 -right-1 bg-white rounded-tl px-0.5 py-px">
      <ChatAltIcon class="h-5 w-5 text-gray-400" aria-hidden="true" />
    </span>
  </div>
  <div class="min-w-0 flex-1 break-words">
    <div>
      <div class="text-sm">
        <a href="#" class="font-medium text-gray-900">{{ data.review.requestedBy.name }}</a>
      </div>
      <p class="mt-0.5 text-sm text-gray-500">
        asked for your feedback on
        <router-link
          class="underline"
          :to="{
            name: 'workspaceHome',
            params: {
              codebaseSlug: codebase_slug,
              id: data.review.workspace.id,
            },
          }"
          @click="$emit('close')"
        >
          <strong>{{ data.review.workspace.name }}</strong>
        </router-link>
        {{ friendly_ago }}
      </p>
    </div>
  </div>
</template>

<script>
import { ChatAltIcon } from '@heroicons/vue/solid'
import Avatar from '../../shared/Avatar.vue'
import time from '../../../time'
import { Slug } from '../../../slug'
import { gql } from '@urql/vue'

export const REQUESTED_REVIEW_NOTIFICATION_FRAGMENT = gql`
  fragment RequestedReviewNotification on RequestedReviewNotification {
    id
    type
    createdAt
    codebase {
      id
      shortID
      name
    }

    review {
      id
      workspace {
        id
        name
      }
      requestedBy {
        id
        name
        avatarUrl
      }
    }
  }
`

export default {
  components: {
    ChatAltIcon,
    Avatar,
  },
  props: ['data', 'now'],
  emits: ['close'],
  computed: {
    friendly_ago() {
      return time.getRelativeTime(new Date(this.data.createdAt * 1000), this.now)
    },
    codebase_slug() {
      return Slug(this.data.codebase.name, this.data.codebase.shortID)
    },
  },
}
</script>
