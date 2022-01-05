<template>
  <div class="flex items-center space-x-2">
    <Tooltip>
      <template #default>
        <EyeIcon
          v-if="watching"
          :class="buttonClass"
          aria-hidden="true"
          @click="unwatchWorkspace"
        />
        <EyeOffIcon v-else :class="buttonClass" aria-hidden="true" @click="watchWorkspace" />
      </template>

      <template v-if="watching" #tooltip>Stop watching this workspace</template>
      <template v-else #tooltip>Start watching this workspace</template>
    </Tooltip>
    <span v-if="watching" class="text-gray-900 text-sm font-medium"
      >You are watching this workspace</span
    >
    <span v-else class="text-gray-900 text-sm font-medium"
      >You are not watching this workspace</span
    >
  </div>
</template>

<script lang="ts">
import { defineComponent, toRefs } from 'vue'
import { gql } from '@urql/vue'
import { WorkspaceWatcherFragment, WorkspaceWatcherUserFragment } from './__generated__/Watching'
import Tooltip from '../../shared/Tooltip.vue'
import { EyeIcon, EyeOffIcon } from '@heroicons/vue/solid'
import { useWatchWorkspace } from '../../../mutations/useWatchWorkspace'
import { useUnwatchWorkspace } from '../../../mutations/useUnwatchWorkspace'

const USER_FRAGMENT = gql`
  fragment WorkspaceWatcherUser on User {
    id
    name
    avatarUrl
  }
`

export const WORKSPACE_WATCHER_FRAGMENT = gql`
  fragment WorkspaceWatcher on WorkspaceWatcher {
    user {
      ...WorkspaceWatcherUser
    }
  }
  ${USER_FRAGMENT}
`

export default defineComponent({
  components: {
    Tooltip,
    EyeIcon,
    EyeOffIcon,
  },
  props: {
    user: {
      type: Object as () => WorkspaceWatcherUserFragment,
      required: true,
    },
    watchers: {
      type: Array as () => WorkspaceWatcherFragment[],
      required: true,
    },
    workspaceId: {
      type: String,
      required: true,
    },
  },
  setup(props) {
    const { workspaceId } = toRefs(props)
    const watchWorkspaceResult = useWatchWorkspace()
    const unwatchWorkspaceResult = useUnwatchWorkspace()
    return {
      async watchWorkspace() {
        await watchWorkspaceResult({ workspaceID: workspaceId.value })
      },
      async unwatchWorkspace() {
        await unwatchWorkspaceResult({ workspaceID: workspaceId.value })
      },
    }
  },
  computed: {
    watching() {
      return this.watchers.some((w) => w.user.id === this.user.id)
    },
    buttonClass() {
      return 'h-5 w-5 text-gray-400 hover:bg-warmgray-200 rounded-md cursor-pointer transition'
    },
  },
})
</script>
