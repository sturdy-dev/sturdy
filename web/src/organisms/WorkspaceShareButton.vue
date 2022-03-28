<template>
  <div class="z-40">
    <WorkspaceMergeGitHubButton
      v-if="shareViaGitHubPR"
      :workspace="workspace"
      :hunk-ids="allHunkIds"
      :disabled="disabled"
      :disabledTooltipMessage="cantSubmitTooltipMessage"
    />

    <OnboardingStep
      v-else-if="shareViaRemote"
      id="SubmittingToRemoteGit"
      :dependencies="['MakingAChange', 'WorkspaceChanges']"
    >
      <template #title>Submit to {{ remote.name }}</template>
      <template #description>
        When you're ready, use this button to push this workspace as a branch to {{ remote.name }}.
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

        <div class="gap-2 flex">
          <ButtonWithDropdown
            color="blue"
            :disabled="disabled || pushingWorkspace"
            :show-tooltip="disabled"
            :tooltip-right="true"
            @click="triggerPushWorkspace"
          >
            <template #default>
              <div v-if="pushingWorkspace && isMergingAndPushing" class="flex items-center">
                <Spinner class="mr-1" />
                <span>Merging and pushing to {{ remote.name }}</span>
              </div>
              <div v-else-if="pushingWorkspace" class="flex items-center">
                <Spinner class="mr-1" />
                <span>Pushing to {{ remote.name }}</span>
              </div>
              <span v-else>Push to {{ remote.name }}</span>
            </template>

            <template #tooltip>
              {{ cantSubmitTooltipMessage }}
            </template>

            <template #dropdown="{ disabled }">
              <MenuItem :disabled="disabled">
                <button
                  class="text-sm text-left py-2 px-4 flex hover:bg-gray-50"
                  @click="triggerPushWorkspaceWithMerge"
                >
                  <ShareIcon class="-ml-1 mr-2 h-5 w-5 text-gray-400" aria-hidden="true" />
                  <span>Merge and push to {{ remote.name }}</span>
                </button>
              </MenuItem>
            </template>
          </ButtonWithDropdown>
        </div>
      </div>
    </OnboardingStep>

    <WorkspaceMergeButton
      v-else
      :workspace-id="workspace.id"
      :disabled="disabled"
      :disabled-tooltip-message="cantSubmitTooltipMessage"
      :hunk-ids="allHunkIds"
    />
  </div>
</template>

<script lang="ts">
import type { PropType } from 'vue'
import { defineComponent, ref } from 'vue'
import { gql } from '@urql/vue'
import type { ShareButtonFragment } from './__generated__/WorkspaceShareButton'
import OnboardingStep from '../components/onboarding/OnboardingStep.vue'
import Spinner from '../components/shared/Spinner.vue'
import { ExternalLinkIcon } from '@heroicons/vue/outline'
import { usePushWorkspace } from '../mutations/usePushWorkspace'
import ButtonWithDropdown from '../components/shared/ButtonWithDropdown.vue'
import { ShareIcon } from '@heroicons/vue/solid'
import { MenuItem } from '@headlessui/vue'
import WorkspaceMergeButton from './WorkspaceMergeButton.vue'
import WorkspaceMergeGitHubButton, {
  WORKSPACE_FRAGMENT as MERGE_GITHUB_BUTTON_WORKSPACE_FRAGMENT,
} from './WorkspaceMergeGitHubButton.vue'

export const SHARE_BUTTON = gql`
  fragment ShareButton on Workspace {
    id
    codebase {
      id
      remote @include(if: $isRemoteEnabled) {
        id
        name
        browserLinkBranch
      }
    }
    ...MergeGitHubButton_Workspace
  }
  ${MERGE_GITHUB_BUTTON_WORKSPACE_FRAGMENT}
`

export enum CANT_SUBMIT_REASON {
  WORKSPACE_NOT_FOUND,
  NO_DIFFS,
  EMPTY_DESCRIPTION,
  HAVE_SELECTED_HUNKS,
}

export default defineComponent({
  components: {
    ButtonWithDropdown,
    Spinner,
    OnboardingStep,
    ExternalLinkIcon,
    ShareIcon,
    MenuItem,
    WorkspaceMergeButton,
    WorkspaceMergeGitHubButton,
  },
  props: {
    workspace: {
      type: Object as PropType<ShareButtonFragment>,
      required: true,
    },
    allHunkIds: {
      type: Array as PropType<string[]>,
      default: () => [],
    },
    cantSubmitReason: {
      type: Number as PropType<CANT_SUBMIT_REASON>,
      default: null,
      required: false,
    },
    disabled: {
      type: Boolean,
      required: true,
    },
  },
  emits: {
    'pre-create-change': () => true,
  },
  setup() {
    const { mutating: pushingWorkspace, pushWorkspace } = usePushWorkspace()
    let pushedWorkspace = ref(false)
    let isMergingAndPushing = ref(false)

    return {
      pushingWorkspace,
      pushedWorkspace,
      isMergingAndPushing,
      pushWorkspace,
    }
  },
  computed: {
    shareViaGitHubPR() {
      if (this.workspace.codebase.gitHubIntegration?.enabled) {
        if (this.workspace.codebase.gitHubIntegration?.gitHubIsSourceOfTruth) {
          return true
        }
      }
      return false
    },
    shareViaRemote() {
      if (this.workspace.codebase.remote?.id) {
        return true
      }
      return false
    },

    cantSubmitTooltipMessage(): string {
      switch (this.cantSubmitReason) {
        case CANT_SUBMIT_REASON.WORKSPACE_NOT_FOUND:
          return 'Error, no workspace found'
        case CANT_SUBMIT_REASON.NO_DIFFS:
          return 'This workspace has no changes.'
        case CANT_SUBMIT_REASON.EMPTY_DESCRIPTION:
          return 'The change must be described before it can be shared.'
        case CANT_SUBMIT_REASON.HAVE_SELECTED_HUNKS:
          return "It's not possible to share a partial change. Deselect all changes before continuing."
        default:
          return 'This change can not be shared.'
      }
    },
    remote() {
      return this.workspace?.codebase?.remote
    },
    gitRemoteBranchURL(): string | null {
      const rem = this.remote
      if (!rem) {
        return null
      }
      return rem.browserLinkBranch.replace('${BRANCH_NAME}', 'sturdy-' + this.workspace.id)
    },
  },
  methods: {
    async triggerPushWorkspace(landOnSturdyAndPushTracked = false) {
      const input = {
        workspaceID: this.workspace.id,
        landOnSturdyAndPushTracked: landOnSturdyAndPushTracked,
      }

      if (landOnSturdyAndPushTracked) {
        this.isMergingAndPushing = true
      }

      this.pushedWorkspace = false

      await this.pushWorkspace(input)
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
      await this.triggerPushWorkspace(true)
    },
  },
})
</script>
