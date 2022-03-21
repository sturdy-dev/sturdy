import { gql, useMutation } from '@urql/vue'
import { Ref } from 'vue'
import { PushWorkspaceInput } from '../__generated__/types'
import { DeepMaybeRef } from '@vueuse/core'
import {
  PushWorkspaceMutation,
  PushWorkspaceMutationVariables,
} from './__generated__/usePushWorkspace'

const PUSH_WORKSPACE = gql<PushWorkspaceMutation, DeepMaybeRef<PushWorkspaceMutationVariables>>`
  mutation PushWorkspace($input: PushWorkspaceInput!) {
    pushWorkspace(input: $input) {
      id
    }
  }
`

export function usePushWorkspace(): {
  mutating: Ref<boolean>
  pushWorkspace(input: DeepMaybeRef<PushWorkspaceInput>): Promise<void>
} {
  const { executeMutation, fetching: mutating } = useMutation(PUSH_WORKSPACE)

  return {
    mutating,
    async pushWorkspace(input: DeepMaybeRef<PushWorkspaceInput>) {
      const result = await executeMutation({ input })
      if (result.error) {
        throw result.error
      }
    },
  }
}
