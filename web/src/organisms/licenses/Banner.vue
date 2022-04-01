<template>
  <template v-for="message in messages" :key="message.text">
    <Banner :status="bannerLevel(message.level)" :message="message.text" class="mb-2" />
  </template>
</template>

<script lang="ts">
import { gql } from '@urql/vue'
import type { PropType } from 'vue'
import type { BannerLicenseMessageFragment } from './__generated__/Banner'
import { Banner } from '../../atoms'
import { LicenseMessageLevel } from '../../__generated__/types'

export const BANNER_MESSAGE_FRAGMENT = gql`
  fragment BannerLicenseMessage on LicenseMessage {
    text
    level
  }
`

export default {
  components: { Banner },
  props: {
    messages: {
      type: Array as PropType<BannerLicenseMessageFragment[]>,
      required: false,
      default: function () {
        return []
      },
    },
  },
  methods: {
    bannerLevel(level: LicenseMessageLevel) {
      switch (level) {
        case LicenseMessageLevel.Info:
          return 'info'
        case LicenseMessageLevel.Warning:
          return 'warning'
        case LicenseMessageLevel.Error:
          return 'error'
        default:
          throw new Error(`Unknown level: ${level}`)
      }
    },
  },
}
</script>
