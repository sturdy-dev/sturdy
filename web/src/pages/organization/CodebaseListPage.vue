<template>
  <PaddedAppLeftSidebar v-if="data" class="bg-white">
    <template #navigation>
      <VerticalNavigation />
    </template>

    <template #header>
      <OrganizationSettingsHeader :name="data.organization.name">
        <RouterLinkButton
          v-if="data.organization.writeable"
          :to="{ name: 'organizationCreateCodebase' }"
        >
          <PlusIcon class="-ml-0.5 mr-2 h-4 w-4" />
          <span>New Codebase</span>
        </RouterLinkButton>
      </OrganizationSettingsHeader>
    </template>

    <template #default>
      <div class="space-y-8 flex flex-col">
        <div v-if="data.organization.codebases.length > 0" class="flex flex-col">
          <div class="align-middle inline-block min-w-full">
            <div class="shadow overflow-hidden border-b border-gray-200 sm:rounded-lg">
              <table class="min-w-full divide-y divide-gray-200 table-fixed">
                <thead class="bg-gray-50">
                  <tr>
                    <th scope="col" />
                    <th
                      scope="col"
                      class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
                    >
                      Codebase
                    </th>
                    <th
                      scope="col"
                      class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider text-right"
                    >
                      Status
                    </th>
                  </tr>
                </thead>
                <tbody class="bg-white divide-y divide-gray-200">
                  <template v-for="cb in data.organization.codebases" :key="cb.id">
                    <tr
                      v-if="cb.isReady"
                      class="hover:bg-gray-100 cursor-pointer"
                      @click="gotoCodebase(cb)"
                    >
                      <td class="md:px-6 md:py-4">
                        <div class="hidden md:block">
                          <AvatarGroup :authors="cb.members" :max="5" />
                        </div>
                      </td>
                      <td class="px-6 py-4 md:px-0 text-sm text-gray-900 w-full">
                        <div class="inline-flex items-center space-x-8">
                          <div class="inline-flex flex-col space-y-1">
                            <span>{{ cb.name }}</span>

                            <Tooltip
                              v-if="cb.gitHubIntegration && cb.gitHubIntegration.enabled"
                              class="text-gray-500 inline-flex space-x-2 group"
                            >
                              <template v-if="cb.gitHubIntegration.gitHubIsSourceOfTruth" #tooltip>
                                GitHub is the source of truth
                              </template>
                              <template v-else #tooltip> Sturdy is the source of truth</template>
                              <template #default>
                                <GitHubIcon class="h-4 w-4 group-hover:text-gray-900" />
                                <span class="group-hover:text-gray-900">
                                  {{ cb.gitHubIntegration.owner }}/{{ cb.gitHubIntegration.name }}
                                </span>
                              </template>
                            </Tooltip>
                          </div>
                          <Pill v-if="cb.isPublic" color="blue">Public</Pill>
                        </div>
                      </td>
                      <td class="px-6 py-4 whitespace-nowrap text-right">
                        <Pill v-if="cb.lastUpdatedAt" color="green">
                          Updated {{ friendly_ago(cb.lastUpdatedAt) }}
                        </Pill>
                        <Pill v-else color="gray"> Empty</Pill>
                      </td>
                    </tr>
                    <tr v-else>
                      <td class="md:px-6 md:py-4">
                        <div class="hidden md:block">
                          <AvatarGroup :authors="cb.members" :max="5" />
                        </div>
                      </td>
                      <td class="px-6 py-4 md:px-0 text-sm text-gray-900 w-full">
                        <div class="inline-flex items-center space-x-4">
                          <Spinner />

                          <div class="inline-flex flex-col space-y-1">
                            <span>{{ cb.name }}</span>

                            <div
                              v-if="cb.gitHubIntegration && cb.gitHubIntegration.enabled"
                              class="text-gray-500 inline-flex space-x-2 group"
                            >
                              <GitHubIcon class="h-4 w-4" />
                              <span>
                                {{ cb.gitHubIntegration.owner }}/{{ cb.gitHubIntegration.name }}
                              </span>
                            </div>
                          </div>
                        </div>
                      </td>
                      <td class="px-6 py-4 whitespace-nowrap text-right">
                        <Pill color="gray">Setting up&hellip;</Pill>
                      </td>
                    </tr>
                  </template>
                </tbody>
              </table>
            </div>
          </div>
        </div>

        <OrganizationNoCodebasesSetupGitHub
          v-if="showGitHubSetupBanner"
          :git-hub-account="data.user.gitHubAccount"
          :git-hub-app="data.gitHubApp"
          :show-start-from-scratch="true"
        />

        <div v-if="showStandaloneSetupBanner" class="bg-gray-100 sm:rounded-lg">
          <div class="flex justify-between px-4 py-5 sm:p-6">
            <div>
              <h3 class="text-lg leading-6 font-medium text-gray-900">
                Create the first codebase 🚀
              </h3>
              <div class="mt-2 max-w-xl text-sm text-gray-500">
                <p>Create a codebase to start coding.</p>
              </div>
            </div>
            <div>
              <CurvedRightIcon class="h-12 w-12 text-gray-700 -rotate-90" />
            </div>
          </div>
        </div>
      </div>
    </template>
  </PaddedAppLeftSidebar>
</template>

<script lang="ts">
import time from '../../time'
import { PlusIcon } from '@heroicons/vue/solid'
import { Slug } from '../../slug'
import Pill from '../../components/shared/Pill.vue'
import AvatarGroup from '../../components/shared/AvatarGroup.vue'
import { gql, useQuery } from '@urql/vue'
import GitHubIcon from '../../components/icons/GitHubIcon.vue'
import Tooltip from '../../components/shared/Tooltip.vue'
import Spinner from '../../components/shared/Spinner.vue'
import { useUpdatedCodebase } from '../../subscriptions/useUpdatedCodebase'
import { computed, defineComponent, inject, ref, type Ref, watch } from 'vue'
import { Feature } from '../../__generated__/types'
import { useRoute } from 'vue-router'
import type {
  CodebaseListPageQuery,
  CodebaseListPageQueryVariables,
} from './__generated__/CodebaseListPage'
import RouterLinkButton from '../../components/shared/RouterLinkButton.vue'
import PaddedAppLeftSidebar from '../../layouts/PaddedAppLeftSidebar.vue'
import VerticalNavigation from '../../organisms/organization/VerticalNavigation.vue'
import OrganizationSettingsHeader from '../../organisms/organization/OrganizationSettingsHeader.vue'
import OrganizationNoCodebasesSetupGitHub from '../../organisms/organization/OrganizationNoCodebasesSetupGitHub.vue'
import CurvedRightIcon from '../../molecules/icons/CurvedRightIcon.vue'
import type { DeepMaybeRef } from '@vueuse/core'

export default defineComponent({
  components: {
    CurvedRightIcon,
    OrganizationSettingsHeader,
    PaddedAppLeftSidebar,
    VerticalNavigation,
    GitHubIcon,
    OrganizationNoCodebasesSetupGitHub,
    PlusIcon,
    Pill,
    AvatarGroup,
    Tooltip,
    Spinner,
    RouterLinkButton,
  },
  setup() {
    const features = inject<Ref<Array<Feature>>>('features', ref([]))
    const isGitHubEnabled = computed(() => features?.value?.includes(Feature.GitHub))
    const isGitHubNotConfigured = computed(() =>
      features?.value?.includes(Feature.GitHubNotConfigured)
    )

    const route = useRoute()
    const organizationSlug = ref(route.params.organizationSlug as string)
    watch(route, (newRoute) => {
      if (
        newRoute.params.organizationSlug &&
        newRoute.params.organizationSlug !== organizationSlug.value
      ) {
        organizationSlug.value = newRoute.params.organizationSlug as string
      }
    })

    const result = useQuery<CodebaseListPageQuery, DeepMaybeRef<CodebaseListPageQueryVariables>>({
      query: gql`
        query CodebaseListPage($organizationID: ID!, $isGitHubEnabled: Boolean!) {
          organization(shortID: $organizationID) {
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
              description
              inviteCode
              createdAt
              archivedAt
              lastUpdatedAt
              isReady
              isPublic
              members {
                id
                name
                avatarUrl
              }
              gitHubIntegration @include(if: $isGitHubEnabled) {
                id
                owner
                name
                enabled
                gitHubIsSourceOfTruth
              }
            }

            writeable
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
      `,
      variables: {
        isGitHubEnabled,
        organizationID: organizationSlug,
      },
    })

    useUpdatedCodebase()

    return {
      fetching: result.fetching,
      data: result.data,
      error: result.error,

      isGitHubEnabled,
      isGitHubNotConfigured,
    }
  },
  computed: {
    serverCanGitHub() {
      return this.isGitHubEnabled || this.isGitHubNotConfigured
    },
    showGitHubSetupBanner() {
      return (
        this.serverCanGitHub &&
        this.data &&
        (this.data.organization.codebases.length === 0 || !this.data.user.gitHubAccount)
      )
    },
    showStandaloneSetupBanner() {
      return !this.serverCanGitHub && this.data && this.data.organization.codebases.length === 0
    },
  },
  watch: {
    error: function (err) {
      if (err) throw err
    },
  },
  methods: {
    friendly_ago(ts) {
      return time.getRelativeTime(new Date(ts * 1000))
    },
    codebaseSlug(cb) {
      return Slug(cb.name, cb.shortID)
    },
    gotoCodebase(cb) {
      this.$router.push({ name: 'codebaseHome', params: { codebaseSlug: this.codebaseSlug(cb) } })
    },
  },
})
</script>
