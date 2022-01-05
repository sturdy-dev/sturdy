import { gql, useMutation } from '@urql/vue'
import { Ref } from 'vue'
import { MergeGitHubPullRequestInput } from '../__generated__/types'
import { DeepMaybeRef } from '@vueuse/core'
import {
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
