<template>
  <OnboardingStep id="SubmittingToRemoteGit" :dependencies="['MakingAChange', 'WorkspaceChanges']">
    <template #title>Submit to {{ workspace.codebase.remote.name }}</template>
    <template #description>
      When you're ready, use this button to push this workspace as a branch to
      {{ workspace.codebase.remote.name }}.
    </template>
    <div class="flex flex-col gap-2 items-end">
      <a
        v-if="pushedWorkspace && gitRemoteBranchURL"
        :href="gitRemoteBranchURL"
        target="_blank"
        class="flex items-center text-sm text-blue-800"
      >
        <span>Go to branch</span>
        <ExternalLinkIcon class="w-4 h-4 ml-1" />
      </a>

      <ButtonWithDropdown
        color="blue"
        :disabled="disabled || pushingWorkspace"
        :show-tooltip="disabled"
        :spinner="pushingWorkspace"
        :tooltip-right="true"
        @click="triggerPushWorkspace"
      >
        <template #default>
          <template v-if="pushingWorkspace">
            <template v-if="isMergingAndPushing"
              >Merging and pushing to {{ workspace.codebase.remote.name }}</template
            >
            <template v-else>Pushing to {{ workspace.codebase.remote.name }}</template>
          </template>
          <template v-else>Push to {{ workspace.codebase.remote.name }}</template>
        </template>

        <template v-if="disabled" #tooltip>
          {{ disabledTooltipMessage }}
        </template>

        <template #dropdown="{ disabled }">
          <MenuItem :disabled="disabled">
            <Button
              class="text-sm text-left py-2 px-4 flex hover:bg-gray-50"
              :icon="shareIcon"
              @click="triggerPushWorkspaceWithMerge"
            >
              Merge and push to {{ remote.name }}
            </Button>
          </MenuItem>
        </template>
      </ButtonWithDropdown>
    </div>
  </OnboardingStep>
</template>

<script lang="ts">
import { defineComponent, type PropType } from 'vue'
import { gql } from '@urql/vue'

import { ShareIcon } from '@heroicons/vue/solid'
import { MenuItem } from '@headlessui/vue'
import { ExternalLinkIcon } from '@heroicons/vue/outline'
import OnboardingStep from '../components/onboarding/OnboardingStep.vue'
import ButtonWithDropdown from '../atoms/ButtonWithDropdown.vue'
import Button from '../atoms/Button.vue'

import type { MergeRemoteButton_WorkspaceFragment } from './__generated__/WorkspaceMergeRemoteButton'

import { usePushWorkspace } from '../mutations/usePushWorkspace'

export const WORKSPACE_FRAGMENT = gql`
  fragment MergeRemoteButton_Workspace on Workspace {
    id
    codebase {
      id
      remote @include(if: $isRemoteEnabled) {
        id
        name
        browserLinkBranch
      }
    }
  }
`

export default defineComponent({
  components: { MenuItem, OnboardingStep, ButtonWithDropdown, Button, ExternalLinkIcon },
  props: {
    workspace: {
      type: Object as PropType<MergeRemoteButton_WorkspaceFragment>,
      required: true,
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    disabledTooltipMessage: {
      type: String,
      required: true,
    },
  },
  setup() {
    const { mutating: pushingWorkspace, pushWorkspace } = usePushWorkspace()
    return {
      pushingWorkspace,
      pushWorkspace,

      shareIcon: ShareIcon,
    }
  },
  data() {
    return {
      isMergingAndPushing: false,
      pushedWorkspace: false,
    }
  },
  computed: {
    gitRemoteBranchURL() {
      return this.workspace.codebase.remote
        ? this.workspace.codebase.remote.browserLinkBranch.replace(
            '${BRANCH_NAME}',
            'sturdy-' + this.workspace.id
          )
        : null
    },
  },
  methods: {
    async triggerPushWorkspace(landOnSturdyAndPushTracked = false) {
      this.pushedWorkspace = false
      await this.pushWorkspace({
        workspaceID: this.workspace.id,
        landOnSturdyAndPushTracked: landOnSturdyAndPushTracked,
      })
        .then(() => {
          this.pushedWorkspace = true
        })
        .catch((e) => {
          const title = 'Failed!'
          const message = 'Failed to push workspace'

          console.error(e)

          this.emitter.emit('notification', {
            title: title,
            message,
            style: 'error',
          })
        })
    },

    async triggerPushWorkspaceWithMerge() {
      this.isMergingAndPushing = true
      await this.triggerPushWorkspace(true).finally(() => {
        this.isMergingAndPushing = false
      })
    },
  },
})
</script>
