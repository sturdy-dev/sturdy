import { ViewStatusState, type OpenWorkspaceOnViewInput } from '../__generated__/types'
import type { OptimisticMutationResolver, UpdateResolver } from '@urql/exchange-graphcache'
import { gql, useMutation } from '@urql/vue'
import type {
  OpenWorkspaceOnViewMutation,
  OpenWorkspaceOnViewMutationVariables,
  UserOpeningWorkspaceQuery,
  UserOpeningWorkspaceQueryVariables,
  WorkspaceToBeOpenedFragment,
} from './__generated__/useOpenWorkspaceOnView'
import type { DeepMaybeRef } from '@vueuse/core'

const OPEN_WORKSPACE_ON_VIEW = gql`
  mutation OpenWorkspaceOnView($input: OpenWorkspaceOnViewInput!) {
    openWorkspaceOnView(input: $input) {
      id
      workspace {
        id
        view {
          id
        }
      }
      status {
        id
        state
      }
    }
  }
`

export function useOpenWorkspaceOnView(): (input: OpenWorkspaceOnViewInput) => Promise<void> {
  const { executeMutation } = useMutation<
    OpenWorkspaceOnViewMutation,
    DeepMaybeRef<OpenWorkspaceOnViewMutationVariables>
  >(OPEN_WORKSPACE_ON_VIEW)
  return async (input) => {
    const result = await executeMutation({ input })
    if (result.error) {
      throw result.error
    }
  }
}

export const openWorkspaceOnViewUpdateResolver: UpdateResolver<
  OpenWorkspaceOnViewMutation,
  OpenWorkspaceOnViewMutationVariables
> = (parent, args, cache, info) => {
  // Manually update cache if necessary
}

const USER_OPENING_WORKSPACE = gql`
  query UserOpeningWorkspace {
    user {
      id
    }
  }
`

const WORKSPACE_TO_BE_OPENED = gql`
  fragment WorkspaceToBeOpened on Workspace {
    id
    author {
      id
    }
  }
`

export const openWorkspaceOnViewOptimisticMutationResolver: OptimisticMutationResolver<
  OpenWorkspaceOnViewMutationVariables,
  OpenWorkspaceOnViewMutation['openWorkspaceOnView']
> = (vars, cache, info): OpenWorkspaceOnViewMutation['openWorkspaceOnView'] => {
  const result: OpenWorkspaceOnViewMutation['openWorkspaceOnView'] = {
    __typename: 'View',
    id: vars.input.viewID,
    workspace: null,
    status: {
      __typename: 'ViewStatus',
      id: vars.input.viewID,
      state: ViewStatusState.Transferring,
    },
  }

  const userQuery = cache.readQuery<UserOpeningWorkspaceQuery, UserOpeningWorkspaceQueryVariables>({
    query: USER_OPENING_WORKSPACE,
  })
  if (userQuery == null) {
    return result
  }
  const workspace = cache.readFragment<
    WorkspaceToBeOpenedFragment,
    { __typename: 'Workspace'; id: string }
  >(WORKSPACE_TO_BE_OPENED, {
    __typename: 'Workspace',
    id: vars.input.workspaceID,
  })

  if (workspace == null) {
    return result
  }

  result.workspace = {
    __typename: 'Workspace',
    id: vars.input.workspaceID,
    view: {
      __typename: 'View',
      id: vars.input.viewID,
    },
  }

  return result
}
