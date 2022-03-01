<template>
  <div class="flow-root">
    <ul role="list" class="-mb-8">
      <li v-for="(c, itemIdx) in changes" :key="c.id">
        <div class="relative pb-8">
          <span
            v-if="itemIdx !== changes.length - 1"
            class="absolute top-5 left-5 -ml-px h-full w-0.5 bg-gray-200"
            aria-hidden="true"
          />

          <ChangelogEntry
            :is-selected="c.id === selectedChangeId"
            :change="c"
            @click="selectChange(c.id)"
            @select="selectChange(c.id)"
          />
        </div>
      </li>
    </ul>
  </div>
</template>

<script lang="ts">
import ChangelogEntry from './ChangelogEntry.vue'
import { Change } from '../../__generated__/types'
import { PropType, toRefs } from 'vue'
import { useUpdatedChangesStatuses } from '../../subscriptions/useUpdatedChangesStatuses'

export default {
  name: 'CodebaseChangelogSidebar',
  components: { ChangelogEntry },
  props: {
    codebaseId: {
      type: String,
      required: true,
    },
    selectedChangeId: {
      type: String,
      required: true,
    },
    changes: {
      type: Array as PropType<Change[]>,
      required: true,
    },
  },
  emits: ['selectCodebaseChange'],
  setup(props) {
    const { changes } = toRefs(props)
    const changeIDs = changes.value.map((c) => c.id)
    // watch top 50 changes, otherwise, message is too big
    useUpdatedChangesStatuses(changeIDs.slice(0, 50))
  },
  methods: {
    selectChange(commitID: string) {
      this.$emit('selectCodebaseChange', {
        commit_id: commitID,
      })
    },
  },
}
</script>
