<template>
  <Tooltip v-if="workspace.statuses.length > 0" :disabled="!isStale">
    <template #tooltip> Draft changed since the last run</template>
    <template #default>
      <StatusDetails :statuses="workspace.statuses" :stale="isStale" />
    </template>
  </Tooltip>

  <Button
    v-if="
      (workspace.statuses.length == 0 && ciEnabled) || isStale || triggeredAndWaitingForStatuses
    "
    :icon="terminalIcon"
    :spinner="triggering || triggeredAndWaitingForStatuses"
    class="border-0 -ml-3"
    @click="onTriggerClicked"
    >Trigger CI
  </Button>
</template>

<script lang="ts">
import { computed, defineComponent, type PropType, toRefs } from 'vue'
import { gql } from '@urql/vue'

import StatusDetails, { STATUS_FRAGMENT } from '../components/statuses/StatusDetails.vue'
import Button from '../atoms/Button.vue'
import Tooltip from '../atoms/Tooltip.vue'
import { TerminalIcon } from '@heroicons/vue/outline'

import type { WorkspaceStatus_WorkspaceFragment } from './__generated__/WorkspaceStatus'
import { IntegrationProvider } from '../__generated__/types'
import { useUpdatedWorkspacesStatuses } from '../subscriptions/useUpdatedWorkspacesStatuses'
import { useTriggerInstantIntegration } from '../mutations/useTriggerInstantIntegration'

export const WORKSPACE_FRAGMENT = gql`
  fragment WorkspaceStatus_Workspace on Workspace {
    id
    statuses {
      id
      stale
      ...Status
    }
    codebase {
      id
      integrations {
        id
        provider
      }
      gitHubIntegration {
        id
        enabled
      }
    }
  }
  ${STATUS_FRAGMENT}
`

const ciProviders = [IntegrationProvider.Buildkite]

export default defineComponent({
  components: { StatusDetails, Button, Tooltip },
  props: {
    workspace: { type: Object as PropType<WorkspaceStatus_WorkspaceFragment>, required: true },
  },
  setup(props) {
    const { workspace } = toRefs(props)

    const ciEnabled = computed(() => {
      const hasCiProvider = workspace.value.codebase.integrations.some(({ provider }) =>
        ciProviders.includes(provider)
      )
      if (hasCiProvider) return true

      // github it not a "integration", requires some special treatment
      if (workspace.value.codebase?.gitHubIntegration?.enabled) {
        return true
      }

      return false
    })

    const workspaceIds = computed(() => [workspace.value.id])
    useUpdatedWorkspacesStatuses(workspaceIds)

    const triggerInstantIntegration = useTriggerInstantIntegration()
    return { triggerInstantIntegration, terminalIcon: TerminalIcon, ciEnabled }
  },
  data() {
    return {
      triggering: false,
      triggered: false,
    }
  },
  computed: {
    isStale() {
      return this.workspace.statuses.some(({ stale }) => stale)
    },
    triggeredAndWaitingForStatuses() {
      return this.triggered && !this.workspace.statuses.length
    },
  },
  methods: {
    onTriggerClicked() {
      this.triggering = true
      this.triggerInstantIntegration({
        workspaceID: this.workspace.id,
      })
        .then(() => {
          this.triggered = true
        })
        .catch(() => {
          this.emitter.emit('notification', {
            title: 'Failed to Trigger CI',
            message: 'Please try again later.',
            style: 'error',
          })
          this.triggering = false
          this.triggered = false
        })
        .finally(() => {
          this.triggering = false
        })
    },
  },
})
</script>
