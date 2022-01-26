import { gql, useMutation } from '@urql/vue'
import { SetupGitHubRepositoryInput } from '../__generated__/types'
import { DeepMaybeRef } from '@vueuse/core'
import { UpdateResolver } from '@urql/exchange-graphcache'
import {
  SetupGitHubRepositoryGitHubRepositoryCodebasesQuery,
  SetupGitHubRepositoryGitHubRepositoryCodebasesQueryVariables,
  SetupGitHubRepositoryMutation,
  SetupGitHubRepositoryMutationVariables,
  SetupGitHubRepositoryOrganizationCodebasesQuery,
  SetupGitHubRepositoryOrganizationCodebasesQueryVariables,
} from './__generated__/useSetupGitHubRepository'

const SETUP_GITHUB_REPOSITORY = gql`
  mutation SetupGitHubRepository($input: SetupGitHubRepositoryInput!) {
    setupGitHubRepository(input: $input) {
      id
      shortID
      name
      isReady
    }
  }
`

const ORGANIZATION_CODEBASES = gql`
  query SetupGitHubRepositoryOrganizationCodebases {
    organizations {
      id
      codebases {
        id
      }
    }
  }
`

const GITHUB_REPOSITORY_CODEBASE = gql`
  query SetupGitHubRepositoryGitHubRepositoryCodebases {
    gitHubRepositories {
      id
      codebase {
        id
      }
    }
  }
`

export function useSetupGitHubRepository(): (
  input: DeepMaybeRef<SetupGitHubRepositoryInput>
) => Promise<SetupGitHubRepositoryMutation> {
  const { executeMutation } = useMutation<
    SetupGitHubRepositoryMutation,
    DeepMaybeRef<SetupGitHubRepositoryMutationVariables>
  >(SETUP_GITHUB_REPOSITORY)

  return async (input): Promise<SetupGitHubRepositoryMutation> => {
    const result = await executeMutation({ input })
    if (result.error) throw result.error
    if (!result.data) throw new Error('No data returned')
    return result.data
  }
}

export const setupGitHubUpdateResolver: UpdateResolver<
  SetupGitHubRepositoryMutation,
  SetupGitHubRepositoryMutationVariables
> = (result, args, cache, info) => {
  // Add codebase to list of codebases in organization
  if (args.input.organizationID) {
    cache.updateQuery<
      SetupGitHubRepositoryOrganizationCodebasesQuery,
      SetupGitHubRepositoryOrganizationCodebasesQueryVariables
    >(
      {
        query: ORGANIZATION_CODEBASES,
      },
      (data) => {
        if (data) {
          data.organizations = data.organizations.map((org) => {
            if (org.id === args.input.organizationID) {
              if (!org.codebases.some((cb) => cb.id === result.setupGitHubRepository.id)) {
                org.codebases.push(result.setupGitHubRepository)
              }
            }
            return org
          })
        }
        return data
      }
    )
  }

  // Add codebase to repository
  cache.updateQuery<
    SetupGitHubRepositoryGitHubRepositoryCodebasesQuery,
    SetupGitHubRepositoryGitHubRepositoryCodebasesQueryVariables
  >(
    {
      query: GITHUB_REPOSITORY_CODEBASE,
    },
    (data) => {
      if (data) {
        data.gitHubRepositories = data.gitHubRepositories.map((r) => {
          if (r.id === args.input.gitHubRepositoryID) {
            r.codebase = result.setupGitHubRepository
          }
          return r
        })
      }
      return data
    }
  )
}
