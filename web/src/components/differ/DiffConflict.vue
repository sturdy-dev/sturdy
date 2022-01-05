<template>
  <div
    class="d2h-file-wrapper bg-white lg:rounded-md border border-gray-200 my-4 z-0 relative overflow-y-hidden overflow-x-auto"
  >
    <DiffHeader
      class="cursor-pointer"
      :conflict-selection="selected"
      :diffs="trunkDiff"
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
            class=""
            :unparsed-diff="trunkDiff.hunks[0]"
          />
        </div>

        <div class="ml-4 flex flex-col flex-grow w-1/2 overflow-x-auto">
          <DiffTable
            :class="selected === 'trunk' ? 'opacity-30' : ''"
            class=""
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
import { defineAsyncComponent } from 'vue'

interface Data {
  selected: string | null
  expanded: boolean
}

export default {
  name: 'DiffConflict',
  components: {
    Button,
    DiffHeader,
    DiffTable: defineAsyncComponent(() => import('../differ/DiffTable.vue')),
  },
  props: {
    conflict: Object,
    liveDiffs: Object,
    filePath: String,
  },
  emits: ['resolveConflict'],
  data(): Data {
    return {
      selected: 'todo',
      expanded: true,
    }
  },
  computed: {
    trunkDiff: function () {
      return this.conflict.trunk_diff
    },
    workspaceDiff: function () {
      return this.conflict.workspace_diff
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
      this.$emit('resolveConflict', {
        conflictingFile: conflictingFile,
        version: version,
      })
    },
  },
}
</script>

<style scoped></style>
