<template>
  <div class="flex flex-col">
    <h2 :class="[accent, 'font-semibold', 'highlight', 'tracking-tight']">
      <div class="flex flex-row space-x-4 text-lg items-center">
        <div class="w-8 h-8">
          <slot name="icon"></slot>
        </div>
        <div>
          {{ title }}
        </div>
      </div>
    </h2>
    <div class="mt-5 text-sm">
      <slot></slot>
    </div>
    <div
      class="mt-3 text-sm grow text-slate-400/80 hover:text-slate-50/80 font-semibold tracking-tight inline-flex group inline-flex items-center"
    >
      <router-link :to="{ hash: '#' + title.toLowerCase() }" hash="expressive">
        <span class="inline-flex items-center">
          Go to
          <ChevronDownIcon class="h-4 w-4 ml-1" />
        </span>
      </router-link>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import { ChevronDownIcon } from '@heroicons/vue/outline'
import UspVideo from './UspVideo.vue'
export default defineComponent({
  components: {
    ChevronDownIcon,
  },
  props: {
    title: {
      type: String,
      required: true,
    },
    accent: {
      type: String,
      required: true,
    },
  },
})
</script>

<style scoped>
.fix-icon {
  transform: scale(-1, 1) rotate(-90deg);
}
</style>
