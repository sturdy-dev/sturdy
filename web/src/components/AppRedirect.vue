<template>
  <div v-if="isLoading" />

  <div v-else-if="showOverlay" class="p-3 pt-32 flex flex-col gap-5 items-center justify-center">
    <img src="../assets/Web/Duck/DuckCap256.png" class="h-16 w-16" alt="Sturdy Duck Logo" />

    <h1 class="text-3xl font-bold">Opening App...</h1>

    <p class="max-w-md text-center">
      We're trying to open this link up in your Sturdy app!<br />Hang tight, or
      <button class="text-yellow-600 underline" @click="showOverlay = false">
        continue in browser
      </button>
      .
    </p>

    <Spinner class="w-7 h-7 text-yellow-600" />
  </div>

  <slot v-else></slot>
</template>

<script lang="ts">
import { defineComponent, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import Spinner from './shared/Spinner.vue'

async function onBlurOnce() {
  await new Promise((resolve) => window.addEventListener('blur', resolve, { once: true }))
}

async function afterTimeout(ms: number) {
  await new Promise((resolve) => setTimeout(resolve, ms))
}

async function checkIfBlurs(f: () => void): Promise<boolean> {
  const timesOut = afterTimeout(300).then(() => false)
  const blurs = onBlurOnce().then(() => true)

  f()

  return Promise.race([timesOut, blurs])
}

export default defineComponent({
  components: { Spinner },
  setup() {
    const showOverlay = ref(false)
    const isLoading = ref(false)

    const currentRoute = useRoute()
    const router = useRouter()
    const ipc = window.ipc

    router.isReady().then(() => {
      const isAppRoute = !currentRoute.meta.nonApp && !currentRoute.meta.neverElectron

      const isMobile = !import.meta.env.SSR && /iPhone|iPad|iPod|Android/i.test(navigator.userAgent)
      const disabled = import.meta.env.VITE_DISABLE_WEB_TO_APP_REDIRECT

      if (!import.meta.env.SSR && isAppRoute && !isMobile && ipc == null && !disabled) {
        isLoading.value = true
        checkIfBlurs(() => {
          if (import.meta.env.DEV) {
            location.assign(`sturdy-dev://${location.pathname}${location.search}`)
          } else {
            location.assign(`sturdy://${location.pathname}${location.search}`)
          }
        })
          .then((blurred) => {
            showOverlay.value = blurred
          })
          .catch((e) => {
            console.error(e)
          })
          .finally(() => {
            isLoading.value = false
          })
      }
    })

    return { showOverlay, isLoading }
  },
})
</script>
