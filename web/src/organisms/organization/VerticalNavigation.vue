<template>
  <VerticalNavigation :navigation="navigation" />
</template>

<script lang="ts">
import VerticalNavigation from '../VerticalNavigation.vue'
import { defineComponent, inject, ref, Ref } from 'vue'
import { useRoute } from 'vue-router'
import { Feature } from '../../__generated__/types'

export default defineComponent({
  components: { VerticalNavigation },
  setup() {
    let route = useRoute()

    const features = inject<Ref<Array<Feature>>>('features', ref([]))
    const isGitHubEnabled = features.value.includes(Feature.GitHub)

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

      isGitHubEnabled
        ? {
            name: 'GitHub',
            linkName: 'organizationSettingsGitHub',
            current: route.name === 'organizationSettingsGitHub',
          }
        : null,

      {
        name: 'Subscriptions',
        linkName: 'organizationCreateSubscription',
        current: route.name === 'organizationCreateSubscription',
      },
    ].filter((nav) => nav)

    return {
      navigation,
    }
  },
})
</script>
