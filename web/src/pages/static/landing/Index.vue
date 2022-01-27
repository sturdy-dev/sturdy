<template>
  <div class="antialiased text-slate-300 bg-slate-900">
    <header class="relative">
      <div
        class="blurr border-b border-slate-700/10 sticky top-0 z-50 py-2 lg:py-3 px-6 lg:px-20 flex items-center justify-between leading-6"
      >
        <img
          src="./logotype.svg"
          alt="Sturdy logotype"
          class="rounded-xl object-contain"
          width="110"
        />

        <div class="flex items-center font-semibold text-sm">
          <nav>
            <ul
              class="flex items-center space-x-8 focus:outline-none focus:ring-2 focus:ring-slate-400 focus:ring-offset-2 focus:ring-offset-slate-50 focus:outline-none focus:ring-2"
            >
              <router-link :to="{ name: 'v2DocsRoot' }" class="hover:text-amber-500">
                Docs
              </router-link>
              <router-link :to="{ name: 'pricing' }" class="hover:text-amber-500">
                Pricing
              </router-link>
              <router-link :to="{ name: 'blog' }" class="hover:text-amber-500"> Blog </router-link>

              <router-link
                :to="{ name: 'download' }"
                class="text-slate-900 bg-amber-500 hover:text-amber-500 hover:bg-slate-900/10 border border-transparent hover:border-amber-500 font-semibold h-9 px-3 rounded-lg w-full flex items-center justify-center sm:w-auto highlight-white/20 focus:outline-none focus:ring-2 focus:ring-slate-400 focus:ring-offset-2 focus:ring-offset-slate-50 focus:outline-none focus:ring-2"
              >
                Get started
              </router-link>
            </ul>
          </nav>
        </div>
      </div>

      <div class="px-4 sm:px-6 md:px-8">
        <div class="absolute inset-0 bottom-10 bg-bottom bg-no-repeat gradient"></div>

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
              class="text-slate-900 bg-amber-500 hover:text-amber-500 hover:bg-slate-900/10 border border-transparent hover:border-amber-500 focus:outline-none focus:ring-2 focus:ring-slate-400 focus:ring-offset-2 focus:ring-offset-slate-50 font-semibold h-12 px-6 rounded-lg w-full flex items-center justify-center sm:w-auto highlight-white/20"
            >
              <DownloadIcon class="h-6 w-6 mr-1" />
              {{ mainDownloadText }}
            </router-link>
            <router-link
              :to="{ name: 'v2DocsRoot' }"
              class="text-slate-300 bg-slate-600 hover:text-slate-300 hover:bg-slate-900/10 border border-transparent hover:border-slate-300 focus:outline-none focus:ring-2 focus:ring-slate-400 focus:ring-offset-2 focus:ring-offset-slate-50 font-semibold h-12 px-6 rounded-lg w-full flex items-center justify-center sm:w-auto highlight-white/20"
            >
              Read the Docs
              <ChevronRightIcon class="h-6 w-6 ml-1" />
            </router-link>
          </div>
        </div>
      </div>

      <div class="max-w-3xl px-4 sm:px-6 md:px-8 mx-auto items-center mt-20 sm:mt-24 lg:mt32">
        <div class="relative row-start-1 col-start-6 xl:col-start-7 col-span-7 xl:col-span-6">
          <UspVideo />
        </div>
      </div>
    </header>

    <section id="features">
      <div class="pt-10 xl:pt-0 mt-10 xl:mt-2 space-y-32 sm:space-y-40 md:space-y-60">
        <div class="pt-16 max-w-7xl px-4 sm:px-6 md:px-8 mx-auto lg:gap-8 space-y-6">
          <h1 class="mt-4 text-3xl sm:text-4xl font-extrabold tracking-tight text-slate-50">
            A new home for your code
          </h1>
          <div class="pt-8 px-4 flex flex-col sm:flex-row sm:space-x-14">
            <MiniUsp class="sm:w-1/4" title="Expressive">
              <p>
                Ship the code to production &mdash; <br />
                instead of the usual dance of creating a branch, staging and committing files,
                pushing, and merging.
              </p>
            </MiniUsp>

            <MiniUsp class="sm:w-1/4" title="Collaborative">
              <div>
                Try each others code and give actionable feedback &mdash; <br />
                instead of guessing what a pull request does and LGTM-ing anyway.
              </div>
            </MiniUsp>

            <MiniUsp class="sm:w-1/4" title="Compatible">
              <div>
                Use Sturdy together with GitHub &mdash; or on it's own. By yourself &mdash; or with
                others on your team. Run your own instance &mdash; or use Sturdy in the cloud.
              </div>
            </MiniUsp>

            <MiniUsp class="sm:w-1/4" title="Streamlined">
              <div>
                Ship small, incremental changes as the default &mdash;<br />
                instead of reminding yourself to "push often" and "keep the PR small".
              </div>
            </MiniUsp>
          </div>
        </div>
      </div>
    </section>

    <div class="mb-20 space-y-20 overflow-hidden sm:mb-32 sm:space-y-32 md:mb-40 md:space-y-40">
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
            <span class="text-amber-500">concrete actions</span>. It allows you to focus on your
            code and the problems you are solving, rather than manipulating Git data structures.
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
            <span class="text-amber-500">together with GitHub</span>, with your team or by yourself,
            and benefit from a leveraged workflow.
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
        subtitle="For the workflow that you're already using."
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
      </Usp>
      <Usp title="Open-source" subtitle="For engineers, by engineers." link="v2DocsRoot">
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
    </div>
  </div>
</template>

<script lang="ts">
import UspVideo from './UspVideo.vue'
import Usp from './Usp.vue'
import MiniUsp from './MiniUsp.vue'
import UspDetailsPair from './UspDetailsPair.vue'
import ThingsYouAreNotDoing from './ThingsYouAreNotDoing.vue'
import { defineComponent } from 'vue'
import { useHead } from '@vueuse/head'
import { DownloadIcon, ChevronRightIcon, PlayIcon, LightBulbIcon } from '@heroicons/vue/outline'

export default defineComponent({
  components: {
    UspVideo,
    Usp,
    MiniUsp,
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
  background-image: #0b1120;
  background-image: radial-gradient(at 23% 84%, hsla(223, 49%, 9%, 1) 0, transparent 56%),
    radial-gradient(at 81% 30%, hsla(223, 49%, 9%, 1) 0, transparent 46%),
    radial-gradient(at 87% 34%, hsla(223, 49%, 9%, 1) 0, transparent 57%),
    radial-gradient(at 16% 71%, hsla(223, 49%, 9%, 1) 0, transparent 40%),
    radial-gradient(at 100% 100%, hsla(223, 49%, 9%, 1) 0, transparent 52%),
    radial-gradient(at 51% 63%, hsla(271, 92%, 66%, 1) 0, transparent 49%),
    radial-gradient(at 0% 5%, hsla(223, 49%, 9%, 1) 0, transparent 49%),
    radial-gradient(at 74% 76%, hsla(44, 97%, 57%, 1) 0, transparent 50%),
    radial-gradient(at 45% 37%, hsla(38, 93%, 51%, 1) 0, transparent 50%);
}

.blurr {
  backdrop-filter: saturate(70%) brightness(80%) blur(5px);
}
</style>
