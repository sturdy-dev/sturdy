<!-- This example requires Tailwind CSS v2.0+ -->
<template>
  <div>
    <div class="space-y-5 sm:space-y-4">
      <h2 class="text-3xl font-extrabold sm:text-4xl">Meet our team</h2>
      <ul role="list" class="space-y-12 sm:grid sm:grid-cols-2 sm:gap-12 sm:space-y-0 lg:gap-x-8">
        <li v-for="person in people" :key="person.name">
          <div class="flex items-center space-x-4 lg:space-x-6">
            <img
              v-if="person.imageUrl"
              class="w-16 h-16 rounded-full lg:w-20 lg:h-20"
              :src="person.imageUrl"
              :alt="person.name + ' ' + person.role + ' at Sturdy'"
            />
            <div v-else class="w-16 h-16 rounded-full lg:w-20 lg:h-20 bg-gray-200"></div>
            <div class="font-medium text-lg leading-6 space-y-1">
              <h3>{{ person.name }}</h3>
              <p class="text-yellow-600">{{ person.role }}</p>
              <p v-if="person.twitter" class="text-sm">
                <a :href="'https://twitter.com/' + person.twitter">@{{ person.twitter }}</a>
              </p>
            </div>
          </div>
        </li>
      </ul>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'

import kiril from './kiril.jpeg?w=500&h=500&fit=cover'
import gustav from './gustav.jpeg?w=500&h=500&fit=cover'
import nikita from './nikita.jpeg?w=500&h=500&fit=cover'

const people = [
  {
    name: 'Kiril Videlov',
    role: 'Co-Founder / CEO',
    imageUrl: kiril,
    twitter: 'krlvi',
  },
  {
    name: 'Gustav Westling',
    role: 'Co-Founder / CTO',
    imageUrl: gustav,
    twitter: 'zegl',
  },
  {
    name: 'Nikita Galaiko',
    role: 'Founding Software Engineer',
    imageUrl: nikita,
  },
  {
    name: 'You?',
    role: "We're expanding, apply now!",
  },
  // More people...
]

export default defineComponent({
  setup() {
    return {
      people,
    }
  },
})
</script>
