import { UpdateResolver } from '@urql/exchange-graphcache'
import { DeepMaybeRef } from '@vueuse/core'
import gql from 'graphql-tag'
import { useMutation } from '@urql/vue'
import { UpdateNotificationPreferenceInput } from '../__generated__/types'
import {
  UpdateNotificationPreferenceMutation,
  UpdateNotificationPreferenceMutationVariables,
} from './__generated__/useUpdateNotificationPreference'

const UPDATE_NOTIFICATION_PREFERENCE = gql`
  mutation UpdateNotificationPreference($input: UpdateNotificationPreferenceInput!) {
    updateNotificationPreference(input: $input) {
      type
      channel
      enabled
    }
  }
`
export function useUpdateNotificationPreference(): (
  input: DeepMaybeRef<UpdateNotificationPreferenceInput>
) => Promise<void> {
  const { executeMutation } = useMutation<
    UpdateNotificationPreferenceMutation,
    DeepMaybeRef<UpdateNotificationPreferenceMutationVariables>
  >(UPDATE_NOTIFICATION_PREFERENCE)
  return async (input) => {
    const result = await executeMutation({ input })

    if (result.error) {
      throw result.error
    }
  }
}

export const updateNotificationPreferenceResolver: UpdateResolver = (parent, args, cache, info) => {
  // Update cache manually if needed
}
