<template>
  <div class="min-h-screen flex items-center justify-center py-12 px-4 sm:px-6 lg:px-8">
    <div class="max-w-md w-full space-y-8">
      <div>
        <h2 class="mt-6 text-center text-3xl font-extrabold text-gray-900">
          {{ headerText }}
        </h2>
        <p class="mt-2 text-center text-sm text-gray-600">
          Or

          <a
            href="#"
            class="font-medium text-yellow-600 hover:text-yellow-500"
            @click.stop.prevent="isLogin = !isLogin"
          >
            {{ subheaderText }}
          </a>
        </p>
      </div>
      <EmailAuth :askName="isLogin" @success="successRedirect" />
    </div>
  </div>
</template>

<script>
import { EmailAuth } from '../organisms/auth'

export default {
  name: 'LoginRegister',
  components: {
    EmailAuth,
  },
  props: {
    user: { type: Object, default: null },
    startWithSignUp: { type: Boolean, default: false },
    navigateTo: { type: String, default: () => '/codebases' },
  },
  data() {
    return {
      isLogin: !this.startWithSignUp && !this.$route.params.email,
    }
  },
  computed: {
    headerText() {
      return this.isLogin ? 'Sign in to your Sturdy account' : 'Create your Sturdy account'
    },
    subheaderText() {
      return this.isLogin ? 'sign up now' : 'login to your existing account'
    },
  },
  watch: {
    user: {
      immediate: true,
      handler: function (newVal) {
        // User is authenticated, redirect to codebase home
        if (newVal && newVal.id) {
          this.successRedirect()
        }
      },
    },
  },
  methods: {
    async successRedirect() {
      const queryParam = this.$route.query.navigateTo
      const to = queryParam ? queryParam : this.navigateTo
      await this.$router.push(to)
    },
  },
}
</script>
