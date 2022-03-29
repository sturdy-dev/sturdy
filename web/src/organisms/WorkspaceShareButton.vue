<template>
  <div class="z-40">
    <WorkspaceMergeGitHubButton
      v-if="shareViaGitHubPR"
      :workspace="workspace"
      :disabled="disabled"
      :disabledTooltipMessage="cantSubmitTooltipMessage"
      :hunk-ids="allHunkIds"
    />

    <WorkspaceMergeRemoteButton
      v-else-if="shareViaRemote"
      :workspace="workspace"
      :disabled="disabled"
      :disabledTooltipMessage="cantSubmitTooltipMessage"
    />

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
import { defineComponent, type PropType } from 'vue'
import { gql } from '@urql/vue'

import WorkspaceMergeButton from './WorkspaceMergeButton.vue'
import WorkspaceMergeGitHubButton, {
  WORKSPACE_FRAGMENT as MERGE_GITHUB_BUTTON_WORKSPACE_FRAGMENT,
} from './WorkspaceMergeGitHubButton.vue'
import WorkspaceMergeRemoteButton, {
  WORKSPACE_FRAGMENT as MERGE_REMOTE_BUTTON_WORKSPACE_FRAGMENT,
} from './WorkspaceMergeRemoteButton.vue'

import type { ShareButtonFragment } from './__generated__/WorkspaceShareButton'

export const SHARE_BUTTON = gql`
  fragment ShareButton on Workspace {
    id
    codebase {
      gitHubIntegration @include(if: $isGitHubEnabled) {
        id
        enabled
        gitHubIsSourceOfTruth
      }
      remote @include(if: $isRemoteEnabled) {
        id
      }
    }
    ...MergeGitHubButton_Workspace
    ...MergeRemoteButton_Workspace
  }
  ${MERGE_GITHUB_BUTTON_WORKSPACE_FRAGMENT}
  ${MERGE_REMOTE_BUTTON_WORKSPACE_FRAGMENT}
`

export enum CANT_SUBMIT_REASON {
  WORKSPACE_NOT_FOUND,
  NO_DIFFS,
  EMPTY_DESCRIPTION,
  HAVE_SELECTED_HUNKS,
}

export default defineComponent({
  components: {
    WorkspaceMergeButton,
    WorkspaceMergeGitHubButton,
    WorkspaceMergeRemoteButton,
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
  computed: {
    shareViaGitHubPR() {
      return (
        this.workspace.codebase.gitHubIntegration?.enabled &&
        this.workspace.codebase.gitHubIntegration?.gitHubIsSourceOfTruth
      )
    },
    shareViaRemote() {
      return !!this.workspace.codebase.remote
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
  },
})
</script>
