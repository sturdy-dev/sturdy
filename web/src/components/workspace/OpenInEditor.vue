<template>
  <div>
    <link-button
      :href="'vscode://file/' + view.mountPath"
      target="_blank"
      class="inline-flex items-center gap-1"
    >
      <VisualStudioCode class="h-5 w-5" />
      <span class="hidden sm:block">Open</span>
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
