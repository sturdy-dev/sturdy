<template>
  <div>
    <StickyTop>
      <div class="text-center bg-yeallow bg-slate-300 text-slate-900 p-4">
        Sturdy is no longer under active development as the team focuses on
        <a class="underline text-slate-900" href="https://codeball.ai">something new</a>.
      </div>
      <NavigationHeader :light="false" />
    </StickyTop>

    <slot></slot>

    <NavigationFooter :light="false" />
  </div>
</template>

<script lang="ts" setup>
import { onMounted, onUnmounted } from 'vue'
import NavigationFooter from '../organisms/NavigationFooter.vue'
import NavigationHeader from '../organisms/NavigationHeader.vue'
import StickyTop from '../molecules/StickyTop.vue'

onMounted(() => {
  document.body.classList.add('bg-slate-900')
})

onUnmounted(() => {
  document.body.classList.remove('bg-slate-900')
})
</script>
