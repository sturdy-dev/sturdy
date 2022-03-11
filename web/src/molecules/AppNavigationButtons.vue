<template>
  <div class="flex items-center gap-0.5">
    <button
      :disabled="!canGoBack"
      class="h-7 w-7 p-1 flex disabled:opacity-30 cursor-auto text-gray-500 disabled:bg-transparent disabled:text-gray-500 hover:bg-gray-300 hover:text-gray-600 rounded-md items-center justify-center"
      @click="back"
    >
      <ArrowLeftIcon />
    </button>
    <button
      :disabled="!canGoForward"
      class="h-7 w-7 p-1 flex disabled:opacity-30 cursor-auto text-gray-500 disabled:bg-transparent disabled:text-gray-500 hover:bg-gray-300 hover:text-gray-600 rounded-md items-center justify-center"
      @click="forward"
    >
      <ArrowRightIcon />
    </button>
  </div>
</template>

<script lang="ts">
import { defineComponent, ref } from 'vue'

import { ArrowRightIcon, ArrowLeftIcon } from '@heroicons/vue/solid'

export default defineComponent({
  components: { ArrowRightIcon, ArrowLeftIcon },
  setup() {
    const { ipc } = window
    const canGoBack = ref(false)
    const canGoForward = ref(false)
    return {
      ipc,
      canGoForward,
      canGoBack,
    }
  },
  watch: {
    $route: function () {
      this.canGoBack = this.ipc.canGoBack().then((canGoBack: boolean) => {
        this.canGoBack = canGoBack
      })
      this.canGoForward = this.ipc.canGoForward().then((canGoForward: boolean) => {
        this.canGoForward = canGoForward
      })
    },
  },
  methods: {
    async back() {
      await this.ipc.goBack()

      this.canGoForward = await this.ipc.canGoForward()
      this.canGoBack = await this.ipc.canGoBack()
    },
    async forward() {
      await this.ipc.goForward()

      this.canGoForward = await this.ipc.canGoForward()
      this.canGoBack = await this.ipc.canGoBack()
    },
  },
})
</script>
