<template>
  <PaddedApp v-if="data" class="bg-white">
    <OrganizationCreate class="max-w-3xl" :with-tier-picker="false" :user="data.user" />
  </PaddedApp>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import PaddedApp from '../../layouts/PaddedApp.vue'
import OrganizationCreate from '../../organisms/organization/OrganizationCreate.vue'
import { gql, useQuery } from '@urql/vue'

export default defineComponent({
  components: { PaddedApp, OrganizationCreate },
  setup() {
    let { data } = useQuery({
      query: gql`
        query CreateOrganizationPage {
          user {
            id
            name
          }
        }
      `,
    })

    return { data }
  },
})
</script>
