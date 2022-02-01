<template>
  <LinkButton v-if="!gitHubAccount" :href="github_oauth_url" :color="color">
    {{ notConnectedText }}
  </LinkButton>
  <div v-else class="space-x-2">
    <slot></slot>
    <LinkButton :href="github_manage_installation_url" :color="color">
      {{ alreadyInstalledText }}
    </LinkButton>
  </div>
</template>

<script lang="ts">
import { defineComponent, PropType } from 'vue'
import { gql } from '@urql/vue'
import LinkButton from '../components/shared/LinkButton.vue'
import { GitHubAccountFragment, GitHubAppFragment } from './__generated__/GitHubConnectButton'

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
  components: { LinkButton },
  props: {
    gitHubApp: {
      type: Object as PropType<GitHubAppFragment>,
      required: true,
    },
    gitHubAccount: {
      type: Object as PropType<GitHubAccountFragment>,
      default: null,
    },
    alreadyInstalledText: {
      type: String,
      required: false,
      default: 'Manage installation',
    },
    notConnectedText: {
      type: String,
      required: false,
      default: 'Connect to GitHub',
    },
    color: {
      type: String,
      required: false,
      default: 'white',
    },
    stateSamePage: {
      type: Boolean,
      required: false,
    },
    statePath: {
      type: String,
      required: false,
    },
  },
  computed: {
    github_oauth_url() {
      const url = new URL('https://github.com/login/oauth/authorize')
      url.searchParams.set('client_id', this.gitHubApp.clientID)
      url.searchParams.set('state', this.state)

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
      const url = new URL('https://github.com/apps/' + this.gitHubApp.name + '/installations/new')
      url.searchParams.set('state', this.state)
      return url.href
    },

    statePrefix() {
      if (typeof ipc !== 'undefined' && import.meta.env.DEV) {
        return 'app-dev'
      } else if (typeof ipc !== 'undefined') {
        return 'app'
      }
      return 'web'
    },

    state() {
      let prefix = this.statePrefix

      let path
      if (this.statePath) {
        path = this.statePath
      } else {
        path = this.$route.fullPath
      }

      return `${prefix}-${path}`
    },
  },
})
</script>
