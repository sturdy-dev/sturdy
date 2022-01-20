<template>
  <PaddedApp>
    <CreateCodebase v-if="data" :create-in-organization-id="data.organization.id">
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
import { defineComponent } from 'vue'
import PaddedApp from '../../layouts/PaddedApp.vue'
import CreateCodebase from '../../organisms/CreateCodebase.vue'
import { useRoute } from 'vue-router'
import { gql, useQuery } from '@urql/vue'
import {
  CreateOrganizationCodebasePageQuery,
  CreateOrganizationCodebasePageQueryVariables,
} from './__generated__/CreateOrganizationCodebasePage'

export default defineComponent({
  components: { PaddedApp, CreateCodebase },
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
