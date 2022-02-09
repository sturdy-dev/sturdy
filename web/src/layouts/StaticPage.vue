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
import { defineComponent } from 'vue'
import DocumentationWithTableOfContents from './DocumentationWithTableOfContents.vue'
import DocsSidebar from '../organisms/docs/DocsSidebar.vue'

export default defineComponent({
  components: { DocumentationWithTableOfContents, DocsSidebar },
  props: {
    title: String,
    subtitle: String,
    category: String,
    narrow: Boolean,
    wide: Boolean,
    metadescription: String,
    image: {
      type: String,
      default() {
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
    imageColClasses() {
      if (this.narrow) {
        return 'lg:col-start-2'
      }
      if (this.wide) {
        return '' // no image
      }
      return 'lg:col-start-3'
    },
    textColClasses() {
      // :class="[narrow ? 'lg:col-end-2' : 'lg:col-end-3']"
      if (this.narrow) {
        return 'lg:col-end-2'
      }
      if (this.wide) {
        return 'lg:col-end-1'
      }
      return 'lg:col-end-3'
    },
  },
})
</script>
