<template>
  <router-link
    class="border rounded-md shadow-sm px-4 py-2 hover:bg-gray-50"
    :to="{
      name: 'codebaseChangelog',
      params: { codebaseSlug, selectedChangeID: change.id },
    }"
  >
    <div class="flex flex-col gap-2">
      <div class="flex items-center justify-between">
        <div class="flex items-center gap-2">
          <Avatar :author="change.author" size="6" />
          <span class="font-semibold text-sm text-black">{{ change.author.name }}</span>
        </div>

        <span class="flex items-center text-xs text-gray-500 gap-1">
          <StatusBadge :statuses="change.statuses" />
          <RelativeTime :date="createdAt" />
        </span>
      </div>

      <span v-if="showDescription" class="text-sm prose line-clamp-6" v-html="change.description" />
      <span v-else class="text-sm flex-1 line-clamp-1">
        {{ change.title }}
      </span>

      <div v-if="change.comments.length > 0" class="flex-none">
        <ChangeCommentsIndicator :change="change" />
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
import RelativeTime from '../atoms/RelativeTime.vue'

import { ChangelogChangeFragment } from './__generated__/ChangelogChange'

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
