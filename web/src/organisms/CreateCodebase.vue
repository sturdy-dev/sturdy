<template>
  <div class="py-8 px-4">
    <Listbox as="span" v-model="selected" v-if="organizations.length > 1">
      <div class="mt-1 relative">
        <ListboxButton
          class="bg-white relative rounded-md pr-10 py-2 text-left cursor-default focus:outline-none font-extrabold"
        >
          <span class="block truncate">{{ selected?.name }}</span>
          <span class="absolute inset-y-0 right-0 flex items-center pr-2 pointer-events-none">
            <SelectorIcon class="h-5 w-5 text-gray-400" aria-hidden="true" />
          </span>
        </ListboxButton>

        <transition
          leave-active-class="transition ease-in duration-100"
          leave-from-class="opacity-100"
          leave-to-class="opacity-0"
        >
          <ListboxOptions
            class="absolute z-10 mt-1 bg-white shadow-lg max-h-60 rounded-md py-1 text-base ring-1 ring-black ring-opacity-5 overflow-auto focus:outline-none"
          >
            <ListboxOption
              as="template"
              v-for="org in organizations"
              :key="org.id"
              :value="org"
              v-slot="{ active, selected }"
            >
              <li
                :class="[
                  active ? 'text-white bg-blue-600' : 'text-gray-900',
                  'cursor-default select-none relative py-2 pl-3 pr-9',
                ]"
              >
                <span :class="[selected ? 'font-semibold' : 'font-normal', 'block truncate']">
                  {{ org.name }}
                </span>

                <span
                  v-if="selected"
                  :class="[
                    active ? 'text-white' : 'text-blue-600',
                    'absolute inset-y-0 right-0 flex items-center pr-4',
                  ]"
                >
                </span>
              </li>
            </ListboxOption>
          </ListboxOptions>
        </transition>
      </div>
    </Listbox>
    <h2 class="text-4xl font-extrabold text-gray-900 sm:text-4xl sm:tracking-tight lg:text-4xl">
      Create a new codebase
    </h2>
    <p class="mt-5 text-xl text-gray-500">You'll soon be ready to code! 📈</p>
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
            <ul class="list-inside mt-2 inline-flex flex-col gap-2">
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
                params: { organizationSlug: selected.shortID },
              }"
              color="blue"
            >
              Create empty
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
            <ul class="list-inside mt-2 inline-flex flex-col gap-2">
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
                <span>Sturdy automatically stays up to date with GitHub</span>
              </li>
              <li class="inline-flex space-x-2">
                <CheckIcon class="h-5 w-5 text-green-400 flex-shrink-0" />
                <span>Gradually migrate to native Sturdy at your own pase</span>
              </li>
            </ul>
          </div>
        </div>
        <div>
          <div class="mt-5 space-x-2 flex">
            <LinkButton
              v-if="isGitHubEnabledNotConfigured"
              href="https://getsturdy.com/v2/docs/self-hosted#setup-github-integration"
              target="_blank"
            >
              Learn how to configure
            </LinkButton>
            <Button v-else-if="!isGitHubEnabled" color="blue" :disabled="true" :show-tooltip="true">
              <template #tooltip> Only available for Sturdy Enterprise </template>
              <template #default> Import from GitHub </template>
            </Button>
            <RouterLinkButton
              v-else
              :to="{
                name: 'organizationCreateGitHubCodebase',
                params: { organizationSlug: selected.shortID },
              }"
              color="blue"
            >
              Import from GitHub
            </RouterLinkButton>
          </div>
        </div>
      </div>
    </div>

    <div class="bg-gray-100 sm:rounded-lg">
      <div class="px-4 py-5 sm:p-6 flex flex-col justify-between h-full">
        <div>
          <h3 class="text-lg leading-6 font-medium text-gray-900">
            Import existing codebase <strong>from any git://</strong> provider
          </h3>
          <div class="mt-2 max-w-xl text-sm text-gray-500">
            <p>
              Sturdy speaks git:// and can connect to any git provider like GitLab, Bitbucket or
              Azure DevOps
            </p>
            <ul class="list-inside mt-2 inline-flex flex-col gap-2">
              <li class="inline-flex space-x-2">
                <CheckIcon class="h-5 w-5 text-green-400 flex-shrink-0" />
                <span>Use Sturdy together with your existing git setup</span>
              </li>
              <li class="inline-flex space-x-2">
                <CheckIcon class="h-5 w-5 text-green-400 flex-shrink-0" />
                <span>Work in Sturdy, create share code when you are done</span>
              </li>
              <li class="inline-flex space-x-2">
                <CheckIcon class="h-5 w-5 text-green-400 flex-shrink-0" />
                <span>Gradually migrate to native Sturdy at your own pase</span>
              </li>
            </ul>
          </div>
        </div>
        <div>
          <div class="mt-5 space-x-2 flex">
            <RouterLinkButton
              v-if="isRemoteAvailable"
              :to="{
                name: 'organizationCreateRemoteCodebase',
                params: { organizationSlug: selected.shortID },
              }"
              color="blue"
            >
              Import from git://
            </RouterLinkButton>
            <Button v-else color="blue" :disabled="true" :show-tooltip="true">
              <template #tooltip> Only available for Sturdy Enterprise </template>
              <template #default> Import from git:// </template>
            </Button>
          </div>
        </div>
      </div>
    </div>

    <div class="bg-gray-100 sm:rounded-lg">
      <div class="px-4 py-5 sm:p-6 flex flex-col justify-between h-full">
        <div>
          <h3 class="text-lg leading-6 font-medium text-gray-900">
            Import <strong>local git</strong> repository
          </h3>
          <div class="mt-2 max-w-xl text-sm text-gray-500">
            <p>Already using git, but don't want to connect your git:// provider yet?</p>
            <ul class="list-inside mt-2 inline-flex flex-col gap-2">
              <li class="inline-flex space-x-2">
                <CheckIcon class="h-5 w-5 text-green-400 flex-shrink-0" />
                <span>Import your existing code</span>
              </li>
              <li class="inline-flex space-x-2">
                <CheckIcon class="h-5 w-5 text-green-400 flex-shrink-0" />
                <span>Host your project on Sturdy</span>
              </li>
            </ul>
          </div>
        </div>
        <div>
          <div class="mt-5 space-x-2 flex">
            <Button :disabled="true" :show-tooltip="true" color="blue">
              <template #tooltip> Coming soon </template>
              <template #default> Import existing </template>
            </Button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { CheckIcon, SelectorIcon } from '@heroicons/vue/solid/esm'
import RouterLinkButton from '../atoms/RouterLinkButton.vue'
import { Feature } from '../__generated__/types'
import { watch, defineProps, inject, computed, type Ref, ref, type PropType } from 'vue'
import LinkButton from '../atoms/LinkButton.vue'
import { gql } from '@urql/vue'
import type { Organization_CreateCodebaseFragment } from './__generated__/CreateCodebase'
import Button from '../atoms/Button.vue'
import { Listbox, ListboxButton, ListboxOption, ListboxOptions } from '@headlessui/vue'
import { useRouter } from 'vue-router'

const props = defineProps({
  organizations: {
    type: Object as PropType<Organization_CreateCodebaseFragment[]>,
    required: true,
  },
  selectedOrganization: {
    type: Object as PropType<Organization_CreateCodebaseFragment>,
    required: false,
  },
})

const features = inject<Ref<Array<Feature>>>('features', ref([]))
const isGitHubEnabled = computed(() => features?.value?.includes(Feature.GitHub))
const isGitHubEnabledNotConfigured = computed(() =>
  features?.value?.includes(Feature.GitHubNotConfigured)
)
const isRemoteAvailable = computed(() => features?.value?.includes(Feature.Remote))

const selected = ref(
  props.selectedOrganization
    ? props.organizations.find(
        (organization) => organization.shortID === props.selectedOrganization?.shortID
      )
    : props.organizations[0]
)

const router = useRouter()
watch(selected, (newValue) => {
  if (newValue) router.push({ params: { organizationSlug: newValue.shortID } })
})
</script>

<script lang="ts">
export const ORGANIZATION_FRAGMENT = gql`
  fragment Organization_CreateCodebase on Organization {
    id
    shortID
    name
  }
`
</script>
