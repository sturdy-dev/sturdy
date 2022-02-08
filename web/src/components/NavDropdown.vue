<template>
  <Menu as="div" class="relative inline-block text-left w-full">
    <div class="flex w-full">
      <MenuButton class="flex w-full px-4 py-4 items-center hover:bg-warmgray-100 transition">
        <img
          src="../assets/Web/Duck/DuckCap256.png"
          class="h-8 w-8 flex-shrink-0"
          alt="Sturdy Duck Logo"
        />
        <div class="text-2xl font-semibold text-warmgray-800 flex-1 text-left pl-4">Sturdy</div>
        <div
          class="flex-shrink-0 flex items-center text-warmgray-400 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-offset-gray-100 focus:ring-blue-500"
        >
          <span class="sr-only">Open options</span>
          <DotsVerticalIcon class="h-5 w-5" aria-hidden="true" />
        </div>
      </MenuButton>
    </div>

    <transition
      enter-active-class="transition ease-out duration-100"
      enter-from-class="transform opacity-0 scale-95"
      enter-to-class="transform opacity-100 scale-100"
      leave-active-class="transition ease-in duration-75"
      leave-from-class="transform opacity-100 scale-100"
      leave-to-class="transform opacity-0 scale-95"
    >
      <MenuItems
        class="origin-top-right absolute left-0 ml-2 mt-2 w-56 rounded-md shadow-lg bg-white ring-1 ring-black ring-opacity-5 divide-y divide-gray-100 focus:outline-none"
      >
        <div class="py-1">
          <MenuItem v-if="authenticated" v-slot="{ active }">
            <router-link
              :to="{ name: 'home' }"
              :class="[
                active ? 'bg-gray-100 text-gray-900' : 'text-gray-700',
                'group flex items-center px-4 py-2 text-sm',
              ]"
            >
              <HomeIcon
                class="mr-3 h-5 w-5 text-gray-400 group-hover:text-gray-500"
                aria-hidden="true"
              />
              Overview
            </router-link>
          </MenuItem>

          <MenuItem v-if="authenticated" v-slot="{ active }">
            <router-link
              :to="{ name: 'codebaseCreate' }"
              :class="[
                active ? 'bg-gray-100 text-gray-900' : 'text-gray-700',
                'group flex items-center px-4 py-2 text-sm',
              ]"
            >
              <PlusIcon
                class="mr-3 h-5 w-5 text-gray-400 group-hover:text-gray-500"
                aria-hidden="true"
              />
              New codebase
            </router-link>
          </MenuItem>

          <MenuItem v-if="!ipc && authenticated" v-slot="{ active }">
            <router-link
              :to="{ name: 'installClient' }"
              :class="[
                active ? 'bg-gray-100 text-gray-900' : 'text-gray-700',
                'group flex items-center px-4 py-2 text-sm',
              ]"
            >
              <DownloadIcon
                class="mr-3 h-5 w-5 text-gray-400 group-hover:text-gray-500"
                aria-hidden="true"
              />
              Install Sturdy
            </router-link>
          </MenuItem>

          <MenuItem v-if="!ipc && !authenticated" v-slot="{ active }">
            <router-link
              :to="{ name: 'download' }"
              :class="[
                active ? 'bg-gray-100 text-gray-900' : 'text-gray-700',
                'group flex items-center px-4 py-2 text-sm',
              ]"
            >
              <DownloadIcon
                class="mr-3 h-5 w-5 text-gray-400 group-hover:text-gray-500"
                aria-hidden="true"
              />
              Download Sturdy
            </router-link>
          </MenuItem>

          <MenuItem v-slot="{ active }">
            <router-link
              :to="{ name: 'resourcesDocs' }"
              :class="[
                active ? 'bg-gray-100 text-gray-900' : 'text-gray-700',
                'group flex items-center px-4 py-2 text-sm',
              ]"
              :target="ipc ? '_blank' : ''"
            >
              <SupportIcon
                class="mr-3 h-5 w-5 text-gray-400 group-hover:text-gray-500"
                aria-hidden="true"
              />
              Help
            </router-link>
          </MenuItem>

          <MenuItem v-slot="{ active }">
            <a
              href="https://getsturdy.com/"
              :class="[
                active ? 'bg-gray-100 text-gray-900' : 'text-gray-700',
                'group flex items-center px-4 py-2 text-sm',
              ]"
              :target="ipc ? '_blank' : ''"
            >
              <ExternalLinkIcon
                class="mr-3 h-5 w-5 text-gray-400 group-hover:text-gray-500"
                aria-hidden="true"
              />
              getsturdy.com
            </a>
          </MenuItem>
        </div>
        <div class="py-1">
          <MenuItem v-if="authenticated" v-slot="{ active }">
            <a
              :class="[
                active ? 'bg-gray-100 text-gray-900' : 'text-gray-700',
                'group flex items-center px-4 py-2 text-sm cursor-pointer',
              ]"
              @click="$emit('logout')"
            >
              <LogoutIcon
                class="mr-3 h-5 w-5 text-gray-400 group-hover:text-gray-500"
                aria-hidden="true"
              />
              Sign out
            </a>
          </MenuItem>
          <MenuItem v-else v-slot="{ active }">
            <a
              :class="[
                active ? 'bg-gray-100 text-gray-900' : 'text-gray-700',
                'group flex items-center px-4 py-2 text-sm cursor-pointer',
              ]"
              @click="toLogin"
            >
              <LoginIcon
                class="mr-3 h-5 w-5 text-gray-400 group-hover:text-gray-500"
                aria-hidden="true"
              />
              Sign in
            </a>
          </MenuItem>
        </div>
      </MenuItems>
    </transition>
  </Menu>
</template>

<script>
import { Menu, MenuButton, MenuItem, MenuItems } from '@headlessui/vue'
import { useRoute } from 'vue-router'
import {
  DotsVerticalIcon,
  DownloadIcon,
  ExternalLinkIcon,
  HomeIcon,
  LogoutIcon,
  LoginIcon,
  PlusIcon,
  SupportIcon,
} from '@heroicons/vue/solid'

export default {
  components: {
    Menu,
    MenuButton,
    MenuItem,
    MenuItems,
    DotsVerticalIcon,
    PlusIcon,
    HomeIcon,
    SupportIcon,
    ExternalLinkIcon,
    LogoutIcon,
    LoginIcon,
    DownloadIcon,
  },
  props: {
    user: {
      type: Object,
    },
  },
  emits: ['logout'],
  setup() {
    return {
      ipc: window.ipc,
    }
  },
  computed: {
    authenticated() {
      return !!this.user
    },
    signInRoute() {
      return { name: 'login', params: { navigateTo: useRoute() } }
    },
  },
  methods: {
    toLogin() {
      this.$router.push({
        name: 'login',
        query: {
          navigateTo: escape(this.$route.path),
        },
      })
    },
  },
}
</script>
