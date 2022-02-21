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
      <EmailAuth v-if="isEmailAuthEnabled" :ask-name="isLogin" @success="successRedirect" />
      <PasswordAuth v-else :sign-up="!isLogin" @success="successRedirect" />
    </div>
  </div>
</template>

<script lang="ts">
import { EmailAuth, PasswordAuth } from '../organisms/auth'
import { Feature } from '../__generated__/types'
import { computed, defineComponent, inject, ref, Ref, PropType, watch } from 'vue'
import { gql, useQuery } from '@urql/vue'

export default defineComponent({
  components: {
    EmailAuth,
    PasswordAuth,
  },
  props: {
    user: { type: Object, default: null },
    startWithSignUp: { type: Boolean, default: false },
    navigateTo: { type: String, default: () => '/codebases' },
    features: { type: Array as PropType<Feature[]>, required: true },
  },
  setup() {
    const features = inject<Ref<Array<Feature>>>('features', ref([]))
    const isEmailAuthEnabled = computed(() => features?.value?.includes(Feature.Emails))

    const isMultiTenancyEnabled = ref(false)
    watch(features, (newVal) => {
      isMultiTenancyEnabled.value = newVal?.includes(Feature.MultiTenancy)
    })

    const { data } = useQuery({
      query: gql`
        query LoginRegister($isMultiTenancyEnabled: Boolean!) {
          installation @skip(if: $isMultiTenancyEnabled) {
            id
            usersCount
          }
        }
      `,
      variables: {
        isMultiTenancyEnabled: isMultiTenancyEnabled,
      },
    })

    return {
      data,
      isEmailAuthEnabled,
    }
  },
  data() {
    return {
      isLogin: !this.startWithSignUp && !this.$route.params.email,
    }
  },
  computed: {
    usersCount() {
      return this.data?.installation?.usersCount
    },
    headerText() {
      return this.isLogin ? 'Sign in to your Sturdy account' : 'Create your Sturdy account'
    },
    subheaderText() {
      return this.isLogin ? 'sign up now' : 'login to your existing account'
    },
  },
  watch: {
    usersCount: function (newValue) {
      // force show register form if there are no users
      if (newValue === 0) this.isLogin = false
    },
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
      const to = queryParam ? (queryParam as string) : this.navigateTo
      await this.$router.push(to)
    },
  },
})
</script>
