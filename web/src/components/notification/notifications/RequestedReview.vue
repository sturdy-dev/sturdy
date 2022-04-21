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
        <span class="font-medium text-gray-900">{{ data.review.requestedBy.name }}</span>
      </div>
      <p class="mt-0.5 text-sm text-gray-500 space-x-1">
        asked for your feedback on
        <router-link
          class="underline"
          :to="{
            name: 'workspaceHome',
            params: {
              codebaseSlug: codebaseSlug,
              id: data.review.workspace.id,
            },
          }"
          @click="$emit('close')"
        >
          <strong>{{ data.review.workspace.name }}</strong>
        </router-link>
        <RelativeTime :date="createdAt" />
      </p>
    </div>
  </div>
</template>

<script lang="ts">
import { ChatAltIcon } from '@heroicons/vue/solid'
import Avatar from '../../../atoms/Avatar.vue'
import RelativeTime from '../../../atoms/RelativeTime.vue'
import { Slug } from '../../../slug'
import { gql } from '@urql/vue'
import { defineComponent, type PropType } from 'vue'
import type { RequestedReviewNotificationFragment } from './__generated__/RequestedReview'

export const REQUESTED_REVIEW_NOTIFICATION_FRAGMENT = gql`
  fragment RequestedReviewNotification on RequestedReviewNotification {
    id
    type
    createdAt
    review {
      id
      workspace {
        id
        name
        codebase {
          id
          shortID
          name
          members {
            id
            name
          }
        }
      }
      requestedBy {
        id
        name
        avatarUrl
      }
    }
  }
`

export default defineComponent({
  components: {
    ChatAltIcon,
    Avatar,
    RelativeTime,
  },
  props: {
    data: {
      type: Object as PropType<RequestedReviewNotificationFragment>,
      required: true,
    },
  },
  emits: ['close'],
  computed: {
    createdAt() {
      return new Date(this.data.createdAt * 1000)
    },
    codebaseSlug() {
      return Slug(
        this.data.review.workspace.codebase.name,
        this.data.review.workspace.codebase.shortID
      )
    },
  },
})
</script>
