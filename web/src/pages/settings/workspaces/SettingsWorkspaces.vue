<template>
  <PaddedAppLeftSidebar v-if="data" class="bg-white">
    <template #navigation>
      <SettingsVerticalNavigation />
    </template>

    <template #default>
      <div class="max-w-7xl space-y-4">
        <Header>Restore Workspaces</Header>
        <SettingsWorkspaces classs="mt-8" :codebase-id="data.codebase.id" />
      </div>
    </template>
  </PaddedAppLeftSidebar>
</template>

<script>
import { gql, useQuery } from '@urql/vue'
import { useRoute } from 'vue-router'
import PaddedAppLeftSidebar from '../../../layouts/PaddedAppLeftSidebar.vue'
import SettingsVerticalNavigation from '../../../components/codebase/settings/SettingsVerticalNavigation.vue'
import Header from '../../../molecules/Header.vue'
import SettingsWorkspaces from '../../../components/codebase/settings/SettingsWorkspaces.vue'

export default {
  name: 'CodebaseSettings',
  components: {
    SettingsWorkspaces,
    PaddedAppLeftSidebar,
    SettingsVerticalNavigation,
    Header,
  },
  setup() {
    let route = useRoute()

    let { data } = useQuery({
      query: gql`
        query SettingsWorkspaces($id: ID, $shortID: ID) {
          codebase(id: $id, shortID: $shortID) {
            id
          }
        }
      `,
      variables: {
        shortID: route.params.codebaseSlug,
      },
    })

    return {
      data,
    }
  },
}
</script>
