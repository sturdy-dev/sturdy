import type { UpdateResolver } from '@urql/exchange-graphcache'
import { useMutation } from '@urql/vue'
import gql from 'graphql-tag'
import type { UpdateWorkspaceInput } from '../__generated__/types'
import type { DeepMaybeRef } from '@vueuse/core'
import type { Ref } from 'vue'

const UPDATE_WORKSPACE = gql`
  mutation UpdateWorkspace($input: UpdateWorkspaceInput!) {
    updateWorkspace(input: $input) {
      id
      name
      draftDescription
    }
  }
`

export function useUpdateWorkspace(): {
  mutating: Ref<boolean>
  updateWorkspace(input: DeepMaybeRef<UpdateWorkspaceInput>): Promise<void>
} {
  const { executeMutation, fetching } = useMutation(UPDATE_WORKSPACE)
  return {
    mutating: fetching,
    updateWorkspace: async (input) => {
      const result = await executeMutation({ input })

      if (result.error) {
        throw result.error
      }
    },
  }
}

export const updateWorkspaceUpdateResolver: UpdateResolver = (parent, args, cache, info) => {
  // Update cache manually if needed
}
