<template>
  <div class="relative">
    <Banner v-if="done" class="rounded-md mt-8" status="success"
      >We will get back to you within two business days!</Banner
    >
    <form
      v-else
      action="#"
      class="mt-8 sm:mx-auto sm:max-w-lg sm:flex"
      @submit.stop.prevent="submit"
    >
      <div class="min-w-0 flex-1">
        <label for="cta-email" class="sr-only">Email address</label>
        <input
          id="cta-email"
          v-model="email"
          :disabled="loading || done"
          type="email"
          class="block w-full border border-transparent rounded-md px-5 py-3 text-base text-gray-900 placeholder-gray-500 shadow-sm focus:outline-none focus:border-transparent focus:ring-2 focus:ring-white focus:ring-offset-2 focus:ring-offset-yellow-600"
          placeholder="Enter your email"
        />
      </div>
      <div class="mt-4 sm:mt-0 sm:ml-3">
        <button
          :disabled="loading || done"
          type="submit"
          class="cursor-pointer select-none block w-full py-3 px-4 rounded-md shadow bg-gradient-to-r from-yellow-100 to-yellow-200 text-yellow-600 font-medium hover:from-yellow-200 hover:to-yellow-300 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-yellow-400 focus:ring-offset-gray-900"
        >
          <span class="flex-1">Request Enterprise</span>
          <Spinner v-if="loading" class="ml-2" />
        </button>
      </div>
    </form>
  </div>
</template>
<script>
import http from '../../../http'
import { Banner } from '../../../atoms'
import Spinner from '../../../components/shared/Spinner.vue'

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

      fetch(http.url('v3/acl-request-enterprise'), {
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
