<template>
  <div class="text-center">
    <div class="inline-block text-blue-800 bg-blue-50 rounded-full w-32 h-32 p-7">
      <ChatAlt2Icon />
    </div>

    <p class="text-gray-500 m-auto text-sm max-w-sm mt-2">
      {{ workspace.author.name }} haven't made any changes in this
      <span class="whitespace-nowrap">workspace yet.</span>
    </p>
  </div>
</template>

<script lang="ts">
import { defineComponent, PropType } from 'vue'
import gql from 'graphql-tag'
import { NoChangesOthersWorkspaceFragment } from './__generated__/NoChangesOthersWorkspace'
import { ChatAlt2Icon } from '@heroicons/vue/outline'

export const NO_CHANGES_OTHERS_WORKSPACE = gql`
  fragment NoChangesOthersWorkspace on Workspace {
    id
    author {
      id
      name
    }
  }
`

export default defineComponent({
  components: { ChatAlt2Icon },
  props: {
    workspace: {
      type: Object as PropType<NoChangesOthersWorkspaceFragment>,
      required: true,
    },
  },
})
</script>
