<template>
  <div v-if="data">
    <slot name="header"></slot>

    <Banner
      v-if="show_failed_message"
      class="my-4"
      status="error"
      message="Could not create a codebase at this time. Try again later!"
    />

    <div class="space-y-6 lg:col-span-9 max-w-4xl">
      <form @submit.stop.prevent="doCreateNewCodebase">
        <div class="shadow sm:rounded-md sm:overflow-hidden">
          <div class="py-6 px-4 space-y-6 sm:p-6" :class="[mainBg]">
            <div class="grid grid-cols-3 gap-6">
              <div class="col-span-3">
                <label for="codebase_name" class="block text-sm font-medium text-gray-700">
                  Codebase name
                </label>
                <div class="mt-1 rounded-md shadow-sm flex">
                  <!--<span class="bg-gray-50 border border-r-0 border-gray-300 rounded-l-md px-3 inline-flex items-center text-gray-500 sm:text-sm">
                    getsturdy.com/
                  </span>-->
                  <input
                    id="codebase_name"
                    v-model="newCodebaseName"
                    type="text"
                    name="codebase_name"
                    autocomplete="off"
                    class="focus:ring-blue-500 focus:border-blue-500 flex-grow block w-full min-w-0 rounded-md sm:text-sm border-gray-300"
                  />
                </div>
              </div>

              <div v-if="false" class="col-span-3">
                <label for="description" class="block text-sm font-medium text-gray-700">
                  Description <span class="text-sm text-gray-500">(optional)</span>
                </label>
                <div class="mt-1">
                  <textarea
                    id="description"
                    v-model="newCodebaseDescription"
                    name="description"
                    rows="3"
                    autocomplete="off"
                    class="shadow-sm focus:ring-blue-500 focus:border-blue-500 mt-1 block w-full sm:text-sm border-gray-300 rounded-md"
                    placeholder="This codebase ..."
                  />
                </div>
              </div>

              <div class="col-span-3">
                <p class="mt-1 text-sm text-gray-500">
                  Codebases are <strong>private</strong>, and only you will be able to see it's
                  contents.
                </p>
                <p class="mt-1 text-sm text-gray-500">
                  You can invite collaborators to the codebase on the next page.
                </p>
              </div>
            </div>
          </div>
          <div class="px-4 py-3 text-right sm:px-6" :class="[bottomBg]">
            <Button type="submit" :disabled="newCodebaseName === ''" size="wider" color="blue">
              Create

              <svg
                v-if="isLoading"
                class="animate-spin ml-3 -mr-1 h-5 w-5 text-white"
                xmlns="http://www.w3.org/2000/svg"
                fill="none"
                viewBox="0 0 24 24"
              >
                <circle
                  class="opacity-25"
                  cx="12"
                  cy="12"
                  r="10"
                  stroke="currentColor"
                  stroke-width="4"
                />
                <path
                  class="opacity-75"
                  fill="currentColor"
                  d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
                />
              </svg>
            </Button>
          </div>
        </div>
      </form>

      <div
        v-if="data.gitHubApp && showSetupGitHub"
        class="shadow sm:rounded-workspaceService, syncService,md sm:overflow-hidden"
      >
        <NoCodebasesGitHubAuth
          :git-hub-account="data.user.gitHubAccount"
          :git-hub-app="data.gitHubApp"
        />
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { Banner } from '../atoms'
import { defineComponent, inject, ref } from 'vue'
import type { Ref } from 'vue'
import { gql, useQuery } from '@urql/vue'
import NoCodebasesGitHubAuth from '../components/codebase/NoCodebasesGitHubAuth.vue'
import Button from '../atoms/Button.vue'
import RandomName from '../components/codebase/create/random-name.js'
import { useCreateCodebase } from '../mutations/useCreateCodebase'
import type {
  CreateCodebasePageQuery,
  CreateCodebasePageQueryVariables,
} from './__generated__/CreateCodebase'
import { Feature } from '../__generated__/types'
import { Slug } from '../slug'

export default defineComponent({
  components: { NoCodebasesGitHubAuth, Banner, Button },
  props: {
    createInOrganizationId: {
      type: String,
      required: false,
    },
    showSetupGitHub: {
      type: Boolean,
      required: false,
    },
    mainBg: {
      type: String,
      required: false,
      default: 'bg-white',
    },
    bottomBg: {
      type: String,
      required: false,
      default: 'bg-gray-50',
    },
  },
  setup() {
    const features = inject<Ref<Array<Feature>>>('features', ref([]))
    const isGitHubEnabled = features.value.includes(Feature.GitHub)

    const result = useQuery<CreateCodebasePageQuery, CreateCodebasePageQueryVariables>({
      query: gql`
        query CreateCodebasePage($isGitHubEnabled: Boolean!) {
          user {
            id
            gitHubAccount @include(if: $isGitHubEnabled) {
              id
              login
            }
          }
          gitHubApp @include(if: $isGitHubEnabled) {
            _id
            name
            clientID
          }
        }
      `,
      variables: {
        isGitHubEnabled,
      },
    })

    const createCodebaseResult = useCreateCodebase()

    return {
      fetching: result.fetching,
      data: result.data,
      error: result.error,

      async createCodebase(name: string, organizationID: string | undefined) {
        return createCodebaseResult({ name, organizationID })
      },
    }
  },
  data() {
    return {
      newCodebaseName: RandomName.generate(),
      newCodebaseDescription: '',
      isLoading: false,
      show_failed_message: false,
    }
  },
  watch: {
    error: function (err) {
      if (err) throw err
    },
  },
  methods: {
    doCreateNewCodebase() {
      this.isLoading = true
      this.show_failed_message = false

      let t0 = +new Date()

      this.createCodebase(this.newCodebaseName, this.createInOrganizationId)
        .then((result) => {
          // Always wait at least 300ms for a nice effect
          let wait = 300 - (+new Date() - t0)
          setTimeout(() => {
            this.isLoading = false

            this.$router.push({
              name: 'codebaseHome',
              params: {
                codebaseSlug: Slug(result.createCodebase.name, result.createCodebase.shortID),
              },
            })
          }, wait)
        })
        .catch(() => {
          this.isLoading = false
          this.show_failed_message = true
          setTimeout(() => {
            this.show_failed_message = false
          }, 8000)
        })
    },
  },
})
</script>
