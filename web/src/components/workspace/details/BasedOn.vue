<template>
  <div class="space-y-5 md:flex-1">
    <router-link
      v-if="workspace.headChange"
      :to="{
        name: 'codebaseChange',
        params: {
          codebaseSlug: codebaseSlug,
          id: workspace.headChange.id,
        },
      }"
      class="flex items-center space-x-2"
    >
      <Avatar size="5" :author="workspace.headChange.author" class="flex-shrink-0" />
      <span class="text-gray-900 text-sm font-medium">
        <span class="text-gray-900 text-sm font-medium">
          Based on
          <span class="text-gray-500 font-medium">
            {{ workspace.headChange.title }}
          </span>
        </span>
      </span>
    </router-link>
  </div>
</template>

<script>
import Avatar from '../../shared/Avatar.vue'

export default {
  components: { Avatar },
  props: ['workspace', 'codebaseSlug'],
}
</script>
