<template>
  <div v-if="workspaces" class="">
    <div v-if="workspaces.length > 0" class="flex flex-col">
      <div class="-my-2 overflow-x-auto sm:-mx-6 lg:-mx-8">
        <div class="py-2 align-middle inline-block min-w-full sm:px-6 lg:px-8">
          <div class="shadow overflow-hidden border-b border-gray-200 sm:rounded-lg">
            <table class="min-w-full divide-y divide-gray-200 table-fixed">
              <thead class="bg-gray-50">
                <tr>
                  <th
                    scope="col"
                    class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
                  >
                    Workspace
                  </th>
                  <th></th>
                </tr>
              </thead>
              <tbody class="bg-white divide-y divide-gray-200">
                <template v-for="ws in sortedWorkspaces" :key="ws.id">
                  <tr class="hover:bg-gray-100 cursor-pointer" @click="goto(ws)">
                    <td class="px-2 py-4 text-sm text-gray-900 w-full flex items-center gap-2">
                      <Avatar :author="ws.author" size="6" class="flex-shrink-0" />
                      <span>{{ ws.name }}</span>
                    </td>
                    <td class="text-sm text-gray-500 whitespace-nowrap px-2 py-4">
                      {{ friendly_ago(ws.lastActivityAt) }}
                    </td>
                  </tr>
                </template>
              </tbody>
            </table>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { gql } from '@urql/vue'
import { defineComponent } from 'vue'
import type { PropType } from 'vue'
import type { WorkspaceListFragment } from './__generated__/WorkspaceList'
import Avatar from '../../atoms/Avatar.vue'
import time from '../../time'

export const WORKSPACE_LIST = gql`
  fragment WorkspaceList on Workspace {
    id
    name
    author {
      id
      name
      avatarUrl
    }
    suggestion {
      id
    }
    lastActivityAt
  }
`

const nonSuggestingWorkspaces = (ws: WorkspaceListFragment) => !ws.suggestion

const workspaceByLastUpdated = (a: WorkspaceListFragment, b: WorkspaceListFragment) =>
  b.lastActivityAt - a.lastActivityAt

export default defineComponent({
  components: {
    Avatar,
  },
  props: {
    workspaces: {
      type: Object as PropType<Array<WorkspaceListFragment>>,
      required: true,
    },
  },
  computed: {
    sortedWorkspaces() {
      let ws = this.workspaces
      return ws?.filter(nonSuggestingWorkspaces).sort(workspaceByLastUpdated)
    },
  },
  methods: {
    goto(ws) {
      this.$router.push({ name: 'workspaceHome', params: { id: ws.id } })
    },
    friendly_ago(ts) {
      return time.getRelativeTime(new Date(ts * 1000))
    },
  },
})
</script>
