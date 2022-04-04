<template>
  <PaddedAppLeftSidebar v-if="data" class="bg-white">
    <template #navigation>
      <SettingsVerticalNavigation />
    </template>

    <template #default>
      <div class="space-y-8">
        <Header>
          <Pill>Beta</Pill>
          <span>Connect Buildkite to {{ data.codebase.name }}</span>
        </Header>

        <p class="text-sm text-gray-500">
          Follow the instructions below to connect <strong>Sturdy</strong> and
          <strong>Buildkite</strong>.
        </p>

        <nav aria-label="Progress" class="max-w-4xl">
          <ol role="list" class="overflow-hidden">
            <Step
              name="What's your Buildkite organization name?"
              :status="buildkiteOrganizationStepStatus"
            >
              <TextInput
                v-if="buildkiteOrganizationStepStatus !== 'pending'"
                v-model="buildkiteOrganizationName"
                placeholder="My Organization Name"
              />
            </Step>

            <Step name="Create API Token" :status="apiTokenStepStatus">
              <p class="my-2 max-w-2xl text-sm text-gray-500">
                Follow the link below, and setup a new API token on Buildkite. This token will be
                used by Sturdy to trigger new builds.
              </p>
              <LinkButton
                href="https://buildkite.com/user/api-access-tokens/new"
                target="_blank"
                class="my-2"
              >
                <span>Create new API token</span>
                <ExternalLinkIcon class="h-5 w-5 ml-1" />
              </LinkButton>
              <Instructions
                description="When asked, enter the following settings:"
                :instructions="createAPITokenInstructions"
              />
              <TextInput
                v-model="buildkiteAPIToken"
                placeholder="Enter your new Buildkite API token"
              />
            </Step>

            <Step name="Create new pipeline" :status="pipelineStatus">
              <CreatePipelineStep
                v-if="pipelineStatus !== 'pending'"
                v-model="buildkitePipelineName"
                :buildkite-organization-slug="buildkiteOrganizationSlug"
                :short-codebase-id="shortCodebaseID"
              />
            </Step>

            <Step name="Setup webhook" :status="webhookStatus">
              <p class="my-2 max-w-2xl text-sm text-gray-500">
                It's time to configure Buildkite to send build statuses back to Sturdy. Follow the
                instructions to setup a webhook from Buildkite to Sturdy.
              </p>
              <LinkButton :href="buildkiteWebhookPage" target="_blank" class="my-2">
                <span>Create new webhook</span>
                <ExternalLinkIcon class="h-5 w-5 ml-1" />
              </LinkButton>
              <Instructions
                description="When asked, enter the following settings:"
                :instructions="createWebhookInstructions"
              />
              <p class="my-2 max-w-2xl text-sm text-gray-500">
                Buildkite has generated a <strong>token</strong> for Sturdy to verify that the
                webhook is authentic. Copy it from Buildkite and enter it below.
              </p>
              <TextInput v-model="buildkiteWebhookToken" placeholder="Webhook Token" />
            </Step>

            <Step name="You are all set!" :status="testIntegrationStatus" :is-last="true">
              <TestIntegration
                v-if="testIntegrationStatus !== 'pending'"
                :api-token="buildkiteAPIToken"
                :organization-name="buildkiteOrganizationName"
                :pipeline-name="buildkitePipelineName"
                :webhook-secret="buildkiteWebhookToken"
                :short-codebase-id="shortCodebaseID"
                :editing-integration-id="selectedIntegrationID"
                :change="headChange"
              />
            </Step>
          </ol>
        </nav>
      </div>
    </template>
  </PaddedAppLeftSidebar>
</template>

<script lang="ts">
import TextInput from '../../../../../molecules/TextInput.vue'
import Step from '../../../../../components/ci/Step.vue'
import LinkButton from '../../../../../atoms/LinkButton.vue'
import { ExternalLinkIcon } from '@heroicons/vue/solid'
import type { Status } from '../../../../../components/ci/StepIndicator.vue'
import Instructions from '../../../../../components/ci/Instructions.vue'
import type { Instruction } from '../../../../../components/ci/Instructions.vue'
import CreatePipelineStep from '../../../../../components/ci/CreatePipeline.vue'
import TestIntegration from '../../../../../components/ci/TestIntegration.vue'
import { gql, useQuery } from '@urql/vue'
import { computed, ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import { IdFromSlug } from '../../../../../slug'
import type {
  GetBuildkiteIntegrationsQuery,
  GetBuildkiteIntegrationsQueryVariables,
} from './__generated__/Buildkite'
import Pill from '../../../../../atoms/Pill.vue'
import PaddedAppLeftSidebar from '../../../../../layouts/PaddedAppLeftSidebar.vue'
import SettingsVerticalNavigation from '../../../../../components/codebase/settings/SettingsVerticalNavigation.vue'
import Header from '../../../../../molecules/Header.vue'
import type { InputMaybe } from '../../../../../__generated__/types'
import { IntegrationProvider } from '../../../../../__generated__/types'

const toKebabCase = (str: string): string => {
  return str.toLowerCase().replace(/ /, '-')
}

export default {
  components: {
    SettingsVerticalNavigation,
    PaddedAppLeftSidebar,
    Instructions,
    Step,
    TextInput,
    LinkButton,
    ExternalLinkIcon,
    CreatePipelineStep,
    TestIntegration,
    Pill,
    Header,
  },
  setup() {
    const route = useRoute()
    const shortCodebaseID = IdFromSlug(route.params.codebaseSlug as string)

    const selectedIntegrationIDInput = computed(() => route.params?.integrationId as InputMaybe)
    const selectedIntegrationID = computed(() => route.params?.integrationId)

    const { data } = useQuery<
      GetBuildkiteIntegrationsQuery,
      GetBuildkiteIntegrationsQueryVariables
    >({
      query: gql`
        query GetBuildkiteIntegrations($shortCodebaseID: ID!, $optionalIntegrationID: ID) {
          codebase(shortID: $shortCodebaseID) {
            id
            name

            integrations(id: $optionalIntegrationID) {
              provider

              ... on BuildkiteIntegration {
                id
                configuration {
                  id
                  organizationName
                  pipelineName
                  apiToken
                  webhookSecret
                }
              }
            }
            changes(input: { limit: 1 }) {
              id
              title
            }
          }
        }
      `,
      variables: {
        shortCodebaseID: shortCodebaseID,
        optionalIntegrationID: selectedIntegrationIDInput,
      },
    })

    const buildkiteOrganizationName = ref('')
    const buildkiteAPIToken = ref('')
    const buildkitePipelineName = ref('')
    const buildkiteWebhookToken = ref('')

    // Set data from API
    watch(
      data,
      (newData) => {
        if (!newData?.codebase?.integrations) {
          return
        }
        for (const provider of newData.codebase.integrations) {
          if (
            provider.provider === IntegrationProvider.Buildkite &&
            provider.configuration &&
            provider.id === route.params?.integrationId
          ) {
            buildkiteOrganizationName.value = provider.configuration.organizationName
            buildkiteAPIToken.value = provider.configuration.apiToken
            buildkitePipelineName.value = provider.configuration.pipelineName
            buildkiteWebhookToken.value = provider.configuration.webhookSecret
          }
        }
      },
      {
        immediate: true,
      }
    )

    return {
      data,
      shortCodebaseID,
      selectedIntegrationID,

      buildkiteOrganizationName,
      buildkiteAPIToken,
      buildkitePipelineName,
      buildkiteWebhookToken,
    }
  },
  computed: {
    buildkiteOrganizationSlug(): string {
      return toKebabCase(this.buildkiteOrganizationName)
    },
    buildkiteWebhookPage(): string {
      return `https://buildkite.com/organizations/${this.buildkiteOrganizationSlug}/services/webhook/new`
    },

    createAPITokenInstructions(): Instruction[] {
      return [
        {
          name: 'Description',
          value: 'Sturdy API Key',
        },
        {
          name: 'Organization Access',
          value: this.buildkiteOrganizationSlug,
        },
        {
          name: 'Rest API Scopes',
          value: 'Modify Builds',
        },
      ]
    },

    createWebhookInstructions(): Instruction[] {
      return [
        { name: 'Webhook URL', value: 'https://api.getsturdy.com/v3/statuses/webhook', pre: true },
        {
          name: 'Token',
          value: 'Choose "Send an HMAC signature in the X-Buildkite-Signature header"',
        },
        {
          name: 'Events',
          value:
            'ping, build.scheduled, build.running, build.finished, job.scheduled, job.running, job.finished, job.activated',
        },
      ]
    },

    buildkiteOrganizationStepStatus(): Status {
      return this.buildkiteOrganizationName ? 'completed' : 'current'
    },
    apiTokenStepStatus(): Status {
      if (this.buildkiteAPIToken) return 'completed'
      return this.buildkiteOrganizationName ? 'current' : 'pending'
    },
    pipelineStatus(): Status {
      if (this.buildkitePipelineName) return 'completed'
      return this.buildkiteAPIToken ? 'current' : 'pending'
    },
    webhookStatus(): Status {
      if (this.buildkiteWebhookToken) return 'completed'
      return this.buildkitePipelineName ? 'current' : 'pending'
    },
    testIntegrationStatus(): Status {
      return this.buildkiteWebhookToken ? 'completed' : 'pending'
    },
    headChange() {
      if (this.data?.codebase?.changes?.length !== 1) {
        return null
      }
      return this.data?.codebase?.changes[0]
    },
  },
}
</script>
