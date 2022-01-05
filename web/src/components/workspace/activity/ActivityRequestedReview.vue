<template>
  <div class="relative flex items-start space-x-3">
    <div>
      <div class="relative px-1">
        <div
          class="h-8 w-8 bg-gray-100 rounded-full ring-8 ring-white flex items-center justify-center"
        >
          <UserCircleIcon class="h-5 w-5 text-gray-500" aria-hidden="true" />
        </div>
      </div>
    </div>
    <div class="min-w-0 flex-1 py-1.5">
      <div class="text-sm text-gray-500" :class="[item.review.dismissedAt ? 'line-through' : '']">
        <a class="font-medium text-gray-900">{{ item.author.name }}</a>
        {{ ' ' }}
        asked for feedback from
        {{ ' ' }}
        <a class="font-medium text-gray-900">{{ item.review.author.name }}</a>
        {{ ' ' }}
        <span class="whitespace-nowrap">{{ friendly_ago(item.createdAt) }}</span>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { UserCircleIcon } from '@heroicons/vue/solid'
import { PropType } from 'vue'
import time from '../../../time'
import { gql } from '@urql/vue'
import { WorkspaceActivityRequestedReviewFragment } from './__generated__/WorkspaceActivityRequestedReview'

export const WORKSPACE_ACTIVITY_REQUESTED_REVIEW_FRAGMENT = gql`
  fragment WorkspaceActivityRequestedReview on WorkspaceRequestedReviewActivity {
    createdAt
    review {
      id
      grade
      createdAt
      dismissedAt
      isReplaced
      author {
        id
        name
        avatarUrl
      }
    }
  }
`

export default {
  components: { UserCircleIcon },
  props: {
    item: {
      type: Object as PropType<WorkspaceActivityRequestedReviewFragment>,
      required: true,
    },
  },
  methods: {
    friendly_ago(ts: number) {
      return time.getRelativeTime(new Date(ts * 1000))
    },
  },
}
</script>
