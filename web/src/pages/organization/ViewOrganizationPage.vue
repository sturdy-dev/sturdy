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
import PaddedApp from '../../layouts/PaddedApp.vue'
import { defineComponent } from 'vue'
import { gql, useQuery } from '@urql/vue'
import { OrganizationQuery, OrganizationQueryVariables } from './__generated__/View'
import { useRoute } from 'vue-router'
import OrganizationMembers from '../../organisms/organization/OrganizationMembers.vue'
import Header from '../../molecules/Header.vue'
import PaddedAppLeftSidebar from '../../layouts/PaddedAppLeftSidebar.vue'
import VerticalNavigation from '../../organisms/organization/VerticalNavigation.vue'

export default defineComponent({
  components: { PaddedAppLeftSidebar, OrganizationMembers, Header, VerticalNavigation },
  setup() {
    let route = useRoute()
    let orgID = route.params.id as string

    let { data } = useQuery<OrganizationQuery, OrganizationQueryVariables>({
      query: gql`
        query Organization($id: ID!) {
          organization(id: $id) {
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
        id: orgID,
      },
    })

    return {
      data,
    }
  },
})
</script>
