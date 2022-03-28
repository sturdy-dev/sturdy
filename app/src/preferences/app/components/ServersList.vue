<template>
  <div class="flex flex-1 flex-col">
    <div class="align-middle inline-block min-w-full">
      <div class="shadow overflow-hidden border-b border-gray-200 sm:rounded-lg">
        <table class="min-w-full divide-y divide-gray-200">
          <thead class="bg-gray-50">
            <tr>
              <th scope="col"></th>
              <th
                scope="col"
                class="px-3 py-1 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
              >
                Title
              </th>
              <th
                scope="col"
                class="px-3 py-1 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
              >
                Host
              </th>
              <th
                scope="col"
                class="px-3 py-1 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
              ></th>
              <th
                scope="col"
                class="px-3 py-1 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
              ></th>
            </tr>
          </thead>
          <tbody>
            <Server v-for="server in servers" :key="server.title" :server="server" />
            <ServerInput @error="onInputError" @success="onInputSuccess" />
          </tbody>
        </table>
        <p v-if="error" class="m-2 text-red-600 font-xs">{{ error }}</p>
      </div>
    </div>
  </div>
</template>

<script>
import Server from './Server.vue'
import ServerInput from './ServerInput.vue'
import { useStore } from '@nanostores/vue'
import { servers, set } from '../stores/servers'
import ipc from '../ipc'

export default {
  components: { Server, ServerInput },
  data() {
    return {
      error: null,
    }
  },
  setup() {
    ipc.listHosts().then((hosts) => set(hosts))
    return {
      servers: useStore(servers),
    }
  },
  methods: {
    onInputError(error) {
      this.error = error
    },
    onInputSuccess(hostConfig) {
      addServer(hostConfig)
      this.error = null
    },
  },
}
</script>
