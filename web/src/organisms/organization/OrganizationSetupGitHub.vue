<template>
  <div class="space-y-2">
    <div class="space-y-1">Setup Sturdy for GitHub</div>

    <GitHubConnectButton
      already-installed-text="Update GitHub-app installation"
      :git-hub-app="gitHubApp"
      :git-hub-account="gitHubAccount"
    />

    <div v-if="data" class="border-b border-gray-200">
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
                repo?.codebase?.organization?.id && repo.codebase.organization.id !== organizationId
              "
              class="text-sm text-gray-500"
            >
              Setup in {{ repo.codebase.organization.name }}
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
    <div v-else class="flex items-center space-x-2">
      <Spinner />
      <span>Refreshing, please wait&hellip;</span>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent, PropType } from 'vue'
import { gql, useQuery } from '@urql/vue'

import {
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

export default defineComponent({
  components: { GitHubConnectButton, Spinner, Button, RouterLinkButton },
  props: {
    organizationId: {
      type: String,
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
    let { data } = useQuery<OrganizationSetupGitHubQuery, OrganizationSetupGitHubQueryVariables>({
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
        this.organizationId,
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
