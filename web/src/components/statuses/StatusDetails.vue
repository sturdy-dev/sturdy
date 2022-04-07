<template>
  <Menu v-if="!unknown" as="div" class="relative inline-block text-left">
    <MenuButton
      class="rounded-full flex items-center space-x-2 text-gray-400 hover:text-gray-600 focus:outline-none"
    >
      <StatusBadge :statuses="statuses" />
      <p v-if="showText" class="text-gray-900 text-sm font-medium cursor-pointer">{{ text }}</p>
    </MenuButton>

    <transition
      enter-active-class="transition ease-out duration-100"
      enter-from-class="transform opacity-0 scale-95"
      enter-to-class="transform opacity-100 scale-100"
      leave-active-class="transition ease-in duration-75"
      leave-from-class="transform opacity-100 scale-100"
      leave-to-class="transform opacity-0 scale-95"
    >
      <MenuItems
        class="origin-top-right absolute right-0 mt-2 min-w-max rounded-md shadow-lg bg-white ring-1 ring-black ring-opacity-5 focus:outline-none z-20"
      >
        <div v-for="status in sortedStatuses" :key="status.id">
          <MenuItem v-slot="{ active }">
            <div
              :class="[
                active ? 'bg-gray-100 text-gray-900' : 'text-gray-700',
                'flex items-center p-3 text-sm cursor-default gap-2',
              ]"
            >
              <StatusBadge class="h-5 w-5" :statuses="[status]" />
              <p>
                {{ status.title }}
                <template v-if="status.description">― {{ status.description }}</template>
              </p>
              <div class="flex-grow"></div>
              <a
                v-if="status.detailsUrl"
                :href="status.detailsUrl"
                class="underline text-indigo-600 font-medium ml-2"
                target="_blank"
                >Details</a
              >
            </div>
          </MenuItem>
        </div>
      </MenuItems>
    </transition>
  </Menu>
</template>

<script lang="ts">
import type { PropType } from 'vue'
import { defineComponent } from 'vue'
import { StatusType } from '../../__generated__/types'
import StatusBadge, { STATUS_FRAGMENT } from './StatusBadge.vue'
import { Menu, MenuButton, MenuItem, MenuItems } from '@headlessui/vue'
import type { StatusFragment } from './__generated__/StatusBadge'

export { STATUS_FRAGMENT }

export default defineComponent({
  components: { StatusBadge, Menu, MenuButton, MenuItem, MenuItems },
  props: {
    statuses: {
      type: Array as PropType<StatusFragment[]>,
      default: () => [],
      required: true,
    },
    showText: {
      type: Boolean,
      required: false,
      default: () => {
        return true
      },
    },
  },
  computed: {
    sortedStatuses() {
      const copy = this.statuses
      return copy.sort((a, b) => a.title.localeCompare(b.title))
    },
    unknown(): boolean {
      return this.statuses.length === 0
    },
    pending(): boolean {
      if (this.unknown) return false
      return this.statuses.some((s) => s.type == StatusType.Pending)
    },
    healthy(): boolean {
      if (this.unknown) return false
      return !this.statuses.some((s) => s.type != StatusType.Healthy)
    },
    failing(): boolean {
      if (this.unknown) return false
      return this.statuses.some((s) => s.type == StatusType.Failing)
    },

    text(): string {
      switch (true) {
        case this.failing:
          return 'Failing'
        case this.pending:
          return 'Pending'
        case this.healthy:
          return 'Healthy'
        default:
          return 'Unknown'
      }
    },
  },
})
</script>
