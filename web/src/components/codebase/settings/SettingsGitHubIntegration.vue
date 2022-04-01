<template>
  <div v-if="gitHubIntegration">
    <HorizontalDivider class="mt-4" bg="bg-white">GitHub settings</HorizontalDivider>
    <div class="mx-4 divide-y divide-gray-200">
      <ul v-if="gitHubIntegration" class="mt-2 divide-y divide-gray-200">
        <Banner v-if="failedToUpdateGitHubSetting" status="error">
          Failed to update the GitHub integration. Please try again later!
        </Banner>

        <li
          class="py-4 flex items-start justify-between flex-col flex-col xl:flex-row space-y-4 xl:space-x-4 xl:space-y-0"
        >
          <div class="flex flex-col pr-8 text-sm text-gray-500 space-y-2">
            <p class="text-sm font-medium text-gray-900">GitHub</p>
            <p class="">
              This codebase is connected to
              <span class="underline">{{ gitHubIntegration.trackedBranch }}</span>
              on
              <a :href="github_repo_url" class="underline whitespace-nowrap">
                {{ gitHubIntegration.owner }}/{{ gitHubIntegration.name }} </a
              >.
            </p>
            <p>
              The settings control if GitHub or Sturdy should be the <i>source of truth</i>. Sturdy
              must be the source of truth to enable the entire Sturdy workflow.
            </p>
            <Banner v-if="gitHubIntegration.lastPushErrorMessage" status="error">
              The last push to GitHub failed:
              {{ gitHubIntegration.lastPushErrorMessage }}
            </Banner>
            <p v-if="gitHubIntegration.lastPushAt">
              Synced with GitHub {{ friendly_ago(gitHubIntegration.lastPushAt) }}.
            </p>
          </div>
          <RadioGroup v-model="selectedGitHubSetting" class="w-full md:w-92">
            <RadioGroupLabel class="sr-only"> Privacy setting </RadioGroupLabel>
            <div class="bg-white rounded-md -space-y-px">
              <RadioGroupOption
                v-for="(setting, settingIdx) in gitHubIntegrationSettings"
                :key="setting.name"
                v-slot="{ checked, active }"
                as="template"
                :value="setting"
              >
                <div
                  :class="[
                    settingIdx === 0 ? 'rounded-tl-md rounded-tr-md' : '',
                    settingIdx === gitHubIntegrationSettings.length - 1
                      ? 'rounded-bl-md rounded-br-md'
                      : '',
                    checked ? 'bg-blue-50 border-blue-200 z-10' : 'border-gray-200',
                    'relative border p-4 flex cursor-pointer focus:outline-none',
                  ]"
                >
                  <span
                    :class="[
                      checked ? 'bg-blue-600 border-transparent' : 'bg-white border-gray-300',
                      active ? 'ring-2 ring-offset-2 ring-blue-500' : '',
                      'h-4 w-4 mt-0.5 cursor-pointer rounded-full border flex items-center justify-center flex-shrink-0',
                    ]"
                    aria-hidden="true"
                  >
                    <span class="rounded-full bg-white w-1.5 h-1.5" />
                  </span>
                  <div class="ml-3 flex flex-col">
                    <RadioGroupLabel
                      as="span"
                      :class="[
                        checked ? 'text-blue-900' : 'text-gray-900',
                        'block text-sm font-medium',
                      ]"
                    >
                      {{ setting.name }}
                    </RadioGroupLabel>
                    <RadioGroupDescription
                      as="span"
                      :class="[checked ? 'text-blue-700' : 'text-gray-500', 'block text-sm']"
                    >
                      {{ setting.description }}
                    </RadioGroupDescription>
                  </div>
                </div>
              </RadioGroupOption>
            </div>
          </RadioGroup>
        </li>
      </ul>
    </div>
  </div>
</template>

<script>
import { gql, useMutation } from '@urql/vue'
import { toRef, ref, watch } from 'vue'
import { Banner } from '../../../atoms'
import HorizontalDivider from '../../../atoms/HorizontalDivider.vue'
import {
  RadioGroup,
  RadioGroupDescription,
  RadioGroupLabel,
  RadioGroupOption,
} from '@headlessui/vue'
import time from '../../../time'

const gitHubIntegrationSettings = [
  {
    name: 'GitHub is the source of truth',
    description:
      'Workspaces on Sturdy are merged through GitHub pull requests. When PRs are merged on GitHub, they will be synced to Sturdy. This mode allows you to incrementally migrate to Sturdy.',
    vals: { enabled: true, gitHubIsSourceOfTruth: true },
  },
  {
    name: 'Sturdy is the source of truth',
    description:
      'Workspaces are merged directly on Sturdy. After a change is created on Sturdy, it will be pushed to the HEAD branch on GitHub. Changes made directly on GitHub will not be synced to Sturdy.',
    vals: { enabled: true, gitHubIsSourceOfTruth: false },
  },
  {
    name: 'Disabled',
    description: 'Disable syncing between Sturdy and GitHub',
    vals: { enabled: false, gitHubIsSourceOfTruth: false },
  },
]

export default {
  name: 'SettingsGitHubIntegration',
  components: {
    Banner,
    HorizontalDivider,
    RadioGroup,
    RadioGroupDescription,
    RadioGroupLabel,
    RadioGroupOption,
  },
  props: {
    gitHubIntegration: Object,
  },
  setup(props) {
    const { executeMutation: updateCodebaseGitHubIntegrationResult } = useMutation(gql`
      mutation SettingsGitHubIntegration(
        $id: ID!
        $enabled: Boolean
        $gitHubIsSourceOfTruth: Boolean
      ) {
        updateCodebaseGitHubIntegration(
          input: { id: $id, enabled: $enabled, gitHubIsSourceOfTruth: $gitHubIsSourceOfTruth }
        ) {
          id
          enabled
          gitHubIsSourceOfTruth
        }
      }
    `)

    let gitHubIntegration = toRef(props, 'gitHubIntegration')

    let getSelected = () => {
      if (!gitHubIntegration.value) {
        return null
      }
      return gitHubIntegrationSettings.filter(
        (s) =>
          s.vals.enabled === gitHubIntegration.value.enabled &&
          s.vals.gitHubIsSourceOfTruth === gitHubIntegration.value.gitHubIsSourceOfTruth
      )[0]
    }

    const selectedGitHubSetting = ref(getSelected())
    watch(gitHubIntegration, () => {
      selectedGitHubSetting.value = getSelected()
    })

    const isUpdatingGitHubSetting = ref(false)
    const failedToUpdateGitHubSetting = ref(false)

    watch(selectedGitHubSetting, (newSetting) => {
      if (!newSetting) {
        return
      }

      isUpdatingGitHubSetting.value = true
      failedToUpdateGitHubSetting.value = false

      const variables = {
        id: gitHubIntegration.value.id,
        enabled: newSetting.vals.enabled,
        gitHubIsSourceOfTruth: newSetting.vals.gitHubIsSourceOfTruth,
      }
      updateCodebaseGitHubIntegrationResult(variables)
        .then((result) => {
          if (result.error) {
            throw new Error(result.error)
          }
        })
        .then(() => {
          isUpdatingGitHubSetting.value = false
        })
        .catch(() => {
          failedToUpdateGitHubSetting.value = true
        })
    })

    return {
      gitHubIntegrationSettings,
      updateCodebaseGitHubIntegrationResult,
      selectedGitHubSetting,
      isUpdatingGitHubSetting,
      failedToUpdateGitHubSetting,
    }
  },
  computed: {
    github_repo_url() {
      return (
        'https://github.com/' + this.gitHubIntegration.owner + '/' + this.gitHubIntegration.name
      )
    },
  },
  methods: {
    friendly_ago(ts) {
      return time.getRelativeTime(new Date(ts * 1000))
    },
  },
}
</script>
