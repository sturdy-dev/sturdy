<template>
  <div class="flex rounded-md mb-1 mr-1">
    <RouterLinkButton
      :to="{
        name: 'installClient',
        params: { codebaseSlug: codebaseSlug },
      }"
    >
      <DownloadIcon class="w-4 h-4 mr-2" />
      Install CLI
    </RouterLinkButton>
  </div>
</template>
<script lang="ts">
import { DownloadIcon } from '@heroicons/vue/solid'
import RouterLinkButton from '../../atoms/RouterLinkButton.vue'
import { Slug } from '../../slug'
import { defineComponent, type PropType } from 'vue'
import { gql } from 'graphql-tag'
import type { InstallCliStepCodebaseFragment } from './__generated__/SetupSturdyInstallCliStep'

export const INSTALL_CLI_STEP_CODEBASE_FRAGMENT = gql`
  fragment InstallCliStepCodebase on Codebase {
    name
    shortID
  }
`

export default defineComponent({
  components: {
    DownloadIcon,
    RouterLinkButton,
  },
  props: {
    codebase: {
      required: true,
      type: Object as PropType<InstallCliStepCodebaseFragment>,
    },
  },
  computed: {
    codebaseSlug() {
      return Slug(this.codebase.name, this.codebase.shortID)
    },
  },
})
</script>
