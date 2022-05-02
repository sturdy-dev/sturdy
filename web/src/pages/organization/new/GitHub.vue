<template>
  <PaddedApp v-if="data" class="bg-white">
    <OrganizationSettingsHeader :name="data.organization.name" />
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
  </PaddedApp>
</template>

<script lang="ts">
import { computed, defineComponent, inject, ref } from 'vue'
import type { Ref } from 'vue'
import { gql, useQuery } from '@urql/vue'
import type {
  OrganizationSetupGitHubPageQuery,
  OrganizationSetupGitHubPageQueryVariables,
} from './__generated__/GitHub'
import { useRoute } from 'vue-router'
import Header from '../../../molecules/Header.vue'
import PaddedApp from '../../../layouts/PaddedApp.vue'
import OrganizationSetupGitHub, {
  ORGANIZATION_SETUP_GITHUB_GITHUB_APP_FRAGMENT,
  ORGANIZATION_SETUP_GITHUB_ORGANIZATION_FRAGMENT,
  ORGANIZATION_SETUP_GITHUB_GITHUB_ACCOUNT_FRAGMENT,
} from '../../../organisms/CreateCodebaseFromGitHub.vue'
import OrganizationSettingsHeader from '../../../organisms/organization/OrganizationSettingsHeader.vue'
import { Feature } from '../../../__generated__/types'
import type { DeepMaybeRef } from '@vueuse/core'

const PAGE_QUERY = gql`
  query OrganizationSetupGitHubPage($shortID: ID!, $isGitHubEnabled: Boolean!) {
    organization(shortID: $shortID) {
      id
      name
      ...OrganizationSetupGitHub_Organization
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
        ...OrganizationSetupGitHub_GitHubAccount
      }
    }
  }

  ${ORGANIZATION_SETUP_GITHUB_GITHUB_ACCOUNT_FRAGMENT}
  ${ORGANIZATION_SETUP_GITHUB_GITHUB_APP_FRAGMENT}
  ${ORGANIZATION_SETUP_GITHUB_ORGANIZATION_FRAGMENT}
`

export default defineComponent({
  components: {
    OrganizationSetupGitHub,
    OrganizationSettingsHeader,
    PaddedApp,
    Header,
  },
  setup() {
    const route = useRoute()

    const features = inject<Ref<Array<Feature>>>('features', ref([]))
    const isGitHubEnabled = computed(() => features?.value?.includes(Feature.GitHub) ?? false)

    const { data } = useQuery<
      OrganizationSetupGitHubPageQuery,
      DeepMaybeRef<OrganizationSetupGitHubPageQueryVariables>
    >({
      query: PAGE_QUERY,
      requestPolicy: 'cache-and-network',
      variables: {
        isGitHubEnabled,
        shortID: route.params.organizationSlug as string,
      },
    })

    return {
      data,
    }
  },
})
</script>
