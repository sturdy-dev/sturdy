<template>
  <PaddedApp class="flex flex-col space-y-4">
    <slot name="header"></slot>
    <div class="flex flex-col md:flex-row md:space-x-4 space-y-8 md:space-y-0">
      <div class="w-full md:w-1/4 lg:w-60">
        <slot name="navigation" />
      </div>
      <div class="flex-1">
        <slot />
      </div>
    </div>
  </PaddedApp>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import PaddedApp from './PaddedApp.vue'

export default defineComponent({
  components: { PaddedApp },
})
</script>
