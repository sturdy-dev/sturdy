<template>
  <div class="flex flex-1 flex-col">
    <Toggle class="p-2" v-model="isDetailed" label="Detailed" />

    <div class="align-middle inline-block min-w-full">
      <div class="shadow overflow-hidden border-b border-gray-200 sm:rounded-lg">
        <table class="min-w-full divide-y divide-gray-200">
          <thead class="bg-gray-50">
            <tr>
              <th
                v-for="header in headers"
                :v-key="header"
                scope="col"
                class="px-3 py-1 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
              >
                {{ header }}
              </th>
            </tr>
          </thead>
          <tbody>
            <Server
              v-for="server in servers"
              :key="server.title"
              :server="server"
              :is-detailed="isDetailed"
            />
            <ServerInput
              @error="onInputError"
              @success="onInputSuccess"
              :is-detailed="isDetailed"
            />
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
import Toggle from './Toggle.vue'
import { useStore } from '@nanostores/vue'
import { servers, set, add as addServer } from '../stores/servers'
import ipc from '../ipc'
import { ref } from 'vue'

export default {
  components: { Server, ServerInput, Toggle },
  data() {
    return {
      detailed: false,
      error: null,
    }
  },
  setup() {
    const isDetailed = ref(false)
    ipc.listHosts().then((hosts) => {
      set(hosts)
      isDetailed.value = hosts.some(({ apiURL }) => !!apiURL)
    })
    return {
      servers: useStore(servers),
      isDetailed,
    }
  },
  computed: {
    headers() {
      return this.isDetailed
        ? [
            '' /* status */,
            'Title',
            'Web Url',
            'API Url',
            'Sync Url',
            '' /* delete */,
            '' /* open */,
          ]
        : ['' /* status */, 'Title', 'Host', '' /* delete */, '' /* open */]
    },
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
