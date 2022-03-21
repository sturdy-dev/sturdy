<template>
  <div class="flex items-center gap-2">
    <CalendarIcon class="h-5 w-5 text-gray-400" aria-hidden="true" />
    <span class="text-gray-900 text-sm font-medium">
      <template v-if="updatedAt > 0"> Updated <RelativeTime :date="updatedAt" /> </template>
      <template v-else> Created <RelativeTime :date="createdAt" /> </template>
    </span>
  </div>
</template>

<script lang="ts">
import { CalendarIcon } from '@heroicons/vue/solid'
import RelativeTime from '../../../atoms/RelativeTime.vue'
import { gql } from '@urql/vue'

export const WORKSPACE_FRAGMENT = gql`
  fragment Updated_Workspace on Workspace {
    id
    updatedAt
    createdAt
  }
`

export default {
  components: { CalendarIcon, RelativeTime },
  props: {
    workspace: {
      type: Object,
      required: true,
    },
  },
  computed: {
    updatedAt() {
      return new Date(this.workspace.updatedAt * 1000)
    },
    createdAt() {
      return new Date(this.workspace.createdAt * 1000)
    },
  },
}
</script>
