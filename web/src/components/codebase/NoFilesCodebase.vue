<template>
  <div class="text-center">
    <div class="inline-block text-blue-800 bg-blue-50 rounded-full w-32 h-32 p-7">
      <FolderIcon />
    </div>

    <div class="text-gray-500 m-auto text-sm max-w-sm mt-2">
      <p>This codebase doesn't have any files yet.</p>
      <p>Connect a draft change to your computer to get started.</p>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import { FolderIcon } from '@heroicons/vue/outline'

export default defineComponent({
  components: {
    FolderIcon,
  },
})
</script>
