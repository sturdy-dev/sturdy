<template>
  <div class="">
    <div class="max-w-7xl mx-auto py-16 px-4 sm:px-6 lg:px-8">
      <div
        class="bg-yellow-400 rounded-lg shadow-xl overflow-hidden lg:grid lg:grid-cols-2 lg:gap-4"
      >
        <div class="pt-10 pb-12 px-6 sm:pt-16 sm:px-16 lg:py-16 lg:pr-0 xl:py-20 xl:px-20">
          <div class="lg:self-center">
            <h2 class="text-3xl font-extrabold text-black sm:text-4xl">
              <span class="block">Ready to dive in?</span>
              <span class="block">Sturdy up your life</span>
            </h2>
            <p class="mt-4 text-lg leading-6 text-black">
              Get started with Sturdy today, use it standalone or on top of your existing Git or
              GitHub setup.
            </p>
            <router-link
              :to="{ name: 'download' }"
              class="mt-8 bg-white border border-transparent rounded-md shadow px-5 py-3 inline-flex items-center text-base font-medium text-yellow-600 hover:bg-yellow-50"
            >
              Download now
            </router-link>
            <p class="mt-4 text-sm leading-6 text-gray-600">
              Pssst, did you know that Sturdy is free and open-source?
            </p>
          </div>
        </div>
        <div class="-mt-6 aspect-w-5 aspect-h-3 md:aspect-w-2 md:aspect-h-1">
          <img
            class="transform translate-x-6 translate-y-6 rounded-md object-cover object-left-top sm:translate-x-16 lg:translate-y-20"
            :src="Image"
            height="655"
            width="1000"
            alt="Sturdy app screenshot"
          />
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'

import Image from '../../pages/download/app.png?w=1000'

export default defineComponent({
  setup() {
    return {
      Image,
    }
  },
})
</script>
