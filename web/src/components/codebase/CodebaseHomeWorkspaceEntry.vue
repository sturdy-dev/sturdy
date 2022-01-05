<template>
  <tr>
    <td class="px-6 py-4 whitespace-nowrap w-auto">
      <div class="flex items-center">
        <div class="flex-shrink-0 h-10 w-10">
          <Avatar
            :author="ws.author"
            size="10"
            :show-online="true"
            online-size="h-3 w-3"
            :online="viewIsOnline"
          />
        </div>
        <div class="ml-4">
          <div class="text-sm text-gray-900">
            {{ ws.name }}
          </div>
          <div class="text-sm text-gray-500">
            {{ ws.author.name }}
          </div>
        </div>
      </div>
    </td>
    <td class="">
      <div v-if="ws.view" class="whitespace-nowrap">
        <div class="text-sm text-gray-900">
          {{ ws.view.shortMountPath }}
        </div>
        <div class="text-sm text-gray-500">
          {{ ws.view.mountHostname }}
        </div>
      </div>
    </td>
    <td class="px-6 py-4 whitespace-nowrap w-full hidden lg:table-cell">
      <span
        v-if="ws.lastActivityAt > 0"
        class="px-2 inline-flex text-xs leading-5 font-semibold rounded-full bg-green-100 text-green-800"
        :class="[ws.createdAt === ws.lastActivityAt ? 'bg-blue-100' : 'bg-green-100']"
        :title="[
          ws.createdAt === ws.lastActivityAt
            ? 'Newly created workspace'
            : 'Last change in this workspace',
        ]"
      >
        {{ friendly_ago }}
      </span>
      <span
        v-else
        class="px-2 inline-flex text-xs leading-5 font-semibold rounded-full bg-gray-100 text-gray-800"
      >
        Inactive
      </span>
    </td>
    <td class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium hidden lg:table-cell">
      <a
        href="#"
        class="text-blue-600 hover:text-blue-900"
        @click.stop.prevent="$emit('openArchiveModal')"
        >Archive</a
      >
    </td>
  </tr>
</template>
<script>
import Avatar from '../shared/Avatar.vue'
import { gql, useSubscription } from '@urql/vue'
import { ref, toRef, watch } from 'vue'
import { useUpdatedView } from '../../subscriptions/useUpdatedView'

export default {
  name: 'CodebaseHomeWorkspaceEntry',
  components: { Avatar },
  props: {
    ws: {},
    now: {},
    friendly_ago: {},
  },
  emits: ['openArchiveModal'],
  setup(props) {
    let ws = toRef(props, 'ws')
    let viewID = ref(ws?.value?.view?.id)
    watch(ws, () => {
      viewID.value = ws?.value?.view?.id
    })

    useUpdatedView(viewID, {
      pause: !viewID.value,
    })
  },
  computed: {
    viewIsOnline() {
      return this.ws.view?.lastUsedAt > this.now / 1000 - 120
    },
  },
}
</script>
