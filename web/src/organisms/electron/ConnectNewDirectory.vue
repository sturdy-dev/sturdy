<template>
  <Select v-if="mutagenIpc && version > 0" id="codebase-open-buttons" color="blue">
    <template #selected="{ option }">
      <component
        color="blue"
        :is="option"
        :disabled="connecting"
        :spinner="connecting"
        :show-tooltip="true"
        tooltip-position="bottom"
      />
    </template>

    <template #options>
      <Button
        :icon="plusIcon"
        class="text-sm text-left py-2 px-4 flex border-0 hover:bg-gray-50"
        @click="onClick(true)"
      >
        Connect new directory
        <template #tooltip>
          Create a new directory on your hard drive with this codebase's contents
        </template>
      </Button>

      <Button
        :icon="folderIcon"
        class="text-sm text-left py-2 px-4 flex border-0 hover:bg-gray-50"
        @click="onClick(false)"
      >
        Open existing directory
        <template #tooltip> Open this codebase in an existing directory on your computer </template>
      </Button>
    </template>
  </Select>

  <Button
    v-else-if="mutagenIpc"
    :disabled="connecting"
    :icon="folderIcon"
    :spinner="connecting"
    @click="onClick(true)"
  >
    Open in
  </Button>
</template>

<script lang="ts">
import { defineComponent, type PropType, ref } from 'vue'
import { FolderOpenIcon, PlusIcon } from '@heroicons/vue/outline'
import { gql } from '@urql/vue'

import Select from '../../atoms/Select.vue'
import Button from '../../atoms/Button.vue'

import type { ConnectNewDirectory_CodebaseFragment } from './__generated__/ConnectNewDirectory'

import { useCreateWorkspace } from '../../mutations/useCreateWorkspace'
import { useArchiveWorkspace } from '../../mutations/useArchiveWorkspace'

export const CODEBASE_FRAGMENT = gql`
  fragment ConnectNewDirectory_Codebase on Codebase {
    id
    name
    slug
  }
`

export default defineComponent({
  components: { Button, Select },
  props: {
    codebase: { type: Object as PropType<ConnectNewDirectory_CodebaseFragment>, required: true },
  },
  setup() {
    const createWorkspaceResult = useCreateWorkspace()
    const archiveWorkspaceResult = useArchiveWorkspace()
    const mutagenIpc = window.mutagenIpc
    const ipc = window.ipc
    const version = ref(0)
    if (mutagenIpc?.version) mutagenIpc.version().then((v: number) => (version.value = v))
    return {
      version,
      mutagenIpc,
      ipc,
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
    return {
      folderIcon: FolderOpenIcon,
      plusIcon: PlusIcon,
      connecting: false,
    }
  },
  methods: {
    async onClick(newDir: boolean) {
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
        .createNewViewWithDialog(workspace.id, this.codebase.name, newDir)
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
            await this.archiveWorkspace(workspace.id)
            return
          } else if (e.message.includes('already exists')) {
            this.emitter.emit('notification', {
              title: 'Directory is already connected',
              message: 'You can connect directory twice in Sturdy.',
              style: 'error',
            })
            await this.archiveWorkspace(workspace.id)
            return
          } else if (e.message.includes('Cancelled')) {
            await this.archiveWorkspace(workspace.id)
            return
          } else {
            await this.archiveWorkspace(workspace.id)
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
