<template>
  <div class="flex py-1.5 px-2.5 sticky z-20 left-0 inline-flex items-start w-full items-center">
    <span class="inline-flex items-center space-x-2 text-sm font-medium text-gray-500">
      <!-- TODO: Different icons for different filetypes -->
      <DocumentTextIcon class="w-5 h-5" />

      <span v-if="compatIsNew" class="d2h-file-name">
        <DifferName :name="compatNewName" :added="isAdded" @addWithPrefix="sendSetWithPrefix" />
      </span>

      <span v-else-if="compatIsDeleted" class="d2h-file-name">
        <DifferName :name="compatOrigName" :added="isAdded" @addWithPrefix="sendSetWithPrefix" />
      </span>

      <span
        v-else-if="compatNewName && compatOrigName && compatNewName !== compatOrigName"
        class="d2h-file-name"
      >
        <DifferName :name="compatOrigName" :added="isAdded" @addWithPrefix="sendSetWithPrefix" />
        →
        <DifferName :name="compatNewName" :added="isAdded" @addWithPrefix="sendSetWithPrefix" />
      </span>

      <span v-else class="d2h-file-name">
        <DifferName :name="compatOrigName" :added="isAdded" @addWithPrefix="sendSetWithPrefix" />
      </span>

      <span v-if="compatIsNew" class="d2h-tag d2h-added d2h-added-tag">ADDED</span>
      <span v-else-if="compatIsDeleted" class="d2h-tag d2h-deleted d2h-deleted-tag"> DELETED </span>
      <span v-else-if="compatIsMoved" class="d2h-tag d2h-moved d2h-moved-tag">MOVED</span>
      <span v-else class="d2h-tag d2h-changed d2h-changed-tag">CHANGED</span>

      <span v-if="compatIsLarge">
        {{ humanSize(compatLargeFileInfo.size) }}
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

    <!-- Chosen conflict resolution -->
    <div v-if="conflictSelection" class="flex text-sm text-gray-500 text-center">
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

    <div class="flex-grow" />
    <RouterLinkButton
      v-if="showFullFileButton"
      size="small"
      :to="{ name: 'browseFile', params: { path: compatNewName.split('/') } }"
    >
      Show file
    </RouterLinkButton>
    <DifferAddButton
      v-if="!showSuggestions && showAddButton"
      :is-added="isAdded"
      :can-ignore-file="canIgnoreFile"
      @add="$emit('add')"
      @ignore="onIgnore(diffs.new_name)"
      @hide="$emit('hide')"
      @undo="$emit('undo')"
      @unhide="$emit('unhide')"
      @showdropdown="$emit('showdropdown')"
      @hidedropdown="$emit('hidedropdown')"
    />
  </div>
</template>

<script lang="js">
import DifferName from "./DifferName.vue";
import DifferAddButton from "./DifferAddButton.vue";
import { CheckCircleIcon, ClockIcon, DocumentTextIcon } from "@heroicons/vue/outline";
import RouterLinkButton from "../shared/RouterLinkButton.vue";
import Avatar from "../shared/Avatar.vue";

const deduplicateAuthors = (authors) => {
  const seen = new Set()
  return authors.filter(author => {
    const isNew = !seen.has(author.id)
    seen.add(author.id)
    return isNew
  })
}

export default {
  name: 'DiffHeader',
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
    isSuggesting: Boolean,
    diffs: {
      type: Object,
      required: true,
    },
    suggestions: Array,
    conflictSelection: String,
    showSuggestions: Boolean,
    isAdded: Boolean,
    fileKey: String,
    haveLiveChanges: Boolean,
    showingSuggestionsByUser: String,
    canIgnoreFile: Boolean,
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
    compatIsNew() {
      return this.diffs.is_new || this.diffs.isNew
    },
    compatIsMoved() {
      return this.diffs.is_moved || this.diffs.isMoved
    },
    compatIsDeleted() {
      return this.diffs.is_deleted || this.diffs.isDeleted
    },
    compatNewName() {
      return this.diffs.new_name || this.diffs.newName
    },
    compatOrigName() {
      return this.diffs.orig_name || this.diffs.origName
    },
    compatIsLarge() {
      return this.diffs.is_large || this.diffs.isLarge
    },
    compatLargeFileInfo() {
      return this.diffs.large_file_info || this.diffs.largeFileInfo
    },
  },
  methods: {
    sendSetWithPrefix(prefix, selected) {
      // Send event to siblings, etc.
      this.emitter.emit('differ-set-hunks-with-prefix', {
        prefix: prefix,
        selected: selected,
      })
    },
    onIgnore(fileName) {
      this.$emit('ignore', fileName)
    },
    round2(num) {
      return Math.round((num + Number.EPSILON) * 100) / 100
    },
    humanSize(bytes) {
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
}
</script>
