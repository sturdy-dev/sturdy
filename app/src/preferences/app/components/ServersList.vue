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
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<script>
import Server from './Server.vue'
import ServerInput from './ServerInput.vue'
import { useStore } from '@nanostores/vue'
import { servers, add as addServer } from '../stores/servers'
import ipc from '../ipc'

export default {
  components: { Server, ServerInput },
  setup() {
    ipc.listHosts().then((hosts) => hosts.forEach(addServer))
    return {
      servers: useStore(servers),
    }
  },
}
</script>
