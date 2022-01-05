import { gql, useMutation } from '@urql/vue'
import { Ref } from 'vue'
import { CreateOrUpdateGitHubPullRequestInput } from '../__generated__/types'
import { DeepMaybeRef } from '@vueuse/core'
import {
  CreateOrUpdateGitHubPullRequestMutation,
  CreateOrUpdateGitHubPullRequestMutationVariables,
} from './__generated__/useCreateOrUpdateGitHubPullRequest'

const CREATE_OR_UPDATE_GITHUB_PULL_REQUEST = gql<
  CreateOrUpdateGitHubPullRequestMutation,
  DeepMaybeRef<CreateOrUpdateGitHubPullRequestMutationVariables>
>`
  mutation CreateOrUpdateGitHubPullRequest($input: CreateOrUpdateGitHubPullRequestInput!) {
    createOrUpdateGitHubPullRequest(input: $input) {
      id
      workspace {
        id
        upToDateWithTrunk
        gitHubPullRequest {
          id
          pullRequestNumber
          open
          merged
          mergedAt
        }
      }
    }
  }
`

export function useCreateOrUpdateGitHubPullRequest(): {
  mutating: Ref<boolean>
  createOrUpdateGitHubPullRequest(
    input: DeepMaybeRef<CreateOrUpdateGitHubPullRequestInput>
  ): Promise<void>
} {
  const { executeMutation, fetching: mutating } = useMutation(CREATE_OR_UPDATE_GITHUB_PULL_REQUEST)

  return {
    mutating,
    async createOrUpdateGitHubPullRequest(
      input: DeepMaybeRef<CreateOrUpdateGitHubPullRequestInput>
    ) {
      const result = await executeMutation({ input })
      if (result.error) {
        throw result.error
      }
    },
  }
}
