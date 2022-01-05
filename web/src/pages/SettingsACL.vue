<template>
  <PaddedAppLeftSidebar v-if="data" class="bg-white">
    <template #navigation>
      <SettingsVerticalNavigation />
    </template>

    <template #default>
      <div class="max-w-7xl">
        <Header>Security and ACLs</Header>
        <SettingsACL
          v-if="data?.codebase?.acl?.id"
          :codebase-id="data.codebase.id"
          :acl-id="data.codebase.acl.id"
          :acl-policy="data.codebase.acl.policy"
        />
      </div>
    </template>
  </PaddedAppLeftSidebar>
</template>

<script>
import { gql, useQuery } from '@urql/vue'
import { useRoute } from 'vue-router'
import PaddedAppLeftSidebar from '../layouts/PaddedAppLeftSidebar.vue'
import SettingsVerticalNavigation from '../components/codebase/settings/SettingsVerticalNavigation.vue'
import Header from '../molecules/Header.vue'
import SettingsACL from '../components/codebase/settings/SettingsACL.vue'

export default {
  name: 'CodebaseSettings',
  components: {
    SettingsACL,
    PaddedAppLeftSidebar,
    SettingsVerticalNavigation,
    Header,
  },
  setup() {
    let route = useRoute()

    let { data } = useQuery({
      query: gql`
        query SettingsAclPage($id: ID, $shortID: ID) {
          codebase(id: $id, shortID: $shortID) {
            id
            acl {
              id
              policy
            }
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
