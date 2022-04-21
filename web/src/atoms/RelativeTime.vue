<template>
  <time :datetime="date.toISOString()" :title="date.toLocaleString()">{{ ago }}</time>
</template>

<script lang="ts">
import time from '../time'
import { defineComponent } from 'vue'

export default defineComponent({
  props: {
    date: {
      type: Date,
      required: true,
    },
  },
  computed: {
    ago(): string {
      return time.getRelativeTime(this.date)
    },
  },
})
</script>
