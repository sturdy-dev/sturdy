<template>
  <div class="space-y-12">
    <h2 class="text-3xl font-bold sm:text-4xl">About Sturdy</h2>
    <div class="max-w-2xl">
      <p class="text-xl text-gray-500 pb-4">Sturdy is built by developers for developers.</p>
      <p class="text-gray-500">
        Join our
        <a class="text-gray-800" href="https://discord.gg/5HnSdzMqtA">Community Discord</a>
        for discussions about Sturdy, and to connect with the developers! If you have a more private
        question, reach out to us at
        <a class="text-gray-800" href="mailto:support@getsturdy.com">support@getsturdy.com</a>
      </p>
    </div>
    <ul
      class="space-y-12 lg:grid lg:grid-cols-2 lg:items-start lg:gap-x-8 lg:gap-y-12 lg:space-y-0"
    >
      <li v-for="employee in employees" :key="employee.name">
        <div class="space-y-40 sm:grid sm:grid-cols-3 sm:gap-6 sm:space-y-0 lg:gap-8">
          <div class="h-0">
            <img
              class="w-32 sm:w-48 shadow-lg rounded-lg"
              :src="employee.image"
              :alt="employee.name"
            />
          </div>
          <div class="sm:col-span-2">
            <div class="space-y-4">
              <div class="text-lg leading-6 font-medium space-y-1">
                <h3>{{ employee.name }}</h3>
                <p class="text-yellow-600">{{ employee.title }}</p>
              </div>
              <div class="text-lg">
                <p class="text-gray-500">
                  {{ employee.description }}
                </p>
              </div>
              <div class="text-s">
                <a class="text-gray-400" :href="'mailto:' + employee.email">{{ employee.email }}</a>
              </div>
              <ul class="flex space-x-5 items-center">
                <li>
                  <a
                    :href="'https://twitter.com/' + employee.twitter"
                    class="text-gray-400 hover:text-gray-500"
                  >
                    <span class="sr-only">Twitter</span>
                    <svg class="w-5 h-5" fill="currentColor" viewBox="0 0 20 20" aria-hidden="true">
                      <path
                        d="M6.29 18.251c7.547 0 11.675-6.253 11.675-11.675 0-.178 0-.355-.012-.53A8.348 8.348 0 0020 3.92a8.19 8.19 0 01-2.357.646 4.118 4.118 0 001.804-2.27 8.224 8.224 0 01-2.605.996 4.107 4.107 0 00-6.993 3.743 11.65 11.65 0 01-8.457-4.287 4.106 4.106 0 001.27 5.477A4.073 4.073 0 01.8 7.713v.052a4.105 4.105 0 003.292 4.022 4.095 4.095 0 01-1.853.07 4.108 4.108 0 003.834 2.85A8.233 8.233 0 010 16.407a11.616 11.616 0 006.29 1.84"
                      />
                    </svg>
                  </a>
                </li>
                <li>
                  <a
                    :href="'https://github.com/' + employee.github"
                    class="text-gray-400 hover:text-gray-500"
                  >
                    <span class="sr-only">GitHub</span>
                    <svg class="w-5 h-5" fill="currentColor" viewBox="0 0 16 16" aria-hidden="true">
                      <path
                        fill-rule="evenodd"
                        d="M8 0C3.58 0 0 3.58 0 8c0 3.54 2.29 6.53 5.47 7.59.4.07.55-.17.55-.38 0-.19-.01-.82-.01-1.49-2.01.37-2.53-.49-2.69-.94-.09-.23-.48-.94-.82-1.13-.28-.15-.68-.52-.01-.53.63-.01 1.08.58 1.23.82.72 1.21 1.87.87 2.33.66.07-.52.28-.87.51-1.07-1.78-.2-3.64-.89-3.64-3.95 0-.87.31-1.59.82-2.15-.08-.2-.36-1.02.08-2.12 0 0 .67-.21 2.2.82.64-.18 1.32-.27 2-.27.68 0 1.36.09 2 .27 1.53-1.04 2.2-.82 2.2-.82.44 1.1.16 1.92.08 2.12.51.56.82 1.27.82 2.15 0 3.07-1.87 3.75-3.65 3.95.29.25.54.73.54 1.48 0 1.07-.01 1.93-.01 2.2 0 .21.15.46.55.38A8.013 8.013 0 0 0 16 8c0-4.42-3.58-8-8-8z"
                        clip-rule="evenodd"
                      />
                    </svg>
                  </a>
                </li>
                <li>
                  <a
                    :href="'https://www.linkedin.com/in/' + employee.linkedin"
                    class="text-gray-400 hover:text-gray-500"
                  >
                    <span class="sr-only">LinkedIn</span>
                    <svg class="w-5 h-5" fill="currentColor" viewBox="0 0 20 20" aria-hidden="true">
                      <path
                        fill-rule="evenodd"
                        d="M16.338 16.338H13.67V12.16c0-.995-.017-2.277-1.387-2.277-1.39 0-1.601 1.086-1.601 2.207v4.248H8.014v-8.59h2.559v1.174h.037c.356-.675 1.227-1.387 2.526-1.387 2.703 0 3.203 1.778 3.203 4.092v4.711zM5.005 6.575a1.548 1.548 0 11-.003-3.096 1.548 1.548 0 01.003 3.096zm-1.337 9.763H6.34v-8.59H3.667v8.59zM17.668 1H2.328C1.595 1 1 1.581 1 2.298v15.403C1 18.418 1.595 19 2.328 19h15.34c.734 0 1.332-.582 1.332-1.299V2.298C19 1.581 18.402 1 17.668 1z"
                        clip-rule="evenodd"
                      />
                    </svg>
                  </a>
                </li>
              </ul>
            </div>
          </div>
        </div>
      </li>
    </ul>
  </div>
</template>

<script lang="ts" setup>
const employees = [
  {
    name: 'Kiril Videlov',
    title: 'Founder & CEO',
    image: 'https://getsturdy.com/assets/kiril.a50ed74d.jpeg',
    description:
      'Kiril is a self proclaimed Software Enthusiast who loves computers, developer tools & automation. He was previously a Staff Software Engineer at Tink.',
    twitter: 'krlvi',
    email: 'kiril@getsturdy.com',
    github: 'krlvi',
    linkedin: 'krlvi',
  },
  {
    name: 'Gustav Westling',
    title: 'Founder & CTO',
    image: 'https://getsturdy.com/assets/gustav.e261d7c0.jpeg',
    description:
      'Gustav has worked at startups since he was 15, and was previously a Staff Software Engineer at Tink where he specialized in developer productivity tools.',
    twitter: 'zegl',
    email: 'gustav@getsturdy.com',
    github: 'zegl',
    linkedin: 'gustavwestling',
  },
]
</script>
