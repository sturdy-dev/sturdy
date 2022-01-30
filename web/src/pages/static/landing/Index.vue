<template>
  <div class="antialiased text-slate-300 bg-slate-900">
    <Nav />

    <div class="px-6 relative -mt-16">
      <div class="absolute inset-0 bottom-10 bg-bottom bg-no-repeat gradient"></div>

      <div class="relative max-w-5xl mx-auto pt-24 sm:pt-28 lg:pt-32">
        <h1
          class="font-extrabold text-4xl sm:text-5xl lg:text-6xl tracking-tight text-center text-slate-50"
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
        <div
          class="flex flex-col space-y-6 sm:flex-row sm:space-y-0 sm:space-x-6 sm:justify-center mt-10 mx-6 text-sm"
        >
          <router-link
            :to="{ name: 'download' }"
            class="text-slate-900 bg-amber-500 hover:text-amber-500 hover:bg-slate-900/10 border border-transparent hover:border-amber-500 focus:outline-none focus:ring-2 focus:ring-slate-400 focus:ring-offset-2 focus:ring-offset-slate-50 font-semibold h-12 px-6 rounded w-full flex items-center justify-center sm:w-auto highlight-white/20"
          >
            <DownloadIcon class="h-6 w-6 mr-1" />
            {{ mainDownloadText }}
          </router-link>
          <router-link
            :to="{ name: 'v2DocsRoot' }"
            class="text-amber-500/70 hover:text-amber-500 bg-slate-900/10 border border-transparent border-amber-500/70 hover:border-amber-500 focus:outline-none focus:ring-2 focus:ring-slate-400 focus:ring-offset-2 focus:ring-offset-slate-50 font-semibold h-12 px-6 rounded w-full flex items-center justify-center sm:w-auto highlight-white/20"
          >
            Read the Docs
            <ChevronRightIcon class="h-6 w-6 ml-1" />
          </router-link>
        </div>

        <div class="flex flex-row space-x-6 justify-center mt-10 mx-6 text-sm">
          <div class="text-lg font-bold tracking-tight text-slate-400/50">backed by</div>
          <div class="flex space-x-10 items-center">
            <a href="https://www.ycombinator.com/companies/sturdy">
              <img
                src="../../../assets/Landing/ycombinator-logo.png"
                alt="Y Combinator logo"
                class="object-contain h-6 grayscale opacity-50"
              />
            </a>
          </div>
        </div>
      </div>

      <div class="max-w-3xl px-4 sm:px-6 md:px-8 mx-auto items-center mt-20 sm:mt-24 lg:mt32">
        <div class="relative row-start-1 col-start-6 xl:col-start-7 col-span-7 xl:col-span-6">
          <UspVideo />
        </div>
      </div>
    </div>

    <section id="features">
      <div class="bg-slate-800/50 border-y border-slate-700 mt-20">
        <div class="my-12 max-w-6xl px-6 md:px-8 mx-auto lg:gap-8 space-y-6">
          <h1
            class="mb-12 sm:text-center mt-4 text-3xl sm:text-4xl font-extrabold tracking-tight text-slate-50"
          >
            A new home for your code
          </h1>
          <div class="flex flex-col md:flex-row md:space-x-20 space-y-10 md:space-y-0">
            <MiniUsp class="md:w-1/3" title="Expressive" accent="text-amber-500">
              <template #icon>
                <FireIcon />
              </template>
              <span>
                Ship your code to production &mdash; <br />
                instead of creating a branches, staging and committing files, pushing, and merging.
              </span>
            </MiniUsp>

            <MiniUsp class="md:w-1/3" title="Collaborative" accent="text-purple-400">
              <template #icon>
                <CogIcon />
              </template>
              <span>
                Try each others code &mdash; <br />
                instead of guessing what a pull request does and LGTM-ing anyway.
              </span>
            </MiniUsp>

            <MiniUsp class="md:w-1/3" title="Streamlined" accent="text-green-500">
              <template #icon>
                <RefreshIcon />
              </template>
              <span>
                Ship small, incremental changes &mdash;<br />
                instead of reminding yourself to "push often" and "keep the PR small".
              </span>
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
        accent="text-amber-500"
      >
        <template #icon>
          <IconDevelop class="" />
        </template>
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
        accent="text-purple-400"
      >
        <template #icon>
          <IconGears></IconGears>
        </template>
        <template #left>
          <p>
            When it comes to working together on code, we believe that
            <span class="text-purple-400">continuous feedback</span>
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
        title="Streamlined"
        subtitle="For the workflow that you're already using."
        link="v2DocsHotToShipSoftwareToProduction"
        accent="text-green-500"
      >
        <template #icon>
          <IconDesktop></IconDesktop>
        </template>
        <template #left>
          <p>
            We believe that, for most software, shipping small and often is more effective than
            developing in long-lived feature branches.
          </p>
          <p>
            We are building Sturdy specifically around
            <span class="text-green-500">trunk-based development</span> and
            <span class="text-green-500">continuous delivery</span>. It optimizes for integration
            frequency and makes shipping of small incremental changes the intuitive default.
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
        <template #icon>
          <IconMobile></IconMobile>
        </template>
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

      <Usp title="Open-source" subtitle="For engineers, by engineers." link="v2DocsRoot">
        <template #icon>
          <IconLaptop></IconLaptop>
        </template>
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
          <div
            class="border border-slate-300/30 relative overflow-hidden shadow-xl flex bg-slate-700/50 rounded-xl max-h-[60vh] sm:max-h-[none] h-72 lg:h-80 xl:h-96"
          >
            <div class="relative w-full flex flex-col">
              <div class="flex-none border-b border-slate-300/30">
                <div class="flex items-center h-8 space-x-1.5 px-3">
                  <div class="w-2.5 h-2.5 bg-red-700 rounded-full"></div>
                  <div class="w-2.5 h-2.5 bg-yellow-600 rounded-full"></div>
                  <div class="w-2.5 h-2.5 bg-green-600 rounded-full"></div>
                </div>

                <div
                  class="-mt-6 mb-1 text-sm font-semibold text-slate-400 text-center select-none"
                >
                  Bash
                </div>
              </div>

              <div class="relative min-h-0 flex-auto flex flex-col">
                <div class="w-full flex-auto flex min-h-0">
                  <div class="w-full flex-auto flex min-h-0 overflow-auto">
                    <div class="w-full relative flex-auto">
                      <div class="flex flex-col min-h-full text-sm leading-6">
                        <code
                          class="inline-flex flex-auto relative block text-slate-50 overflow-auto p-4"
                        >
                          <span class="text-slate-400 mr-4 select-none">$</span>
                          <pre id="oneliner"
                            >{{ dockerOneliner }}
                          </pre>
                        </code>
                        <div class="flex flex-row-reverse">
                          <DuplicateIcon
                            class="cursor-pointer m-4 w-8 h-8 hover:text-slate-50"
                            @click="copyDockerOneliner"
                          />
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </template>
      </Usp>

      <section id="featured">
        <div class="pt-10 xl:pt-0 mt-10 xl:mt-2 space-y-32 sm:space-y-40 md:space-y-60">
          <div class="pt-16 max-w-6xl px-6 mx-auto lg:gap-8 space-y-6">
            <h1 class="mt-4 text-lg sm:text-xl font-extrabold tracking-tight text-slate-300">
              Featured at
            </h1>
            <div
              class="flex flex-col space-y-6 sm:flex-row sm:space-y-0 sm:space-x-10 sm:items-center"
            >
              <a href="https://websummit.com/schedule/ws21/timeslot/breakout-startups-1104-0927">
                <img
                  src="./web-summit.png"
                  alt="Web Summit logotype"
                  class="rounded-xl object-contain h-4"
                />
              </a>
              <a href="https://devsecops.fm/episodes/36-sturdy/">
                <img
                  src="./devsecops.png"
                  alt="DevSecOps logotype"
                  class="rounded-xl object-contain h-6"
                />
              </a>
              <a href="">
                <img
                  src="./devcon.png"
                  alt="Devcon logotype"
                  class="rounded-xl object-contain h-6"
                />
              </a>
            </div>
          </div>
        </div>
      </section>
    </div>

    <div id="footer" class="text-sm leading-6">
      <Footer />
    </div>
  </div>
</template>

<script lang="ts">
import UspVideo from './UspVideo.vue'
import Usp from './Usp.vue'
import MiniUsp from './MiniUsp.vue'
import UspDetailsPair from './UspDetailsPair.vue'
import ThingsYouAreNotDoing from './ThingsYouAreNotDoing.vue'
import IconDevelop from './IconDevelop.vue'
import IconGears from './IconGears.vue'
import IconDesktop from './IconDesktop.vue'
import IconLaptop from './IconLaptop.vue'
import IconMobile from './IconMobile.vue'
import Footer from './Footer.vue'
import Nav from './Nav.vue'
import { defineComponent } from 'vue'
import { useHead } from '@vueuse/head'
import {
  DownloadIcon,
  ChevronRightIcon,
  PlayIcon,
  LightBulbIcon,
  DuplicateIcon,
  ExternalLinkIcon,
  FireIcon,
  CogIcon,
  RefreshIcon,
} from '@heroicons/vue/outline'

export default defineComponent({
  components: {
    Footer,
    Nav,
    UspVideo,
    Usp,
    MiniUsp,
    ThingsYouAreNotDoing,
    UspDetailsPair,
    IconDevelop,
    IconGears,
    IconDesktop,
    IconLaptop,
    IconMobile,
    DownloadIcon,
    ChevronRightIcon,
    PlayIcon,
    LightBulbIcon,
    DuplicateIcon,
    FireIcon,
    CogIcon,
    RefreshIcon,
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
      year: new Date().getFullYear(),
      dockerOneliner: `docker run --interactive --tty \\
             --publish 30080:80 --publish 30022:22 \\
             --volume "$(pwd)/sturdydata:/var/data" \\
             getsturdy/server`,
    }
  },

  methods: {
    copyDockerOneliner() {
      navigator.clipboard.writeText(this.dockerOneliner)
    },
  },
})
</script>

<style>
body {
  background-color: #0f172a;
}
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
</style>
