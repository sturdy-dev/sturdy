<template>
  <Public>
    <JobHeader title="Come and work with us!" />

    <div class="max-w-7xl mx-auto py-16 px-4 sm:px-6 lg:py-20 lg:px-8 flex flex-col space-y-32">
      <div class="flex flex-col-reverse lg:flex-row">
        <div class="lg:flex-1">
          <div class="prose">
            <h2 class="text-3xl font-extrabold text-gray-900">About Sturdy</h2>
            <p>At Sturdy we're replacing Git with real-time collaboration in the cloud!</p>

            <p>
              This is a great opportunity to be one of the first employees at a fast moving startup,
              funded by <a href="https://www.ycombinator.com/companies/sturdy/">Y Combinator</a>!
              We’re building tech for techies, and naturally really take pride in both the software
              that we’re writing and the product that we’re building, and hope that you do too!
            </p>

            <h2>Sturdy is growing</h2>

            <p>
              We're looking for Backend Engineers with Go experience to join our team, together with
              the team, you'll play a crucial role in designing Sturdy, and making us known to the
              world.
            </p>

            <h2>What we do</h2>

            <p>
              We're building the future of how coders like you and me are collaborating on code.
              With Sturdy developers don't have to fight branches, pull requests, or crazy command
              line arguments to get things done.
            </p>

            <p>
              We regularly post about tech and product on our
              <router-link :to="{ name: 'blog' }">blog</router-link>.
            </p>
          </div>
        </div>
        <div class="lg:w-80">
          <div class="bg-gray-50 p-8 rounded-xl prose">
            <h2>Open Positions</h2>
            <ul>
              <li>
                <router-link :to="{ name: 'careersFoundingBackendEngineer' }">
                  Founding Backend Engineer
                </router-link>
              </li>
              <li>
                <router-link :to="{ name: 'careersFoundingFrontendEngineer' }">
                  Founding Frontend Engineer
                </router-link>
              </li>
            </ul>
          </div>
        </div>
      </div>
      <JobTeam />
    </div>
  </Public>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import Public from '../../layouts/Public.vue'
import JobHeader from '../../organisms/careers/JobHeader.vue'
import JobTeam from '../../organisms/careers/JobTeam.vue'
import { useHead } from '@vueuse/head'

export default defineComponent({
  components: {
    Public,
    JobTeam,
    JobHeader,
  },
  setup() {
    useHead({
      title: 'Careers at Sturdy',
    })
  },
})
</script>
