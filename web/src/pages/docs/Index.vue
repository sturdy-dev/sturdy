<template>
  <Documentation>
    <template #default>
      <div class="prose p-4 max-w-[800px]">
        <h1 id="quick-start">All pages</h1>
      </div>

      <div class="p-4 max-w-[800px]">
        <DocsSidebar />
      </div>
    </template>
  </Documentation>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import Documentation from '../../layouts/Documentation.vue'
import { useHead } from '@vueuse/head'
import DocsSidebar from '../../organisms/docs/DocsSidebar.vue'

export default defineComponent({
  components: { DocsSidebar, Documentation },
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
          content: 'Index of the Sturdy Documentation',
        },
        {
          name: 'keywords',
          content: 'study learn documentation index',
        },
      ],
      title: 'Documentation Index | Sturdy',
    })
  },
})
</script>
