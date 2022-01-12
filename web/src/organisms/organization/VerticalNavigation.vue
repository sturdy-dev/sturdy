<template>
  <VerticalNavigation :navigation="navigation" />
</template>

<script lang="ts">
import VerticalNavigation from '../VerticalNavigation.vue'
import { defineComponent } from 'vue'
import { useRoute } from 'vue-router'

export default defineComponent({
  components: { VerticalNavigation },
  setup() {
    let route = useRoute()

    const navigation = [
      {
        name: 'Settings',
        linkName: 'organizationView',
        current: route.name === 'organizationView',
      },
      {
        name: 'Subscriptions',
        linkName: 'organizationCreateSubscription',
        current: route.name === 'organizationCreateSubscription',
      },
    ]

    return {
      navigation,
    }
  },
})
</script>
