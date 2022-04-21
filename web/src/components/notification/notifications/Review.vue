<template>
  <div class="relative">
    <Avatar
      class="rounded-full bg-gray-400 flex items-center justify-center ring-8 ring-white"
      size="10"
      :author="data.review.author"
    />

    <span class="absolute -bottom-0.5 -right-1 bg-white rounded px-0.5 py-px">
      <ThumbUpIcon v-if="isApproved" class="h-5 w-5 text-gray-400" title="Approved" />
      <InformationCircleIcon
        v-else-if="isRejected"
        class="h-5 w-5 text-gray-400"
        title="Rejected"
      />
      <ChatAltIcon v-else class="h-5 w-5 text-gray-400" aria-hidden="true" />
    </span>
  </div>
  <div class="min-w-0 flex-1 break-words">
    <div>
      <div class="text-sm">
        <span class="font-medium text-gray-900">{{ data.review.author.name }}</span>
      </div>
      <p class="mt-0.5 text-sm text-gray-500 space-x-1">
        <span :class="{ 'line-through	': data.review.dismissedAt || data.review.isReplaced }">
          <template v-if="isApproved">approved&nbsp;</template>
          <template v-else-if="isRejected">has questions about&nbsp;</template>
          <template v-else>reviewed&nbsp;</template>

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
        </span>

        <RelativeTime :date="createdAt" />
      </p>
    </div>
  </div>
</template>

<script lang="ts">
import { ChatAltIcon, InformationCircleIcon, ThumbUpIcon } from '@heroicons/vue/solid'
import Avatar from '../../../atoms/Avatar.vue'
import RelativeTime from '../../../atoms/RelativeTime.vue'
import { Slug } from '../../../slug'
import { gql } from '@urql/vue'
import { defineComponent, type PropType } from 'vue'
import type { ReviewNotificationFragment } from './__generated__/Review'
import { ReviewGrade } from '../../../__generated__/types'

export const REVIEW_NOTIFICATION_FRAGMENT = gql`
  fragment ReviewNotification on ReviewNotification {
    id
    type
    createdAt

    review {
      id
      grade
      isReplaced
      dismissedAt

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

      author {
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
    ThumbUpIcon,
    InformationCircleIcon,
    Avatar,
    RelativeTime,
  },
  props: {
    data: {
      type: Object as PropType<ReviewNotificationFragment>,
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
    isApproved() {
      return this.data.review.grade === ReviewGrade.Approve
    },
    isRejected() {
      return this.data.review.grade === ReviewGrade.Reject
    },
  },
})
</script>
