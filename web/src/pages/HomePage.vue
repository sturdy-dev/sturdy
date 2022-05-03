<template>
  <PaddedApp v-if="data" class="bg-white">
    <div class="space-y-8">
      <div v-if="data.organizations.length === 0">
        <h2>You don't have any organizations:</h2>
        <RouterLinkButton :to="{ name: 'organizationCreate' }">Get started now 🚀</RouterLinkButton>
      </div>

      <CreateCodebase v-else-if="codebasesCount == 0" :organization="data.organizations[0]" />

      <ul v-else role="list" class="divide-y divide-gray-200">
        <li v-for="org in data.organizations" :key="org.id">
          <router-link
            :to="{ name: 'organizationListCodebases', params: { organizationSlug: org.shortID } }"
            class="block hover:bg-gray-50"
          >
            <div class="px-4 py-4 flex items-center sm:px-6">
              <div class="min-w-0 flex-1 sm:flex sm:items-center sm:justify-between">
                <div class="truncate">
                  <div class="flex text-sm">
                    <p class="font-medium text-gray-600 truncate">{{ org.name }}</p>
                  </div>
                </div>
                <div class="mt-4 flex-shrink-0 sm:mt-0 sm:ml-5">
                  <div class="flex overflow-hidden -space-x-1">
                    <AvatarGroup :authors="org.members" />
                  </div>
                </div>
              </div>
              <div class="ml-5 flex-shrink-0">
                <ChevronRightIcon class="h-5 w-5 text-gray-400" aria-hidden="true" />
              </div>
            </div>
          </router-link>
        </li>
      </ul>
    </div>
  </PaddedApp>
</template>

<script lang="ts">
import AvatarGroup from '../atoms/AvatarGroup.vue'
import { gql, useQuery } from '@urql/vue'
import PaddedApp from '../layouts/PaddedApp.vue'
import { computed, defineComponent, inject, ref } from 'vue'
import type { Ref } from 'vue'
import { Feature } from '../__generated__/types'
import { ChevronRightIcon } from '@heroicons/vue/solid'
import type { HomePageQuery, HomePageQueryVariables } from './__generated__/HomePage'
import RouterLinkButton from '../atoms/RouterLinkButton.vue'
import type { DeepMaybeRef } from '@vueuse/shared'
import CreateCodebase, { ORGANIZATION_FRAGMENT } from '../organisms/CreateCodebase.vue'

const PAGE_QUERY = gql`
  query HomePage($isGitHubEnabled: Boolean!) {
    organizations {
      id
      name
      shortID
      ...Organization_CreateCodebase

      members {
        id
        name
        avatarUrl
      }

      codebases {
        id
      }
    }

    gitHubApp @include(if: $isGitHubEnabled) {
      _id
      name
      clientID
    }

    user {
      id
      name
      gitHubAccount @include(if: $isGitHubEnabled) {
        id
        login
      }
    }
  }
  ${ORGANIZATION_FRAGMENT}
`

export default defineComponent({
  components: {
    PaddedApp,
    AvatarGroup,
    ChevronRightIcon,
    RouterLinkButton,
    CreateCodebase,
  },
  setup() {
    const features = inject<Ref<Array<Feature>>>('features', ref([]))
    const isGitHubEnabled = computed(() => features?.value?.includes(Feature.GitHub))

    const result = useQuery<HomePageQuery, DeepMaybeRef<HomePageQueryVariables>>({
      query: PAGE_QUERY,
      variables: {
        isGitHubEnabled,
      },
    })

    return {
      fetching: result.fetching,
      data: result.data,
      error: result.error,
    }
  },
  watch: {
    error: function (err) {
      if (err) throw err
    },
  },
  computed: {
    codebasesCount() {
      return this.data?.organizations.reduce((acc, org) => acc + org.codebases.length, 0)
    },
  },
})
</script>
