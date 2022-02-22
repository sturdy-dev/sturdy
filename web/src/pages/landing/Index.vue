<template>
  <DocumentationStickyHeader class="antialiased text-slate-300 bg-slate-900">
    <div class="px-6 relative -mt-16">
      <div class="absolute inset-0 bottom-20 bg-bottom bg-no-repeat gradient"></div>

      <div class="relative max-w-5xl mx-auto pt-24 sm:pt-28 lg:pt-32">
        <h1
          class="font-extrabold text-4xl sm:text-5xl lg:text-6xl tracking-tight text-center text-slate-50"
        >
          Real-time code collaboration.
        </h1>
        <p class="mt-6 text-lg text-slate-400 text-center max-w-3xl mx-auto">
          Sturdy is an
          <code class="font-mono font-semibold text-amber-500"
            ><a class="underline" href="https://getsturdy.com/sturdy-zyTDsnY">open-source</a>
            version control platform</code
          >
          that allows you to interact with your code at a higher abstraction level.
        </p>
        <div
          class="flex flex-col space-y-6 sm:flex-row sm:space-y-0 sm:space-x-6 sm:justify-center mt-10 mx-6 text-sm"
        >
          <router-link
            :to="{ name: 'v2download' }"
            class="text-slate-900 bg-amber-500 hover:text-amber-500 hover:bg-slate-900/10 border border-transparent hover:border-amber-500 focus:outline-none focus:ring-2 focus:ring-slate-400 focus:ring-offset-2 focus:ring-offset-slate-50 font-semibold h-12 px-6 rounded w-full flex items-center justify-center sm:w-auto highlight-white/20"
          >
            <DownloadIcon class="h-6 w-6 mr-1" />
            {{ mainDownloadText }}
          </router-link>
          <router-link
            :to="{ name: 'v2DocsRoot' }"
            class="text-amber-500/80 hover:text-amber-500 bg-slate-900/10 border border-transparent border-amber-500/80 hover:border-amber-500 focus:outline-none focus:ring-2 focus:ring-slate-400 focus:ring-offset-2 focus:ring-offset-slate-50 font-semibold h-12 px-6 rounded w-full flex items-center justify-center sm:w-auto highlight-white/20"
          >
            Read the Docs
            <ChevronRightIcon class="h-6 w-6 ml-1" />
          </router-link>
        </div>

        <div class="flex flex-row space-x-6 justify-center mt-10 mx-6 text-sm">
          <div class="text-lg font-bold tracking-tight text-slate-400/70">backed by</div>
          <div class="flex space-x-10 items-center">
            <a href="https://www.ycombinator.com/companies/sturdy">
              <img
                src="../../assets/Landing/ycombinator-logo.png"
                alt="Y Combinator logo"
                class="object-contain h-6 grayscale opacity-70"
                height="24"
                width="24"
              />
            </a>
          </div>
        </div>
      </div>

      <div class="max-w-4xl mx-auto items-center mt-16">
        <div class="relative" style="filter: drop-shadow(0 0 1.5rem rgba(192, 132, 252, 0.14))">
          <Video :src="videoCollaborate" />
        </div>
      </div>
    </div>

    <section id="features">
      <div
        class="-mt-20 border-t border-slate-700/30 block top-0 inset-x-0 bg-gradient-to-b from-[#0c1322] top-0 xl:top-8"
      >
        <div class="mt-48 mb-12 max-w-6xl px-6 md:px-8 mx-auto lg:gap-8 space-y-6">
          <h1
            class="mb-12 sm:text-center mt-4 text-3xl sm:text-4xl font-extrabold tracking-tight text-slate-50"
          >
            Ship faster and focus on your code
          </h1>
          <div class="flex flex-col md:flex-row md:space-x-20 space-y-10 md:space-y-0">
            <MiniUsp class="md:w-1/3" title="Expressive" accent="text-amber-500">
              <template #icon>
                <FireIcon />
              </template>
              <span>
                <span class="text-amber-500">Ship your code to production &mdash;</span> <br />
                instead of creating branches, staging and committing files, pushing, and merging.
              </span>
            </MiniUsp>

            <MiniUsp class="md:w-1/3" title="Collaborative" accent="text-purple-400">
              <template #icon>
                <CogIcon />
              </template>
              <span>
                <span class="text-purple-400">Try each other's code &mdash;</span> <br />
                instead of guessing what a pull request does and LGTM-ing anyway.
              </span>
            </MiniUsp>

            <MiniUsp class="md:w-1/3" title="Streamlined" accent="text-green-500">
              <template #icon>
                <RefreshIcon />
              </template>
              <span>
                <span class="text-green-500">Ship small, incremental changes &mdash;</span><br />
                instead of reminding yourself to "push often" and "keep the PR small".
              </span>
            </MiniUsp>
          </div>
        </div>
        <div class="mt-48">
          <h1
            class="pb-6 text-center text-slate-200/80 uppercase text-sm font-bold tracking-widest"
          >
            What developers are saying
          </h1>
          <div class="mt-6 flex gap-6 overflow-x-auto fade max-w-6xl px-6 mx-auto">
            <Testimonial name="Kevin Simons" title="CTO, Kitemaker">
              <template #logo
                ><img src="./assets/kitemaker.svg" alt="Kitemaker logo" height="26"
              /></template>
              <template #quote>
                Sturdy makes it extremely easy for us to explore and try each other’s ideas during
                development. I can thoroughly recommend it for fast moving teams.
              </template>
              <template #face
                ><img
                  src="./assets/kevin_simons.jpeg"
                  alt="Kevin Simons"
                  class="h-10 w-10 border-2 border-amber-500/50 rounded-full"
                  height="40"
                  width="40"
              /></template>
            </Testimonial>
            <Testimonial name="&mdash; Will Stensvold">
              <template #quote>
                I am just starting out on the adventure of being a freelance programmer and was
                looking for a tool with the functionality of Git but without the overhead of having
                to learn cryptic commands that obfuscate what I wanted to do with my code.
              </template>
            </Testimonial>
            <Testimonial name="&mdash; Slamet Hendry">
              <template #quote>
                We expect our docs and pictures to sync to the cloud automatically to Dropbox or
                iCloud or whatever service we use.
                <br />
                Similarly, Sturdy automatically sync our code to the repo and lets us focus on
                coding and collaboration.
              </template>
            </Testimonial>
          </div>
        </div>
      </div>
    </section>
    <div class="mb-20 space-y-20 overflow-hidden sm:mb-32 sm:space-y-32 md:mb-40 md:space-y-40">
      <Usp
        title="Expressive"
        subtitle="Don't hand-hold your version control."
        link="v2DocsHowToEditCode"
        accent="amber-500"
      >
        <template #icon>
          <IconDevelop class="border-amber-500/80 border-2 rounded-full w-20 h-20" />
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
        accent="purple-400"
        :bigger-right-side="true"
      >
        <template #icon>
          <IconGears class="border-purple-400/80 border-2 rounded-full w-20 h-20"></IconGears>
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
        <template #right>
          <Video
            :src="collabVideo"
            style="filter: drop-shadow(0 0 1.5rem rgba(192, 132, 252, 0.14))"
          />
        </template>
        <template #bottom>
          <div class="mt-8 flex flex-col sm:flex-row gap-10 tracking-tight">
            <div
              :class="[
                collabSelection === 'trying-code'
                  ? 'border-purple-400/30'
                  : 'border-transparent cursor-pointer hover:bg-purple-400/5',
                'flex flex-col p-2 border rounded',
              ]"
              @click="collabSelection = 'trying-code'"
            >
              <h1 class="flex items-center text-purple-400 font-semibold tracking-tight highlight">
                <PlayIcon class="w-8 h-8 mr-2" />
                Trying code
              </h1>
              <p class="mt-2">Get your your teammates' code with a single click.</p>

              <div
                :class="[
                  'mt-2 text-sm text-slate-400/80 font-semibold tracking-tight inline-flex group inline-flex items-center h-9',
                ]"
              >
                See it in action
                <ChevronRightIcon class="h-4 w-4 ml-1" />
              </div>
            </div>
            <div
              :class="[
                collabSelection === 'suggesting-code'
                  ? 'border-purple-400/30'
                  : 'border-transparent cursor-pointer hover:bg-purple-400/5',
                'flex flex-col p-2 border rounded',
              ]"
              @click="collabSelection = 'suggesting-code'"
            >
              <h1 class="flex items-center text-purple-400 font-semibold tracking-tight highlight">
                <LightBulbIcon class="w-8 h-8 mr-2" />
                Suggesting code
              </h1>
              <p class="mt-2">Give code suggestions by simpy typing in your IDE.</p>

              <div
                :class="[
                  'mt-2 text-sm text-slate-400/80 font-semibold tracking-tight inline-flex group inline-flex items-center h-9',
                ]"
              >
                See it in action
                <ChevronRightIcon class="h-4 w-4 ml-1" />
              </div>
            </div>
            <div
              :class="[
                collabSelection === 'taking-code'
                  ? 'border-purple-400/30'
                  : 'border-transparent cursor-pointer hover:bg-purple-400/5',
                'flex flex-col p-2 border rounded',
              ]"
              @click="collabSelection = 'taking-code'"
            >
              <h1 class="flex items-center text-purple-400 font-semibold tracking-tight highlight">
                <CheckCircleIcon class="w-8 h-8 mr-2" />
                Taking code
              </h1>
              <p class="mt-2">Accepted suggestions show up in your IDE immediately.</p>

              <div
                :class="[
                  'mt-2 text-sm text-slate-400/80 font-semibold tracking-tight inline-flex group inline-flex items-center h-9',
                ]"
              >
                See it in action
                <ChevronRightIcon class="h-4 w-4 ml-1" />
              </div>
            </div>
          </div>
        </template>
      </Usp>

      <Usp
        title="Streamlined"
        subtitle="For the workflow that you're already using."
        link="v2DocsHotToShipSoftwareToProduction"
        accent="green-500"
        :bigger-right-side="true"
      >
        <template #icon>
          <IconDesktop class="border-green-500/80 border-2 rounded-full w-20 h-20"></IconDesktop>
        </template>
        <template #left>
          <p>
            We believe that, for most software, shipping small and often is more effective than
            developing in long-lived feature branches.
          </p>
          <p>
            We built Sturdy specifically around
            <span class="text-green-500">trunk-based development</span> and
            <span class="text-green-500">continuous delivery</span>. It optimizes for integration
            frequency and makes shipping of small incremental changes the intuitive default.
          </p>
        </template>
        <template #right>
          <Video
            :src="videoWorkflow"
            style="filter: drop-shadow(0 0 1.5rem rgba(34, 197, 94, 0.14))"
          ></Video>
        </template>
      </Usp>

      <Usp
        title="Compatible"
        subtitle="Worried about compatibility? Don't be."
        link="v2DocsHowSturdyInteractsWithGit"
        accent="amber-500"
        :bigger-right-side="true"
      >
        <template #icon>
          <IconMobile class="border-amber-500/80 border-2 rounded-full w-20 h-20"></IconMobile>
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
        <template #bottomLeft></template>
        <template #right>
          <Video
            :src="videoCreatePr"
            style="filter: drop-shadow(0 0 1.5rem rgba(245, 158, 11, 0.14))"
          ></Video>
        </template>
      </Usp>

      <Usp title="Open-source" subtitle="Join the community." accent="purple-400">
        <template #icon>
          <IconLaptop class="border-purple-400/80 border-2 rounded-full w-20 h-20"></IconLaptop>
        </template>
        <template #left>
          <p>
            Sturdy is made by engineers for other engineers. As
            <span class="text-purple-400">tool builders</span> at heart we are passionately
            exploring what is possible.
          </p>
          <p>
            Join our community. Request features. Ask questions. Get help. Hang out on Discord with
            our devs. Contribute code. Or even

            <router-link :to="{ name: 'careers' }" class="text-purple-400 hover:text-purple-300">
              join the Sturdy team.
            </router-link>
          </p>

          <a
            class="hover:text-slate-300 hover:border-slate-300 hover:bg-slate-900/10 bg-slate-600/80 text-slate-200 border-transparent border focus:outline-none focus:ring-2 focus:ring-slate-400 focus:ring-offset-2 focus:ring-offset-slate-50 focus:outline-none focus:ring-2 group inline-flex items-center h-9 rounded text-sm font-semibold whitespace-nowrap px-3"
            href="https://discord.gg/fQcH9QAVpX"
          >
            <svg class="h-6 w-6" fill="currentColor" viewBox="0 0 71 55" aria-hidden="true">
              <path
                fill-rule="evenodd"
                d="M60.1045 4.8978C55.5792 2.8214 50.7265 1.2916 45.6527 0.41542C45.5603 0.39851 45.468 0.440769 45.4204 0.525289C44.7963 1.6353 44.105 3.0834 43.6209 4.2216C38.1637 3.4046 32.7345 3.4046 27.3892 4.2216C26.905 3.0581 26.1886 1.6353 25.5617 0.525289C25.5141 0.443589 25.4218 0.40133 25.3294 0.41542C20.2584 1.2888 15.4057 2.8186 10.8776 4.8978C10.8384 4.9147 10.8048 4.9429 10.7825 4.9795C1.57795 18.7309 -0.943561 32.1443 0.293408 45.3914C0.299005 45.4562 0.335386 45.5182 0.385761 45.5576C6.45866 50.0174 12.3413 52.7249 18.1147 54.5195C18.2071 54.5477 18.305 54.5139 18.3638 54.4378C19.7295 52.5728 20.9469 50.6063 21.9907 48.5383C22.0523 48.4172 21.9935 48.2735 21.8676 48.2256C19.9366 47.4931 18.0979 46.6 16.3292 45.5858C16.1893 45.5041 16.1781 45.304 16.3068 45.2082C16.679 44.9293 17.0513 44.6391 17.4067 44.3461C17.471 44.2926 17.5606 44.2813 17.6362 44.3151C29.2558 49.6202 41.8354 49.6202 53.3179 44.3151C53.3935 44.2785 53.4831 44.2898 53.5502 44.3433C53.9057 44.6363 54.2779 44.9293 54.6529 45.2082C54.7816 45.304 54.7732 45.5041 54.6333 45.5858C52.8646 46.6197 51.0259 47.4931 49.0921 48.2228C48.9662 48.2707 48.9102 48.4172 48.9718 48.5383C50.038 50.6034 51.2554 52.5699 52.5959 54.435C52.6519 54.5139 52.7526 54.5477 52.845 54.5195C58.6464 52.7249 64.529 50.0174 70.6019 45.5576C70.6551 45.5182 70.6887 45.459 70.6943 45.3942C72.1747 30.0791 68.2147 16.7757 60.1968 4.9823C60.1772 4.9429 60.1437 4.9147 60.1045 4.8978ZM23.7259 37.3253C20.2276 37.3253 17.3451 34.1136 17.3451 30.1693C17.3451 26.225 20.1717 23.0133 23.7259 23.0133C27.308 23.0133 30.1626 26.2532 30.1066 30.1693C30.1066 34.1136 27.28 37.3253 23.7259 37.3253ZM47.3178 37.3253C43.8196 37.3253 40.9371 34.1136 40.9371 30.1693C40.9371 26.225 43.7636 23.0133 47.3178 23.0133C50.9 23.0133 53.7545 26.2532 53.6986 30.1693C53.6986 34.1136 50.9 37.3253 47.3178 37.3253Z"
                clip-rule="evenodd"
              ></path>
            </svg>
            <span class="ml-2">Discord server</span>
          </a>
        </template>
        <template #right></template>
      </Usp>

      <section id="featured">
        <div class="relative">
          <div
            class="absolute border-t border-slate-700/30 block top-0 inset-x-0 h-[37.5rem] bg-gradient-to-b from-[#0c1322] top-0 xl:top-8"
          ></div>
          <div class="relative pt-16 max-w-6xl px-6 mx-auto lg:gap-8 space-y-6">
            <h1 class="pb-6 text-center text-slate-50 uppercase text-sm font-bold tracking-widest">
              Featured at
            </h1>
            <div
              class="sm:px-20 flex flex-col items-center space-y-6 sm:flex-row sm:justify-between sm:space-y-0 sm:space-x-10 sm:items-center"
            >
              <a href="https://websummit.com/schedule/ws21/timeslot/breakout-startups-1104-0927">
                <img
                  src="./assets/web-summit.png"
                  alt="Web Summit logotype"
                  class="rounded-xl object-contain h-6"
                  height="24"
                />
              </a>
              <a href="https://devsecops.fm/episodes/36-sturdy/">
                <img
                  src="./assets/devsecops.png"
                  alt="DevSecOps logotype"
                  class="rounded-xl object-contain h-8"
                  height="24"
                />
              </a>
              <a href="https://dev-con.dev/">
                <img
                  src="./assets/devcon.png"
                  alt="Devcon logotype"
                  class="rounded-xl object-contain h-8"
                  height="24"
                />
              </a>
            </div>
          </div>
        </div>
      </section>

      <section id="get-started" class="relative">
        <div class="px-6 pt-16 max-w-6xl mx-auto items-center text-center">
          <h1 class="text-slate-50 uppercase text-sm font-bold tracking-widest">
            Begin your Sturdy journey
          </h1>

          <div class="mt-12 flex flex-col md:flex-row gap-20">
            <div class="flex flex-col md:w-1/2 border rounded border-slate-800 bg-slate-800/50">
              <h2 class="mx-6 mt-4 flex mt-6 text-2xl font-semibold text-amber-500 tracking-tight">
                Sturdy in the Cloud
              </h2>
              <p
                class="h-10 my-4 mx-6 flex text-center tracking-tight text-slate-400/80 font-semibold text-sm"
              >
                The easiest way of trying Sturdy.
              </p>
              <div
                class="shadow-[0_0px_0.5rem_2px_rgba(245,158,11,0.24)] flex flex-col h-60 grow shadow-xl border border-slate-300/30 bg-slate-700/50 rounded-xl mx-6 mx-auto"
              >
                <div class="my-10">
                  <p class="tracking-tight mx-10">
                    Try Sturdy with <span class="text-amber-500">zero setup</span> using our
                    free-forever tier.
                  </p>
                  <div
                    class="bg-gradient-to-t from-amber-500 to-amber-300/0 text-sm w-px h-16 text-slate-900 font-bold mx-auto flex flex-col justify-center items-center"
                  ></div>
                  <div class="w-2/3 mx-auto">
                    <router-link
                      :to="{ name: 'v2download' }"
                      class="text-slate-900 bg-amber-500 hover:text-amber-500 hover:bg-slate-900/10 border border-transparent hover:border-amber-500 focus:outline-none focus:ring-2 focus:ring-slate-400 focus:ring-offset-2 focus:ring-offset-slate-50 font-semibold h-12 px-6 rounded w-full flex items-center justify-center sm:w-auto highlight-white/20"
                    >
                      <DownloadIcon class="h-6 w-6 mr-1" />
                      {{ mainDownloadText }}
                    </router-link>
                  </div>
                </div>
              </div>

              <router-link
                :to="{ name: 'v2DocsRoot' }"
                class="mx-6 my-4 text-sm cursor-pointer text-slate-400/80 hover:text-slate-50/80 font-semibold tracking-tight group items-center h-9 flex"
              >
                Read the docs
                <ChevronRightIcon class="h-5 w-5" />
              </router-link>
            </div>
            <div class="md:w-1/2 border rounded border-slate-800 bg-slate-800/50">
              <div class="">
                <h2
                  class="mx-6 mt-4 flex mt-6 text-2xl font-semibold text-purple-400 tracking-tight"
                >
                  Run your own instance
                </h2>
                <p
                  class="h-10 my-4 flex mx-6 tracking-tight text-slate-400/80 font-semibold text-sm"
                >
                  Alternatively, run your own Sturdy instance.
                </p>

                <div class="mx-6 rounded-xl shadow-[0_0px_0.5rem_2px_rgba(192,132,252,0.24)]">
                  <CommandLineCopy :content="dockerOneliner" />
                </div>
              </div>

              <router-link
                :to="{ name: 'v2DocsSelfHosted' }"
                class="mx-6 my-4 text-sm cursor-pointer text-slate-400/80 hover:text-slate-50/80 font-semibold tracking-tight group items-center h-9 flex"
              >
                Read the docs
                <ChevronRightIcon class="h-5 w-5" />
              </router-link>
            </div>
          </div>
        </div>
      </section>
    </div>
  </DocumentationStickyHeader>
</template>

<script lang="ts">
import Video from '../../molecules/landing/Video.vue'
import Usp from '../../organisms/landing/Usp.vue'
import MiniUsp from '../../organisms/landing/MiniUsp.vue'
import Testimonial from '../../organisms/landing/Testimonial.vue'
import ThingsYouAreNotDoing from '../../organisms/landing/ThingsYouAreNotDoing.vue'
import IconDevelop from '../../molecules/landing/IconDevelop.vue'
import IconGears from '../../molecules/landing/IconGears.vue'
import IconDesktop from '../../molecules/landing/IconDesktop.vue'
import IconLaptop from '../../molecules/landing/IconLaptop.vue'
import IconMobile from '../../molecules/landing/IconMobile.vue'
import { defineComponent } from 'vue'
import { useHead } from '@vueuse/head'
import {
  CheckCircleIcon,
  ChevronRightIcon,
  CogIcon,
  DownloadIcon,
  FireIcon,
  LightBulbIcon,
  PlayIcon,
  RefreshIcon,
} from '@heroicons/vue/outline'
import CommandLineCopy from '../../molecules/CommandLineCopy.vue'
import DocumentationStickyHeader from '../../layouts/DocumentationStickyHeader.vue'

import videoWorkflow from './assets/workflow.mp4'
import videoCreatePr from './assets/create-pr.mp4'
import videoAcceptSuggestion from './assets/accept-suggestion.mp4'
import videoGiveSuggestion from './assets/give-suggestion.mp4'
import videoTryCode from './assets/try-code.mp4'
import videoCollaborate from './assets/collaborate.mp4'
import { DOCKER_ONELINER } from '../../docker'

export default defineComponent({
  components: {
    DocumentationStickyHeader,
    Video,
    Usp,
    MiniUsp,
    Testimonial,
    ThingsYouAreNotDoing,
    IconDevelop,
    IconGears,
    IconDesktop,
    IconLaptop,
    IconMobile,
    DownloadIcon,
    ChevronRightIcon,
    CheckCircleIcon,
    PlayIcon,
    LightBulbIcon,
    FireIcon,
    CogIcon,
    RefreshIcon,
    CommandLineCopy,
  },

  props: {
    light: {
      type: Boolean,
      required: true,
    },
  },

  setup() {
    useHead({
      title: 'Sturdy - Real-time code collaboration',
      meta: [
        {
          name: 'description',
          content:
            'Sturdy is an open-source version control platform that allows you to interact with your code at a higher abstraction level.',
        },
        {
          name: 'keywords',
          content: 'study code collaboration git vcs',
        },
        { property: 'og:title', content: 'Real-time code collaboration.' },
        {
          property: 'og:description',
          content:
            'Sturdy is an open-source version control platform that allows you to interact with your code at a higher abstraction level.',
        },
        // { property: 'og:image', content: '' },
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
      dockerOneliner: DOCKER_ONELINER,

      videoWorkflow,
      videoCreatePr,
      videoCollaborate,
    }
  },

  data: function () {
    return {
      collabSelection: 'trying-code',
    }
  },

  computed: {
    collabVideo() {
      switch (this.collabSelection) {
        case 'trying-code':
          return videoTryCode
        case 'suggesting-code':
          return videoGiveSuggestion
        case 'taking-code':
          return videoAcceptSuggestion
        default:
          return videoAcceptSuggestion
      }
    },
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

.fade {
  mask-image: linear-gradient(90deg, #000 80%, transparent);
}
</style>
