<template v-if="show">
  <Tooltip y-direction="down" x-direction="left">
    <template #tooltip>
      <p>
        <b>Search </b>
        <label v-if="darwin || linux"> ⌘ + {{ keyToSearch }} </label>
        <label v-else> Ctrl + {{ keyToSearch }} </label> or /
      </p>
    </template>
    <button
      v-if="show"
      class="h-7 w-7 p-1 flex disabled:opacity-30 cursor-auto text-gray-500 disabled:bg-transparent disabled:text-gray-500 hover:bg-gray-200 hover:text-gray-600 rounded-md items-center justify-center"
      @click="open"
    >
      <SearchIcon class="w-5" />
    </button>
  </Tooltip>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import { SearchIcon } from '@heroicons/vue/solid'
import Tooltip from '../atoms/Tooltip.vue'

export default defineComponent({
  components: { Tooltip, SearchIcon },
  data() {
    const { appEnvironment } = window
    const ipc = window.ipc
    return {
      ipc,
      appEnvironment,
      show: false,
    }
  },
  computed: {
    isApp() {
      return !!this.ipc
    },
    darwin() {
      return this.appEnvironment?.platform === 'darwin'
    },
    windows() {
      return this.appEnvironment?.platform === 'win32'
    },
    linux() {
      return this.appEnvironment?.platform === 'linux'
    },
    keyToSearch() {
      return this.isApp ? 'F' : 'K'
    },
  },
  mounted() {
    this.emitter.on('search-toolbar-button-visible', this.searchToolbarButtonVisible)
  },
  unmounted() {
    this.emitter.on('search-toolbar-button-visible', this.searchToolbarButtonVisible)
  },
  methods: {
    open() {
      this.emitter.emit('open-search-toolbar')
    },
    searchToolbarButtonVisible(visible: boolean) {
      this.show = visible
    },
  },
})
</script>
