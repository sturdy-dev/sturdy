<template>
  <AppTitleBar :show-sidebar="false">
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
            >Take me back to safety
          </RouterLinkButton>
        </div>
      </div>
      <Banner v-else status="info">
        <div class="inline-flex">
          <span>Setting up...</span>
          <Spinner class="ml-3" />
        </div>
      </Banner>
    </PaddedApp>
  </AppTitleBar>
</template>

<script lang="ts">
import { defineComponent, onMounted, ref } from 'vue'
import type { HttpError } from '../http'
import http from '../http'
import { useRoute, useRouter } from 'vue-router'
import { Banner } from '../atoms'
import Spinner from '../atoms/Spinner.vue'
import PaddedApp from '../layouts/PaddedApp.vue'
import RouterLinkButton from '../atoms/RouterLinkButton.vue'
import AppTitleBar from '../organisms/AppTitleBar.vue'

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
  components: { AppTitleBar, PaddedApp, Banner, Spinner, RouterLinkButton },
  setup() {
    /*
    SetupGithub handles redirects from GitHub back to Sturdy.

    It has three paths that are in use:
    1) OAuth on Web:                Web -> GitHub -> Web (this component)
    2) OAUth in App:                App -> GitHub -> App (this component)
    3) Manage installation on web:
    3) App -> GitHub -> Web (this component) -> App (this component)
    */

    const route = useRoute()
    const router = useRouter()

    let errorMessage = ref('')
    let show_redirected_to_app = ref(false)

    // Open this page in the app (web to app redirect)
    // http://localhost:8080/setup-github?code=XXXXX&installation_id=XXXXX&setup_action=update&state=app-dev-%2Fuser

    onMounted(() => {
      if (typeof window?.ipc === 'undefined' && route.query.state) {
        let state = route.query.state as string
        if (state.startsWith('app-') || state.startsWith('app-dev-')) {
          let url = openInAppUrl()
          if (url) {
            location.assign(url)
            show_redirected_to_app.value = true
            return
          }
        }
      }

      // if no code, redirect away
      if (!route.query.code) {
        router.push({ name: 'home' })
        return
      }

      fetch(http.url('v3/github/oauth'), {
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
          } else {
            // fallback
            router.push({ name: 'home' })
          }
        })
        .catch((e) => {
          const err = e as HttpError
          if (err.statusCode === 400) {
            err.response?.json().then((msg) => {
              errorMessage.value = `Failed to connect: ${msg.error}`
            })
          } else {
            errorMessage.value = 'Something went wrong connecting to GitHub, please try again!'
          }
        })
    })

    return {
      show_redirected_to_app,
      errorMessage,
    }
  },
})
</script>
