<template>
  <div v-if="sortedWorkspaces.length === 0">
    <EmptyPlaceholder
      title="Nothing"
      description="This codebase have no archived workspaces. If you
    archive one later, you'll find it here!"
      :show-icon="false"
    />
  </div>

  <div v-if="sortedWorkspaces.length > 0" class="flex flex-col">
    <div class="-my-2 overflow-x-auto sm:-mx-6 lg:-mx-8">
      <div class="py-2 align-middle inline-block min-w-full sm:px-6 lg:px-8">
        <div class="shadow overflow-hidden border-b border-gray-200 sm:rounded-lg">
          <table class="min-w-full divide-y divide-gray-200">
            <thead class="bg-gray-50">
              <tr>
                <th
                  scope="col"
                  class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
                >
                  Name
                </th>
                <th
                  scope="col"
                  class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
                >
                  Archived
                </th>
                <th scope="col" class="relative px-6 py-3">
                  <span class="sr-only">Edit</span>
                </th>
              </tr>
            </thead>
            <tbody class="bg-white divide-y divide-gray-200">
              <tr
                v-for="ws in sortedWorkspaces"
                :key="ws.id"
                :class="[!ws.archivedAt ? 'hover:bg-gray-100 cursor-pointer' : '']"
                @click="maybeGoto(ws)"
              >
                <td class="px-6 py-4 whitespace-nowrap">
                  <div class="flex items-center">
                    <div class="flex-shrink-0 h-10 w-10">
                      <Avatar :author="ws.author" size="10" />
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
                <td class="px-6 py-4 whitespace-nowrap">
                  <span
                    v-if="ws.archivedAt > 0"
                    class="px-2 inline-flex text-xs leading-5 font-semibold rounded-full bg-gray-100 text-gray-800"
                  >
                    {{ friendly_ago(ws.archivedAt) }}
                  </span>
                  <span
                    v-else
                    class="px-2 inline-flex text-xs leading-5 font-semibold rounded-full bg-green-100 text-green-800"
                  >
                    Restored
                  </span>
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                  <span v-if="ws.unarchivedAt">Open</span>
                  <a
                    v-else
                    href="#"
                    class="text-blue-600 hover:text-blue-900"
                    @click.stop.prevent="restore(ws.id)"
                  >
                    Restore
                  </a>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { gql, useMutation, useQuery } from '@urql/vue'
import time from '../../time'
import Avatar from '../shared/Avatar.vue'
import { defineComponent } from 'vue'
import {
  RecentlyArchivedWorkspaceFragment,
  RecentlyArchivedWorkspacesQuery,
  RecentlyArchivedWorkspacesQueryVariables,
} from './__generated__/RecentlyArchivedWorkspaces'
import EmptyPlaceholder from '../../molecules/EmptyPlaceholder.vue'

export const RECENTLY_ARCHIVED_WORKSPACE_FRAGMENT = gql`
  fragment RecentlyArchivedWorkspace on Workspace {
    id
    archivedAt
    unarchivedAt
  }
`

export default defineComponent({
  name: 'RecentlyArchivedWorkspaces',
  components: { EmptyPlaceholder, Avatar },
  props: {
    codebaseId: { type: String, required: true },
  },
  setup(props) {
    let { data } = useQuery<
      RecentlyArchivedWorkspacesQuery,
      RecentlyArchivedWorkspacesQueryVariables
    >({
      query: gql`
        query RecentlyArchivedWorkspaces($codebaseID: ID!) {
          workspaces(codebaseID: $codebaseID, includeArchived: true) {
            id
            name
            author {
              id
              name
              avatarUrl
            }
            ...RecentlyArchivedWorkspace
          }
        }
        ${RECENTLY_ARCHIVED_WORKSPACE_FRAGMENT}
      `,
      variables: {
        codebaseID: props.codebaseId,
      },
    })

    const { executeMutation: unarchiveWorkspaceResult } = useMutation(gql`
      mutation RecentlyArchivedWorkspacesUnarchive($id: ID!) {
        unarchiveWorkspace(id: $id) {
          id
          archivedAt
          unarchivedAt
          view {
            id
          }
        }
      }
    `)

    return {
      data,

      async unarchiveWorkspace(id: string) {
        const variables = { id }
        await unarchiveWorkspaceResult(variables).then((result) => {
          if (result.error) {
            throw new Error(result.error)
          }
        })
      },
    }
  },
  data() {
    return { preRestoreArchivedAt: new Map<string, number>() }
  },
  computed: {
    sortedWorkspaces() {
      if (this.data?.workspaces) {
        let ws = this.data.workspaces.filter((w) => w.archivedAt || w.unarchivedAt)
        ws.sort((a, b) => this.sortVal(b) - this.sortVal(a))
        return ws
      }
      return []
    },
  },
  methods: {
    friendly_ago(ts: number) {
      return time.getRelativeTime(new Date(ts * 1000))
    },
    restore(id: string) {
      let preArchivedAt = this.data?.workspaces.filter((w) => w.id === id)[0].archivedAt
      if (preArchivedAt) {
        this.preRestoreArchivedAt.set(id, preArchivedAt)
      }
      this.unarchiveWorkspace(id)
    },
    sortVal(ws: RecentlyArchivedWorkspaceFragment): number {
      let p = this.preRestoreArchivedAt.get(ws.id)
      if (p) {
        return p
      }

      return Math.max(ws.archivedAt || 0, ws.unarchivedAt || 0)
    },
    maybeGoto(ws: RecentlyArchivedWorkspaceFragment) {
      if (ws.archivedAt) {
        return
      }

      this.$router.push({
        name: 'workspaceHome',
        params: { id: ws.id },
      })
    },
  },
})
</script>
