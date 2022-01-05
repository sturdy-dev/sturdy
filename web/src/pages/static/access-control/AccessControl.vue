<template>
  <div class="bg-white">
    <div class="relative overflow-hidden">
      <ClientOnly>
        <Popover as="header" class="relative">
          <transition
            enter-active-class="duration-150 ease-out"
            enter-from-class="opacity-0 scale-95"
            enter-to-class="opacity-100 scale-100"
            leave-active-class="duration-100 ease-in"
            leave-from-class="opacity-100 scale-100"
            leave-to-class="opacity-0 scale-95"
          >
            <PopoverPanel
              focus
              class="absolute top-0 inset-x-0 p-2 transition transform origin-top md:hidden"
            >
              <div
                class="rounded-lg shadow-md bg-white ring-1 ring-black ring-opacity-5 overflow-hidden"
              >
                <div class="px-5 pt-4 flex items-center justify-between">
                  <div>
                    <img
                      class="h-8 w-auto"
                      src="https://tailwindui.com/img/logos/workflow-mark-teal-200-cyan-400.svg"
                      alt=""
                    />
                  </div>
                  <div class="-mr-2">
                    <PopoverButton
                      class="bg-white rounded-md p-2 inline-flex items-center justify-center text-gray-400 hover:bg-gray-100 focus:outline-none focus:ring-2 focus:ring-inset focus:ring-yellow-600"
                    >
                      <span class="sr-only">Close menu</span>
                      <XIcon class="h-6 w-6" aria-hidden="true" />
                    </PopoverButton>
                  </div>
                </div>
                <div class="pt-5 pb-6">
                  <div class="px-2 space-y-1">
                    <a
                      v-for="item in navigation"
                      :key="item.name"
                      :href="item.href"
                      class="block px-3 py-2 rounded-md text-base font-medium text-gray-900 hover:bg-gray-50"
                      >{{ item.name }}</a
                    >
                  </div>
                  <div class="mt-6 px-5">
                    <a
                      href="#"
                      class="block text-center w-full py-3 px-4 rounded-md shadow bg-gradient-to-r from-yellow-500 to-yellow-600 text-white font-medium hover:from-yellow-600 hover:to-yellow-700"
                      >Start free trial</a
                    >
                  </div>
                  <div class="mt-6 px-5">
                    <p class="text-center text-base font-medium text-gray-500">
                      Existing customer? <a href="#" class="text-gray-900 hover:underline">Login</a>
                    </p>
                  </div>
                </div>
              </div>
            </PopoverPanel>
          </transition>
        </Popover>
      </ClientOnly>
      <main>
        <div class="pt-10 bg-gray-900 sm:pt-16 lg:pt-8 lg:pb-14 lg:overflow-hidden">
          <div class="mx-auto max-w-7xl lg:px-8">
            <div class="lg:grid lg:grid-cols-2 lg:gap-8">
              <div
                class="mx-auto max-w-md px-4 sm:max-w-2xl sm:px-6 sm:text-center lg:px-0 lg:text-left lg:flex lg:items-center"
              >
                <div class="lg:py-24">
                  <a
                    href="#"
                    class="inline-flex items-center text-white bg-black rounded-full p-1 pr-2 sm:text-base lg:text-sm xl:text-base hover:text-gray-200"
                  >
                    <span
                      class="px-3 py-0.5 text-white text-xs font-semibold leading-5 uppercase tracking-wide bg-gradient-to-r from-yellow-500 to-yellow-600 rounded-full"
                      >API Reference</span
                    >
                    <span class="ml-4 text-sm">Explore our GraphQL API</span>
                    <ChevronRightIcon class="ml-2 w-5 h-5 text-gray-500" aria-hidden="true" />
                  </a>
                  <h1
                    class="mt-4 text-4xl tracking-tight font-extrabold text-white sm:mt-5 sm:text-6xl lg:mt-6 xl:text-6xl"
                  >
                    <span class="block">Share a subfolder</span>
                    <span
                      class="pb-3 block bg-clip-text text-transparent bg-gradient-to-r from-yellow-200 to-yellow-400 sm:pb-5"
                      >of your Git repo</span
                    >
                  </h1>
                  <p class="text-base text-gray-300 sm:text-xl lg:text-lg xl:text-xl">
                    With <strong>{{ productName }}</strong> you can get code contributions to
                    specific files or folders without the need to share your entire Git repository.
                  </p>
                  <div class="mt-10 sm:mt-12">
                    <form action="#" class="sm:max-w-xl sm:mx-auto lg:mx-0">
                      <div class="sm:flex">
                        <div class="flex">
                          <div class="sm:flex sm:justify-center">
                            <div class="rounded-md shadow">
                              <router-link
                                :to="{ name: 'signup' }"
                                class="w-full py-4 px-4 rounded-md shadow bg-gradient-to-r from-yellow-500 to-yellow-600 text-white font-medium hover:from-yellow-600 hover:to-yellow-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-yellow-400 focus:ring-offset-gray-900"
                              >
                                Get started in the cloud
                              </router-link>

                              <p class="mt-10 text-gray-300 sm:mt-10">
                                Or contact us for enterprise setup.
                              </p>

                              <ACLRequestAccess />
                            </div>
                          </div>
                        </div>
                      </div>
                    </form>
                  </div>
                </div>
              </div>
              <div class="lg:py-24 mt-12 lg:-mb-16 lg:m-0 lg:relative">
                <div class="mx-auto max-w-md px-4 sm:max-w-sm sm:px-6 lg:max-w-none lg:px-0">
                  <img src="./sturdy-proxy.png" alt="" />
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Alternating Feature Sections -->
        <div class="relative pt-16 pb-32 overflow-hidden">
          <!-- img to the right-->
          <div
            aria-hidden="true"
            class="absolute inset-x-0 top-0 h-48 bg-gradient-to-b from-gray-100"
          />
          <div class="relative">
            <div
              class="lg:mx-auto lg:max-w-7xl lg:px-2 lg:grid lg:grid-cols-2 lg:grid-flow-col-dense lg:gap-12"
            >
              <div class="px-4 max-w-xl mx-auto sm:px-6 lg:py-16 lg:max-w-none lg:mx-0 lg:px-0">
                <div>
                  <div>
                    <span
                      class="h-12 w-12 rounded-md flex items-center justify-center bg-gradient-to-r from-yellow-500 to-yellow-600"
                    >
                      <InboxIcon class="h-6 w-6 text-white" aria-hidden="true" />
                    </span>
                  </div>
                  <div class="mt-6">
                    <h2 class="text-3xl font-extrabold tracking-tight text-gray-900">
                      {{ usps[0].title }}
                    </h2>
                    <p class="mt-4 text-lg text-gray-500">
                      {{ usps[0].subtitle }}
                    </p>
                    <div class="mt-6">
                      <router-link
                        :to="{ name: 'docsAccessControl' }"
                        class="inline-flex bg-gradient-to-r from-yellow-500 to-yellow-600 bg-origin-border px-4 py-2 border border-transparent text-base font-medium rounded-md shadow-sm text-white hover:from-yellow-700 hover:to-yellow-700"
                      >
                        Read more
                      </router-link>
                    </div>
                  </div>
                </div>
              </div>
              <div class="mt-12 sm:mt-16 lg:mt-0">
                <div class="sm:px-6 lg:py-16 lg:px-0 lg:m-0 lg:relative lg:h-full">
                  <video
                    autoplay
                    loop
                    muted
                    playsinline
                    alt="Navigating between Sturdy workspaces with one click"
                    class="rounded-lg object-cover object-center"
                  >
                    <source
                      src="../../../assets/Landing/WorkspaceSwitching.webm"
                      type="video/webm"
                    />
                    <source src="../../../assets/Landing/WorkspaceSwitching.mp4" type="video/mp4" />
                  </video>
                </div>
              </div>
            </div>
          </div>

          <!-- img to the left-->
          <div class="mt-12">
            <div
              class="lg:mx-auto lg:max-w-7xl lg:px-2 lg:grid lg:grid-cols-2 lg:grid-flow-col-dense lg:gap-12"
            >
              <div
                class="px-4 max-w-xl mx-auto sm:px-6 lg:py-32 lg:max-w-none lg:mx-0 lg:px-0 lg:col-start-2"
              >
                <div>
                  <div>
                    <span
                      class="h-12 w-12 rounded-md flex items-center justify-center bg-gradient-to-r from-yellow-500 to-yellow-600"
                    >
                      <SparklesIcon class="h-6 w-6 text-white" aria-hidden="true" />
                    </span>
                  </div>
                  <div class="mt-6">
                    <h2 class="text-3xl font-extrabold tracking-tight text-gray-900">
                      {{ usps[1].title }}
                    </h2>
                    <p class="mt-4 text-lg text-gray-500">
                      {{ usps[1].subtitle }}
                    </p>
                    <div v-if="false" class="mt-6">
                      <router-link
                        :to="{ name: 'signup' }"
                        class="inline-flex bg-gradient-to-r from-yellow-500 to-yellow-600 bg-origin-border px-4 py-2 border border-transparent text-base font-medium rounded-md shadow-sm text-white hover:from-yellow-700 hover:to-yellow-700"
                      >
                        Get started
                      </router-link>
                    </div>
                  </div>
                </div>
              </div>
              <div class="mt-12 sm:mt-16 lg:mt-20">
                <div class="sm:px-6 lg:py-16 lg:px-0 lg:m-0 lg:relative lg:h-full">
                  <img
                    class="rounded-lg object-cover object-center w-full"
                    src="./sturdy-acl.png"
                    alt="Sturdy ACL"
                  />
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Alternating Feature Sections -->
        <div class="relative pt-16 pb-32 overflow-hidden">
          <!-- img to the right-->
          <div
            aria-hidden="true"
            class="absolute inset-x-0 top-0 h-48 bg-gradient-to-b from-gray-100"
          />
          <div class="relative">
            <div
              class="lg:mx-auto lg:max-w-7xl lg:px-2 lg:grid lg:grid-cols-2 lg:grid-flow-col-dense lg:gap-12"
            >
              <div class="px-4 max-w-xl mx-auto sm:px-6 lg:py-16 lg:max-w-none lg:mx-0 lg:px-0">
                <div>
                  <div>
                    <span
                      class="h-12 w-12 rounded-md flex items-center justify-center bg-gradient-to-r from-yellow-500 to-yellow-600"
                    >
                      <InboxIcon class="h-6 w-6 text-white" aria-hidden="true" />
                    </span>
                  </div>
                  <div class="mt-6">
                    <h2 class="text-3xl font-extrabold tracking-tight text-gray-900">
                      {{ usps[2].title }}
                    </h2>
                    <p class="mt-4 text-lg text-gray-500">
                      {{ usps[2].subtitle }}
                    </p>
                    <div class="mt-6">
                      <router-link
                        :to="{ name: 'featuresIntegrations' }"
                        class="inline-flex bg-gradient-to-r from-yellow-500 to-yellow-600 bg-origin-border px-4 py-2 border border-transparent text-base font-medium rounded-md shadow-sm text-white hover:from-yellow-700 hover:to-yellow-700"
                      >
                        Read more
                      </router-link>
                    </div>
                  </div>
                </div>
              </div>
              <div class="mt-12 sm:mt-16 lg:mt-0">
                <div class="sm:px-6 lg:py-16 lg:px-0 lg:m-0 lg:relative lg:h-full">
                  <img
                    class="rounded-lg object-cover object-center w-full"
                    src="./sturdy-integrations.png"
                    alt="Sturdy integrations"
                  />
                </div>
              </div>
            </div>
          </div>

          <!-- img to the left-->
          <div class="mt-12">
            <div
              class="lg:mx-auto lg:max-w-7xl lg:px-2 lg:grid lg:grid-cols-2 lg:grid-flow-col-dense lg:gap-12"
            >
              <div
                class="px-4 max-w-xl mx-auto sm:px-6 lg:py-32 lg:max-w-none lg:mx-0 lg:px-0 lg:col-start-2"
              >
                <div>
                  <div>
                    <span
                      class="h-12 w-12 rounded-md flex items-center justify-center bg-gradient-to-r from-yellow-500 to-yellow-600"
                    >
                      <SparklesIcon class="h-6 w-6 text-white" aria-hidden="true" />
                    </span>
                  </div>
                  <div class="mt-6">
                    <h2 class="text-3xl font-extrabold tracking-tight text-gray-900">
                      {{ usps[3].title }}
                    </h2>
                    <p class="mt-4 text-lg text-gray-500">
                      {{ usps[3].subtitle }}
                    </p>
                    <div class="mt-6">
                      <router-link
                        :to="{ name: 'resourcesMigrateFromGitHub' }"
                        class="inline-flex bg-gradient-to-r from-yellow-500 to-yellow-600 bg-origin-border px-4 py-2 border border-transparent text-base font-medium rounded-md shadow-sm text-white hover:from-yellow-700 hover:to-yellow-700"
                      >
                        Read more
                      </router-link>
                    </div>
                  </div>
                </div>
              </div>
              <div class="mt-12 sm:mt-16 lg:mt-20">
                <div class="sm:px-6 lg:py-16 lg:px-0 lg:m-0 lg:relative lg:h-full">
                  <img
                    class="rounded-lg object-cover object-center w-full"
                    src="../../../assets/Web/GitHubPullRequest.png"
                    alt="Sturdy Pull request integration"
                  />
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- How it works section -->
        <div class="relative bg-gray-50 py-16 sm:py-24 lg:py-32">
          <div class="relative">
            <div
              class="text-center mx-auto max-w-md px-4 sm:max-w-3xl sm:px-6 lg:px-8 lg:max-w-7xl"
            >
              <h2 class="text-base font-semibold tracking-wider text-yellow-600 uppercase">
                Architecture
              </h2>
              <p class="mt-2 text-3xl font-extrabold text-gray-900 tracking-tight sm:text-4xl">
                How it works
              </p>
              <p class="mt-5 mx-auto max-w-prose text-xl text-gray-500">
                Sturdy sits between your Git repository and contributors in order to provide access
                control for individual files and directories.
              </p>
            </div>
            <div
              class="mt-12 mx-auto max-w-md px-4 grid gap-8 sm:max-w-lg sm:px-6 lg:px-8 lg:grid-cols-3 lg:max-w-7xl"
            >
              <div class="flex flex-col rounded-lg shadow-lg overflow-hidden">
                <div class="flex-shrink-0">
                  <img class="h-48 object-cover" src="./sturdy-role.png" height="420" width="713" />
                </div>
                <div class="flex-1 bg-white p-6 flex flex-col justify-between">
                  <div class="flex-1">
                    <p class="text-sm font-medium text-yellow-600">
                      <a class="hover:underline"> Step One </a>
                    </p>
                    <a class="block mt-2">
                      <p class="text-xl font-semibold text-gray-900">
                        Create a role with suitable access
                      </p>
                      <p class="mt-3 text-base text-gray-500">
                        After installing Sturdy for a Git repository you create a role with the
                        access level that is appropriate for your use case (eg. share only /web or
                        everything except /corp/sensitive). You can then invite collaborators with
                        that role.
                      </p>
                    </a>
                  </div>
                </div>
              </div>
              <div class="flex flex-col rounded-lg shadow-lg overflow-hidden">
                <div class="flex-shrink-0">
                  <img
                    class="h-48 w-full object-cover"
                    src="../../../pages/blog/redesign/newsletter-workspace.png"
                    alt=""
                  />
                </div>
                <div class="flex-1 bg-white p-6 flex flex-col justify-between">
                  <div class="flex-1">
                    <p class="text-sm font-medium text-yellow-600">
                      <a class="hover:underline"> Step Two</a>
                    </p>
                    <a class="block mt-2">
                      <p class="text-xl font-semibold text-gray-900">
                        Contributors develop with the access provided
                      </p>
                      <p class="mt-3 text-base text-gray-500">
                        Sturdy provides a codebase directory on the collaborator's workstation with
                        only the permitted content. Developers can use any tools, IDEs, and text
                        editors they are used to.
                      </p>
                    </a>
                  </div>
                </div>
              </div>
              <div class="flex flex-col rounded-lg shadow-lg overflow-hidden">
                <div class="flex-shrink-0">
                  <img
                    class="h-48 w-full object-cover"
                    src="../../../assets/Web/GitHubPullRequest.png"
                    alt=""
                    height="498"
                    width="937"
                  />
                </div>
                <div class="flex-1 bg-white p-6 flex flex-col justify-between">
                  <div class="flex-1">
                    <p class="text-sm font-medium text-yellow-600">
                      <a class="hover:underline"> Step Three</a>
                    </p>
                    <a class="block mt-2">
                      <p class="text-xl font-semibold text-gray-900">
                        Sturdy generates Pull Requests towards your Git repository
                      </p>
                      <p class="mt-3 text-base text-gray-500">
                        When the collaborator is ready they can generate a Pull Request towards your
                        Git repository through Sturdy. This way contributions can be integrated
                        without providing direct access to your Git repository.
                      </p>
                    </a>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Feature section with grid -->
        <div class="relative bg-white py-16 sm:py-24 lg:py-32">
          <div class="mx-auto max-w-md px-4 text-center sm:max-w-3xl sm:px-6 lg:px-8 lg:max-w-7xl">
            <h2 class="text-base font-semibold tracking-wider text-yellow-600 uppercase">
              Features
            </h2>
            <p class="mt-2 text-3xl font-extrabold text-gray-900 tracking-tight sm:text-4xl">
              Designed for your enterprise needs
            </p>
            <p class="mt-5 max-w-prose mx-auto text-xl text-gray-500">
              Sturdy combines the strengths of Git compatibility with robust access control and
              monitoring.
            </p>
            <div class="mt-12">
              <div class="grid grid-cols-1 gap-8 sm:grid-cols-2 lg:grid-cols-3">
                <div v-for="feature in features" :key="feature.name" class="pt-6">
                  <div class="flow-root bg-gray-50 rounded-lg px-6 pb-8">
                    <div class="-mt-6">
                      <div>
                        <span
                          class="inline-flex items-center justify-center p-3 bg-gradient-to-r from-yellow-500 to-yellow-600 rounded-md shadow-lg"
                        >
                          <component
                            :is="feature.icon"
                            class="h-6 w-6 text-white"
                            aria-hidden="true"
                          />
                        </span>
                      </div>
                      <h3 class="mt-8 text-lg font-medium text-gray-900 tracking-tight">
                        {{ feature.name }}
                      </h3>
                      <p class="mt-5 text-base text-gray-500">
                        {{ feature.description }}
                      </p>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>

        <div class="">
          <div class="max-w-7xl mx-auto py-16 px-4 sm:px-6 lg:px-8">
            <div
              class="bg-gradient-to-r from-yellow-500 to-yellow-600 rounded-lg shadow-xl overflow-hidden lg:grid lg:grid-cols-2 lg:gap-4"
            >
              <div class="pt-10 pb-12 px-6 sm:pt-16 sm:px-16 lg:py-16 lg:pr-0 xl:py-20 xl:px-20">
                <div class="lg:self-center">
                  <h2 class="text-3xl font-extrabold text-black sm:text-4xl">
                    <span class="block">Ready to dive in?</span>
                    <span class="block">Start for free today.</span>
                  </h2>
                  <p class="mt-4 text-lg leading-6 text-black">Try our cloud service</p>
                  <router-link
                    :to="{ name: 'signup' }"
                    class="mt-8 bg-white border border-transparent rounded-md shadow px-5 py-3 inline-flex items-center text-base font-medium text-yellow-600 hover:bg-yellow-50"
                    >Get started in the cloud
                  </router-link>

                  <p class="mt-10 sm:mt-10">Or contact us for enterprise setup.</p>

                  <ACLRequestAccess />
                </div>
              </div>
              <div class="-mt-6 aspect-w-5 aspect-h-3 md:aspect-w-2 md:aspect-h-1">
                <img
                  class="transform translate-x-6 translate-y-6 rounded-md object-cover object-left-top sm:translate-x-16 lg:translate-y-20"
                  src="../../../pages/blog/redesign/newsletter-workspace.png"
                  height="2620"
                  width="4000"
                  alt="Sturdy app screenshot"
                />
              </div>
            </div>
          </div>
        </div>

        <div class="bg-gray-50 py-16 sm:py-24 lg:py-32 mx-auto">
          <span class="relative text-center mx-auto flex justify-center items-center">
            <h2
              class="block bg-clip-text text-transparent bg-gradient-to-r from-yellow-500 to-yellow-600 text-2xl md:text-4xl font-extrabold tracking-tight"
            >
              Backed by
            </h2>
            <img
              class="ml-8 h-12"
              src="../../../assets/Landing/ycombinator-logo.png"
              alt="Y Combinator"
            />
          </span>
        </div>
      </main>
    </div>
  </div>
</template>

<script>
import { Popover, PopoverButton, PopoverPanel } from '@headlessui/vue'
import {
  CloudUploadIcon,
  CogIcon,
  InboxIcon,
  LockClosedIcon,
  MenuIcon,
  RefreshIcon,
  ServerIcon,
  ShieldCheckIcon,
  SparklesIcon,
  XIcon,
} from '@heroicons/vue/outline'
import { ChevronRightIcon, ExternalLinkIcon } from '@heroicons/vue/solid'
import { ClientOnly } from 'vite-ssr/vue'
import ACLRequestAccess from './ACLRequestAccess.vue'

const navigation = [
  { name: 'Product', href: '#' },
  { name: 'Features', href: '#' },
  { name: 'Marketplace', href: '#' },
  { name: 'Company', href: '#' },
]
const features = [
  {
    name: 'Least Privilege Principle',
    description:
      'Protect your Intellectual Property by staying in control over access to sources. Limit the exposure in the event of a developer workstation being compromised.',
    icon: CloudUploadIcon,
  },
  {
    name: 'Easy setup',
    description:
      'No need to split your repository and deal with code or artefact dependencies. No need to configure Git submodules or subtrees. No need to create new CI pipelines. Sturdy plugs into your existing workflow.',
    icon: CogIcon,
  },
  {
    name: 'Code review and collaboration',
    description:
      'Sturdy offers a number of review and collaboration tools for sharing feedback and helping one another.',
    icon: RefreshIcon,
  },
  {
    name: 'Run on your hardware',
    description: 'Sturdy can be installed on your own hardware that you control.',
    icon: ShieldCheckIcon,
  },
  {
    name: 'Identity providers',
    description: 'Integrates with a number of identity providers: OAuth, SAML, Active Directory.',
    icon: LockClosedIcon,
  },
  {
    name: 'Monitor access',
    description:
      'Sturdy provides usage and access logs appropriate for compliance and regulatory needs.',
    icon: ServerIcon,
  },
]

const productName = 'Sturdy'

const usps = [
  {
    title: 'Share specific files and folders from your Git repository',
    subtitle:
      productName +
      ' allows you to share only parts of your Git repository. This way you can get external help with your project without the need to share all sources within a repository. ' +
      'Sturdy acts as a filter between Git and contributors eliminating the need for full clones of the repository on developer workstations. Each local change is securely integrated with the rest of the repository in real time.',
  },
  {
    title: 'Role-based access control',
    subtitle:
      'Setup roles and access control lists that suit your organization. ' +
      productName +
      ' allows you to create inclusion and exclusion filters per role. ' +
      'Collaborators will only see files and folders matching your access control lists.',
  },
  {
    title: 'Plugs into your existing workflow',
    subtitle:
      productName +
      ' integrates with your existing Git workflow.' +
      ' There is no need to split your repository in order to protect the sensitive parts of your codebase. ' +
      productName +
      ' integrates with GitHub and sits between your repository and contributors, enabling partial access to the Git repository.' +
      ' Such granular authorization would normally be impossible with pure Git.',
  },
  {
    title: 'Pull Request contributions without access to your repository',
    subtitle:
      productName +
      ' creates Pull Requests on behalf of your contributors. ' +
      ' This way code changes can be integrated and reviewed without the need to add third party developers to your codebase.' +
      ' At the same time,  collaborators can see test results from your existing Continuous Integration pipelines, no additional setup needed.',
  },
  {
    title: 'Test results from your existing CI can be securely shared',
    subtitle:
      'External developers can securely receive integration feedback from your existing pipelines',
  },
  {
    title: 'USP',
    subtitle: 'usp',
  },
]

export default {
  components: {
    ACLRequestAccess,
    Popover,
    PopoverButton,
    PopoverPanel,
    ChevronRightIcon,
    ExternalLinkIcon,
    MenuIcon,
    XIcon,
    InboxIcon,
    SparklesIcon,
    ClientOnly,
  },
  setup() {
    return {
      productName,
      usps,
      navigation,
      features,
    }
  },
}
</script>
