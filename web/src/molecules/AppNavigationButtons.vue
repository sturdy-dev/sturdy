<template>
  <div class="flex items-center gap-0.5">
    <button
      :disabled="!canGoBack"
      class="h-7 w-7 p-1 flex disabled:opacity-30 cursor-auto text-gray-500 disabled:bg-transparent disabled:text-gray-500 hover:bg-gray-300 hover:text-gray-600 rounded-md items-center justify-center"
      @click="ipc.goBack()"
    >
      <ArrowLeftIcon />
    </button>
    <button
      :disabled="!canGoForward"
      class="h-7 w-7 p-1 flex disabled:opacity-30 cursor-auto text-gray-500 disabled:bg-transparent disabled:text-gray-500 hover:bg-gray-300 hover:text-gray-600 rounded-md items-center justify-center"
      @click="ipc.goForward()"
    >
      <ArrowRightIcon />
    </button>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'

import { ArrowRightIcon, ArrowLeftIcon } from '@heroicons/vue/solid'

export default defineComponent({
  components: { ArrowRightIcon, ArrowLeftIcon },

  data() {
    const { ipc } = window
    return {
      ipc,
      canGoBack: false,
      canGoForward: false,
    }
  },
  watch: {
    $route: {
      immediate: true,
      handler: function () {
        this.ipc.canGoBack().then((can: boolean) => (this.canGoBack = can))
        this.ipc.canGoForward().then((can: boolean) => (this.canGoForward = can))
      },
    },
  },
})
</script>
