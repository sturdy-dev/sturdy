import { CreateWorkspaceInput } from '../__generated__/types'
import { DeepMaybeRef } from '@vueuse/core'
import { UpdateResolver } from '@urql/exchange-graphcache'
import { gql, useMutation } from '@urql/vue'
import {
  CreateWorkspaceMutation,
  CreateWorkspaceMutationVariables,
} from './__generated__/useCreateWorkspace'

const CREATE_WORKSPACE = gql`
  mutation CreateWorkspace($input: CreateWorkspaceInput!) {
    createWorkspace(input: $input) {
      id
      name
    }
  }
`

export function useCreateWorkspace(): (
  input: DeepMaybeRef<CreateWorkspaceInput>
) => Promise<CreateWorkspaceMutation> {
  const { executeMutation } = useMutation<
    CreateWorkspaceMutation,
    DeepMaybeRef<CreateWorkspaceMutationVariables>
  >(CREATE_WORKSPACE)
  return async (input) => {
    const result = await executeMutation({ input })

    if (result.error) {
      throw result.error
    }

    console.log('create workspace', result)

    if (result.data) {
      return result.data
    }

    throw new Error('unexpected result')
  }
}

export const createWorkspaceUpdateResolver: UpdateResolver<
  CreateWorkspaceMutation,
  CreateWorkspaceMutationVariables
> = (parent, args, cache, info) => {
  // Update cache manually if needed
}
