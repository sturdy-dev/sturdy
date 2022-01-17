<template>
  <div class="p-4 sm:p-8">
    <div v-if="displaySelfHostedBanner" class="bg-green-500 m-4 p-2">
      <strong>Debug: This is a self-hosted instance of Sturdy.</strong>
    </div>

    <FirstTimeSetup v-if="data && data?.serverStatus?.needsFirstTimeSetup" />
    <FirstTimeUserNoNameTakeover v-else-if="data && data.user && !data.user.name" :user="user" />
    <slot v-else></slot>
  </div>
</template>

<script lang="ts">
import { defineComponent, inject, ref, Ref } from 'vue'
import { gql, useQuery } from '@urql/vue'
import FirstTimeUserNoNameTakeover from '../components/user/FirstTimeUserNoNameTakeover.vue'
import FirstTimeSetup from '../organisms/serverstatus/FirstTimeSetup.vue'
import { Feature } from '../__generated__/types'

export default defineComponent({
  components: { FirstTimeUserNoNameTakeover, FirstTimeSetup },
  setup() {
    const features = inject<Ref<Array<Feature>>>('features', ref([]))
    const isMultiTenancyEnabled = features.value.includes(Feature.MultiTenancy)

    const { data, error } = useQuery({
      query: gql`
        query PaddedApp($isMultiTenancyEnabled: Boolean!) {
          user {
            id
            name
          }

          serverStatus @skip(if: $isMultiTenancyEnabled) {
            _id
            needsFirstTimeSetup
            version
          }
        }
      `,
      requestPolicy: 'cache-and-network',
      variables: {
        isMultiTenancyEnabled,
      },
    })

    return {
      displaySelfHostedBanner: !isMultiTenancyEnabled,
      data,
      error,
      features,
    }
  },
})
</script>
