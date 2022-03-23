<template>
  <component :is="badge.icon" v-if="badge" :class="badge.class"></component>
</template>

<script lang="ts">
import { gql } from '@urql/vue'
import { defineComponent } from 'vue'
import type { PropType } from 'vue'
import { StatusType } from '../../__generated__/types'
import { CheckCircleIcon, ClockIcon, XCircleIcon } from '@heroicons/vue/solid'
import type { StatusFragment } from './__generated__/StatusBadge'

export const STATUS_FRAGMENT = gql`
  fragment Status on Status {
    id
    title
    description
    type
    timestamp
    detailsUrl
  }
`

export default defineComponent({
  components: { XCircleIcon, CheckCircleIcon, ClockIcon },
  props: {
    statuses: {
      type: Array as PropType<StatusFragment[]>,
      default: () => [],
    },
  },
  computed: {
    badge() {
      if (this.statuses.length === 0) {
        return null
      }

      const common = 'inline h-5 w-5'

      // If any has failed, mark as failed
      if (this.statuses.some((s) => s.type == StatusType.Failing)) {
        return {
          icon: XCircleIcon,
          class: `${common} text-red-400`,
        }
      }

      // If any is pending, mark as pending
      if (this.statuses.some((s) => s.type == StatusType.Pending)) {
        return {
          icon: ClockIcon,
          class: `${common} animate-pulse text-blue-400`,
        }
      }

      // If all are healthy, mark as healthy
      if (!this.statuses.some((s) => s.type != StatusType.Healthy)) {
        return {
          icon: CheckCircleIcon,
          class: `${common} text-green-400`,
        }
      }

      return null
    },
  },
})
</script>
