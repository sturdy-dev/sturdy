<template>
  <RecentlyArchivedWorkspaces :codebase-id="codebaseId" />
</template>

<script lang="ts">
import RecentlyArchivedWorkspaces from '../RecentlyArchivedWorkspaces.vue'
import { defineComponent } from 'vue'

export default defineComponent({
  name: 'SettingsWorkspaces',
  components: {
    RecentlyArchivedWorkspaces,
  },
  props: {
    codebaseId: {
      type: String,
      required: true,
    },
  },
})
</script>
