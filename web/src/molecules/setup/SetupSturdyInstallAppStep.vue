<template>
  <div class="flex rounded-md mb-1 mr-1">
    <LinkButton href="https://getsturdy.com/download" target="_blank">
      <DownloadIcon class="w-4 h-4 mr-2" />
      Download and install Sturdy
    </LinkButton>
  </div>
</template>
<script lang="ts">
import { DownloadIcon } from '@heroicons/vue/solid'
import LinkButton from '../../atoms/LinkButton.vue'
import { defineComponent } from 'vue'

export default defineComponent({
  components: {
    DownloadIcon,
    LinkButton,
  },
})
</script>
