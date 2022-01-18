<template>
  <div class="">
    <ClientOnly>
      <IndexNavbar :user="user" />
    </ClientOnly>

    <slot></slot>
    <DocumentationFooter />
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import IndexNavbar from '../components/IndexNavbar.vue'
import IndexFooter from '../components/IndexFooter.vue'
import { ClientOnly } from 'vite-ssr/vue'
import DocumentationFooter from "../organisms/DocumentationFooter.vue";

export default defineComponent({
  components: {
    DocumentationFooter,
    IndexNavbar,
    ClientOnly,
  },
  setup() {
    return {
      user: null,
    }
  },
})
</script>
