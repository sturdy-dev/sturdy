<template>
  <div class="flex space-y-4 xl:space-y-0 xl:space-x-4 flex-col xl:flex-row">
    <div class="bg-gray-100 sm:rounded-lg">
      <div class="px-4 py-5 sm:p-6 flex flex-col justify-between h-full">
        <div>
          <h3 class="text-lg leading-6 font-medium text-gray-900">Create a new codebase</h3>

          <div class="mt-2 max-w-xl text-sm text-gray-500">
            <p>Working on something new? Create a new codebase on Sturdy.</p>
            <ul class="list-inside mt-2 inline-flex flex-col">
              <li class="inline-flex space-x-2">
                <CheckIcon class="h-5 w-5 text-green-400 flex-shrink-0" />
                <span>Start from scratch</span>
              </li>
              <li class="inline-flex space-x-2">
                <CheckIcon class="h-5 w-5 text-green-400 flex-shrink-0" />
                <span>Host your project on Sturdy</span>
              </li>
              <li class="inline-flex space-x-2">
                <CheckIcon class="h-5 w-5 text-green-400 flex-shrink-0" />
                <span>Use the full Sturdy workflow</span>
              </li>
            </ul>
          </div>
        </div>
        <div>
          <div class="mt-5 space-x-2 flex">
            <RouterLinkButton
              v-if="showStartFromScratch"
              :to="{ name: 'organizationCreateCodebase' }"
              color="blue"
            >
              Create a empty codebase
            </RouterLinkButton>
          </div>
        </div>
      </div>
    </div>

    <div class="bg-gray-100 sm:rounded-lg">
      <div class="px-4 py-5 sm:p-6 flex flex-col justify-between h-full">
        <div>
          <h3 class="text-lg leading-6 font-medium text-gray-900">
            Setup <strong>Sturdy for GitHub</strong>
          </h3>
          <div class="mt-2 max-w-xl text-sm text-gray-500">
            <p>
              Install the bridge between Sturdy and GitHub, to use Sturdy on top of existing
              GitHub-repositories.
            </p>
            <ul class="list-inside mt-2 inline-flex flex-col">
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
        </div>
        <div>
          <div class="mt-5 space-x-2 flex">
            <GitHubConnectButton
              v-if="gitHubApp"
              :git-hub-account="gitHubAccount"
              :git-hub-app="gitHubApp"
              already-installed-text="Update installation"
              not-connected-text="Login with GitHub"
              color="blue"
              :state-path="gitHubRedirect"
            />
            <LinkButton
              v-else
              href="https://getsturdy.com/v2/docs/self-hosted#setup-github-integration"
              target="_blank"
            >
              Read the docs
            </LinkButton>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { CheckIcon } from '@heroicons/vue/solid/esm'
import GitHubConnectButton from '../../molecules/GitHubConnectButton.vue'
import { defineComponent } from 'vue'
import RouterLinkButton from '../../atoms/RouterLinkButton.vue'
import LinkButton from '../../atoms/LinkButton.vue'

export default defineComponent({
  components: { GitHubConnectButton, CheckIcon, RouterLinkButton, LinkButton },
  props: {
    showStartFromScratch: {
      type: Boolean,
      default: false,
    },
    gitHubApp: {
      type: Object,
      required: false,
      default: null,
    },
    gitHubAccount: {
      type: Object,
      default: null,
      required: false,
    },
  },
  computed: {
    gitHubRedirect() {
      return this.$route.fullPath + '/settings/github'
    },
  },
})
</script>
