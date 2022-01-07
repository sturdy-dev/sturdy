<template>
  <LinkButton v-if="data && !data.user.gitHubAccount" :href="github_oauth_url" :color="color">
    Connect to GitHub
  </LinkButton>
  <div v-else-if="data" class="space-x-2">
    <slot></slot>
    <LinkButton :href="github_manage_installation_url" :color="color">
      {{ alreadyInstalledText }}
    </LinkButton>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import { gql, useQuery } from '@urql/vue'
import LinkButton from '../components/shared/LinkButton.vue'
import {
  GitHubConnectButtonQuery,
  GitHubConnectButtonQueryVariables,
} from './__generated__/GitHubConnectButton'

export default defineComponent({
  components: { LinkButton },
  props: {
    alreadyInstalledText: {
      type: String,
      required: false,
      default: 'Manage installation',
    },
    color: {
      type: String,
      required: false,
      default: 'white',
    },
    gitHubRedirectState: {
      type: String,
      required: false,
      default: 'user-settings',
    },
  },
  setup() {
    const GitHubConnectButtonQuery = gql`
      query GitHubConnectButton {
        user {
          id
          gitHubAccount {
            id
            login
          }
        }

        gitHubApp {
          _id
          name
          clientID
        }
      }
    `

    let { data } = useQuery<GitHubConnectButtonQuery, GitHubConnectButtonQueryVariables>({
      query: GitHubConnectButtonQuery,
      requestPolicy: 'cache-and-network',
    })

    return {
      data,
    }
  },

  computed: {
    github_oauth_url() {
      if (!this.data) {
        return '#'
      }

      const url = new URL('https://github.com/login/oauth/authorize')
      url.searchParams.set('client_id', this.data.gitHubApp.clientID)
      url.searchParams.set('state', this.gitHubRedirectState)

      if (typeof ipc !== 'undefined') {
        const callbackURL = new URL('sturdy:///setup-github')
        if (import.meta.env.DEV) {
          callbackURL.protocol = 'sturdy-dev:'
        }
        url.searchParams.set('redirect_uri', callbackURL.href)
      }

      return url.href
    },
    github_manage_installation_url() {
      if (!this.data) {
        return '#'
      }

      const url = new URL(
        'https://github.com/apps/' + this.data.gitHubApp.name + '/installations/new'
      )

      let state = 'install'
      if (typeof ipc !== 'undefined') {
        state = 'install-app'
      }

      url.searchParams.set('state', state)

      return url.href
    },
  },
})
</script>
