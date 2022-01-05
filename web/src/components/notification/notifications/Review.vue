<template>
  <div class="relative">
    <Avatar
      class="rounded-full bg-gray-400 flex items-center justify-center ring-8 ring-white"
      size="10"
      :author="data.review.author"
    />

    <span class="absolute -bottom-0.5 -right-1 bg-white rounded-tl px-0.5 py-px">
      <ChatAltIcon class="h-5 w-5 text-gray-400" aria-hidden="true" />
    </span>
  </div>
  <div class="min-w-0 flex-1 break-words">
    <div>
      <div class="text-sm">
        <a href="#" class="font-medium text-gray-900">{{ data.review.author.name }}</a>
      </div>
      <p class="mt-0.5 text-sm text-gray-500">
        reviewed

        <span
          class="relative inline-flex items-center rounded-full border border-gray-300 px-3 py-0.5 mx-1"
        >
          <ThumbUpIcon
            v-if="data.review.grade === 'Approve'"
            class="h-5 w-5"
            :class="[
              data.review.dismissedAt || data.review.isReplaced
                ? 'text-gray-500'
                : 'text-green-400',
            ]"
            title="Approved"
          />
          <InformationCircleIcon
            v-else-if="data.review.grade === 'Reject'"
            e
            class="h-5 w-5"
            :class="[
              data.review.dismissedAt || data.review.isReplaced
                ? 'text-gray-500'
                : 'text-orange-400',
            ]"
            title="Rejected"
          />
        </span>

        on

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
import { ChatAltIcon, InformationCircleIcon, ThumbUpIcon } from '@heroicons/vue/solid'
import Avatar from '../../shared/Avatar.vue'
import time from '../../../time'
import { Slug } from '../../../slug'
import { gql } from '@urql/vue'

export const REVIEW_NOTIFICATION_FRAGMENT = gql`
  fragment ReviewNotification on ReviewNotification {
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
      grade
      isReplaced
      dismissedAt

      workspace {
        id
        name
      }

      author {
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
    ThumbUpIcon,
    InformationCircleIcon,
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
