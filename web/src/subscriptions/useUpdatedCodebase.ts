import { gql, useSubscription } from '@urql/vue'
import {
  UpdatedCodebaseCodebasesQuery,
  UpdatedCodebaseCodebasesQueryVariables,
  UpdatedCodebaseSubscription,
  UpdatedCodebaseSubscriptionVariables,
} from './__generated__/useUpdatedCodebase'
import { UpdateResolver } from '@urql/exchange-graphcache'
import { DeepMaybeRef } from '@vueuse/core'

const UPDATED_CODEBASE = gql`
  subscription UpdatedCodebase {
    updatedCodebase {
      id
      shortID
      name
      description
      inviteCode
      createdAt
      archivedAt
      lastUpdatedAt
      isReady

      members {
        id
        name
        avatarUrl
      }

      workspaces {
        id
      }
    }
  }
`

export function useUpdatedCodebase() {
  useSubscription<
    UpdatedCodebaseSubscription,
    UpdatedCodebaseSubscription,
    DeepMaybeRef<UpdatedCodebaseSubscriptionVariables>
  >({
    query: UPDATED_CODEBASE,
  })
}

const CODEBASES = gql`
  query UpdatedCodebaseCodebases {
    codebases {
      id
    }
  }
`

export const updatedCodebaseUpdateResolver: UpdateResolver<
  UpdatedCodebaseSubscription,
  UpdatedCodebaseSubscriptionVariables
> = (result, args, cache, info) => {
  cache.updateQuery<UpdatedCodebaseCodebasesQuery, UpdatedCodebaseCodebasesQueryVariables>(
    {
      query: CODEBASES,
    },
    (data) => {
      // Add codebase if not exists
      if (data && !data.codebases.some((c) => c.id === result.updatedCodebase.id)) {
        data.codebases.push(result.updatedCodebase)
      }
      return data
    }
  )
}
