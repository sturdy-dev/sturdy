<template>
  <VerticalNavigation :navigation="navigation" />
</template>

<script lang="ts">
import VerticalNavigation from '../VerticalNavigation.vue'
import { computed, defineComponent, inject, ref } from 'vue'
import type { Ref } from 'vue'
import { useRoute } from 'vue-router'
import { Feature } from '../../__generated__/types'

export default defineComponent({
  components: { VerticalNavigation },
  setup() {
    let route = useRoute()

    const features = inject<Ref<Array<Feature>>>('features', ref([]))

    const isOrganizationSubscriptionsEnabled = computed(() =>
      features?.value?.includes(Feature.OrganizationSubscriptions)
    )
    const isSelfHostedLicenseEnabled = computed(() =>
      features?.value?.includes(Feature.SelfHostedLicense)
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

        isOrganizationSubscriptionsEnabled.value
          ? {
              name: 'Subscriptions',
              linkName: 'organizationListSubscription',
              current: route.name === 'organizationListSubscription',
            }
          : null,

        isSelfHostedLicenseEnabled.value
          ? {
              name: 'Manage Server',
              linkName: 'organizationManageInstallation',
              current: route.name === 'organizationManageInstallation',
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
