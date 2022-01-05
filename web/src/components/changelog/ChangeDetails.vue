<template>
  <div v-if="change">
    <div class="flex xl:justify-between xl:space-x-4 items-start flex-col-reverse xl:flex-row mb-4">
      <Editor
        v-if="description"
        ref="editor"
        :model-value="description"
        :editable="true"
        @updated="onUpdatedDescription"
      />

      <div v-if="isAuthorized" class="flex space-x-4 flex-shrink-0 mb-4">
        <Button size="wider" @click="createWorkspaceHandler(true)">
          <MinusIcon class="-ml-1 mr-2 h-5 w-5 text-gray-400" aria-hidden="true" />
          <span>Revert</span>
        </Button>

        <Button size="wider" @click="createWorkspaceHandler(false)">
          <PlusIcon class="-ml-1 mr-2 h-5 w-5 text-gray-400" aria-hidden="true" />
          <span>Create Workspace</span>
        </Button>
      </div>
    </div>

    <div v-if="justSaved" class="text-gray-400 text-sm my-2">Saved!</div>

    <div v-if="change?.is_landed" class="mt-2 text-sm text-gray-500 flex items-center italic">
      This change is landed.
    </div>

    <Change
      :change="change"
      :comments="nonArchivedComments"
      :user="user"
      :members="members"
      :show-full-file-button="true"
      @submittedNewComment="onSubmittedNewComment"
    />
  </div>
</template>

<script>
import http from '../../http'
import Change from '../differ/Change.vue'
import time from '../../time'
import { IsFocusChildOfElementWithClass } from '../../focus'
import debounce from '../../debounce'
import { defineAsyncComponent } from 'vue'
import Button from '../shared/Button.vue'
import { MinusIcon, PlusIcon } from '@heroicons/vue/outline'
import { useCreateWorkspace } from '../../mutations/useCreateWorkspace'

export default {
  name: 'ChangeDetails',
  components: {
    Change,
    Button,
    PlusIcon,
    MinusIcon,
    Editor: defineAsyncComponent(() => import('../workspace/Editor.vue')),
  },
  props: {
    codebaseSlug: {
      type: String,
      required: true,
    },
    codebaseId: {
      type: String,
      required: true,
    },
    change: {},
    user: {
      type: Object,
    },
    members: {
      type: Array,
      required: true,
    },
  },
  emits: ['ready', 'approve', 'land', 'drop', 'revert', 'undo', 'commented'],
  setup() {
    const createWorkspaceResult = useCreateWorkspace()

    return {
      createWorkspace(codebaseID, onTopOfChange, onTopOfChangeWithRevert) {
        return createWorkspaceResult({
          codebaseID,
          onTopOfChange,
          onTopOfChangeWithRevert,
        })
      },
    }
  },
  data() {
    return {
      description: null,
      editingDescription: false,
      justSaved: false,
      unsetJustSavedFunc: null,
      updateDescriptionDebounceFunc: null,
    }
  },
  computed: {
    isAuthenticated() {
      return !!this.user
    },
    isAuthorized() {
      const isMember = this.members.some(({ id }) => id === this.user?.id)
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

    window.addEventListener('keydown', this.onkey)

    this.description = this.change?.description
  },
  methods: {
    onUpdatedDescription(ev) {
      this.description = ev.content
      this.updateDescriptionDebounceFunc()
    },

    friendly_ago(ts) {
      return time.getRelativeTime(new Date(ts * 1000))
    },
    async createWorkspaceHandler(withRevert) {
      let onTopOfChange = null
      let onTopOfChangeWithRevert = null

      if (withRevert) {
        onTopOfChangeWithRevert = this.change.id
      } else {
        onTopOfChange = this.change.id
      }

      await this.createWorkspace(this.codebaseId, onTopOfChange, onTopOfChangeWithRevert).then(
        (result) => {
          this.$router.push({
            name: 'workspaceHome',
            params: { id: result.createWorkspace.id },
          })
        }
      )
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
          this.unsetJustSavedFunc()
        })
        .catch((e) => {
          console.log(e)
        })
    },
    descriptionKeydown(e) {
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
        this.$refs.changeDescription.focus()
      })
    },
    descriptionBlur(event) {
      // Don't blur here if clicked the edit button (which does the same action)
      if (IsFocusChildOfElementWithClass(event, 'change-details-edit-button')) {
        return
      }

      this.editingDescription = false
    },
    onkey(e) {
      // Cmd + Z
      if ((e.metaKey || e.ctrlKey) && e.keyCode === 90) {
        if (this.canUndo) {
          this.$emit('undo', this.commitID)
        }
      }
    },
    focusEdit() {
      this.$refs.editor.editor.view.dom.focus()
    },
    onSubmittedNewComment() {
      this.$emit('commented')
    },
  },
}
</script>
