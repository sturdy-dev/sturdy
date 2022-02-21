<template>
  <ElectronNavigation>
    <PaddedApp>
      <div v-if="show_redirected_to_app">
        <div class="p-3 pt-32 flex flex-col gap-5 items-center justify-center">
          <img src="../assets/Web/Duck/DuckCap256.png" class="h-16 w-16" alt="Sturdy Duck Logo" />

          <h1 class="text-3xl font-bold">Opening App...</h1>

          <p class="max-w-md text-center">
            We're trying to open this link up in your Sturdy app!<br />Please hang tight!
            <router-link class="text-yellow-600 underline" :to="{ name: 'home' }">
              Continue to Sturdy
            </router-link>
            .
          </p>

          <Spinner class="w-7 h-7 text-yellow-600" />
        </div>
      </div>
      <div v-else-if="errorMessage" class="flex flex-col space-y-4">
        <Banner status="error">
          {{ errorMessage }}
        </Banner>
        <div>
          <RouterLinkButton color="green" :to="{ name: 'home' }"
            >Take me back to safety</RouterLinkButton
          >
        </div>
      </div>
      <Banner v-else status="info">
        <div class="inline-flex">
          <span>Setting up...</span>
          <Spinner class="ml-3" />
        </div>
      </Banner>
    </PaddedApp>
  </ElectronNavigation>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import http, { HttpError } from '../http'
import { useRoute, useRouter } from 'vue-router'
import { Banner } from '../atoms'
import Spinner from '../components/shared/Spinner.vue'
import PaddedApp from '../layouts/PaddedApp.vue'
import ElectronNavigation from '../layouts/ElectronNavigation.vue'
import RouterLinkButton from '../components/shared/RouterLinkButton.vue'

const openInAppUrl = (): string | undefined => {
  let searchParams = new URLSearchParams(location.search)
  let state = searchParams.get('state')

  if (!state) {
    return undefined
  }

  let path
  let proto
  if (state.startsWith('app-dev-')) {
    path = state.substring(8)
    proto = 'sturdy-dev'
  } else if (state.startsWith('app-')) {
    path = state.substring(4)
    proto = 'sturdy'
  }

  searchParams.set('state', 'web-' + path)
  let search = searchParams.toString()
  return `${proto}://${location.pathname}?${search}`
}

export default defineComponent({
  components: { ElectronNavigation, PaddedApp, Banner, Spinner, RouterLinkButton },
  data() {
    return {
      show_redirected_to_app: false,

      errorMessage: undefined as string | undefined,
    }
  },
  async mounted() {
    const route = useRoute()
    const router = useRouter()
    // Open this page in the app
    if (typeof ipc === 'undefined' && route.query.state) {
      let state = route.query.state as string
      if (state.startsWith('app-') || state.startsWith('app-dev-')) {
        let url = openInAppUrl()
        if (url) {
          location.assign(url)
          this.show_redirected_to_app = true
          return
        }
      }
    }

    // TODO(gustav): remove the following paths
    // If this request is rendered on the web, open in the app.
    if (typeof ipc === 'undefined' && route.query.state === 'install-app') {
      if (import.meta.env.DEV) {
        location.assign(`sturdy-dev://${location.pathname}${location.search}`)
      } else {
        location.assign(`sturdy://${location.pathname}${location.search}`)
      }
      this.show_redirected_to_app = true
      return
    }

    await fetch(http.url('v3/github/oauth'), {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        code: route.query.code,
      }),
      credentials: 'include',
    })
      .then(http.checkStatus)
      .then(() => {
        const state = route.query.state as string
        if (state === 'user-settings') {
          router.push({ name: 'user' })
        } else if (state === 'codebase-overview') {
          router.push({ name: 'codebaseOverview' })
        } else if (state.startsWith('web-')) {
          router.push(state.substring(4))
        } else if (state.startsWith('app-dev-')) {
          router.push(state.substring(8))
        } else if (state.startsWith('app-')) {
          router.push(state.substring(4))
        } else {
          // fallback
          router.push({ name: 'codebaseOverview' })
        }
      })
      .catch(async (e) => {
        const err = e as HttpError
        if (err.statusCode === 400) {
          this.errorMessage = `Failed to connect: ${(await err.response?.json()).error}`
        } else {
          this.errorMessage = 'Something went wrong connecting to GitHub, please try again!'
        }
      })
  },
})
</script>
