<template>
  <PaddedApp v-if="data" class="bg-white">
    <div v-if="!selected?.writeable">
      <p class="text-sm text-gray-500">
        You don't have permissions to create a new codebase in this organization. Ask an
        administrator for help if you want to create a new codebase in
        <strong>{{ selected?.name }}</strong
        >.
      </p>
    </div>

    <CreateCodebase v-else :selected-organization="selected" :organizations="data.organizations" />
  </PaddedApp>
</template>

<script lang="ts" setup>
import { computed } from 'vue'
import { gql, useQuery } from '@urql/vue'
import type { NewCodebasePageQuery, NewCodebasePageQueryVariables } from './__generated__/Index'
import PaddedApp from '../../../layouts/PaddedApp.vue'
import type { DeepMaybeRef } from '@vueuse/core'
import CreateCodebase, { ORGANIZATION_FRAGMENT } from '../../../organisms/CreateCodebase.vue'
import { useRoute } from 'vue-router'

const PAGE_QUERY = gql`
  query NewCodebasePage {
    organizations {
      id
      name
      writeable
      shortID
      ...Organization_CreateCodebase
    }
  }
  ${ORGANIZATION_FRAGMENT}
`

const { data } = useQuery<NewCodebasePageQuery, DeepMaybeRef<NewCodebasePageQueryVariables>>({
  query: PAGE_QUERY,
  requestPolicy: 'cache-and-network',
})

const route = useRoute()

const selected = computed(
  () =>
    data.value?.organizations.find(({ shortID }) => shortID === route.params.organizationSlug) ||
    data.value?.organizations[0]
)
</script>
