<template>
  <div class="flex flex-col gap-4">
    <p class="text-gray-500">
      Sturdy Desktop on Linux is coming soon. Until then, use Sturdy in the browser.
    </p>

    <router-link
      :to="{ name: 'login' }"
      class="appearance-none flex gap-1 items-center rounded-md bg-gradient-to-b from-green-500 to-green-600 border border-green-600 text-green-50 font-medium shadow-md py-2 px-3 hover:from-green-400 hover:to-green-500 hover:border-green-500"
    >
      Use Sturdy in the browser
    </router-link>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'

export default defineComponent({ components: {} })
</script>
