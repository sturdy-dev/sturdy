<template>
  <div class="h-full flex flex-col">
    <header
      class="titlebar flex items-center bg-gray-50 justify-between bg-gray-50 border-b border-gray-300"
      :style="{
        height: height,
      }"
    >
      <div
        :class="{
          'w-20': darwin,
          'w-32': !darwin,
        }"
      />

      <slot name="header" />

      <WindowsControls v-if="isFrameless && (windows || linux)" />
      <div v-if="!isFrameless && windows" class="w-32" />
      <div v-if="darwin" class="w-20" />
    </header>

    <main class="overflow-auto" :style="{ height: mainHeight }">
      <slot />
    </main>
  </div>
</template>

<script>
import WindowsControls from './WindowsControls.vue'

export default {
  components: { WindowsControls },
  data() {
    const { appEnvironment } = window
    return {
      appEnvironment,
    }
  },
  computed: {
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
    isFrameless() {
      return this.appEnvironment?.frameless
    },
  },
}
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
