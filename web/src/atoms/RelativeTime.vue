<template>
  <time :datetime="date.toISOString()" :title="date.toLocaleString()">{{ ago }}</time>
</template>

<script lang="ts">
import time from '../time'
import { defineComponent } from 'vue'

export default defineComponent({
  props: {
    date: { type: Date, required: true },
  },
  data() {
    return {
      from: new Date(),
      interval: null as null | ReturnType<typeof setInterval>,
    }
  },
  mounted() {
    this.interval = setInterval(() => {
      this.from = new Date()
    }, 1000)
  },
  unmounted() {
    if (this.interval) clearInterval(this.interval)
  },
  computed: {
    ago(): string {
      return time.getRelativeTime(this.date, this.from)
    },
  },
})
</script>
