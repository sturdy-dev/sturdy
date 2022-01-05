<template>
  <Tooltip y-direction="down" x-direction="left">
    <template #tooltip>Copy a link to the current page to clipboard</template>

    <button
      class="h-7 w-7 p-1 flex disabled:opacity-30 cursor-auto text-gray-500 disabled:bg-transparent disabled:text-gray-500 hover:bg-gray-200 hover:text-gray-600 rounded-md items-center justify-center"
      @click="share"
    >
      <LinkIcon class="w-5" />
    </button>
  </Tooltip>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import { LinkIcon } from '@heroicons/vue/outline'
import Tooltip from './shared/Tooltip.vue'

export default defineComponent({
  components: { Tooltip, LinkIcon },
  methods: {
    share() {
      navigator.clipboard.writeText(location.href)
      this.emitter.emit('notification', {
        title: 'Link Copied',
      })
    },
  },
})
</script>
