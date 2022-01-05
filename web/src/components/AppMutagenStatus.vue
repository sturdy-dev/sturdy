<template>
  <div class="flex gap-2 items-center justify-center">
    <template v-if="appState === 'online'"></template>

    <template v-else-if="appState === 'offline'">
      <span class="text-sm text-gray-800">Disconnected!</span>
      <button
        :disabled="restarting"
        class="disabled:opacity-70 text-sm px-2 py-0.5 bg-blue-600 text-white rounded-md border-blue-700"
        @click="restart"
      >
        Reconnect
      </button>
    </template>

    <template v-else>
      <span class="text-sm text-gray-800">
        <span v-if="appState === 'starting'">Starting</span>
        <span v-else-if="appState === 'creating-ssh-key'">First time setup</span>
        <span v-else-if="appState === 'uploading-ssh-key'">Authorizing Device</span>
        <span v-else>Loading</span>
      </span>

      <div class="inline-flex space-x-1 rounded-full items-center -ml-1 mt-2.5">
        <div
          class="bg-gray-500 w-1 h-1 rounded-full animate-bounce"
          style="animation-delay: 0.1s"
        ></div>
        <div
          class="bg-gray-500 w-1 h-1 rounded-full animate-bounce"
          style="animation-delay: 0.2s"
        ></div>
        <div
          class="bg-gray-500 w-1 h-1 rounded-full animate-bounce"
          style="animation-delay: 0.3s"
        ></div>
      </div>
    </template>
  </div>
</template>

<script lang="ts">
import { defineComponent, onUnmounted, ref } from 'vue'

export default defineComponent({
  setup() {
    let appState = ref('starting')
    let restarting = ref(false)

    const fetchSetState = async () => {
      // This API was added on 2021-11-24
      if (window.ipc && window.ipc.state) {
        appState.value = await window.ipc.state()
      }
      // This API was removed on 2021-11-24
      else if (window.mutagenIpc && window.mutagenIpc.isReady) {
        if (await window.mutagenIpc.isReady()) {
          appState.value = 'online'
        } else {
          appState.value = 'starting'
        }
      } else {
        appState.value = 'online'
      }
    }

    fetchSetState()
    let interval = setInterval(fetchSetState, 1000)

    onUnmounted(() => {
      clearInterval(interval)
    })

    return {
      appState,
      restarting,
    }
  },
  methods: {
    async restart() {
      try {
        this.restarting = true
        await window.ipc.forceRestartMutagen()
      } catch (e) {
        console.error(e)
      } finally {
        this.restarting = false
      }
    },
  },
})
</script>
