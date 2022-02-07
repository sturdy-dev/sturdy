<template>
  <router-link
    class="flex-grow block w-full bg-slate-800 border border-slate-200 rounded-md my-4 py-2 text-sm font-semibold text-white text-center hover:bg-gray-900"
  >
    <slot></slot>
  </router-link>
</template>

<script lang="ts">
import { defineComponent } from 'vue'

export default defineComponent({})
</script>
