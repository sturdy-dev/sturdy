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
        <div class="inline-flex items-center">
          <a class="font-medium text-gray-900">{{ item.author.name }}</a>
          {{ ' ' }}

          <span
            class="relative inline-flex items-center rounded-full border border-gray-300 px-3 py-0.5 mx-1"
          >
            <ThumbUpIcon
              v-if="item.review.grade === 'Approve'"
              class="h-5 w-5"
              :class="[
                item.review.dismissedAt || item.review.isReplaced
                  ? 'text-gray-500'
                  : 'text-green-400',
              ]"
              title="Approved"
            />
            <InformationCircleIcon
              v-else-if="item.review.grade === 'Reject'"
              e
              class="h-5 w-5"
              :class="[
                item.review.dismissedAt || item.review.isReplaced
                  ? 'text-gray-500'
                  : 'text-orange-400',
              ]"
              title="Rejected"
            />
          </span>
        </div>
        {{ ' ' }}
        this change
        <span class="whitespace-nowrap">{{ friendly_ago(item.createdAt) }}</span>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { InformationCircleIcon, ThumbUpIcon, UserCircleIcon } from '@heroicons/vue/solid'
import time from '../../time'
import { PropType } from 'vue'
import { gql } from '@urql/vue'
import { WorkspaceReviewedActivityFragment } from './__generated__/WorkspaceActivityReviewed'

export const WORKSPACE_ACTIVITY_REVIEWED_FRAGMENT = gql`
  fragment WorkspaceReviewedActivity on WorkspaceReviewedActivity {
    createdAt
    author {
      id
      name
    }
    review {
      id
      grade
      createdAt
      isReplaced
      dismissedAt
      author {
        id
        name
        avatarUrl
      }
    }
  }
`

export default {
  name: 'WorkspaceActivityComment',
  components: {
    UserCircleIcon,
    ThumbUpIcon,
    InformationCircleIcon,
  },
  props: {
    item: {
      type: Object as PropType<WorkspaceReviewedActivityFragment>,
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
