<template>
  <div class="flex flex-col gap-2 flex-1">
    <template v-for="change in changes" :key="change.id">
      <ChangelogChange :codebase-slug="codebaseSlug" :change="change" />
    </template>
  </div>
</template>

<script lang="ts">
import { PropType } from 'vue'

import ChangelogChange, { CHANGELOG_CHANGE_FRAGMENT } from '../../molecules/ChangelogChange.vue'

import { Changelog_ChangeFragment } from './__generated__/ChangeList'

export { CHANGELOG_CHANGE_FRAGMENT }

export default {
  components: { ChangelogChange },
  props: {
    codebaseSlug: {
      type: String,
      required: true,
    },
    changes: {
      type: Array as PropType<Changelog_ChangeFragment[]>,
      required: true,
    },
  },
}
</script>
