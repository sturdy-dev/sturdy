<template>
  <div v-if="codebase.changes.length > 0" class="grid grid-cols-2 gap-2">
    <ChangelogChange
      :codebase-slug="codebase.shortID"
      :change="latestChange"
      :show-description="true"
    />

    <div v-if="codebase.changes.length > 1" class="flex flex-col gap-2 flex-1">
      <template v-for="change in codebase.changes.slice(1)" :key="change.id">
        <ChangelogChange :codebase-slug="codebase.shortID" :change="change" />
      </template>

      <router-link
        :to="{ name: 'codebaseChanges', params: { codebaseSlug: codebase.shortID } }"
        class="text-sm flex flex-row items-center gap-1 text-gray-500 hover:bg-gray-100 self-end px-2 py-1 rounded-md"
      >
        Older Changes
        <div class="h-4 w-4">
          <ArrowRightIcon />
        </div>
      </router-link>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent, PropType, toRefs } from 'vue'
import gql from 'graphql-tag'
import { TopOfChangelogFragment } from './__generated__/TopOfChangelogWidget'
import { ArrowRightIcon } from '@heroicons/vue/outline'
import time from '../time'
import { useUpdatedChangesStatuses } from '../subscriptions/useUpdatedChangesStatuses'
import ChangelogChange, { CHANGELOG_CHANGE_FRAGMENT } from '../molecules/ChangelogChange.vue'

export const TOP_OF_CHANGELOG = gql`
  fragment TopOfChangelog on Codebase {
    id
    shortID
    changes(input: { limit: 4 }) {
      ...ChangelogChange
    }
  }
  ${CHANGELOG_CHANGE_FRAGMENT}
`

type Change = TopOfChangelogFragment['changes'][number]

export default defineComponent({
  name: 'TopOfChangelogWidget',
  components: { ArrowRightIcon, ChangelogChange },
  props: {
    codebase: {
      type: Object as PropType<TopOfChangelogFragment>,
      required: true,
    },
  },
  setup(props) {
    const { codebase } = toRefs(props)
    const changeIDs = codebase.value.changes.map((c) => c.id)
    useUpdatedChangesStatuses(changeIDs)
  },
  computed: {
    latestChange(): Change {
      return this.codebase.changes[0]
    },
  },
  methods: {
    timeAgo(timestamp: number) {
      return time.getRelativeTime(new Date(timestamp * 1000))
    },
  },
})
</script>
