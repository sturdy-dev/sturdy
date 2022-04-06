<template>
  <StatusDetails v-if="workspace.statuses.length > 0" :statuses="workspace.statuses" />
  <Button
    v-else-if="ciEnabled"
    :icon="terminalIcon"
    :spinner="triggering"
    class="border-0"
    @click="onTriggerClicked"
    >Trigger CI</Button
  >
</template>

<script lang="ts">
import { defineComponent, type PropType, toRefs } from 'vue'
import { gql } from '@urql/vue'

import StatusDetails, { STATUS_FRAGMENT } from '../components/statuses/StatusDetails.vue'
import Button from '../atoms/Button.vue'
import { TerminalIcon } from '@heroicons/vue/outline'

import type { WorkspaceStatus_WorkspaceFragment } from './__generated__/WorkspaceStatus'
import { IntegrationProvider } from '../__generated__/types'
import { useUpdatedWorkspacesStatuses } from '../subscriptions/useUpdatedWorkspacesStatuses'
import { useTriggerInstantIntegration } from '../mutations/useTriggerInstantIntegration'

export const WORKSPACE_FRAGMENT = gql`
  fragment WorkspaceStatus_Workspace on Workspace {
    id
    statuses {
      ...Status
    }
    codebase {
      id
      integrations {
        id
        provider
      }
    }
  }
  ${STATUS_FRAGMENT}
`

const ciProviders = [IntegrationProvider.Buildkite]

export default defineComponent({
  components: { StatusDetails, Button },
  props: {
    workspace: { type: Object as PropType<WorkspaceStatus_WorkspaceFragment>, required: true },
  },
  setup(props) {
    const { workspace } = toRefs(props)
    const ciEnabled = workspace.value.codebase.integrations.some(({ provider }) =>
      ciProviders.includes(provider)
    )
    if (ciEnabled) {
      useUpdatedWorkspacesStatuses([workspace.value.id])
    }
    const triggerInstantIntegration = useTriggerInstantIntegration()
    return { triggerInstantIntegration, terminalIcon: TerminalIcon, ciEnabled }
  },
  data() {
    return {
      triggering: false,
    }
  },
  methods: {
    onTriggerClicked() {
      this.triggering = true
      this.triggerInstantIntegration({
        workspaceID: this.workspace.id,
      }).finally(() => (this.triggering = false))
    },
  },
})
</script>
