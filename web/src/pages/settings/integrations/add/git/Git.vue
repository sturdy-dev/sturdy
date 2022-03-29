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
          <ol role="list">
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
              <div v-if="!keyPairID" class="space-y-2">
                <p class="my-2 max-w-2xl text-sm text-gray-500">
                  Authenticate with <strong>{{ gitRemoteName }}</strong> using Basic Auth.
                </p>

                <TextInput v-model="basicAuthUsername" placeholder="Username" />

                <TextInput v-model="basicAuthPassword" placeholder="Password" />

                <!-- Disable SSH, we're having some issues with it -->
                <Button
                  v-if="false"
                  :spinner="generatingPrivateKey"
                  :disabled="generatingPrivateKey"
                  @click="generateKeyPair"
                >
                  <span v-if="generatingPrivateKey">Generating a new keypair...</span>
                  <span v-else>Switch to SSH keypair auth</span>
                </Button>
              </div>

              <div v-if="keyPairID" class="space-y-2 w-full">
                <p class="my-2 max-w-2xl text-sm text-gray-500">
                  Use this Public Key to authenticate Sturdy and
                  <strong>{{ gitRemoteName }}</strong>
                </p>

                <InputCopyToClipboard :value="keyPairPublicKey" class="w-full" />
                <Button @click="unsetKeyPair">
                  Use Basic Auth (username and password) instead
                </Button>
              </div>
            </Step>

            <Step name="Links" :status="gitAuthStepStatus">
              <div class="space-y-4">
                <p class="text-sm text-gray-500">
                  Links to use when linking from Sturdy to <strong>{{ gitRemoteName }}</strong>
                </p>

                <div>
                  <p class="text-sm text-gray-500">Link to repository</p>
                  <TextInput
                    v-model="browserLinkRepo"
                    placeholder="https://my-host.com/repo/name"
                  />

                  <div
                    v-if="recommendedLinkRepo && recommendedLinkRepo !== browserLinkRepo"
                    class="text-sm text-gray-500 border-l-2 border-green-400 p-2 my-4 bg-green-50"
                  >
                    <p>
                      Found a recommended link. "<code class="underline">{{
                        recommendedLinkRepo
                      }}</code
                      >". Do you want to use it?
                    </p>
                    <Button size="small" class="mt-2" @click="browserLinkRepo = recommendedLinkRepo"
                      >Yes!
                    </Button>
                  </div>
                </div>

                <div>
                  <p class="text-sm text-gray-500">
                    Link to branch (template). Use <code>${BRANCH_NAME}</code> as a variable for the
                    branch name.
                  </p>
                  <TextInput
                    v-model="browserLinkBranch"
                    placeholder="https://my-host.com/repo/name/branch/${BRANCH_NAME}"
                  />

                  <div
                    v-if="recommendedLinkBranch && recommendedLinkBranch !== browserLinkBranch"
                    class="text-sm text-gray-500 border-l-2 border-green-400 p-2 my-4 bg-green-50"
                  >
                    <p>
                      Found a recommended link. "<code>{{ recommendedLinkRepo }}</code
                      >". Do you want to use it?
                    </p>
                    <Button
                      size="small"
                      class="mt-2"
                      @click="browserLinkBranch = recommendedLinkBranch"
                      >Yes
                    </Button>
                  </div>
                </div>
              </div>
            </Step>

            <Step name="Save" :status="saveUpdateStepStatus">
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

            <Step
              v-if="webhookTrigger"
              name="Webhooks (optional)"
              :is-last="true"
              :status="saveUpdateStepStatus"
            >
              <p class="text-sm text-gray-500">
                For a better (and faster) experience, configure {{ gitRemoteName }} to send webhooks
                to the following URL on pushes and merges.
              </p>
              <InputCopyToClipboard :value="webhookTrigger" />
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
import type { Status } from '../../../../../components/ci/StepIndicator.vue'
import { gql, useQuery } from '@urql/vue'
import { defineComponent, ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import { IdFromSlug } from '../../../../../slug'
import Pill from '../../../../../atoms/Pill.vue'
import PaddedAppLeftSidebar from '../../../../../layouts/PaddedAppLeftSidebar.vue'
import SettingsVerticalNavigation from '../../../../../components/codebase/settings/SettingsVerticalNavigation.vue'
import Header from '../../../../../molecules/Header.vue'
import { Banner } from '../../../../../atoms'
import type { GetGitIntegrationsQuery, GetGitIntegrationsQueryVariables } from './__generated__/Git'
import { useCreateOrUpdateCodebaseRemote } from '../../../../../mutations/useCreateOrUpdateGitRemote'
import Button from '../../../../../atoms/Button.vue'
import { defaultLinkBranch, defaultLinkRepo } from './Links'
import InputCopyToClipboard from '../../../../../organisms/InputCopyToClipboard.vue'
import http from '../../../../../http'
import { useGenerateKeyPair } from '../../../../../mutations/useGenerateKeyPair'
import { KeyPairType } from '../../../../../__generated__/types'

export default defineComponent({
  components: {
    InputCopyToClipboard,
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
              keyPair {
                id
                publicKey
              }
              browserLinkRepo
              browserLinkBranch
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
    const basicAuthUsername = ref<string | null | undefined>(undefined)
    const basicAuthPassword = ref<string | null | undefined>(undefined)
    const browserLinkRepo = ref('')
    const browserLinkBranch = ref('')
    const keyPairID = ref<string | null | undefined>(undefined)
    const keyPairPublicKey = ref<string | null | undefined>('')
    const showSuccess = ref(false)
    const error = ref<Error | string | null>(null)

    // Set data from API (only once)
    let didLoad = false
    watch(
      data,
      (newData) => {
        if (!newData?.codebase?.remote || didLoad) {
          return
        }
        gitRemoteURL.value = newData.codebase.remote.url
        gitRemoteName.value = newData.codebase.remote.name
        trackedBranch.value = newData.codebase.remote.trackedBranch
        basicAuthUsername.value = newData.codebase.remote.basicAuthUsername
        basicAuthPassword.value = newData.codebase.remote.basicAuthPassword
        browserLinkRepo.value = newData.codebase.remote.browserLinkRepo
        browserLinkBranch.value = newData.codebase.remote.browserLinkBranch
        keyPairID.value = newData.codebase.remote?.keyPair?.id
        keyPairPublicKey.value = newData.codebase.remote?.keyPair?.publicKey
        didLoad = true
      },
      {
        immediate: true,
      }
    )

    const createOrUpdateCodebaseRemoteFunc = useCreateOrUpdateCodebaseRemote()

    watch(gitRemoteURL, () => {
      if (gitRemoteURL.value && !browserLinkRepo.value) {
        let n = defaultLinkRepo(gitRemoteURL.value)
        if (n) {
          browserLinkRepo.value = n
        }
      }

      if (gitRemoteURL.value && !browserLinkBranch.value) {
        let n = defaultLinkBranch(gitRemoteURL.value)
        if (n) {
          browserLinkBranch.value = n
        }
      }
    })

    const { mutating: generatingPrivateKey, generateKeyPair: generateKeyPairFunc } =
      useGenerateKeyPair()

    return {
      data,
      shortCodebaseID,
      showSuccess,

      gitRemoteURL,
      gitRemoteName,
      trackedBranch,
      basicAuthUsername,
      basicAuthPassword,
      browserLinkRepo,
      browserLinkBranch,

      keyPairID,
      keyPairPublicKey,

      error,

      generatingPrivateKey,

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
          browserLinkRepo: browserLinkRepo.value,
          browserLinkBranch: browserLinkBranch.value,
          keyPairID: keyPairID.value,
        }

        if (vars.keyPairID) {
          vars.basicAuthUsername = undefined
          vars.basicAuthPassword = undefined
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

      async generateKeyPair() {
        await generateKeyPairFunc({ keyPairType: KeyPairType.Rsa_4096 }).then((kp) => {
          keyPairID.value = kp.generateKeyPair.id
          keyPairPublicKey.value = kp.generateKeyPair.publicKey
        })
      },
    }
  },
  computed: {
    gitRemoteUrlStatus(): Status {
      return this.gitRemoteURL && this.trackedBranch && this.gitRemoteName ? 'completed' : 'current'
    },
    gitAuthStepStatus(): Status {
      if (this.basicAuthUsername && this.basicAuthPassword) return 'completed'
      if (this.keyPairID && this.keyPairPublicKey) return 'completed'
      return this.gitRemoteUrlStatus === 'completed' ? 'current' : 'pending'
    },
    saveUpdateStepStatus(): Status {
      if (this.data?.codebase?.remote?.id) return 'completed'
      return this.gitAuthStepStatus === 'completed' ? 'current' : 'pending'
    },
    recommendedLinkRepo() {
      return defaultLinkRepo(this.gitRemoteURL)
    },
    recommendedLinkBranch() {
      return defaultLinkBranch(this.gitRemoteURL)
    },
    webhookTrigger(): string {
      if (!window.location || !this?.data?.codebase?.id) {
        return ''
      }
      const base = http.url('/v3/remotes/webhook/sync-codebase/' + this.data.codebase.id)
      // using the current browser location as the base, used if url() returns a relative url
      return new URL(base, new URL(window.location.href)).href
    },
  },
  methods: {
    unsetKeyPair() {
      this.keyPairID = undefined
      this.keyPairPublicKey = undefined
    },
  },
})
</script>
