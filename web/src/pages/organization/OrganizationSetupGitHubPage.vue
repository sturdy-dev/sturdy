<template>
  <PaddedAppLeftSidebar v-if="data" class="bg-white">
    <template #navigation>
      <VerticalNavigation />
    </template>

    <template #default>
      <div class="max-w-7xl">
        <Header>
          <span>Setup codebase from GitHub in {{ data.organization.name }}</span>
        </Header>

        <OrganizationSetupGitHub
          :organization-id="data.organization.id"
          :git-hub-app="data.gitHubApp"
          :git-hub-account="data.user.gitHubAccount"
        />
      </div>
    </template>
  </PaddedAppLeftSidebar>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import { gql, useQuery } from '@urql/vue'
import {
  OrganizationSetupGitHubPageQuery,
  OrganizationSetupGitHubPageQueryVariables,
} from './__generated__/OrganizationSetupGitHubPage'
import { useRoute } from 'vue-router'
import Header from '../../molecules/Header.vue'
import PaddedAppLeftSidebar from '../../layouts/PaddedAppLeftSidebar.vue'
import VerticalNavigation from '../../organisms/organization/VerticalNavigation.vue'
import OrganizationSetupGitHub from '../../organisms/organization/OrganizationSetupGitHub.vue'

export default defineComponent({
  components: {
    OrganizationSetupGitHub,
    PaddedAppLeftSidebar,
    Header,
    VerticalNavigation,
  },
  setup() {
    let route = useRoute()

    let { data } = useQuery<
      OrganizationSetupGitHubPageQuery,
      OrganizationSetupGitHubPageQueryVariables
    >({
      query: gql`
        query OrganizationSetupGitHubPage($shortID: ID!) {
          organization(shortID: $shortID) {
            id
            name
          }

          gitHubApp {
            _id
            clientID
            name
          }

          user {
            id
            gitHubAccount {
              id
              login
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
