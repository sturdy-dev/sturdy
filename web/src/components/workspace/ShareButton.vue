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
            :disabled="creatingOrUpdatingPR"
            :class="[creatingOrUpdatingPR ? 'cursor-default' : '']"
            @click="createOrUpdatePR"
          >
            <div v-if="creatingOrUpdatingPR" class="flex items-center">
              <Spinner class="mr-1" />
              <span v-if="!hasOpenGitHubPR">Creating pull request</span>
              <span v-else> Updating pull request </span>
            </div>
            <span v-else-if="!hasOpenGitHubPR">Create pull request</span>
            <span v-else> Update pull request </span>
          </Button>

          <Button
            v-if="hasGitHubPR && hasOpenGitHubPR"
            color="green"
            :disabled="mergingGitHubPullRequest"
            @click="triggerMergePullRequest"
          >
            <div v-if="mergingGitHubPullRequest" class="flex items-center">
              <Spinner class="mr-1" />
              <span>Merging</span>
            </div>
            <span v-else>Merge</span>
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
        :disabled="landing"
        :class="[landing ? 'cursor-default' : '']"
        @click="shareChange"
      >
        <div v-if="landing" class="flex items-center">
          <Spinner class="mr-1" />
          <span>Sharing</span>
        </div>
        <span v-else>Share</span>
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
      open
    }
  }
`

export default defineComponent({
  name: 'ShareButton',
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
      type: String as PropType<string | undefined>,
      default: undefined,
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
      return this.workspace.gitHubPullRequest?.open
    },
    gitHubPRLink() {
      const { owner, name } = this.workspace.codebase.gitHubIntegration ?? {}
      const { pullRequestNumber } = this.workspace.gitHubPullRequest ?? {}
      return `https://github.com/${owner}/${name}/pull/${pullRequestNumber}`
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
