<template>
  <div v-if="change">
    <div class="flex flex-1 items-start flex-col-reverse md:flex-row gap-2 mb-4">
      <Editor
        v-if="description"
        ref="editor"
        :model-value="description"
        :editable="true"
        @updated="onUpdatedDescription"
      />

      <div v-if="isAuthorized" class="flex md:flex-col flex-shrink-0 gap-2">
        <Button size="wider" :icon="plusIcon" color="white" @click="createWorkspaceHandler(false)">
          New draft on this
        </Button>

        <Button size="wider" :icon="minusIcon" color="white" @click="createWorkspaceHandler(true)">
          Revert change
        </Button>
      </div>
    </div>

    <div v-if="justSaved" class="text-gray-400 text-sm my-2">Saved!</div>

    <Change
      :change="change"
      :comments="nonArchivedComments"
      :user="user"
      :members="change.codebase.members"
      :show-full-file-button="true"
    />
  </div>
</template>

<script lang="ts">
import http from '../../http'
import Change from '../differ/Change.vue'
import debounce from '../../debounce'
import { defineAsyncComponent } from 'vue'
import type { PropType } from 'vue'
import Button from '../shared/Button.vue'
import { MinusIcon, PlusIcon } from '@heroicons/vue/outline'
import { useCreateWorkspace } from '../../mutations/useCreateWorkspace'
import { gql } from '@urql/vue'
import type { ChangeDetails_ChangeFragment } from './__generated__/ChangeDetails'
import { MEMBER_FRAGMENT } from '../../components/shared/TextareaMentions.vue'

export const CHANGE_DETAILS_CHANGE_FRAGMENT = gql`
  fragment ChangeDetails_Change on Change {
    id
    description

    codebase {
      id
      members {
        ...Member
      }
    }

    diffs {
      id
      origName
      newName
      preferredName

      isDeleted
      isNew
      isMoved

      isLarge
      largeFileInfo {
        id
        size
      }

      isHidden

      hunks {
        id
        patch
      }
    }

    comments {
      id
      message
      codeContext {
        id
        path
        lineEnd
        lineStart
        lineIsNew
      }
      createdAt
      deletedAt
      author {
        id
        name
        avatarUrl
      }
      replies {
        id
        message
        createdAt
        author {
          id
          name
          avatarUrl
        }
      }
    }
  }
  ${MEMBER_FRAGMENT}
`

type Member = ChangeDetails_ChangeFragment['codebase']['members'][number]

export default {
  components: {
    Change,
    Button,
    Editor: defineAsyncComponent(() => import('../workspace/Editor.vue')),
  },
  props: {
    change: {
      type: Object as PropType<ChangeDetails_ChangeFragment>,
      required: true,
    },
    user: {
      type: Object as PropType<Member>,
    },
  },
  emits: ['ready', 'approve', 'land', 'drop', 'revert', 'undo', 'commented'],
  setup() {
    const createWorkspaceResult = useCreateWorkspace()

    return {
      createWorkspace(
        codebaseID: string,
        onTopOfChange?: string,
        onTopOfChangeWithRevert?: string
      ) {
        return createWorkspaceResult({
          codebaseID,
          onTopOfChange,
          onTopOfChangeWithRevert,
        })
      },

      plusIcon: PlusIcon,
      minusIcon: MinusIcon,
    }
  },
  data() {
    return {
      description: null as string | null,
      editingDescription: false,
      justSaved: false,
      unsetJustSavedFunc: null as (() => void) | null,
      updateDescriptionDebounceFunc: null as (() => void) | null,
    }
  },
  computed: {
    isAuthenticated() {
      return !!this.user
    },
    isAuthorized() {
      const isMember = this.change.codebase.members.some(({ id }) => id === this.user?.id)
      return this.isAuthenticated && isMember
    },
    nonArchivedComments() {
      return this.change?.comments.filter((c) => !c.deletedAt)
    },
  },
  watch: {
    'change.id': function () {
      this.description = this.change?.description
    },
  },
  mounted() {
    this.updateDescriptionDebounceFunc = debounce(this.save, 300)
    this.unsetJustSavedFunc = debounce(() => {
      this.justSaved = false
    }, 4000)

    this.description = this.change?.description
  },
  methods: {
    onUpdatedDescription(ev: { content: string }) {
      this.description = ev.content
      if (this.updateDescriptionDebounceFunc) this.updateDescriptionDebounceFunc()
    },

    async createWorkspaceHandler(withRevert: boolean) {
      const onTopOfChangeWithRevert = withRevert ? this.change.id : undefined
      const onTopOfChange = withRevert ? undefined : this.change.id

      await this.createWorkspace(
        this.change.codebase.id,
        onTopOfChange,
        onTopOfChangeWithRevert
      ).then((result) => {
        this.$router.push({
          name: 'workspaceHome',
          params: { id: result.createWorkspace.id },
        })
      })
    },
    save() {
      fetch(http.url('v3/changes/' + this.change.id + '/update'), {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include',
        body: JSON.stringify({
          updated_description: this.description,
        }),
      })
        .then(http.checkStatus)
        .then((response) => response.json())
        .then(() => {
          this.justSaved = true
          if (this.unsetJustSavedFunc) this.unsetJustSavedFunc()
        })
    },
    descriptionKeydown(e: KeyboardEvent) {
      // Cmd + Enter
      if ((e.metaKey || e.ctrlKey) && e.keyCode === 13) {
        this.editingDescription = false
        this.save()
      }
    },
    startStopEdit() {
      // Stop editing
      if (this.editingDescription) {
        this.editingDescription = false
        return
      }

      this.editingDescription = true

      // Focus the input field
      this.$nextTick(() => {
        this.focusEdit()
      })
    },
    focusEdit() {
      const editor = this.$refs['editor'] as HTMLElement
      editor.focus()
    },
  },
}
</script>
