<template>
  <div id="differ-root" class="w-full z-10 relative ph-no-capture">
    <TooManyFilesChanged v-if="hideDiffsTooMany" :count="diffs.length" />
    <div v-else class="overflow-x-auto flex-grow flex flex-col gap-4">
      <DifferFile
        v-for="(fileDiffs, fileKey) in diffsByFile"
        :key="fileKey"
        :is-suggesting="isSuggesting"
        :file-key="fileKey"
        :diffs="fileDiffs"
        :extra-classes="extraClasses"
        :comments="commentsByFile[fileDiffs.newName] ?? []"
        :suggestions="suggestionsForFile(fileDiffs)"
        :new-comment-avatar-url="newCommentAvatarUrl"
        :can-comment="canComment"
        :init-show-suggestions-by-user="initShowSuggestionsByUser"
        :hovering-comment-i-d="hoveringCommentID"
        :show-full-file-button="showFullFileButton"
        :search-result="searchResult"
        :search-current-selected-id="searchCurrentSelectedId"
        :user="user"
        :members="members"
        :workspace="workspace"
        :view="view"
        :change="change"
        :comments-state="commentsState"
        :show-add-button="showAddButton"
        :differ-state="getDifferState(fileKey)"
        @fileSelectedHunks="updateSelectedHunks"
        @applyHunkedSuggestion="(e) => $emit('applyHunkedSuggestion', e)"
        @dismissHunkedSuggestion="(e) => $emit('dismissHunkedSuggestion', e)"
        @set-comment-expanded="onSetCommentExpanded"
        @set-comment-composing-reply="onSetCommentComposingReply"
        @set-is-hidden="onSetFileIsHidden"
      />
    </div>
  </div>
</template>

<script lang="ts">
import { defineAsyncComponent, defineComponent, PropType } from 'vue'
import { MemberFragment, UserFragment } from '../shared/__generated__/TextareaMentions'
import {
  ReviewNewCommentChangeFragment,
  ReviewNewCommentViewFragment,
  ReviewNewCommentWorkspaceFragment,
} from './__generated__/ReviewNewComment'
import {
  CommentState,
  SetCommentComposingReply,
  SetCommentExpandedEvent,
} from '../comments/CommentState'
import TooManyFilesChanged from './TooManyFilesChanged.vue'
import { DifferState, SetFileIsHiddenEvent } from './DifferState'
import { gql } from '@urql/vue'
import {
  DIFFER_FILE_TOP_COMMENT,
  DIFFER_FILE_FILE_DIFF,
  DIFFER_FILE_SUGGESTION,
} from './DifferFile.vue'
import type {
  Differ_TopCommentFragment,
  Differ_FileDiffFragment,
  Differ_SuggestionFragment,
} from './__generated__/Differ'

export const DIFFER_TOP_COMMENT_FRAGMENT = gql`
  fragment Differ_TopComment on TopComment {
    id
    ...DifferFile_TopComment
  }
  ${DIFFER_FILE_TOP_COMMENT}
`

export const DIFFER_FILE_DIFF = gql`
  fragment Differ_FileDiff on FileDiff {
    id
    ...DifferFile_FileDiff
  }
  ${DIFFER_FILE_FILE_DIFF}
`

export const DIFFER_SUGGESTION = gql`
  fragment Differ_Suggestion on Suggestion {
    id
    ...DifferFile_Suggestion
  }
  ${DIFFER_FILE_SUGGESTION}
`

const getIndicesOf = function (searchStr: string, str: string, caseSensitive: boolean): number[] {
  let searchStrLen = searchStr.length
  if (searchStrLen === 0) {
    return []
  }
  let startIndex = 0,
    index,
    indices = []
  if (!caseSensitive) {
    str = str.toLowerCase()
    searchStr = searchStr.toLowerCase()
  }
  while ((index = str.indexOf(searchStr, startIndex)) > -1) {
    indices.push(index)
    startIndex = index + searchStrLen
  }
  return indices
}

export default defineComponent({
  components: {
    DifferFile: defineAsyncComponent(() => import('./DifferFile.vue')),
    TooManyFilesChanged,
  },
  props: {
    isSuggesting: Boolean,
    diffs: {
      type: Array as PropType<Differ_FileDiffFragment[]>,
      required: true,
    },
    comments: {
      type: Array as PropType<Differ_TopCommentFragment[]>,
      required: true,
    },
    extraClasses: String,
    newCommentAvatarUrl: String,
    canComment: Boolean,

    workspace: {
      type: Object as PropType<ReviewNewCommentWorkspaceFragment>,
      required: false,
    },
    view: {
      type: Object as PropType<ReviewNewCommentViewFragment>,
      required: false,
    },
    change: {
      type: Object as PropType<ReviewNewCommentChangeFragment>,
    },
    suggestionsByFile: {
      type: Object as PropType<{ [key: string]: Differ_SuggestionFragment[] }>,
      required: false,
    },
    initShowSuggestionsByUser: {
      type: String,
      required: false,
    },
    // The logged in user
    user: {
      type: Object as PropType<UserFragment>,
      required: false,
    },
    // members of the selected codebase
    members: {
      type: Array as PropType<MemberFragment[]>,
      required: true,
    },
    showFullFileButton: Boolean,

    showAddButton: {
      type: Boolean,
      required: true,
    },
  },
  emits: ['selectedHunks', 'applyHunkedSuggestion', 'dismissHunkedSuggestion'],
  data() {
    return {
      selectedHunksIDsByFile: new Map(),
      hoveringCommentID: null,
      showComments: false,

      searchQuery: null as string | null,
      searchCurrentSelectedId: null as string | null,

      // Client side state of comments is stored here
      // Which comments that are expanded, and which ones that we're currently replying to (and with their contents)
      // This makes commenting resume-able if a child component is re-mounted
      commentsState: new Map<string, CommentState>(),

      differState: new Map<string, DifferState>(),
    }
  },
  computed: {
    hideDiffsTooMany() {
      return this.diffs.length > 250
    },
    searchResult() {
      let result = new Map<string, number[]>()
      if (!this.searchQuery || !this.diffs) {
        this.emitter.emit('search-result', { matchesCount: 0 })
        return result
      }

      let matchesCount = 0

      for (const diff of this.diffs) {
        for (const hunk of diff.hunks) {
          const idx = getIndicesOf(this.searchQuery, hunk.patch, false)
          if (idx.length > 0) {
            result.set(hunk.id, idx)
            matchesCount += idx.length
          }
        }
      }

      this.emitter.emit('search-result', { matchesCount })
      return result
    },
    commentsByFile() {
      const result = {} as Record<string, Differ_TopCommentFragment[]>
      if (!this.comments) {
        return result
      }
      return this.comments.reduce((acc, comment) => {
        if (!comment.codeContext) return acc
        if (!acc[comment.codeContext.path]) acc[comment.codeContext.path] = []
        acc[comment.codeContext.path].push(comment)
        return acc
      }, result)
    },
    diffsByFile() {
      return this.diffs.reduce((acc, diff) => {
        const fileKey = `${diff.origName}//${diff.newName}`
        acc[fileKey] = diff
        return acc
      }, {} as Record<string, Differ_FileDiffFragment>)
    },
  },
  mounted() {
    this.emitter.on('differ-select-all-hunks', this.onSelectAllHunks)
    this.emitter.on('search', this.onSearch)
  },
  unmounted() {
    this.emitter.off('differ-select-all-hunks', this.onSelectAllHunks)
    this.emitter.off('search', this.onSearch)
  },
  methods: {
    updateSelectedHunks(event: { fileKey: string; patchIDs: string[] }) {
      this.selectedHunksIDsByFile.set(event.fileKey, event.patchIDs)

      // Forward upwards!
      // Build combined set
      let combinedIDs = new Set()
      this.selectedHunksIDsByFile.forEach((fileSet) => {
        for (const patch of fileSet) {
          combinedIDs.add(patch)
        }
      })

      this.emitter.emit('differ-selected-hunk-ids', combinedIDs)
      this.$emit('selectedHunks', { patchIDs: combinedIDs })
    },

    onSelectAllHunks() {
      this.emitter.emit('differ-set-hunks-with-prefix', {
        prefix: null,
        selected: true,
      })
    },
    suggestionsForFile(diff: Differ_FileDiffFragment) {
      if (!this.suggestionsByFile) return []
      return this.suggestionsByFile[diff.preferredName] ?? []
    },
    onSearch(event: { searchQuery: string; searchCurrentSelectedId: string }) {
      this.searchQuery = event.searchQuery
      this.searchCurrentSelectedId = event.searchCurrentSelectedId
    },
    onSetCommentExpanded(e: SetCommentExpandedEvent) {
      const current = this.commentsState.get(e.commentId)
      if (current) {
        current.isExpanded = e.isExpanded
        return
      }
      const state = {
        isExpanded: e.isExpanded,
        composingReply: undefined,
      }
      this.commentsState.set(e.commentId, state)
    },
    onSetCommentComposingReply(e: SetCommentComposingReply) {
      const current = this.commentsState.get(e.commentId)
      if (current) {
        current.composingReply = e.composingReply
        return
      }

      const state = {
        isExpanded: true,
        composingReply: e.composingReply,
      }
      this.commentsState.set(e.commentId, state)
    },
    onSetFileIsHidden(e: SetFileIsHiddenEvent) {
      const current = this.differState.get(e.fileKey)
      if (current) {
        current.isHidden = e.isHidden
        return
      }
      const state = {
        isHidden: e.isHidden,
      }
      this.differState.set(e.fileKey, state)
    },
    getDifferState(fileKey: string): DifferState {
      const current = this.differState.get(fileKey)
      if (current) {
        return current
      }
      return {
        isHidden: false,
      }
    },
  },
})
</script>
