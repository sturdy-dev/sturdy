<template>
  <PaddedAppLeftSidebar v-if="data" class="bg-white">
    <template #navigation>
      <SettingsVerticalNavigation />
    </template>

    <template #default>
      <div class="max-w-7xl">
        <Header>Team and collaborators</Header>
        <SettingsCollaborators :codebase-id="data.codebase.id" />
      </div>
    </template>
  </PaddedAppLeftSidebar>
</template>

<script>
import SettingsCollaborators from '../components/codebase/settings/SettingsCollaborators.vue'
import { gql, useQuery } from '@urql/vue'
import { useRoute } from 'vue-router'
import PaddedAppLeftSidebar from '../layouts/PaddedAppLeftSidebar.vue'
import SettingsVerticalNavigation from '../components/codebase/settings/SettingsVerticalNavigation.vue'
import Header from '../molecules/Header.vue'

export default {
  name: 'CodebaseSettings',
  components: {
    PaddedAppLeftSidebar,
    SettingsVerticalNavigation,
    SettingsCollaborators,
    Header,
  },
  setup() {
    let route = useRoute()

    let { data } = useQuery({
      query: gql`
        query SettingsTeam($id: ID, $shortID: ID) {
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
