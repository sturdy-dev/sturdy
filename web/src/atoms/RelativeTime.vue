<template>
  <time :datetime="date.toISOString()" :title="date.toLocaleString()">{{ ago }}</time>
</template>

<script lang="ts">
import { getRelativeTimeStrict, getRelativeTime } from '../time'
import { defineComponent } from 'vue'

export default defineComponent({
  props: {
    date: { type: Date, required: true },
    strict: { type: Boolean, default: false },
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
      return this.strict
        ? getRelativeTimeStrict(this.date, this.from)
        : getRelativeTime(this.date, this.from)
    },
  },
})
</script>
