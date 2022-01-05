<template>
  <HorizontalDivider bg="bg-white">
    <template #default>Resolve Conflicts</template>
    <template #right>
      Step {{ rebasing.progress_current + 1 }} of {{ rebasing.progress_total }}
    </template>
  </HorizontalDivider>

  <p class="text-sm text-gray-500 pb-4 text-center">
    You have a conflict between the workspace and the trunk, pick the version of the code that you
    want to keep.
  </p>
  <div v-for="cf in rebasing.conflicting_files" :key="cf.path">
    <DiffConflict
      :conflict="cf"
      :live-diffs="conflictDiffs.filter((d) => d.orig_name === cf.path)"
      :file-path="cf.path"
      @resolveConflict="resolveConflict"
    />
  </div>
</template>

<script lang="ts">
import { defineComponent, PropType } from 'vue'
import DiffConflict from '../differ/DiffConflict.vue'
import HorizontalDivider from '../shared/HorizontalDivider.vue'

export default defineComponent({
  components: { HorizontalDivider, DiffConflict },
  emits: ['resolveConflict'],
  data() {
    return {}
  },
  props: {
    rebasing: {
      type: Object as PropType<any>,
      required: true,
    },
    conflictDiffs: {
      required: false,
    },
  },
  methods: {
    resolveConflict(event) {
      let conflict = event.conflictingFile
      let version = event.version

      this.$emit('resolveConflict', {
        conflictingFile: conflict,
        version: version,
      })
    },
  },
})
</script>
