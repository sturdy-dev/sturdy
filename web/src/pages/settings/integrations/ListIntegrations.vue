<template>
  <PaddedAppLeftSidebar v-if="data?.codebase" class="bg-white">
    <template #navigation>
      <SettingsVerticalNavigation />
    </template>

    <template #default>
      <Header>
        <Pill>Beta</Pill>
        <span>Integrate {{ data.codebase.name }} with other services</span>
      </Header>

      <div class="bg-white shadow overflow-hidden rounded-md mt-8">
        <ul role="list" class="divide-y divide-gray-200">
          <template v-for="item in list" :key="item.name">
            <li
              class="flex items-center space-x-4 px-6 py-4"
              :class="!item.enabled ? ['opacity-50'] : []"
            >
              <img :src="item.logo" class="h-10 w-10" />
              <div class="flex-1">
                <h3>{{ item.name }}</h3>
                <p class="text-gray-500 text-sm">{{ item.description }}</p>
                <Pill
                  v-if="configuredProviders.has(item.name)"
                  color="green"
                  class="text-green-500 text-sm"
                >
                  Installed
                </Pill>
                <Pill v-else color="gray">Not Installed</Pill>
              </div>



              <RouterLinkButton
                :disabled="!item.enabled"
                :to="{
                  name: item.page,
                }"
              >
                <span v-if="item.supportMulti">Add</span>
                <span v-else>Edit</span>
              </RouterLinkButton>
            </li>
            <template v-if="configuredProviders.has(item.name)">
              <li
                v-for="instance in configuredProviders.get(item.name)"
                :key="instance.id"
                class="px-6 py-4 flex space-x-4 hover:cursor-pointer hover:bg-gray-50 items-center"
              >
                <p class="flex-1 pl-14">
                  Pipeline
                  <strong>{{ instance.configuration.pipelineName }}</strong> in
                  <strong>{{ instance.configuration.organizationName }}</strong>
                </p>

                <RouterLinkButton
                  :to="{
                    name: 'codebaseSettingsEditBuildkite',
                    params: { integrationId: instance.id },
                  }"
                >
                  Edit
                </RouterLinkButton>
                <Button color="red" @click="doDeleteIntegration(instance.id)">Delete</Button>
              </li>
            </template>
          </template>
        </ul>
      </div>
    </template>
  </PaddedAppLeftSidebar>
</template>

<script lang="ts">
import { gql, useQuery } from '@urql/vue'
import { useRoute } from 'vue-router'
import { IdFromSlug } from '../../../slug'
import {
  GetIntegrationsQuery,
  GetIntegrationsQueryVariables,
  IntegrationListItemFragment,
} from './__generated__/ListIntegrations'
import Pill from '../../../components/shared/Pill.vue'
import Button from '../../../components/shared/Button.vue'
import PaddedAppLeftSidebar from '../../../layouts/PaddedAppLeftSidebar.vue'
import SettingsVerticalNavigation from '../../../components/codebase/settings/SettingsVerticalNavigation.vue'
import Header from '../../../molecules/Header.vue'
import RouterLinkButton from '../../../components/shared/RouterLinkButton.vue'
import { useDeleteIntegration } from '../../../mutations/useDeleteIntegration'
import { computed, defineComponent, inject, ref, Ref } from 'vue'
import { Feature } from '../../../__generated__/types'

import buildkiteLogo from '../../../components/ci/logos/BuildkiteLogo.svg'
import gitLogo from '../../../components/ci/logos/GitLogo.svg'

const INTEGRATION_FRAGMENT = gql`
  fragment IntegrationListItem on Integration {
    id
    provider
    deletedAt
    ... on BuildkiteIntegration {
      id
      configuration {
        id
        organizationName
        pipelineName
      }
    }
  }
`

export default defineComponent({
  components: {
    SettingsVerticalNavigation,
    PaddedAppLeftSidebar,
    Pill,
    Button,
    Header,
    RouterLinkButton,
  },
  setup() {
    const route = useRoute()
    const shortCodebaseID = IdFromSlug(route.params.codebaseSlug as string)

    const features = inject<Ref<Array<Feature>>>('features', ref([]))
    const isGitHubEnabled = computed(() => features?.value?.includes(Feature.GitHub))

    const { data } = useQuery<GetIntegrationsQuery, GetIntegrationsQueryVariables>({
      query: gql`
        query GetIntegrations($shortCodebaseID: ID!) {
          codebase(shortID: $shortCodebaseID) {
            id
            name
            integrations {
              ...IntegrationListItem
            }
            remote {
              id
            }
          }
        }

        ${INTEGRATION_FRAGMENT}
      `,
      variables: {
        shortCodebaseID: shortCodebaseID,
      },
      requestPolicy: 'cache-and-network',
    })

    const deleteIntegration = useDeleteIntegration()
    return {
      isGitHubEnabled,

      data,
      shortCodebaseID,

      doDeleteIntegration: function (id: string) {
        deleteIntegration({ id: id })
      },
    }
  },
  computed: {
    list() {
      return [
        {
          name: 'Buildkite',
          description: 'Setup CI/CD with Buildkite',
          page: 'codebaseSettingsAddBuildkite',
          enabled: this.isGitHubEnabled,
          logo: buildkiteLogo,
          supportMulti: true,
        },

        // TODO: Uncomment when ready!
        /*{
          name: 'Git',
          description: 'Sync Sturdy with any Git Provider (GitLab, Azure DevOps, etc)',
          page: 'codebaseSettingsAddGit',
          enabled: this.isGitHubEnabled,
          logo: gitLogo,
          supportMulti: false,
        },*/
      ]
    },
    nonDeletedIntegrations(): Array<IntegrationListItemFragment> {
      let res = this.data?.codebase?.integrations.filter((i) => !i.deletedAt)
      if (!res) {
        return []
      }
      return res
    },
    configuredProviders() {
      let res = new Map<string, Array<IntegrationListItemFragment>>()

      for (const provider of this.nonDeletedIntegrations) {
        let existing = res.get(provider.provider)
        if (existing) {
          existing.push(provider)
        } else {
          res.set(provider.provider, new Array<IntegrationListItemFragment>(provider))
        }
      }

      if (this.data?.codebase?.remote?.id) {
        res.set('Git', new Array<IntegrationListItemFragment>())
      }

      return res
    },
  },
})
</script>
