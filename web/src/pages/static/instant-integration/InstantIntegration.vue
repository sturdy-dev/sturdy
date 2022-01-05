<template>
  <div class="bg-white">
    <div class="overflow-hidden">
      <div class="wrapper bg-green-200">
        <div class="mx-auto max-w-7xl lg:px-8">
          <div class="lg:grid lg:grid-cols-2 lg:gap-8">
            <div
              class="mx-auto max-w-md px-4 sm:max-w-2xl sm:px-6 sm:text-center lg:px-0 lg:text-left lg:flex lg:items-center"
            >
              <div class="lg:py-24">
                <h1
                  class="text-4xl tracking-tight font-extrabold text-white pt-5 sm:text-6xl xl:text-6xl"
                >
                  <span class="strike">Continuous</span>
                  <span
                    class="pb-3 block bg-clip-text text-transparent bg-gradient-to-r from-yellow-200 to-yellow-400 sm:pb-5"
                  >
                    Instant Integration
                  </span>
                </h1>
                <p class="text-base text-gray-300 sm:text-xl lg:text-lg xl:text-xl">
                  Builds and tests your code as you type it.
                </p>
                <div class="pt-20">
                  <InstantIntegrationRequestAccess />
                  <p class="mt-6 py-6 text-gray-300 sm:mt-6 sm:py-0">
                    Instant Integration is currently in private beta.
                  </p>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- mini usps -->
      <div class="relative bg-gray-50 py-16">
        <div
          class="mx-auto max-w-md px-4 grid gap-8 sm:max-w-lg sm:px-6 lg:px-8 lg:grid-cols-4 lg:max-w-7xl"
        >
          <div
            class="flex flex-col rounded-lg shadow-lg overflow-hidden transition transition-all transform hover:scale-105 hover:shadow-xl"
          >
            <div class="flex-1 bg-white p-6 flex flex-col justify-between">
              <div class="flex-shrink-0 flex items-center justify-center p-3">
                <FireIcon class="h-14 w-14 text-yellow-600" aria-hidden="true" />
              </div>
              <div class="flex-1">
                <a class="block mt-2">
                  <p class="text-xl font-semibold text-gray-900">Instant test feedback</p>
                  <p class="mt-3 text-base text-gray-500">
                    As you code, your tests are run on re-usable and hot test runners. You get
                    results from the most relevant tests first, giving you a tighter feedback loop.
                    Runners are kept warm while you develop, giving you up to 90% less waiting time.
                  </p>
                </a>
              </div>
            </div>
          </div>

          <div
            class="flex flex-col rounded-lg shadow-lg overflow-hidden transition transition-all transform hover:scale-105 hover:shadow-xl"
          >
            <div class="flex-1 bg-white p-6 flex flex-col justify-between">
              <div class="flex-shrink-0 flex items-center justify-center p-3">
                <HomeIcon class="h-14 w-14 text-yellow-600" aria-hidden="true" />
              </div>
              <div class="flex-1">
                <a class="block mt-2">
                  <p class="text-xl font-semibold text-gray-900">On your existing runners</p>
                  <p class="mt-3 text-base text-gray-500">
                    II runs on your existing CI/CD provider, like CircleCI or GitHub Actions.
                    Compatible with Docker and Docker Compose. The setup is simply wrapping your
                    existing commands with ii.
                  </p>
                  <p class="font-mono mt-3 text-base text-gray-500">$ ii npm test</p>
                </a>
              </div>
            </div>
          </div>

          <div
            class="flex flex-col rounded-lg shadow-lg overflow-hidden transition transition-all transform hover:scale-105 hover:shadow-xl"
          >
            <div class="flex-1 bg-white p-6 flex flex-col justify-between">
              <div class="flex-shrink-0 flex items-center justify-center p-3">
                <VariableIcon class="h-14 w-14 text-yellow-600" aria-hidden="true" />
              </div>
              <div class="flex-1">
                <a class="block mt-2">
                  <p class="text-xl font-semibold text-gray-900">Intelligent scheduling</p>
                  <p class="mt-3 text-base text-gray-500">
                    Sturdy triggers builds intelligently, and it optimizes resource utilization and
                    how often your tests run. Tests are ran progressively, starting with the most
                    relevant ones first, based on historical data and other heuristics.
                  </p>
                </a>
              </div>
            </div>
          </div>

          <div
            class="flex flex-col rounded-lg shadow-lg overflow-hidden transition transition-all transform hover:scale-105 hover:shadow-xl"
          >
            <div class="flex-1 bg-white p-6 flex flex-col justify-between">
              <div class="flex-shrink-0 flex items-center justify-center p-3">
                <CollectionIcon class="h-14 w-14 text-yellow-600" aria-hidden="true" />
              </div>
              <div class="flex-1">
                <a class="block mt-2">
                  <p class="text-xl font-semibold text-gray-900">It's all cached</p>
                  <p class="mt-3 text-base text-gray-500">
                    Incremental testing on Sturdy, with hot-caches for all of your dependencies
                    makes sure that your tests can run really fast.
                  </p>
                </a>
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
            class="lg:mx-auto lg:max-w-7xl lg:px-2 lg:grid lg:grid-flow-col-dense lg:gap-12 lg:grid-cols-3"
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
                    Tests your code as you write it
                  </h2>
                  <p class="mt-4 text-lg text-gray-500">
                    Designed to create a faster development feedback loop, Sturdy’s Instant
                    Integration builds and runs your tests progressively as you develop. This means
                    that you don't have to run tests locally, and can simply glance at the latest II
                    status in your workspace.
                  </p>
                  <p class="mt-4 text-lg text-gray-500">
                    Thanks to our hot-caches, the feedback is not only continuous but also
                    <strong>instant</strong>.
                  </p>
                </div>
              </div>
            </div>

            <div class="col-span-2 mt-12 sm:mt-16 lg:mt-0 mx-4 md:mx-0">
              <div class="sm:px-6 lg:py-16 lg:px-0 lg:m-0 lg:relative lg:h-full">
                <video
                  autoplay
                  loop
                  muted
                  playsinline
                  alt="Instant Integration builds and tests your code as it's written"
                  class="rounded-lg object-cover object-center"
                >
                  <source src="./ii.webm" type="video/webm" />
                  <source src="./ii.mp4" type="video/mp4" />
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
                    Dead simple setup
                  </h2>
                  <p class="mt-4 text-lg text-gray-500">
                    Firstly, Instant Integration runs on existing CI providers like CircleCI,
                    Buildkite and Github Actions. It slots right into your
                    <strong>build</strong> or <strong>test</strong> steps, for example:
                  </p>
                  <p class="mt-4 text-lg text-gray-500 font-mono">$ ii npm test</p>
                  <p class="mt-4 text-lg text-gray-500">
                    At the same time, some of your existing tests might just be slow or flaky.
                    <strong>II</strong> keeps track of test performance and highlights where you
                    could get the biggest benefit from making improvements.
                    <strong>The faster a test is, the earlier it gets executed</strong> - the
                    fastest ones run and have the results ready every time a file is saved. This
                    creates a direct incentive for investing in making tests modular.
                  </p>
                </div>
              </div>
            </div>
            <div class="mt-12 sm:mt-16 lg:mt-20 mx-3">
              <div class="sm:px-6 lg:py-16 lg:px-0 lg:m-0 lg:relative lg:h-full mx-auto max-w-lg">
                <div></div>
                <img
                  class="rounded-lg object-cover object-center w-full"
                  src="./cicd-providers.png"
                  alt="Sturdy's Instant Integration provides granular results"
                />
              </div>
            </div>
          </div>
        </div>
      </div>
      <!-- How it works section -->
      <div class="relative bg-gray-50 py-16 sm:py-24 lg:py-32">
        <div class="relative">
          <div class="mx-auto max-w-md px-4 sm:max-w-3xl sm:px-6 lg:px-8 lg:max-w-7xl">
            <h2
              class="text-center text-base font-semibold tracking-wider text-yellow-600 uppercase"
            >
              Architecture
            </h2>
            <p
              class="text-center mt-2 text-3xl font-extrabold text-gray-900 tracking-tight sm:text-4xl"
            >
              How it works
            </p>
            <p class="text-left md:text-center mt-5 mx-auto max-w-prose text-xl text-gray-500">
              Instant Integration (or or just <strong>II</strong>)
              <a
                class="text-yellow-500 hover:bg-gray-50"
                href="https://en.wikipedia.org/wiki/Shift-left_testing"
                >shifts left</a
              >
              and gives you instant test feedback. It works by continuously streaming code changes
              to warmed-up runners with hot-caches and intelligently scheduling the execution.
            </p>

            <div class="mt-16 max-w-xl mx-auto px-4 sm:px-6 lg:max-w-7xl lg:px-8">
              <dl class="space-y-10 lg:space-y-0 lg:grid lg:grid-cols-3 lg:gap-8">
                <div>
                  <dt>
                    <p class="text-center mt-5 text-lg leading-6 font-medium text-gray-900">
                      Real-time code sync
                    </p>
                  </dt>
                  <dd class="mt-2 text-base text-gray-500">
                    Instant Integration builds on top of Sturdy
                    <router-link
                      :to="{ name: 'featuresLiveFeedback' }"
                      class="text-yellow-500 hover:bg-gray-50"
                    >
                      workspaces </router-link
                    >. Changes you make locally are available in real-time for your team, and now
                    also for your test runners. Only the code deltas are transmitted, which means
                    there is virtually no overhead from copying sources.
                  </dd>
                </div>
                <div>
                  <dt>
                    <p class="text-center mt-5 text-lg leading-6 font-medium text-gray-900">
                      Caching
                    </p>
                  </dt>
                  <dd class="mt-2 text-base text-gray-500">
                    Sturdy provides hot caches of all your dependencies (yes, your databases as
                    well!). While you are actively working, a runner will be kept warm and ready to
                    execute tests.
                  </dd>
                </div>
                <div>
                  <dt>
                    <p class="text-center mt-5 text-lg leading-6 font-medium text-gray-900">
                      Smart Scheduling
                    </p>
                  </dt>
                  <dd class="mt-2 text-base text-gray-500">
                    By incrementally running your tests, you'll get the test results really fast. II
                    optimizes resource usage by intelligently detecting test relevance in relation
                    to code changes. Over time II learns failure patterns, further improving the
                    scheduling.
                  </dd>
                </div>
              </dl>
            </div>

            <div class="mt-12 sm:mt-16 lg:mt-0">
              <div class="sm:px-6 lg:py-16 lg:px-0 lg:relative lg:h-full max-w-3xl mx-auto">
                <img
                  class="rounded-lg object-cover object-center w-full"
                  src="./concurrency.png"
                  alt="Technical diagram illustrating how Instant Integration works"
                />
              </div>
            </div>
          </div>
        </div>
      </div>

      <div class="bg-gray-800">
        <div class="max-w-7xl mx-auto py-12 px-4 sm:px-6 lg:py-16 lg:px-8 lg:flex lg:items-center">
          <div class="lg:w-0 lg:flex-1">
            <h2
              id="newsletter-headline"
              class="text-3xl font-extrabold tracking-tight text-white sm:text-4xl"
            >
              Join the waitlist
            </h2>
            <p class="mt-3 max-w-3xl text-lg leading-6 text-gray-300">
              Instant integration is currently in private beta. If you’d like to be notified when it
              launches, and have a chance to get early access, put your email below.
            </p>
          </div>
          <div class="mt-8 lg:mt-0 lg:ml-8">
            <InstantIntegrationRequestAccess />

            <p class="mt-3 text-sm text-gray-300">
              We care about the protection of your data. Read our
              {{ ' ' }}
              <router-link :to="{ name: 'privacy' }" class="text-white font-medium underline">
                Privacy Policy
              </router-link>
            </p>
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
    </div>
  </div>
</template>

<script>
import { InboxIcon, SparklesIcon } from '@heroicons/vue/outline'
import { FireIcon, VariableIcon, HomeIcon, CollectionIcon } from '@heroicons/vue/solid'
import InstantIntegrationRequestAccess from './InstantIntegrationRequestAccess.vue'
export default {
  name: 'InstantIntegration',
  components: {
    InboxIcon,
    SparklesIcon,
    FireIcon,
    VariableIcon,
    HomeIcon,
    CollectionIcon,
    InstantIntegrationRequestAccess,
  },
}
</script>

<style scoped>
.wrapper {
  background-image: url('./interstellar.svg');
  background-size: cover;
  background-position: 50% 50%;
}

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

.strike::after {
  content: ' ';
  position: absolute;
  top: 50%;
  left: 0;
  height: 6px;
  background: #d97706;
  animation-name: strike;
  animation-duration: 0.5s;
  animation-timing-function: linear;
  animation-iteration-count: 1;
  animation-fill-mode: forwards;
  animation-delay: 0.5s;
}
</style>
