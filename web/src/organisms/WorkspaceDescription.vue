<template>
  <div>
    <h2 class="sr-only">Description</h2>
    <OnboardingStep id="MakingAChange" :dependencies="['FindingYourWorkspace', 'WorkspaceChanges']">
      <template #title>Publishing a Change</template>
      <template #description>
        When you've made edits to your files and feel ready to make a checkpoint, write a
        description of your change(s) here.
      </template>

      <Editor
        :model-value="description"
        :editable="isAuthorized"
        placeholder="Describe your changes&hellip;"
        @updated="onUpdatedDescription"
      >
        <transition
          v-if="isAuthenticated"
          enter-active-class="transition ease-out duration-75"
          enter-from-class="opacity-0 scale-75"
          enter-to-class="opacity-100 scale-100"
          leave-active-class="transition ease-in duration-75"
          leave-from-class="opacity-100 scale-100"
          leave-to-class="opacity-0 scale-75"
        >
          <ShareButton
            :workspace="workspace"
            :all-hunk-ids="diffIds"
            :disabled="!canSubmitChange"
            :cant-submit-reason="cantSubmitChangeReason"
          />
        </transition>
      </Editor>
    </OnboardingStep>

    <div
      v-if="saving"
      class="hidden xl:block text-gray-400 text-sm absolute bottom-full translate-y-1 right-0 origin-bottom-right animate-pulse"
    >
      Saving...
    </div>
    <transition
      enter-active-class="transition ease-out duration-50"
      enter-from-class="opacity-0 scale-75"
      enter-to-class="opacity-100 scale-100"
      leave-active-class="transition ease-in duration-25"
      leave-from-class="opacity-100 scale-100"
      leave-to-class="opacity-0 scale-75"
    >
      <div
        v-if="justSaved"
        class="hidden xl:block text-gray-400 text-sm absolute bottom-full translate-y-1 right-0 origin-bottom-right"
      >
        Saved
      </div>
    </transition>
  </div>
</template>

<script lang="ts">
import { defineComponent, defineAsyncComponent } from 'vue'
import type { PropType } from 'vue'
import { gql } from '@urql/vue'

import ShareButton, { CANT_SUBMIT_REASON, SHARE_BUTTON } from './WorkspaceShareButton.vue'
import OnboardingStep from '../components/onboarding/OnboardingStep.vue'
import { AUTHOR } from '../components/shared/AvatarHelper'

import type { WorkspaceDescription_WorkspaceFragment } from './__generated__/WorkspaceDescription'

import { useUpdateWorkspace } from '../mutations/useUpdateWorkspace'

export const WORKSPACE_FRAGMENT = gql`
  fragment WorkspaceDescription_Workspace on Workspace {
    id
    draftDescription
    ...ShareButton
    codebase {
      id
      members {
        ...Author
      }
    }
  }
  ${SHARE_BUTTON}
  ${AUTHOR}
`

type Author = WorkspaceDescription_WorkspaceFragment['codebase']['members'][number]

export default defineComponent({
  props: {
    user: {
      type: Object as PropType<Author>,
    },
    workspace: {
      type: Object as PropType<WorkspaceDescription_WorkspaceFragment>,
      required: true,
    },
    diffIds: {
      type: Array as PropType<string[]>,
      required: true,
    },
    selectedHunkIds: {
      type: Object as PropType<Set<string>>,
      required: true,
    },
  },
  components: {
    ShareButton,
    OnboardingStep,
    Editor: defineAsyncComponent(() => import('../components/workspace/Editor.vue')),
  },
  setup() {
    const { updateWorkspace, mutating: saving } = useUpdateWorkspace()
    return {
      saving,
      updateWorkspace,
    }
  },
  data() {
    return {
      justSaved: false,
      justSavedTimeout: null as null | ReturnType<typeof setTimeout>,

      description: this.workspace.draftDescription,
      descriptionTimeout: null as null | ReturnType<typeof setTimeout>,
    }
  },
  computed: {
    isAuthenticated() {
      return !!this.user
    },
    isAuthorized() {
      const isMember = this.workspace.codebase.members.some(({ id }) => id === this.user?.id)
      return this.isAuthenticated && isMember
    },
    cantSubmitChangeReason() {
      if (this.workspace == null) {
        return CANT_SUBMIT_REASON.WORKSPACE_NOT_FOUND
      }
      if (this.diffIds.length === 0) {
        return CANT_SUBMIT_REASON.NO_DIFFS
      }
      // Have to have a change description before sharing
      if (this.workspace.draftDescription.length === 0) {
        return CANT_SUBMIT_REASON.EMPTY_DESCRIPTION
      }
      // Disallow users from sharing when they have selected hunks
      // (since it might lead them to think they're doing a partial share)
      if (this.selectedHunkIds.size > 0) {
        return CANT_SUBMIT_REASON.HAVE_SELECTED_HUNKS
      }
      return null
    },
    canSubmitChange() {
      return !this.cantSubmitChangeReason
    },
  },
  watch: {
    'workspace.draftDescription': function (newDescription) {
      this.description = newDescription
    },
  },
  methods: {
    clearJustSaved() {
      this.justSaved = false
      if (this.justSavedTimeout) clearTimeout(this.justSavedTimeout)
    },
    setJustSaved() {
      this.justSaved = true
      this.justSavedTimeout = setTimeout(() => {
        this.clearJustSaved()
      }, 1000)
    },
    onUpdatedDescription(event: {
      content: string
      shouldSaveImmediately: boolean
      isInteractiveUpdate: boolean
    }) {
      this.clearJustSaved()
      if (event.shouldSaveImmediately) {
        this.saveDraftDescription(event.content)
      } else {
        this.scheduleSaveDraftDescription(event.content)
      }
    },
    scheduleSaveDraftDescription(draftDescription: string) {
      if (this.descriptionTimeout) clearTimeout(this.descriptionTimeout)
      this.descriptionTimeout = setTimeout(() => this.saveDraftDescription(draftDescription), 300)
    },
    saveDraftDescription(draftDescription: string) {
      if (this.workspace.draftDescription === draftDescription) return
      this.updateWorkspace({ id: this.workspace.id, draftDescription }).then(this.setJustSaved)
    },
  },
})
</script>
