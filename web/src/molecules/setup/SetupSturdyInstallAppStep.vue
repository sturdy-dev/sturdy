<template>
  <div class="flex rounded-md mb-1 mr-1">
    <RouterLinkButton
      :to="{
        name: 'download',
      }"
    >
      <DownloadIcon class="w-4 h-4 mr-2" />
      Download and install Sturdy
    </RouterLinkButton>
  </div>
</template>
<script lang="ts">
import { DownloadIcon } from '@heroicons/vue/solid'
import RouterLinkButton from '../../components/shared/RouterLinkButton.vue'
import { Slug } from '../../slug'
import { defineComponent } from 'vue'

export default defineComponent({
  components: {
    DownloadIcon,
    RouterLinkButton,
  },
  props: ['codebase'],
  computed: {
    codebaseSlug() {
      return Slug(this.codebase.name, this.codebase.shortID)
    },
  },
})
</script>
