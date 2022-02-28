<template>
  <Disclosure
    v-slot="{ open }"
    as="nav"
    :class="[
      light
        ? ''
        : 'bg-slate-900 border-slate-800/10 sm:bg-transparent sm:backdrop-blur-[5px] sm:bg-slate-900/50',
      ' border-b',
    ]"
  >
    <div class="max-w-6xl mx-auto px-6">
      <div class="flex items-center justify-between h-14">
        <div class="flex items-center">
          <div class="flex-shrink-0">
            <router-link :to="{ name: 'v2Index' }">
              <img
                src="../pages/landing/assets/logotype.svg"
                alt="Sturdy logotype"
                class="h-7"
                height="28"
                width="110"
              />
            </router-link>
          </div>
        </div>

        <div class="hidden sm:ml-6 sm:block">
          <div class="flex items-center">
            <div class="hidden sm:block sm:ml-6">
              <div
                class="flex items-center space-x-8 font-semibold text-sm focus:outline-none focus:ring-2 focus:ring-slate-400 focus:ring-offset-2 focus:ring-offset-slate-50 focus:outline-none focus:ring-2"
              >
                <router-link
                  :to="{ name: 'v2DocsRoot' }"
                  :class="[light ? 'hover:text-gray-500' : 'hover:text-amber-500']"
                >
                  Docs
                </router-link>
                <router-link
                  :to="{ name: 'v2pricing' }"
                  :class="[light ? 'hover:text-gray-500' : 'hover:text-amber-500']"
                >
                  Pricing
                </router-link>
                <router-link
                  :to="{ name: 'blog' }"
                  :class="[light ? 'hover:text-gray-500' : 'hover:text-amber-500']"
                >
                  Blog
                </router-link>

                <GitHubButton :is-light="light" />

                <ClientOnly>
                  <router-link
                    v-if="user"
                    :to="{ name: 'home' }"
                    :class="[
                      light
                        ? 'hover:bg-amber-400 text-slate-800'
                        : 'hover:text-amber-500 hover:bg-transparent',
                    ]"
                    class="text-slate-900 bg-amber-500 border border-transparent hover:border-amber-500 font-semibold h-9 px-3 rounded flex items-center justify-center sm:w-auto highlight-white/20 focus:outline-none focus:ring-2 focus:ring-slate-400 focus:ring-offset-2 focus:ring-offset-slate-50 focus:outline-none focus:ring-2"
                  >
                    Go to codebases
                  </router-link>
                  <router-link
                    v-else
                    :to="{ name: 'download' }"
                    :class="[
                      light
                        ? 'hover:bg-amber-400 text-slate-800'
                        : 'hover:text-amber-500 hover:bg-transparent',
                    ]"
                    class="text-slate-900 bg-amber-500 border border-transparent hover:border-amber-500 font-semibold h-9 px-3 rounded flex items-center justify-center sm:w-auto highlight-white/20 focus:outline-none focus:ring-2 focus:ring-slate-400 focus:ring-offset-2 focus:ring-offset-slate-50 focus:outline-none focus:ring-2"
                  >
                    Get started
                  </router-link>
                </ClientOnly>
              </div>
            </div>
          </div>
        </div>

        <div class="-mr-2 flex sm:hidden">
          <!-- Mobile menu button -->
          <DisclosureButton
            class="inline-flex items-center justify-center p-2 rounded text-gray-400 focus:outline-none focus:ring-0"
          >
            <span class="sr-only">Open main menu</span>
            <MenuIcon v-if="!open" class="block h-6 w-6" aria-hidden="true" />
            <XIcon v-else class="block h-6 w-6" aria-hidden="true" />
          </DisclosureButton>
        </div>
      </div>
    </div>

    <DisclosurePanel class="sm:hidden mt-10">
      <div class="px-2 pt-2 pb-3 h-screen">
        <router-link
          :to="{ name: 'download' }"
          :class="[
            light
              ? 'hover:bg-amber-400 text-slate-800'
              : 'hover:text-amber-500 hover:bg-transparent',
          ]"
          class="mx-3 block text-slate-900 bg-amber-500 border border-transparent hover:border-amber-500 font-semibold h-9 px-3 rounded flex items-center justify-center sm:w-auto highlight-white/20 focus:outline-none focus:ring-2 focus:ring-slate-400 focus:ring-offset-2 focus:ring-offset-slate-50 focus:outline-none focus:ring-2"
        >
          Get started
        </router-link>

        <router-link
          :to="{ name: 'v2DocsRoot' }"
          :class="[
            light
              ? 'text-gray-800 hover:text-gray-500'
              : 'text-gray-300 hover:bg-slate-800/50 hover:text-white',
          ]"
          class="block mx-3 py-4 text-base font-medium border-b border-slate-700"
        >
          Docs
        </router-link>
        <router-link
          :to="{ name: 'v2pricing' }"
          :class="[
            light
              ? 'text-gray-800 hover:text-gray-500'
              : 'text-gray-300 hover:bg-slate-800/50 hover:text-white',
          ]"
          class="block mx-3 py-4 text-base font-medium border-b border-slate-700"
        >
          Pricing
        </router-link>
        <router-link
          :to="{ name: 'blog' }"
          :class="[
            light
              ? 'text-gray-800 hover:text-gray-500'
              : 'text-gray-300 hover:bg-slate-800/50 hover:text-white',
          ]"
          class="block mx-3 py-4 text-base font-medium border-b border-slate-700"
        >
          Blog
        </router-link>
        <a
          href="https://github.com/sturdy-dev/sturdy"
          :class="[
            light
              ? 'text-gray-800 hover:text-gray-500'
              : 'text-gray-300 hover:bg-slate-800/50 hover:text-white',
          ]"
          class="block mx-3 py-4 text-base font-medium border-b border-slate-700"
        >
          GitHub Repository
        </a>
      </div>
    </DisclosurePanel>
  </Disclosure>
</template>

<script lang="ts" setup>
import { defineProps, inject, withDefaults } from 'vue'
import { Disclosure, DisclosureButton, DisclosurePanel } from '@headlessui/vue'
import { MenuIcon, XIcon } from '@heroicons/vue/outline'
import GitHubButton from '../molecules/GitHubButton.vue'

interface Props {
  light?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  light: false,
})

let user = inject('user')
</script>
