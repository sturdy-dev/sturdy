<template>
  <PaddedApp class="flex flex-col md:flex-row md:space-x-4 space-y-8 md:space-y-0">
    <div class="md:w-60">
      <slot name="navigation" />
    </div>
    <div class="md:flex-1 md:max-w-7xl">
      <slot />
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
