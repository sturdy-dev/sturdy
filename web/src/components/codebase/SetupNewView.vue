<template>
  <div class="bg-white border overflow-hidden sm:rounded-lg">
    <div class="px-4 py-5 sm:px-6 flex justify-between bg-blue-50">
      <div>
        <h3 class="text-lg leading-6 font-medium text-gray-900">Let's get started!</h3>
        <p class="mt-1 max-w-2xl text-sm text-gray-500">
          Setup {{ codebase.name }} on your computer
        </p>
      </div>
    </div>
    <div class="border-t border-gray-200 px-4 py-5 sm:px-6 text-sm">
      <nav aria-label="Progress">
        <ol role="list" class="overflow-hidden">
          <li
            v-for="(step, stepIdx) in steps"
            :key="step.name"
            :class="[stepIdx !== steps.length - 1 ? 'pb-10' : '', 'relative']"
          >
            <template v-if="step.status === 'complete'">
              <div
                v-if="stepIdx !== steps.length - 1"
                class="-ml-px absolute mt-0.5 top-4 left-4 w-0.5 h-full bg-blue-600"
                aria-hidden="true"
              />
              <a class="relative flex items-start group">
                <span class="h-9 flex items-center">
                  <span
                    class="relative z-10 w-8 h-8 flex items-center justify-center bg-blue-600 rounded-full group-hover:bg-blue-800"
                  >
                    <CheckIcon class="w-5 h-5 text-white" aria-hidden="true" />
                  </span>
                </span>
                <span class="ml-4 flex flex-col w-full">
                  <span class="text-xs font-semibold tracking-wide uppercase">{{ step.name }}</span>
                  <span class="text-sm text-gray-500">{{ step.description }}</span>
                  <div class="mt-2">
                    <component
                      :is="step.component"
                      v-if="step.component"
                      class="w-full"
                      :codebase="codebase"
                      :user="data?.user"
                      :codebase-slug="codebaseSlug"
                      :codebase-id="codebase.id"
                    />
                  </div>
                </span>
              </a>
            </template>
            <template v-else-if="step.status === 'current'">
              <div
                v-if="stepIdx !== steps.length - 1"
                class="-ml-px absolute mt-0.5 top-4 left-4 w-0.5 h-full bg-gray-300"
                aria-hidden="true"
              />
              <a class="relative flex items-start group" aria-current="step">
                <span class="h-9 flex items-center" aria-hidden="true">
                  <span
                    class="relative z-10 w-8 h-8 flex items-center justify-center bg-white border-2 border-blue-600 rounded-full"
                  >
                    <span class="h-2.5 w-2.5 bg-blue-600 rounded-full" />
                  </span>
                </span>
                <span class="ml-4 flex flex-col w-full">
                  <span class="text-xs font-semibold tracking-wide uppercase text-blue-600">
                    {{ step.name }}
                  </span>
                  <span class="text-sm text-gray-500">{{ step.description }}</span>
                  <div class="mt-2">
                    <component
                      :is="step.component"
                      v-if="step.component"
                      :codebase="codebase"
                      :user="data?.user"
                      :codebase-slug="codebaseSlug"
                      :codebase-id="codebase.id"
                    />
                  </div>
                </span>
              </a>
            </template>
            <template v-else>
              <div
                v-if="stepIdx !== steps.length - 1"
                class="-ml-px absolute mt-0.5 top-4 left-4 w-0.5 h-full bg-gray-300"
                aria-hidden="true"
              />
              <a class="relative flex items-start group">
                <span class="h-9 flex items-center" aria-hidden="true">
                  <span
                    class="relative z-10 w-8 h-8 flex items-center justify-center bg-white border-2 border-gray-300 rounded-full group-hover:border-gray-400"
                  >
                    <span class="h-2.5 w-2.5 bg-transparent rounded-full group-hover:bg-gray-300" />
                  </span>
                </span>
                <span class="ml-4 flex flex-col w-full">
                  <span class="text-xs font-semibold tracking-wide uppercase text-gray-500">
                    {{ step.name }}
                  </span>
                  <span class="text-sm text-gray-500">{{ step.description }}</span>
                  <div class="mt-2">
                    <component
                      :is="step.component"
                      v-if="step.component"
                      :codebase="codebase"
                      :user="data?.user"
                      :codebase-slug="codebaseSlug"
                      :codebase-id="codebase.id"
                    />
                  </div>
                </span>
              </a>
            </template>
          </li>
        </ol>
      </nav>
    </div>
  </div>
</template>

<script>
import Button from '../shared/Button.vue'
import { CheckIcon, DownloadIcon } from '@heroicons/vue/solid'
import SetupSturdyInitStep from './SetupSturdyInitStep.vue'
import SetupSturdyInstallStep from './SetupSturdyInstallStep.vue'
import { gql, useQuery } from '@urql/vue'
import SetupSturdyGoToWorkspaceStep, { SETUP_USER_VIEWS } from './SetupSturdyGoToWorkspaceStep.vue'
import CreateViewAndWorkspace from './CreateViewAndWorkspace.vue'

export default {
  name: 'SetupNewView',
  components: { Button, DownloadIcon, CheckIcon },
  props: {
    codebase: {
      type: Object,
      required: true,
    },
    codebaseSlug: {
      type: String,
      required: true,
    },
    currentUserHasAView: {
      type: Boolean,
      required: true,
    },
  },
  setup() {
    const { data } = useQuery({
      query: gql`
        query SetupNewView {
          user {
            id
            views {
              id
            }
            ...SetupUserViews
          }
        }
        ${SETUP_USER_VIEWS}
      `,
      requestPolicy: 'cache-and-network',
    })
    return {
      data,
    }
  },
  computed: {
    isApp() {
      return !!window.ipc
    },
    haveAnyViewsAnyCodebase() {
      return this.data?.user?.views.length > 0
    },
    currentStep() {
      const codebaseHasChanges = this.codebase.changes.length > 0
      if (this.isApp) {
        // app steps
        if (this.currentUserHasAView && codebaseHasChanges) {
          return 2
        } else if (this.currentUserHasAView) {
          return 1
        }
        return 0
      } else {
        // web steps
        const codebaseHasChanges = this.codebase.changes.length > 0
        const visitedInstallationPage = localStorage.getItem('visitedInstallClient')
        if (this.currentUserHasAView && codebaseHasChanges) {
          return 3
        } else if (this.currentUserHasAView) {
          return 2
        } else if (this.haveAnyViewsAnyCodebase || visitedInstallationPage) {
          return 1
        }
        return 0
      }
    },
    steps() {
      if (this.isApp) {
        return [
          {
            name: 'Setup directory',
            description: 'Connect Sturdy app to a local directory',
            status:
              this.currentStep === 0 ? 'current' : this.currentStep > 1 ? 'complete' : 'upcoming',
            component: CreateViewAndWorkspace,
          },
          {
            name: 'Start coding',
            description: 'Make your first change to the codebase',
            status:
              this.currentStep === 1 ? 'current' : this.currentStep > 2 ? 'complete' : 'upcoming',
            component: SetupSturdyGoToWorkspaceStep,
          },
        ]
      }
      return [
        {
          name: 'Install Sturdy',
          description: 'Install the Sturdy app on your computer',
          status:
            this.currentStep === 0 ? 'current' : this.currentStep > 0 ? 'complete' : 'upcoming',
          component: SetupSturdyInstallStep,
        },
        {
          name: 'Setup directory',
          description: 'Run this command to connect this codebase to a directory',
          status:
            this.currentStep === 1 ? 'current' : this.currentStep > 1 ? 'complete' : 'upcoming',
          component: SetupSturdyInitStep,
        },
        {
          name: 'Start coding',
          description: 'Make your first change to the codebase',
          status:
            this.currentStep === 2 ? 'current' : this.currentStep > 2 ? 'complete' : 'upcoming',
          component: SetupSturdyGoToWorkspaceStep,
        },
      ]
    },
  },
}
</script>
