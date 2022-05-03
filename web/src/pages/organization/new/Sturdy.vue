<template>
  <PaddedApp v-if="data" class="bg-white">
    <OrganizationSettingsHeader :name="data.organization.name" />

    <div v-if="!data.organization.writeable">
      <p class="text-sm text-gray-500">
        You don't have permissions to create a new codebase in this organization. Ask an
        administrator for help if you want to create a new codebase in
        <strong>{{ data.organization.name }}</strong
        >.
      </p>
    </div>

    <CreateCodebase
      v-if="data.organization.writeable"
      :create-in-organization-id="data.organization.id"
      bottom-bg="bg-gray-100"
      main-bg="bg-gray-100"
      :user="data.user"
      :git-hub-app="data.gitHubApp"
    >
      <template #header>
        <div class="py-8 px-4">
          <div class="">
            <h2
              class="text-4xl font-extrabold text-gray-900 sm:text-4xl sm:tracking-tight lg:text-4xl"
            >
              Create a new codebase in <span class="underline">{{ data.organization.name }}</span>
            </h2>
            <p class="mt-5 text-xl text-gray-500">You'll soon be ready to code! 📈</p>
          </div>
        </div>
      </template>
    </CreateCodebase>
  </PaddedApp>
</template>

<script lang="ts">
import { defineComponent, inject, ref, type Ref } from 'vue'
import CreateCodebase, {
  USER_FRAGMENT,
  GITHUB_APP_FRAGMENT,
} from '../../../organisms/CreateCodebaseOnSturdy.vue'
import { useRoute } from 'vue-router'
import { gql, useQuery } from '@urql/vue'
import type {
  CreateOrganizationCodebasePageQuery,
  CreateOrganizationCodebasePageQueryVariables,
} from './__generated__/Sturdy'
import OrganizationSettingsHeader from '../../../organisms/organization/OrganizationSettingsHeader.vue'
import PaddedApp from '../../../layouts/PaddedApp.vue'
import type { DeepMaybeRef } from '@vueuse/core'
import { Feature } from '../../../__generated__/types'

const PAGE_QUERY = gql`
  query CreateOrganizationCodebasePage($shortID: ID!, $isGitHubEnabled: Boolean!) {
    organization(shortID: $shortID) {
      id
      name
      writeable
    }

    user {
      id
      ...User_CreateCodebaseOnSturdy
    }

    gitHubApp @include(if: $isGitHubEnabled) {
      _id
      ...GitHubApp_CreateCodebaseOnSturdy
    }
  }
  ${USER_FRAGMENT}
  ${GITHUB_APP_FRAGMENT}
`

export default defineComponent({
  components: {
    PaddedApp,
    CreateCodebase,
    OrganizationSettingsHeader,
  },
  setup() {
    const route = useRoute()
    const features = inject<Ref<Array<Feature>>>('features', ref([]))
    const isGitHubEnabled = features.value.includes(Feature.GitHub)

    const { data } = useQuery<
      CreateOrganizationCodebasePageQuery,
      DeepMaybeRef<CreateOrganizationCodebasePageQueryVariables>
    >({
      query: PAGE_QUERY,
      requestPolicy: 'cache-and-network',
      variables: {
        shortID: route.params.organizationSlug as string,
        isGitHubEnabled,
      },
    })

    return {
      data,
    }
  },
})
</script>
