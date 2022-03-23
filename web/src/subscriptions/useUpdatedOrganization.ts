import { gql, useSubscription } from '@urql/vue'
import type { DeepMaybeRef, MaybeRef } from '@vueuse/core'
import type {
  UpdatedOrganizationSubscription,
  UpdatedOrganizationSubscriptionVariables,
} from './__generated__/useUpdatedOrganization'
import type { UpdateResolver } from '@urql/exchange-graphcache'

const UPDATED_ORGANIZATION = gql`
  subscription UpdatedOrganization($organizationID: ID) {
    updatedOrganization(organizationID: $organizationID) {
      id
      name
    }
  }
`

export function useUpdatedOrganization(
  organizationID?: MaybeRef<string>,
  opts?: { pause?: MaybeRef<boolean> }
) {
  return useUpdatedOrganizationSubscription({ organizationID: organizationID }, opts)
}

function useUpdatedOrganizationSubscription(
  variables: DeepMaybeRef<UpdatedOrganizationSubscriptionVariables>,
  { pause = false }: { pause?: MaybeRef<boolean> } = {}
) {
  useSubscription<
    UpdatedOrganizationSubscription,
    UpdatedOrganizationSubscription,
    DeepMaybeRef<UpdatedOrganizationSubscriptionVariables>
  >({
    query: UPDATED_ORGANIZATION,
    variables,
    pause,
  })
}

export const updatedOrganizationUpdateResolver: UpdateResolver<
  UpdatedOrganizationSubscription,
  UpdatedOrganizationSubscriptionVariables
> = (parent, args, cache, info) => {
  // not doing anything
}
