<template>
  <div
    v-if="ipc != null"
    class="spacer flex items-center"
    :class="[
      fixed ? 'fixed top-0 z-50' : '',
      padLeft ? 'pad-left' : '',
      padRight ? 'pad-right' : '',
    ]"
    :style="{
      '--min-padding-left': padLeft || '0px',
      '--min-padding-right': padRight || '0px',
    }"
  >
    <slot :ipc="ipc" :app-environment="appEnvironment"></slot>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'

export default defineComponent({
  props: {
    fixed: {
      type: Boolean,
      default: false,
    },
    padRight: {
      type: String,
      default: '',
    },
    padLeft: {
      type: String,
      default: '',
    },
  },
  setup() {
    return {
      ipc: window.ipc,
      appEnvironment: window.appEnvironment,
    }
  },
})
</script>

<style>
.spacer {
  height: calc(env(titlebar-area-height, 2rem) + 1px);
  -webkit-app-region: drag;
}

.spacer button {
  -webkit-app-region: no-drag;
}

.pad-left {
  padding-left: max(var(--min-padding-left), env(titlebar-area-x, 0));
}

.pad-right {
  padding-right: max(
    var(--min-padding-right),
    calc(100vw - env(titlebar-area-width, 100vw) - env(titlebar-area-x, 0))
  );
}
</style>
