<template>
  <div>
    <ClientOnly>
      <StickyTop>
        <NavigationHeader :user="user" :light="false" />
      </StickyTop>
    </ClientOnly>

    <slot></slot>

    <NavigationFooter :light="false" />
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import { ClientOnly } from 'vite-ssr/vue'
import NavigationFooter from '../organisms/NavigationFooter.vue'
import NavigationHeader from '../organisms/NavigationHeader.vue'
import StickyTop from '../molecules/StickyTop.vue'

export default defineComponent({
  components: {
    StickyTop,
    NavigationFooter,
    NavigationHeader,
    ClientOnly,
  },
  setup() {
    return {
      user: null,
    }
  },
})
</script>
