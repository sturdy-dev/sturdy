<template>
  <Button
    v-if="mutagenIpc"
    :disabled="connecting"
    :icon="icon"
    :spinner="connecting"
    @click="createViewInDirectory"
  >
    Open in
  </Button>
</template>

<script lang="ts">
import { defineComponent, type PropType } from 'vue'
import { FolderOpenIcon } from '@heroicons/vue/outline'
import { useCreateWorkspace } from '../../mutations/useCreateWorkspace'
import { useArchiveWorkspace } from '../../mutations/useArchiveWorkspace'
import Button from '../../atoms/Button.vue'
import { gql } from '@urql/vue'
import type { ConnectNewDirectory_CodebaseFragment } from './__generated__/ConnectNewDirectory'

export const CODEBASE_FRAGMENT = gql`
  fragment ConnectNewDirectory_Codebase on Codebase {
    id
    name
    slug
  }
`

export default defineComponent({
  components: { Button },
  props: {
    codebase: { type: Object as PropType<ConnectNewDirectory_CodebaseFragment>, required: true },
  },
  setup() {
    const createWorkspaceResult = useCreateWorkspace()
    const archiveWorkspaceResult = useArchiveWorkspace()
    return {
      async createWorkspace(codebaseId: string) {
        return createWorkspaceResult({ codebaseID: codebaseId }).then(
          ({ createWorkspace }) => createWorkspace
        )
      },
      async archiveWorkspace(id: string) {
        return archiveWorkspaceResult({ id }).then(({ archiveWorkspace }) => archiveWorkspace)
      },
    }
  },
  data() {
    const mutagenIpc = window.mutagenIpc
    const ipc = window.ipc
    return {
      mutagenIpc,
      ipc,
      icon: FolderOpenIcon,
      connecting: false,
    }
  },
  methods: {
    async createViewInDirectory() {
      const oldIsReady = this.mutagenIpc?.isReady && (await this.mutagenIpc.isReady())
      const newIsReady = this.ipc?.state && (await this.ipc.state()) === 'online'

      const mutagenReady = oldIsReady || newIsReady

      if (!mutagenReady) {
        this.emitter.emit('notification', {
          title: 'Sturdy is not running',
          message: 'Sturdy is still starting, please wait.',
          style: 'error',
        })
        return
      }

      this.connecting = true
      const workspace = await this.createWorkspace(this.codebase.id)

      await this.mutagenIpc
        .createNewViewWithDialog(workspace.id, this.codebase.name)
        .then(() =>
          this.$router.push({
            name: 'workspaceHome',
            params: { codebaseSlug: this.codebase.slug, id: workspace.id },
          })
        )
        .catch(async (e: any) => {
          if (e.message.includes('non-empty')) {
            this.emitter.emit('notification', {
              title: 'Directory is not empty',
              message: 'Please select an empty directory.',
              style: 'error',
            })
            return
          } else if (e.message.includes('Cancelled')) {
            await this.archiveWorkspace(workspace.id)
            return
          } else {
            throw e
          }
        })
        .finally(() => {
          this.connecting = false
        })
    },
  },
})
</script>
