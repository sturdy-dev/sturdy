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
import { defineComponent } from 'vue'

export default defineComponent({
  data() {
    const { ipc, mutagenIpc } = window
    return {
      appState: 'starting',
      restarting: false,
      ipc,
      mutagenIpc,
      interval: undefined as undefined | ReturnType<typeof setInterval>,
    }
  },
  unmounted() {
    if (this.interval) clearInterval(this.interval)
  },
  mounted() {
    this.fetchSetState()
    this.interval = setInterval(this.fetchSetState, 1000)
  },
  methods: {
    async fetchSetState() {
      // This API was added on 2021-11-24
      if (this.ipc && this.ipc.state) {
        this.appState = await this.ipc.state()
      }
      // This API was removed on 2021-11-24
      else if (this.mutagenIpc && this.mutagenIpc.isReady) {
        if (await this.mutagenIpc.isReady()) {
          this.appState = 'online'
        } else {
          this.appState = 'starting'
        }
      } else {
        this.appState = 'online'
      }
    },
    async restart() {
      try {
        this.restarting = true
        await this.ipc.forceRestartMutagen()
      } catch (e) {
        console.error(e)
      } finally {
        this.restarting = false
      }
    },
  },
})
</script>
