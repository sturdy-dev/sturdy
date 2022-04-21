<template>
  <DocumentationWithTableOfContents>
    <template #sidebar>
      <DocsSidebar />
    </template>
    <template #default>
      <div class="prose p-4 max-w-[800px]">
        <div>
          <h1>
            {{ title }}
          </h1>
          <h4 v-if="subtitle">
            {{ subtitle }}
          </h4>
        </div>

        <slot></slot>
      </div>
    </template>
  </DocumentationWithTableOfContents>
</template>

<script lang="ts">
import { useHead } from '@vueuse/head'
import { defineComponent, type PropType } from 'vue'
import DocumentationWithTableOfContents from './DocumentationWithTableOfContents.vue'
import DocsSidebar from '../organisms/docs/DocsSidebar.vue'

export default defineComponent({
  components: { DocumentationWithTableOfContents, DocsSidebar },
  props: {
    title: {
      type: String,
      required: true,
    },
    subtitle: {
      type: String,
      required: false,
      default: null,
    },
    category: {
      type: String,
      required: false,
      default: null,
    },
    narrow: {
      type: Boolean,
      required: false,
    },
    wide: {
      type: Boolean,
      required: false,
    },
    metadescription: {
      type: String,
      required: true,
    },
    image: {
      type: String,
      default: () => {
        return 'https://images.unsplash.com/photo-1553877522-43269d4ea984?ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&ixlib=rb-1.2.1&auto=format&fit=crop&w=1950&q=80'
      },
    },
  },
  setup(props) {
    useHead({
      title: props.title + ' | Sturdy',
      meta: [
        { property: 'og:title', content: props.title + ' | Sturdy' },
        { property: 'description', content: props.metadescription },
        { property: 'og:description', content: props.metadescription },
      ],
    })
  },
  computed: {
    grid() {
      if (this.narrow) {
        return 'lg:grid-cols-2'
      }
      if (this.wide) {
        return 'lg:grid-cols-1'
      }
      return 'lg:grid-cols-3'
    },
  },
})
</script>
