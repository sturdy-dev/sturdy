<template>
  <tr class="bg-white">
    <td class="px-3 py-2 flex items-center">
      <ServerStatus :server="server" />
    </td>
    <td
      v-for="field in fields"
      :key="field"
      class="px-3 py-2 whitespace-nowrap text-sm font-medium text-gray-900"
    >
      {{ field }}
    </td>
    <td class="px-3 py-2 whitespace-nowrap text-sm font-medium text-gray-900">
      <button type="button" class="text-red-600 hover:text-red-900" @click.prevent="handleDelete">
        Delete
      </button>
    </td>
    <td class="px-3 py-2 whitespace-nowrap text-sm font-medium text-gray-900">
      <button
        type="button"
        class="text-indigo-600 hover:text-indigo-900"
        @click.prevent="handleOpen"
      >
        Open
      </button>
    </td>
  </tr>
</template>

<script>
import ServerStatus from './ServerStatus.vue'
import ipc from '../ipc'
import { remove as deleteServer } from '../stores/servers'

export default {
  components: {
    ServerStatus,
  },
  props: {
    server: {
      type: Object,
      required: true,
    },
  },
  computed: {
    fields() {
      return [
        this.server.title,
        this.server.webURL,
        this.server.apiURL,
        this.server.syncURL,
        this.server.reposBasePath,
      ]
    },
  },
  methods: {
    handleOpen() {
      ipc.isHostUp(this.server).then((isUp) => {
        if (isUp) ipc.openHost(this.server)
      })
    },
    handleDelete() {
      deleteServer(this.server)
      ipc.deleteHost(this.server)
    },
  },
}
</script>
