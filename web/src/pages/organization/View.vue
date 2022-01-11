<template>
  <PaddedApp>
    <pre>{{ data }}</pre>
  </PaddedApp>
</template>

<script lang="ts">
import PaddedApp from '../../layouts/PaddedApp.vue'
import { defineComponent } from 'vue'
import { gql, useQuery } from '@urql/vue'
import { OrganizationQuery, OrganizationQueryVariables } from './__generated__/View'
import { useRoute } from 'vue-router'

export default defineComponent({
  components: { PaddedApp },
  setup() {
    let route = useRoute()
    let orgID = route.params.id as string

    let { data } = useQuery<OrganizationQuery, OrganizationQueryVariables>({
      query: gql`
        query Organization($id: ID!) {
          organization(id: $id) {
            id
            name
            members {
              id
              name
              email
              avatarUrl
            }
            codebases {
              id
              shortID
              name
            }
          }
        }
      `,
      requestPolicy: 'cache-and-network',
      variables: {
        id: orgID,
      },
    })

    return {
      data,
    }
  },
})
</script>
