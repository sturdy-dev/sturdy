<template>
  <div>
    <link-button
      :href="'vscode://file/' + view.mountPath"
      target="_blank"
      class="inline-flex items-center space-x-1"
    >
      <span>Open in</span>
      <VisualStudioCode class="h-5 w-5" />
    </link-button>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import LinkButton from '../shared/LinkButton.vue'
import VisualStudioCode from './icons/VisualStudioCode.vue'

export default defineComponent({
  components: { VisualStudioCode, LinkButton },
  props: ['view'],
})
</script>
