<template>
  <PaddedAppLeftSidebar v-if="data" class="bg-white">
    <template #navigation>
      <SettingsVerticalNavigation />
    </template>

    <template #default>
      <div class="space-y-8">
        <Header>
          <Pill>Beta</Pill>
          <span>Connect Git to {{ data.codebase.name }}</span>
        </Header>

        <p class="text-sm text-gray-500">
          Follow the instructions below to connect <strong>Sturdy</strong> and <strong>Git</strong>.
        </p>

        <nav aria-label="Progress" class="max-w-4xl">
          <ol role="list" class="overflow-hidden">
            <Step name="Configure host" :status="gitRemoteUrlStatus">
              <p class="my-2 max-w-2xl text-sm text-gray-500">Remote URL (for push and pull)</p>

              <TextInput
                v-if="gitRemoteUrlStatus !== 'pending'"
                v-model="gitRemoteURL"
                placeholder="https://"
              />

              <p class="my-2 max-w-2xl text-sm text-gray-500">
                Which branch should Sturdy import as the trunk?
              </p>

              <TextInput
                v-if="gitRemoteUrlStatus !== 'pending'"
                v-model="trackedBranch"
                placeholder="main"
              />

              <p class="my-2 max-w-2xl text-sm text-gray-500">Name this integration</p>

              <TextInput
                v-if="gitRemoteUrlStatus !== 'pending'"
                v-model="gitRemoteName"
                placeholder="Eg: GitLab, Azure, etc..."
              />
            </Step>

            <Step name="Authenticate" :status="gitAuthStepStatus">
              <p class="my-2 max-w-2xl text-sm text-gray-500">
                Authenticate with the Git Host using Basic Auth.
              </p>

              <div class="space-y-2">
                <TextInput v-model="basicAuthUsername" placeholder="Username" />

                <TextInput v-model="basicAuthPassword" placeholder="Password" />
              </div>
            </Step>

            <Step name="Save" :status="saveUpdateStepStatus" :is-last="true">
              <div class="flex flex-col space-y-2">
                <Banner v-if="error && error.length > 0" status="error">{{ error }}</Banner>
                <Banner v-if="showSuccess" status="success">Saved!</Banner>

                <div>
                  <Button
                    v-if="data.codebase?.remote?.id"
                    color="green"
                    @click="createOrUpdateCodebaseRemote"
                  >
                    Update
                  </Button>
                  <Button v-else color="green" @click="createOrUpdateCodebaseRemote">Create</Button>
                </div>
              </div>
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
import { Status } from '../../../../../components/ci/StepIndicator.vue'
import { gql, useQuery } from '@urql/vue'
import { ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import { IdFromSlug } from '../../../../../slug'
import Pill from '../../../../../components/shared/Pill.vue'
import PaddedAppLeftSidebar from '../../../../../layouts/PaddedAppLeftSidebar.vue'
import SettingsVerticalNavigation from '../../../../../components/codebase/settings/SettingsVerticalNavigation.vue'
import Header from '../../../../../molecules/Header.vue'
import { Banner } from '../../../../../atoms'
import { GetGitIntegrationsQuery, GetGitIntegrationsQueryVariables } from './__generated__/Git'
import { useCreateOrUpdateCodebaseRemote } from '../../../../../mutations/useCreateOrUpdateGitRemote'
import Button from '../../../../../components/shared/Button.vue'

export default {
  components: {
    SettingsVerticalNavigation,
    PaddedAppLeftSidebar,
    Step,
    TextInput,
    Pill,
    Header,
    Banner,
    Button,
  },
  setup() {
    const route = useRoute()
    const shortCodebaseID = IdFromSlug(route.params.codebaseSlug as string)

    const { data } = useQuery<GetGitIntegrationsQuery, GetGitIntegrationsQueryVariables>({
      query: gql`
        query GetGitIntegrations($shortCodebaseID: ID!) {
          codebase(shortID: $shortCodebaseID) {
            id
            name

            remote {
              id
              name
              url
              trackedBranch
              basicAuthUsername
              basicAuthPassword
            }
          }
        }
      `,
      variables: {
        shortCodebaseID: shortCodebaseID,
      },
    })

    const gitRemoteURL = ref('')
    const gitRemoteName = ref('')
    const trackedBranch = ref('')
    const basicAuthUsername = ref('')
    const basicAuthPassword = ref('')
    const showSuccess = ref(false)
    const error = ref<Error | string | null>(null)

    // Set data from API
    watch(
      data,
      (newData) => {
        if (!newData?.codebase?.remote) {
          return
        }
        gitRemoteURL.value = newData.codebase.remote.url
        gitRemoteName.value = newData.codebase.remote.name
        trackedBranch.value = newData.codebase.remote.trackedBranch
        basicAuthUsername.value = newData.codebase.remote.basicAuthUsername
        basicAuthPassword.value = newData.codebase.remote.basicAuthPassword
      },
      {
        immediate: true,
      }
    )

    const createOrUpdateCodebaseRemoteFunc = useCreateOrUpdateCodebaseRemote()

    return {
      data,
      shortCodebaseID,
      showSuccess,

      gitRemoteURL,
      gitRemoteName,
      trackedBranch,
      basicAuthUsername,
      basicAuthPassword,

      error,

      async createOrUpdateCodebaseRemote() {
        if (!data.value?.codebase?.id) {
          error.value = 'Could not update, please try again later...'
          return
        }

        const vars = {
          name: gitRemoteName.value,
          codebaseID: data.value.codebase.id,
          url: gitRemoteURL.value,
          trackedBranch: trackedBranch.value,
          basicAuthUsername: basicAuthUsername.value,
          basicAuthPassword: basicAuthPassword.value,
        }

        await createOrUpdateCodebaseRemoteFunc(vars)
          .then(() => {
            showSuccess.value = true
            setTimeout(() => (showSuccess.value = false), 5000)
          })
          .catch((e) => {
            error.value = e
          })
      },
    }
  },
  computed: {
    gitRemoteUrlStatus(): Status {
      return this.gitRemoteURL ? 'completed' : 'current'
    },
    gitAuthStepStatus(): Status {
      if (this.basicAuthUsername && this.basicAuthPassword) return 'completed'
      return this.gitRemoteUrlStatus === 'completed' ? 'current' : 'pending'
    },
    saveUpdateStepStatus(): Status {
      if (this.data?.codebase?.remote?.id) return 'completed'
      return this.gitAuthStepStatus === 'completed' ? 'current' : 'pending'
    },
  },
}
</script>
