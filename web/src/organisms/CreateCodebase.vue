<template>
  <div class="py-8 px-4">
    <div class="">
      <h2 class="text-4xl font-extrabold text-gray-900 sm:text-4xl sm:tracking-tight lg:text-4xl">
        Create a new codebase in <span class="underline">{{ organization.name }}</span>
      </h2>
      <p class="mt-5 text-xl text-gray-500">You'll soon be ready to code! 📈</p>
    </div>
  </div>

  <div class="flex space-y-4 xl:space-y-0 xl:space-x-4 flex-col xl:flex-row">
    <div class="bg-gray-100 sm:rounded-lg">
      <div class="px-4 py-5 sm:p-6 flex flex-col justify-between h-full">
        <div>
          <h3 class="text-lg leading-6 font-medium text-gray-900">
            Create new codebase on <strong>Strudy</strong>
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
            <RouterLinkButton
              :to="{
                name: 'organizationCreateSturdyCodebase',
                params: { organizationSlug: organization.shortID },
              }"
              color="blue"
            >
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
            Import existing codebase <strong>from GitHub</strong>
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
            <RouterLinkButton
              v-if="isGitHubAvailable"
              :to="{
                name: 'organizationCreateGitHubCodebase',
                params: { organizationSlug: organization.shortID },
              }"
              color="blue"
            >
              Import from GitHub
            </RouterLinkButton>
            <LinkButton
              v-else
              href="https://getsturdy.com/v2/docs/self-hosted#setup-github-integration"
              target="_blank"
            >
              Learn how to configure
            </LinkButton>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { CheckIcon } from '@heroicons/vue/solid/esm'
import RouterLinkButton from '../atoms/RouterLinkButton.vue'
import { Feature } from '../__generated__/types'
import { defineComponent, inject, computed, type Ref, ref, type PropType } from 'vue'
import LinkButton from '../atoms/LinkButton.vue'
import { gql } from '@urql/vue'
import type { Organization_CreateCodebaseFragment } from './__generated__/CreateCodebase'

export const ORGANIZATION_FRAGMENT = gql`
  fragment Organization_CreateCodebase on Organization {
    id
    shortID
    name
  }
`

export default defineComponent({
  components: { CheckIcon, RouterLinkButton, LinkButton },
  props: {
    organization: {
      type: Object as PropType<Organization_CreateCodebaseFragment>,
      required: true,
    },
  },
  setup() {
    const features = inject<Ref<Array<Feature>>>('features', ref([]))
    const isGitHubEnabled = computed(() => features?.value?.includes(Feature.GitHub))
    const isGitHubEnabledNotConfigured = computed(() =>
      features?.value?.includes(Feature.GitHubNotConfigured)
    )
    return {
      isGitHubEnabled,
      isGitHubEnabledNotConfigured,
    }
  },
  computed: {
    gitHubRedirect() {
      return this.$route.fullPath + '/settings/github'
    },
    isGitHubAvailable() {
      return this.isGitHubEnabled && !this.isGitHubEnabledNotConfigured
    },
  },
})
</script>
