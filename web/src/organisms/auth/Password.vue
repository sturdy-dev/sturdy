<template>
  <form class="mt-8 space-y-6">
    <Banner v-if="error" :message="error" status="error" :show-icon="false" />

    <div class="isolate flex flex-col">
      <template v-if="signUp">
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
      </template>

      <label for="email-address" class="sr-only">Email address</label>
      <input
        id="email-address"
        v-model="email"
        class="appearance-none relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 focus:outline-none focus:z-10 sm:text-sm"
        :class="[
          isEmailWarning
            ? 'bg-yellow-50 focus:ring-yellow-500 focus:border-yellow-500'
            : 'focus:ring-green-500 focus:border-green-500',
          !signUp ? 'rounded-t-md' : '',
        ]"
        name="email"
        type="email"
        autocomplete="email"
        required
        placeholder="Email address"
      />

      <label for="passowrd" class="sr-only">Password</label>
      <input
        id="password"
        v-model="password"
        class="appearance-none relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 focus:outline-none focus:z-10 sm:text-sm focus:ring-green-500 focus:border-green-500"
        :class="[
          isPasswordWarning
            ? 'bg-yellow-50 focus:ring-yellow-500 focus:border-yellow-500'
            : 'focus:ring-green-500 focus:border-green-500',
          !signUp ? 'rounded-b-md' : '',
        ]"
        name="password"
        type="password"
        autocomplete="password"
        required
        placeholder="Password"
      />

      <template v-if="signUp">
        <label for="repeat-passowrd" class="sr-only">Repeat passowrd</label>
        <input
          id="repeat-password"
          v-model="passwordRepeat"
          class="appearance-none relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 focus:outline-none focus:z-10 sm:text-sm rounded-b-md"
          name="repeat-password"
          type="password"
          autocomplete="repeat-password"
          required
          placeholder="Repeat password"
          :class="[
            isPasswordRepeatWarning
              ? 'bg-yellow-50 focus:ring-yellow-500 focus:border-yellow-500'
              : 'focus:ring-green-500 focus:border-green-500',
          ]"
        />
      </template>
    </div>

    <button
      :disabled="!canSubmit"
      type="submit"
      class="disabled:opacity-50 opacity-100 group relative w-full flex justify-center py-2 px-4 border border-transparent text-sm font-medium rounded-md text-white bg-yellow-600 hover:bg-yellow-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-yellow-500"
      @click.stop.prevent="onButtonClick"
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
  </form>
</template>

<script lang="ts">
import http from '../../http'
import { Banner } from '../../atoms'

export default {
  components: {
    Banner,
  },
  props: {
    signUp: {
      type: Boolean,
      default: false,
    },
  },
  emits: ['success'],
  data() {
    return {
      email: '',
      emailActivated: false,

      name: '',
      nameActivated: false,

      password: '',
      passwordActivated: false,

      passwordRepeat: '',
      passwordRepeatActivated: false,

      error: undefined as string | undefined,
    }
  },
  computed: {
    buttonText() {
      return this.signUp ? 'Sign up' : 'Sign in'
    },
    isNameWarning() {
      return this.nameActivated && this.name.length === 0
    },
    isEmailWarning() {
      return this.emailActivated && this.email.indexOf('@') === -1
    },
    isPasswordWarning() {
      return this.passwordActivated && this.password.length === 0
    },
    isPasswordRepeatWarning() {
      return this.passwordRepeatActivated && this.password !== this.passwordRepeat
    },
    canSubmit() {
      const password = this.passwordActivated && !this.isPasswordWarning
      const email = this.emailActivated && !this.isEmailWarning
      const name = this.nameActivated && !this.isNameWarning
      const passwordRepeat = this.passwordRepeatActivated && !this.isPasswordRepeatWarning
      return this.signUp ? email && password && name && passwordRepeat : email && password
    },
  },
  watch: {
    name: function () {
      this.nameActivated = true
      this.error = undefined
    },
    email: function () {
      this.emailActivated = true
      this.error = undefined
    },
    password: function () {
      this.passwordActivated = true
      this.error = undefined
    },
    passwordRepeat: function () {
      this.passwordRepeatActivated = true
      this.error = undefined
    },
  },
  methods: {
    onButtonClick() {
      if (this.signUp) {
        this.signup(this.name, this.email, this.password)
      } else {
        this.login(this.email, this.password)
      }
    },
    async login(email: string, password: string) {
      const resp = await fetch(http.url('v3/auth'), {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          email: email,
          password: password,
        }),
        credentials: 'include',
      })

      if (resp.status === 200) {
        this.$emit('success')
      } else if (resp.status >= 400 && resp.status <= 500) {
        const json = await resp.json()
        this.error = json.error as string
      } else {
        this.error = 'Sorry, something went wrong on our side.'
      }
    },
    async signup(name: string, email: string, password: string) {
      const resp = await fetch(http.url('v3/users'), {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          name: name,
          email: email,
          password: password,
        }),
        credentials: 'include',
      })

      if (resp.status === 200) {
        this.$emit('success')
      } else if (resp.status >= 400 && resp.status <= 500) {
        const json = await resp.json()
        this.error = json.error as string
      } else {
        this.error = 'Sorry, something went wrong on our side.'
      }
    },
  },
}
</script>
