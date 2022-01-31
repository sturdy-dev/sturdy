<template>
  <PaddedAppLeftSidebar v-if="data" class="bg-white">
    <template #navigation>
      <VerticalNavigation />
    </template>

    <template #header>
      <OrganizationSettingsHeader :name="data.organization.name" />
    </template>

    <template #default>
      <div class="max-w-7xl">
        <Header>
          <span>Settings</span>
        </Header>

        <OrganizationMembers :organization="data.organization" :user="data.user" />
      </div>
    </template>
  </PaddedAppLeftSidebar>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import { gql, useQuery } from '@urql/vue'
import {
  OrganizationSettingsPageQuery,
  OrganizationSettingsPageQueryVariables,
} from './__generated__/OrganizationSettingsPage'
import { useRoute } from 'vue-router'
import OrganizationMembers from '../../organisms/organization/OrganizationMembers.vue'
import Header from '../../molecules/Header.vue'
import PaddedAppLeftSidebar from '../../layouts/PaddedAppLeftSidebar.vue'
import VerticalNavigation from '../../organisms/organization/VerticalNavigation.vue'
import OrganizationSettingsHeader from '../../organisms/organization/OrganizationSettingsHeader.vue'

export default defineComponent({
  components: {
    PaddedAppLeftSidebar,
    OrganizationMembers,
    Header,
    VerticalNavigation,
    OrganizationSettingsHeader,
  },
  setup() {
    let route = useRoute()

    let { data } = useQuery<OrganizationSettingsPageQuery, OrganizationSettingsPageQueryVariables>({
      query: gql`
        query OrganizationSettingsPage($shortID: ID!) {
          organization(shortID: $shortID) {
            id
            name
            members {
              id
              name
              email
              avatarUrl
            }
            writeable
          }

          user {
            id
          }
        }
      `,
      requestPolicy: 'cache-and-network',
      variables: {
        shortID: route.params.organizationSlug,
      },
    })

    return {
      data,
    }
  },
})
</script>
