<template>
  <span
    class="px-2 inline-flex text-xs leading-5 font-semibold rounded-full whitespace-nowrap"
    :class="classes"
  >
    <slot>Pull</slot>
  </span>
</template>

<script lang="ts">
import { defineComponent } from 'vue'

export default defineComponent({
  props: {
    color: {
      type: String,
      default: 'green',
    },
  },
  computed: {
    classes() {
      return [this.colors]
    },
    colors() {
      if (this.color === 'green') {
        return 'bg-green-100 text-green-800'
      } else if (this.color === 'gray') {
        return 'bg-gray-100 text-gray-800'
      } else if (this.color === 'blue') {
        return 'bg-blue-100 text-blue-800'
      } else {
        return ''
      }
    },
  },
})
</script>
