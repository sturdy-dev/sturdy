<template>
  <PaddedApp>
    <div v-if="show_redirected_to_app">
      <div class="p-3 pt-32 flex flex-col gap-5 items-center justify-center">
        <img src="../assets/Web/Duck/DuckCap256.png" class="h-16 w-16" alt="Sturdy Duck Logo" />

        <h1 class="text-3xl font-bold">Opening App...</h1>

        <p class="max-w-md text-center">
          We're trying to open this link up in your Sturdy app!<br />Hang tight, or
          <router-link class="text-yellow-600 underline" :to="{ name: 'codebaseOverview' }">
            continue in browser
          </router-link>
          .
        </p>

        <Spinner class="w-7 h-7 text-yellow-600" />
      </div>
    </div>
    <Banner v-else-if="show_oauth_went_wrong" status="error">
      Something went wrong connecting to GitHub, please try again!
    </Banner>
    <Banner v-else status="info">
      <div class="inline-flex">
        <span>Setting up...</span>
        <Spinner class="ml-3" />
      </div>
    </Banner>
  </PaddedApp>
</template>

<script lang="ts">
import { defineComponent, onMounted, ref } from 'vue'
import http from '../http'
import { useRoute, useRouter } from 'vue-router'
import Banner from '../components/shared/Banner.vue'
import Spinner from '../components/shared/Spinner.vue'
import PaddedApp from '../layouts/PaddedApp.vue'

export default defineComponent({
  components: { PaddedApp, Banner, Spinner },
  setup() {
    let route = useRoute()
    let router = useRouter()

    let show_oauth_went_wrong = ref(false)
    let show_redirected_to_app = ref(false)

    onMounted(async () => {
      // If this request is rendered on the web, open in the app.
      if (typeof ipc === 'undefined' && route.query.state === 'install-app') {
        if (import.meta.env.DEV) {
          location.assign(`sturdy-dev://${location.pathname}${location.search}`)
        } else {
          location.assign(`sturdy://${location.pathname}${location.search}`)
        }
        show_redirected_to_app.value = true
        return
      }

      if (
        route.query.state === 'install' ||
        route.query.state === 'install-app' ||
        route.query.setup_action === 'update'
      ) {
        await router.push({ name: 'codebaseOverview' })
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
          if (route.query.state === 'user-settings') {
            router.push({ name: 'user' })
          } else if (route.query.state === 'codebase-overview') {
            router.push({ name: 'codebaseOverview' })
          } else if (route.query.state.startsWith('web-')) {
            router.push(route.query.state.substring(4))
          }
        })
        .catch(() => {
          show_oauth_went_wrong.value = true
        })
    })

    return {
      show_oauth_went_wrong,
      show_redirected_to_app,
    }
  },
})
</script>
