<template>
  <Disclosure
    as="nav"
    class="backdrop-blur-[5px] backdrop-brightness-[.80] backdrop-saturate-[.70] border-b border-slate-700/10 sticky top-0 z-50"
    v-slot="{ open }"
  >
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      <div class="flex items-center justify-between h-10 sm:h-14">
        <div class="flex items-center">
          <div class="flex-shrink-0">
            <img src="./logotype.svg" alt="Sturdy logotype" class="h-7" />
          </div>
        </div>

        <div class="hidden sm:ml-6 sm:block">
          <div class="flex items-center">
            <div class="hidden sm:block sm:ml-6">
              <div
                class="flex items-center space-x-8 font-semibold text-sm focus:outline-none focus:ring-2 focus:ring-slate-400 focus:ring-offset-2 focus:ring-offset-slate-50 focus:outline-none focus:ring-2"
              >
                <router-link :to="{ name: 'v2DocsRoot' }" class="hover:text-amber-500">
                  Docs
                </router-link>
                <router-link :to="{ name: 'pricing' }" class="hover:text-amber-500">
                  Pricing
                </router-link>
                <router-link :to="{ name: 'blog' }" class="hover:text-amber-500">
                  Blog
                </router-link>

                <router-link
                  :to="{ name: 'download' }"
                  class="text-slate-900 bg-amber-500 hover:text-amber-500 hover:bg-slate-900/10 border border-transparent hover:border-amber-500 font-semibold h-9 px-3 rounded flex items-center justify-center sm:w-auto highlight-white/20 focus:outline-none focus:ring-2 focus:ring-slate-400 focus:ring-offset-2 focus:ring-offset-slate-50 focus:outline-none focus:ring-2"
                >
                  Get started
                </router-link>
              </div>
            </div>
          </div>
        </div>

        <div class="-mr-2 flex sm:hidden">
          <!-- Mobile menu button -->
          <DisclosureButton
            class="inline-flex items-center justify-center p-2 rounded text-gray-400 border border-transparent hover:border-slate-500 focus:outline-none focus:ring-0"
          >
            <span class="sr-only">Open main menu</span>
            <MenuIcon v-if="!open" class="block h-6 w-6" aria-hidden="true" />
            <XIcon v-else class="block h-6 w-6" aria-hidden="true" />
          </DisclosureButton>
        </div>
      </div>
    </div>

    <DisclosurePanel class="sm:hidden">
      <div class="px-2 pt-2 pb-3 space-y-1 h-screen">
        <router-link
          :to="{ name: 'download' }"
          class="block mx-3 py-2 text-slate-900 bg-amber-500 hover:text-amber-500 hover:bg-slate-900/10 border border-transparent hover:border-amber-500 font-semibold h-9 px-3 rounded flex items-center justify-center sm:w-auto highlight-white/20 focus:outline-none focus:ring-2 focus:ring-slate-400 focus:ring-offset-2 focus:ring-offset-slate-50 focus:outline-none focus:ring-2"
        >
          Get started
        </router-link>

        <router-link
          :to="{ name: 'v2DocsRoot' }"
          class="text-gray-300 block mx-3 py-2 text-base font-medium border-b border-slate-700 hover:bg-slate-800/50 hover:text-white"
        >
          Docs
        </router-link>
        <router-link
          :to="{ name: 'pricing' }"
          class="text-gray-300 block mx-3 py-2 text-base font-medium border-b border-slate-700 hover:bg-slate-800/50 hover:text-white"
        >
          Pricing
        </router-link>
        <router-link
          :to="{ name: 'blog' }"
          class="text-gray-300 block mx-3 py-2 text-base font-medium border-b border-slate-700 hover:bg-slate-800/50 hover:text-white"
        >
          Blog
        </router-link>
      </div>
    </DisclosurePanel>
  </Disclosure>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import { Disclosure, DisclosureButton, DisclosurePanel } from '@headlessui/vue'
import { MenuIcon, XIcon } from '@heroicons/vue/outline'

export default defineComponent({
  components: {
    Disclosure,
    DisclosureButton,
    DisclosurePanel,
    MenuIcon,
    XIcon,
  },
})
</script>

<style scoped></style>
