<template>
  <div class="text-center">
    <div class="inline-block text-blue-800 bg-blue-50 rounded-full w-32 h-32 p-7">
      <FolderIcon v-if="workspace.view == null" />
      <FolderOpenIcon v-else />
    </div>

    <div class="text-gray-500 m-auto text-sm max-w-sm mt-2">
      <p v-if="workspace.view == null">
        You haven't made any changes yet. Connect a local directory to this draft change to
        <span class="whitespace-nowrap">get started!</span>
      </p>

      <p v-else>
        You haven't made any changes yet. Start coding in
        <span class="text-black font-medium">{{ workspace.view.shortMountPath }}</span
        >, and the changes will appear here.
      </p>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent, PropType } from 'vue'
import gql from 'graphql-tag'
import { NoChangesOwnWorkspaceFragment } from './__generated__/NoChangesOwnWorkspace'
import { FolderIcon, FolderOpenIcon } from '@heroicons/vue/outline'

export const NO_CHANGES_OWN_WORKSPACE = gql`
  fragment NoChangesOwnWorkspace on Workspace {
    id
    view {
      id
      shortMountPath
    }
  }
`

export default defineComponent({
  components: {
    FolderIcon,
    FolderOpenIcon,
  },
  props: {
    workspace: {
      type: Object as PropType<NoChangesOwnWorkspaceFragment>,
      required: true,
    },
  },
})
</script>
