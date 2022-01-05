import { DeepMaybeRef, MaybeRef } from '@vueuse/core'
import { gql, useSubscription } from '@urql/vue'
import {
  UpdatedViewAuthorViewsQuery,
  UpdatedViewAuthorViewsQueryVariables,
  UpdatedViewsSubscription,
  UpdatedViewsSubscriptionVariables,
} from './__generated__/useUpdatedViews'
import { UpdateResolver } from '@urql/exchange-graphcache'

const UPDATED_VIEWS = gql`
  subscription UpdatedViews {
    updatedViews {
      id
      lastUsedAt
      author {
        id
      }
      status {
        id
        state
        progressPath
        progressReceived
        progressTotal
        lastError
        sturdyVersion
      }
      workspace {
        id
      }
    }
  }
`

export function useUpdatedViews({ pause = false }: { pause?: MaybeRef<boolean> } = {}) {
  useSubscription<
    UpdatedViewsSubscription,
    UpdatedViewsSubscription,
    DeepMaybeRef<UpdatedViewsSubscriptionVariables>
  >({
    query: UPDATED_VIEWS,
    pause,
  })
}

const UPDATED_VIEWS_AUTHOR_VIEWS = gql`
  query UpdatedViewAuthorViews {
    user {
      id
      views {
        id
      }
    }
  }
`

export const updatedViewsUpdateResolver: UpdateResolver<
  UpdatedViewsSubscription,
  UpdatedViewsSubscriptionVariables
> = (parent, args, cache, info) => {
  // When a view is created, add it to the author's list of views (if they're the current user)
  cache.updateQuery<UpdatedViewAuthorViewsQuery, UpdatedViewAuthorViewsQueryVariables>(
    {
      query: UPDATED_VIEWS_AUTHOR_VIEWS,
    },
    (data) => {
      if (
        data?.user.id === parent.updatedViews.author.id &&
        !data.user.views.some((view) => view.id === parent.updatedViews.id)
      ) {
        data.user.views.push(parent.updatedViews)
      }
      return data
    }
  )
}
