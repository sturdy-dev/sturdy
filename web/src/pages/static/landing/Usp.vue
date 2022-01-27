<template>
  <section :id="title.toLowerCase()">
    <div class="relative pt-10 xl:pt-0 mt-10 xl:mt-2 space-y-32 sm:space-y-40 md:space-y-60">
      <div
        class="hidden dark:block absolute top-0 inset-x-0 h-[37.5rem] bg-gradient-to-b from-[#0c1120] top-0 xl:top-8"
      ></div>
      <div
        class="relative pt-16 max-w-7xl px-4 sm:px-6 md:px-8 mx-auto flex flex-col items-center md:flex-row lg:gap-8 space-y-6"
      >
        <div class="md:w-1/2">
          <h2 class="mt-8 font-semibold text-amber-500">{{ title }}</h2>
          <p class="mt-4 text-3xl sm:text-4xl font-extrabold tracking-tight text-slate-50">
            {{ subtitle }}
          </p>
          <div class="mt-4 max-w-3xl space-y-6">
            <slot name="left"></slot>
            <router-link
              class="text-slate-900 bg-amber-500 hover:text-amber-500 hover:bg-slate-900/10 border border-transparent hover:border-amber-500 focus:outline-none focus:ring-2 focus:ring-slate-400 focus:ring-offset-2 focus:ring-offset-slate-50 focus:outline-none focus:ring-2 group inline-flex items-center h-9 rounded-full text-sm font-semibold whitespace-nowrap px-3"
              :to="{ name: link }"
            >
              Learn more
              <ChevronRightIcon class="h-4 w-4 ml-1" />
            </router-link>
          </div>
          <div class="mt-12 max-w-3xl space-y-6">
            <slot name="bottomLeft"></slot>
          </div>
        </div>
        <div class="md:w-1/2">
          <slot name="right"></slot>
        </div>
      </div>
    </div>
  </section>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import { ChevronRightIcon } from '@heroicons/vue/outline'

export default defineComponent({
  components: { ChevronRightIcon },
  props: {
    title: {
      type: String,
      required: true,
    },
    subtitle: {
      type: String,
      required: true,
    },
    link: {
      type: String,
      required: true,
    },
  },
})
</script>

<style scoped></style>
