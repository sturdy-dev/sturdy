<template>
  <PaddedAppLeftSidebar v-if="data" class="bg-white">
    <template #navigation>
      <VerticalNavigation />
    </template>

    <template #header>
      <OrganizationSettingsHeader :name="data.organization.name" />
    </template>

    <template #default>
      <div v-if="!data.organization.writeable">
        <p class="text-sm text-gray-500">
          You don't have permissions to create a new codebase in this organization. Ask an
          administrator for help if you want to create a new codebase in
          <strong>{{ data.organization.name }}</strong
          >.
        </p>
      </div>

      <CreateCodebase
        v-if="data && data.organization.writeable"
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
import CreateCodebase from '../../organisms/CreateCodebaseOnSturdy.vue'
import { useRoute } from 'vue-router'
import { gql, useQuery } from '@urql/vue'
import type {
  CreateOrganizationCodebasePageQuery,
  CreateOrganizationCodebasePageQueryVariables,
} from './__generated__/CreateOrganizationCodebasePage'
import VerticalNavigation from '../../organisms/organization/VerticalNavigation.vue'
import OrganizationSettingsHeader from '../../organisms/organization/OrganizationSettingsHeader.vue'
import PaddedAppLeftSidebar from '../../layouts/PaddedAppLeftSidebar.vue'
import type { DeepMaybeRef } from '@vueuse/core'

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
      DeepMaybeRef<CreateOrganizationCodebasePageQueryVariables>
    >({
      query: gql`
        query CreateOrganizationCodebasePage($shortID: ID!) {
          organization(shortID: $shortID) {
            id
            name
            writeable
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
