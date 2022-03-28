<template>
  <TitleBar>
    <template #header>
      <h1>Preferences</h1>
    </template>

    <template #default>
      <div class="h-full flex gap-4">
        <nav class="min-w-max flex flex-col bg-gray-50 gap-2 p-2 border-r border-gray-100">
          <span
            v-for="page in pages"
            :key="page.id"
            :href="`#${page.id}`"
            class="flex flex-inline font-medium hover:bg-gray-200 rounded-md p-2 cursor-pointer"
            @click="() => (this.selected = page)"
          >
            {{ page.title }}
          </span>
        </nav>

        <div class="flex p-2">
          <component :is="selected.component" />
        </div>
      </div>
    </template>
  </TitleBar>
</template>

<script lang="ts">
import TitleBar from './components/TitleBar.vue'
import ServersList from './components/ServersList.vue'
import SoftwareUpdates from './components/SoftwareUpdates.vue'
import { ref } from 'vue'

const pages = [
  { id: 'servers', title: 'Servers', component: ServersList },
  { id: 'software-updates', title: 'Software Updates', component: SoftwareUpdates },
]

export default {
  components: {
    TitleBar,
    ServersList,
    SoftwareUpdates,
  },
  setup() {
    const selected = ref(pages[0])
    return {
      pages,
      selected,
    }
  },
}
</script>
