<template>
  <router-link
    class="border rounded-md shadow-sm px-4 py-2 hover:bg-gray-50"
    :to="{
      name: 'codebaseChange',
      params: { codebaseSlug, id: change.id },
    }"
  >
    <div class="flex flex-col gap-2 h-full">
      <div class="flex items-center justify-between">
        <div class="flex items-center gap-2">
          <Avatar :author="change.author" size="6" />
          <span class="font-semibold text-sm text-black">{{ change.author.name }}</span>
        </div>

        <RelativeTime class="text-sm text-gray-500" :date="createdAt" />
      </div>

      <div class="flex w-full h-full">
        <span
          v-if="showDescription"
          class="flex-1 text-sm line-clamp-12"
          v-html="change.description"
        />
        <span v-else class="flex-1 text-sm line-clamp-1">
          {{ change.title }}
        </span>

        <div class="flex h-5 items-center gap-1">
          <ChangeCommentsIndicator v-if="change.comments.length > 0" :change="change" />
          <StatusBadge :statuses="change.statuses" />
        </div>
      </div>
    </div>
  </router-link>
</template>

<script lang="ts">
import { gql } from '@urql/vue'
import type { PropType } from 'vue'

import { AUTHOR } from '../components/shared/AvatarHelper'
import Avatar from '../components/shared/Avatar.vue'
import StatusBadge, { STATUS_FRAGMENT } from '../components/statuses/StatusBadge.vue'
import ChangeCommentsIndicator, {
  CHANGE_COMMENTS,
} from '../components/changelog/ChangeCommentsIndicator.vue'
import RelativeTime from '../atoms/RelativeTime.vue'

import type { ChangelogChangeFragment } from './__generated__/ChangelogChange'

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
  components: { Avatar, StatusBadge, ChangeCommentsIndicator, RelativeTime },
  props: {
    change: {
      type: Object as PropType<ChangelogChangeFragment>,
      required: true,
    },
    codebaseSlug: {
      type: String,
      required: true,
    },
    showDescription: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    createdAt() {
      return new Date(this.change.createdAt * 1000)
    },
  },
}
</script>
