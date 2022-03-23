<template>
  <div class="relative group">
    <slot></slot>

    <div
      v-if="!disabled"
      :class="[
        yDirection === 'up'
          ? 'bottom-full -translate-y-1 group-hover:-translate-y-2'
          : 'top-full translate-y-1 group-hover:translate-y-2',
        xDirection === 'left' ? 'right-0' : 'left-0',
      ]"
      class="absolute flex flex-col transition transition-all group-hover:delay-500 opacity-0 transform group-hover:opacity-100 z-20 pointer-events-none"
    >
      <span
        class="relative z-10 px-2 py-1.5 text-xs rounded text-white bg-black shadow-lg w-max max-w-md"
      >
        <slot name="tooltip"></slot>
      </span>
      <div
        :class="[
          yDirection === 'up' ? 'top-full -translate-y-2' : 'bottom-full translate-y-2',
          xDirection === 'left' ? 'right-2' : 'left-2',
        ]"
        class="w-3 h-3 transform rotate-45 bg-black absolute"
      ></div>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import type { PropType } from 'vue'

export default defineComponent({
  name: 'Tooltip',
  props: {
    disabled: Boolean,
    yDirection: {
      type: String as PropType<'up' | 'down'>,
      default: 'up',
    },
    xDirection: {
      type: String as PropType<'left' | 'right'>,
      default: 'right',
    },
  },
})
</script>
