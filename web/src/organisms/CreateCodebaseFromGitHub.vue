<template>
  <div class="space-y-2">
    <div class="mt-2 bg-gray-100 p-8 rounded">
      <p v-if="!gitHubAccount">
        Authenticate with GitHub and install <strong>Sturdy for GitHub</strong> to use Sturdy on top
        of your existing repositories.
      </p>
      <p v-else>
        Install <strong>Sturdy for GitHub</strong> to use Sturdy on top of your existing
        repositories.
      </p>

      <ul class="list-inside mt-2 block inline-flex flex-col text-gray-800">
        <li class="inline-flex space-x-2">
          <CheckIcon class="h-5 w-5 text-green-400 flex-shrink-0" />
          <span>Use Sturdy on top of GitHub</span>
        </li>
        <li class="inline-flex space-x-2">
          <CheckIcon class="h-5 w-5 text-green-400 flex-shrink-0" />
          <span>Fine grained permissions</span>
        </li>
        <li class="inline-flex space-x-2">
          <CheckIcon class="h-5 w-5 text-green-400 flex-shrink-0" />
          <span>Install on only the repositories that you want to use Sturdy on</span>
        </li>
        <li class="inline-flex space-x-2">
          <CheckIcon class="h-5 w-5 text-green-400 flex-shrink-0" />
          <span>Work in Sturdy, create pull requests with your code when you're done</span>
        </li>
        <li class="inline-flex space-x-2">
          <CheckIcon class="h-5 w-5 text-green-400 flex-shrink-0" />
          <span>Sturdy automatically syncs data from GitHub</span>
        </li>
      </ul>
    </div>

    <div v-if="!organization.writeable">
      <p class="text-sm text-gray-500">
        You don't have permissions to connect GitHub repositories in this organization, ask an admin
        for help if you want to setup a GitHub connection.
      </p>
    </div>

    <template v-if="organization.writeable">
      <GitHubConnectButton
        v-if="isGitHubEnabled"
        already-installed-text="Update GitHub-app installation"
        not-connected-text="Login with GitHub"
        color="blue"
        :git-hub-app="gitHubApp"
        :git-hub-account="gitHubAccount"
      />
      <LinkButton
        v-else
        href="https://getsturdy.com/v2/docs/self-hosted#setup-github-integration"
        target="_blank"
      >
        Read the docs
      </LinkButton>

      <template v-if="data && data.gitHubRepositories.length > 0 && gitHubApp.validation.ok">
        <div class="border-b border-gray-200">
          <ul role="list" class="divide-y divide-gray-200">
            <li
              v-for="repo in data.gitHubRepositories"
              :key="repo.id"
              class="py-4 flex justify-between items-center"
            >
              <div class="ml-3 flex flex-col">
                <span class="font-medium text-gray-900">
                  {{ repo.gitHubOwner }}/{{ repo.gitHubName }}
                </span>
                <span
                  v-if="
                    repo?.codebase?.organization?.id &&
                    repo.codebase.organization.id === organization.id
                  "
                  class="text-sm text-gray-500"
                >
                  {{ repo.gitHubName }} is connected to {{ organization.name }}
                </span>
                <span
                  v-else-if="
                    repo?.codebase?.organization?.id &&
                    repo.codebase.organization.id !== organization.id
                  "
                  class="text-sm text-gray-500"
                >
                  Connected to: {{ repo.codebase.organization.name }}
                </span>
              </div>
              <div v-if="repo?.codebase?.isReady">
                <RouterLinkButton
                  :to="{ name: 'codebaseHome', params: { codebaseSlug: slug(repo.codebase) } }"
                  color="green"
                >
                  Open
                </RouterLinkButton>
              </div>
              <div v-else-if="repo?.codebase" class="flex items-center space-x-2">
                <Spinner />
                <span>Getting ready&hellip;</span>
              </div>
              <div v-else>
                <Button @click="installRepo(repo)">Import</Button>
              </div>
            </li>
          </ul>
        </div>
      </template>

      <div v-else-if="fetching" class="flex items-center space-x-2">
        <Spinner />
        <span>Loading repositories, please wait&hellip;</span>
      </div>

      <Banner status="info">
        Not seeing the repository you want to install setup? Update the app installation above to
        install <em>Sturdy for GitHub</em> on more organizations or repositories.
      </Banner>
    </template>
  </div>
</template>

<script lang="ts">
import { computed, defineComponent, inject, ref } from 'vue'
import type { PropType, Ref } from 'vue'
import { gql, useQuery } from '@urql/vue'

import type {
  OrganizationSetupGitHub_GitHubAccountFragment,
  OrganizationSetupGitHub_GitHubAppFragment,
  OrganizationSetupGitHub_OrganizationFragment,
  OrganizationSetupGitHubQuery,
  OrganizationSetupGitHubQueryVariables,
} from './__generated__/CreateCodebaseFromGitHub'
import Spinner from '../atoms/Spinner.vue'
import Button from '../atoms/Button.vue'
import { Slug } from '../slug'
import RouterLinkButton from '../atoms/RouterLinkButton.vue'
import GitHubConnectButton, {
  GITHUB_CONNECT_BUTTON_GITHUB_APP_FRAGMENT,
} from '../molecules/GitHubConnectButton.vue'
import { useSetupGitHubRepository } from '../mutations/useSetupGitHubRepository'
import { CheckIcon } from '@heroicons/vue/solid'
import { Feature } from '../__generated__/types'
import LinkButton from '../atoms/LinkButton.vue'
import Banner from '../atoms/Banner.vue'

export const ORGANIZATION_SETUP_GITHUB_GITHUB_APP_FRAGMENT = gql`
  fragment OrganizationSetupGitHub_GitHubApp on GitHubApp {
    _id
    name
    clientID
    ...GitHubConnectButton_GitHubApp
  }
  ${GITHUB_CONNECT_BUTTON_GITHUB_APP_FRAGMENT}
`

export const ORGANIZATION_SETUP_GITHUB_GITHUB_ACCOUNT_FRAGMENT = gql`
  fragment OrganizationSetupGitHub_GitHubAccount on GitHubAccount {
    id
    login
    isValid
  }
`

export const ORGANIZATION_SETUP_GITHUB_ORGANIZATION_FRAGMENT = gql`
  fragment OrganizationSetupGitHub_Organization on Organization {
    id
    name
    writeable
  }
`

export default defineComponent({
  components: {
    GitHubConnectButton,
    Spinner,
    Button,
    RouterLinkButton,
    CheckIcon,
    LinkButton,
    Banner,
  },
  props: {
    organization: {
      type: Object as PropType<OrganizationSetupGitHub_OrganizationFragment>,
      required: true,
    },
    gitHubApp: {
      type: Object as PropType<OrganizationSetupGitHub_GitHubAppFragment>,
      default: null,
      required: false,
    },
    gitHubAccount: {
      type: Object as PropType<OrganizationSetupGitHub_GitHubAccountFragment>,
      default: null,
      required: false,
    },
  },
  setup() {
    const features = inject<Ref<Array<Feature>>>('features', ref([]))
    const isGitHubEnabled = computed(() => features?.value?.includes(Feature.GitHub))

    let { data, fetching } = useQuery<
      OrganizationSetupGitHubQuery,
      OrganizationSetupGitHubQueryVariables
    >({
      query: gql`
        query OrganizationSetupGitHub {
          gitHubRepositories {
            id
            gitHubInstallationID
            gitHubRepositoryID
            gitHubOwner
            gitHubName
            codebase {
              id
              shortID
              name
              isReady
              organization {
                id
                name
              }
            }
          }
        }
      `,
    })

    let execSetupGitHubRepository = useSetupGitHubRepository()

    return {
      data,
      fetching,

      isGitHubEnabled,

      async setupGitHubRepository(
        organizationID: string,
        gitHubInstallationID: string,
        gitHubRepositoryID: string
      ) {
        const variables = { organizationID, gitHubInstallationID, gitHubRepositoryID }
        return execSetupGitHubRepository(variables)
      },
    }
  },
  data() {
    return {
      showInvitedBanner: false,
      showFailedBanner: false,
    }
  },
  methods: {
    async installRepo(repo: { gitHubInstallationID: string; gitHubRepositoryID: string }) {
      await this.setupGitHubRepository(
        this.organization.id,
        repo.gitHubInstallationID,
        repo.gitHubRepositoryID
      )
    },
    slug(cb: { name: string; shortID: string }) {
      return Slug(cb.name, cb.shortID)
    },
  },
})
</script>
