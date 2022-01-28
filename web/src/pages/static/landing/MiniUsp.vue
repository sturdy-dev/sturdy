<template>
  <div class="flex flex-col">
    <h2 class="text-slate-50 font-semibold highlight tracking-tight">
      <div class="flex flex-col space-y-2">
        <div class="w-16 h-16">
          <slot name="icon"></slot>
        </div>
        <div>
          {{ title }}
        </div>
      </div>
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
