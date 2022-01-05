<template>
  <div>
    <div class="min-h-screen flex items-center justify-center py-12 px-4 sm:px-6 lg:px-8">
      <div v-if="invalid_code" class="max-w-md w-full space-y-8">
        <Banner status="error" message="Sorry, this code is invalid or has expired." />
        <router-link
          :to="{ name: 'codebaseOverview' }"
          type="submit"
          class="disabled:opacity-50 opacity-100 group relative w-full flex justify-center py-2 px-4 border border-transparent text-sm font-medium rounded-md text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
        >
          Go to codebases
        </router-link>
      </div>

      <div v-if="codebase" class="max-w-md w-full space-y-8">
        <div>
          <h2 class="mt-6 text-center text-3xl font-extrabold text-gray-900">
            You've been invited to join {{ codebase.name }}
          </h2>
        </div>

        <Banner
          v-if="failed_to_join"
          status="error"
          message="You where not able to join the codebase. Please try again later."
        />

        <template v-if="user">
          <button
            type="submit"
            class="disabled:opacity-50 opacity-100 group relative w-full flex justify-center py-2 px-4 border border-transparent text-sm font-medium rounded-md text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
            @click="joinNow"
          >
            <span class="absolute left-0 inset-y-0 flex items-center pl-3">
              <svg
                class="h-5 w-5 text-blue-500 group-hover:text-blue-400"
                xmlns="http://www.w3.org/2000/svg"
                viewBox="0 0 20 20"
                fill="currentColor"
              >
                <path
                  fill-rule="evenodd"
                  d="M6.672 1.911a1 1 0 10-1.932.518l.259.966a1 1 0 001.932-.518l-.26-.966zM2.429 4.74a1 1 0 10-.517 1.932l.966.259a1 1 0 00.517-1.932l-.966-.26zm8.814-.569a1 1 0 00-1.415-1.414l-.707.707a1 1 0 101.415 1.415l.707-.708zm-7.071 7.072l.707-.707A1 1 0 003.465 9.12l-.708.707a1 1 0 001.415 1.415zm3.2-5.171a1 1 0 00-1.3 1.3l4 10a1 1 0 001.823.075l1.38-2.759 3.018 3.02a1 1 0 001.414-1.415l-3.019-3.02 2.76-1.379a1 1 0 00-.076-1.822l-10-4z"
                  clip-rule="evenodd"
                />
              </svg>
            </span>
            Join now!
          </button>
        </template>
        <LoginRegister v-else :redirect-to="redirectTo" :start-with-sign-up="true" />
      </div>
    </div>
  </div>
</template>

<script>
import http from '../../http'
import LoginRegister from '../../pages/LoginRegister.vue'
import Banner from '../shared/Banner.vue'
import { Slug } from '../../slug'

export default {
  name: 'Join',
  components: { Banner, LoginRegister },
  props: ['user'],
  data() {
    return {
      codebase: null,
      redirectTo: this.$route.fullPath,
      invalid_code: false,
      failed_to_join: false,
    }
  },
  mounted() {
    this.getCodebase()
  },
  methods: {
    getCodebase() {
      fetch(http.url('v3/join/get-codebase/' + this.$route.params.code), { credentials: 'include' })
        .then(http.checkStatus)
        .then((response) => response.json())
        .then((data) => {
          this.codebase = data
        })
        .catch(() => {
          this.invalid_code = true
        })
    },
    joinNow() {
      fetch(http.url('v3/join/codebase/' + this.$route.params.code), {
        method: 'POST',
        credentials: 'include',
      })
        .then(http.checkStatus)
        .then((response) => response.json())
        .then((data) => {
          this.$router.push({
            name: 'codebaseHome',
            params: { codebaseSlug: Slug(data.name, data.short_id) },
          })
        })
        .catch(() => {
          this.failed_to_join = true
        })
    },
  },
}
</script>

<style scoped></style>
