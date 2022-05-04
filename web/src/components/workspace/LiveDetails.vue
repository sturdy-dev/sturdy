<template>
  <div v-if="data">
    <Banner
      v-if="isAuthorized && conflictsData && conflictsData.workspace.conflicts"
      status="warning"
      class="my-2"
    >
      <h3 class="text-sm font-medium text-yellow-800">
        The changes are conflicting with changes already landed on the codebase.
      </h3>
      <div class="mt-2 text-sm text-yellow-700">
        <p>You need to sync and fix the conflicts before the changes can be landed.</p>
      </div>
    </Banner>

    <Banner
      v-if="hasHiddenChanges"
      class="mb-2"
      status="warning"
      message="This draft change has more changes, but you don't have access to see them."
    />

    <div v-if="hasLiveChanges && mutable" class="relative ml-3 flex space-x-4">
      <label class="inline-flex items-center gap-1.5 text-sm font-medium">
        <input
          type="checkbox"
          class="focus:ring-red-500 h-4 w-4 text-red-600 border-gray-300 rounded"
          :checked="hunkCount === selectedHunkIDs.size"
          @change.prevent="$event.target.checked ? selectAll() : deselectAll()"
        />

        Select All
      </label>

      <Spinner v-if="isStale" />

      <div class="grow" />
      <WorkspaceUndoRedo class="self-end" :workspace="workspace" />
    </div>

    <template v-if="!isSuggesting">
      <Banner
        v-for="(files, suggestingUserID) in suggestedFilesByUser"
        :key="suggestingUserID"
        status="success"
        :show-icon="false"
        class="my-2"
      >
        <div class="w-full flex items-center">
          <div class="flex items-center flex-grow overflow-hidden text-ellipsis">
            <Avatar class="mr-2" size="6" :author="suggestingUsers[suggestingUserID]" />
            <span class="flex-none"
              >{{ suggestingUsers[suggestingUserID].name }} has suggested changes to&nbsp;</span
            >
            <span class="whitespace-nowrap">
              {{ Array.from(files).join(',&nbsp;') }}
            </span>
          </div>

          <Button
            v-if="lastShowSuggestionsByUser === suggestingUserID"
            class="whitespace-nowrap"
            size="small"
            @click="onClickShowSuggestionsByUser(null)"
          >
            Hide suggestions
          </Button>
          <Button
            v-else
            class="whitespace-nowrap"
            size="small"
            @click="onClickShowSuggestionsByUser(suggestingUserID)"
          >
            Show suggestions
          </Button>

          <div class="ml-4 flex-shrink-0 flex" title="Dismiss suggestion">
            <button
              class="rounded-md inline-flex text-gray-400 hover:text-gray-500 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
              @click="dismissSuggestionByUser(suggestingUserID)"
            >
              <span class="sr-only">Dismiss</span>
              <XIcon class="h-5 w-5" aria-hidden="true" />
            </button>
          </div>
        </div>
      </Banner>
    </template>

    <div v-if="mutable && !userIsConnectedToView && view" class="rounded-md p-4 bg-blue-50 my-4">
      <div class="flex">
        <div class="flex-shrink-0">
          <svg
            class="h-5 w-5 text-blue-400"
            xmlns="http://www.w3.org/2000/svg"
            viewBox="0 0 20 20"
            fill="currentColor"
            aria-hidden="true"
          >
            <path
              fill-rule="evenodd"
              d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a1 1 0 000 2v3a1 1 0 001 1h1a1 1 0 100-2v-3a1 1 0 00-1-1H9z"
              clip-rule="evenodd"
            />
          </svg>
        </div>
        <div class="ml-3">
          <p class="text-sm font-medium text-blue-800">
            The computer with the
            <code>{{ view.shortMountPath }}</code> directory isn't running Sturdy right now.
          </p>
        </div>
      </div>
    </div>

    <div class="relative">
      <div v-if="isFetching" class="space-y-8 mt-4">
        <div
          v-for="k in [1, 2, 3]"
          :key="k"
          class="border border-gray-200 shadow rounded-md p-4 max-w-lg w-full"
        >
          <div class="animate-pulse flex space-x-4">
            <div class="flex-1 space-y-1 py-1">
              <div class="h-4 bg-gray-200 rounded w-3/4"></div>
              <div class="h-4 bg-gray-200 rounded"></div>
              <div class="h-4 bg-gray-200 rounded w-5/6"></div>
              <div class="h-4 bg-gray-200 rounded w-3/4"></div>
              <div class="h-4 bg-gray-200 rounded"></div>
              <div class="h-4 bg-gray-200 rounded w-5/6"></div>
              <div class="h-4 bg-gray-200 rounded w-3/4"></div>
              <div class="h-4 bg-gray-200 rounded"></div>
            </div>
          </div>
        </div>
      </div>

      <template v-else-if="!hasLiveChanges">
        <div
          v-if="!hasHiddenChanges && !hasLiveChanges && workspace?.author?.id === user?.id"
          class="mt-8"
        >
          <NoChangesOwnWorkspace :workspace="workspace" />
        </div>
        <div
          v-else-if="!hasHiddenChanges && !hasLiveChanges && workspace?.author?.id !== user?.id"
          class="mt-8"
        >
          <NoChangesOthersWorkspace :workspace="workspace" />
        </div>

        <div class="flex items-center flex-col mt-8">
          <Spinner v-if="isStale" />
        </div>
      </template>

      <template v-else>
        <div class="mt-2" />

        <OnboardingStep
          id="WorkspaceChanges"
          :dependencies="['FindingYourWorkspace']"
          :enabled="combinedDiffTypes.length > 0"
        >
          <template #title>Workspace Changes</template>

          <template #description>
            These are the changes that currently reside within this workspace. Until these changes
            have been published to the changelog, you can review them here.
          </template>

          <Differ
            :diffs="combinedDiffTypes"
            :suggestions-by-file="suggestionsByFile"
            :init-show-suggestions-by-user="lastShowSuggestionsByUser"
            :is-suggesting="isSuggesting"
            :can-comment="isAuthorized"
            :comments="comments"
            :user="user"
            :members="members"
            :workspace="workspace"
            :view="view"
            :show-add-button="isAuthorized"
            @selectedHunks="updateSelectedHunks"
            @applyHunkedSuggestion="onApplyHunkedSuggestion"
            @dismissHunkedSuggestion="onDismissHunkedSuggestion"
          />
        </OnboardingStep>
      </template>
    </div>
  </div>
</template>

<script lang="ts">
import Differ, { DIFFER_FILE_DIFF, DIFFER_SUGGESTION } from '../differ/Differ.vue'
import http from '../../http'
import { XIcon } from '@heroicons/vue/outline'
import { Banner } from '../../atoms'
import Avatar from '../../atoms/Avatar.vue'
import Button from '../../atoms/Button.vue'
import { CombinedError, gql, useMutation, useQuery } from '@urql/vue'
import { useRoute } from 'vue-router'
import { computed, defineComponent, inject, ref } from 'vue'
import type { PropType, Ref } from 'vue'
import type { DeepMaybeRef } from '@vueuse/core'
import { useUpdatedWorkspace } from '../../subscriptions/useUpdatedWorkspace'
import { useUpdatedGitHubPullRequest } from '../../subscriptions/useUpdatedGitHubPullRequest'
import OnboardingStep from '../onboarding/OnboardingStep.vue'
import NoChangesOwnWorkspace, { NO_CHANGES_OWN_WORKSPACE } from './NoChangesOwnWorkspace.vue'
import NoChangesOthersWorkspace, {
  NO_CHANGES_OTHERS_WORKSPACE,
} from './NoChangesOthersWorkspace.vue'
import { useApplySuggestionHunks } from '../../mutations/useApplySuggestionHunks'
import { useDismissSuggestionHunks } from '../../mutations/useDismissSuggestionHunks'
import { useUpdatedSuggestion } from '../../subscriptions/useUpdatedSuggestion'
import { useDismissSuggestion } from '../../mutations/useDismissSuggestion'
import { useRemovePatches } from '../../mutations/useRemovePatches'
import { Feature } from '../../__generated__/types'
import type {
  LiveDetailsDiffsFragment,
  LiveDetailsMemberFragment,
  LiveDetailsQuery,
  LiveDetailsQueryVariables,
  LiveDetailsViewFragment,
} from './__generated__/LiveDetails'
import Spinner from '../../atoms/Spinner.vue'
import type { Differ_SuggestionFragment } from '../differ/__generated__/Differ'
import WorkspaceUndoRedo, {
  WORKSPACE_FRAGMENT as WORKSPACE_UNDO_REDO_FRAGMENT,
} from '../../molecules/WorkspaceUndoRedo.vue'
import type { LiveDetailsWorkspaceFragment } from './__generated__/LiveDetails'

export const LIVE_DETAILS_WORKSPACE = gql`
  fragment LiveDetailsWorkspace on Workspace {
    id
    ...NoChangesOwnWorkspace
    ...NoChangesOthersWorkspace
    ...WorkspaceUndoRedo_Workspace
  }
  ${NO_CHANGES_OWN_WORKSPACE}
  ${NO_CHANGES_OTHERS_WORKSPACE}
  ${WORKSPACE_UNDO_REDO_FRAGMENT}
`

export const LIVE_DETAILS_DIFFS = gql`
  fragment LiveDetailsDiffs on FileDiff {
    id

    origName
    newName
    preferredName

    isDeleted
    isNew
    isMoved
    isHidden

    hunks {
      _id
      hunkID
      patch

      isOutdated
      isApplied
      isDismissed
    }

    ...Differ_FileDiff
  }

  ${DIFFER_FILE_DIFF}
`

export const LIVE_DETAILS_VIEW = gql`
  fragment LiveDetailsView on View {
    id
    shortMountPath
    lastUsedAt
  }
`

export const LIVE_DETAILS_MEMBER = gql`
  fragment LiveDetailsMember on Author {
    id
    name
    avatarUrl
  }
`

type Suggestion = LiveDetailsQuery['workspace']['suggestions'][number]
type SuggestingAuthor = LiveDetailsQuery['workspace']['suggestions'][number]['author']
type SuggestionFileDiff = LiveDetailsQuery['workspace']['suggestions'][number]['diffs'][number]

interface Data {
  lastShowSuggestionsByUser: string | null
  selectedHunkIDs: Set<string>
}

export default defineComponent({
  components: {
    Spinner,
    XIcon,
    NoChangesOthersWorkspace,
    NoChangesOwnWorkspace,
    OnboardingStep,
    Banner,
    Differ,
    Avatar,
    Button,
    WorkspaceUndoRedo,
  },
  props: {
    view: {
      type: Object as PropType<LiveDetailsViewFragment>,
      required: false,
      default: null,
    },
    user: {
      type: Object,
      required: false,
      default: null,
    },
    members: {
      type: Array as PropType<Array<LiveDetailsMemberFragment>>,
      required: true,
    },
    isOnAuthoritativeView: {
      type: Boolean,
      required: true,
    },
    workspace: {
      type: Object as PropType<LiveDetailsWorkspaceFragment>,
      required: true,
    },
    mutable: {
      type: Boolean,
      required: true,
    },
    isSuggesting: {
      type: Boolean,
      required: true,
    },

    diffs: {
      type: Object as PropType<Array<LiveDetailsDiffsFragment>>,
      required: true,
    },

    comments: {
      type: Array,
      required: true,
    },

    isStale: {
      type: Boolean,
      required: true,
    },

    isFetching: {
      type: Boolean,
      required: true,
    },
  },
  emits: ['codebase-updated'],
  setup() {
    const features = inject<Ref<Array<Feature>>>('features', ref([]))
    const isGitHubEnabled = computed(() => features?.value?.includes(Feature.GitHub))

    let route = useRoute()
    let workspaceID = computed(() => route.params.id as string)

    let { data, fetching, error, executeQuery } = useQuery<
      LiveDetailsQuery,
      DeepMaybeRef<LiveDetailsQueryVariables>
    >({
      query: gql`
        query LiveDetails($workspaceID: ID!, $isGitHubEnabled: Boolean!) {
          workspace(id: $workspaceID) {
            id
            upToDateWithTrunk
            headChange {
              id
              createdAt
            }
            suggestions {
              id
              dismissedAt
              author {
                id
                name
                avatarUrl
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
                isHidden

                hunks {
                  _id

                  hunkID
                  patch

                  isOutdated
                  isApplied
                  isDismissed
                }
              }

              ...Differ_Suggestion
            }

            codebase {
              id
              gitHubIntegration @include(if: $isGitHubEnabled) {
                id
                owner
                name
                enabled
                gitHubIsSourceOfTruth
                lastPushErrorMessage
                lastPushAt
              }
            }
            gitHubPullRequest @include(if: $isGitHubEnabled) {
              id
              pullRequestNumber
              open
              merged
              mergedAt
            }
          }
        }

        ${DIFFER_SUGGESTION}
      `,
      requestPolicy: 'cache-and-network',
      variables: { workspaceID: workspaceID, isGitHubEnabled },
    })

    useUpdatedGitHubPullRequest(workspaceID)
    useUpdatedSuggestion(workspaceID)

    const { executeMutation: extractWorkspace } = useMutation(gql`
      mutation LiveDetailsExtract($workspaceID: ID!, $patchIDs: [String!]!) {
        extractWorkspace(input: { workspaceID: $workspaceID, patchIDs: $patchIDs }) {
          id
          name
        }
      }
    `)

    // Check if the workspace conflicts with trunk - run this 2s after going to the workspace page (just to delay it a bit?)
    let readyToCheckIfConflicts = ref(false)
    setTimeout(() => {
      readyToCheckIfConflicts.value = true
    }, 2000)
    let {
      data: conflictsData,
      error: conflictsError,
      executeQuery: executeConflictsQuery,
    } = useQuery({
      query: gql`
        query LiveDetailsConflicts($workspaceID: ID!) {
          workspace(id: $workspaceID) {
            id
            conflicts
          }
        }
      `,
      requestPolicy: 'cache-and-network',
      variables: { workspaceID: workspaceID },
      pause: computed(() => !readyToCheckIfConflicts.value),
    })

    // Subscription for if the workspace is conflicting
    useUpdatedWorkspace(workspaceID, {
      pause: computed(() => !readyToCheckIfConflicts.value || !workspaceID.value),
    })

    const applySuggestionHunksResult = useApplySuggestionHunks()
    const dismissSuggestionHunksResult = useDismissSuggestionHunks()
    const dismissSuggestionResult = useDismissSuggestion()
    const removePatchesResult = useRemovePatches()
    return {
      async removePatches(workspaceID: string, hunkIds: string[]) {
        await removePatchesResult({
          workspaceID,
          hunkIDs: hunkIds,
        })
      },

      async dismissSuggestion(suggestionId: string) {
        await dismissSuggestionResult({
          id: suggestionId,
        })
      },

      async applySuggestionHunks(suggestionId: string, hunkIds: string[]) {
        await applySuggestionHunksResult({
          id: suggestionId,
          hunkIDs: hunkIds,
        })
      },

      async dismissSuggestionHunks(suggestionId: string, hunkIds: string[]) {
        await dismissSuggestionHunksResult({
          id: suggestionId,
          hunkIDs: hunkIds,
        })
      },

      async extractWorkspace(workspaceID: string, patchIDs: string[]) {
        const variables = { workspaceID, patchIDs }
        return await extractWorkspace(variables).then((result) => {
          if (result.error) {
            throw new CombinedError(result.error)
          }
          return result.data.extractWorkspace
        })
      },

      fetching: fetching,
      data: data,
      error: error,
      refresh() {
        executeQuery({
          requestPolicy: 'network-only',
        })
      },

      conflictsData,
      conflictsError,
      executeConflictsQuery,
    }
  },
  data(): Data {
    return {
      lastShowSuggestionsByUser: null,
      selectedHunkIDs: new Set<string>(),
    }
  },
  computed: {
    filesWithDiffs(): Set<string> {
      const set = new Set<string>()
      this.visible_diffs.map((d) => d.preferredName).forEach((f) => set.add(f))
      return set
    },
    isAuthenticated() {
      return !!this.user
    },
    isAuthorized() {
      const isMember = this.members.some(({ id }) => id === this.user?.id)
      return this.isAuthenticated && isMember
    },
    visible_diffs(): Array<LiveDetailsDiffsFragment> {
      if (!this.diffs) {
        return []
      }
      return this.diffs.filter((d) => !d.isHidden)
    },
    hasHiddenChanges() {
      if (!this.diffs) {
        return false
      }
      return this.diffs.length > this.visible_diffs.length
    },
    hasLiveChanges() {
      return this.visible_diffs.length > 0
    },
    hunkCount() {
      return this.visible_diffs.reduce((acc, current) => {
        if (current.hunks) {
          return acc + current.hunks.length
        }
        return acc
      }, 0)
    },
    userIsConnectedToView() {
      let t = +new Date() / 1000 - 120 // 2 minutes ago
      if (this.view && this.view.lastUsedAt && this.view.lastUsedAt > t) {
        return true
      }
      return false
    },
    combinedDiffTypes() {
      return this.extraSuggestionOnlyFiles.concat(this.visible_diffs)
    },
    gitHubPRLink() {
      return (
        'https://github.com/' +
        this.data?.workspace?.codebase?.gitHubIntegration?.owner +
        '/' +
        this.data?.workspace?.codebase?.gitHubIntegration?.name +
        '/pull/' +
        this.data?.workspace?.gitHubPullRequest?.pullRequestNumber
      )
    },

    openSuggestions(): Suggestion[] {
      if (!this.data) {
        return []
      }
      return this.data.workspace.suggestions.filter((s) => {
        return !s.dismissedAt
      })
    },
    suggestingUsers() {
      let acc: { [index: string]: SuggestingAuthor } = {}
      return this.openSuggestions.reduce((acc, current) => {
        acc[current.author.id] = current.author
        return acc
      }, acc)
    },
    suggestionsByUser() {
      let acc: { [index: string]: Suggestion[] } = {}
      return this.openSuggestions.reduce((acc, current) => {
        if (!acc[current.author.id]) {
          acc[current.author.id] = []
        }
        acc[current.author.id].push(current)
        return acc
      }, acc)
    },
    suggestedFilesByUser() {
      const suggestedFilesByUser: { [index: string]: Set<string> } = {}
      this.openSuggestions.forEach((suggestion) => {
        suggestion.diffs.forEach((diff) => {
          if (!suggestedFilesByUser[suggestion.author.id]) {
            suggestedFilesByUser[suggestion.author.id] = new Set()
          }
          suggestedFilesByUser[suggestion.author.id].add(diff.preferredName)
        })
      })
      return suggestedFilesByUser
    },
    suggestionsByFile(): { [key: string]: Differ_SuggestionFragment[] } {
      const suggestionsByFile: { [key: string]: Differ_SuggestionFragment[] } = {}
      this.openSuggestions.forEach((suggestion) => {
        suggestion.diffs.forEach((diff) => {
          if (!suggestionsByFile[diff.preferredName]) {
            suggestionsByFile[diff.preferredName] = []
          }

          let sug = { ...suggestion }
          sug.diffs = [diff]
          suggestionsByFile[diff.preferredName].push(sug)
        })
      })
      return suggestionsByFile
    },
    extraSuggestionOnlyFiles(): LiveDetailsDiffsFragment[] {
      // Files that have suggestions to them, but no diffs in the original view
      const extraSuggestionOnlyFiles: LiveDetailsDiffsFragment[] = []
      const userID = this.lastShowSuggestionsByUser

      if (this.suggestionsByUser && userID && this.suggestionsByUser[userID]) {
        this.suggestionsByUser[userID].forEach((suggestion) => {
          // Copy the diff, but remove the hunks
          suggestion.diffs.forEach((diff) => {
            // This file is already visible
            let name = diff.preferredName
            if (this.filesWithDiffs.has(name)) {
              return
            }
            let cloned = { ...diff }
            cloned.hunks = []
            extraSuggestionOnlyFiles.push(cloned)
          })
        })
      }

      return extraSuggestionOnlyFiles
    },
  },
  watch: {
    'data.workspace.upToDateWithTrunk': function (newVal) {
      // If no longer up to date
      if (!newVal) {
        this.executeConflictsQuery()
      }
    },
  },
  mounted() {
    window.addEventListener('keydown', this.onkey)

    this.emitter.on('ignore-file', this.onIgnoreFileEvent)
    this.emitter.on('undo-file', this.onUndoFileEvent)
    this.emitter.on('undo-selected', this.deleteSelected)
    this.emitter.on('copy-selected-to-new-workspace', this.copySelected)
  },
  unmounted() {
    this.emitter.off('ignore-file', this.onIgnoreFileEvent)
    this.emitter.off('undo-file', this.onUndoFileEvent)
    this.emitter.off('undo-selected', this.deleteSelected)
    this.emitter.off('copy-selected-to-new-workspace', this.copySelected)
  },
  methods: {
    copySelected() {
      const selectedHunkIDs = Array.from(this.selectedHunkIDs)
      if (selectedHunkIDs.length === 0) return

      this.extractWorkspace(this.workspace.id, selectedHunkIDs).then((result) => {
        this.emitter.emit('notification', {
          title: 'Changes copied',
          message: `Selected changes are copied to ${result.name}`,
        })
        this.emitter.emit('differ-deselect-all-hunks', {})
      })
    },

    deleteSelected() {
      this.removePatches(this.workspace.id, Array.from(this.selectedHunkIDs)).then(() => {
        this.emitter.emit('differ-deselect-all-hunks', {})
      })
    },
    undoFile(patch_ids: Set<string>) {
      this.removePatches(this.workspace.id, Array.from(patch_ids))
    },
    selectAll() {
      this.emitter.emit('differ-select-all-hunks', {})
    },
    deselectAll() {
      this.emitter.emit('differ-deselect-all-hunks', {})
    },
    ignoreFile(path: string) {
      if (!this.view) {
        return
      }

      fetch(http.url('v3/views/' + this.view.id + '/ignore-file'), {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          path: path,
        }),
        credentials: 'include',
      })
    },
    onkey(e: KeyboardEvent) {
      const cmd = e.metaKey || e.ctrlKey
      const shift = e.shiftKey
      const a = e.keyCode == 65
      const enter = e.keyCode === 13

      switch (true) {
        case cmd && shift && a:
          this.emitter.emit('differ-deselect-all-hunks', {})
          e.stopPropagation()
          e.preventDefault()
          break
        case cmd && shift && enter:
          e.stopPropagation()
          e.preventDefault()
          break
        case cmd && enter:
          e.stopPropagation()
          e.preventDefault()
          break
      }
    },
    updateSelectedHunks(ev: { patchIDs: string[] }) {
      this.selectedHunkIDs = new Set(ev.patchIDs)
    },
    onIgnoreFileEvent(ignoreFile: { fileName: string }) {
      this.ignoreFile(ignoreFile.fileName)
    },
    onUndoFileEvent(undoFile: { patch_ids: Set<string> }) {
      this.undoFile(undoFile.patch_ids)
    },
    onClickShowSuggestionsByUser(userID: string | null) {
      // Event is picked up by downstream components
      this.lastShowSuggestionsByUser = userID
      this.emitter.emit('show-suggestions-by-user', userID)
    },
    onDismissHunkedSuggestion(e: { suggestionId: string; hunks: string[] }) {
      this.dismissSuggestionHunks(e.suggestionId, e.hunks)
    },
    onApplyHunkedSuggestion(e: { suggestionId: string; hunks: string[] }) {
      this.applySuggestionHunks(e.suggestionId, e.hunks)
    },
    dismissSuggestionByUser(userId: string) {
      this.suggestionsByUser[userId].map((s) => s.id).forEach(this.dismissSuggestion)
    },
  },
})
</script>
