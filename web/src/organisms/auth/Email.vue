<template>
  <form class="mt-8 space-y-6">
    <div class="isolate -space-y-px gap-4 flex flex-col">
      <div>
        <div v-if="!askName && !waitingForEmailCode">
          <label for="name" class="sr-only">Name</label>
          <input
            id="name"
            v-model="name"
            class="appearance-none rounded-t-md relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 focus:outline-none focus:z-10 sm:text-sm"
            :class="[
              isNameWarning
                ? 'bg-yellow-50 focus:ring-yellow-500 focus:border-yellow-500'
                : 'focus:ring-green-500 focus:border-green-500',
            ]"
            name="name"
            type="text"
            autocomplete="name"
            required
            placeholder="Name"
          />
        </div>

        <div v-if="!waitingForEmailCode" class="rounded-md shadow-sm -space-y-px">
          <label for="email-address" class="sr-only">Email address</label>
          <input
            id="email-address"
            v-model="email"
            class="appearance-none relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 focus:outline-none focus:z-10 sm:text-sm"
            :class="[
              isEmailWarning
                ? 'bg-yellow-50 focus:ring-yellow-500 focus:border-yellow-500'
                : 'focus:ring-green-500 focus:border-green-500',
              askName ? 'rounded-md' : 'rounded-b-md',
            ]"
            name="email"
            type="email"
            autocomplete="email"
            required
            placeholder="Email address"
          />
        </div>
      </div>

      <article v-if="showFailedEmailSent" class="">
        <div class="font-medium">
          Sorry, something went wrong on our side while sending an email.
        </div>
      </article>

      <div v-if="waitingForEmailCode" class="rounded text-center">
        <div class="flex flex-col">
          <span>Enter the code you received at</span>
          <span class="font-bold">{{ email }}</span>
        </div>

        <OtpInput
          class="mt-5"
          :input-classes="
            otpInvalid
              ? 'ring-2 ring-red-300 focus:ring-red-300 m-2 h-2 border h-10 w-10 text-center form-control rounded'
              : 'm-2 h-2 border h-10 w-10 text-center form-control rounded'
          "
          :num-inputs="otpLength"
          @complete="verifyOTP"
          @change="verifyOTP"
        />

        <div class="text-gray-500 text-sm mt-12">
          Entered the wrong email?
          <a href="#" class="text-yellow-600" @click="waitingForEmailCode = false">Try again</a>.
        </div>
      </div>

      <button
        v-if="!waitingForEmailCode"
        :disabled="!canSubmit"
        type="submit"
        class="disabled:opacity-50 opacity-100 group relative w-full flex justify-center py-2 px-4 border border-transparent text-sm font-medium rounded-md text-white bg-yellow-600 hover:bg-yellow-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-yellow-500"
        @click.stop.prevent="sendOTP"
      >
        <span class="absolute left-0 inset-y-0 flex items-center pl-3">
          <svg
            class="h-5 w-5 text-yellow-500 group-hover:text-yellow-400"
            xmlns="http://www.w3.org/2000/svg"
            viewBox="0 0 20 20"
            fill="currentColor"
            aria-hidden="true"
          >
            <path
              fill-rule="evenodd"
              d="M5 9V7a5 5 0 0110 0v2a2 2 0 012 2v5a2 2 0 01-2 2H5a2 2 0 01-2-2v-5a2 2 0 012-2zm8-2v2H7V7a3 3 0 016 0z"
              clip-rule="evenodd"
            />
          </svg>
        </span>
        {{ buttonText }}
      </button>

      <p
        v-if="!waitingForEmailCode"
        class="text-sm p-4 text-gray-500 bg-gray-200 rounded-lg inline-flex gap-4"
      >
        <SparklesIcon class="h-5 w-5 text-gray-400" />
        We’ll email you a magic code for a password-free sign in.
      </p>
    </div>
  </form>
</template>

<script lang="ts">
import { SparklesIcon } from '@heroicons/vue/solid'
import { OtpInput } from '../../molecules/auth'
import http from '../../http'

export default {
  components: { SparklesIcon, OtpInput },
  props: {
    askName: {
      type: Boolean,
      default: false,
    },
  },
  emits: ['success'],
  data() {
    return {
      name: '',
      email: this.$route.params.email ?? '',

      showFailedEmailSent: false,
      waitingForEmailCode: false,

      otpLength: 6,
      otpInvalid: false,
    }
  },

  computed: {
    buttonText() {
      return !this.askName ? 'Sign up' : 'Sign in with email'
    },
    canSubmit() {
      if (!this.askName && this.name.length === 0) return false

      const emailValid = this.email !== '' && this.email.indexOf('@') > 0
      return emailValid
    },
    isNameWarning() {
      return this.name.length === 0
    },
    isEmailWarning() {
      if (this.email !== '' && this.email.indexOf('@') > 0) {
        return false
      }
      if (this.email.length === 0) {
        return false
      }
      return true
    },
  },
  methods: {
    verifyOTP(value: string) {
      if (value.length !== this.otpLength) return
      this.otpInvalid = false

      fetch(http.url('v3/auth/magic-link/verify'), {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          email: this.email,
          code: value,
        }),
        credentials: 'include',
      })
        .then(http.checkStatus)
        .then(() => this.$emit('success'))
        .catch(() => {
          this.otpInvalid = true
        })
    },
    sendOTP() {
      this.showFailedEmailSent = false
      fetch(http.url('v3/auth/magic-link/send'), {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          name: this.name,
          email: this.email,
        }),
        credentials: 'include',
      })
        .then(http.checkStatus)
        .then(() => {
          this.waitingForEmailCode = true
        })
        .catch(() => {
          this.showFailedEmailSent = true
        })
    },
  },
}
</script>
