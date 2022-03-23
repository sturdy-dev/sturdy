import { gql, useMutation } from '@urql/vue'
import type {
  TriggerInstantIntegrationMutation,
  TriggerInstantIntegrationMutationVariables,
} from './__generated__/useTriggerInstantIntegration'
import type { DeepMaybeRef } from '@vueuse/core'
import type { UpdateResolver } from '@urql/exchange-graphcache'
import type { TriggerInstantIntegrationInput, Status } from '../__generated__/types'

const TRIGGER_INSTANT_INTEGRATION = gql`
  mutation TriggerInstantIntegration($input: TriggerInstantIntegrationInput!) {
    triggerInstantIntegration(input: $input) {
      id
      type
      title
      description
      timestamp
      detailsUrl
    }
  }
`

export function useTriggerInstantIntegration(): (
  input: DeepMaybeRef<TriggerInstantIntegrationInput>
) => Promise<Status[]> {
  const { executeMutation } = useMutation<
    TriggerInstantIntegrationMutation,
    DeepMaybeRef<TriggerInstantIntegrationMutationVariables>
  >(TRIGGER_INSTANT_INTEGRATION)
  return async (input): Promise<Status[]> => {
    const result = await executeMutation({ input })
    if (result.error) throw result.error
    return result.data ? result.data.triggerInstantIntegration : []
  }
}

export const triggerInstantIntegrationUpdateResolver: UpdateResolver<
  TriggerInstantIntegrationMutation,
  TriggerInstantIntegrationMutationVariables
> = (result, args, cache, info) => {
  // update cache manually if needed
}
