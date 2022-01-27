<template>
  <div class="flex flex-col">
    <h2
      class="-ml-4 pl-4 border-l-2 border-purple-500 text-slate-50 font-semibold highlight tracking-tight"
    >
      {{ title }}
    </h2>
    <div class="my-4 text-sm">
      <slot></slot>
    </div>
    <div
      class="text-sm grow text-amber-500/80 hover:text-slate-50/80 font-semibold tracking-tight inline-flex group inline-flex items-center"
    >
      <router-link :to="{ hash: '#' + title.toLowerCase() }" hash="expressive">
        Go to ⤵️
      </router-link>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
export default defineComponent({
  props: {
    title: {
      type: String,
      required: true,
    },
  },
})
</script>

<style scoped></style>
