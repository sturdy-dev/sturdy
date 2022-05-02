<template>
  <div class="flex space-y-4 xl:space-y-0 xl:space-x-4 flex-col xl:flex-row">
    <div class="bg-gray-100 sm:rounded-lg">
      <div class="px-4 py-5 sm:p-6 flex flex-col justify-between h-full">
        <div>
          <h3 class="text-lg leading-6 font-medium text-gray-900">
            Create codebase on <strong>Strudy</strong>
          </h3>

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
            <RouterLinkButton :to="{ name: 'organizationCreateSturdyCodebase' }" color="blue">
              Create an empty codebase
            </RouterLinkButton>
          </div>
        </div>
      </div>
    </div>

    <div class="bg-gray-100 sm:rounded-lg">
      <div class="px-4 py-5 sm:p-6 flex flex-col justify-between h-full">
        <div>
          <h3 class="text-lg leading-6 font-medium text-gray-900">
            Import codebase <strong>from GitHub</strong>
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
            <RouterLinkButton :to="{ name: 'organizationCreateGitHubCodebase' }" color="blue">
              Import from GitHub
            </RouterLinkButton>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { CheckIcon } from '@heroicons/vue/solid/esm'
import { defineComponent } from 'vue'
import RouterLinkButton from '../atoms/RouterLinkButton.vue'

export default defineComponent({
  components: { CheckIcon, RouterLinkButton },
  computed: {
    gitHubRedirect() {
      return this.$route.fullPath + '/settings/github'
    },
  },
})
</script>
