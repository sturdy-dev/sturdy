import { gql, useMutation } from '@urql/vue'
import type { Ref } from 'vue'
import type { MergeGitHubPullRequestInput } from '../__generated__/types'
import type { DeepMaybeRef } from '@vueuse/core'
import type {
  MergeGitHubPullRequestMutation,
  MergeGitHubPullRequestMutationVariables,
} from './__generated__/useMergeGitHubPullRequest'

const MERGE_GITHUB_PULL_REQUEST = gql<
  MergeGitHubPullRequestMutation,
  DeepMaybeRef<MergeGitHubPullRequestMutationVariables>
>`
  mutation MergeGitHubPullRequest($input: MergeGitHubPullRequestInput!) {
    mergeGitHubPullRequest(input: $input) {
      id
      open
      merged
      mergedAt
    }
  }
`

export function useMergeGitHubPullRequest(): {
  mutating: Ref<boolean>
  mergeGitHubPullRequest(input: DeepMaybeRef<MergeGitHubPullRequestInput>): Promise<void>
} {
  const { executeMutation, fetching: mutating } = useMutation(MERGE_GITHUB_PULL_REQUEST)

  return {
    mutating,
    async mergeGitHubPullRequest(input: DeepMaybeRef<MergeGitHubPullRequestInput>) {
      const result = await executeMutation({ input })
      if (result.error) {
        throw result.error
      }
    },
  }
}
