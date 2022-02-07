<template>
  <li class="inline-flex items-center space-x-2">
    <CheckIcon class="h-6 w-6 text-green-400" />
    <span><slot></slot></span>
  </li>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import { CheckIcon } from '@heroicons/vue/solid'

export default defineComponent({
  components: { CheckIcon },
})
</script>
