<template>
  <AppRedirect>
    <div class="bg-gray-100 min-h-screen flex flex-col">
      <template v-if="haveSelfContainedLayout">
        <router-view :user="user" />
      </template>
      <template v-else-if="!isApp">
        <ClientOnly>
          <IndexNavbar v-if="showNavigation" :user="user" />
          <AppTitleBarSpacer
            fixed
            pad-right="1rem"
            class="bg-gray-50 border-b left-0 right-0"
            :class="[appEnvironment?.platform === 'win32' ? 'border-black' : '']"
          >
            <div class="flex-1 flex flex-row justify-between gap-8 items-center">
              <!-- Navigation buttons when the sidebar is hidden -->
              <AppTitleBarSpacer v-slot="{ ipc }" fixed pad-left="1rem">
                <AppHistoryNavigationButtons :ipc="ipc" class="flex md:hidden" />
              </AppTitleBarSpacer>
            </div>
          </AppTitleBarSpacer>
        </ClientOnly>

        <router-view :user="user" :features="features" />
        <IndexFooter v-if="showNavigation" />
      </template>
      <template v-else>
        <!-- Bottom section -->
        <ClientOnly>
          <div class="z-0 flex-1 flex">
            <!-- Narrow sidebar-->
            <StackedMenu
              v-if="showSidebar"
              class="hidden md:flex w-64"
              :user="user"
              :features="features"
              @logout="logout"
            />

            <!-- Main area -->
            <main class="md:pl-64 flex flex-1 flex-col">
              <AppTitleBarSpacer
                fixed
                pad-right="1rem"
                class="bg-gray-50 border-b left-0 md:left-64 right-0"
                :class="[appEnvironment?.platform === 'win32' ? 'border-black' : '']"
              >
                <div class="flex-1 flex flex-row justify-between gap-8 items-center">
                  <!-- Navigation buttons when the sidebar is hidden -->
                  <AppTitleBarSpacer v-slot="{ ipc }" pad-left="1rem">
                    <AppHistoryNavigationButtons :ipc="ipc" class="flex md:hidden" />
                  </AppTitleBarSpacer>

                  <AppMutagenStatus class="flex-1" />
                  <AppShareButton class="flex-shrink-0" />
                </div>
              </AppTitleBarSpacer>

              <!-- Primary column -->
              <section
                class="flex-1 flex flex-col overflow-x-auto"
                :class="[appEnvironment ? 'spacer-padding' : '']"
              >
                <router-view v-if="showRoute" :user="user" class="flex-1" />
                <Error v-else-if="error" :error="error" @reset-error="error = null" />
                <ComingSoon v-else class="pt-2 px-2" />
              </section>
            </main>
          </div>
        </ClientOnly>
      </template>

      <ClientOnly>
        <div class="fixed inset-0 flex items-start justify-end px-4 py-6 pointer-events-none z-10">
          <!-- Top margin so that the notification is placed below the "current view" banner -->
          <div class="max-w-sm w-full lg:mt-24 self-end">
            <ToastNotification
              v-for="n in toastNotifications"
              :id="n.id"
              :key="n.id"
              :title="n.title"
              :style="n.style"
              :message="n.message"
            />
          </div>
        </div>
      </ClientOnly>
    </div>
  </AppRedirect>

  <div id="teleported-position"></div>

  <ClientOnly>
    <Onboarding />
  </ClientOnly>
</template>

<script lang="ts">
import { computed, defineComponent, provide, watch } from 'vue'
import http from './http'
import posthog from 'posthog-js'
import ToastNotification from './components/ToastNotification.vue'
import { uuidv4 } from './uuid'
import IndexNavbar from './components/IndexNavbar.vue'
import { CombinedError, gql, useQuery } from '@urql/vue'
import ComingSoon from './components/ComingSoon.vue'
import IndexFooter from './components/IndexFooter.vue'
import { useHead } from '@vueuse/head'
import { ClientOnly } from 'vite-ssr/vue'
import * as Sentry from '@sentry/browser'
import { Integrations } from '@sentry/tracing'
import StackedMenu from './components/menu/StackedMenu.vue'
import Error from './components/Error.vue'
import Onboarding from './components/onboarding/Onboarding.vue'
import AppTitleBarSpacer from './components/AppTitleBarSpacer.vue'
import AppHistoryNavigationButtons from './components/AppHistoryNavigationButtons.vue'
import AppShareButton from './components/AppShareButton.vue'
import AppMutagenStatus from './components/AppMutagenStatus.vue'
import AppRedirect from './components/AppRedirect.vue'
import { RouteLocationNormalizedLoaded, useRoute, useRouter } from 'vue-router'
import { User, Feature } from './__generated__/types'
import {
  UserQueryQuery,
  UserQueryQueryVariables,
  FeaturesQuery,
  FeaturesQueryVariables,
} from './__generated__/Sturdy'

type ToastNotificationMessage = {
  id: string
  style: string
  title: string
  message: string
}

export default defineComponent({
  components: {
    AppRedirect,
    AppShareButton,
    AppHistoryNavigationButtons,
    AppTitleBarSpacer,
    AppMutagenStatus,
    Onboarding,
    StackedMenu,
    IndexFooter,
    ComingSoon,
    IndexNavbar,
    ToastNotification,
    ClientOnly,
    Error,
  },
  setup() {
    let route = useRoute()
    useHead({
      title: 'Sturdy - Code collaboration',
      meta: [
        {
          property: 'description',
          content:
            'Collaborate on code without git. Sturdy is free and works great for school projects, startups, and teams.',
        },
        {
          property: 'keywords',
          content: 'sturdy git collaborate group projects team startup free github gitlab',
        },
        { property: 'og:title', content: 'Sturdy' },
        { property: 'og:url', content: 'https://getsturdy.com/' },
        { property: 'og:type', content: 'website' },
        { property: 'og:locale', content: 'en_US' },
        {
          property: 'og:description',
          content: 'Sturdy - Software collaboration. Reimagined.',
        },
        {
          property: 'og:image',
          content: 'https://getsturdy.com/assets/Colour.png',
        },
      ],
      link: [
        {
          rel: 'canonical',
          href: computed(() => 'https://getsturdy.com' + route.fullPath),
        },
      ],
    })

    if (import.meta.env.SSR) {
      return { data: null, fetching: false }
    }

    const router = useRouter()

    function onChangeRoute(currentRoute: RouteLocationNormalizedLoaded) {
      if (window.ipc && currentRoute.meta.nonApp && !currentRoute.meta.isAuth) {
        if (currentRoute.path !== '/') {
          window.open(new URL(currentRoute.path, location.href).href)
        }
        // Use replace instead of push to make it possible to use the browser back to "skip" over this broken route
        router.replace('/login')
      }
    }

    watch(route, onChangeRoute)
    onChangeRoute(route)

    // TODO: If we can make user optional (in the GraphQL schema), the features and user query could be merged to the same query
    const { data: featuresData } = useQuery<FeaturesQuery, FeaturesQueryVariables>({
      query: gql`
        query Features {
          features
        }
      `,
      requestPolicy: 'cache-and-network',
    })

    provide(
      'features',
      computed(() => featuresData.value?.features)
    )

    const { data, fetching, executeQuery } = useQuery<UserQueryQuery, UserQueryQueryVariables>({
      query: gql`
        query UserQuery {
          user {
            id
            name
            avatarUrl
            email
            notificationPreferences {
              type
              channel
              enabled
            }
          }
        }
      `,
      requestPolicy: 'cache-and-network',
    })

    provide(
      'user',
      computed(() => data.value?.user)
    )

    return {
      data,
      fetching,
      refreshUser: async () => {
        await executeQuery({
          requestPolicy: 'network-only',
        })
      },
      appEnvironment: window.appEnvironment,
      featuresData,
    }
  },
  data(): {
    postHogEnabled: boolean
    toastNotifications: ToastNotificationMessage[]
    error: CombinedError | null
  } {
    return {
      postHogEnabled: false,
      toastNotifications: [] as ToastNotificationMessage[],
      error: null,
    }
  },
  computed: {
    features(): Feature[] {
      return this.featuresData?.features ?? []
    },
    user(): User | null {
      return this.data?.user
    },
    authenticated(): boolean {
      return !!this.user
    },
    authenticationIsLoaded(): boolean {
      if (!this.fetching) {
        return true
      }
      return false
    },
    showRoute(): boolean {
      return this.authenticationIsLoaded && !this.isNotFound
    },
    isUnauthenticated(): boolean {
      return this.errIsUnauthenticated(this.error)
    },
    isForbidden(): boolean {
      return this.errIsForbidden(this.error)
    },
    isNotFound(): boolean {
      return this.errIsNotFound(this.error)
    },
    showSidebar(): boolean {
      return this.isApp && !this.isAuthPage && !this.error
    },
    showNavigation(): boolean {
      const r = this.$route
      const hideNavigation = r.meta && r.meta.hideNavigation
      return !hideNavigation && !this.error
    },
    isApp(): boolean {
      const r = this.$route
      const nonApp = r.meta && r.meta.nonApp
      return !nonApp
    },
    isAuthPage(): boolean {
      const r = this.$route
      return r.meta && !!r.meta.isAuth
    },
    haveSelfContainedLayout(): boolean {
      const r = this.$route
      return r.meta && !!r.meta.selfContainedLayout
    },
  },
  watch: {
    'data.user.id': function () {
      if (this.data?.user?.id && this.postHogEnabled) {
        posthog.identify(this.data.user.id)
      }
    },
  },
  unmounted() {
    this.emitter.off('authed', this.refreshUser)
  },
  errorCaptured(err) {
    if (err as CombinedError) {
      if (this.errIsNotFound(err as CombinedError)) {
        this.error = err as CombinedError
        // do not propogate the error further up
        return false
      }
    }

    if (err as Error) {
      if (this.errIsNotFound(err as Error)) {
        this.error = err as Error
        // do not propogate the error further up
        return false
      }
    }

    if (this.errIsForbidden(err as CombinedError)) {
      this.toAuth()
      // do not propogate the error further up
      return false
    }

    if (this.errIsUnauthenticated(err as CombinedError)) {
      this.toAuth()
      // do not propogate the error further up
      return false
    }

    // 500-style errors, show as notification
    this.emitter.emit('notification', {
      title: 'Ooops! Something went wrong...',
      message: 'Please try again later.',
      style: 'error',
    })

    console.error(err)

    // Track with Sentry
    Sentry.captureException(err)

    // do not propogate the error further up
    return false
  },
  created() {
    this.emitter.on('logout', this.logout)
    this.emitter.on('reload-user', this.refreshAndIdentifyUser)

    // Add id before passing along
    this.emitter.on('notification', (n: ToastNotificationMessage) => {
      n.id = uuidv4()
      // default style
      if (!n.style) {
        n.style = 'success'
      }
      this.toastNotification(n)
    })

    this.emitter.on('notification-close', (n: ToastNotificationMessage) => {
      // remove notification with id
      this.toastNotifications = this.toastNotifications.filter((notif) => notif.id !== n.id)
    })
  },
  mounted() {
    this.emitter.on('authed', this.refreshUser)
    this.emitter.on('codebase', (id: string) => {
      if (this.postHogEnabled) {
        posthog.group('codebase', id)
      }
    })

    this.configureSentry()
    this.configurePostHog()

    this.$router.afterEach(() => {
      if (this.postHogEnabled) {
        posthog.capture('$pageview')
      }
    })
  },
  methods: {
    errIsUnauthenticated(err: null | CombinedError): boolean {
      if (!err) return false
      if (!err.graphQLErrors) return false
      return (
        err.graphQLErrors.filter(({ message }) => message === 'UnauthenticatedError').length > 0
      )
    },
    errIsNotFound(err: null | CombinedError | Error): boolean {
      if (!err) return false
      if (err?.message === 'SturdyCodebaseNotFoundError') return true
      if (!err.graphQLErrors) return false
      return err.graphQLErrors.filter(({ message }) => message === 'NotFoundError').length > 0
    },
    errIsForbidden(err: null | CombinedError): boolean {
      if (!err) return false
      if (!err.graphQLErrors) return false
      return err.graphQLErrors.filter(({ message }) => message === 'ForbiddenError').length > 0
    },

    toAuth() {
      // use replace instead of push to not break the browser history
      this.$router.replace({
        name: 'login',
        query: {
          navigateTo: escape(this.$route.path),
        },
      })
    },

    logout() {
      // Make a request to the backend so that we can delete the httpOnly cookie
      fetch(http.url('v3/auth/destroy'), {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include',
      })
        .then(http.checkStatus)
        // Navigate away from here
        .then(() => this.$router.push('/login'))
        .then(() => this.refreshUser())
        .then(() => {
          if (this.postHogEnabled) {
            posthog.reset()
          }
        })
        .then(() => {
          location.reload()
        })
        .catch(() => alert('something went wrong'))
    },

    async refreshAndIdentifyUser() {
      await this.refreshUser()

      // identify the user to posthog
      // this connects the web and backend identities
      if (this.data?.user?.id) {
        posthog.identify(this.data.user.id)
      }
    },

    configurePostHog() {
      if (
        window.location.href.indexOf('127.0.0.1') === -1 &&
        window.location.href.indexOf('localhost') === -1
      ) {
        posthog.init('ZuDRoGX9PgxGAZqY4RF9CCJJLpx14h3szUPzm7XBWSg', {
          api_host: 'https://app.posthog.com',
        })
        this.postHogEnabled = true
      } else {
        console.info('Ignoring PostHog on localhost')
      }
    },

    configureSentry() {
      Sentry.init({
        dsn: 'https://868feaf6fee74c368f2375232e045e5a@o952367.ingest.sentry.io/5901793',
        integrations: [new Integrations.BrowserTracing()],
        tracesSampleRate: 0.1, // TODO: Lower this if it turns out that this is too much
      })
    },

    toastNotification(notif: ToastNotificationMessage) {
      this.toastNotifications.push(notif)
    },
  },
})
</script>

<style type="text/css">
.spacer-padding {
  padding-top: env(titlebar-area-height, 2rem);
}
</style>
