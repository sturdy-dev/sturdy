<template>
  <div>
    <NavigationHeader :light="true" />

    <div class="max-w-6xl mx-auto">
      <slot></slot>
    </div>

    <NavigationFooter :light="true" />
  </div>
</template>

<script lang="ts" setup>
import { ClientOnly } from 'vite-ssr/vue'
import NavigationFooter from '../organisms/NavigationFooter.vue'
import NavigationHeader from '../organisms/NavigationHeader.vue'
</script>
