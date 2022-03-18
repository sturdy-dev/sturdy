<template>
  <div
    class="d2h-file-wrapper bg-white lg:rounded-md border border-gray-200 my-4 z-0 relative overflow-y-hidden overflow-x-auto"
  >
    <DiffHeader
      class="cursor-pointer"
      :conflict-selection="selected"
      :diffs="trunkDiff"
      :show-add-button="false"
      @click="toggleExpand"
    />

    <div v-if="expanded">
      <div>
        <div class="m-4 flex space-between">
          <div class="w-full flex flex-grow justify-center">
            <Button
              :color="selected === 'trunk' ? 'blue' : 'white'"
              @click="selectResolution('trunk')"
            >
              Version from trunk
            </Button>
          </div>
          <div class="">
            <Button
              :color="selected === 'custom' ? 'blue' : 'white'"
              @click="selectResolution('custom')"
            >
              Custom
            </Button>
          </div>
          <div class="w-full flex flex-grow justify-center">
            <Button
              :color="selected === 'workspace' ? 'blue' : 'white'"
              @click="selectResolution('workspace')"
            >
              Version from workspace
            </Button>
          </div>
        </div>
      </div>

      <div v-if="selected !== 'custom'" class="flex space-between">
        <div class="flex flex-col flex-grow w-1/2 overflow-x-auto">
          <DiffTable
            :class="selected === 'workspace' ? 'opacity-30' : ''"
            :unparsed-diff="trunkDiff.hunks[0]"
          />
        </div>

        <div class="ml-4 flex flex-col flex-grow w-1/2 overflow-x-auto">
          <DiffTable
            :class="selected === 'trunk' ? 'opacity-30' : ''"
            :unparsed-diff="workspaceDiff.hunks[0]"
          />
        </div>
      </div>
      <div v-else class="overflow-x-auto">
        <div class="flex flex-col">
          <div v-if="liveDiffs && liveDiffs.length > 0">
            <p class="text-sm text-gray-500 pb-2 text-center">
              <span>Resolve this conflict by editing </span>
              <span class="text-gray-700">{{ filePath }}</span>
              <span> on your computer</span>
            </p>
            <div v-for="h in liveDiffs[0].hunks" :key="h" class="mt-2">
              <DiffTable class="" :unparsed-diff="h" />
            </div>
          </div>
          <div v-else>
            <p class="text-sm text-gray-500 pb-4 text-center">
              There are no changes to show for file {{ filePath }}
            </p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import DiffHeader from './DiffHeader.vue'
import Button from '../shared/Button.vue'
import { defineAsyncComponent, defineComponent, PropType } from 'vue'
import { gql } from '@urql/vue'
import { DiffConflictDiffFragment } from './__generated__/DiffConflict'

interface Data {
  selected: string | null
  expanded: boolean
}

export const DIFF_CONFLICT_DIFF = gql`
  fragment DiffConflictDiff on FileDiff {
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
  }
`

export default defineComponent({
  components: {
    Button,
    DiffHeader,
    DiffTable: defineAsyncComponent(() => import('../differ/DiffTable.vue')),
  },
  props: {
    conflict: {
      type: Object,
      required: true,
    },
    liveDiffs: {
      type: Object as PropType<Array<DiffConflictDiffFragment>>,
      required: true,
    },
    filePath: {
      type: String,
      required: true,
    },
  },
  emits: ['resolve-conflict'],
  data(): Data {
    return {
      selected: 'todo',
      expanded: true,
    }
  },
  computed: {
    trunkDiff: function () {
      return this.conflict.trunkDiff
    },
    workspaceDiff: function () {
      return this.conflict.workspaceDiff
    },
  },
  methods: {
    selectResolution(res: string) {
      if (res == this.selected) {
        res = 'todo'
      }
      this.selected = res
      this.resolveConflict(this.conflict, res)
    },
    toggleExpand() {
      this.expanded = !this.expanded
    },
    resolveConflict(conflictingFile: any, version: string) {
      this.$emit('resolve-conflict', {
        conflictingFile: conflictingFile,
        version: version,
      })
    },
  },
})
</script>
