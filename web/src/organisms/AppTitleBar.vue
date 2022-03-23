<template>
  <div v-if="!isApp">
    <slot />
  </div>

  <div v-else class="h-full flex flex-col">
    <header
      class="titlebar flex items-center bg-gray-50"
      :style="{
        height: height,
      }"
    >
      <!-- this is a left side, above the sidebar -->
      <div
        class="h-full w-64 flex titlebar border-b"
        :class="{
          'md:bg-gray-200 md:border-r md:border-b-0 border-gray-300': showSidebar,
          'md:border-r-0': windows && showSidebar,

          'border-b': !showSidebar,
        }"
      >
        <!-- reserve space for traffic lights on mac os -->
        <div v-if="darwin" class="w-20" />
        <AppNavigationButtons class="p-2" />
      </div>

      <!-- this is the rest of the title bar -->
      <div class="h-full flex flex-1 items-center border-b titlebar">
        <OpenSearchToolbarButton class="px-2" />
        <AppMutagenStatus class="grow bg-gray-50" />
        <AppShareButton class="px-2" />

        <WindowsControls v-if="isFrameless && (windows || linux)" />
        <!-- reserve space for traffic lights on windows -->
        <div v-if="!isFrameless && windows" class="w-32" />
      </div>
    </header>

    <div class="overflow-auto" :style="{ height: mainHeight }">
      <slot />
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'

import AppShareButton from '../molecules/AppShareButton.vue'
import AppNavigationButtons from '../molecules/AppNavigationButtons.vue'
import AppMutagenStatus from '../molecules/AppMutagenStatus.vue'
import WindowsControls from '../molecules/WindowsControls.vue'
import OpenSearchToolbarButton from '../molecules/OpenSearchToolbarButton.vue'

export default defineComponent({
  components: {
    OpenSearchToolbarButton,
    AppShareButton,
    AppNavigationButtons,
    AppMutagenStatus,
    WindowsControls,
  },
  props: {
    showSidebar: {
      type: Boolean,
      required: true,
    },
  },
  data() {
    const { appEnvironment } = window
    return {
      appEnvironment,
    }
  },
  computed: {
    isApp() {
      return !!this.appEnvironment
    },
    isFrameless() {
      return this.appEnvironment?.frameless ?? false
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
    height() {
      return this.darwin ? '3rem' : 'calc(2rem - 1px)'
    },
    mainHeight() {
      return this.darwin ? 'calc(100vh - 3rem)' : 'calc(100vh - 2rem - 1px)'
    },
  },
})
</script>

<style>
.titlebar {
  -webkit-app-region: drag;
  -webkit-user-select: none;
  user-select: none;
}

.titlebar button {
  -webkit-app-region: no-drag;
}
</style>
