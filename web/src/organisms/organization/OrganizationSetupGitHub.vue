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

    <GitHubConnectButton
      already-installed-text="Update GitHub-app installation"
      not-connected-text="Login with GitHub"
      color="blue"
      :git-hub-app="gitHubApp"
      :git-hub-account="gitHubAccount"
    />

    <template v-if="data && data.gitHubRepositories.length > 0">
      <p class="mt-4 text-sm p-4">
        Click <em>Setup</em> to import the repository to Sturdy, and connect it to
        <em>{{ organization.name }}</em
        >.
      </p>

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
              <Button @click="installRepo(repo)">Setup</Button>
            </div>
          </li>
        </ul>
      </div>
    </template>
    <div v-else-if="fetching" class="flex items-center space-x-2">
      <Spinner />
      <span>Loading repositories, please wait&hellip;</span>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent, PropType } from 'vue'
import { gql, useQuery } from '@urql/vue'

import {
  OrganizationSetupGitHubOrganizationFragment,
  OrganizationSetupGitHubQuery,
  OrganizationSetupGitHubQueryVariables,
} from './__generated__/OrganizationSetupGitHub'
import Spinner from '../../components/shared/Spinner.vue'
import Button from '../../components/shared/Button.vue'
import { Slug } from '../../slug'
import RouterLinkButton from '../../components/shared/RouterLinkButton.vue'
import GitHubConnectButton from '../../molecules/GitHubConnectButton.vue'
import {
  GitHubAccountFragment,
  GitHubAppFragment,
} from '../../molecules/__generated__/GitHubConnectButton'
import { useSetupGitHubRepository } from '../../mutations/useSetupGitHubRepository'
import { CheckIcon } from '@heroicons/vue/solid'

export const GITHUB_APP_FRAGMENT = gql`
  fragment GitHubApp on GitHubApp {
    _id
    name
    clientID
  }
`

export const GITHUB_ACCOUNT_FRAGMENT = gql`
  fragment GitHubAccount on GitHubAccount {
    id
    login
  }
`

export const ORGANIZATION_FRAGMENT = gql`
  fragment OrganizationSetupGitHubOrganization on Organization {
    id
    name
  }
`

export default defineComponent({
  components: { GitHubConnectButton, Spinner, Button, RouterLinkButton, CheckIcon },
  props: {
    organization: {
      type: Object as PropType<OrganizationSetupGitHubOrganizationFragment>,
      required: true,
    },
    gitHubApp: {
      type: Object as PropType<GitHubAppFragment>,
      required: true,
    },
    gitHubAccount: {
      type: Object as PropType<GitHubAccountFragment>,
      default: null,
    },
  },
  setup() {
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
    async installRepo(repo) {
      await this.setupGitHubRepository(
        this.organization.id,
        repo.gitHubInstallationID,
        repo.gitHubRepositoryID
      )
    },
    slug(cb) {
      return Slug(cb.name, cb.shortID)
    },
  },
})
</script>
