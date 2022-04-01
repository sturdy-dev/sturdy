<template>
  <ul class="max-w-3xl mx-auto mt-16 flex flex-col gap-2">
    <li v-for="message in messages" :key="message.text">
      {{ message.text }}
    </li>
  </ul>
</template>

<script lang="ts">
import { gql } from '@urql/vue'
import type { PropType } from 'vue'
import type { FullscreenLicenseMessageFragment } from './__generated__/Fullscreen'

export const FULLSCREEN_MESSAGE_FRAGMENT = gql`
  fragment FullscreenLicenseMessage on LicenseMessage {
    text
  }
`

export default {
  props: {
    messages: {
      type: Array as PropType<FullscreenLicenseMessageFragment[]>,
      required: false,
      default: function () {
        return []
      },
    },
  },
}
</script>
