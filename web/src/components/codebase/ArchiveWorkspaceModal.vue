<template>
  <ConfirmModal
    title="Archive workspace"
    subtitle="Are you sure you want to archive this workspace?"
    :show="isActive"
    @confirmed="deleteWorkspace"
    @close="close"
  />
</template>

<script lang="ts">
import { gql, useMutation } from '@urql/vue'
import { defineComponent } from 'vue'
import ConfirmModal from '../../molecules/ConfirmModal.vue'

export default defineComponent({
  components: {
    ConfirmModal,
  },
  props: {
    isActive: {
      type: Boolean,
      required: true,
    },
    workspaceId: {
      type: String,
      required: true,
    },
  },
  emits: ['close', 'archived', 'archiving'],
  setup() {
    const { executeMutation: archiveWorkspaceResult } = useMutation(gql`
      mutation ArchiveWorkspaceModal($id: ID!) {
        archiveWorkspace(id: $id) {
          id
          archivedAt
        }
      }
    `)

    return {
      async archiveWorkspace(id) {
        const variables = { id }
        await archiveWorkspaceResult(variables).then((result) => {
          if (result.error) {
            throw new Error(result.error)
          }
        })
      },
    }
  },
  computed: {
    modal() {
      return {
        'is-active': this.isActive,
        modal: true,
      }
    },
  },
  mounted() {
    window.addEventListener('keydown', this.onkey)
  },
  unmounted() {
    window.addEventListener('keydown', this.onkey)
  },
  methods: {
    onkey(e) {
      if (!this.isActive) {
        return
      }

      const escape = e.keyCode === 27
      const enter = e.keyCode === 13

      e.stopPropagation()
      e.preventDefault()

      switch (true) {
        case escape:
          this.close()
          break
        case enter:
          this.deleteWorkspace()
          break
      }
    },
    deleteWorkspace() {
      this.$emit('archiving')
      this.archiveWorkspace(this.workspaceId).then(() => {
        this.$emit('archived')
        this.close()
      })
    },
    close() {
      this.$emit('close')
    },
  },
})
</script>
