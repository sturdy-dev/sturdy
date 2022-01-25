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
        <span>Buy subscription</span>
      </Header>

      <p class="my-2">A subscription allows you to self-host Sturdy Enterprise.</p>

      <div class="flex w-full">
        <div class="w-80">
          <label class="text-xl">Seats</label>

          <div class="w-42">
            <input v-model="seats" class="w-full" type="range" min="20" max="500" step="5" />
          </div>

          <p class="text-sm">One seat equals to one user, minimum 20 users.</p>
        </div>

        <div class="bg-gray-100 ml-16 p-2 rounded-md flex flex-col">
          <div class="inline-flex">
            <UserGroupIcon class="w-5 h-5 mr-2" />
            <span>
              <strong>{{ seats }}</strong> seats
            </span>
          </div>
          <div class="inline-flex">
            <CurrencyDollarIcon class="w-5 h-5 mr-2" />
            <span>
              <strong>${{ cost }}</strong> / month (billed annually)
            </span>
          </div>
          <div class="inline-flex pl-7 text-sm">
            <em>${{ annualCost }} annually</em>
          </div>
        </div>
      </div>

      <p class="my-2">Reach out to support@getsturdy.com to get started.</p>
    </template>
  </PaddedAppLeftSidebar>
</template>

<script lang="ts">
import { computed, defineComponent } from 'vue'
import Header from '../../molecules/Header.vue'
import VerticalNavigation from '../../organisms/organization/VerticalNavigation.vue'
import PaddedAppLeftSidebar from '../../layouts/PaddedAppLeftSidebar.vue'
import { CurrencyDollarIcon, UserGroupIcon } from '@heroicons/vue/solid'
import { gql, useQuery } from '@urql/vue'
import { useRoute } from 'vue-router'
import {
  ListOrganizationSubscriptionsQuery,
  ListOrganizationSubscriptionsQueryVariables,
} from './__generated__/CreateSubscriptionPage'
import OrganizationSettingsHeader from '../../organisms/organization/OrganizationSettingsHeader.vue'

export default defineComponent({
  components: {
    Header,
    VerticalNavigation,
    PaddedAppLeftSidebar,
    OrganizationSettingsHeader,
    UserGroupIcon,
    CurrencyDollarIcon,
  },
  setup() {
    let route = useRoute()

    let { data } = useQuery<
      ListOrganizationSubscriptionsQuery,
      ListOrganizationSubscriptionsQueryVariables
    >({
      query: gql`
        query ListOrganizationSubscriptions($shortID: ID!) {
          organization(shortID: $shortID) {
            id
            name
            licenseSubscriptions {
              id
              seats
              usedSeats
              licenseKey
            }
          }
        }
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
