<template>
  <Documentation>
    <div>
      <div class="mx-auto py-12 px-4 max-w-7xl sm:px-6 lg:px-8 lg:py-24 space-y-24">
        <Team />
        <Press />
      </div>
    </div>
  </Documentation>
</template>

<script lang="ts" setup>
import Press from './Press.vue'
import Team from './Team.vue'
import { useHead } from '@vueuse/head'
import Documentation from '../../layouts/Documentation.vue'

useHead({
  title: 'About Sturdy',
})
</script>
