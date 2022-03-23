<template>
  <nav class="space-y-1" aria-label="Sidebar">
    <router-link
      v-for="item in navigation"
      :key="item.name"
      :to="{ name: item.linkName }"
      :class="[
        item.current
          ? 'bg-gray-100 text-gray-900'
          : 'text-gray-600 hover:bg-gray-50 hover:text-gray-900',
        'flex items-center px-3 py-2 text-sm font-medium rounded-md',
      ]"
      :aria-current="item.current ? 'page' : undefined"
    >
      <span class="truncate">
        {{ item.name }}
      </span>
    </router-link>
  </nav>
</template>

<script lang="ts">
import type { PropType } from 'vue'
import { defineComponent } from 'vue'

interface NavigationItem {
  name: string
  linkName: string
  current: boolean
}

export default defineComponent({
  props: {
    navigation: {
      type: Object as PropType<Array<NavigationItem>>,
      required: true,
    },
  },
})
</script>
