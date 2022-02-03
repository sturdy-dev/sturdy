<template>
  <a :class="[light ? 'text-gray-800 hover:text-gray-600' : 'text-slate-300 hover:text-slate-100']">
    <slot></slot>
  </a>
</template>

<script lang="ts">
import { defineComponent } from 'vue'

export default defineComponent({
  props: {
    light: {
      type: Boolean,
      required: true,
    },
  },
})
</script>
