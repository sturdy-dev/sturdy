<template>
  <div
    class="d2h-file-wrapper bg-white rounded-md border border-gray-200 z-0 relative overflow-y-hidden overflow-x-auto"
    :class="extraClasses"
    :style="[fileDropdownOpen ? 'min-height: 180px' : '']"
    @mouseleave="hideMakeNewCommentPill"
  >
    <DiffHeader
      :diffs="diffs"
      :is-suggesting="isSuggesting"
      :suggestions="suggestions"
      :show-suggestions="showSuggestions"
      :is-added="isAdded"
      :have-live-changes="haveLiveChanges"
      :showing-suggestions-by-user="showingSuggestionsByUser"
      :can-ignore-file="!diffs.new_name || !diffs.is_new"
      :can-take-suggestions="canTakeSuggestions"
      :show-full-file-button="showFullFileButton"
      :show-add-button="showAddButton"
      @add="toggleAdd"
      @hide="toggleHideFile"
      @undo="undoFile"
      @ignore="ignoreFile"
      @showdropdown="fileDropdownOpen = true"
      @hidedropdown="fileDropdownOpen = false"
      @unhide="emitIsHidden(false)"
      @showSuggestionsByUser="onSuggestionsAvatarClick"
    />
    <div v-if="differState.isHidden && isHiddenTooManyChanges">
      <div class="bg-white">
        <div class="px-4 py-5 sm:px-6">
          <h3 class="text-lg leading-6 font-medium text-gray-900">
            This file has a lot of changes...
          </h3>
          <p class="mt-1 max-w-2xl text-sm text-gray-500">
            It's been hidden by default to save you some power...
            <a href="#" class="text-blue-500" @click.stop.prevent="forceShow">show now!</a>
          </p>
        </div>
      </div>
    </div>

    <template v-if="!differState.isHidden && showSuggestions && !diffs.is_large && !diffs.isLarge">
      <template v-for="suggestion in suggestions">
        <template v-if="suggestion.author.id === showingSuggestionsByUser">
          <DiffTable
            v-for="(hunk, hunkIndex) in suggestion.diff.hunks"
            :key="hunkIndex"
            :unparsed-diff="hunk"
            :grayed-out="hunk.isApplied || hunk.isOutdated || hunk.isDismissed"
          >
            <template #blockIndexAction>
              <div class="relative flex items-start justify-center w-full">
                <span v-if="hunk.isApplied" class="text-sm font-medium text-green-600">Taken</span>
                <span v-else-if="hunk.isDismissed" class="text-sm font-medium text-red-600"
                  >Dismissed</span
                >
                <span v-else-if="hunk.isOutdated" class="text-sm font-medium text-gray-500"
                  >Outdated</span
                >
                <template v-else-if="showSuggestions && canTakeSuggestions">
                  <Button
                    size="tiny"
                    color="green"
                    :grouped="true"
                    :first="true"
                    @click="onClickApplyHunkedSuggestion(hunk, suggestion)"
                  >
                    <CheckIcon class="w-4 h-4" />
                  </Button>

                  <Button
                    size="tiny"
                    color="red"
                    :grouped="true"
                    :last="true"
                    @click="onClickDismissHunkedSuggestion(hunk, suggestion)"
                  >
                    <XIcon class="w-4 h-4" />
                  </Button>
                </template>
              </div>
            </template>
          </DiffTable>
        </template>
      </template>
    </template>

    <div
      v-if="
        isReadyToDisplay &&
        !differState.isHidden &&
        !showSuggestions &&
        !diffs.is_large &&
        !diffs.isLarge
      "
    >
      <div class="d2h-code-wrapper">
        <table
          class="d2h-diff-table leading-4 z-0"
          style="border-collapse: separate; border-spacing: 0"
        >
          <tbody
            v-for="(hunk, hunkIndex) in parsedHunks"
            :key="hunkIndex"
            :class="[
              'd2h-diff-tbody d2h-file-diff z-0',
              checkedHunks.get(diffs.hunks[hunkIndex].id) ? 'opacity-70' : '',
              differState.isHidden ? 'hidden' : '',
            ]"
          >
            <template
              v-for="(block, blockIndex) in highlightedBlocks(hunk.blocks, hunk.language)"
              :key="block.header"
            >
              <tr class="h-full overflow-hidden z-0">
                <td
                  class="d2h-code-linenumber d2h-info h-full sticky left-0 z-20 bg-white min-w-[80px]"
                >
                  <label
                    v-if="showAddButton"
                    class="ml-2.5 inline-flex items-center gap-1.5 font-sans text-sm font-medium"
                  >
                    <input
                      :id="'add-' + fileKey + '-' + hunkIndex"
                      :checked="checkedHunks.get(diffs.hunks[hunkIndex].id)"
                      :value="hunkIndex"
                      type="checkbox"
                      class="focus:ring-red-500 h-4 w-4 text-red-600 border-gray-300 rounded"
                      @change="updatedHunkSelection"
                    />

                    Select
                  </label>
                </td>
                <td class="bg-blue-50" />
                <td class="d2h-info h-full bg-blue-50 left-0 z-0 w-full">
                  <div class="flex items-center sticky left-0">
                    <div class="d2h-code-line d2h-info text-gray-500">
                      &nbsp;&nbsp;{{ block.header }}
                    </div>
                  </div>
                </td>
              </tr>

              <template v-for="(row, rowIndex) in block.lines" :key="rowIndex">
                <tr
                  :data-row-index="rowIndex"
                  class="z-0"
                  :data-preferred-name="diffs.preferred_name || diffs.preferredName"
                  :data-line-oldnum="row.oldNumber"
                  :data-line-newnum="row.newNumber"
                  @mouseover="setShowCommentPill(hunkIndex, blockIndex, rowIndex)"
                >
                  <td
                    :class="[
                      'd2h-code-linenumber bg-white sticky left-0 z-20',
                      checkedHunks.get(diffs.hunks[hunkIndex].id) && row.type === 'insert'
                        ? 'bg-green-50'
                        : '',
                      checkedHunks.get(diffs.hunks[hunkIndex].id) && row.type === 'delete'
                        ? 'bg-red-50'
                        : '',
                      searchMatchesHunk(diffs.hunks[hunkIndex].id) ? '!bg-yellow-100' : '',
                      row.type === 'insert' ? 'border-r border-l border-green-500' : '',
                      row.type === 'delete' ? 'border-r border-l border-red-500' : '',
                    ]"
                  >
                    <label
                      :for="'add-' + fileKey + '-' + hunkIndex"
                      class="cursor-pointer select-none text-gray-600 flex"
                    >
                      <div class="line-num">{{ row.oldNumber }}</div>
                      <div class="line-num">{{ row.newNumber }}</div>
                    </label>
                  </td>

                  <td
                    :class="[
                      row.type === 'insert' ? 'bg-green-50' : '',
                      row.type === 'delete' ? 'bg-red-50' : '',
                    ]"
                  >
                    <button
                      v-if="canComment && showMakeNewCommentPillAt(hunkIndex, blockIndex, rowIndex)"
                      class="absolute -mt-2 -ml-2 z-40 w-4 h-4 inline-flex items-center rounded-md border border-blue-500 bg-blue-400 text-sm font-medium text-gray-500 hover:bg-blue-500 focus:outline-none focus:ring-0"
                      @click.stop.prevent="composeNewComment(hunkIndex, blockIndex, rowIndex)"
                    >
                      <PlusIcon class="text-white" />
                    </button>
                  </td>

                  <td
                    :id="diffs.hunks[hunkIndex].id + '-' + rowIndex"
                    :class="[
                      'code-row-wrapper relative z-10',
                      row.type === 'insert' ? 'bg-green-50' : '',
                      row.type === 'delete' ? 'bg-red-50' : '',
                      row.newNumber && newRowsWithComments.has(row.newNumber) ? '!bg-blue-100' : '',
                      row.oldNumber && oldRowsWithComments.has(row.oldNumber) ? '!bg-blue-100' : '',
                      hasCommentHighlightOnRow(row.oldNumber, row.newNumber, hoveringCommentID)
                        ? '!bg-blue-200'
                        : '',

                      searchIsCurrentSelected(diffs.hunks[hunkIndex].id, rowIndex)
                        ? '!bg-yellow-400  font-bold sturdy-searchmatch'
                        : hasMatchingSearchOnRow(diffs.hunks[hunkIndex].id, rowIndex)
                        ? '!bg-yellow-200 font-bold sturdy-searchmatch'
                        : '',
                    ]"
                  >
                    <div class="d2h-code-line relative z-0 px-4">
                      <span v-if="row.type === 'context'" class="d2h-code-line-prefix">&nbsp;</span>
                      <span v-else class="d2h-code-line-prefix">{{ row.prefix }}</span>
                      <span
                        v-if="row.content"
                        class="d2h-code-line-ctn whitespace-pre"
                        v-html="row.content"
                      />
                      <span v-else class="d2h-code-line-ctn whitespace-pre">{{
                        row.originalContent
                      }}</span>
                    </div>
                  </td>
                </tr>
                <tr
                  v-for="comment in commentsOnRow(row.oldNumber, row.newNumber)"
                  :key="comment.id"
                >
                  <td class="font-sans p-2" colspan="3">
                    <ReviewComment
                      :comment="comment"
                      :members="members"
                      :user="user"
                      :comment-state="getCommentState(comment.id)"
                      @set-comment-expanded="$emit('set-comment-expanded', $event)"
                      @set-comment-composing-reply="$emit('set-comment-composing-reply', $event)"
                    />
                  </td>
                </tr>
                <tr v-if="isComposingNewCommentAt(hunkIndex, blockIndex, rowIndex)">
                  <td class="font-sans p-2" colspan="3">
                    <ReviewNewComment
                      :members="members"
                      :user="user"
                      :path="diffs.preferred_name || diffs.preferredName"
                      :old-path="diffs.orig_name || diffs.origName"
                      :line-is-new="!!row.newNumber"
                      :line-start="row.oldNumber || row.newNumber"
                      :line-end="row.oldNumber || row.newNumber"
                      :change="change"
                      :view="view"
                      :workspace="workspace"
                      :comments-state="commentsState"
                      @cancel="onCancelComposeNewComment"
                      @submitted="onSubmitNewComment"
                      @set-comment-composing-reply="$emit('set-comment-composing-reply', $event)"
                    />
                  </td>
                </tr>
              </template>
            </template>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { CheckIcon, PlusIcon, XIcon } from '@heroicons/vue/solid'
import { defineComponent, PropType, reactive } from 'vue'
import { Block, DifferSetHunksWithPrefix, HighlightedBlock } from './event'
import DiffTable from './DiffTable.vue'
import Button from '../shared/Button.vue'
import * as Diff2Html from 'diff2html'
import DiffHeader from './DiffHeader.vue'
import '../../highlight/highlight_common_languages'
import highlight from '../../highlight/highlight'
import { MemberFragment, UserFragment } from '../shared/__generated__/TextareaMentions'
import ReviewComment from './ReviewComment.vue'
import ReviewNewComment from './ReviewNewComment.vue'
import {
  ReviewNewCommentChangeFragment,
  ReviewNewCommentViewFragment,
  ReviewNewCommentWorkspaceFragment,
} from './__generated__/ReviewNewComment'
import { CommentState } from '../comments/CommentState'
import { DifferState, SetFileIsHiddenEvent } from './DifferState'

interface Data {
  isHiddenTooManyChanges: boolean
  isReadyToDisplay: boolean
  isAdded: boolean
  checkedHunks: Map<string, boolean>

  fileDropdownOpen: boolean // TODO: Set and update

  showMakeNewCommentPillPos: Pos | null
  newCommentComposePos: Pos | null

  showingSuggestionsByUser: string | null

  parsedHunks: any

  enableSyntaxHighlight: boolean
}

interface Pos {
  hunkIndex: number
  blockIndex: number
  rowIndex: number
}

export interface FileDiff {
  // unidiff.FileDiff
  orig_name: string
  new_name: string
  hunks: Array<Hunk>
}

export interface Hunk {
  id: string
  patch: string
}

export default defineComponent({
  components: {
    ReviewNewComment,
    ReviewComment,
    DiffHeader,
    DiffTable,
    PlusIcon,
    CheckIcon,
    XIcon,
    Button,
  },
  emits: [
    'fileSelectedHunks',
    'applyHunkedSuggestion',
    'dismissHunkedSuggestion',
    'set-comment-expanded',
    'set-comment-composing-reply',
    'set-is-hidden',
  ],
  data(): Data {
    return {
      isHiddenTooManyChanges: false,
      isReadyToDisplay: false,
      isAdded: false,
      checkedHunks: reactive(new Map<string, boolean>()),
      fileDropdownOpen: false,

      newCommentComposePos: null,
      showMakeNewCommentPillPos: null,

      showingSuggestionsByUser: null,

      parsedHunks: [], // Diff2Html objects

      enableSyntaxHighlight: false,
    }
  },
  props: {
    isSuggesting: Boolean,
    fileKey: {
      type: String,
      required: true,
    },
    extraClasses: String,
    diffs: {
      type: Object as PropType<FileDiff>,
      required: true,
    },
    comments: Object,
    newCommentAvatarUrl: String,
    canComment: Boolean,
    suggestions: Object,
    initShowSuggestionsByUser: String,
    hoveringCommentID: String,
    showFullFileButton: Boolean,

    searchResult: {
      type: Object as PropType<Map<string, number[]>>,
    },
    searchCurrentSelectedId: String,

    // The logged in user
    user: {
      type: Object as PropType<UserFragment>,
    },
    // members of the selected codebase
    members: {
      type: Array as PropType<MemberFragment[]>,
      required: true,
    },

    workspace: {
      type: Object as PropType<ReviewNewCommentWorkspaceFragment>,
    },
    view: {
      type: Object as PropType<ReviewNewCommentViewFragment>,
    },
    change: {
      type: Object as PropType<ReviewNewCommentChangeFragment>,
    },

    commentsState: {
      type: Object as PropType<Map<string, CommentState>>,
      required: true,
    },

    showAddButton: {
      type: Boolean,
      required: true,
    },

    differState: {
      type: Object as PropType<DifferState>,
      required: true,
    },
  },
  computed: {
    haveLiveChanges(): boolean {
      return this.diffs.hunks.length > 0
    },
    suggestionByUser() {
      // if have suggestions by this user
      if (this.showingSuggestionsByUser && this.suggestions) {
        for (const suggestion of this.suggestions) {
          if (suggestion.author.id === this.showingSuggestionsByUser) {
            return suggestion
          }
        }
      }
      return null
    },
    showSuggestions(): boolean {
      return !!this.suggestionByUser
    },
    canTakeSuggestions(): boolean {
      const s = this.suggestionByUser
      if (!s) {
        return false
      }

      // If at least one hunk is both non-outdated and non-applied, we can take suggestions
      for (const hunk of s.diff.hunks) {
        if (!hunk.isApplied && !hunk.isOutdated && !hunk.isDismissed) {
          return true
        }
      }

      return false
    },
    newRowsWithComments() {
      let res = new Set()
      if (this.comments) {
        for (const comment of this.comments) {
          if (comment.codeContext.lineIsNew && comment.codeContext.lineStart > 0) {
            res.add(comment.codeContext.lineStart)
          }
        }
      }
      return res
    },
    oldRowsWithComments() {
      let res = new Set()
      if (this.comments) {
        for (const comment of this.comments) {
          if (!comment.codeContext.lineIsNew && comment.lineStart > 0) {
            res.add(comment.codeContext.lineStart)
          }
        }
      }
      return res
    },
    rowsWithSearchMatches() {
      let res = new Set<string>()

      if (!this.searchResult || this.searchResult.size === 0) {
        return res
      }

      let endsAt = new Map<string, number[]>()

      for (let hunk of this.diffs.hunks) {
        let ends = 0
        let started = false

        // TODO: Read line by line instead of allocating a list
        let lines = hunk.patch.split('\n')

        for (let line of lines) {
          ends += line.length + 1 // add trimmed newline
          if (!started && line.startsWith('@@ ')) {
            started = true
          }
          if (!started) {
            continue
          }
          if (!endsAt.has(hunk.id)) {
            endsAt.set(hunk.id, [ends])
          } else {
            let ref = endsAt.get(hunk.id)
            if (ref) {
              ref.push(ends)
            }
          }
        }
      }

      for (const [hunkID, foundIndexes] of this.searchResult) {
        let ends = endsAt.get(hunkID)
        if (!ends) {
          continue
        }
        for (let foundIndex of foundIndexes) {
          for (let i = 0; i < ends.length; i++) {
            if (foundIndex >= ends[i] && foundIndex < ends[i + 1]) {
              res.add(hunkID + '-' + i)
            }
          }
        }
      }

      return res
    },
  },
  watch: {
    diffs() {
      this.loadDiffsAndParse()
    },
    fileKey() {
      this.reset()
    },
  },
  created() {
    if (this.initShowSuggestionsByUser) {
      this.showingSuggestionsByUser = this.initShowSuggestionsByUser
    }

    this.loadDiffsAndParse()
    this.emitter.on('differ-deselect-all-hunks', this.reset)
    this.emitter.on('differ-set-hunks-with-prefix', this.onDifferSetHunksWithPrefix)
    this.emitter.on('show-suggestions-by-user', this.showSuggestionsByUser)
  },
  unmounted() {
    // Make sure to send deregister events upwards
    this.reset()
    this.emitter.off('differ-deselect-all-hunks', this.reset)
    this.emitter.off('differ-set-hunks-with-prefix', this.onDifferSetHunksWithPrefix)
    this.emitter.off('show-suggestions-by-user', this.showSuggestionsByUser)
  },
  methods: {
    loadDiffsAndParse() {
      let sumPatchLength = this.diffs.hunks
        .map((hunk) => hunk.patch.length)
        .reduce((a, b) => a + b, 0)

      // Is 20k a reasonable limit?
      // There is definitely slowness happening at ~80k (with syntax highlighting)
      // 10k is approx 125 LOC

      // Syntax highlighting is enabled for files where the diff is less than 10k, and can _not_ be enabled by the user
      this.enableSyntaxHighlight = sumPatchLength < 10000

      // Files that are larger than 20k are hidden by default, but can be enabled
      if (sumPatchLength > 20000) {
        this.isHiddenTooManyChanges = true
        this.isReadyToDisplay = true
        this.emitIsHidden(true)
        this.parsedHunks = []
        return
      }

      // The diff is ok, parse and render it!
      this.parse()
    },

    parse() {
      let res = []

      this.diffs.hunks.forEach((hunk) => {
        let parsed = Diff2Html.parse(hunk.patch, {
          matching: 'lines',
          outputFormat: 'line-by-line',
        })
        res = res.concat(parsed)
      })

      this.parsedHunks = res
      this.updatedSelection()

      // Show!
      if (this.differState.isHidden) {
        // this.emitIsHidden(false)
      }
      this.isHiddenTooManyChanges = false
      this.isReadyToDisplay = true
    },
    forceShow() {
      this.emitIsHidden(false)
      this.isHiddenTooManyChanges = false
      this.isReadyToDisplay = false
      this.parse()
    },
    onDifferSetHunksWithPrefix(event: DifferSetHunksWithPrefix) {
      // If prefix is provided, and this file does not match the prefix
      if (event.prefix) {
        if (!this.diffs.preferred_name.startsWith(event.prefix)) {
          return
        }
      }

      this.setAllHunks(event.selected)
      this.isAdded = event.selected

      // Forward upwards
      this.updatedSelection()
    },
    reset() {
      this.isReadyToDisplay = false
      this.checkedHunks.clear()
      this.isAdded = false
      this.loadDiffsAndParse()
      this.updatedSelection()
    },
    setShowCommentPill(hunkIndex: number, blockIndex: number, rowIndex: number) {
      this.showMakeNewCommentPillPos = {
        hunkIndex: hunkIndex,
        blockIndex: blockIndex,
        rowIndex: rowIndex,
      }
    },
    composeNewComment(hunkIndex: number, blockIndex: number, rowIndex: number) {
      this.newCommentComposePos = {
        hunkIndex: hunkIndex,
        blockIndex: blockIndex,
        rowIndex: rowIndex,
      }
    },

    isComposingNewCommentAt(hunkIndex: number, blockIndex: number, rowIndex: number): boolean {
      if (!this.newCommentComposePos) return false
      if (
        this.newCommentComposePos.hunkIndex === hunkIndex &&
        this.newCommentComposePos.blockIndex === blockIndex &&
        this.newCommentComposePos.rowIndex === rowIndex
      ) {
        return true
      }
    },
    onCancelComposeNewComment() {
      this.newCommentComposePos = undefined
    },
    onSubmitNewComment() {
      this.newCommentComposePos = undefined
    },
    showMakeNewCommentPillAt(hunkIndex: number, blockIndex: number, rowIndex: number) {
      return (
        this.showMakeNewCommentPillPos &&
        this.showMakeNewCommentPillPos.hunkIndex === hunkIndex &&
        this.showMakeNewCommentPillPos.blockIndex === blockIndex &&
        this.showMakeNewCommentPillPos.rowIndex === rowIndex
      )
    },
    hideMakeNewCommentPill() {
      this.showMakeNewCommentPillPos = null
    },
    commentsOnRow(oldRow: number, newRow: number) {
      let res = []

      if (!this.comments) {
        return res
      }

      for (const comment of this.comments) {
        if (
          newRow !== undefined &&
          comment.codeContext.lineIsNew === true &&
          comment.codeContext.lineStart === newRow
        ) {
          res.push(comment)
        }
        if (
          oldRow !== undefined &&
          comment.codeContext.lineIsNew === false &&
          comment.codeContext.lineStart === oldRow
        ) {
          res.push(comment)
        }
      }

      return res
    },
    hasCommentHighlightOnRow(oldRow: number, newRow: number, id: string): boolean {
      if (!this.comments) {
        return false
      }
      for (const comment of this.comments) {
        if (
          newRow !== undefined &&
          comment.codeContext.lineIsNew === true &&
          comment.codeContext.lineStart === newRow &&
          comment.id === id
        ) {
          return true
        }
        if (
          oldRow !== undefined &&
          comment.codeContext.lineIsNew === false &&
          comment.codeContext.lineStart === oldRow &&
          comment.id === id
        ) {
          return true
        }
      }
      return false
    },
    toggleHideFile() {
      this.emitIsHidden(!this.differState.isHidden)
    },
    toggleAdd() {
      this.isAdded = !this.isAdded

      // Hidden follows added
      this.emitIsHidden(this.isAdded)

      // Set all hunks as added / un-added
      this.setAllHunks(this.isAdded)

      // Forward upwards
      this.updatedSelection()
    },
    setAllHunks(setTo: boolean) {
      this.diffs.hunks.forEach((val) => {
        this.checkedHunks.set(val.id, setTo)
      })
    },
    updatedHunkSelection(event: Event) {
      let el = event.target as HTMLInputElement
      let hunkIndex = parseInt(el.value)

      this.checkedHunks.set(this.diffs.hunks[hunkIndex].id, el.checked)

      // Update isAdded
      let total = this.parsedHunks.length
      let count = 0
      this.checkedHunks.forEach((val) => {
        count += val ? 1 : 0
      })

      this.isAdded = total === count

      // Forward upwards
      this.updatedSelection()
    },
    updatedSelection() {
      let hunkPatchIDs = new Set<string>()

      let currentlyExistingHunkIds = new Set<string>(this.diffs.hunks.map((h) => h.id))

      this.checkedHunks.forEach((isChecked, patchID) => {
        // un-set if hunk no longer exists
        if (!currentlyExistingHunkIds.has(patchID)) {
          console.log('un-registering', { patchID })
          this.checkedHunks.delete(patchID)
          return
        }

        if (isChecked) {
          hunkPatchIDs.add(patchID)
        }
      })

      this.$emit('fileSelectedHunks', {
        fileKey: this.fileKey,
        patchIDs: hunkPatchIDs,
      })
    },
    ignoreFile(fileName: string) {
      this.emitter.emit('ignore-file', { fileName: '/' + fileName })
    },
    undoFile() {
      let ids = new Set<string>(this.diffs.hunks.map((h) => h.id))
      this.emitter.emit('undo-file', { patch_ids: ids })
    },
    sendSetWithPrefix(prefix: string, selected: boolean) {
      // Send event to siblings, etc.
      this.emitter.emit('differ-set-hunks-with-prefix', {
        prefix: prefix,
        selected: selected,
      })
    },
    highlightedBlocks(input: Array<Block>, lang: string): Array<HighlightedBlock> {
      return highlight(input, lang, this.enableSyntaxHighlight)
    },
    onClickDismissHunkedSuggestion(hunk: any, suggestion: any) {
      this.$emit('dismissHunkedSuggestion', {
        suggestionId: suggestion.id,
        hunks: [hunk.id],
      })
    },
    onClickApplyHunkedSuggestion(hunk: any, suggestion: any) {
      this.$emit('applyHunkedSuggestion', {
        suggestionId: suggestion.id,
        hunks: [hunk.id],
      })
    },
    onSuggestionsAvatarClick(userID: string) {
      if (this.showingSuggestionsByUser == userID) {
        this.showSuggestionsByUser(null) // hide
      } else {
        this.showSuggestionsByUser(userID) // show
      }
    },
    showSuggestionsByUser(userID: string | null) {
      // Check if this component have any suggestions for this user
      if (!this.suggestions) {
        this.showingSuggestionsByUser = null
        return
      }

      if (userID) {
        let haveSuggestion = false
        for (const suggestion of this.suggestions) {
          if (suggestion.author.id == userID) {
            haveSuggestion = true
            break
          }
        }
        if (!haveSuggestion) {
          this.showingSuggestionsByUser = null
          return
        }
      }

      this.showingSuggestionsByUser = userID
    },
    searchMatchesHunk(hunkID: string): boolean {
      if (!this.searchResult) {
        return false
      }
      let r = this.searchResult.get(hunkID)
      return r && r.length > 0
    },
    hasMatchingSearchOnRow(hunkID: string, rowIndex: number): boolean {
      let matchingRows = this.rowsWithSearchMatches
      return matchingRows.has(hunkID + '-' + rowIndex)
    },
    searchIsCurrentSelected(hunkID: string, rowIndex: number): boolean {
      return this.searchCurrentSelectedId === hunkID + '-' + rowIndex
    },

    getCommentState(id: string): CommentState {
      const state = this.commentsState.get(id)
      if (state) {
        return state
      }

      return {
        isExpanded: false,
        composingReply: undefined,
      }
    },

    emitIsHidden(val: boolean) {
      const ev: SetFileIsHiddenEvent = {
        fileKey: this.fileKey,
        isHidden: val,
      }
      this.$emit('set-is-hidden', ev)
    },
  },
})
</script>
