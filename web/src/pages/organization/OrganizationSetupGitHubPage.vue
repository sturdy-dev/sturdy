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
          <span>Sturdy for GitHub</span>
        </Header>

        <OrganizationSetupGitHub
          :organization="data.organization"
          :git-hub-app="data.gitHubApp"
          :git-hub-account="data.user.gitHubAccount"
        />
      </div>
    </template>
  </PaddedAppLeftSidebar>
</template>

<script lang="ts">
import { computed, defineComponent, inject, ref } from 'vue'
import type { Ref } from 'vue'
import { gql, useQuery } from '@urql/vue'
import type {
  OrganizationSetupGitHubPageQuery,
  OrganizationSetupGitHubPageQueryVariables,
} from './__generated__/OrganizationSetupGitHubPage'
import { useRoute } from 'vue-router'
import Header from '../../molecules/Header.vue'
import PaddedAppLeftSidebar from '../../layouts/PaddedAppLeftSidebar.vue'
import VerticalNavigation from '../../organisms/organization/VerticalNavigation.vue'
import OrganizationSetupGitHub, {
  ORGANIZATION_SETUP_GITHUB_GITHUB_APP_FRAGMENT,
} from '../../organisms/organization/OrganizationSetupGitHub.vue'
import OrganizationSettingsHeader from '../../organisms/organization/OrganizationSettingsHeader.vue'
import { Feature } from '../../__generated__/types'

export default defineComponent({
  components: {
    OrganizationSetupGitHub,
    OrganizationSettingsHeader,
    PaddedAppLeftSidebar,
    Header,
    VerticalNavigation,
  },
  setup() {
    let route = useRoute()

    const features = inject<Ref<Array<Feature>>>('features', ref([]))
    const isGitHubEnabled = computed(() => features?.value?.includes(Feature.GitHub))

    let { data } = useQuery<
      OrganizationSetupGitHubPageQuery,
      OrganizationSetupGitHubPageQueryVariables
    >({
      query: gql`
        query OrganizationSetupGitHubPage($shortID: ID!, $isGitHubEnabled: Boolean!) {
          organization(shortID: $shortID) {
            id
            name
          }

          gitHubApp @include(if: $isGitHubEnabled) {
            _id
            clientID
            name
            ...OrganizationSetupGitHub_GitHubApp
          }

          user {
            id
            gitHubAccount @include(if: $isGitHubEnabled) {
              id
              login
              isValid
            }
          }
        }

        ${ORGANIZATION_SETUP_GITHUB_GITHUB_APP_FRAGMENT}
      `,
      requestPolicy: 'cache-and-network',
      variables: {
        isGitHubEnabled,
        shortID: route.params.organizationSlug,
      },
    })

    return {
      data,
    }
  },
})
</script>
