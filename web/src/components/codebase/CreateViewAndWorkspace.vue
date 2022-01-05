<template>
  <Button v-if="mutagenIpc" @click="createViewInDirectory">
    <div class="flex items-center px-1">
      <DesktopComputerIcon class="-ml-1 mr-2 h-5 w-5 text-gray-400" aria-hidden="true" />
      Connect directory
    </div>
  </Button>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import { DesktopComputerIcon } from '@heroicons/vue/outline'
import { useCreateWorkspace } from '../../mutations/useCreateWorkspace'
import Button from '../shared/Button.vue'

export default defineComponent({
  components: { DesktopComputerIcon, Button },
  props: ['codebaseId', 'codebaseSlug'],
  setup() {
    const createWorkspaceResult = useCreateWorkspace()
    return {
      mutagenIpc: window.mutagenIpc,
      ipc: window.ipc,
      createWorkspace(codebaseID: string) {
        return createWorkspaceResult({
          codebaseID,
        })
      },
    }
  },
  methods: {
    async createViewInDirectory() {
      if (!this.codebaseId || !this.codebaseSlug) {
        return
      }

      let oldIsReady = this.mutagenIpc?.isReady && (await this.mutagenIpc.isReady())
      let newIsReady = this.ipc?.state && (await this.ipc.state()) === 'online'

      let mutagenReady = oldIsReady || newIsReady

      if (!mutagenReady) {
        this.emitter.emit('notification', {
          title: 'Sturdy is not running',
          message: 'Sturdy is still starting, please wait.',
          style: 'error',
        })
        return
      }

      const res = await this.createWorkspace(this.codebaseId)

      try {
        await this.mutagenIpc.createNewViewWithDialog(res.createWorkspace.id)
      } catch (e) {
        if (e.message.includes('non-empty')) {
          this.emitter.emit('notification', {
            title: 'Directory is not empty',
            message: 'Please select an empty directory.',
            style: 'error',
          })
          return
        } else if (e.message.includes('Cancelled')) {
          // User cancelled the dialog. Do nothing
          return
        } else {
          throw e
        }
      }

      await this.$router.push({
        name: 'workspaceHome',
        params: { codebaseSlug: this.codebaseSlug, id: res.createWorkspace.id },
      })
    },
  },
})
</script>
