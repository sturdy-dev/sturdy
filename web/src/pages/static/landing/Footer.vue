<template>
  <div class="bg-slate-800/50 border-t border-slate-700 pb-6 pt-9 divide-y divide-slate-700/30">
    <div class="px-6 max-w-6xl mx-auto flex flex-col md:flex-row space-y-10 md:space-y-0">
      <div
        class="flex items-center justify-between md:items-start md:flex-col md:justify-start md:w-1/3 md:space-y-8"
      >
        <div class="w-1/2">
          <img src="./logotype.svg" alt="Sturdy logotype" class="h-7" />
        </div>
        <SocialLinks />
      </div>
      <!-- columns -->
      <div class="flex md:w-2/3">
        <!-- left -->
        <div class="flex-none w-1/2 space-y-8 sm:space-y-8 lg:flex lg:space-y-0">
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
                <router-link :to="{ name: '' }" class="hover:text-slate-100"> Values </router-link>
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
        <div class="flex-none w-1/2 space-y-8 sm:space-y-8 lg:flex lg:space-y-0">
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
                <router-link :to="{ name: '' }" class="hover:text-slate-100"> Sources </router-link>
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
    <div class="px-6 max-w-6xl mx-auto mt-10">
      <p class="pt-6 text-slate-300/50">&copy; {{ year }} Sturdy Sweden AB</p>
    </div>
  </div>
</template>

<script lang="ts">
import SocialLinks from './SocialLinks.vue'
import { defineComponent } from 'vue'
import { ExternalLinkIcon } from '@heroicons/vue/outline'

export default defineComponent({
  components: {
    ExternalLinkIcon,
    SocialLinks,
  },
  setup() {
    return {
      year: new Date().getFullYear(),
    }
  },
})
</script>

<style scoped></style>
