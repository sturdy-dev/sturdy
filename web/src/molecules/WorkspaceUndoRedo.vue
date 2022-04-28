<template>
  <div class="flex items-center" v-if="visible">
    <Button
      size="sm"
      class="flex items-center border-0 px-0.5 py-0.5"
      :spinner="undoing"
      :disabled="!previous"
      :icon="undoIcon"
      :class="{ 'opacity-50 cursor-not-allowed': !previous }"
      :show-tooltip="!!previous"
      @click.prevent.stop="onUndoClicked"
    >
      <template #tooltip>
        rewind to <RelativeTime v-if="previousDate" :date="previousDate" />
      </template>
      {{ '' }}
    </Button>

    <Button
      size="sm"
      class="flex items-center border-0 px-0.5 py-0.5"
      :disabled="!next"
      :spinner="redoing"
      :icon="redoIcon"
      :class="{ 'opacity-50 cursor-not-allowed': !next }"
      :show-tooltip="!!next"
      @click.prevent.stop="onRedoClicked"
    >
      <template #tooltip> rewind to <RelativeTime v-if="nextDate" :date="nextDate" /> </template>
      {{ '' }}
    </Button>
  </div>
</template>

<script lang="ts">
import { defineComponent, type PropType } from 'vue'
import { gql } from '@urql/vue'

import type { WorkspaceUndoRedo_WorkspaceFragment } from './__generated__/WorkspaceUndoRedo'

import { Undo, Redo } from '../atoms/icons'
import Button from '../atoms/Button.vue'
import RelativeTime from '../atoms/RelativeTime.vue'

import { useSetWorkspaceSnapshot } from '../mutations/useSetWorkspaceSnapshot'

export const WORKSPACE_FRAGMENT = gql`
  fragment WorkspaceUndoRedo_Workspace on Workspace {
    id
    snapshot {
      id
      previous {
        id
        createdAt
      }
      next {
        id
        createdAt
      }
    }
  }
`

export default defineComponent({
  components: {
    Button,
    RelativeTime,
  },
  props: {
    workspace: {
      type: Object as PropType<WorkspaceUndoRedo_WorkspaceFragment>,
      required: true,
    },
  },
  setup() {
    const setWorkspaceSnapshot = useSetWorkspaceSnapshot()
    return {
      setWorkspaceSnapshot,
    }
  },
  data() {
    return {
      visible: false,
      undoing: false,
      redoing: false,
      undoIcon: Undo,
      redoIcon: Redo,
    }
  },
  mounted() {
    window.onkeydown = (e) => {
      if (e.key === 'i' && e.ctrlKey) {
        this.visible = !this.visible
      }
    }
  },
  computed: {
    previous() {
      return this.workspace.snapshot?.previous
    },
    next() {
      return this.workspace.snapshot?.next
    },
    previousDate() {
      if (!this.previous) return null
      return new Date(this.previous.createdAt * 1000)
    },
    nextDate() {
      if (!this.next) return null
      return new Date(this.next.createdAt * 1000)
    },
  },
  methods: {
    onUndoClicked() {
      if (!this.previous) return
      this.undoing = true
      this.setWorkspaceSnapshot({
        workspaceID: this.workspace.id,
        snapshotID: this.previous.id,
      }).finally(() => {
        this.undoing = false
      })
    },
    onRedoClicked() {
      if (!this.next) return
      this.redoing = true
      this.setWorkspaceSnapshot({
        workspaceID: this.workspace.id,
        snapshotID: this.next.id,
      }).finally(() => {
        this.redoing = false
      })
    },
  },
})
</script>