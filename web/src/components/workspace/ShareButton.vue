<template>
  <div>
    <OnboardingStep
      v-if="shareViaGitHubPR"
      id="SubmittingAPullRequest"
      :dependencies="['MakingAChange', 'WorkspaceChanges']"
    >
      <template #title>Submitting to GitHub</template>
      <template #description>
        When you're ready, use this button to create a pull request on GitHub. If you don't want to
        create a PR, but rather have GitHub just receive updates from Sturdy, update your GitHub
        Integration settings.
      </template>
      <div class="flex flex-col gap-2 items-end">
        <a
          v-if="hasGitHubPR && hasOpenGitHubPR"
          :href="gitHubPRLink"
          target="_blank"
          class="flex items-center text-sm text-blue-800"
        >
          <span>Go to pull request</span>
          <ExternalLinkIcon class="w-4 h-4 ml-1" />
        </a>

        <div class="gap-2 flex">
          <Button
            color="blue"
            size="wider"
            :disabled="creatingOrUpdatingPR || disabled || isMerging"
            :class="[creatingOrUpdatingPR || disabled ? 'cursor-default' : '']"
            :show-tooltip="disabled"
            :tooltip-right="true"
            @click="createOrUpdatePR"
          >
            <template #default>
              <div v-if="creatingOrUpdatingPR" class="flex items-center">
                <Spinner class="mr-1" />
                <span v-if="!hasOpenGitHubPR">Creating pull request</span>
                <span v-else> Updating pull request </span>
              </div>
              <span v-else-if="!hasOpenGitHubPR">Create pull request</span>
              <span v-else> Update pull request </span>
            </template>

            <template #tooltip>
              {{ cantSubmitTooltipMessage }}
            </template>
          </Button>

          <Button
            v-if="hasGitHubPR && hasOpenGitHubPR"
            color="green"
            :disabled="isMerging"
            :show-tooltip="isMerging"
            :tooltip-right="true"
            @click="triggerMergePullRequest"
          >
            <template #tooltip>Hang on, we are waiting for GitHub to call us back...</template>
            <template #default>
              <div v-if="isMerging" class="flex items-center">
                <Spinner class="mr-1" />
                <span>Merging</span>
              </div>
              <span v-else>Merge</span>
            </template>
          </Button>
        </div>
      </div>
    </OnboardingStep>

    <OnboardingStep
      v-else
      id="LandingAChange"
      :dependencies="['MakingAChange', 'WorkspaceChanges']"
    >
      <template #title>Publishing a Change</template>
      <template #description>
        When you're ready, use this button to save the changes you've made so far.
      </template>
      <Button
        color="blue"
        :disabled="landing || disabled"
        :class="[landing || disabled ? 'cursor-default' : '']"
        :show-tooltip="disabled"
        :tooltip-right="true"
        @click="shareChange"
      >
        <template #default>
          <div v-if="landing" class="flex items-center">
            <Spinner class="mr-1" />
            <span>Sharing</span>
          </div>
          <span v-else>Share</span>
        </template>

        <template #tooltip>{{ cantSubmitTooltipMessage }} </template>
      </Button>
    </OnboardingStep>
  </div>
</template>

<script lang="ts">
import { defineComponent, nextTick, PropType } from 'vue'
import { gql } from '@urql/vue'
import { ShareButtonFragment } from './__generated__/ShareButton'
import OnboardingStep from '../onboarding/OnboardingStep.vue'
import Button from '../shared/Button.vue'
import Spinner from '../shared/Spinner.vue'
import { ExternalLinkIcon } from '@heroicons/vue/outline'
import { useLandWorkspaceChange } from '../../mutations/useLandWorkspaceChange'
import { useCreateOrUpdateGitHubPullRequest } from '../../mutations/useCreateOrUpdateGitHubPullRequest'
import { useMergeGitHubPullRequest } from '../../mutations/useMergeGitHubPullRequest'
import { GitHubPullRequestState } from '../../__generated__/types'

export const SHARE_BUTTON = gql`
  fragment ShareButton on Workspace {
    id
    codebase {
      id
      gitHubIntegration {
        id
        enabled
        gitHubIsSourceOfTruth
        owner
        name
      }
    }
    gitHubPullRequest {
      id
      pullRequestNumber
      state
    }
  }
`

export enum CANT_SUBMIT_REASON {
  WORKSPACE_NOT_FOUND,
  NO_DIFFS,
  EMPTY_DESCRIPTION,
  HAVE_SELECTED_HUNKS,
}

export default defineComponent({
  components: { Spinner, Button, OnboardingStep, ExternalLinkIcon },
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
    const { mutating: landing, landWorkspaceChange } = useLandWorkspaceChange()
    const { mutating: creatingOrUpdatingPR, createOrUpdateGitHubPullRequest } =
      useCreateOrUpdateGitHubPullRequest()
    const { mutating: mergingGitHubPullRequest, mergeGitHubPullRequest } =
      useMergeGitHubPullRequest()

    return {
      landing,
      landWorkspaceChange,

      creatingOrUpdatingPR,
      createOrUpdateGitHubPullRequest,

      mergingGitHubPullRequest,
      mergeGitHubPullRequest,
    }
  },
  computed: {
    isMerging() {
      return (
        this.mergingGitHubPullRequest ||
        this.workspace.gitHubPullRequest?.state === GitHubPullRequestState.Merging
      )
    },
    shareViaGitHubPR() {
      if (this.workspace.codebase.gitHubIntegration?.enabled) {
        if (this.workspace.codebase.gitHubIntegration?.gitHubIsSourceOfTruth) {
          return true
        }
      }
      return false
    },
    hasGitHubPR() {
      return Boolean(this.workspace.gitHubPullRequest?.pullRequestNumber)
    },
    hasOpenGitHubPR() {
      const isOpen = this.workspace.gitHubPullRequest?.state == GitHubPullRequestState.Open
      const isMerging = this.workspace.gitHubPullRequest?.state == GitHubPullRequestState.Merging
      return isOpen || isMerging
    },
    gitHubPRLink() {
      const { owner, name } = this.workspace.codebase.gitHubIntegration ?? {}
      const { pullRequestNumber } = this.workspace.gitHubPullRequest ?? {}
      return `https://github.com/${owner}/${name}/pull/${pullRequestNumber}`
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
  methods: {
    async shareChange() {
      // Triggers WorkspaceHome to flush the draft description
      await this.$emit('pre-create-change')

      const input = {
        workspaceID: this.workspace.id,
        patchIDs: this.allHunkIds,
      }

      nextTick(() =>
        this.landWorkspaceChange(input).catch((e) => {
          console.error(e)
          this.emitter.emit('notification', {
            title: 'Failed sharing changes',
            message: 'Sorry about that! You might need to sync first!',
            style: 'error',
          })
        })
      )
    },
    async createOrUpdatePR() {
      // Triggers WorkspaceHome to flush the draft description
      await this.$emit('pre-create-change')

      const input = {
        workspaceID: this.workspace.id,
        patchIDs: this.allHunkIds,
      }

      nextTick(() =>
        this.createOrUpdateGitHubPullRequest(input).catch((e) => {
          let message = 'Failed to create or update pull request'

          // Server generated error if the push fails (due to branch protection rules, etc)
          if (e.graphQLErrors && e.graphQLErrors.length > 0) {
            if (e.graphQLErrors[0].extensions?.pushFailure) {
              message = e.graphQLErrors[0].extensions.pushFailure
            } else if (e.graphQLErrors[0].extensions?.createPullRequestFailure) {
              message = e.graphQLErrors[0].extensions.createPullRequestFailure
            } else if (e.graphQLErrors[0].extensions?.getPullRequestFailure) {
              message = e.graphQLErrors[0].extensions.getPullRequestFailure
            } else if (e.graphQLErrors[0].extensions?.updatePullRequestFailure) {
              message = e.graphQLErrors[0].extensions.updatePullRequestFailure
            } else if (e.graphQLErrors[0].extensions?.message) {
              message = e.graphQLErrors[0].extensions.message
            } else {
              console.error(e)
            }
          } else {
            console.error(e)
          }

          this.emitter.emit('notification', {
            title: 'Failed!',
            message,
            style: 'error',
          })
        })
      )
    },
    async triggerMergePullRequest() {
      const input = {
        workspaceID: this.workspace.id,
      }

      await this.mergeGitHubPullRequest(input).catch((e) => {
        let title = 'Failed!'
        let message = 'Failed to merge pull request'

        // Server generated error if the push fails (due to branch protection rules, etc)
        if (e.graphQLErrors && e.graphQLErrors.length > 0) {
          if (e.graphQLErrors[0].extensions?.message) {
            title = 'GitHub error'
            message = e.graphQLErrors[0].extensions.message
          } else {
            console.error(e)
          }
        } else {
          console.error(e)
        }

        this.emitter.emit('notification', {
          title: title,
          message,
          style: 'error',
        })
      })
    },
  },
})
</script>
