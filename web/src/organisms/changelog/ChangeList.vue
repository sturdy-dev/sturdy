<template>
  <div class="flex flex-col gap-2 flex-1">
    <template v-for="change in changes" :key="change.id">
      <ChangelogChange :codebase-slug="codebaseSlug" :change="change" />
    </template>

    <button
      class="text-sm flex flex-row items-center gap-1 text-gray-500 hover:bg-gray-100 self-end px-2 py-1 rounded-md"
      :disabled="!hasNextPage"
      :class="{
        'cursor-not-allowed': !hasNextPage,
        'opacity-50': !hasNextPage,
        'hover:bg-inherit': !hasNextPage,
      }"
      @click="() => $emit('next-page')"
    >
      older changes
      <ArrowRightIcon class="h-4 w-4" />
    </button>
  </div>
</template>

<script lang="ts">
import type { PropType } from 'vue'

import { ArrowRightIcon } from '@heroicons/vue/solid'
import ChangelogChange, { CHANGELOG_CHANGE_FRAGMENT } from '../../molecules/ChangelogChange.vue'

import type { Changelog_ChangeFragment } from './__generated__/ChangeList'

export { CHANGELOG_CHANGE_FRAGMENT }

export default {
  components: { ChangelogChange, ArrowRightIcon },
  props: {
    codebaseSlug: {
      type: String,
      required: true,
    },
    changes: {
      type: Array as PropType<Changelog_ChangeFragment[]>,
      required: true,
    },
    hasNextPage: {
      type: Boolean,
      default: false,
    },
  },
  emits: ['next-page'],
}
</script>
