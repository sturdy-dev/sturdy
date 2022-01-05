<template>
  <div class="">
    <ClientOnly>
      <IndexNavbar :user="user" />
    </ClientOnly>

    <slot></slot>
    <IndexFooter />
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import IndexNavbar from '../components/IndexNavbar.vue'
import IndexFooter from '../components/IndexFooter.vue'
import { ClientOnly } from 'vite-ssr/vue'

export default defineComponent({
  components: {
    IndexNavbar,
    IndexFooter,
    ClientOnly,
  },
  setup() {
    return {
      user: null,
    }
  },
})
</script>
