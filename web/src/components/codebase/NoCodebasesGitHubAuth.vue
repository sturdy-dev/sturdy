<template>
  <div class="bg-white sm:rounded-lg">
    <div class="px-4 py-5 sm:p-6">
      <h3 class="text-lg leading-6 font-medium text-gray-900">Coming from GitHub?</h3>
      <div class="mt-2 max-w-xl text-sm text-gray-500">
        <p>Import your existing repositories to Sturdy, without losing any history!</p>
        <ul class="list-inside mt-2">
          <li class="inline-flex space-x-2">
            <CheckIcon class="h-5 w-5 text-green-400" />
            <span>Keep your history</span>
          </li>
          <li class="inline-flex space-x-2">
            <CheckIcon class="h-5 w-5 text-green-400" />
            <span>Easy migration path, you can migrate your team members to Sturdy one by one</span>
          </li>
          <li class="inline-flex space-x-2">
            <CheckIcon class="h-5 w-5 text-green-400" />
            <span>Compatible with Pull Requests</span>
          </li>
        </ul>
      </div>
      <div class="mt-5 space-x-2 flex">
        <Button v-if="showStartFromScratch" @click="$router.push({ name: 'codebaseCreate' })">
          Start from scratch
        </Button>

        <GitHubConnectButton
          :gitHubAccount="gitHubAccount"
          :gitHubApp="gitHubApp"
          already-installed-text="Import repositories"
          color="blue"
          git-hub-redirect-state="codebase-overview"
        />
      </div>
    </div>
  </div>
</template>

<script>
import { CheckIcon } from '@heroicons/vue/solid/esm'
import Button from '../shared/Button.vue'
import GitHubConnectButton from '../../molecules/GitHubConnectButton.vue'

export default {
  name: 'NoCodebasesGitHubAuth',
  components: { GitHubConnectButton, CheckIcon, Button },
  props: {
    showStartFromScratch: {
      type: Boolean,
      default: false,
    },
    gitHubApp: {
      type: Object,
      required: true,
    },
    gitHubAccount: {
      type: Object,
      default: null,
    },
  },
}
</script>
