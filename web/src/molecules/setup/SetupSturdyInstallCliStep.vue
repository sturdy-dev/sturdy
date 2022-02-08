<template>
  <div class="flex rounded-md mb-1 mr-1">
    <RouterLinkButton
      :to="{
        name: 'installClient',
        params: { codebaseSlug: codebaseSlug },
      }"
    >
      <DownloadIcon class="w-4 h-4 mr-2" />
      Install CLI
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
