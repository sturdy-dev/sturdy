<template>
  <div class="flex py-1.5 px-2.5 sticky z-20 left-0 inline-flex items-start w-full items-center">
    <span class="inline-flex items-center space-x-2 text-sm font-medium text-gray-500">
      <!-- TODO: Different icons for different filetypes -->
      <DocumentTextIcon class="w-5 h-5" />

      <span v-if="diffs.isNew" class="d2h-file-name">
        <DifferName :name="diffs.newName" :added="isAdded" @addWithPrefix="sendSetWithPrefix" />
      </span>

      <span v-else-if="diffs.isDeleted" class="d2h-file-name">
        <DifferName :name="diffs.origName" :added="isAdded" @addWithPrefix="sendSetWithPrefix" />
      </span>

      <span
        v-else-if="diffs.newName && diffs.origName && diffs.newName !== diffs.origName"
        class="d2h-file-name"
      >
        <DifferName :name="diffs.origName" :added="isAdded" @addWithPrefix="sendSetWithPrefix" />
        →
        <DifferName :name="diffs.origName" :added="isAdded" @addWithPrefix="sendSetWithPrefix" />
      </span>

      <span v-else class="d2h-file-name">
        <DifferName :name="diffs.origName" :added="isAdded" @addWithPrefix="sendSetWithPrefix" />
      </span>

      <span v-if="diffs.isNew" class="d2h-tag d2h-added d2h-added-tag">ADDED</span>
      <span v-else-if="diffs.isDeleted" class="d2h-tag d2h-deleted d2h-deleted-tag"> DELETED </span>
      <span v-else-if="diffs.isMoved" class="d2h-tag d2h-moved d2h-moved-tag">MOVED</span>
      <span v-else class="d2h-tag d2h-changed d2h-changed-tag">CHANGED</span>

      <span v-if="diffs.isLarge">
        {{ humanSize(diffs.largeFileInfo.size) }}
      </span>
    </span>

    <!-- Suggestion heads -->
    <div v-if="!isSuggesting" class="ml-8">
      <div v-for="author in suggestingAuthors" :key="author.id">
        <div
          :class="[
            showingSuggestionsByUser === author.id ? 'bg-gray-200' : 'bg-transparent',
            'w-8 h-8 rounded-full items-center justify-center inline-flex',
          ]"
        >
          <!-- Avatar if suggestions can be hidden (user have their own modifications to this file) -->
          <Avatar
            v-if="haveLiveChanges"
            :author="author"
            size="6"
            class="border-2 border-green-300 cursor-pointer rounded-full"
            :title="'Show suggestions by ' + author.name"
            @click="$emit('showSuggestionsByUser', author.id)"
          />
          <Avatar
            v-else
            :author="author"
            size="6"
            class="border-2 border-green-300 rounded-full"
            :title="'Show suggestions by ' + author.name"
          />
        </div>
      </div>
    </div>

    <div class="flex-grow" />
    <RouterLinkButton
      v-if="showFullFileButton && diffs.newName"
      size="small"
      :to="{ name: 'browseFile', params: { path: diffs.newName.split('/') } }"
    >
      Show file
    </RouterLinkButton>
    <DifferAddButton
      v-if="!showSuggestions && showAddButton"
      :is-added="isAdded"
      :can-ignore-file="canIgnoreFile"
      :is-hidden="false"
      @add="$emit('add')"
      @ignore="onIgnore(diffs.newName)"
      @hide="$emit('hide')"
      @undo="$emit('undo')"
      @unhide="$emit('unhide')"
      @showdropdown="$emit('showdropdown')"
      @hidedropdown="$emit('hidedropdown')"
    />
  </div>
  <!-- Chosen conflict resolution -->
  <div v-if="conflictSelection" class="w-full border-t border-gray-200">
    <div class="pt-1 w-full flex flex-grow justify-center text-sm text-gray-500 text-center">
      <p v-if="conflictSelection === 'trunk'">Using version from trunk</p>
      <p v-else-if="conflictSelection === 'workspace'">Using version from workspace</p>
      <p v-else-if="conflictSelection === 'custom'">Using custom version</p>
      <p v-else>Please choose conflict resolution</p>
      <CheckCircleIcon
        v-if="['trunk', 'workspace', 'custom'].includes(conflictSelection)"
        class="-mr-1 ml-1 h-5 w-5 text-green-700"
      />
      <ClockIcon v-else class="-mr-1 ml-1 h-5 w-5 text-yellow-500" />
    </div>
  </div>
</template>

<script lang="ts">
import DifferName from './DifferName.vue'
import DifferAddButton from './DifferAddButton.vue'
import { CheckCircleIcon, ClockIcon, DocumentTextIcon } from '@heroicons/vue/outline'
import RouterLinkButton from '../shared/RouterLinkButton.vue'
import Avatar from '../shared/Avatar.vue'
import { defineComponent, PropType } from 'vue'
import { gql } from '@urql/vue'
import {
  DiffHeader_FileDiffFragment,
  DiffHeader_SuggestionsFragment,
} from './__generated__/DiffHeader'

const deduplicateAuthors = (authors: Author[]) => {
  const seen = new Set()
  return authors.filter((author) => {
    const isNew = !seen.has(author.id)
    seen.add(author.id)
    return isNew
  })
}

type Author = DiffHeader_SuggestionsFragment['author']

export const DIFF_HEADER_FILE_DIFF = gql`
  fragment DiffHeader_FileDiff on FileDiff {
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
  }
`

export const DIFF_HEADER_SUGGESTIONS = gql`
  fragment DiffHeader_Suggestions on Suggestion {
    id
    author {
      id
      name
      avatarUrl
    }
  }
`

export default defineComponent({
  components: {
    DifferName,
    DifferAddButton,
    CheckCircleIcon,
    ClockIcon,
    Avatar,
    RouterLinkButton,
    DocumentTextIcon: DocumentTextIcon,
  },
  props: {
    diffs: {
      type: Object as PropType<DiffHeader_FileDiffFragment>,
      required: true,
    },

    isSuggesting: Boolean,
    showSuggestions: Boolean,
    suggestions: {
      type: Object as PropType<Array<DiffHeader_SuggestionsFragment>>,
      required: false,
      default: function () {
        return new Array<DiffHeader_SuggestionsFragment>()
      },
    },
    showingSuggestionsByUser: {
      type: String,
      required: false,
      default: function () {
        return null
      },
    },

    conflictSelection: String,

    isAdded: Boolean,
    haveLiveChanges: Boolean,

    canIgnoreFile: {
      type: Boolean,
      required: true,
    },
    canTakeSuggestions: Boolean,
    showFullFileButton: Boolean,
    showAddButton: {
      type: Boolean,
      required: true,
    },
  },
  emits: [
    'add',
    'hide',
    'undo',
    'ignore',
    'showdropdown',
    'hidedropdown',
    'unhide',
    'showSuggestionsByUser',
  ],
  computed: {
    suggestingAuthors() {
      if (!this.suggestions) return []
      const authors = this.suggestions.map((suggestion) => suggestion.author)
      return deduplicateAuthors(authors)
    },
  },
  methods: {
    sendSetWithPrefix(prefix: string, selected: boolean) {
      // Send event to siblings, etc.
      this.emitter.emit('differ-set-hunks-with-prefix', {
        prefix: prefix,
        selected: selected,
      })
    },
    onIgnore(fileName: string) {
      this.$emit('ignore', fileName)
    },
    round2(num: number) {
      return Math.round((num + Number.EPSILON) * 100) / 100
    },
    humanSize(bytes: number) {
      if (bytes < 1024) {
        return bytes + 'B'
      }
      let kBytes = bytes / 1024
      if (kBytes < 1014) {
        return this.round2(kBytes) + 'kB'
      }
      let mBytes = kBytes / 1024
      if (mBytes < 1014) {
        return this.round2(mBytes) + 'MB'
      }
      let gBytes = mBytes / 1024
      return this.round2(gBytes) + 'GB'
    },
  },
})
</script>
