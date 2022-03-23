<template>
  <a
    v-if="pullRequestLink"
    :href="pullRequestLink"
    class="flex items-center space-x-2"
    target="_blank"
  >
    <StatusBadge v-if="displayStatus" :statuses="gitHubPullRequest.statuses" />
    <PuzzleIcon
      v-else
      class="h-5 w-5"
      aria-hidden="true"
      :class="[gitHubPullRequest.open ? 'text-green-500' : 'text-gray-400']"
    />

    <span
      class="text-sm font-medium"
      :class="[gitHubPullRequest.open ? 'text-green-700' : 'text-gray-900']"
    >
      <span v-if="gitHubPullRequest.open">Pull request</span>
      <span v-else-if="gitHubPullRequest.merged">Merged pull request</span>
      <span v-else>Pull request</span>

      <span
        class="font-medium"
        :class="[gitHubPullRequest.open ? 'text-green-700' : 'text-gray-500']"
      >
        #{{ gitHubPullRequest.pullRequestNumber }}
      </span>
    </span>
  </a>
</template>

<script lang="ts">
import { PuzzleIcon } from '@heroicons/vue/solid'
import { gql } from '@urql/vue'
import { toRefs } from 'vue'
import type { PropType } from 'vue'
import type {
  GitHubPullRequestFragment,
  CodebaseGitHubIntegrationFragment,
} from './__generated__/GitHubPullRequest'
import StatusBadge, { STATUS_FRAGMENT } from '../../statuses/StatusBadge.vue'
import { useUpdatedGitHubPullRequestStatuses } from '../../../subscriptions/useUpdatedGitHubPullRequestStatuses'

export const GITHUB_PULL_REQUEST_FRAGMENT = gql`
  fragment GitHubPullRequest on GitHubPullRequest {
    id
    pullRequestNumber
    open
    merged
    statuses {
      ...Status
    }
  }
  ${STATUS_FRAGMENT}
`

export const CODEBASE_GITHUB_INTEGRATION_FRAGMENT = gql`
  fragment CodebaseGitHubIntegration on CodebaseGitHubIntegration {
    id
    enabled
    owner
    name
    gitHubIsSourceOfTruth
  }
`

export default {
  components: { PuzzleIcon, StatusBadge },
  props: {
    gitHubIntegration: {
      type: Object as PropType<CodebaseGitHubIntegrationFragment>,
      default: null,
    },
    gitHubPullRequest: {
      type: Object as PropType<GitHubPullRequestFragment>,
      default: null,
    },
  },
  setup(props) {
    const { gitHubPullRequest } = toRefs(props)
    if (gitHubPullRequest && gitHubPullRequest.value) {
      useUpdatedGitHubPullRequestStatuses(gitHubPullRequest.value.id)
    }
  },
  computed: {
    displayStatus(): boolean {
      if (!this.gitHubPullRequest) return false
      return this.gitHubPullRequest.open && this.gitHubPullRequest.statuses?.length > 0
    },
    pullRequestLink(): string | null {
      if (!this.gitHubIntegration) return null
      if (!this.gitHubPullRequest) return null
      if (!this.gitHubIntegration.gitHubIsSourceOfTruth) return null
      if (!this.gitHubIntegration.enabled) return null

      const owner = this.gitHubIntegration.owner
      const name = this.gitHubIntegration.name
      const number = this.gitHubPullRequest.pullRequestNumber
      return `https://github.com/${owner}/${name}/pull/${number}`
    },
  },
}
</script>
