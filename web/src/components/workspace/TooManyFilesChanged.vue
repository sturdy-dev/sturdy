<template>
  <div class="bg-white shadow overflow-hidden sm:rounded-lg my-4">
    <div class="px-4 py-5 sm:px-6">
      <h3 class="text-lg leading-6 font-medium text-gray-900">
        Ouch! You have a lot of changes...
      </h3>
      <p class="mt-1 max-w-2xl text-sm text-gray-500">
        We're sorry, but this workspace have too many changes to show.
      </p>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'

export default defineComponent({
  data() {
    return {}
  },
})
</script>
