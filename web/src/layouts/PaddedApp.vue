<template>
  <div class="p-4 sm:p-8">
    <FirstTimeUserNoNameTakeover v-if="data && data.user && !data.user.name" :user="data.user" />
    <slot v-else></slot>
  </div>
</template>

<script lang="ts">
import { computed, defineComponent, inject, ref, Ref } from 'vue'
import { gql, useQuery } from '@urql/vue'
import FirstTimeUserNoNameTakeover from '../components/user/FirstTimeUserNoNameTakeover.vue'
import { Feature } from '../__generated__/types'

export default defineComponent({
  components: { FirstTimeUserNoNameTakeover },
  setup() {
    const features = inject<Ref<Array<Feature>>>('features', ref([]))
    const isMultiTenancyEnabled = computed(() => features?.value?.includes(Feature.MultiTenancy))
    const isLicenseEnabled = computed(() => features?.value?.includes(Feature.License))

    const { data, error } = useQuery({
      query: gql`
        query PaddedApp($isMultiTenancyEnabled: Boolean!, $isLicenseEnabled: Boolean!) {
          user {
            id
            name
          }

          installation @skip(if: $isMultiTenancyEnabled) {
            id
            needsFirstTimeSetup
            version

            license @include(if: $isLicenseEnabled) {
              id
            }
          }
        }
      `,
      requestPolicy: 'cache-and-network',
      variables: {
        isMultiTenancyEnabled,
        isLicenseEnabled,
      },
    })

    return {
      displaySelfHostedBanner: !isMultiTenancyEnabled.value,
      data,
      error,
      features,
    }
  },
})
</script>
