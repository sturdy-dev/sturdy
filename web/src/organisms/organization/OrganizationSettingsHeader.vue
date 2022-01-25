<template>
  <div class="relative">
    <div class="absolute inset-0 flex items-center" aria-hidden="true">
      <div class="w-full border-t border-gray-300" />
    </div>

    <div class="relative flex items-center justify-between h-10">
      <span class="pr-2 bg-white text-sm text-gray-500">
        {{ name }}
      </span>
      <div v-if="hasDefaultSlot" class="pl-2 bg-white">
        <slot></slot>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'

export default defineComponent({
  props: {
    name: {
      type: String,
      required: true,
    },
  },
  computed: {
    hasDefaultSlot() {
      return !!this.$slots.default
    },
  },
})
</script>
