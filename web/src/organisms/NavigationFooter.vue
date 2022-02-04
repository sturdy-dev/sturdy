<template>
  <div
    class="border-t pb-6 pt-9 divide-y mt-12"
    :class="[light ? '' : 'bg-slate-800/50 border-slate-700 divide-slate-700/30']"
  >
    <div class="px-6 max-w-6xl mx-auto flex flex-col md:flex-row space-y-10 md:space-y-0">
      <div
        class="flex items-center justify-between md:items-start md:flex-col md:justify-start md:w-1/3 md:space-y-8"
      >
        <div class="w-1/2">
          <router-link :to="{ name: 'v2Index' }">
            <img src="../pages/landing/assets/logotype.svg" alt="Sturdy logotype" class="h-7" />
          </router-link>
        </div>
        <SocialLinks :light="light" />
      </div>
      <!-- columns -->
      <div class="grid md:w-2/3 grid-cols-2 lg:grid-cols-4 gap-x-4 gap-y-12">
        <div v-for="(group, idx) in navigation" :key="idx" class="text-left text-sm leading-6">
          <h3
            class="text-sm font-semibold tracking-wider uppercase"
            :class="[light ? 'text-gray-800' : 'text-slate-100']"
          >
            {{ group.title }}
          </h3>
          <ul role="list" class="mt-4 space-y-2">
            <li v-for="item in group.links" :key="item.name">
              <a
                v-if="item.href"
                :href="item.href"
                class="inline-flex items-center"
                :class="[light ? 'text-gray-600' : 'hover:text-slate-100']"
              >
                {{ item.name }}
                <ExternalLinkIcon class="w-4 h-4 ml-1.5" />
              </a>
              <router-link
                v-else
                :to="{ name: item.route }"
                :class="[light ? 'text-gray-600' : 'hover:text-slate-100']"
              >
                {{ item.name }}
              </router-link>
            </li>
          </ul>
        </div>

        <template v-if="false">
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
                  <router-link :to="{ name: '' }" class="hover:text-slate-100">
                    Sources</router-link
                  >
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
        </template>
      </div>
    </div>
    <div class="px-6 max-w-6xl mx-auto mt-10 text-sm">
      <p class="pt-6" :class="[light ? 'text-gray-400' : 'text-slate-300/50']">
        &copy; {{ year }} Sturdy Sweden AB
      </p>
    </div>
  </div>
</template>

<script lang="ts">
import SocialLinks from './SocialLinks.vue'
import { defineComponent } from 'vue'
import { ExternalLinkIcon } from '@heroicons/vue/outline'

const navigation = [
  {
    title: 'Product',
    links: [
      { name: 'Get Started', route: 'v2DocsRoot' },
      { name: 'Download', route: 'download' },
      { name: 'Status', href: 'https://sturdy.statuspage.io/', external: true },
      { name: 'Pricing', route: 'pricing' },
    ],
  },
  {
    title: 'Developers',
    links: [
      { name: 'Docs', route: 'v2DocsRoot' },
      { name: 'API', route: 'resourcesApi' },
      { name: 'Source Code (soon)', href: '#' },
    ],
  },
  {
    title: 'Company',
    links: [
      { name: 'Blog', route: 'blog' },
      { name: 'Jobs', route: 'careers' },
      { name: 'Contact', route: 'contact' },
    ],
  },
  {
    title: 'Legal',
    links: [
      { name: 'Privacy', route: 'privacy' },
      { name: 'Terms', route: 'termsOfService' },
      { name: 'Security', route: 'resourcesSecurity' },
    ],
  },
]

export default defineComponent({
  components: {
    ExternalLinkIcon,
    SocialLinks,
  },
  props: {
    light: {
      type: Boolean,
      required: true,
    },
  },
  setup() {
    return {
      year: new Date().getFullYear(),
      navigation,
    }
  },
})
</script>
