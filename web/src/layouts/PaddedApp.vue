<template>
  <div class="p-4 sm:p-8">
    <div v-if="displaySelfHostedBanner" class="bg-green-500 m-4 p-2">
      <strong>Debug: This is a self-hosted instance of Sturdy.</strong>
    </div>

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

    const { data, error } = useQuery({
      query: gql`
        query PaddedApp($isMultiTenancyEnabled: Boolean!) {
          user {
            id
            name
          }

          installation @skip(if: $isMultiTenancyEnabled) {
            id
            needsFirstTimeSetup
            version

            license {
              id
              key
              createdAt
              expiresAt
              status
              messages {
                level
                text
                type
              }
            }
          }
        }
      `,
      requestPolicy: 'cache-and-network',
      variables: {
        isMultiTenancyEnabled,
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
