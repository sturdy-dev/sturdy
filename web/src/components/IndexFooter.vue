<template>
  <footer class="bg-black">
    <div class="max-w-7xl mx-auto py-12 px-4 overflow-hidden sm:px-6 lg:px-8 relative">
      <img
        src="../assets/Web/Duck/DuckCap92.png"
        class="hidden md:block absolute right-8 h-24"
        title="Quack!"
      />

      <nav class="-mx-5 -my-2 flex flex-wrap justify-center" aria-label="Footer">
        <div class="px-5 py-2">
          <router-link :to="{ name: 'index' }" class="text-base text-gray-200 hover:text-gray-400">
            Home
          </router-link>
        </div>

        <div class="px-5 py-2">
          <router-link :to="{ name: 'about' }" class="text-base text-gray-200 hover:text-gray-400">
            About
          </router-link>
        </div>

        <div class="px-5 py-2">
          <router-link :to="{ name: 'blog' }" class="text-base text-gray-200 hover:text-gray-400">
            Blog
          </router-link>
        </div>

        <div class="px-5 py-2">
          <router-link
            :to="{ name: 'resourcesApi' }"
            class="text-base text-gray-200 hover:text-gray-400"
          >
            API
          </router-link>
        </div>

        <div class="px-5 py-2">
          <a
            href="https://discord.gg/5HnSdzMqtA"
            class="text-base text-gray-200 hover:text-gray-400"
          >
            Discord
          </a>
        </div>

        <div class="px-5 py-2">
          <a
            href="https://sturdy.statuspage.io/"
            class="text-base text-gray-200 hover:text-gray-400"
          >
            Status
          </a>
        </div>

        <div class="px-5 py-2">
          <a
            href="https://twitter.com/getsturdy"
            class="text-base text-gray-200 hover:text-gray-400"
          >
            Twitter
          </a>
        </div>

        <div class="px-5 py-2">
          <router-link
            :to="{ name: 'privacy' }"
            class="text-base text-gray-200 hover:text-gray-400"
          >
            Privacy Policy
          </router-link>
        </div>

        <div class="px-5 py-2">
          <router-link
            :to="{ name: 'termsOfService' }"
            class="text-base text-gray-200 hover:text-gray-400"
          >
            Terms of Service
          </router-link>
        </div>
      </nav>
      <p class="mt-8 text-center text-base text-gray-400">&copy; {{ year }} Sturdy Sweden AB</p>
    </div>
  </footer>
</template>
<script lang="ts">
import { defineComponent } from 'vue'

export default defineComponent({
  setup() {
    return {
      year: new Date().getFullYear(),
    }
  },
})
</script>
