<template>
  <div class="flex items-center space-x-2">
    <CalendarIcon class="h-5 w-5 text-gray-400" aria-hidden="true" />
    <span v-if="data.workspace.updatedAt > 0" class="text-gray-900 text-sm font-medium">
      Updated {{ friendly_ago(data.workspace.updatedAt) }}
    </span>
    <span v-else class="text-gray-900 text-sm font-medium">
      Created {{ friendly_ago(data.workspace.createdAt) }}
    </span>
  </div>
</template>

<script>
import { CalendarIcon } from '@heroicons/vue/solid'
import time from '../../../time'

export default {
  components: { CalendarIcon },
  props: ['data'],
  methods: {
    friendly_ago(ts) {
      return time.getRelativeTime(new Date(ts * 1000))
    },
  },
}
</script>
