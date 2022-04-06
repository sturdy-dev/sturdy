<template>
  <OnboardingStep id="SubmittingAPullRequest" :dependencies="['MakingAChange', 'WorkspaceChanges']">
    <template #title>Submitting to GitHub</template>
    <template #description>
      When you're ready, use this button to create a pull request on GitHub. If you don't want to
      create a PR, but rather have GitHub just receive updates from Sturdy, update your GitHub
      Integration settings.
    </template>
    <div class="flex flex-col gap-2 items-end">
      <a
        v-if="hasOpenGitHubPR"
        :href="gitHubPRLink"
        target="_blank"
        class="flex items-center text-sm text-blue-800"
      >
        <span>Go to pull request</span>
        <ExternalLinkIcon class="w-4 h-4 ml-1" />
      </a>

      <Select v-if="!hasOpenGitHubPR || creatingAndMergingPR" id="merge-github-method" color="blue">
        <template #selected="{ option }">
          <component
            :is="option"
            size="wider"
            color="blue"
            :disabled="disabled || isMerging || creatingOrUpdatingPR || creatingAndMergingPR"
            :show-tooltip="disabled"
            :tooltip-right="true"
            :spinner="isMerging || creatingOrUpdatingPR || creatingAndMergingPR"
            class="rounded-r-none"
          />
        </template>
        <template #options>
          <Button
            class="text-sm text-left py-2 px-4 flex border-0 hover:bg-gray-50"
            @click="createOrUpdatePR"
          >
            <template #default>
              {{ creatingOrUpdatingPR ? 'Creating pull request' : 'Create pull request' }}
            </template>

            <template v-if="disabled" #tooltip>
              {{ disabledTooltipMessage }}
            </template>
          </Button>

          <Button
            class="text-sm text-left py-2 px-4 flex border-0 hover:bg-gray-50"
            @click="createAndMergePR"
          >
            <template #default>
              {{
                creatingAndMergingPR && creatingOrUpdatingPR
                  ? 'Creating pull request'
                  : creatingAndMergingPR && isMerging
                  ? 'Merging pull request'
                  : 'Create and merge pull request'
              }}
            </template>

            <template v-if="disabled" #tooltip>
              {{ disabledTooltipMessage }}
            </template>
          </Button>
        </template>
      </Select>

      <div v-else class="gap-2 flex">
        <Button
          color="blue"
          size="wider"
          :disabled="creatingOrUpdatingPR || disabled || isMerging"
          :show-tooltip="disabled"
          :tooltip-right="true"
          :spinner="creatingOrUpdatingPR"
          @click="createOrUpdatePR"
        >
          <template #default>
            {{ creatingOrUpdatingPR ? 'Updating pull request' : 'Update pull request' }}
          </template>

          <template v-if="disabled" #tooltip>
            {{ disabledTooltipMessage }}
          </template>
        </Button>

        <Button
          color="green"
          :disabled="isMerging || creatingOrUpdatingPR"
          :show-tooltip="isMerging"
          :submitting-a-pull-request="isMerging"
          :tooltip-right="true"
          :spinner="isMerging"
          @click="triggerMergePullRequest"
        >
          <template #tooltip>Hang on, we are waiting for GitHub to call back to us...</template>
          <template #default>
            <template v-if="isMerging"> Merging</template>
            <template v-else> Merge</template>
          </template>
        </Button>
      </div>
    </div>
  </OnboardingStep>
</template>

<script lang="ts">
import { defineComponent, inject, type PropType } from 'vue'
import { gql } from '@urql/vue'

import OnboardingStep from '../components/onboarding/OnboardingStep.vue'
import Button from '../atoms/Button.vue'
import Select from '../atoms/Select.vue'
import { ExternalLinkIcon } from '@heroicons/vue/outline'

import { useCreateOrUpdateGitHubPullRequest } from '../mutations/useCreateOrUpdateGitHubPullRequest'
import { useMergeGitHubPullRequest } from '../mutations/useMergeGitHubPullRequest'

import type { MergeGitHubButton_WorkspaceFragment } from './__generated__/WorkspaceMergeGitHubButton'
import { GitHubPullRequestState } from '../__generated__/types'

export const WORKSPACE_FRAGMENT = gql`
  fragment MergeGitHubButton_Workspace on Workspace {
    id
    codebase {
      id
      gitHubIntegration @include(if: $isGitHubEnabled) {
        id
        owner
        name
      }
    }
    gitHubPullRequest @include(if: $isGitHubEnabled) {
      id
      pullRequestNumber
      state
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
      type: Object as PropType<MergeGitHubButton_WorkspaceFragment>,
      required: true,
    },
    disabled: {
      type: Boolean,
      required: false,
    },
    disabledTooltipMessage: {
      type: String,
      required: true,
    },
  },
  setup() {
    const { mutating: creatingOrUpdatingPR, createOrUpdateGitHubPullRequest } =
      useCreateOrUpdateGitHubPullRequest()
    const { mutating: mergingGitHubPullRequest, mergeGitHubPullRequest } =
      useMergeGitHubPullRequest()

    return {
      creatingOrUpdatingPR,
      createOrUpdateGitHubPullRequest,

      mergingGitHubPullRequest,
      mergeGitHubPullRequest,
    }
  },
  data() {
    return {
      creatingAndMergingPR: false,
    }
  },
  computed: {
    gitHubPRLink() {
      const { owner, name } = this.workspace.codebase.gitHubIntegration ?? {}
      const { pullRequestNumber } = this.workspace.gitHubPullRequest ?? {}
      return `https://github.com/${owner}/${name}/pull/${pullRequestNumber}`
    },
    isMerging() {
      return (
        this.mergingGitHubPullRequest ||
        this.workspace.gitHubPullRequest?.state === GitHubPullRequestState.Merging
      )
    },
    hasOpenGitHubPR() {
      const isOpen = this.workspace.gitHubPullRequest?.state == GitHubPullRequestState.Open
      return isOpen || this.isMerging
    },
  },
  methods: {
    async createAndMergePR() {
      this.creatingAndMergingPR = true
      await this.createOrUpdatePR()
        .then(this.triggerMergePullRequest)
        .catch(() => {
          this.creatingAndMergingPR = false
        })
    },

    async createOrUpdatePR() {
      await this.createOrUpdateGitHubPullRequest({
        workspaceID: this.workspace.id,
      }).catch((e) => {
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
    },
    async triggerMergePullRequest() {
      await this.mergeGitHubPullRequest({
        workspaceID: this.workspace.id,
      }).catch((e: any) => {
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
