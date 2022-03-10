<template>
  <Menu as="div" class="relative inline-block text-left w-full">
    <div v-if="currentOrganization" class="flex w-full">
      <MenuButton class="flex w-full px-4 py-4 items-center hover:bg-warmgray-100 transition">
        <div
          class="font-semibold text-warmgray-800 flex-1 text-left whitespace-nowrap overflow-x-hidden ellipsis text-ellipsis"
          :class="[fontSize(currentOrganization.name)]"
        >
          {{ currentOrganization.name }}
        </div>
        <div
          class="flex-shrink-0 flex items-center text-warmgray-400 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-offset-gray-100 focus:ring-blue-500"
        >
          <span class="sr-only">Open options</span>
          <DotsVerticalIcon class="h-5 w-5" aria-hidden="true" />
        </div>
      </MenuButton>
    </div>

    <div v-else class="flex w-full">
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
        <div v-if="authenticated && currentOrganization" class="py-1">
          <MenuItem v-slot="{ active }">
            <router-link
              :to="{
                name: 'organizationListCodebases',
                params: { organizationSlug: currentOrganization.shortID },
              }"
              :class="[
                active ? 'bg-gray-100 text-gray-900' : 'text-gray-700',
                'group flex items-center px-4 py-2 text-sm',
              ]"
            >
              <HomeIcon
                class="mr-3 h-5 w-5 text-gray-400 group-hover:text-gray-500"
                aria-hidden="true"
              />
              {{ currentOrganization.name }}
            </router-link>
          </MenuItem>

          <MenuItem v-slot="{ active }">
            <router-link
              :to="{
                name: 'organizationCreateCodebase',
                params: { organizationSlug: currentOrganization.shortID },
              }"
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

          <MenuItem v-slot="{ active }">
            <router-link
              :to="{
                name: 'organizationSettings',
                params: { organizationSlug: currentOrganization.shortID },
              }"
              :class="[
                active ? 'bg-gray-100 text-gray-900' : 'text-gray-700',
                'group flex items-center px-4 py-2 text-sm',
              ]"
            >
              <CogIcon
                class="mr-3 h-5 w-5 text-gray-400 group-hover:text-gray-500"
                aria-hidden="true"
              />
              Settings
            </router-link>
          </MenuItem>
        </div>

        <div
          v-if="authenticated && nonCurrentOrganizations && nonCurrentOrganizations.length > 0"
          class="py-1"
        >
          <MenuItem v-for="org in nonCurrentOrganizations" :key="org.id" v-slot="{ active }">
            <router-link
              :to="{ name: 'organizationListCodebases', params: { organizationSlug: org.shortID } }"
              :class="[
                active ? 'bg-gray-100 text-gray-900' : 'text-gray-700',
                'group flex items-center px-4 py-2 text-sm',
              ]"
            >
              <SwitchHorizontalIcon
                class="mr-3 h-5 w-5 text-gray-400 group-hover:text-gray-500"
                aria-hidden="true"
              />
              {{ org.name }}
            </router-link>
          </MenuItem>
        </div>

        <div class="py-1">
          <MenuItem v-if="authenticated && isMultiTenancyEnabled" v-slot="{ active }">
            <router-link
              :to="{ name: 'organizationCreate' }"
              :class="[
                active ? 'bg-gray-100 text-gray-900' : 'text-gray-700',
                'group flex items-center px-4 py-2 text-sm',
              ]"
            >
              <UserGroupIcon
                class="mr-3 h-5 w-5 text-gray-400 group-hover:text-gray-500"
                aria-hidden="true"
              />
              New organization
            </router-link>
          </MenuItem>

          <MenuItem v-if="showInstallCLI" v-slot="{ active }">
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

          <MenuItem v-if="showInstallDownloadApp" v-slot="{ active }">
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
            <a
              href="https://getsturdy.com/docs"
              :class="[
                active ? 'bg-gray-100 text-gray-900' : 'text-gray-700',
                'group flex items-center px-4 py-2 text-sm',
              ]"
              :target="isApp ? '_blank' : ''"
            >
              <SupportIcon
                class="mr-3 h-5 w-5 text-gray-400 group-hover:text-gray-500"
                aria-hidden="true"
              />
              Help
            </a>
          </MenuItem>

          <MenuItem v-slot="{ active }">
            <router-link
              to="/"
              :class="[
                active ? 'bg-gray-100 text-gray-900' : 'text-gray-700',
                'group flex items-center px-4 py-2 text-sm',
              ]"
              :target="isApp ? '_blank' : ''"
            >
              <ExternalLinkIcon
                class="mr-3 h-5 w-5 text-gray-400 group-hover:text-gray-500"
                aria-hidden="true"
              />
              getsturdy.com
            </router-link>
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

<script lang="ts">
import { Menu, MenuButton, MenuItem, MenuItems } from '@headlessui/vue'
import { useRoute } from 'vue-router'
import {
  CogIcon,
  DotsVerticalIcon,
  DownloadIcon,
  ExternalLinkIcon,
  HomeIcon,
  LoginIcon,
  LogoutIcon,
  PlusIcon,
  SupportIcon,
  SwitchHorizontalIcon,
  UserGroupIcon,
} from '@heroicons/vue/solid'
import { gql } from '@urql/vue'
import { computed, defineComponent, inject, PropType, ref, Ref } from 'vue'
import { NavigationOrganizationPickerMenuFragment } from './__generated__/NavigationOrganizationPickerMenu'
import { Feature } from '../__generated__/types'

export const ORGANIZATION_FRAGMENT = gql`
  fragment NavigationOrganizationPickerMenu on Organization {
    id
    name
    shortID
  }
`

export default defineComponent({
  components: {
    Menu,
    MenuButton,
    MenuItem,
    MenuItems,
    DotsVerticalIcon,
    PlusIcon,
    HomeIcon,
    SupportIcon,
    SwitchHorizontalIcon,
    ExternalLinkIcon,
    LogoutIcon,
    LoginIcon,
    CogIcon,
    UserGroupIcon,
    DownloadIcon,
  },
  props: {
    user: {
      type: Object,
      required: false,
    },
    organizations: {
      type: Object as PropType<Array<NavigationOrganizationPickerMenuFragment>>,
      required: false,
    },
    currentOrganization: {
      type: Object as PropType<NavigationOrganizationPickerMenuFragment>,
      required: false,
    },
  },
  emits: ['logout'],
  setup() {
    const features = inject<Ref<Array<Feature>>>('features', ref([]))
    const isMultiTenancyEnabled = computed(() => features.value.includes(Feature.MultiTenancy))

    return {
      isMultiTenancyEnabled,
    }
  },
  data() {
    const ipc = window.ipc
    return {
      ipc,
    }
  },
  computed: {
    isApp() {
      return !!this.ipc
    },
    authenticated() {
      return !!this.user
    },
    signInRoute() {
      return { name: 'login', params: { navigateTo: useRoute() } }
    },
    nonCurrentOrganizations() {
      return this.organizations.filter((org) => org.id !== this.currentOrganization.id)
    },
    showInstallCLI() {
      return !this.isApp && this.authenticated && this.isMultiTenancyEnabled
    },
    showInstallDownloadApp() {
      return !this.isApp && !this.authenticated
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
    fontSize(title: string): string {
      if (title.length < 10) {
        return 'text-2xl'
      }
      if (title.length < 15) {
        return 'text-xl'
      }
      return 'text-lg'
    },
  },
})
</script>
