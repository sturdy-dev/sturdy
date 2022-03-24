import { gql, useMutation } from '@urql/vue'
import type { UpdateResolver } from '@urql/exchange-graphcache'
import type { DeepMaybeRef } from '@vueuse/core'
import type {
  CreateOrUpdateCodebaseRemoteMutation,
  CreateOrUpdateCodebaseRemoteMutationVariables,
} from './__generated__/useCreateOrUpdateGitRemote'
import type { CreateOrUpdateCodebaseRemoteInput } from '../__generated__/types'

const CREATE_GIT_REMOTE = gql`
  mutation CreateOrUpdateCodebaseRemote($input: CreateOrUpdateCodebaseRemoteInput!) {
    createOrUpdateCodebaseRemote(input: $input) {
      id
      name
      url
      basicAuthUsername
      basicAuthPassword
      trackedBranch
      browserLinkRepo
      browserLinkBranch
      keyPair {
        id
        publicKey
      }
    }
  }
`

export function useCreateOrUpdateCodebaseRemote(): (
  input: DeepMaybeRef<CreateOrUpdateCodebaseRemoteInput>
) => Promise<CreateOrUpdateCodebaseRemoteMutation> {
  const { executeMutation } = useMutation<
    CreateOrUpdateCodebaseRemoteMutation,
    DeepMaybeRef<CreateOrUpdateCodebaseRemoteMutationVariables>
  >(CREATE_GIT_REMOTE)
  return async (input) => {
    const result = await executeMutation({ input })
    if (result.error) throw result.error
    if (result.data) {
      return result.data
    }
    throw new Error('unexpected result')
  }
}

export const createOrUpdateCodebaseRemoteUpdateResolver: UpdateResolver<
  CreateOrUpdateCodebaseRemoteMutation,
  CreateOrUpdateCodebaseRemoteMutationVariables
> = (result, args, cache, info) => {
  if (result.createOrUpdateCodebaseRemote.__typename) {
    const codebaseKey = cache.keyOfEntity({ __typename: 'Codebase', id: args.input.codebaseID })
    const remoteKey = cache.keyOfEntity({
      __typename: result.createOrUpdateCodebaseRemote.__typename,
      id: result.createOrUpdateCodebaseRemote.id,
    })
    cache.link(codebaseKey, 'remote', remoteKey)
  }
}
