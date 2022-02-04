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
import { onMounted, onUnmounted, defineComponent } from 'vue'
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
    onMounted(() => {
      document.body.classList.add('bg-slate-900')
    })
    onUnmounted(() => {
      document.body.classList.remove('bg-slate-900')
    })
    return {
      user: null,
    }
  },
})
</script>
