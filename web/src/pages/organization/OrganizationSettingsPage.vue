<template>
  <PaddedAppLeftSidebar v-if="data" class="bg-white">
    <template #navigation>
      <VerticalNavigation />
    </template>

    <template #default>
      <div class="max-w-7xl">
        <Header>
          <span>Manage {{ data.organization.name }}</span>
        </Header>

        <OrganizationMembers
          :members="data.organization.members"
          :organization-id="data.organization.id"
        />
      </div>
    </template>
  </PaddedAppLeftSidebar>
</template>

<script lang="ts">
import { computed, defineComponent } from 'vue'
import { gql, useQuery } from '@urql/vue'
import {
  OrganizationSettingsPageQuery,
  OrganizationSettingsPageQueryVariables,
} from './__generated__/View'
import { useRoute } from 'vue-router'
import OrganizationMembers from '../../organisms/organization/OrganizationMembers.vue'
import Header from '../../molecules/Header.vue'
import PaddedAppLeftSidebar from '../../layouts/PaddedAppLeftSidebar.vue'
import VerticalNavigation from '../../organisms/organization/VerticalNavigation.vue'

export default defineComponent({
  components: { PaddedAppLeftSidebar, OrganizationMembers, Header, VerticalNavigation },
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
            codebases {
              id
              shortID
              name
            }
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
