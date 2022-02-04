<template>
  <DocumentationWithTableOfContents>
    <template #sidebar>
      <DocsSidebar />
    </template>
    <template #default>
      <div class="prose p-4 max-w-[800px]">
        <h1 id="working-in-the-open">Self-hosted Sturdy</h1>

        <p>TODO</p>
      </div>
    </template>
  </DocumentationWithTableOfContents>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import DocumentationWithTableOfContents from '../../layouts/DocumentationWithTableOfContents.vue'
import { useHead } from '@vueuse/head'
import DocsSidebar from '../../organisms/docs/DocsSidebar.vue'

export default defineComponent({
  components: { DocsSidebar, DocumentationWithTableOfContents },
  setup() {
    useHead({
      meta: [
        // TODO: Remove when we're launching!
        {
          name: 'robots',
          content: 'noindex',
        },
        {
          name: 'description',
          content: 'Working in the open (within your team)',
        },
        {
          name: 'keywords',
          content: 'study learn documentation working code in the open publicly public',
        },
      ],
      title: 'Self-hosted Sturdy',
    })
  },
})
</script>
