<template>
  <div class="flex items-start justify-center min-h-screen pt-4 px-4 pb-20 text-center sm:p-0">
    <div
      class="inline-block bg-white rounded-lg px-4 pt-5 pb-4 text-left overflow-hidden shadow-xl transform transition-all sm:my-8 sm:align-middle sm:max-w-sm sm:w-full sm:p-6"
    >
      <div>
        <div
          class="mx-auto flex items-center justify-center h-12 w-12 rounded-full"
          :class="[
            failed ? 'bg-red-100' : '',
            unsubscribed ? 'bg-green-100' : '',
            !failed && !unsubscribed ? 'bg-blue-100' : '',
          ]"
        >
          <XIcon v-if="failed" class="h-6 w-6 text-red-600" aria-hidden="true" />
          <CheckIcon v-else-if="unsubscribed" class="h-6 w-6 text-green-600" aria-hidden="true" />
          <Spinner v-else class="h-6 w-6 text-blue-600" aria-hidden="true" />
        </div>
        <div class="mt-3 text-center sm:mt-5">
          <h3 v-if="failed" class="text-lg leading-6 font-medium text-gray-900">Failed</h3>
          <h3 v-else-if="unsubscribed" class="text-lg leading-6 font-medium text-gray-900">
            Unsubscribed
          </h3>
          <h3 v-else class="text-lg leading-6 font-medium text-gray-900">Please wait...</h3>
          <div class="mt-2">
            <p v-if="failed" class="text-sm text-gray-500">
              We where unable to unsubscribe you at this moment. Please try again.
            </p>
            <p v-else-if="unsubscribed" class="text-sm text-gray-500">
              You've been unsubscribed and won't receive any more newsletters from Sturdy.
            </p>
            <p v-else class="text-sm text-gray-500">You'll soon be unsubscribed</p>
          </div>
        </div>
      </div>
      <div class="mt-5 sm:mt-6">
        <a
          href="https://getsturdy.com/"
          class="inline-flex justify-center w-full rounded-md border border-transparent shadow-sm px-4 py-2 bg-indigo-600 text-base font-medium text-white hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 sm:text-sm"
        >
          Go to Sturdy
        </a>
      </div>
    </div>
  </div>
</template>

<script>
import http from '../../http'
import { CheckIcon, XIcon } from '@heroicons/vue/outline'
import Spinner from '../../components/shared/Spinner.vue'

export default {
  components: { CheckIcon, XIcon, Spinner },
  data() {
    return {
      unsubscribed: false,
      failed: false,
    }
  },
  mounted() {
    this.failed = false
    this.unsubscribed = false

    fetch(http.url('v3/unsubscribe'), {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        email: atob(this.$route.params.email),
      }),
      credentials: 'include',
    })
      .then(http.checkStatus)
      .then(() => {
        this.unsubscribed = true
      })
      .catch(() => {
        this.failed = true
      })
  },
}
</script>
