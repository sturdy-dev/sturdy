import type { DeepMaybeRef, MaybeRef } from '@vueuse/core'
import { gql, useSubscription } from '@urql/vue'
import type { UpdateResolver } from '@urql/exchange-graphcache'
import type {
  UpdatedSuggestionSubscription,
  UpdatedSuggestionSubscriptionVariables,
  UpdatedWorkspaceSuggestionsQuery,
  UpdatedWorkspaceSuggestionsQueryVariables,
} from './__generated__/useUpdatedSuggestion'

const UPDATED_SUGGESTION = gql`
  subscription UpdatedSuggestion($workspaceID: ID!) {
    updatedSuggestion(workspaceID: $workspaceID) {
      id

      author {
        id
        name
        avatarUrl
      }

      dismissedAt

      workspace {
        id
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
  }
`

const WORKSPACE_SUGGESTIONS = gql`
  query UpdatedWorkspaceSuggestions($workspaceID: ID!) {
    workspace(id: $workspaceID) {
      id
      suggestions {
        id
        dismissedAt
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

export function useUpdatedSuggestion(workspaceID: MaybeRef<string>) {
  useSubscription<
    UpdatedSuggestionSubscription,
    UpdatedSuggestionSubscription,
    DeepMaybeRef<UpdatedSuggestionSubscriptionVariables>
  >({
    query: UPDATED_SUGGESTION,
    variables: { workspaceID },
  })
}

export const updatedSuggestionResolver: UpdateResolver<
  UpdatedSuggestionSubscription,
  UpdatedSuggestionSubscriptionVariables
> = (result, args, cache, info) => {
  cache.updateQuery<UpdatedWorkspaceSuggestionsQuery, UpdatedWorkspaceSuggestionsQueryVariables>(
    {
      query: WORKSPACE_SUGGESTIONS,
      variables: { workspaceID: args.workspaceID },
    },
    (data) => {
      if (!data) return data

      data.workspace.suggestions = [
        ...data.workspace.suggestions.filter(({ id }) => id != result.updatedSuggestion.id),
        result.updatedSuggestion,
      ]

      return data
    }
  )
}
