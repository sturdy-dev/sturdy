<template>
  <PaddedAppLeftSidebar v-if="data" class="bg-white">
    <template #navigation>
      <VerticalNavigation />
    </template>

    <template #header>
      <OrganizationSettingsHeader :name="data.organization.name" />
    </template>

    <template #default>
      <Header>
        <span>Subscriptions</span>
      </Header>

      <template v-if="data.organization.licenses && data.organization.licenses.length > 0">
        <p class="my-2 text-sm">
          A subscription allows you to self-host Sturdy Enterprise. Reach out to
          support@getsturdy.com to make updates to your subscriptions.
        </p>
        <OrganizationListLicenses :licenses="data.organization.licenses" />
      </template>
      <div v-else>
        <p class="my-2 text-sm">
          A subscription allows you to self-host Sturdy Enterprise. Reach out to
          support@getsturdy.com to get started.
        </p>
      </div>
    </template>
  </PaddedAppLeftSidebar>
</template>

<script lang="ts">
import { computed, defineComponent } from 'vue'
import Header from '../../molecules/Header.vue'
import VerticalNavigation from '../../organisms/organization/VerticalNavigation.vue'
import PaddedAppLeftSidebar from '../../layouts/PaddedAppLeftSidebar.vue'
import { gql, useQuery } from '@urql/vue'
import { useRoute } from 'vue-router'
import OrganizationSettingsHeader from '../../organisms/organization/OrganizationSettingsHeader.vue'
import type {
  ListSubscriptionsPageQuery,
  ListSubscriptionsPageQueryVariables,
} from './__generated__/ListSubscriptionPage'
import OrganizationListLicenses, {
  ORGANIZATION_LIST_SINGLE_LICENSE,
} from '../../organisms/organization/OrganizationListLicenses.vue'

export default defineComponent({
  components: {
    OrganizationListLicenses,
    Header,
    VerticalNavigation,
    PaddedAppLeftSidebar,
    OrganizationSettingsHeader,
  },
  setup() {
    let route = useRoute()

    let { data } = useQuery<ListSubscriptionsPageQuery, ListSubscriptionsPageQueryVariables>({
      query: gql`
        query ListSubscriptionsPage($shortID: ID!) {
          organization(shortID: $shortID) {
            id
            name
            licenses {
              ...OrganizationListSingleLicense
            }
          }
        }
        ${ORGANIZATION_LIST_SINGLE_LICENSE}
      `,
      requestPolicy: 'cache-and-network',
      variables: {
        shortID: computed(() => route.params.organizationSlug as string),
      },
    })

    return {
      data,
    }
  },
  data() {
    return {
      seats: 25,
    }
  },
  computed: {
    cost() {
      return this.seats * 10
    },
    annualCost() {
      return this.cost * 12
    },
  },
})
</script>
