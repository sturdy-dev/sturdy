<template>
  <VerticalNavigation :navigation="navigation" />
</template>

<script lang="ts">
import VerticalNavigation from '../VerticalNavigation.vue'
import { computed, defineComponent, inject, ref, Ref } from 'vue'
import { useRoute } from 'vue-router'
import { Feature } from '../../__generated__/types'

export default defineComponent({
  components: { VerticalNavigation },
  setup() {
    let route = useRoute()

    const features = inject<Ref<Array<Feature>>>('features', ref([]))
    const isGitHubEnabled = computed(() => features?.value?.includes(Feature.GitHub))
    const isOrganizationSubscriptionsEnabled = computed(() =>
      features?.value?.includes(Feature.OrganizationSubscriptions)
    )

    const navigation = computed(() =>
      [
        {
          name: 'Codebases',
          linkName: 'organizationListCodebases',
          current:
            route.name === 'organizationListCodebases' ||
            route.name === 'organizationCreateCodebase',
        },

        {
          name: 'Settings',
          linkName: 'organizationSettings',
          current: route.name === 'organizationSettings',
        },

        isGitHubEnabled.value
          ? {
              name: 'GitHub',
              linkName: 'organizationSettingsGitHub',
              current: route.name === 'organizationSettingsGitHub',
            }
          : null,

        isOrganizationSubscriptionsEnabled.value
          ? {
              name: 'Subscriptions',
              linkName: 'organizationListSubscription',
              current: route.name === 'organizationListSubscription',
            }
          : null,
      ].filter((nav) => nav)
    )

    return {
      navigation,
    }
  },
})
</script>
