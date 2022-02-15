<template>
  <div class="flex flex-col">
    <div class="-my-2 overflow-x-auto sm:-mx-6 lg:-mx-8">
      <div class="py-2 align-middle inline-block min-w-full sm:px-6 lg:px-8">
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
                  Web URL
                </th>
                <th
                  scope="col"
                  class="px-3 py-1 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
                >
                  API URL
                </th>
                <th
                  scope="col"
                  class="px-3 py-1 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
                >
                  Sync URL
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
            <tbody id="servers">
              <template v-for="server in servers" :key="server.title">
                <Server :server="server" />
              </template>
              <ServerInput />
            </tbody>
          </table>
        </div>
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
    ipc.getConfig().then((config) => config.hosts.forEach(addServer))
    return {
      servers: useStore(servers),
    }
  },
}
</script>
