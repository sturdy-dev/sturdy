<template>
  <div class="antialiased text-slate-300 bg-slate-900">
    <div
      class="blurr border-b border-slate-700/10 sticky top-0 z-50 py-2 lg:py-3 px-6 lg:px-20 flex items-center justify-between leading-6"
    >
      <img
        src="./logotype.svg"
        alt="Sturdy logotype"
        class="rounded-xl object-contain"
        width="130"
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

    <div class="relative -mt-6 lg:-mt-20">
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
    </div>

    <div class="pt-10 xl:pt-0 mt-10 xl:mt-2 space-y-32 sm:space-y-40 md:space-y-60">
      <div class="pt-16 max-w-7xl px-4 sm:px-6 md:px-8 mx-auto lg:gap-8 space-y-6">
        <div
          class="mt-4 text-lg sm:text-xl font-extrabold tracking-tight text-slate-300/70 inline-flex space-x-4 items-center"
        >
          <div>Backed by</div>
          <div class="flex space-x-10 items-center">
            <a href="https://www.ycombinator.com/companies/sturdy">
              <img
                src="../../../assets/Landing/ycombinator-logo.png"
                alt="Y Combinator logo"
                class="object-contain h-6 grayscale"
              />
            </a>
          </div>
        </div>
      </div>
    </div>

    <section id="features">
      <div class="pt-10 xl:pt-0 mt-10 xl:mt-2 space-y-32 sm:space-y-40 md:space-y-60">
        <div class="pt-16 max-w-7xl px-4 sm:px-6 md:px-8 mx-auto lg:gap-8 space-y-6">
          <h1
            class="text-center mt-4 text-3xl sm:text-4xl font-extrabold tracking-tight text-slate-50"
          >
            A new home for your code
          </h1>
          <div class="pt-10 px-4 flex flex-col sm:flex-row sm:space-x-20 space-y-10 md:space-y-0">
            <div class="md:w-1/2 flex flex-col md:flex-row md:space-x-20 space-y-10 md:space-y-0">
              <MiniUsp class="md:w-1/2" title="Expressive">
                <template #icon>
                  <IconDevelop />
                </template>
                <p>
                  Ship the code to production &mdash; <br />
                  instead of the usual dance of creating a branch, staging and committing files,
                  pushing, and merging.
                </p>
              </MiniUsp>

              <MiniUsp class="md:w-1/2" title="Collaborative">
                <template #icon>
                  <IconGears />
                </template>
                <div>
                  Try each others code and give actionable feedback &mdash; <br />
                  instead of guessing what a pull request does and LGTM-ing anyway.
                </div>
              </MiniUsp>
            </div>

            <div class="md:w-1/2 flex flex-col md:flex-row md:space-x-20 space-y-10 md:space-y-0">
              <MiniUsp class="md:w-1/2" title="Compatible">
                <template #icon>
                  <IconMobile />
                </template>
                <div>
                  Use Sturdy together with GitHub &mdash; or on it's own. By yourself &mdash; or
                  with others on your team. Run your own instance &mdash; or use Sturdy in the
                  cloud.
                </div>
              </MiniUsp>

              <MiniUsp class="md:w-1/2" title="Streamlined">
                <template #icon>
                  <IconDesktop />
                </template>
                <div>
                  Ship small, incremental changes as the default &mdash;<br />
                  instead of reminding yourself to "push often" and "keep the PR small".
                </div>
              </MiniUsp>
            </div>
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
      >
        <template #icon>
          <IconGears></IconGears>
        </template>
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

      <Usp
        title="Streamlined"
        subtitle="For the workflow that you're already using."
        link="v2DocsHotToShipSoftwareToProduction"
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
          <div class="pt-16 max-w-7xl px-4 sm:px-6 md:px-8 mx-auto lg:gap-8 space-y-6">
            <h1 class="mt-4 text-lg sm:text-xl font-extrabold tracking-tight text-slate-300">
              Featured at
            </h1>
            <div class="flex space-x-10 items-center">
              <a href="https://websummit.com/schedule/ws21/timeslot/breakout-startups-1104-0927">
                <img
                  src="./web-summit.png"
                  alt="Web Summit logotype"
                  class="rounded-xl object-contain h-6"
                />
              </a>
              <a href="https://devsecops.fm/episodes/36-sturdy/">
                <img
                  src="./devsecops.png"
                  alt="DevSecOps logotype"
                  class="rounded-xl object-contain h-8"
                />
              </a>
              <img src="./devcon.png" alt="Devcon logotype" class="rounded-xl object-contain h-8" />
            </div>
          </div>
        </div>
      </section>
    </div>

    <div id="footer" class="pb-16 text-sm leading-6">
      <div class="max-w-7xl mx-auto divide-y px-4 sm:px-6 md:px-8 divide-slate-700">
        <div class="flex flex-row">
          <div class="w-1/3 space-y-8">
            <!--          class="blurr border-b border-slate-700/10 sticky top-0 z-50 py-2 lg:py-3 px-6 lg:px-20 flex items-center justify-between leading-6"-->
            <img
              src="./logotype.svg"
              alt="Sturdy logotype"
              class="rounded-xl object-contain"
              width="130"
            />
            <div class="mt-4 flex items-center space-x-6">
              <a class="text-slate-300 hover:text-slate-100" href="https://twitter.com/getsturdy">
                <span class="sr-only">Twitter</span>
                <svg class="h-6 w-6" fill="currentColor" viewBox="0 0 24 24" aria-hidden="true">
                  <path
                    d="M8.29 20.251c7.547 0 11.675-6.253 11.675-11.675 0-.178 0-.355-.012-.53A8.348 8.348 0 0022 5.92a8.19 8.19 0 01-2.357.646 4.118 4.118 0 001.804-2.27 8.224 8.224 0 01-2.605.996 4.107 4.107 0 00-6.993 3.743 11.65 11.65 0 01-8.457-4.287 4.106 4.106 0 001.27 5.477A4.072 4.072 0 012.8 9.713v.052a4.105 4.105 0 003.292 4.022 4.095 4.095 0 01-1.853.07 4.108 4.108 0 003.834 2.85A8.233 8.233 0 012 18.407a11.616 11.616 0 006.29 1.84"
                  ></path>
                </svg>
              </a>
              <a class="text-slate-300 hover:text-slate-100" href="https://discord.gg/fQcH9QAVpX">
                <span class="sr-only">Discord</span>
                <svg class="h-6 w-6" fill="currentColor" viewBox="0 0 71 55" aria-hidden="true">
                  <path
                    fill-rule="evenodd"
                    d="M60.1045 4.8978C55.5792 2.8214 50.7265 1.2916 45.6527 0.41542C45.5603 0.39851 45.468 0.440769 45.4204 0.525289C44.7963 1.6353 44.105 3.0834 43.6209 4.2216C38.1637 3.4046 32.7345 3.4046 27.3892 4.2216C26.905 3.0581 26.1886 1.6353 25.5617 0.525289C25.5141 0.443589 25.4218 0.40133 25.3294 0.41542C20.2584 1.2888 15.4057 2.8186 10.8776 4.8978C10.8384 4.9147 10.8048 4.9429 10.7825 4.9795C1.57795 18.7309 -0.943561 32.1443 0.293408 45.3914C0.299005 45.4562 0.335386 45.5182 0.385761 45.5576C6.45866 50.0174 12.3413 52.7249 18.1147 54.5195C18.2071 54.5477 18.305 54.5139 18.3638 54.4378C19.7295 52.5728 20.9469 50.6063 21.9907 48.5383C22.0523 48.4172 21.9935 48.2735 21.8676 48.2256C19.9366 47.4931 18.0979 46.6 16.3292 45.5858C16.1893 45.5041 16.1781 45.304 16.3068 45.2082C16.679 44.9293 17.0513 44.6391 17.4067 44.3461C17.471 44.2926 17.5606 44.2813 17.6362 44.3151C29.2558 49.6202 41.8354 49.6202 53.3179 44.3151C53.3935 44.2785 53.4831 44.2898 53.5502 44.3433C53.9057 44.6363 54.2779 44.9293 54.6529 45.2082C54.7816 45.304 54.7732 45.5041 54.6333 45.5858C52.8646 46.6197 51.0259 47.4931 49.0921 48.2228C48.9662 48.2707 48.9102 48.4172 48.9718 48.5383C50.038 50.6034 51.2554 52.5699 52.5959 54.435C52.6519 54.5139 52.7526 54.5477 52.845 54.5195C58.6464 52.7249 64.529 50.0174 70.6019 45.5576C70.6551 45.5182 70.6887 45.459 70.6943 45.3942C72.1747 30.0791 68.2147 16.7757 60.1968 4.9823C60.1772 4.9429 60.1437 4.9147 60.1045 4.8978ZM23.7259 37.3253C20.2276 37.3253 17.3451 34.1136 17.3451 30.1693C17.3451 26.225 20.1717 23.0133 23.7259 23.0133C27.308 23.0133 30.1626 26.2532 30.1066 30.1693C30.1066 34.1136 27.28 37.3253 23.7259 37.3253ZM47.3178 37.3253C43.8196 37.3253 40.9371 34.1136 40.9371 30.1693C40.9371 26.225 43.7636 23.0133 47.3178 23.0133C50.9 23.0133 53.7545 26.2532 53.6986 30.1693C53.6986 34.1136 50.9 37.3253 47.3178 37.3253Z"
                    clip-rule="evenodd"
                  ></path>
                </svg>
              </a>
              <a class="text-slate-300 hover:text-slate-100" href="https://github.com/sturdy-dev">
                <span class="sr-only">GitHub</span>
                <svg
                  class="h-6 h-6"
                  fill="currentColor"
                  viewBox="0 0 30 30"
                  width="30px"
                  height="30px"
                >
                  <path
                    d="M15,3C8.373,3,3,8.373,3,15c0,5.623,3.872,10.328,9.092,11.63C12.036,26.468,12,26.28,12,26.047v-2.051 c-0.487,0-1.303,0-1.508,0c-0.821,0-1.551-0.353-1.905-1.009c-0.393-0.729-0.461-1.844-1.435-2.526 c-0.289-0.227-0.069-0.486,0.264-0.451c0.615,0.174,1.125,0.596,1.605,1.222c0.478,0.627,0.703,0.769,1.596,0.769 c0.433,0,1.081-0.025,1.691-0.121c0.328-0.833,0.895-1.6,1.588-1.962c-3.996-0.411-5.903-2.399-5.903-5.098 c0-1.162,0.495-2.286,1.336-3.233C9.053,10.647,8.706,8.73,9.435,8c1.798,0,2.885,1.166,3.146,1.481C13.477,9.174,14.461,9,15.495,9 c1.036,0,2.024,0.174,2.922,0.483C18.675,9.17,19.763,8,21.565,8c0.732,0.731,0.381,2.656,0.102,3.594 c0.836,0.945,1.328,2.066,1.328,3.226c0,2.697-1.904,4.684-5.894,5.097C18.199,20.49,19,22.1,19,23.313v2.734 c0,0.104-0.023,0.179-0.035,0.268C23.641,24.676,27,20.236,27,15C27,8.373,21.627,3,15,3z"
                  />
                </svg>
              </a>
            </div>
          </div>
          <!-- columns -->
          <div class="flex w-2/3">
            <!-- left -->
            <div class="flex-none w-1/2 space-y-10 sm:space-y-8 lg:flex lg:space-y-0">
              <!-- One -->
              <div class="lg:flex-none lg:w-1/2">
                <h2 class="font-semibold text-slate-100">Product</h2>
                <ul class="mt-3 space-y-2">
                  <li>
                    <router-link :to="{ name: 'pricing' }" class="hover:text-slate-100">
                      Pricing
                    </router-link>
                  </li>
                  <li>
                    <router-link :to="{ name: 'blog' }" class="hover:text-slate-100">
                      Blog
                    </router-link>
                  </li>
                  <li>
                    <router-link :to="{ name: 'blog' }" class="hover:text-slate-100">
                      <span class="inline-flex items-center">
                        Status <ExternalLinkIcon class="w-4 h-4 ml-1.5" />
                      </span>
                    </router-link>
                  </li>
                </ul>
              </div>
              <!-- Two -->
              <div class="lg:flex-none lg:w-1/2">
                <h2 class="font-semibold text-slate-100">Company</h2>
                <ul class="mt-3 space-y-2">
                  <li>
                    <router-link :to="{ name: 'about' }" class="hover:text-slate-100">
                      About
                    </router-link>
                  </li>
                  <li>
                    <router-link :to="{ name: '' }" class="hover:text-slate-100">
                      Values
                    </router-link>
                  </li>
                  <li>
                    <router-link :to="{ name: 'careers' }" class="hover:text-slate-100">
                      Careers
                    </router-link>
                  </li>
                </ul>
              </div>
            </div>
            <!-- right -->
            <div class="flex-none w-1/2 space-y-10 sm:space-y-8 lg:flex lg:space-y-0">
              <!-- Three -->
              <div class="lg:flex-none lg:w-1/2">
                <h2 class="font-semibold text-slate-100">Developers</h2>
                <ul class="mt-3 space-y-2">
                  <li>
                    <router-link :to="{ name: 'v2DocsRoot' }" class="hover:text-slate-100">
                      Docs
                    </router-link>
                  </li>
                  <li>
                    <router-link :to="{ name: 'resourcesApi' }" class="hover:text-slate-100">
                      API
                    </router-link>
                  </li>
                  <li>
                    <router-link :to="{ name: '' }" class="hover:text-slate-100">
                      Sources
                    </router-link>
                  </li>
                </ul>
              </div>
              <!-- Four -->
              <div class="lg:flex-none lg:w-1/2">
                <h2 class="font-semibold text-slate-100">Legal</h2>
                <ul class="mt-3 space-y-2">
                  <li>
                    <router-link :to="{ name: 'termsOfService' }" class="hover:text-slate-100">
                      Terms of Service
                    </router-link>
                  </li>
                  <li>
                    <router-link :to="{ name: 'privacy' }" class="hover:text-slate-100">
                      Privacy
                    </router-link>
                  </li>
                </ul>
              </div>
            </div>
          </div>
        </div>
        <div class="mt-24 border-t">
          <p class="mt-4 text-slate-300/50">&copy; {{ year }} Sturdy Sweden AB</p>
        </div>
      </div>
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
import { defineComponent } from 'vue'
import { useHead } from '@vueuse/head'
import {
  DownloadIcon,
  ChevronRightIcon,
  PlayIcon,
  LightBulbIcon,
  DuplicateIcon,
  ExternalLinkIcon,
} from '@heroicons/vue/outline'
import time from '../../../time'

export default defineComponent({
  components: {
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
    ExternalLinkIcon,
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
