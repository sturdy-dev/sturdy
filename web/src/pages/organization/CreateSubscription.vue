<template>
  <PaddedAppLeftSidebar class="bg-white">
    <template #navigation>
      <VerticalNavigation />
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
            <input class="w-full" type="range" min="20" max="500" step="5" v-model="seats" />
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

      <pre>data={{ data }}</pre>
    </template>
  </PaddedAppLeftSidebar>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import Header from '../../molecules/Header.vue'
import VerticalNavigation from '../../organisms/organization/VerticalNavigation.vue'
import PaddedAppLeftSidebar from '../../layouts/PaddedAppLeftSidebar.vue'
import { CurrencyDollarIcon, UserGroupIcon } from '@heroicons/vue/solid'
import { gql, useQuery } from '@urql/vue'
import { useRoute } from 'vue-router'
import {
  ListOrganizationSubscriptionsQuery,
  ListOrganizationSubscriptionsQueryVariables,
} from './__generated__/CreateSubscription'

export default defineComponent({
  components: {
    Header,
    VerticalNavigation,
    PaddedAppLeftSidebar,
    UserGroupIcon,
    CurrencyDollarIcon,
  },
  setup() {
    let route = useRoute()
    let orgID = route.params.id as string

    let { data } = useQuery<
      ListOrganizationSubscriptionsQuery,
      ListOrganizationSubscriptionsQueryVariables
    >({
      query: gql`
        query ListOrganizationSubscriptions($id: ID!) {
          organization(id: $id) {
            id
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
        id: orgID,
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
