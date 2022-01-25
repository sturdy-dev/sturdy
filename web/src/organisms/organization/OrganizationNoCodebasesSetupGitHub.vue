<template>
  <div class="bg-gray-100 sm:rounded-lg">
    <div class="px-4 py-5 sm:p-6">
      <h3 class="text-lg leading-6 font-medium text-gray-900">
        Setup <strong>Sturdy for GitHub</strong>
      </h3>
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
            <span>Compatible with pull requests</span>
          </li>
        </ul>
      </div>

      <div class="mt-5 space-x-2 flex">
        <RouterLinkButton v-if="showStartFromScratch" :to="{ name: 'organizationCreateCodebase' }">
          Use Sturdy without GitHub
        </RouterLinkButton>

        <GitHubConnectButton
          :git-hub-account="gitHubAccount"
          :git-hub-app="gitHubApp"
          already-installed-text="Update installation"
          color="blue"
          :git-hub-redirect-state="gitHubRedirectState"
        />
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { CheckIcon } from '@heroicons/vue/solid/esm'
import GitHubConnectButton from '../../molecules/GitHubConnectButton.vue'
import { defineComponent } from 'vue'
import RouterLinkButton from '../../components/shared/RouterLinkButton.vue'

export default defineComponent({
  components: { GitHubConnectButton, CheckIcon, RouterLinkButton },
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
  computed: {
    gitHubRedirectState() {
      return 'web-' + this.$route.fullPath + '/settings/github'
    },
  },
})
</script>
