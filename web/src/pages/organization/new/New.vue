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

    <CreateCodebase v-else />
  </PaddedApp>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import { useRoute } from 'vue-router'
import { gql, useQuery } from '@urql/vue'
import type { NewCodebasePageQuery, NewCodebasePageQueryVariables } from './__generated__/New'
import OrganizationSettingsHeader from '../../../organisms/organization/OrganizationSettingsHeader.vue'
import PaddedApp from '../../../layouts/PaddedApp.vue'
import type { DeepMaybeRef } from '@vueuse/core'
import CreateCodebase from '../../../organisms/CreateCodebase.vue'

const PAGE_QUERY = gql`
  query NewCodebasePage($shortID: ID!) {
    organization(shortID: $shortID) {
      id
      name
      writeable
    }
  }
`

export default defineComponent({
  components: {
    PaddedApp,
    CreateCodebase,
    OrganizationSettingsHeader,
  },
  setup() {
    const route = useRoute()

    const { data } = useQuery<NewCodebasePageQuery, DeepMaybeRef<NewCodebasePageQueryVariables>>({
      query: PAGE_QUERY,
      requestPolicy: 'cache-and-network',
      variables: {
        shortID: route.params.organizationSlug as string,
      },
    })

    return {
      data,
    }
  },
})
</script>
