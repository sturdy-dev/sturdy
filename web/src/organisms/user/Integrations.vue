<template>
  <div>
    <div>
      <h2 class="text-lg leading-6 font-medium text-gray-900">Integrations</h2>
      <p class="mt-1 text-sm text-gray-500">Integrate Sturdy with other tools and platforms.</p>
    </div>
    <ul class="mt-2 divide-y divide-gray-200">
      <li class="py-4 flex items-center justify-between">
        <div class="flex flex-col">
          <p class="text-sm font-medium text-gray-900">GitHub</p>
          <p v-if="user.gitHubAccount" class="text-sm text-gray-500">
            You're connected to GitHub as
            <a :href="'https://github.com/' + user.gitHubAccount.login">
              {{ user.gitHubAccount.login }}
            </a>
          </p>
          <p v-else class="text-sm text-gray-500">
            Integrate Sturdy with GitHub, to easily work with GitHub repositories from Sturdy.
          </p>
        </div>

        <GitHubConnectButton :git-hub-account="user.gitHubAccount" :git-hub-app="gitHubApp">
          <Button :disabled="fetchingRefreshGitHubCodebases" @click="refreshGitHubCodebases">
            {{ fetchingRefreshGitHubCodebases ? 'Loading' : 'Reload' }}
          </Button>
        </GitHubConnectButton>
      </li>
    </ul>
  </div>
</template>

<script lang="ts">
import { gql, useMutation } from '@urql/vue'
import { PropType } from 'vue'
import GitHubConnectButton, {
  GITHUB_ACCOUNT_FRAGMENT,
  GITHUB_APP_FRAGMENT,
} from '../../molecules/GitHubConnectButton.vue'
import Button from '../../components/shared/Button.vue'
import {
  IntegrationsGitHubAppFragment,
  IntegrationsUserFragment,
} from './__generated__/Integrations'

export const INTEGRATIONS_GITHUB_APP_FRAGMENT = gql`
  fragment IntegrationsGitHubApp on GitHubApp {
    ...GitHubApp
  }
  ${GITHUB_APP_FRAGMENT}
`

export const INTEGRATIONS_USER_FRAGMENT = gql`
  fragment IntegrationsUser on User {
    id
    gitHubAccount {
      id
      login
      ...GitHubAccount
    }
  }
  ${GITHUB_ACCOUNT_FRAGMENT}
`

export default {
  components: {
    GitHubConnectButton,
    Button,
  },
  props: {
    user: {
      type: Object as PropType<IntegrationsUserFragment>,
      required: true,
    },
    gitHubApp: {
      type: Object as PropType<IntegrationsGitHubAppFragment>,
      required: true,
    },
  },
  setup() {
    const {
      executeMutation: refreshGitHubCodebasesResult,
      fetching: fetchingRefreshGitHubCodebases,
    } = useMutation(gql`
      mutation RefreshGitHubCodebases {
        refreshGitHubCodebases {
          id
        }
      }
    `)
    return {
      fetchingRefreshGitHubCodebases,
      async refreshGitHubCodebases() {
        const variables = {}
        await refreshGitHubCodebasesResult(variables).then((result) => {
          if (result.error) {
            throw new Error(result.error.toString())
          }
          console.log('refreshGitHubCodebases', result)
        })
      },
    }
  },
}
</script>
