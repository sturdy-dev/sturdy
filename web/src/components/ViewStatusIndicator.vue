<template>
  <div class="flex items-center">
    <div
      class="p-0.5 relative flex-none flex items-center justify-center"
      :class="[compact ? 'h-5 w-5' : 'h-7 w-7']"
    >
      <ChevronDoubleUpIcon
        v-if="isConnected && view.status?.state === 'Transferring'"
        class="animate-pulse"
      />
      <Spinner v-else-if="isConnected && view.status?.state === 'Finishing'" />
      <Spinner v-else-if="isConnected && view.status?.state === 'Reconciling'" />
      <Spinner v-else-if="isConnected && view.status?.state === 'Connecting'" />
      <Spinner v-else-if="isConnected && view.status?.state === 'Scanning'" />
      <template v-else>
        <DesktopComputerIcon />

        <div
          v-if="view.status?.state"
          class="absolute top-0 right-0 rounded-full border"
          :class="[compact ? 'w-2 h-2' : 'w-3 h-3', indicatorBorderColor, indicatorColor]"
        />
      </template>
    </div>

    <div v-if="!compact && view.status" class="ml-3 flex flex-col">
      <div
        v-if="view.status.progressPath && isConnected"
        class="inline-block font-medium w-64 whitespace-nowrap overflow-hidden text-ellipsis"
      >
        {{ view.status.progressPath?.split('/').pop() }}
      </div>
      <div v-else class="inline-block font-medium">
        {{ view.shortMountPath }}
      </div>
      <div class="text-sm text-gray-500">
        <span v-if="view.status.state === 'Ready'">
          {{ view.mountHostname }}
        </span>
        <span
          v-else-if="
            view.status.state === 'Transferring' &&
            view.status.progressReceived != null &&
            isConnected
          "
        >
          <span
            >Syncing file {{ view.status.progressReceived + 1 }} of
            {{ view.status.progressTotal }}</span
          >
          <DotDotDot class="ml-1" />
        </span>
        <span v-else-if="!isConnected">Disconnected</span>
        <span v-else>
          {{ view.status.state }}
        </span>
      </div>
    </div>

    <div
      v-if="!compact && isSuggesting"
      title="You're connected to someone else's draft change. All edits you make will appear to them as suggestions."
      class="rounded-full bg-green-200 text-green-800 text-sm px-2 py-0.5 ml-3"
    >
      Suggesting
    </div>
  </div>
</template>

<script lang="ts">
import gql from 'graphql-tag'
import { defineComponent, onUnmounted, ref } from 'vue'
import type { PropType } from 'vue'
import type { ViewStatusIndicatorFragment } from './__generated__/ViewStatusIndicator'
import { ViewStatusState } from '../__generated__/types'
import { ChevronDoubleUpIcon, DesktopComputerIcon } from '@heroicons/vue/outline'
import Spinner from './shared/Spinner.vue'
import DotDotDot from '../molecules/DotDotDot.vue'

export const VIEW_STATUS_INDICATOR = gql`
  fragment ViewStatusIndicator on View {
    id
    shortMountPath
    mountHostname
    lastUsedAt
    status {
      id
      state
      progressPath
      progressReceived
      progressTotal
    }
    workspace {
      suggestion {
        id
      }
    }
  }
`

export default defineComponent({
  name: 'ViewStatusIndicator',
  components: {
    DotDotDot,
    Spinner,
    DesktopComputerIcon,
    ChevronDoubleUpIcon,
  },
  props: {
    view: {
      type: Object as PropType<ViewStatusIndicatorFragment>,
      required: true,
    },
    compact: {
      type: Boolean,
      default: false,
    },
  },
  setup() {
    let ts = ref(+new Date() / 1000)

    let stop = setInterval(() => {
      ts.value = +new Date() / 1000
    }, 2 * 60 * 1000) // Every 2 minutes, this timing is not very critical

    onUnmounted(() => {
      clearInterval(stop)
    })

    return {
      ts,
    }
  },
  computed: {
    // Switch is checked to be exhaustive, so this rule is superfluous
    // eslint-disable-next-line vue/return-in-computed-property
    indicatorColor(): string {
      if (this.view?.status == null) {
        return ''
      }

      if (this.isConnected) {
        switch (this.view.status.state) {
          case ViewStatusState.Connecting:
            return 'bg-yellow-300'
          case ViewStatusState.Downloading:
          case ViewStatusState.Finishing:
          case ViewStatusState.Reconciling:
          case ViewStatusState.Scanning:
          case ViewStatusState.Transferring:
          case ViewStatusState.Uploading:
            return 'bg-blue-300'
          case ViewStatusState.Ready:
            return 'bg-green-300'
        }
      }

      // Disconnected
      return 'bg-red-300'
    },
    // Switch is checked to be exhaustive, so this rule is superfluous
    // eslint-disable-next-line vue/return-in-computed-property
    indicatorBorderColor(): string {
      if (this.view?.status == null) {
        return ''
      }

      if (this.isConnected) {
        switch (this.view.status.state) {
          case ViewStatusState.Connecting:
            return 'border-yellow-400'
          case ViewStatusState.Downloading:
          case ViewStatusState.Finishing:
          case ViewStatusState.Reconciling:
          case ViewStatusState.Scanning:
          case ViewStatusState.Transferring:
          case ViewStatusState.Uploading:
            return 'border-blue-300'
          case ViewStatusState.Ready:
            return 'border-green-400'
        }
      }
      // Disconnected
      return 'border-red-400'
    },
    isSuggesting(): boolean {
      return !!this.view.workspace?.suggestion
    },
    isConnected(): boolean {
      if (!this.view?.lastUsedAt) {
        return false
      }
      // Always mark as disconnected if the view has not connected to the server for 2 minutes (no matter the status)
      if (this.view.lastUsedAt < this.ts - 60 * 2) {
        return false
      }
      if (this.view.status == null) {
        return false
      }
      if (this.view?.status?.state === ViewStatusState.Disconnected) {
        return false
      }
      return true
    },
  },
})
</script>
