<template>
  <div class="flex rounded-md mb-1 mr-1">
    <Button
      @click="
        $router.push({
          name: 'installClient',
          params: { codebaseSlug: codebaseSlug },
        })
      "
    >
      <DownloadIcon class="w-4 h-4 mr-2" />
      Install
    </Button>
  </div>
</template>
<script>
import { DownloadIcon } from '@heroicons/vue/solid'
import Button from '../shared/Button.vue'
import { Slug } from '../../slug'

export default {
  name: 'SetupSturdyInstallStep',
  components: {
    DownloadIcon,
    Button,
  },
  props: ['codebase'],
  computed: {
    codebaseSlug() {
      return Slug(this.codebase.name, this.codebase.shortID)
    },
  },
}
</script>
