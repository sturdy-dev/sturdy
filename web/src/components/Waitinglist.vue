<template>
  <div>
    <div class="mx-auto max-w-md sm:max-w-3xl lg:max-w-7xl">
      <div
        class="relative rounded-2xl px-6 py-10 bg-yellow-600 overflow-hidden shadow-xl sm:px-12 sm:py-20"
      >
        <div aria-hidden="true" class="absolute inset-0 -mt-72 sm:-mt-32 md:mt-0">
          <svg
            class="absolute inset-0 h-full w-full"
            preserveAspectRatio="xMidYMid slice"
            xmlns="http://www.w3.org/2000/svg"
            fill="none"
            viewBox="0 0 1463 360"
          >
            <path
              class="text-yellow-500 text-opacity-40"
              fill="currentColor"
              d="M-82.673 72l1761.849 472.086-134.327 501.315-1761.85-472.086z"
            />
            <path
              class="text-yellow-700 text-opacity-40"
              fill="currentColor"
              d="M-217.088 544.086L1544.761 72l134.327 501.316-1761.849 472.086z"
            />
          </svg>
        </div>
        <div class="relative">
          <div class="sm:text-center">
            <h2 class="text-3xl font-extrabold !text-white tracking-tight sm:text-4xl">
              Want product news and updates?
            </h2>
            <p class="mt-6 mx-auto max-w-2xl text-lg text-white">
              Sign up for our newsletter to stay up to date.
            </p>
          </div>
          <Banner v-if="done" status="success"
            >👍🏻 &nbsp;&nbsp;You're signed up to our newsletter!</Banner
          >
          <form
            v-else
            action="#"
            class="mt-12 sm:mx-auto sm:max-w-lg sm:flex"
            @submit.stop.prevent="submit"
          >
            <div class="min-w-0 flex-1">
              <label for="cta-email" class="sr-only">Email address</label>
              <input
                id="cta-email"
                v-model="email"
                :disabled="loading || done"
                type="email"
                class="block w-full border border-transparent rounded-md px-5 py-3 text-base text-gray-900 placeholder-gray-500 shadow-sm focus:outline-none focus:border-transparent focus:ring-2 focus:ring-white focus:ring-offset-2 focus:ring-offset-indigo-600"
                placeholder="Enter your email"
              />
            </div>
            <div class="mt-4 sm:mt-0 sm:ml-3">
              <button
                :disabled="loading || done"
                type="submit"
                class="block w-full rounded-md border border-transparent px-5 py-3 bg-yellow-500 text-base font-medium text-white shadow hover:bg-yellow-400 focus:outline-none focus:ring-2 focus:ring-white focus:ring-offset-2 focus:ring-offset-indigo-600 sm:px-10 inline-flex"
              >
                <span class="flex-1">Notify me</span>
                <Spinner v-if="loading" class="ml-2" />
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>
  </div>
</template>
<script>
import http from '../http'
import { Banner } from '../atoms'
import Spinner from './shared/Spinner.vue'

export default {
  components: { Spinner, Banner },
  data() {
    return {
      email: null,
      loading: false,
      done: false,
    }
  },
  methods: {
    submit() {
      this.loading = true

      fetch(http.url('v3/waitinglist'), {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include',
        body: JSON.stringify({ email: this.email }),
      })
        .then(http.checkStatus)
        .then(() => {
          this.done = true
          this.loading = false
        })
    },
  },
}
</script>
