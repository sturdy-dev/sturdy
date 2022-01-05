<template>
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
        {{ title }}
      </h1>
      <p class="mt-6 text-xl text-cyan-100 max-w-3xl">
        Come join Sturdy as a {{ title }}, and join us on our mission to build great developer
        tooling!
      </p>
    </div>
  </header>
</template>

<script lang="ts">
import { defineComponent } from 'vue'

import headerImage from './united-spaces.jpeg'

export default defineComponent({
  props: {
    title: {
      type: String,
      required: true,
    },
  },
  setup() {
    return {
      headerImage,
    }
  },
})
</script>
