<template>
  <template v-for="message in messages" :key="message.text">
    <Banner :status="bannerLevel(message.level)" :message="message.text" />
  </template>
</template>

<script lang="ts">
import { gql } from '@urql/vue'
import { PropType } from 'vue'
import { BannerLicenseMessageFragment } from './__generated__/Banner'
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
