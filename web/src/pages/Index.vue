<template>
  <div class="bg-gray-800">
    <div class="pt-12">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div
          class="w-full text-white text-center md:text-left flex flex-col items-center space-y-4"
        >
          <div class="space-y-2">
            <router-link
              :to="{ name: 'featuresInstantIntegration' }"
              class="inline-flex items-center text-white bg-black rounded-full p-1 pr-2 sm:text-base lg:text-sm xl:text-base hover:text-gray-200"
            >
              <span
                class="px-3 py-0.5 text-white text-xs font-semibold leading-5 uppercase tracking-wide bg-yellow-500 rounded-full"
              >
                Introducing
              </span>
              <span class="ml-4 text-sm hidden lg:inline-block">
                Instant Integration &mdash; Builds and tests your code as you type it.</span
              >
              <span class="ml-4 text-sm lg:hidden inline-block"> Instant Integration</span>
              <ChevronRightIcon class="ml-2 w-5 h-5 text-gray-500" aria-hidden="true" />
            </router-link>

            <h1 class="text-4xl lg:text-6xl font-bold">Real-Time Version Control</h1>
          </div>

          <h3 class="text-xl lg:text-2xl">
            Share, try, and iterate on code together using Sturdy’s supercharged workspaces. It’s
            like Slack instead of email.
          </h3>
        </div>
      </div>
    </div>

    <div v-if="false" class="my-16">
      <div v-if="waitingListMessage" class="justify-center flex">
        <div class="max-w-3xl flex-1 text-center pb-4 text-white font-bold">
          <span>{{ waitingListMessage }}</span>
        </div>
      </div>
      <div v-if="showWaitingListButton" class="justify-center flex">
        <button
          class="inline-flex items-center text-black bg-yellow-400 hover:bg-yellow-500 font-semibold focus:outline-none focus:ring-2 focus:ring-amber-600 focus:ring-opacity-50 rounded p-3 md:p-2 lg:p-3 shadow-sm"
          @click="onClickRequestAccess"
        >
          <span class="pr-4">Request Access</span>
          <svg
            class="h-4 w-4"
            xmlns="http://www.w3.org/2000/svg"
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M14 5l7 7m0 0l-7 7m7-7H3"
            />
          </svg>
        </button>
      </div>

      <div v-if="showWaitingListForm" class="justify-center flex">
        <div class="max-w-3xl flex-1">
          <form>
            <div class="md:flex align-middle">
              <div class="flex-grow">
                <input
                  id="email"
                  v-model="email"
                  type="email"
                  name="email"
                  required="required"
                  autocomplete="email"
                  class="block w-full text-black shadow-sm focus:ring-red-500 focus:border-red-500 sm:text-sm border-gray-300 rounded-md"
                  placeholder="Email address"
                />
              </div>
              <div class="ml-3 mt-6 md:mt-0">
                <button
                  class="inline-flex items-center text-white bg-amber-400 hover:bg-amber-500 font-semibold focus:outline-none focus:ring-2 focus:ring-amber-600 focus:ring-opacity-50 rounded px-4 py-2"
                  @click.stop.prevent="joinWaitinglist"
                >
                  <span>Join the waitlist</span>
                </button>
              </div>
            </div>
          </form>
        </div>
      </div>
    </div>

    <div class="justify-center flex my-16">
      <div class="sm:flex sm:justify-center">
        <div class="rounded-md shadow">
          <router-link
            :to="{ name: 'download' }"
            class="w-full flex items-center justify-center px-8 py-3 border border-transparent text-base font-medium rounded-md text-black bg-yellow-400 hover:bg-yellow-500 md:py-4 md:text-lg md:px-10"
          >
            <DownloadIcon class="h-6 w-6 mr-1" />
            {{ mainDownloadText }}
          </router-link>
        </div>
        <div class="mt-3 rounded-md shadow sm:mt-0 sm:ml-3">
          <router-link
            :to="{ name: 'resourcesDocs' }"
            class="w-full flex items-center justify-center px-8 py-3 border border-transparent text-base font-medium rounded-md text-yellow-500 bg-white hover:bg-gray-50 md:py-4 md:text-lg md:px-10"
          >
            Read the docs
          </router-link>
        </div>
      </div>
    </div>

    <div class="justify-center flex my-16">
      <img
        src="../assets/Web/Duck/SwimmingDuck512.png"
        style="height: 200px; width: auto"
        width="512"
        height="512"
      />
    </div>

    <div class="mx-auto px-4 max-w-5xl">
      <div>
        <div
          v-if="!videoClicked"
          class="aspect-w-16 aspect-h-10 cursor-pointer"
          @click="clickVideo"
        >
          <img src="../assets/Landing/SturdyDemoScreenShot.png" class="rounded-xl object-contain" />
        </div>
        <div v-else class="aspect-w-16 aspect-h-10">
          <iframe
            class="rounded-xl"
            width="960"
            height="600"
            src="https://www.youtube-nocookie.com/embed/D54x6wBRCH0?autoplay=1"
            title="Sturdy demo"
            frameborder="0"
            allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture"
            allowfullscreen
          ></iframe>
        </div>
      </div>
    </div>

    <div class="relative py-16 sm:py-24 lg:py-32">
      <div class="mx-auto max-w-md px-4 text-center sm:max-w-3xl sm:px-6 lg:px-8 lg:max-w-7xl">
        <p class="mt-2 text-4xl md:text-5xl font-extrabold text-yellow-400 tracking-tight">
          Develop locally &mdash; with extra power from the cloud
        </p>
        <p class="my-8 max-w-prose mx-auto text-xl text-gray-300">
          Sturdy automatically keeps track of the changes you make locally and makes them instantly
          available for you and your team. This gives you the power of the cloud, while keeping the
          convenience of your existing development tools.
        </p>
      </div>
      <div class="max-w-4xl mx-auto px-2 md:px-0">
        <img
          src="../assets/Landing/Sturdy-code-collaboration.svg"
          class="rounded-xl object-contain"
          width="1963"
          height="1218"
        />
      </div>

      <div
        class="mt-10 mx-auto max-w-md px-4 text-cente lg:space-x-5 flex flex-col items-center lg:flex-row space-y-5 lg:space-y-0"
      >
        <router-link
          :to="{ name: 'download' }"
          class="inline-flex px-4 py-2 border border-transparent text-base font-medium rounded-md shadow-sm text-black bg-yellow-400 hover:bg-yellow-500"
        >
          Get started
        </router-link>
        <router-link
          :to="{ name: 'featuresLiveFeedback' }"
          class="inline-flex px-4 py-2 border border-transparent text-base font-medium rounded-md text-yellow-500 bg-white hover:bg-gray-50"
        >
          Learn about workspaces
        </router-link>
      </div>
    </div>

    <div ref="target" class="relative py-16 sm:py-24 lg:py-32">
      <div class="mx-auto max-w-md px-4 text-center sm:max-w-3xl sm:px-6 lg:px-8 lg:max-w-7xl">
        <p class="mt-2 text-5xl md:text-5xl font-extrabold text-yellow-400 tracking-tight">
          Look at all the things you're not doing!
        </p>
        <div
          class="my-8 max-w-prose mx-auto overflow-hidden text-lg sm:text-2xl md:text-2xl lg:text-4xl text-gray-100"
        >
          <div>Have an idea</div>
          <div>
            <span :class="[targetIsVisible ? 'strike' : '', 'text-gray-500 truncate font-mono']">
              git pull
            </span>
          </div>
          <div>
            <span
              :class="[
                targetIsVisible ? 'strikedelay-1 strike' : '',
                'text-gray-500 truncate font-mono',
              ]"
            >
              git checkout -b feature
            </span>
          </div>

          <div>Write the code</div>

          <div>
            <span
              :class="[
                targetIsVisible ? 'strikedelay-2 strike' : '',
                'text-gray-500 truncate font-mono',
              ]"
            >
              git add .
            </span>
          </div>
          <div>
            <span
              :class="[
                targetIsVisible ? 'strikedelay-3 strike' : '',
                'text-gray-500 truncate font-mono',
              ]"
            >
              git commit -m "My change"
            </span>
          </div>

          <div>
            <span
              :class="[
                targetIsVisible ? 'strikedelay-4 strike' : '',
                'text-gray-500 truncate font-mono hidden lg:inline-block',
              ]"
            >
              git push --set-upstream origin feature
            </span>
            <span
              :class="[
                targetIsVisible ? 'strikedelay-4 strike' : '',
                'text-gray-500 truncate font-mono lg:hidden',
              ]"
            >
              git push -u origin feature
            </span>
          </div>

          <div>
            <span
              :class="[
                targetIsVisible ? 'strikedelay-5 strike' : '',
                'text-gray-500 truncate font-mono',
              ]"
            >
              Create Pull Request
            </span>
          </div>
        </div>
        <div
          class="mt-10 mx-auto max-w-md px-4 text-center lg:space-x-5 flex flex-col items-center lg:flex-row space-y-5 lg:space-y-0"
        >
          <router-link
            :to="{ name: 'download' }"
            class="px-4 py-2 border border-transparent text-base font-medium rounded-md shadow-sm text-black bg-yellow-400 hover:bg-yellow-500"
          >
            Get started
          </router-link>
          <router-link
            :to="{ name: 'featuresLiveFeedback' }"
            class="px-4 py-2 border border-transparent text-base font-medium rounded-md text-yellow-500 bg-white hover:bg-gray-50"
          >
            Read about the workflow
          </router-link>
        </div>
      </div>
    </div>

    <div class="relative py-16 sm:py-24 lg:py-32">
      <div class="mx-auto max-w-md px-4 text-center sm:max-w-3xl sm:px-6 lg:px-8 lg:max-w-7xl">
        <p class="mt-2 text-4xl md:text-5xl font-extrabold text-yellow-400 tracking-tight">
          Humane code review
        </p>
        <p class="my-8 max-w-prose mx-auto text-xl text-gray-300">
          Let’s face it &mdash; a Pull Request review is either a nitpick fest or a blind LGTM. Our
          goal is to make meaningful collaboration the default. Getting somebody else’s code to your
          computer is a one click step. Creating code suggestions is just a matter of typing in your
          local editor.
        </p>
      </div>
      <div class="mx-auto flex flex-col lg:flex-row w-3/4">
        <div class="max-w-4xl mx-auto px-2 md:px-0">
          <img
            src="../assets/Web/SuggestionDarkBackground.png"
            class="rounded-md object-contain"
            height="393"
            width="859"
          />
        </div>
        <div class="max-w-4xl mx-auto px-2 md:px-0 relative top-6">
          <img
            src="../assets/Web/CommentDarkBackground.png"
            class="rounded-md object-contain"
            height="531"
            width="1003"
          />
        </div>
      </div>

      <div
        class="mt-10 mx-auto max-w-md px-4 text-center lg:space-x-5 flex flex-col items-center lg:flex-row space-y-5 lg:space-y-0"
      >
        <router-link
          :to="{ name: 'download' }"
          class="inline-flex px-4 py-2 border border-transparent text-base font-medium rounded-md shadow-sm text-black bg-yellow-400 hover:bg-yellow-500"
        >
          Get started
        </router-link>
        <router-link
          :to="{ name: 'featuresLiveFeedback' }"
          class="inline-flex px-4 py-2 border border-transparent text-base font-medium rounded-md text-yellow-500 bg-white hover:bg-gray-50"
        >
          Read about collaboration
        </router-link>
      </div>
    </div>

    <div class="relative py-16 sm:py-24 lg:py-32">
      <div class="mx-auto max-w-md px-4 text-center sm:max-w-3xl sm:px-6 lg:px-8 lg:max-w-7xl">
        <p class="mt-2 text-4xl md:text-5xl font-extrabold text-yellow-400 tracking-tight">
          Manual branching & pushing is so 2009
        </p>
        <p class="my-8 max-w-prose mx-auto text-xl text-gray-300">
          Your tools should empower you and allow you to focus on what you do best. This is why
          Sturdy is getting out of its way &mdash; to get out of your way. Sturdy saves and shares
          your changes automatically. Create as many workspaces as you want and switch between them
          effortlessly.
        </p>
      </div>
      <div class="max-w-4xl mx-auto px-2 md:px-0">
        <video
          autoplay
          loop
          muted
          playsinline
          alt="Navigating between Sturdy workspaces with one click"
          class="rounded-lg object-cover object-center"
        >
          <source src="../assets/Landing/WorkspaceSwitching.webm" type="video/webm" />
          <source src="../assets/Landing/WorkspaceSwitching.mp4" type="video/mp4" />
        </video>
      </div>

      <div
        class="mt-10 mx-auto max-w-md px-4 text-center lg:space-x-5 flex flex-col items-center lg:flex-row space-y-5 lg:space-y-0"
      >
        <router-link
          :to="{ name: 'download' }"
          class="inline-flex px-4 py-2 border border-transparent text-base font-medium rounded-md shadow-sm text-black bg-yellow-400 hover:bg-yellow-500"
        >
          Get started
        </router-link>
        <router-link
          :to="{ name: 'featuresWorkspaceNavigation' }"
          class="inline-flex px-4 py-2 border border-transparent text-base font-medium rounded-md text-yellow-500 bg-white hover:bg-gray-50"
        >
          Learn about workspaces
        </router-link>
      </div>
    </div>

    <div class="relative py-16 sm:py-24 lg:py-32">
      <div class="mx-auto max-w-md px-4 text-center sm:max-w-3xl sm:px-6 lg:px-8 lg:max-w-7xl">
        <p class="mt-2 text-4xl md:text-5xl font-extrabold text-yellow-400 tracking-tight">
          Works nicely with your existing tools
        </p>
        <p class="my-8 max-w-prose mx-auto text-xl text-gray-300">
          Sturdy and GitHub can be used side by side by different people on the same project. Your
          existing integrations and CI/CD will continue to work just like before. At the same time,
          staying in sync with your team will be easier than ever.
        </p>
      </div>
      <div
        class="flex mx-auto flex-col lg:flex-row w-3/4 space-y-4 lg:space-y-0 lg:space-x-2 items-center lg:items-start"
      >
        <div class="max-w-4xl mx-auto px-2 md:px-0">
          <img
            src="../assets/Web/GitHubPullRequestDark.png"
            class="rounded-md object-contain"
            height="1332"
            width="1890"
          />
        </div>
        <div class="max-w-4xl mx-auto px-2 md:px-0 relative lg:top-4 lg:w-10/12">
          <img
            src="../assets/Web/GitHubConfig.png"
            class="rounded-md object-contain"
            width="1292"
            height="630"
          />
        </div>
      </div>

      <div
        class="mt-10 mx-auto max-w-md px-4 text-center lg:space-x-5 flex flex-col items-center lg:flex-row space-y-5 lg:space-y-0"
      >
        <router-link
          :to="{ name: 'download' }"
          class="inline-flex px-4 py-2 border border-transparent text-base font-medium rounded-md shadow-sm text-black bg-yellow-400 hover:bg-yellow-500"
        >
          Get started
        </router-link>
        <router-link
          :to="{ name: 'resourcesMigrateFromGitHub' }"
          class="inline-flex px-4 py-2 border border-transparent text-base font-medium rounded-md text-yellow-500 bg-white hover:bg-gray-50"
        >
          How it works with GitHub
        </router-link>
      </div>
    </div>

    <div class="relative py-16 sm:py-24 lg:py-32">
      <div class="mx-auto max-w-md px-4 text-center sm:max-w-3xl sm:px-6 lg:px-8 lg:max-w-7xl">
        <h2 class="text-base font-semibold tracking-wider text-yellow-100 uppercase">
          Ship faster
        </h2>
        <p class="mt-2 text-3xl font-extrabold text-yellow-400 tracking-tight sm:text-4xl">
          Spend more time doing, less time undoing
        </p>
        <p class="mt-5 max-w-prose mx-auto text-xl text-white">
          Sturdy is the perfect workflow for group projects, teams, and startups.
        </p>
        <div class="mt-12">
          <div class="grid grid-cols-1 gap-8 sm:grid-cols-2 lg:grid-cols-3">
            <div v-for="usp in uspGrid" :key="usp.title" class="pt-6">
              <div class="flow-root bg-gray-50 rounded-lg px-6 pb-8 h-full">
                <div class="-mt-6">
                  <div class="">
                    <span
                      class="inline-flex items-center justify-center p-3 bg-yellow-400 rounded-md shadow-lg"
                    >
                      <component :is="usp.icon" class="h-6 w-6 text-black" aria-hidden="true" />
                    </span>
                  </div>
                  <h3 class="mt-8 text-lg font-medium text-gray-900 tracking-tight">
                    {{ usp.title }}
                  </h3>
                  <p class="mt-5 text-base text-gray-500">
                    {{ usp.text }}
                  </p>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <DiveIn />

    <div class="py-16 sm:py-24 lg:py-32 mx-auto">
      <span class="relative text-center mx-auto flex justify-center items-center">
        <h2 class="text-2xl md:text-4xl font-extrabold tracking-tight text-yellow-400">
          Backed by
        </h2>
        <img
          class="ml-8 h-12 w-12"
          src="../assets/Landing/ycombinator-logo.png"
          alt="Y Combinator"
          height="257"
          width="257"
        />
      </span>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import http from '../http'
import IndexInteractiveEditor from '../components/IndexInteractiveEditor.vue'
import DiveIn from '../components/index/DiveIn.vue'
import {
  AcademicCapIcon,
  CloudUploadIcon,
  CodeIcon,
  CollectionIcon,
  RefreshIcon,
  RewindIcon,
  ServerIcon,
  AnnotationIcon,
  EyeIcon,
  ChipIcon,
  LinkIcon,
  ChevronRightIcon,
  DownloadIcon,
} from '@heroicons/vue/outline'

import workspacesImage from '../assets/Web/Workspaces.png'
import snapshotsImage from '../assets/Web/Snapshots.png'
import fizzBuzzCommentImage from '../assets/Web/CommentDarkBackground.png'
import suggestionImage from '../assets/Web/SuggestionDarkBackground.png'
import pullRequestImage from '../assets/Web/GitHubPullRequest.png'
import { useHead } from '@vueuse/head'
import { ref } from 'vue'
import { useElementVisibility } from '@vueuse/core'

// import emilTullstedtAvatar from '../assets/Web/testimonial-emil-tullstedt.jpeg'

interface Data {
  showWaitingListButton: boolean
  showWaitingListForm: boolean
  waitingListMessage: string
  email: string
}

const largeUsps = [
  {
    title: 'Early feedback',
    description:
      "The earlier you give and take feedback the better. Sturdy live-streams everyone's code, so that you can give and take instant feedback.",
    icon: AnnotationIcon,
    image: fizzBuzzCommentImage,
    getStarted: true,
    // testimonial: {
    //   text: 'Sturdy allows me and my team to iterate and help each other faster than ever, I love it!',
    //   name: 'Emil Tullstedt, Software Engineer at Grafana Labs',
    //   image: emilTullstedtAvatar,
    // },
  },
  {
    title: 'Suggestions',
    description:
      'Helping a teammate out should not be hard, copy all of their code to your editor with one click 🤯<br>Suggest improvements by writing them in your editor.',
    icon: EyeIcon,
    image: suggestionImage,
  },
  {
    title: 'Build your project with Sturdy',
    description:
      'Sturdy is built for teams that are working together. Fast feedback on Sturdy makes it easy to ship changes many times per day.',
    icon: CodeIcon,
    image: workspacesImage,
  },
  {
    title: 'Undo anything',
    description:
      "If you make a mistake, or write some code that you want to back out of, it shouldn't be impossible to get back to where you started. Sturdy makes automatic snapshots in the background, so that you can go back to just how things were.",
    icon: CollectionIcon,
    image: snapshotsImage,
  },
  {
    title: 'Migrate from GitHub',
    description:
      'Easy migration from GitHub, where you can use Sturdy and GitHub side-by-side, without losing any work or getting in your way. Your existing integrations and CI/CD will continue to work like before.',
    icon: LinkIcon,
    image: pullRequestImage,
    width: 489,
    height: 937,
  },
]

const uspGrid = [
  {
    title: 'Automatic sync',
    text: "Sturdy automatically syncs your work. There's no push or pull. Sync your work with the team when you're ready.",
    icon: CloudUploadIcon,
  },

  {
    title: 'Use any editor',
    text: 'Use Visual Studio Code, IntelliJ, vim, or any other editor. Sturdy supports them all, and does not force you to use a web based editor.',
    icon: CodeIcon,
  },

  {
    title: 'Full history',
    text: 'Sturdy keeps track of all changes in your project, just like git. And easily allows you to restore or revert changes.',
    icon: ServerIcon,
  },

  {
    title: 'Undo',
    text: "It's OK to make mistakes. If something goes wrong, fixing it is always 1 click away. Revert to any point in history with a single click.",
    icon: RewindIcon,
  },

  {
    title: 'Real-time',
    text: "We want you to ship code as fast as possible. Sturdy removes everything that's in your way. Comment, suggest, and review code in real-time.",
    icon: RefreshIcon,
  },

  {
    title: 'No command line',
    text: "Control your code from the web. We know what you want to do, so that you don't have to remember.",
    icon: ChipIcon,
  },
]

export default defineComponent({
  components: {
    IndexInteractiveEditor,
    AcademicCapIcon,
    RewindIcon,
    RefreshIcon,
    ServerIcon,
    CloudUploadIcon,
    CollectionIcon,
    CodeIcon,
    ChipIcon,
    DiveIn,
    ChevronRightIcon,
    DownloadIcon,
  },
  props: ['user'],
  setup() {
    useHead({
      title: 'Sturdy - Real-time Version Control & Collaboration',
    })
    const target = ref(null)
    const targetIsVisible = useElementVisibility(target)

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
      largeUsps,
      uspGrid,
      target,
      targetIsVisible,
      mainDownloadText,
    }
  },
  data(): Data {
    return {
      showWaitingListButton: true, // Show large "Request Access" button
      showWaitingListForm: false,
      waitingListMessage: '',
      email: '',
      videoClicked: false,
    }
  },
  methods: {
    clickVideo() {
      this.videoClicked = true
    },
    joinWaitinglist() {
      fetch(http.url('v3/waitinglist'), {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          email: this.email,
        }),
      })
        .then((r) => {
          if (r.status == 200) {
            this.waitingListMessage = "Thanks! We'll soon be in touch with your invite."
            this.showWaitingListForm = false
          } else {
            this.waitingListMessage = 'Sorry, we where not able to sign you up. Please try again.'
          }
        })
        .catch(() => {
          this.waitingListMessage = 'Sorry, we where not able to sign you up. Please try again.'
        })
    },
    onClickRequestAccess() {
      this.showWaitingListButton = false
      this.showWaitingListForm = true
    },
  },
})
</script>
<style scoped>
@keyframes strike {
  0% {
    width: 0;
  }
  100% {
    width: 100%;
  }
}

.strike {
  position: relative;
}

.strikedelay-1::after {
  animation-delay: 0.5s;
}
.strikedelay-2::after {
  animation-delay: 1s;
}
.strikedelay-3::after {
  animation-delay: 1.5s;
}
.strikedelay-4::after {
  animation-delay: 2s;
}
.strikedelay-5::after {
  animation-delay: 2.5s;
}
.strikedelay-6::after {
  animation-delay: 3s;
}
.strike::after {
  content: ' ';
  position: absolute;
  top: 50%;
  left: 0;
  height: 3px;
  background: gray;
  animation-name: strike;
  animation-duration: 0.5s;
  animation-timing-function: linear;
  animation-iteration-count: 1;
  animation-fill-mode: forwards;
}

@media screen and (max-width: 992px) {
  .strike::after {
    height: 2px;
  }
}
</style>
