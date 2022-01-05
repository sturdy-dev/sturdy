<template>
  <div class="flex items-center space-x-2">
    <ChatAltIcon class="h-5 w-5 text-gray-400" aria-hidden="true" />
    <span v-if="data.workspace.comments.length === 0" class="text-gray-900 text-sm font-medium">
      No comments
    </span>
    <span v-else class="text-gray-900 text-sm font-medium">
      {{ data.workspace.comments.length }}
      {{ data.workspace.comments.length === 1 ? 'comment' : 'comments' }}
    </span>
  </div>
</template>

<script>
import { ChatAltIcon } from '@heroicons/vue/solid'

export default {
  components: { ChatAltIcon },
  props: ['data'],
}
</script>
