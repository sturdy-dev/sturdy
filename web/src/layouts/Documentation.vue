<template>
  <div class="">
    <NavigationHeader :light="true" />

    <slot></slot>

    <NavigationFooter :light="true" />
  </div>
</template>

<script lang="ts" setup>
import { ClientOnly } from 'vite-ssr/vue'
import NavigationFooter from '../organisms/NavigationFooter.vue'
import NavigationHeader from '../organisms/NavigationHeader.vue'
</script>
