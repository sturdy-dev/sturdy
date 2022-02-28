<template>
  <time :datetime="date.toISOString()" :title="date.toLocaleString()">{{ ago }}</time>
</template>

<script lang="ts">
import time from '../time'

export default {
  props: {
    date: {
      type: Date,
      required: true,
    },
  },
  computed: {
    ago() {
      return time.getRelativeTime(this.date)
    },
  },
}
</script>
