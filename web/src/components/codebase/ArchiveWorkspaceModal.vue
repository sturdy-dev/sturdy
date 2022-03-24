<template>
  <ConfirmModal
    title="Archive draft change"
    subtitle="Are you sure you want to archive this draft?"
    :show="isActive"
    @confirmed="deleteWorkspace"
    @close="close"
  />
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import { useArchiveWorkspace } from '../../mutations/useArchiveWorkspace'
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
    const archiveWorkspaceResult = useArchiveWorkspace()
    return {
      async archiveWorkspace(id: string) {
        return archiveWorkspaceResult({ id }).then(({ archiveWorkspace }) => archiveWorkspace)
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
    onkey(e: KeyboardEvent) {
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
