<template>
  <div class="p-4 sm:p-8">
    <Banner :messages="bannerMessages" />
    <Fullscreen v-if="fullscreenMessages.length > 0" :messages="fullscreenMessages" />

    <FirstTimeUserNoNameTakeover v-if="data && data.user && !data.user.name" />
    <FirstTimeSetupSelfHosted
      v-else-if="data && data?.serverStatus?.needsFirstTimeSetup && data.user"
      :user="data.user"
    />
    <FirstTimeCreateOrganizationTakeover
      v-else-if="data && data.organizations.length === 0"
      :user="data.user"
    />
    <slot v-else></slot>
  </div>
</template>

<script lang="ts">
import { defineComponent, inject, ref, Ref } from 'vue'
import { gql, useQuery } from '@urql/vue'
import FirstTimeUserNoNameTakeover from '../components/user/FirstTimeUserNoNameTakeover.vue'
import { Feature, LicenseMessageType } from '../__generated__/types'
import Banner, { BANNER_MESSAGE_FRAGMENT } from '../organisms/licenses/Banner.vue'
import Fullscreen, { FULLSCREEN_MESSAGE_FRAGMENT } from '../organisms/licenses/Fullscreen.vue'
import { PaddedAppQuery, PaddedAppQueryVariables } from './__generated__/PaddedApp'
import FirstTimeSetupSelfHosted from '../organisms/serverstatus/FirstTimeSetupSelfHosted.vue'
import FirstTimeCreateOrganizationTakeover from '../components/user/FirstTimeCreateOrganizationTakeover.vue'

export default defineComponent({
  components: {
    FirstTimeUserNoNameTakeover,
    FirstTimeSetupSelfHosted,
    FirstTimeCreateOrganizationTakeover,
    Banner,
    Fullscreen,
  },

  setup() {
    const features = inject<Ref<Array<Feature>>>('features', ref([]))
    const isMultiTenancyEnabled = features?.value?.includes(Feature.MultiTenancy)
    const isSelfHostedLicenseEnabled = features?.value?.includes(Feature.SelfHostedLicense)

    const { data, error } = useQuery<PaddedAppQuery, PaddedAppQueryVariables>({
      query: gql`
        query PaddedApp($isMultiTenancyEnabled: Boolean!, $isSelfHostedLicenseEnabled: Boolean!) {
          user {
            id
            name
          }

          installation @skip(if: $isMultiTenancyEnabled) {
            id
            needsFirstTimeSetup
            version

            license @include(if: $isSelfHostedLicenseEnabled) {
              id
              messages {
                type
                ...BannerLicenseMessage
                ...FullscreenLicenseMessage
              }
            }
          }

          organizations {
            id
          }
        }
        ${BANNER_MESSAGE_FRAGMENT}
        ${FULLSCREEN_MESSAGE_FRAGMENT}
      `,
      requestPolicy: 'cache-and-network',
      variables: {
        isMultiTenancyEnabled,
        isSelfHostedLicenseEnabled,
      },
    })

    return {
      displaySelfHostedBanner: !isMultiTenancyEnabled,
      data,
      error,
      features,
    }
  },

  computed: {
    fullscreenMessages() {
      return (
        this.data?.installation?.license?.messages?.filter(
          (message) => message.type === LicenseMessageType.Fullscreen
        ) || []
      )
    },
    bannerMessages() {
      return (
        this.data?.installation?.license?.messages?.filter(
          (message) => message.type === LicenseMessageType.Banner
        ) || []
      )
    },
  },
})
</script>
