import type { DeepMaybeRef } from '@vueuse/core'
import { gql, useMutation } from '@urql/vue'
import type { UpdateResolver } from '@urql/exchange-graphcache'
import type { CreateSuggestionInput } from '../__generated__/types'
import type {
  CreateSuggestionMutation,
  CreateSuggestionMutationVariables,
  CreateSuggestionSuggestionFragment,
  WorkspaceSuggestionsQuery,
  WorkspaceSuggestionsQueryVariables,
} from './__generated__/useCreateSuggestion'

const CREATE_SUGGESTION_WORKSPACE_FRAGMENT = gql`
  fragment CreateSuggestionWorkspace on Workspace {
    id
  }
`

const CREATE_SUGGESTION_SUGGESTION_FRAGMENT = gql`
  fragment CreateSuggestionSuggestion on Suggestion {
    id
    author {
      id
      name
      avatarUrl
    }
    workspace {
      ...CreateSuggestionWorkspace
    }
    diffs {
      id

      origName
      newName
      preferredName

      isDeleted
      isNew
      isMoved

      hunks {
        id
        patch

        isOutdated
        isApplied
        isDismissed
      }
    }
  }
  ${CREATE_SUGGESTION_WORKSPACE_FRAGMENT}
`

const CREATE_SUGGESTION = gql`
  mutation CreateSuggestion($input: CreateSuggestionInput!) {
    createSuggestion(input: $input) {
      ...CreateSuggestionSuggestion
    }
  }
  ${CREATE_SUGGESTION_SUGGESTION_FRAGMENT}
`

const WORKSPACE_SUGGESTIONS = gql`
  query WorkspaceSuggestions($workspaceID: ID!) {
    workspace(id: $workspaceID) {
      id
      suggestions {
        id
        author {
          id
        }
        diffs {
          id
        }
      }
    }
  }
`

export function useCreateSuggestion(): (
  input: DeepMaybeRef<CreateSuggestionInput>
) => Promise<CreateSuggestionSuggestionFragment> {
  const { executeMutation } = useMutation<
    CreateSuggestionMutation,
    DeepMaybeRef<CreateSuggestionMutationVariables>
  >(CREATE_SUGGESTION)
  return async (input) => {
    const result = await executeMutation({ input })
    if (result.error) throw result.error
    if (!result.data) throw new Error('No data returned')
    return result.data.createSuggestion
  }
}

export const createSuggestionUpdateResolver: UpdateResolver<
  CreateSuggestionMutation,
  CreateSuggestionMutationVariables
> = (result, args, cache, info) => {
  if (!result) {
    return
  }

  cache.updateQuery<WorkspaceSuggestionsQuery, WorkspaceSuggestionsQueryVariables>(
    {
      query: WORKSPACE_SUGGESTIONS,
      variables: {
        workspaceID: result.createSuggestion.workspace.id,
      },
    },
    (data) => {
      if (!data) return data

      const suggestionExists = data.workspace.suggestions.some(
        ({ id }) => id === result.createSuggestion.id
      )

      if (suggestionExists) return data

      data.workspace.suggestions.push(result.createSuggestion)

      return data
    }
  )
}
