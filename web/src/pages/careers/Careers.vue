<template>
  <Documentation>
    <header class="relative pb-24 bg-sky-300 sm:pb-42">
      <div class="absolute inset-0">
        <img class="w-full h-full object-cover" :src="headerImage" alt="" />
        <div
          class="absolute inset-0 bg-gradient-to-l from-sky-800 to-cyan-700 mix-blend-hard-light"
          aria-hidden="true"
        />
      </div>

      <div
        class="relative max-w-md mx-auto px-4 sm:max-w-3xl sm:py-[12rem] sm:px-6 lg:max-w-7xl lg:px-8"
      >
        <h1 class="text-4xl font-extrabold tracking-tight text-white sm:text-5xl lg:text-6xl">
          Careers at Sturdy
        </h1>
        <p class="mt-6 text-xl text-cyan-100 max-w-3xl">
          Let's build some amazing developer tooling together!
        </p>
      </div>
    </header>

    <div class="max-w-7xl mx-auto py-16 px-4 sm:px-6 lg:py-20 lg:px-8 flex flex-col space-y-32">
      <div class="flex flex-col-reverse lg:flex-row">
        <div class="lg:flex-1">
          <div class="prose">
            <h2 class="text-3xl font-extrabold text-gray-900">About Sturdy</h2>
            <p>At Sturdy we're replacing Git with real-time collaboration in the cloud!</p>

            <p>
              This is a great opportunity to be one of the first employees at a fast moving and
              <strong>open-source</strong> startup, funded by
              <a href="https://www.ycombinator.com/companies/sturdy/">Y Combinator</a>! We’re
              building tech for techies, and naturally really take pride in both the software that
              we’re writing and the product that we’re building, and hope that you do too!
            </p>

            <h2>Sturdy is growing</h2>

            <p>
              We're looking for experienced engineers to join our team, who wants to play a crucial
              role in designing Sturdy, and making us known to the world!
            </p>

            <h2>What we do</h2>

            <p>
              We're building the future of how coders like you and me are collaborating on code.
              With Sturdy developers don't have to fight branches, pull requests, or crazy command
              line arguments to get things done.
            </p>

            <p>
              We regularly post about tech and product on our
              <router-link :to="{ name: 'blog' }">blog</router-link>
              .
            </p>
          </div>
        </div>
        <div class="lg:w-[30rem]">
          <div class="bg-gray-50 p-8 rounded-xl prose">
            <h2>Open Positions</h2>
            <ul>
              <li>
                <router-link :to="{ name: 'careersSeniorBackendSoftwareEngineer' }">
                  Senior Backend Software Engineer
                </router-link>
              </li>
              <li>
                <router-link :to="{ name: 'careersSeniorFrontendSoftwareEngineer' }">
                  Senior Frontend Software Engineer
                </router-link>
              </li>
              <li>
                <router-link :to="{ name: 'careersFullStackEngineer' }">
                  Full Stack Engineer
                </router-link>
              </li>
            </ul>
          </div>
        </div>
      </div>
      <JobTeam />
    </div>
  </Documentation>
</template>

<script lang="ts" setup>
import JobHeader from '../../organisms/careers/JobHeader.vue'
import JobTeam from '../../organisms/careers/JobTeam.vue'
import { useHead } from '@vueuse/head'
import Documentation from '../../layouts/Documentation.vue'
import headerImage from '../../organisms/careers/united-spaces.jpeg'

useHead({
  title: 'Careers at Sturdy',
})
</script>
