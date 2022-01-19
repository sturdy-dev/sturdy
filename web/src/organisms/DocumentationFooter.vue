<template>
  <footer class="bg-gray-100" aria-labelledby="footer-heading">
    <h2 id="footer-heading" class="sr-only">Footer</h2>
    <div class="max-w-7xl mx-auto py-12 px-4 sm:px-6 lg:py-16 lg:px-8">
      <div class="mt-12 grid sm:grid-cols-2 xl:grid-cols-5 gap-12 sm:gap-8 xl:mt-0 xl:col-span-1">
        <div v-for="(group, idx) in navigation" :key="idx" class="text-center sm:text-left">
          <h3 class="text-sm font-semibold text-gray-400 tracking-wider uppercase">
            {{ group.title }}
          </h3>
          <ul role="list" class="mt-4 space-y-4">
            <li v-for="item in group.links" :key="item.name">
              <a
                v-if="item.href"
                :href="item.href"
                class="text-base text-gray-500 hover:text-gray-900"
              >
                {{ item.name }}
              </a>
              <router-link
                v-else
                :to="{ name: item.route }"
                class="text-base text-gray-500 hover:text-gray-900"
              >
                {{ item.name }}
              </router-link>
            </li>
          </ul>
        </div>
      </div>
      <div class="mt-12 border-t border-gray-200 pt-8">
        <p class="text-base text-gray-400 text-center">&copy; {{ year }} Sturdy Sweden AB</p>
      </div>
    </div>
  </footer>
</template>

<script lang="ts">
import { defineComponent } from 'vue'

const navigation = [
  {
    title: 'Product',
    links: [
      // {name: 'Pricing', route: 'pricing'},
      // {name: 'Documentation', route: 'v2DocsProductIntro'},
      // // { name: 'Guides', href: '#' },
      // {name: 'API Status', href: 'https://sturdy.statuspage.io/'},
      { name: 'Get Started', route: '' },
      { name: 'Download', route: 'download' },
      { name: 'Open Source', href: '#' },
      { name: 'API', route: 'resourcesApi' },
    ],
  },
  {
    title: 'Support',
    links: [
      { name: 'Pricing', route: 'pricing' },
      { name: 'Documentation', route: 'v2DocsProductIntro' },
      { name: 'API Status', href: 'https://sturdy.statuspage.io/' },
      { name: 'Contact', route: 'contact' },
      { name: 'Discord', href: 'https://discord.gg/5HnSdzMqtA' },
    ],
  },
  {
    title: 'Community',
    links: [
      { name: 'Discord', href: 'https://discord.gg/5HnSdzMqtA' },
      { name: 'Twitter', href: 'https://twitter.com/getsturdy' },
      { name: 'Open Source', href: '#' },
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
  setup() {
    return {
      year: new Date().getFullYear(),
      navigation,
    }
  },
})
</script>
