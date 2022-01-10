<template>
  <PaddedApp>
    <div v-if="data">
      <div class="py-8 px-4">
        <div class="">
          <h2
            class="text-4xl font-extrabold text-gray-900 sm:text-4xl sm:tracking-tight lg:text-4xl"
          >
            Create a new codebase
          </h2>
          <p class="mt-5 text-xl text-gray-500">
            Create your first Sturdy codebase. You'll soon be ready to code!
          </p>
        </div>
      </div>

      <Banner
        v-if="show_failed_message"
        class="my-4"
        status="error"
        message="Could not create a codebase at this time. Try again later!"
      />

      <div class="space-y-6 lg:col-span-9 max-w-4xl">
        <form @submit.stop.prevent="doCreateNewCodebase">
          <div class="shadow sm:rounded-md sm:overflow-hidden">
            <div class="bg-white py-6 px-4 space-y-6 sm:p-6">
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
                      v-model="newCodebaseName"
                      type="text"
                      name="codebase_name"
                      autocomplete="off"
                      class="focus:ring-blue-500 focus:border-blue-500 flex-grow block w-full min-w-0 rounded-md sm:text-sm border-gray-300"
                    />
                  </div>
                </div>

                <div class="col-span-3">
                  <label for="description" class="block text-sm font-medium text-gray-700">
                    Description <span class="text-sm text-gray-500">(optional)</span>
                  </label>
                  <div class="mt-1">
                    <textarea
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
            <div class="px-4 py-3 bg-gray-50 text-right sm:px-6">
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

        <div class="shadow sm:rounded-workspaceService, syncService,md sm:overflow-hidden">
          <NoCodebasesGitHubAuth
            :git-hub-account="data.user.gitHubAccount"
            :git-hub-app="data.gitHubApp"
          />
        </div>
      </div>
    </div>
  </PaddedApp>
</template>

<script>
import http from '../http'
import Banner from '../components/shared/Banner.vue'
import { Slug } from '../slug'
import { toRefs } from 'vue'
import { gql, useQuery } from '@urql/vue'
import NoCodebasesGitHubAuth from '../components/codebase/NoCodebasesGitHubAuth.vue'
import Button from '../components/shared/Button.vue'
import RandomName from '../components/codebase/create/random-name.js'
import PaddedApp from '../layouts/PaddedApp.vue'

export default {
  name: 'CreateEmpty',
  components: { PaddedApp, NoCodebasesGitHubAuth, Banner, Button },
  props: {
    features: {
      type: Array,
      required: true,
    },
  },
  setup(props) {
    const { features } = toRefs(props)
    const isGitHubEnabled = features.value.includes('GitHub')

    const result = useQuery({
      query: gql`
        query CreateEmpty($isGitHubEnabled: Boolean!) {
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

    return {
      fetching: result.fetching,
      data: result.data,
      error: result.error,
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

      let t0 = new Date()

      fetch(http.url('v3/codebases'), {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include',
        body: JSON.stringify({
          name: this.newCodebaseName,
          description: this.newCodebaseDescription,
        }),
      })
        .then(http.checkStatus)
        .then((response) => response.json())
        .then((data) => {
          // Always wait at least 300ms for a nice effect
          let wait = 300 - (new Date() - t0)
          setTimeout(() => {
            this.isLoading = false
            this.$router.push({
              name: 'codebaseHome',
              params: { codebaseSlug: Slug(data.name, data.short_id) },
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
}
</script>
