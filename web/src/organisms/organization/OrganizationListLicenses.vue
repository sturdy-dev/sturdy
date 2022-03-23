<template>
  <div v-if="licenses" class="bg-white shadow sm:rounded-md">
    <ul role="list" class="divide-y divide-gray-200">
      <li v-for="license in licenses" :key="license.id">
        <div class="space-y-1">
          <div class="px-4 py-4 sm:px-6">
            <div class="flex items-center justify-between">
              <code class="text-sm font-medium text-blue-600 truncate">
                {{ license.key }}
              </code>
              <div class="ml-2 flex-shrink-0 flex">
                <Pill v-if="license.status === 'Valid'"> {{ license.status }}</Pill>
                <Pill v-else color="yellow"> {{ license.status }}</Pill>
              </div>
            </div>
            <div class="mt-2 sm:flex sm:justify-between">
              <div class="sm:flex">
                <p class="flex items-center text-sm text-gray-500">
                  <UsersIcon
                    class="flex-shrink-0 mr-1.5 h-5 w-5 text-gray-400"
                    aria-hidden="true"
                  />
                  {{ license.seats }} seats
                </p>
                <p
                  v-if="false"
                  class="mt-2 flex items-center text-sm text-gray-500 sm:mt-0 sm:ml-6"
                >
                  <LocationMarkerIcon
                    class="flex-shrink-0 mr-1.5 h-5 w-5 text-gray-400"
                    aria-hidden="true"
                  />
                  asdasd
                </p>
              </div>
              <div class="mt-2 flex items-center text-sm text-gray-500 sm:mt-0">
                <CalendarIcon
                  class="flex-shrink-0 mr-1.5 h-5 w-5 text-gray-400"
                  aria-hidden="true"
                />
                <p>
                  Expires at
                  {{ ' ' }}
                  <time :datetime="license.expiresAt">{{ unixToHuman(license.expiresAt) }}</time>
                </p>
              </div>
            </div>
          </div>
          <div v-for="msg in license.messages" :key="msg.text">
            <Banner :status="bannerStatus(msg)">
              {{ msg.text }}
            </Banner>
          </div>
        </div>
      </li>
    </ul>
  </div>
</template>

<script lang="ts">
import { CalendarIcon, LocationMarkerIcon, UsersIcon } from '@heroicons/vue/solid'
import type { PropType } from 'vue'
import { defineComponent } from 'vue'
import { gql } from '@urql/vue'
import type { OrganizationListSingleLicenseFragment } from './__generated__/OrganizationListLicenses'
import Pill from '../../components/shared/Pill.vue'
import { Banner } from '../../atoms'

export const ORGANIZATION_LIST_SINGLE_LICENSE = gql`
  fragment OrganizationListSingleLicense on License {
    id
    key
    createdAt
    expiresAt
    status
    seats

    messages {
      level
      text
      type
    }
  }
`

export default defineComponent({
  components: {
    CalendarIcon,
    LocationMarkerIcon,
    UsersIcon,
    Pill,
    Banner,
  },
  props: {
    licenses: {
      type: Object as PropType<OrganizationListSingleLicenseFragment>,
      required: true,
    },
  },
  methods: {
    unixToHuman(ts) {
      return new Date(ts * 1000).toLocaleString()
    },
    bannerStatus(msg) {
      if (msg.level === 'Warning') {
        return 'warning'
      }
      if (msg.level === 'Error') {
        return 'error'
      }
      return ''
    },
  },
})
</script>
