<template>
  <HorizontalDivider bg="bg-white">
    <template #default>Resolve Conflicts</template>
  </HorizontalDivider>
  <p class="text-sm text-gray-500 pb-4 text-center">
    You have a conflict between the workspace and the trunk, pick the version of the code that you
    want to keep.
  </p>

  <div v-for="cf in rebasing.conflictingFiles" :key="cf.path">
    <DiffConflict
      :conflict="cf"
      :live-diffs="conflictDiffs.filter((d) => d.origName === cf.path)"
      :file-path="cf.path"
      @resolve-conflict="resolveConflict"
    />
  </div>
</template>

<script lang="ts">
import { defineComponent, PropType } from 'vue'
import DiffConflict, {DIFF_CONFLICT_DIFF} from '../differ/DiffConflict.vue'
import HorizontalDivider from '../shared/HorizontalDivider.vue'
import { gql } from '@urql/vue'
import { ResolveConflictDiffFragment } from './__generated__/ResolveConflict'

export const RESOLVE_CONFLICT_DIFF = gql`
  fragment ResolveConflictDiff on FileDiff {
    id

    origName
    newName
    preferredName

    isDeleted
    isNew
    isMoved

    hunks {
      id
      patch

      isOutdated
      isApplied
      isDismissed
    }
    ...DiffConflictDiff
  }
  ${DIFF_CONFLICT_DIFF}
`

export default defineComponent({
  components: { HorizontalDivider, DiffConflict },
  emits: ['resolve-conflict'],
  props: {
    rebasing: {
      type: Object as PropType<any>,
      required: true,
    },
    conflictDiffs: {
      type: Object as PropType<Array<ResolveConflictDiffFragment>>,
      required: true,
    },
  },
  methods: {
    resolveConflict(event) {
      let conflict = event.conflictingFile
      let version = event.version

      this.$emit('resolve-conflict', {
        conflictingFile: conflict,
        version: version,
      })
    },
  },
})
</script>
