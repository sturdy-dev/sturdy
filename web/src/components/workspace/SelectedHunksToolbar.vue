<template>
  <transition
    enter-active-class="transition ease-out duration-100"
    enter-from-class="transform opacity-0 -translate-y-1"
    enter-to-class="transform opacity-100 translate-y-none"
    leave-active-class="transition ease-in duration-75"
    leave-from-class="transform opacity-100 translate-y-none"
    leave-to-class="transform opacity-0 -translate-y-1"
  >
    <div
      v-if="selectedHunkIDs.size > 0"
      class="bg-blue-600 text-white shadow-lg flex flex-row items-center border-b border-blue-700 px-5"
    >
      <div class="text-blue-100 ml-3">
        {{ selectedHunkIDs.size }} change{{ selectedHunkIDs.size === 1 ? '' : 's' }} selected
      </div>
      <div class="flex-1" />
      <div class="flex flex-row divide-x divide-blue-700">
        <button
          class="flex flex-row gap-2 items-center px-3 py-2 focus:outline-none focus:bg-blue-500 hover:bg-blue-500 active:bg-blue-700"
          @click="copyToNewWorkspace"
        >
          <DocumentDuplicateIcon class="h-4" />
          <span class="flex-1 whitespace-nowrap">Copy to New Workspace</span>
        </button>

        <button
          class="flex flex-row gap-2 items-center px-3 py-2 focus:outline-none focus:bg-blue-500 hover:bg-blue-500 active:bg-blue-700"
          @click="undo"
        >
          <TrashIcon class="h-4" />
          <span class="flex-1 whitespace-nowrap">Undo</span>
        </button>

        <button
          class="flex flex-row gap-2 items-center px-3 py-2 focus:outline-none focus:bg-blue-500 hover:bg-blue-500 active:bg-blue-700"
          @click="deselect"
        >
          <XIcon class="h-4" />
          <span class="flex-1 whitespace-nowrap">Unselect</span>
        </button>
      </div>
    </div>
  </transition>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import { DocumentDuplicateIcon, TrashIcon, XIcon } from '@heroicons/vue/outline'

export default defineComponent({
  components: { XIcon, TrashIcon, DocumentDuplicateIcon },
  data() {
    return {
      selectedHunkIDs: new Set(),
    }
  },
  mounted() {
    this.emitter.on('differ-selected-hunk-ids', this.onSelectedHunkIDs)
  },
  unmounted() {
    this.emitter.off('differ-selected-hunk-ids', this.onSelectedHunkIDs)
  },
  methods: {
    deselect() {
      this.emitter.emit('differ-deselect-all-hunks', {})
    },
    undo() {
      this.emitter.emit('undo-selected')
    },
    copyToNewWorkspace() {
      this.emitter.emit('copy-selected-to-new-workspace')
    },
    onSelectedHunkIDs(hunkIDs: Set<string>) {
      this.selectedHunkIDs = hunkIDs
    },
  },
})
</script>
