<template>
  <div class="flex items-center space-x-2">
    <ChatAltIcon class="h-5 w-5 text-gray-400" aria-hidden="true" />
    <span v-if="commentsCount === 0" class="text-gray-900 text-sm font-medium">No comments</span>
    <span v-else class="text-gray-900 text-sm font-medium">
      {{ commentsCount }}
      {{ commentsCount === 1 ? 'comment' : 'comments' }}
    </span>
  </div>
</template>

<script lang="ts">
import { ChatAltIcon } from '@heroicons/vue/solid'
import { gql } from '@urql/vue'
import { CommentsCount_WorkspaceFragment } from './__generated__/Comments'
import { defineComponent, PropType } from 'vue'

export const WORKSPACE_FRAGMENT = gql`
  fragment CommentsCount_Workspace on Workspace {
    id
    commentsCount
  }
`

export default defineComponent({
  components: { ChatAltIcon },
  props: {
    workspace: {
      type: Object as PropType<CommentsCount_WorkspaceFragment>,
      required: true,
    },
  },
  computed: {
    commentsCount() {
      return this.workspace.commentsCount
    },
  },
})
</script>
