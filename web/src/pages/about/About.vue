<template>
  <div class="bg-gray-100">
    <div>
      <div class="mx-auto py-12 px-4 max-w-7xl sm:px-6 lg:px-8 lg:py-24 space-y-24">
        <Team />
        <Press />
      </div>
    </div>
  </div>
</template>

<script>
import Press from './Press.vue'
import Team from './Team.vue'
import { useHead } from '@vueuse/head'

export default {
  name: 'About',
  components: { Press, Team },
  setup() {
    useHead({
      title: 'About Sturdy',
    })
  },
}
</script>

<style scoped></style>
