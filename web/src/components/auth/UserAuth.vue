<template>
  <div>
    <div class="min-h-screen flex items-center justify-center py-12 px-4 sm:px-6 lg:px-8">
      <div class="max-w-md w-full space-y-8">
        <LoginRegister :start-with-sign-up="startWithSignUp" :navigate-to="navigateTo" />
      </div>
    </div>
  </div>
</template>

<script>
import LoginRegister from '../../pages/LoginRegister.vue'

export default {
  components: { LoginRegister },
  props: {
    user: {
      type: Object,
      default: null,
    },
    navigateTo: {
      type: String,
      default: () => '/codebases',
    },
  },
  data() {
    return {
      startWithSignUp: false,
    }
  },
  computed: {
    authenticated() {
      return !!this.user
    },
  },
  watch: {
    $route: {
      immediate: true,
      deep: true,
      handler(newRoute) {
        this.startWithSignUp =
          newRoute.name === 'getStartedGitHub' ||
          newRoute.name === 'signup' ||
          newRoute.name === 'getStartedYC'
      },
    },
    user: function () {
      if (this.authenticated) this.redirect()
    },
  },
  mounted() {
    if (this.authenticated) this.redirect()
  },
  methods: {
    redirect() {
      const queryParam = this.$route.query.navigateTo
      const to = queryParam ? queryParam : this.navigateTo
      this.$router.push(to)
    },
  },
}
</script>
