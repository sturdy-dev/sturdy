<template>
  <section :id="title.toLowerCase()">
    <div
      class="max-w-7xl mx-auto px-4 sm:px-6 md:px-8 flex flex-col items-center md:flex-row lg:gap-8 space-y-6"
    >
      <div class="md:w-1/2">
        <h2 class="mt-8 font-semibold text-amber-500">{{ title }}</h2>
        <p class="mt-4 text-3xl sm:text-4xl font-extrabold tracking-tight text-slate-50">
          {{ subtitle }}
        </p>
        <div class="mt-4 max-w-3xl space-y-6">
          <slot name="left"></slot>
          <router-link
            class="bg-slate-600 text-slate-100 hover:bg-slate-500 hover:text-white focus:outline-none focus:ring-2 focus:ring-slate-400 focus:ring-offset-2 focus:ring-offset-slate-50 group inline-flex items-center h-9 rounded-full text-sm font-semibold whitespace-nowrap px-3 focus:outline-none focus:ring-2"
            :to="{ name: link }"
          >
            Learn more
            <ChevronRightIcon class="h-4 w-4 ml-1 text-slate-400" />
          </router-link>
        </div>
      </div>
      <div class="md:w-1/2">
        <slot name="right"></slot>
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
