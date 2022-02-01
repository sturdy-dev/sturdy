<template>
  <div v-if="!ssr" class="relative bg-gray-50 w-full">
    <Popover
      class="relative bg-gray-800 flex items-center"
      :class="narrow ? 'h-16 border-b  border-warmgray-300 ' : 'h-24'"
    >
      <div class="max-w-7xl mx-auto px-4 sm:px-6 flex-1">
        <div class="flex justify-between items-center md:justify-start md:space-x-10">
          <!-- Logo on md -->
          <div v-if="!narrow" class="md:flex justify-start hidden lg:w-0 lg:flex-1">
            <router-link :to="{ name: 'index' }" class="flex space-x-2">
              <span class="sr-only">Sturdy</span>
              <img
                class="h-8 w-auto sm:h-10"
                src="../assets/Web/Logo/Yellow482x.png"
                height="96"
                width="253"
                alt="Sturdy Logo"
              />
            </router-link>
          </div>
          <!-- Logo on < md -->
          <div class="flex justify-start md:hidden">
            <router-link :to="{ name: 'index' }" class="flex space-x-2">
              <span class="sr-only">Sturdy</span>
              <img
                class="h-8 w-auto sm:h-10"
                src="../assets/Web/Logo/Yellow482x.png"
                height="96"
                width="253"
                alt="Sturdy Logo"
              />
            </router-link>
          </div>
          <div v-if="narrow" class="hidden lg:flex-1"></div>
          <div class="-mr-2 -my-2 md:hidden">
            <PopoverButton
              class="bg-white rounded-md p-2 inline-flex items-center justify-center text-gray-400 hover:text-gray-500 hover:bg-gray-100 focus:outline-none focus:ring-2 focus:ring-inset focus:ring-yellow-500"
            >
              <span class="sr-only">Open menu</span>
              <MenuIcon class="h-6 w-6" aria-hidden="true" />
            </PopoverButton>
          </div>
          <PopoverGroup as="nav" class="hidden md:flex space-x-10">
            <Popover v-slot="{ open }" class="relative">
              <PopoverButton
                :class="[
                  open ? 'text-gray-300' : 'text-gray-100',
                  'border-0 group bg-gray-800 rounded-md inline-flex items-center text-base font-medium hover:text-gray-300 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-yellow-500 focus:ring-opacity-50 focus:ring-offset-transparent',
                ]"
              >
                <span>Solutions</span>
                <ChevronDownIcon
                  :class="[
                    open ? 'text-gray-300' : 'text-gray-100',
                    'ml-2 h-5 w-5 group-hover:text-gray-300',
                  ]"
                  aria-hidden="true"
                />
              </PopoverButton>

              <transition
                enter-active-class="transition ease-out duration-200"
                enter-from-class="opacity-0 translate-y-1"
                enter-to-class="opacity-100 translate-y-0"
                leave-active-class="transition ease-in duration-150"
                leave-from-class="opacity-100 translate-y-0"
                leave-to-class="opacity-0 translate-y-1"
              >
                <PopoverPanel
                  v-slot="{ close }"
                  class="absolute -ml-4 mt-3 transform z-50 px-2 w-screen max-w-md sm:px-0"
                >
                  <div
                    class="rounded-lg shadow-lg ring-1 ring-black ring-opacity-5 overflow-hidden"
                  >
                    <div class="relative grid gap-6 bg-white px-5 py-6 sm:gap-8 sm:p-8">
                      <template v-for="item in features" :key="item.name">
                        <span
                          v-if="!item.routeName"
                          class="-m-3 p-3 flex items-start rounded-lg hover:bg-gray-50"
                        >
                          <component
                            :is="item.icon"
                            class="flex-shrink-0 h-6 w-6 text-yellow-500"
                            aria-hidden="true"
                          />
                          <div class="ml-4">
                            <div class="inline-flex items-center justify-between w-full">
                              <p class="text-base font-medium text-gray-900">
                                {{ item.name }}
                              </p>
                              <Pill v-if="item.comingSoon">Coming soon</Pill>
                            </div>
                            <p class="mt-1 text-sm text-gray-500">
                              {{ item.description }}
                            </p>
                          </div>
                        </span>
                        <router-link
                          v-else
                          :to="{ name: item.routeName }"
                          class="-m-3 p-3 flex items-start rounded-lg hover:bg-gray-50"
                          @click="close()"
                        >
                          <component
                            :is="item.icon"
                            class="flex-shrink-0 h-6 w-6 text-yellow-500"
                            aria-hidden="true"
                          />
                          <div class="ml-4">
                            <div class="inline-flex items-center justify-between w-full">
                              <p class="text-base font-medium text-gray-900">
                                {{ item.name }}
                              </p>
                              <Pill v-if="item.comingSoon">Coming soon</Pill>
                            </div>
                            <p class="mt-1 text-sm text-gray-500">
                              {{ item.description }}
                            </p>
                          </div>
                        </router-link>
                      </template>
                    </div>
                    <div
                      v-if="false"
                      class="px-5 py-5 bg-gray-50 space-y-6 sm:flex sm:space-y-0 sm:space-x-10 sm:px-8"
                    >
                      <div v-for="item in callsToAction" :key="item.name" class="flow-root">
                        <a
                          :href="item.href"
                          class="-m-3 p-3 flex items-center rounded-md text-base font-medium text-gray-900 hover:bg-gray-100"
                        >
                          <component
                            :is="item.icon"
                            class="flex-shrink-0 h-6 w-6 text-gray-400"
                            aria-hidden="true"
                          />
                          <span class="ml-3">{{ item.name }}</span>
                        </a>
                      </div>
                    </div>
                  </div>
                </PopoverPanel>
              </transition>
            </Popover>

            <router-link
              v-for="item in midLinks"
              :key="item.id"
              class="text-base font-medium text-gray-100 hover:text-gray-300 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-yellow-500 focus:ring-opacity-50 focus:ring-offset-transparent rounded-md"
              :to="{ name: item.routeName }"
            >
              {{ item.name }}
            </router-link>

            <Popover v-slot="{ open }" class="relative">
              <PopoverButton
                :class="[
                  open ? 'text-gray-300' : 'text-gray-100',
                  'border-0 group bg-gray-800 rounded-md inline-flex items-center text-base font-medium hover:text-gray-300 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-yellow-500 focus:ring-opacity-50 focus:ring-offset-transparent',
                ]"
              >
                <span>More</span>
                <ChevronDownIcon
                  :class="[
                    open ? 'text-gray-600' : 'text-gray-400',
                    'ml-2 h-5 w-5 group-hover:text-gray-500',
                  ]"
                  aria-hidden="true"
                />
              </PopoverButton>

              <transition
                enter-active-class="transition ease-out duration-200"
                enter-from-class="opacity-0 translate-y-1"
                enter-to-class="opacity-100 translate-y-0"
                leave-active-class="transition ease-in duration-150"
                leave-from-class="opacity-100 translate-y-0"
                leave-to-class="opacity-0 translate-y-1"
              >
                <PopoverPanel
                  v-slot="{ close }"
                  class="absolute left-1/2 z-50 transform -translate-x-1/4 mt-3 px-2 w-screen max-w-md sm:px-0"
                >
                  <div
                    class="rounded-lg shadow-lg ring-1 ring-black ring-opacity-5 overflow-hidden"
                  >
                    <div class="relative grid gap-6 bg-white px-5 py-6 sm:gap-8 sm:p-8">
                      <template v-for="item in resources" :key="item.name">
                        <router-link
                          v-if="item.routeName"
                          :to="{ name: item.routeName }"
                          class="-m-3 p-3 flex items-start rounded-lg hover:bg-gray-50"
                          @click="close()"
                        >
                          <component
                            :is="item.icon"
                            class="flex-shrink-0 h-6 w-6 text-yellow-500"
                            aria-hidden="true"
                          />
                          <div class="ml-4">
                            <p class="text-base font-medium text-gray-900">
                              {{ item.name }}
                            </p>
                            <p class="mt-1 text-sm text-gray-500">
                              {{ item.description }}
                            </p>
                          </div>
                        </router-link>
                        <a
                          v-else
                          :href="item.href"
                          class="-m-3 p-3 flex items-start rounded-lg hover:bg-gray-50"
                        >
                          <component
                            :is="item.icon"
                            class="flex-shrink-0 h-6 w-6 text-yellow-500"
                            aria-hidden="true"
                          />
                          <div class="ml-4">
                            <p class="text-base font-medium text-gray-900">
                              {{ item.name }}
                            </p>
                            <p class="mt-1 text-sm text-gray-500">
                              {{ item.description }}
                            </p>
                          </div>
                        </a>
                      </template>
                    </div>
                    <div class="px-5 py-5 bg-gray-50 sm:px-8 sm:py-8">
                      <div>
                        <h3 class="text-sm tracking-wide font-medium text-gray-500 uppercase">
                          Recent Posts
                        </h3>
                        <ul class="mt-4 space-y-4">
                          <li
                            v-for="item in recentPosts"
                            :key="item.name"
                            class="text-base truncate"
                          >
                            <router-link
                              :to="{ name: item.name }"
                              class="font-medium text-gray-900 hover:text-gray-700"
                              @click="close()"
                            >
                              {{ item.meta.blog.title }}
                            </router-link>
                          </li>
                        </ul>
                      </div>
                      <div class="mt-5 text-sm">
                        <router-link
                          :to="{ name: 'blog' }"
                          class="font-medium text-yellow-500 hover:text-yellow-600"
                          @click="close()"
                        >
                          View all posts <span aria-hidden="true">&rarr;</span>
                        </router-link>
                      </div>
                    </div>
                  </div>
                </PopoverPanel>
              </transition>
            </Popover>
          </PopoverGroup>
          <div class="hidden md:flex items-center justify-end md:flex-1 lg:w-0">
            <router-link
              v-if="user"
              :to="{ name: 'home' }"
              type="button"
              class="inline-flex justify-center px-4 py-2 border border-gray-300 shadow-sm text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-gray-900"
            >
              <span>Open Sturdy</span>
              <ArrowCircleRightIcon class="-mr-1 ml-2 h-5 w-5 text-green-700" />
            </router-link>
            <template v-else>
              <router-link
                :to="{ name: 'download' }"
                class="ml-8 whitespace-nowrap inline-flex items-center justify-center px-4 py-2 border border-transparent rounded-md shadow-sm text-base font-medium text-black bg-yellow-400 hover:bg-yellow-500"
              >
                Get started
              </router-link>
            </template>
          </div>
        </div>
      </div>

      <transition
        enter-active-class="duration-200 ease-out"
        enter-from-class="opacity-0 scale-95"
        enter-to-class="opacity-100 scale-100"
        leave-active-class="duration-100 ease-in"
        leave-from-class="opacity-100 scale-100"
        leave-to-class="opacity-0 scale-95"
      >
        <PopoverPanel
          v-slot="{ close }"
          focus
          class="absolute top-0 inset-x-0 z-50 p-2 transition transform origin-top-right md:hidden"
        >
          <div
            class="rounded-lg shadow-lg ring-1 ring-black ring-opacity-5 bg-white divide-y-2 divide-gray-50"
          >
            <div class="pt-5 pb-6 px-5">
              <div class="flex items-center justify-between">
                <div>
                  <img
                    class="h-8 w-auto"
                    src="../assets/Web/Logo/Yellow482x.png"
                    height="96"
                    width="253"
                    alt="Sturdy Logo"
                  />
                </div>
                <div class="-mr-2">
                  <PopoverButton
                    class="bg-wh rounded-md p-2 inline-flex items-center justify-center text-gray-400 hover:text-gray-500 hover:bg-gray-100 focus:outline-none focus:ring-2 focus:ring-inset focus:ring-yellow-500"
                  >
                    <span class="sr-only">Close menu</span>
                    <XIcon class="h-6 w-6" aria-hidden="true" />
                  </PopoverButton>
                </div>
              </div>
              <div class="mt-6">
                <nav class="grid gap-y-8">
                  <router-link
                    v-for="item in allLiveFeatures"
                    :key="item.name"
                    :to="{ name: item.routeName }"
                    class="-m-3 p-3 flex items-center rounded-md hover:bg-gray-50"
                    @click="close()"
                  >
                    <component
                      :is="item.icon"
                      class="flex-shrink-0 h-6 w-6 text-yellow-500"
                      aria-hidden="true"
                    />
                    <span class="ml-3 text-base font-medium text-gray-900">
                      {{ item.name }}
                    </span>
                  </router-link>
                </nav>
              </div>
            </div>
            <div class="py-6 px-5 space-y-6">
              <div class="grid grid-cols-2 gap-y-4 gap-x-8">
                <router-link
                  v-for="item in midLinks"
                  :key="item.id"
                  class="text-base font-medium text-gray-900 hover:text-gray-700"
                  :to="{ name: item.routeName }"
                  @click="close()"
                >
                  {{ item.name }}
                </router-link>

                <router-link
                  v-for="item in resources"
                  :key="item.name"
                  :to="{ name: item.routeName }"
                  class="text-base font-medium text-gray-900 hover:text-gray-700"
                  @click="close()"
                >
                  {{ item.name }}
                </router-link>
              </div>
              <div v-if="user" class="w-full">
                <router-link
                  :to="{ name: 'home' }"
                  type="button"
                  class="inline-flex justify-center px-4 py-2 border border-gray-300 shadow-sm text-sm font-medium rounded-md block w-full text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-gray-900"
                >
                  <span>Open Sturdy</span>
                  <ArrowCircleRightIcon class="-mr-1 ml-2 h-5 w-5 text-green-700" />
                </router-link>
              </div>
              <div v-else>
                <router-link
                  :to="{ name: 'download' }"
                  class="w-full flex items-center justify-center px-4 py-2 border border-transparent rounded-md shadow-sm text-base font-medium text-white bg-yellow-600 hover:bg-yellow-700"
                  @click="close()"
                >
                  Download now
                </router-link>
              </div>
            </div>
          </div>
        </PopoverPanel>
      </transition>
    </Popover>
  </div>
</template>

<script>
import { Popover, PopoverButton, PopoverGroup, PopoverPanel } from '@headlessui/vue'
import {
  ChartBarIcon,
  CursorClickIcon,
  MenuIcon,
  PhoneIcon,
  PlayIcon,
  RefreshIcon,
  ShieldCheckIcon,
  SupportIcon,
  BriefcaseIcon,
  ViewGridIcon,
  XIcon,
  AtSymbolIcon,
  LightningBoltIcon,
  LightBulbIcon,
  PuzzleIcon,
  PhotographIcon,
  LockClosedIcon,
} from '@heroicons/vue/outline'
import { ChevronDownIcon, ArrowCircleRightIcon } from '@heroicons/vue/solid'
import Pill from './shared/Pill.vue'
import { useRouter } from 'vue-router'

const features = [
  {
    name: 'Live feedback',
    routeName: 'featuresLiveFeedback',
    description: 'Share and review in real-time',
    icon: ChartBarIcon,
  },
  {
    name: 'Instant workspace switching',
    routeName: 'featuresWorkspaceNavigation',
    description: 'Explore code with one click',
    icon: CursorClickIcon,
  },

  {
    name: 'Conflicts',
    routeName: 'featuresConflicts',
    description: 'Conflicts happen! How to prevent and resolve them in Sturdy',
    icon: LightBulbIcon,
  },
  // {
  //   name: 'A workflow for teams',
  //   routeName: 'featuresWorkflow',
  //   comingSoon: true,
  //   description: 'Sturdy is built after the way modern software teams are working.',
  //   icon: CursorClickIcon,
  // },
  {
    name: 'Integrations',
    routeName: 'featuresIntegrations',
    description: "Connect with third-party tools that you're already using.",
    icon: ViewGridIcon,
  },
  {
    name: 'Large files',
    routeName: 'featuresLargeFiles',
    description: 'Work with gigabyte sized files on Sturdy',
    icon: PhotographIcon,
  },
  {
    name: 'Access Control',
    routeName: 'featuresAccessControl',
    description: 'Principle of least privilege, for files!',
    icon: LockClosedIcon,
  },
  {
    name: 'Instant Integration',
    comingSoon: true,
    routeName: 'featuresInstantIntegration',
    description: 'The future of Continuous Integration, coming soon to Sturdy!',
    icon: RefreshIcon,
  },
]
const callsToAction = [
  { name: 'Watch Demo', href: '#', icon: PlayIcon },
  { name: 'Contact Sales', href: '#', icon: PhoneIcon },
]
const resources = [
  {
    name: 'Help Center',
    description: 'Get all of your questions answered from our documentation or contact support.',
    routeName: 'resourcesDocs',
    icon: SupportIcon,
  },
  {
    name: 'Migrate from GitHub',
    description: 'Learn how to migrate from GitHub to Sturdy without disrupting the team',
    routeName: 'resourcesMigrateFromGitHub',
    icon: LightningBoltIcon,
  },
  {
    name: 'Security',
    description: 'Understand how we take your privacy seriously.',
    routeName: 'resourcesSecurity',
    icon: ShieldCheckIcon,
  },
  {
    name: 'API',
    description: 'Integrate anything with Sturdy using the GraphQL API',
    routeName: 'resourcesApi',
    icon: PuzzleIcon,
  },
  {
    name: 'About Sturdy',
    description: 'Learn more about Sturdy, and how to contact us',
    routeName: 'about',
    icon: AtSymbolIcon,
  },
  {
    name: 'Careers',
    description: 'Come and join the Sturdy team in Stockholm!',
    routeName: 'careers',
    icon: BriefcaseIcon,
  },
]

// const recentPosts = [
//   {
//     id: 1,
//     name: 'Humane Code Review',
//     href: 'https://sturdy.substack.com/p/005-humane-code-review',
//   },
//   {
//     id: 2,
//     name: 'Importing from Git!',
//     href: 'https://sturdy.substack.com/p/004-importing-from-git',
//   },
//   { id: 3, name: 'Share now!', href: 'https://sturdy.substack.com/p/003-share-now' },
// ]

const midLinks = [
  {
    id: 1,
    name: 'Pricing',
    routeName: 'pricing',
  },
  // {
  //   id: 2,
  //   name: 'Docs',
  //   routeName: 'resourcesDocs',
  // },
]

export default {
  components: {
    Popover,
    PopoverButton,
    PopoverGroup,
    PopoverPanel,
    ChevronDownIcon,
    ArrowCircleRightIcon,
    BriefcaseIcon,
    MenuIcon,
    XIcon,
    Pill,
    PhotographIcon,
  },
  props: ['user', 'narrow'],
  setup() {
    let routes = useRouter().getRoutes()

    return {
      features,
      callsToAction,
      resources,
      // recentPosts,
      midLinks,
      recentPosts: routes.filter((r) => r.meta.blog).slice(0, 3),

      ssr: import.meta.env.SSR,
    }
  },
  computed: {
    allLiveFeatures() {
      return this.features.filter((f) => !f.comingSoon)
    },
  },
}
</script>
