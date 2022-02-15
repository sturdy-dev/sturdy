<template>
  <tr class="bg-white">
    <td class="px-3 py-2 whitespace-nowrap text-sm font-medium text-gray-900"></td>

    <td class="px-3 py-2 whitespace-nowrap text-sm font-medium text-gray-900">
      <input
        type="text"
        class="focus:ring-indigo-500 focus:border-indigo-500 relative block w-full rounded-none rounded-t-md bg-transparent focus:z-10 sm:text-sm border-gray-300"
        placeholder="My Sturdy"
        v-model="title"
        :class="[
          isTitleWarning
            ? 'bg-yellow-50 focus:ring-yellow-500 focus:border-yellow-500'
            : 'focus:ring-indigo-500 focus:border-indigo-500 ',
        ]"
      />
    </td>

    <td class="px-3 py-2 whitespace-nowrap text-sm font-medium text-gray-900">
      <input
        type="text"
        class="focus:ring-indigo-500 focus:border-indigo-500 relative block w-full rounded-none rounded-t-md bg-transparent focus:z-10 sm:text-sm border-gray-300"
        placeholder="https://example.com"
        v-model="webURL"
        :class="[
          isWebURLWarning
            ? 'bg-yellow-50 focus:ring-yellow-500 focus:border-yellow-500'
            : 'focus:ring-indigo-500 focus:border-indigo-500 ',
        ]"
      />
    </td>

    <td class="px-3 py-2 whitespace-nowrap text-sm font-medium text-gray-900">
      <input
        type="text"
        class="focus:ring-indigo-500 focus:border-indigo-500 relative block w-full rounded-none rounded-t-md bg-transparent focus:z-10 sm:text-sm border-gray-300"
        placeholder="https://api.example.com"
        v-model="apiURL"
        :class="[
          isApiURLWarning
            ? 'bg-yellow-50 focus:ring-yellow-500 focus:border-yellow-500'
            : 'focus:ring-indigo-500 focus:border-indigo-500 ',
        ]"
      />
    </td>

    <td class="px-3 py-2 whitespace-nowrap text-sm font-medium text-gray-900">
      <input
        type="text"
        class="focus:ring-indigo-500 focus:border-indigo-500 relative block w-full rounded-none rounded-t-md bg-transparent focus:z-10 sm:text-sm border-gray-300"
        placeholder="ssh://example.com"
        v-model="syncURL"
        :class="[
          isSyncURLWarning
            ? 'bg-yellow-50 focus:ring-yellow-500 focus:border-yellow-500'
            : 'focus:ring-indigo-500 focus:border-indigo-500 ',
        ]"
      />
    </td>

    <td class="px-3 py-2 whitespace-nowrap text-sm font-medium text-gray-900">
      <button
        type="button"
        class="text-indigo-600 hover:text-indigo-900"
        @click.prevent="handleAdd"
      >
        Add
      </button>
    </td>
  </tr>
</template>

<script>
import { add as addServer } from '../stores/servers'
import ipc from '../ipc'

export default {
  data() {
    return {
      title: '',
      titleActivated: false,

      webURL: '',
      webURLActivated: false,

      apiURL: '',
      apiURLActivated: false,

      syncURL: '',
      syncURLActivated: false,
    }
  },
  watch: {
    title(n) {
      if (n.length > 0) {
        this.titleActivated = true
      }
    },
    webURL(n) {
      if (n.length > 0) {
        this.webURLActivated = true
      }
    },
    apiURL(n) {
      if (n.length > 0) {
        this.apiURLActivated = true
      }
    },
    syncURL(n) {
      if (n.length > 0) {
        this.syncURLActivated = true
      }
    },
  },
  computed: {
    isTitleWarning() {
      return this.titleActivated && !this.isTitleValid
    },
    isTitleValid() {
      return this.title.length > 0
    },

    isWebURLWarning() {
      return this.webURLActivated && !this.isWebURLValid
    },
    isWebURLValid() {
      try {
        new URL(this.webURL)
        return true
      } catch (e) {
        return false
      }
    },

    isApiURLWarning() {
      return this.apiURLActivated && !this.isApiURLValid
    },
    isApiURLValid() {
      try {
        new URL(this.apiURL)
        return true
      } catch (e) {
        return false
      }
    },

    isSyncURLWarning() {
      return this.syncURLActivated && !this.isSyncURLValid
    },
    isSyncURLValid() {
      try {
        new URL(this.syncURL)
        return true
      } catch (e) {
        return false
      }
    },

    isValid() {
      return this.isTitleValid && this.isWebURLValid && this.isApiURLValid && this.isSyncURLValid
    },
    hostConfig() {
      return {
        title: this.title,
        webURL: this.webURL,
        apiURL: this.apiURL,
        syncURL: this.syncURL,
      }
    },
  },
  methods: {
    handleAdd() {
      if (!this.isValid) return

      addServer(this.hostConfig)
      ipc.addHostConfig(this.hostConfig)

      this.resetValues()
      this.resetActivated()
    },
    resetActivated() {
      this.titleActivated = false
      this.webURLActivated = false
      this.apiURLActivated = false
      this.syncURLActivated = false
    },
    resetValues() {
      this.title = ''
      this.webURL = ''
      this.apiURL = ''
      this.syncURL = ''
    },
  },
}
</script>
