<template>
  <div class="relative py-4 mx-4">
    <div class="absolute inset-0 flex items-center" aria-hidden="true">
      <div class="w-full border-t" :class="lineStyle" />
    </div>
    <div class="relative flex items-center justify-between">
      <span class="pr-2 text-sm" :class="textStyle"><slot></slot></span>
      <span class="pl-2 text-sm" :class="textStyle"><slot name="right"></slot></span>
    </div>
  </div>
</template>

<script>
import { defineComponent } from 'vue'

export default defineComponent({
  name: 'HorizontalDivider',
  props: {
    color: {
      type: String,
      default: 'gray', // or "red"
    },
    bg: {
      type: String,
      default: 'bg-gray-100',
    },
  },
  computed: {
    lineStyle() {
      if (this.color === 'red') {
        return 'border-red-300'
      }
      return 'border-gray-300'
    },
    textStyle() {
      if (this.color === 'red') {
        return ['text-red-500', this.bg]
      }
      return ['text-gray-500', this.bg]
    },
  },
})
</script>
