<template>
  <router-link
    class="border rounded-md shadow-sm px-4 py-2 hover:bg-gray-50"
    :to="{
      name: 'codebaseChangelog',
      params: { codebaseSlug, selectedChangeID: change.id },
    }"
  >
    <div class="flex flex-row gap-2 items-center justify-between">
      <div class="flex-none">
        <Avatar :author="change.author" size="6" />
      </div>
      <div class="text-sm flex-1">{{ change.title }}</div>
      <div v-if="change.comments.length > 0" class="flex-none">
        <ChangeCommentsIndicator :change="change" />
      </div>
      <div class="flex text-sm text-gray-500">
        <div class="mr-1">
          <StatusBadge :statuses="change.statuses" />
        </div>
        {{ timeAgo }}
      </div>
    </div>
  </router-link>
</template>

<script lang="ts">
import { gql } from '@urql/vue'
import { PropType } from 'vue'

import { AUTHOR } from '../components/shared/AvatarHelper'
import Avatar from '../components/shared/Avatar.vue'
import StatusBadge, { STATUS_FRAGMENT } from '../components/statuses/StatusBadge.vue'
import ChangeCommentsIndicator, {
  CHANGE_COMMENTS,
} from '../components/changelog/ChangeCommentsIndicator.vue'

import { ChangelogChangeFragment } from './__generated__/ChangelogChange'

import time from '../time'

export const CHANGELOG_CHANGE_FRAGMENT = gql`
  fragment ChangelogChange on Change {
    id
    title
    description
    createdAt
    author {
      ...Author
    }
    ...ChangeComments
    statuses {
      ...Status
    }
  }
  ${AUTHOR}
  ${STATUS_FRAGMENT}
  ${CHANGE_COMMENTS}
`

export default {
  components: { Avatar, StatusBadge, ChangeCommentsIndicator },
  props: {
    change: {
      type: Object as PropType<ChangelogChangeFragment>,
      required: true,
    },
    codebaseSlug: {
      type: String,
      required: true,
    },
  },
  computed: {
    timeAgo() {
      return time.getRelativeTime(new Date(this.change.createdAt * 1000))
    },
  },
}
</script>
