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

    <td class="px-3 py-2 whitespace-nowrap text-sm font-medium text-gray-900">
      <input
        type="text"
        class="focus:ring-blue-500 focus:border-blue-500 relative block w-full rounded-none rounded-t-md bg-transparent focus:z-10 sm:text-sm border-gray-300"
        placeholder="127.0.0.1:30080"
        v-model="host"
        :class="[
          isHostWarning
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
  data() {
    return {
      title: '',
      titleActivated: false,

      host: '',
      hostActivated: false,
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
  },
  computed: {
    isTitleWarning() {
      return this.titleActivated && !this.isTitleValid
    },
    isTitleValid() {
      return this.title.length > 0
    },

    isHostWarning() {
      return this.hostActivated && !this.isHostValid
    },
    isHostValid() {
      return this.host.length > 0
    },

    isValid() {
      return this.isTitleValid && this.isHostValid
    },
    hostConfig() {
      return {
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
    },
    resetValues() {
      this.title = ''
      this.host = ''
    },
  },
}
</script>
