import type { UpdateResolver } from '@urql/exchange-graphcache'
import { gql, useMutation } from '@urql/vue'
import type { DeepMaybeRef } from '@vueuse/core'
import type {
  ArchiveWorkspaceMutation,
  ArchiveWorkspaceMutationVariables,
} from './__generated__/useArchiveWorkspace'

const ARCHIVE_WORKSPACE = gql`
  mutation ArchiveWorkspace($id: ID!) {
    archiveWorkspace(id: $id) {
      id
      archivedAt
    }
  }
`

export function useArchiveWorkspace(): (
  input: DeepMaybeRef<ArchiveWorkspaceMutationVariables>
) => Promise<ArchiveWorkspaceMutation> {
  const { executeMutation } = useMutation<
    ArchiveWorkspaceMutation,
    DeepMaybeRef<ArchiveWorkspaceMutationVariables>
  >(ARCHIVE_WORKSPACE)

  return async (input): Promise<ArchiveWorkspaceMutation> => {
    const result = await executeMutation(input)
    if (result.error) {
      throw result.error
    }
    if (!result.data) throw new Error('no data returned')
    return result.data
  }
}

export const archiveWorkspaceUpdateResolver: UpdateResolver<
  ArchiveWorkspaceMutation,
  ArchiveWorkspaceMutationVariables
> = (result, args, cache, info) => {
  // not doing anything
}
