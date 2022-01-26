<template>
  <PaddedAppLeftSidebar v-if="data" class="bg-white">
    <template #navigation>
      <VerticalNavigation />
    </template>

    <template #header>
      <OrganizationSettingsHeader :name="data.organization.name" />
    </template>

    <template #default>
      <CreateCodebase
        v-if="data"
        :create-in-organization-id="data.organization.id"
        bottom-bg="bg-gray-100"
        main-bg="bg-gray-100"
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
    </template>
  </PaddedAppLeftSidebar>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import CreateCodebase from '../../organisms/CreateCodebase.vue'
import { useRoute } from 'vue-router'
import { gql, useQuery } from '@urql/vue'
import {
  CreateOrganizationCodebasePageQuery,
  CreateOrganizationCodebasePageQueryVariables,
} from './__generated__/CreateOrganizationCodebasePage'
import VerticalNavigation from '../../organisms/organization/VerticalNavigation.vue'
import OrganizationSettingsHeader from '../../organisms/organization/OrganizationSettingsHeader.vue'
import PaddedAppLeftSidebar from '../../layouts/PaddedAppLeftSidebar.vue'

export default defineComponent({
  components: {
    PaddedAppLeftSidebar,
    CreateCodebase,
    VerticalNavigation,
    OrganizationSettingsHeader,
  },
  setup() {
    let route = useRoute()

    let { data } = useQuery<
      CreateOrganizationCodebasePageQuery,
      CreateOrganizationCodebasePageQueryVariables
    >({
      query: gql`
        query CreateOrganizationCodebasePage($shortID: ID!) {
          organization(shortID: $shortID) {
            id
            name
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
