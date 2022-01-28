<template>
  <div class="blurr border-b border-slate-700/10 sticky top-0 z-50">
    <div
      class="flex items-center justify-between leading-6 py-2 max-w-7xl mx-auto px-4 sm:px-6 md:px-8"
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
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'

export default defineComponent({})
</script>

<style scoped>
.blurr {
  backdrop-filter: saturate(70%) brightness(80%) blur(5px);
}
</style>
