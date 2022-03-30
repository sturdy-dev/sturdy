<template>
  <OnboardingStep id="SubmittingToRemoteGit" :dependencies="['MakingAChange', 'WorkspaceChanges']">
    <template #title>Submit to {{ workspace.codebase.remote?.name }}</template>
    <template #description>
      When you're ready, use this button to push this workspace as a branch to
      {{ workspace.codebase.remote?.name }}.
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

      <Select id="merge-remote-method" color="blue">
        <template #selected="{ item }">
          <component
            color="blue"
            :is="item"
            :show-tooltip="disabled"
            :disabled="disabled || pushingWorkspace || isMergingAndPushing"
            :spinner="pushingWorkspace || isMergingAndPushing"
            :tooltip-right="true"
            class="rounded-r-none"
          />
        </template>

        <template #options>
          <Button
            class="text-sm text-left py-2 px-4 flex border-0 hover:bg-gray-50"
            @click="() => triggerPushWorkspace()"
          >
            <template #default>
              {{
                pushingWorkspace
                  ? `Pushing to ${workspace.codebase.remote?.name}`
                  : `Push to ${workspace.codebase.remote?.name}`
              }}
            </template>
            <template v-if="disabled" #tooltip>
              {{ disabledTooltipMessage }}
            </template>
          </Button>

          <Button
            class="text-sm text-left py-2 px-4 flex border-0 hover:bg-gray-50"
            @click="triggerPushWorkspaceWithMerge"
          >
            <template #default>
              {{
                isMergingAndPushing
                  ? `Merging and pushing to ${workspace.codebase.remote?.name}`
                  : `Merge and push to ${workspace.codebase.remote?.name}`
              }}
            </template>

            <template v-if="disabled" #tooltip>
              {{ disabledTooltipMessage }}
            </template>
          </Button>
        </template>
      </Select>
    </div>
  </OnboardingStep>
</template>

<script lang="ts">
import { defineComponent, type PropType } from 'vue'
import { gql } from '@urql/vue'
import { ShareIcon } from '@heroicons/vue/solid'
import { ExternalLinkIcon } from '@heroicons/vue/outline'
import OnboardingStep from '../components/onboarding/OnboardingStep.vue'
import Button from '../atoms/Button.vue'
import Select from '../atoms/Select.vue'

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
  components: {
    OnboardingStep,
    Button,
    ExternalLinkIcon,
    Select,
  },
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
  watch: {
    workspace: function (a, b) {
      // reset local state on navigation
      if (a?.id !== b?.id) {
        this.isMergingAndPushing = false
        this.pushedWorkspace = false
      }
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
