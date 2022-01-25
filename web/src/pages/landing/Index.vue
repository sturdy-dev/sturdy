<template>
  <div class="antialiased text-slate-300 bg-slate-900">
    <!--    <div class="absolute inset-0 bottom-10 bg-bottom bg-no-repeat bg-[#0B1120]">-->
    <!--    </div>-->
    <div class="relative max-w-5xl mx-auto pt-20 sm:pt-24 lg:pt-32">
      <h1
        class="font-extrabold text-4xl sm:text-5xl lg:text-6xl tracking-tight text-center text-white"
      >
        Real-time code collaboration.
      </h1>
      <p class="mt-6 text-lg text-slate-400 text-center max-w-3xl mx-auto">
        Sturdy is an
        <code class="font-mono font-semibold text-amber-500"
          >open-source version control platform</code
        >
        that allows you to interact with your code at a higher abstraction level.
      </p>
      <div class="mt-6 sm:mt-10 flex justify-center space-x-6 text-sm">
        <router-link
          :to="{ name: 'download' }"
          class="text-white bg-amber-500 hover:bg-amber-400 focus:outline-none focus:ring-2 focus:ring-slate-400 focus:ring-offset-2 focus:ring-offset-slate-50 font-semibold h-12 px-6 rounded-lg w-full flex items-center justify-center sm:w-auto highlight-white/20"
        >
          <DownloadIcon class="h-6 w-6 mr-1" />
          {{ mainDownloadText }}
        </router-link>
        <router-link
          :to="{ name: 'v2DocsRoot' }"
          class="text-slate-300 bg-slate-700 hover:bg-slate-600 focus:outline-none focus:ring-2 focus:ring-slate-400 focus:ring-offset-2 focus:ring-offset-slate-50 font-semibold h-12 px-6 rounded-lg w-full flex items-center justify-center sm:w-auto highlight-white/20"
        >
          Read the Docs
          <ChevronRightIcon class="h-6 w-6 ml-1" />
        </router-link>
      </div>
    </div>

    <div
      class="pt-20 mb-20 space-y-20 overflow-hidden sm:pt-32 sm:mb-32 sm:space-y-32 md:pt-40 md:mb-40 md:space-y-40"
    >
      <Usp
        title="Expressive"
        subtitle="Don't hand-hold your version control."
        link="v2DocsHowToEditCode"
      >
        <template #left>
          <p>
            Sturdy maps developer intent to concrete actions, rather than manipulating Git data
            structures. It allows you to focus on your code and the problems you are solving.
          </p>
          <p>
            Think &mdash; "Ship this code to production" instead of "Create a branch, stage files,
            make a commit, push to remote, merge, etc.".
          </p>
        </template>
        <template #right>
          <ThingsYouAreNotDoing />
        </template>
      </Usp>

      <Usp
        title="Collaborative"
        subtitle="Collaborative everything."
        link="v2DocsHowToCollaborateWithOthers"
      >
        <template #left>
          <p>
            When it comes to working together on code, we believe that continuous feedback is better
            than formal code "reviews".
          </p>
          <p>
            Sturdy allows you to quickly try your teammates' code on your machine. Give code
            suggestions by just typing in your text editor.
          </p>
        </template>
        <template #right>
          <UspVideo />
        </template>
      </Usp>

      <Usp
        title="Compatible"
        subtitle="Worried about compatibility? Don't be."
        link="v2DocsHowSturdyAugmentsGit"
      >
        <template #left>
          <p>
            Sturdy utilizes low-level Git data structures, which means your code is stored in a
            compatible format.
          </p>
          <p>
            Use Sturdy together with GitHub, either with your team or by yourself, and benefit from
            a leveraged workflow.
          </p>
        </template>
        <template #right>
          <UspVideo />
        </template>
      </Usp>

      <Usp
        title="Streamlined"
        subtitle="The best feedback comes from production."
        link="v2DocsHotToShipSoftwareToProduction"
      >
        <template #left>
          <p>
            We believe that, for most software, shipping small and often is more effective than
            developing in long-lived feature branches.
          </p>
          <p>
            We are building Sturdy specifically around trunk-based development and continuous
            delivery. It optimizes for integration frequency and makes shipping of small incremental
            changes the intuitive default.
          </p>
        </template>
        <template #right>
          <UspVideo />
        </template>
      </Usp>
    </div>
  </div>
</template>

<script lang="ts">
import UspVideo from './UspVideo.vue'
import Usp from './Usp.vue'
import ThingsYouAreNotDoing from './ThingsYouAreNotDoing.vue'
import { defineComponent } from 'vue'
import { useHead } from '@vueuse/head'
import { DownloadIcon, ChevronRightIcon } from '@heroicons/vue/outline'

export default defineComponent({
  components: { UspVideo, Usp, ThingsYouAreNotDoing, DownloadIcon, ChevronRightIcon },
  setup() {
    // TODO: Remove when we're launching!
    useHead({
      meta: [
        {
          name: 'robots',
          content: 'noindex',
        },
      ],
    })

    let mainDownloadText = 'Download'

    if (!import.meta.env.SSR) {
      if (navigator.userAgent.includes('Win')) {
        mainDownloadText = 'Download for Windows'
      } else if (navigator.userAgent.includes('Mac')) {
        mainDownloadText = 'Download for Mac'
      } else if (navigator.userAgent.includes('Linux')) {
        mainDownloadText = 'Download for Linux'
      }
    }

    return {
      mainDownloadText,
    }
  },
})
</script>

<style scoped></style>
