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
        name: 'Codebases',
        linkName: 'organizationListCodebases',
        current: route.name === 'organizationListCodebases',
      },
      {
        name: 'Settings',
        linkName: 'organizationSettings',
        current: route.name === 'organizationSettings',
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
