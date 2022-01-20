<template>
  <PaddedApp v-if="data">
    <div class="space-y-8">
      <div class="relative">
        <div class="absolute inset-0 flex items-center" aria-hidden="true">
          <div class="w-full border-t border-gray-300" />
        </div>

        <div class="relative flex items-center justify-between">
          <span class="pr-2 bg-gray-100 text-sm text-gray-500"> Your codebases on Sturdy </span>

          <div class="pl-2 bg-gray-100 text-sm text-gray-500">
            <router-link
              :to="{ name: 'codebaseCreate' }"
              class="inline-flex items-center px-3 py-2 border border-gray-300 shadow-sm text-sm leading-4 font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
            >
              <PlusIcon class="-ml-0.5 mr-2 h-4 w-4" />
              <span>New</span>
            </router-link>
          </div>
        </div>
      </div>

      <div v-if="data.codebases.length > 0" class="flex flex-col">
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
                <template v-for="cb in data.codebases" :key="cb.id">
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
                            <template v-else #tooltip> Sturdy is the source of truth </template>
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
    </div>

    <NoCodebasesGitHubAuth
      v-if="isGitHubEnabled && data && (data.codebases.length === 0 || !data.user.gitHubAccount)"
      class="mt-4"
      :git-hub-account="data.user.gitHubAccount"
      :git-hub-app="data.gitHubApp"
      :show-start-from-scratch="true"
    />
  </PaddedApp>
</template>

<script lang="ts">
import time from '../time'
import { PlusIcon } from '@heroicons/vue/solid'
import { Slug } from '../slug'
import Pill from '../components/shared/Pill.vue'
import AvatarGroup from '../components/shared/AvatarGroup.vue'
import { gql, useQuery } from '@urql/vue'
import NoCodebasesGitHubAuth from '../components/codebase/NoCodebasesGitHubAuth.vue'
import GitHubIcon from '../components/icons/GitHubIcon.vue'
import Tooltip from '../components/shared/Tooltip.vue'
import Spinner from '../components/shared/Spinner.vue'
import { useUpdatedCodebase } from '../subscriptions/useUpdatedCodebase'
import PaddedApp from '../layouts/PaddedApp.vue'
import { defineComponent, toRefs } from 'vue'
import { Feature } from '../__generated__/types'
import {
  CodebaseOverviewQuery,
  CodebaseOverviewQueryVariables,
} from './__generated__/CodebaseOverview'

export default defineComponent({
  components: {
    PaddedApp,
    GitHubIcon,
    NoCodebasesGitHubAuth,
    PlusIcon,
    Pill,
    AvatarGroup,
    Tooltip,
    Spinner,
  },
  props: {
    features: {
      type: Array,
      required: true,
    },
  },
  setup(props) {
    const { features } = toRefs(props)
    const isGitHubEnabled = features.value.includes(Feature.GitHub)

    const result = useQuery<CodebaseOverviewQuery, CodebaseOverviewQueryVariables>({
      query: gql`
        query CodebaseOverview($isGitHubEnabled: Boolean!) {
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
      },
    })

    useUpdatedCodebase()

    return {
      fetching: result.fetching,
      data: result.data,
      error: result.error,

      isGitHubEnabled,
    }
  },
  data() {
    return {}
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
