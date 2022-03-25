<template>
  <div>
    <div class="fixed z-50 md:pr-64 w-full">
      <slot name="toolbar" />
    </div>

    <PaddedApp>
      <div class="grid grid-cols-1 xl:grid-cols-4">
        <div class="xl:col-span-3 xl:pr-8 xl:border-r xl:border-gray-200">
          <slot />
        </div>
        <div class="xl:pl-8 hidden xl:block">
          <slot name="sidebar" />
        </div>
      </div>
    </PaddedApp>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import PaddedApp from './PaddedApp.vue'

export default defineComponent({
  components: { PaddedApp },
})
</script>
