<template>
  <tr class="bg-white">
    <td class="px-3 py-2 whitespace-nowrap text-sm font-medium text-gray-900"></td>

    <td class="px-3 py-2 whitespace-nowrap text-sm font-medium text-gray-900">
      <input
        type="text"
        class="focus:ring-blue-500 focus:border-blue-500 relative block w-full rounded-none rounded-t-md bg-transparent focus:z-10 sm:text-sm border-gray-300"
        placeholder="My Sturdy"
        v-model="title"
        :class="[
          isTitleWarning
            ? 'bg-yellow-50 focus:ring-yellow-500 focus:border-yellow-500'
            : 'focus:ring-blue-500 focus:border-blue-500 ',
        ]"
      />
    </td>

    <td v-if="!isDetailed" class="px-3 py-2 whitespace-nowrap text-sm font-medium text-gray-900">
      <input
        type="text"
        class="focus:ring-blue-500 focus:border-blue-500 relative block w-full rounded-none rounded-t-md bg-transparent focus:z-10 sm:text-sm border-gray-300"
        placeholder="http://localhost:30080"
        v-model="host"
        :class="[
          isHostWarning
            ? 'bg-yellow-50 focus:ring-yellow-500 focus:border-yellow-500'
            : 'focus:ring-blue-500 focus:border-blue-500 ',
        ]"
      />
    </td>

    <td v-if="isDetailed" class="px-3 py-2 whitespace-nowrap text-sm font-medium text-gray-900">
      <input
        type="text"
        class="focus:ring-blue-500 focus:border-blue-500 relative block w-full rounded-none rounded-t-md bg-transparent focus:z-10 sm:text-sm border-gray-300"
        placeholder="https://getsturdy.com"
        v-model="webURL"
        :class="[
          isWebURLWarning
            ? 'bg-yellow-50 focus:ring-yellow-500 focus:border-yellow-500'
            : 'focus:ring-blue-500 focus:border-blue-500 ',
        ]"
      />
    </td>

    <td v-if="isDetailed" class="px-3 py-2 whitespace-nowrap text-sm font-medium text-gray-900">
      <input
        type="text"
        class="focus:ring-blue-500 focus:border-blue-500 relative block w-full rounded-none rounded-t-md bg-transparent focus:z-10 sm:text-sm border-gray-300"
        placeholder="https://api.getsturdy.com"
        v-model="apiURL"
        :class="[
          isAPIURLWarning
            ? 'bg-yellow-50 focus:ring-yellow-500 focus:border-yellow-500'
            : 'focus:ring-blue-500 focus:border-blue-500 ',
        ]"
      />
    </td>

    <td v-if="isDetailed" class="px-3 py-2 whitespace-nowrap text-sm font-medium text-gray-900">
      <input
        type="text"
        class="focus:ring-blue-500 focus:border-blue-500 relative block w-full rounded-none rounded-t-md bg-transparent focus:z-10 sm:text-sm border-gray-300"
        placeholder="ssh://sync.getsturdy.com"
        v-model="syncURL"
        :class="[
          isSyncURLWarning
            ? 'bg-yellow-50 focus:ring-yellow-500 focus:border-yellow-500'
            : 'focus:ring-blue-500 focus:border-blue-500 ',
        ]"
      />
    </td>

    <td class="px-3 py-2 whitespace-nowrap text-sm font-medium text-gray-900">
      <button type="button" class="text-blue-600 hover:text-blue-900" @click.prevent="handleAdd">
        Add
      </button>
    </td>
  </tr>
</template>

<script>
import ipc from '../ipc'

export default {
  props: {
    isDetailed: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    return {
      title: '',
      titleActivated: false,

      host: '',
      hostActivated: false,

      webURL: '',
      webURLActivated: false,

      apiURL: '',
      apiURLActivated: false,

      syncURL: '',
      syncURLActivated: false,
    }
  },
  emits: ['error', 'success'],
  watch: {
    title(n) {
      if (n.length > 0) {
        this.titleActivated = true
      }
    },
    host(n) {
      if (n.length > 0) {
        this.hostActivated = true
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
    activated() {
      return (
        this.titleActivated ||
        this.hostActivated ||
        this.webURLActivated ||
        this.apiURLActivated ||
        this.syncURLActivated
      )
    },

    isTitleWarning() {
      return this.activated && !this.isTitleValid
    },
    isTitleValid() {
      return this.title.length > 0
    },

    isHostWarning() {
      return this.activated && !this.isHostValid
    },
    isHostValid() {
      try {
        new URL(this.host)
        return true
      } catch {
        return false
      }
    },

    isWebURLWarning() {
      return this.activated && !this.isWebURLValid
    },
    isWebURLValid() {
      try {
        new URL(this.webURL)
        return true
      } catch {
        return false
      }
    },

    isAPIURLWarning() {
      return this.activated && !this.isAPIURLValid
    },
    isAPIURLValid() {
      try {
        new URL(this.apiURL)
        return true
      } catch {
        return false
      }
    },

    isSyncURLWarning() {
      return this.activated && !this.isSyncURLValid
    },
    isSyncURLValid() {
      try {
        new URL(this.syncURL)
        return true
      } catch {
        return false
      }
    },

    isValid() {
      return this.isDetailed
        ? this.isTitleValid && this.isWebURLValid && this.isAPIURLValid && this.isSyncURLValid
        : this.isTitleValid && this.isHostValid
    },
    hostConfig() {
      return this.isDetailed
        ? {
            title: this.title,
            webURL: this.webURL,
            apiURL: this.apiURL,
            syncURL: this.syncURL,
          }
        : {
            title: this.title,
            host: this.host,
          }
    },
  },
  methods: {
    async handleAdd() {
      if (!this.isValid) return

      try {
        await ipc.addHostConfig(this.hostConfig)
      } catch (e) {
        const errorMessage = e.message.split('Error:').slice(-1)[0]
        this.$emit('error', errorMessage)
        return
      }

      this.$emit('success', this.hostConfig)
      this.resetValues()
      this.resetActivated()
    },
    resetActivated() {
      this.titleActivated = false
      this.hostActivated = false
      this.webURLActivated = false
      this.apiURLActivated = false
      this.syncURLActivated = false
    },
    resetValues() {
      this.title = ''
      this.host = ''
      this.webURL = ''
      this.apiURL = ''
      this.syncURL = ''
    },
  },
}
</script>
