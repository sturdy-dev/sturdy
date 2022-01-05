<template>
  <div>
    <div class="flex flex-row-reverse">
      <div v-for="p in presenceToShow" :key="p.id" :class="['-mr-2 z-10 bg-white rounded-full']">
        <Tooltip>
          <template #default>
            <Avatar
              :author="p.author"
              size="8"
              :class="[
                p.state === 'Coding' && p.lastActiveAt >= now - 600 ? 'ring-2 ring-green-500' : '',
                p.state === 'Viewing' && p.lastActiveAt >= now - 600 ? '' : '',
                p.state === 'Idle' || p.lastActiveAt < now - 600 ? 'opacity-25' : '',
              ]"
            />
          </template>
          <template v-if="p.state === 'Coding' && p.lastActiveAt >= now - 600" #tooltip>
            {{ p.author.name }} is coding
          </template>
          <template v-else-if="p.state === 'Viewing' && p.lastActiveAt >= now - 600" #tooltip>
            {{ p.author.name }} is here
          </template>
          <template v-else #tooltip>{{ p.author.name }} is idle</template>
        </Tooltip>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import Avatar from '../shared/Avatar.vue'
import Tooltip from '../shared/Tooltip.vue'
import { gql, useMutation } from '@urql/vue'
import { computed, defineComponent, onUnmounted, PropType, ref, toRefs, watch } from 'vue'
import { PresencePartsFragment, WorkspacePartsFragment } from './__generated__/Presence'
import { useUpdatedWorkspacePresence } from '../../subscriptions/useUpdatedWorkspacePresence'

export const PRESENCE_FRAGMENT_QUERY = gql`
  fragment PresenceParts on WorkspacePresence {
    id
    author {
      id
      name
      avatarUrl
    }
    state
    lastActiveAt
  }
`

export const WORKSPACE_FRAGMENT_QUERY = gql`
  fragment WorkspaceParts on Workspace {
    id
  }
`

export default defineComponent({
  components: { Tooltip, Avatar },
  props: {
    presence: {
      type: Object as PropType<PresencePartsFragment[]>,
      required: true,
    },
    workspace: {
      type: Object as PropType<WorkspacePartsFragment>,
      required: true,
    },
    user: {
      type: Object,
    },
  },
  setup(props) {
    const { workspace, user } = toRefs(props)

    let now = ref(new Date() / 1000)
    let interval = setInterval(() => (now.value = new Date() / 1000), 1000 * 5)
    onUnmounted(() => clearInterval(interval))

    const { executeMutation: reportWorkspacePresence } = useMutation(gql`
      mutation ReportWorkspacePresence($workspaceID: ID!, $state: WorkspacePresenceState!) {
        reportWorkspacePresence(input: { workspaceID: $workspaceID, state: $state }) {
          id
        }
      }
    `)

    let isVisible = document.visibilityState === 'visible'
    let hasFocus = document.hasFocus()

    let report = () => {
      if (!user.value) return
      let vars = {
        workspaceID: workspace.value.id,
        state: isVisible && hasFocus ? 'Viewing' : 'Idle',
      }
      reportWorkspacePresence(vars)
    }

    let visibilityListener = () => {
      isVisible = document.visibilityState === 'visible'
      hasFocus = document.hasFocus()
      report()
    }

    document.addEventListener('visibilitychange', visibilityListener)
    window.addEventListener('focus', visibilityListener)
    window.addEventListener('blur', visibilityListener)
    onUnmounted(() => {
      document.removeEventListener('visibilitychange', visibilityListener)
      window.removeEventListener('focus', visibilityListener)
      window.removeEventListener('blur', visibilityListener)
    })

    visibilityListener()
    let reportInterval = setInterval(visibilityListener, 60 * 1000) // every minute
    onUnmounted(() => clearInterval(reportInterval))
    watch(workspace, (o, n) => {
      if (o?.id !== n?.id) {
        visibilityListener()
      }
    })

    let workspaceID = ref(workspace.value.id)
    watch(workspace, () => {
      if (workspaceID?.value?.id) {
        workspaceID = workspaceID?.value?.id
      }
    })

    useUpdatedWorkspacePresence(workspaceID, {
      pause: computed(() => !workspaceID.value),
    })

    return {
      now,
    }
  },
  computed: {
    presenceToShow() {
      // Show presences with activity in the last 15 minutes
      return this.presence
        .filter((p) => p?.lastActiveAt && p.lastActiveAt >= this.now - 60 * 15) // Hide old entries
        .filter((p) => p.author.id !== this.user?.id) // Hide yourself
        .sort((a, b) => a.author.name.localeCompare(b.author.name))
    },
  },
})
</script>
