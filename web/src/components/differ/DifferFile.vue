<template>
  <div
    :id="diffs.id"
    class="d2h-file-wrapper bg-white rounded-md border border-gray-200 z-0 relative overflow-y-hidden overflow-x-auto"
    :class="[
      extraClasses,
    ]"
    :style="[fileDropdownOpen ? 'min-height: 180px' : '']"
    @mouseleave="hideMakeNewCommentPill"
  >
    <DiffHeader
      :diffs="diffs"
      :class="[searchIsCurrentSelectedFilename(diffs.id)
        ? '!bg-yellow-400 font-bold sturdy-searchmatch'
        : searchMatchesFiles(diffs.id)
        ? '!bg-yellow-200 font-bold sturdy-searchmatch'
        : '']"
      :file-key="fileKey"
      :is-suggesting="isSuggesting"
      :suggestions="suggestions"
      :show-suggestions="showSuggestions"
      :is-added="isAdded"
      :have-live-changes="haveLiveChanges"
      :showing-suggestions-by-user="showingSuggestionsByUser"
      :can-ignore-file="canIgnoreFile"
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

    <template v-if="!differState.isHidden && showSuggestions && !diffs.isLarge">
      <template v-for="suggestion in suggestions">
        <template v-if="suggestion.author.id === showingSuggestionsByUser">
          <template v-for="(diff, diffIdx) in suggestion.diffs" :key="diffIdx">
            <DiffTable
              v-for="hunk in diff.hunks"
              :key="hunk.id"
              :unparsed-diff="hunk"
              :grayed-out="hunk.isApplied || hunk.isOutdated || hunk.isDismissed"
            >
              <template #blockIndexAction>
                <div class="relative flex items-start justify-center w-full">
                  <span v-if="hunk.isApplied" class="text-sm font-medium text-green-600"
                    >Taken</span
                  >
                  <span v-else-if="hunk.isDismissed" class="text-sm font-medium text-red-600">
                    Dismissed
                  </span>
                  <span v-else-if="hunk.isOutdated" class="text-sm font-medium text-gray-500">
                    Outdated
                  </span>
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
    </template>

    <div v-if="isReadyToDisplay && !differState.isHidden && !showSuggestions && !diffs.isLarge">
      <div class="d2h-code-wrapper">
        <table
          class="d2h-diff-table leading-4"
          style="border-collapse: separate; border-spacing: 0"
        >
          <tbody
            v-for="(hunk, hunkIndex) in parsedHunks"
            :key="hunkIndex"
            :class="[
              'd2h-diff-tbody d2h-file-diff',
              checkedHunks.get(diffs.hunks[hunkIndex].id) ? 'opacity-70' : '',
              differState.isHidden ? 'hidden' : '',
            ]"
          >
            <template
              v-for="(block, blockIndex) in highlightedBlocks(hunk.blocks, hunk.language)"
              :key="block.header"
            >
              <tr class="h-full overflow-hidden">
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
                <td class="d2h-info h-full bg-blue-50 left-0 w-full">
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
                  :data-preferred-name="diffs.preferredName"
                  :data-line-oldnum="row.oldNumber"
                  :data-line-newnum="row.newNumber"
                  @mouseover="setShowCommentPill(hunkIndex, blockIndex, rowIndex)"
                >
                  <td
                    class="d2h-code-linenumber bg-white sticky left-0 z-20"
                    :class="[
                      row.type === 'insert' ? 'bg-green-50 border-r border-l border-green-500' : '',
                      row.type === 'delete' ? 'bg-red-50 border-r border-l border-red-500' : '',
                      searchMatchesHunk(diffs.hunks[hunkIndex].id) ? '!bg-yellow-100' : '',
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
                      class="absolute z-40 -mt-2 -ml-2 w-4 h-4 inline-flex items-center rounded-md border border-blue-500 bg-blue-400 text-sm font-medium text-gray-500 hover:bg-blue-500 focus:outline-none focus:ring-0"
                      @click.stop.prevent="composeNewComment(hunkIndex, blockIndex, rowIndex)"
                    >
                      <PlusIcon class="text-white" />
                    </button>
                  </td>

                  <td
                    :id="diffs.hunks[hunkIndex].id + '-' + rowIndex"
                    class="code-row-wrapper relative z-10"
                    :class="[
                      row.type === 'insert' ? 'bg-green-50' : '',
                      row.type === 'delete' ? 'bg-red-50' : '',
                      row.newNumber && newRowsWithComments.has(row.newNumber) ? '!bg-blue-100' : '',
                      row.oldNumber && oldRowsWithComments.has(row.oldNumber) ? '!bg-blue-100' : '',
                      hasCommentHighlightOnRow(row.oldNumber, row.newNumber, hoveringCommentID)
                        ? '!bg-blue-200'
                        : '',

                      searchIsCurrentSelected(diffs.hunks[hunkIndex].id, rowIndex)
                        ? '!bg-yellow-400 font-bold sturdy-searchmatch'
                        : hasMatchingSearchOnRow(diffs.hunks[hunkIndex].id, rowIndex)
                        ? '!bg-yellow-200 font-bold sturdy-searchmatch'
                        : '',
                    ]"
                  >
                    <div class="d2h-code-line relative px-4">
                      <span class="d2h-code-line-prefix">
                        <template v-if="row.type === 'context'">&nbsp;</template>
                        <template v-else>{{ row.prefix }}</template>
                      </span>

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
                  <td
                    class="d2h-code-linenumber d2h-info h-full sticky left-0 z-20 bg-white min-w-[80px]"
                    colspan="2"
                  />
                  <td class="font-sans p-2">
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
                  <td
                    class="d2h-code-linenumber d2h-info h-full sticky left-0 z-20 bg-white min-w-[80px]"
                    colspan="2"
                  />
                  <td class="font-sans p-2">
                    <ReviewNewComment
                      :members="members"
                      :user="user"
                      :path="diffs.preferredName"
                      :old-path="diffs.origName"
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
import DiffTable, { HUNK_FRAGMENT as DIFF_TABLE_HUNK_FRAGMENT } from './DiffTable.vue'
import Button from '../shared/Button.vue'
import * as Diff2Html from 'diff2html'
import DiffHeader, { DIFF_HEADER_FILE_DIFF, DIFF_HEADER_SUGGESTIONS } from './DiffHeader.vue'
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
import { gql } from '@urql/vue'
import {
  DifferFile_TopCommentFragment,
  DifferFile_FileDiffFragment,
  DifferFile_SuggestionFragment,
} from './__generated__/DifferFile'
import { searchMatches } from './DifferHelper'

type SuggestionHunk = DifferFile_SuggestionFragment['diffs'][number]['hunks'][number]

export const DIFFER_FILE_SUGGESTION = gql`
  fragment DifferFile_Suggestion on Suggestion {
    id
    author {
      id
    }
    diffs {
      id
      isLarge
      hunks {
        id
        ...DiffTable_Hunk
        isApplied
        isOutdated
        isDismissed
      }
    }

    ...DiffHeader_Suggestions
  }

  ${DIFF_TABLE_HUNK_FRAGMENT}
  ${DIFF_HEADER_SUGGESTIONS}
`

export const DIFFER_FILE_FILE_DIFF = gql`
  fragment DifferFile_FileDiff on FileDiff {
    id

    origName
    newName
    preferredName

    isDeleted
    isNew
    isMoved
    isLarge

    hunks {
      id
      patch

      isOutdated
      isApplied
      isDismissed
    }
    ...DiffHeader_FileDiff
  }
  ${DIFF_HEADER_FILE_DIFF}
`

export const DIFFER_FILE_TOP_COMMENT = gql`
  fragment DifferFile_TopComment on TopComment {
    id
    message
    codeContext {
      id
      lineStart
      lineEnd
      lineIsNew
      context
      contextStartsAtLine
      path
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
`

type Position = {
  hunkIndex: number
  blockIndex: number
  rowIndex: number
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
  data() {
    return {
      isHiddenTooManyChanges: false,
      isReadyToDisplay: false,
      isAdded: false,
      checkedHunks: reactive(new Map<string, boolean>()),
      fileDropdownOpen: false,

      newCommentComposePos: undefined as Position | undefined,
      showMakeNewCommentPillPos: undefined as Position | undefined,

      showingSuggestionsByUser: null as string | null,

      enableSyntaxHighlight: false,
    }
  },
  props: {
    isSuggesting: {
      type: Boolean,
      required: true,
    },
    fileKey: {
      type: String,
      required: true,
    },
    extraClasses: String,

    diffs: {
      type: Object as PropType<DifferFile_FileDiffFragment>,
      required: true,
    },

    comments: {
      type: Object as PropType<DifferFile_TopCommentFragment[]>,
      required: true,
    },

    newCommentAvatarUrl: String,
    canComment: Boolean,
    suggestions: {
      type: Array as PropType<Array<DifferFile_SuggestionFragment>>,
      required: true,
    },
    initShowSuggestionsByUser: String,
    hoveringCommentID: String,
    showFullFileButton: {
      type: Boolean,
      default: false,
    },

    searchResult: {
      type: Object as PropType<Map<string, number[]>>,
      required: false,
      default: undefined,
    },

    searchCurrentSelectedId: {
      type: String,
      required: false,
      default: '',
    },

    // The logged in user
    user: {
      type: Object as PropType<UserFragment>,
      required: false,
      default: undefined,
    },
    // members of the selected codebase
    members: {
      type: Array as PropType<MemberFragment[]>,
      required: true,
    },

    workspace: {
      type: Object as PropType<ReviewNewCommentWorkspaceFragment>,
      required: false,
      default: undefined,
    },
    view: {
      type: Object as PropType<ReviewNewCommentViewFragment>,
      required: false,
      default: undefined,
    },
    change: {
      type: Object as PropType<ReviewNewCommentChangeFragment>,
      required: false,
      default: undefined,
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
    hunkIds() {
      return new Set([...this.diffs.hunks.flatMap((hunk) => hunk.id)])
    },
    parsedHunks() {
      return this.diffs.hunks.flatMap(({ patch }) =>
        Diff2Html.parse(patch, {
          matching: 'lines',
          outputFormat: 'line-by-line',
        })
      )
    },
    haveLiveChanges(): boolean {
      return this.diffs.hunks.length > 0
    },
    suggestionByUser(): DifferFile_SuggestionFragment | undefined {
      if (!this.showingSuggestionsByUser) return undefined
      return this.suggestions.find(({ author }) => author.id === this.showingSuggestionsByUser)
    },
    showSuggestions(): boolean {
      return !!this.suggestionByUser
    },
    canTakeSuggestions(): boolean {
      if (!this.suggestionByUser) {
        return false
      }

      return this.suggestionByUser?.diffs?.some(({ hunks }) =>
        hunks.some(
          ({ isApplied, isOutdated, isDismissed }) => !isApplied && !isOutdated && !isDismissed
        )
      )
    },
    newRowsWithComments() {
      return new Set([
        ...this.comments
          .filter(({ codeContext }) => codeContext)
          .map(({ codeContext }) => codeContext!)
          .filter(({ lineIsNew, lineStart }) => lineIsNew && lineStart > 0)
          .map(({ lineStart }) => lineStart),
      ])
    },
    oldRowsWithComments() {
      return new Set([
        ...this.comments
          .filter(({ codeContext }) => codeContext)
          .map(({ codeContext }) => codeContext!)
          .filter(({ lineIsNew, lineStart }) => !lineIsNew && lineStart > 0)
          .map(({ lineStart }) => lineStart),
      ])
    },
    rowsWithSearchMatches() {
      return searchMatches(this.searchResult, this.diffs.hunks)
    },
    canIgnoreFile() {
      return this.diffs.isNew && this.diffs.newName && !this.diffs.newName.endsWith('.gitignore')
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
    this.emitter.off('differ-deselect-all-hunks', this.reset)
    this.emitter.off('differ-set-hunks-with-prefix', this.onDifferSetHunksWithPrefix)
    this.emitter.off('show-suggestions-by-user', this.showSuggestionsByUser)
  },
  methods: {
    loadDiffsAndParse() {
      const sumPatchLength = this.diffs.hunks
        .map(({ patch }) => patch.length)
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
        return
      }

      // The diff is ok, parse and render it!
      this.parse()
    },

    parse() {
      this.updatedSelection()
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
        if (!this.diffs.preferredName.startsWith(event.prefix)) {
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
      return (
        this.newCommentComposePos.hunkIndex === hunkIndex &&
        this.newCommentComposePos.blockIndex === blockIndex &&
        this.newCommentComposePos.rowIndex === rowIndex
      )
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
      this.showMakeNewCommentPillPos = undefined
    },
    commentsOnRow(oldRow: number, newRow: number): Array<DifferFile_TopCommentFragment> {
      let res = Array<DifferFile_TopCommentFragment>()

      if (!this.comments) {
        return res
      }

      for (const comment of this.comments) {
        if (!comment.codeContext) {
          continue
        }
        if (
          newRow != undefined &&
          comment.codeContext.lineIsNew &&
          comment.codeContext.lineStart === newRow
        ) {
          res.push(comment)
        }
        if (
          oldRow != undefined &&
          !comment.codeContext.lineIsNew &&
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
        if (!comment.codeContext) {
          continue
        }
        if (
          newRow != undefined &&
          comment.codeContext.lineIsNew === true &&
          comment.codeContext.lineStart === newRow &&
          comment.id === id
        ) {
          return true
        }
        if (
          oldRow != undefined &&
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
      const hunkPatchIDs = new Set<string>()

      this.checkedHunks.forEach((isChecked, patchID) => {
        // un-set if hunk no longer exists
        if (!this.hunkIds.has(patchID)) {
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
      this.emitter.emit('undo-file', { patch_ids: this.hunkIds })
    },
    sendSetWithPrefix(prefix: string, selected: boolean) {
      // Send event to siblings, etc.
      this.emitter.emit('differ-set-hunks-with-prefix', {
        prefix: prefix,
        selected: selected,
      })
    },
    highlightedBlocks(input: Block[], lang: string): HighlightedBlock[] {
      return highlight(input, lang, this.enableSyntaxHighlight)
    },
    onClickDismissHunkedSuggestion(
      hunk: SuggestionHunk,
      suggestion: DifferFile_SuggestionFragment
    ) {
      this.$emit('dismissHunkedSuggestion', {
        suggestionId: suggestion.id,
        hunks: [hunk.id],
      })
    },
    onClickApplyHunkedSuggestion(hunk: SuggestionHunk, suggestion: DifferFile_SuggestionFragment) {
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
      // Check if this component has any suggestions for this user
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
    searchMatchesFiles(fileName: string): boolean {
      if (!this.searchResult) {
        return false
      }

      return this.searchResult.has(fileName)
    },
    searchIsCurrentSelectedFilename(filename: string): boolean {
      return this.searchCurrentSelectedId === filename
    },
    searchMatchesHunk(hunkID: string): boolean {
      if (!this.searchResult) {
        return false
      }
      let r = this.searchResult.get(hunkID)
      if (r && r.length > 0) {
        return true
      }
      return false
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
