<template>
  <div>
    <ul v-if="changes.length > 0" class="changes-timeline relative">
      <ChangeLogEntry
        v-for="c in changes"
        :key="c.change_id"
        :is-selected="c.commit_id === selectedCommitID"
        :title="c.title"
        :alt-color="c.is_landed"
        :num-comments="c.num_comments"
        @click="selectChange(c.commit_id)"
        @select="selectChange(c.commit_id)"
      />
    </ul>
    <div v-else class="text-gray-500 mt-2">This workspace has no saved changes yet...</div>
  </div>
</template>

<script>
import ChangeLogEntry from '../change/ChangeLogEntry.vue'
import '../../changelog.css'

export default {
  name: 'WorkspaceHomeSidebar',
  components: { ChangeLogEntry },
  props: {
    workspaceID: String,
    selectedCommitID: String,
    changes: Object,
  },
  methods: {
    selectChange(commitID) {
      this.$emit('selectWorkspaceChange', {
        commit_id: commitID,
        is_head: this.changes[0].commit_id === commitID,
      })
    },
  },
}
</script>
