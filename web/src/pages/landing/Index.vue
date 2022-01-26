<template>
  <div class="antialiased text-slate-300 bg-slate-900">
    <div class="relative">
      <div class="px-4 sm:px-6 md:px-8">
        <div class="absolute inset-0 bottom-10 bg-bottom bg-no-repeat bg-[#0B1120] gradient">
          <div class="absolute inset-0 bg-bottom border-b border-slate-100/5 bestagons" />
        </div>
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
              class="text-slate-900 bg-amber-500 hover:bg-amber-400 focus:outline-none focus:ring-2 focus:ring-slate-400 focus:ring-offset-2 focus:ring-offset-slate-50 font-semibold h-12 px-6 rounded-lg w-full flex items-center justify-center sm:w-auto highlight-white/20"
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
      </div>

      <div
        class="max-w-7xl mx-auto px-4 sm:px-6 md:px-8 mt-20 sm:mt-24 lg:mt-32 lg:grid lg:gap-8 lg:grid-cols-12 lg:items-center"
      >
        <div class="relative row-start-1 col-start-6 xl:col-start-7 col-span-7 xl:col-span-6">
          <UspVideo />
        </div>
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
            Sturdy maps developer
            <span class="text-amber-500">intent</span>
            to
            <span class="text-amber-500">concrete actions</span>.
            It allows you to focus on your code and the problems you are solving, rather than manipulating Git data structures.
          </p>
          <p>
            Think "Ship this code to production" instead of "Create a branch, stage files, make a
            commit, push to remote, merge" etc.
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
            When it comes to working together on code, we believe that
            <span class="text-amber-500">continuous feedback</span>
            is better than formal code "reviews".
          </p>
          <p>
            Sturdy allows you to quickly try your teammates' code on your machine. Give code
            suggestions by just editing the code in your IDE.
          </p>
        </template>
        <template #bottomLeft>
          <UspDetailsPair
            left-title="Trying code"
            left-body="Get your your teammates' code with a single click."
            right-title="Suggesting code"
            right-body="Give code suggestions by simpy typing in your IDE."
          >
            <template #leftIcon>
              <PlayIcon class="h-10 w-10 mb-4 text-purple-500" />
            </template>
            <template #leftLink>
              <span>See it in action</span>
            </template>
            <template #rightIcon>
              <LightBulbIcon class="h-10 w-10 mb-4 text-purple-500" />
            </template>
            <template #rightLink>
              <span>See it in action</span>
            </template>
          </UspDetailsPair>
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
            Use Sturdy
            <span class="text-amber-500">together with GitHub</span>, with your team or by
            yourself, and benefit from a leveraged workflow.
          </p>
        </template>
        <template #bottomLeft>
          <UspDetailsPair
            left-title="Connect GitHub"
            left-body="Setup a bridge between Sturdy and GitHub for one or multiple repositories."
            right-title="Create pull requests"
            right-body="Submit and merge pull requests without leaving Sturdy."
          >
            <template #leftIcon>
              <PlayIcon class="h-10 w-10 mb-4 text-purple-500" />
            </template>
            <template #leftLink>
              <span>See it in action</span>
            </template>
            <template #rightIcon>
              <LightBulbIcon class="h-10 w-10 mb-4 text-purple-500" />
            </template>
            <template #rightLink>
              <span>See it in action</span>
            </template>
          </UspDetailsPair>
        </template>
        <template #right>
          <UspVideo />
        </template>
      </Usp>

      <Usp
        title="Streamlined"
        subtitle="Built for the workflow that you're already using."
        link="v2DocsHotToShipSoftwareToProduction"
      >
        <template #left>
          <p>
            We believe that, for most software, shipping small and often is more effective than
            developing in long-lived feature branches.
          </p>
          <p>
            We are building Sturdy specifically around
            <span class="text-amber-500">trunk-based development</span> and
            <span class="text-amber-500">continuous delivery</span>. It optimizes for integration
            frequency and makes shipping of small incremental changes the intuitive default.
          </p>
        </template>
        <template #right>
          <UspVideo />
        </template>

      <Usp
        title="Open"
        subtitle="By engineers, for engineers."
        link="v2DocsRoot"
      >
        <template #left>
          <p>

            We believe that, for most software, shipping small and often is more effective than
            developing in long-lived feature branches.
          </p>
          <p>
            We are building Sturdy specifically around
            <span class="text-amber-500">trunk-based development</span>, and
            <span class="text-amber-500">continuous delivery</span>. It optimizes for integration
            frequency and makes shipping of small incremental changes the intuitive default.
          </p>
        </template>
        <template #right>
          <UspVideo />
        </template>

      </Usp>

      </Usp>
    </div>
  </div>
</template>

<script lang="ts">
import UspVideo from './UspVideo.vue'
import Usp from './Usp.vue'
import UspDetailsPair from './UspDetailsPair.vue'
import ThingsYouAreNotDoing from './ThingsYouAreNotDoing.vue'
import { defineComponent } from 'vue'
import { useHead } from '@vueuse/head'
import { DownloadIcon, ChevronRightIcon, PlayIcon, LightBulbIcon } from '@heroicons/vue/outline'

export default defineComponent({
  components: {
    UspVideo,
    Usp,
    ThingsYouAreNotDoing,
    UspDetailsPair,
    DownloadIcon,
    ChevronRightIcon,
    PlayIcon,
    LightBulbIcon,
  },
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

<style scoped>
.gradient {
  background-image: url('./gradient.svg');
  background-size: 86.125rem;
}

.bestagons {
  background-image: url('./bestagons.svg');
  background-size: 86.125rem;
}
</style>
