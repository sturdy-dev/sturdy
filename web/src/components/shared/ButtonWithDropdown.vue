<template>
  <Menu v-slot="{ open }" class="inline-block">
    <div class="relative">
      <div
        ref="buttonGroup"
        class="rounded-md border leading-4 relative inline-flex items-center font-medium flex-shrink-0 group leading-5 items-stretch divide-x"
        :class="[wrapperStyle, disabled ? 'opacity-80' : 'cursor-pointer']"
      >
        <button
          class="px-4 py-2 focus:z-10 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
          :class="[$slots.dropdown ? 'rounded-l-md' : 'rounded-md', buttonStyle]"
          :disabled="disabled"
          @click="$emit('click')"
        >
          <div
            v-if="showTooltip"
            class="absolute bottom-0 flex-col mb-8 hidden group-hover:flex z-50"
            :class="[tooltipRight ? 'right-1' : 'left-1']"
          >
            <span
              class="relative p-2 text-xs leading-none rounded text-white whitespace-nowrap bg-black shadow-lg"
            >
              <slot name="tooltip"></slot>
            </span>
            <div
              class="w-3 h-3 transform rotate-45 bg-black absolute bottom-0 -mb-1"
              :class="[tooltipRight ? 'right-3' : 'left-3']"
            ></div>
          </div>

          <span class="text-sm font-medium contents">
            <slot></slot>
          </span>
        </button>
        <template v-if="$slots.dropdown">
          <MenuButton
            class="px-1.5 py-2 rounded-r-md focus:z-10 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
            :class="[menuButtonStyle(open), disabled ? 'opacity-80' : '']"
          >
            <ChevronDownIcon class="w-5 h-5" />
          </MenuButton>
        </template>
      </div>
      <div v-if="$slots.dropdown">
        <transition
          enter-active-class="transition ease-out duration-100"
          enter-from-class="transform opacity-0 scale-95"
          enter-to-class="transform opacity-100 scale-100"
          leave-active-class="transition ease-in duration-75"
          leave-from-class="transform opacity-100 scale-100"
          leave-to-class="transform opacity-0 scale-95"
        >
          <MenuItems
            class="origin-top absolute z-10 top-full left-0 bg-white shadow-lg rounded-md border mt-1 flex flex-col divide-y divide-gray-200 whitespace-nowrap"
            :style="{ minWidth: `${$refs.buttonGroup?.getBoundingClientRect()?.width || 200}px` }"
          >
            <slot name="dropdown" :disabled="disabled"></slot>
          </MenuItems>
        </transition>
      </div>
    </div>
  </Menu>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import type { PropType } from 'vue'
import { ChevronDownIcon } from '@heroicons/vue/solid'
import { Menu, MenuItems, MenuButton } from '@headlessui/vue'

export default defineComponent({
  components: { ChevronDownIcon, Menu, MenuItems, MenuButton },
  props: {
    color: {
      type: String as PropType<'default' | 'green' | 'blue'>,
      default: 'default',
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    showTooltip: {
      type: Boolean,
      default: false,
    },
    tooltipRight: {
      type: Boolean,
      default: false,
    },
  },
  emits: ['click'],
  computed: {
    // eslint-disable-next-line vue/return-in-computed-property
    wrapperStyle(): string {
      switch (this.color) {
        case 'default':
          return 'divide-gray-200 bg-white border-gray-200 text-gray-700'
        case 'green':
          return 'divide-green-200 bg-green-100 border-green-200'
        case 'blue':
          return 'divide-blue-700 bg-blue-600 border-blue-700 text-white'
      }
    },
    // eslint-disable-next-line vue/return-in-computed-property
    buttonStyle(): string {
      switch (this.color) {
        case 'default':
          return 'hover:bg-gray-50'
        case 'green':
          return 'hover:bg-green-200'
        case 'blue':
          return 'hover:bg-blue-700'
      }
    },
  },
  methods: {
    menuButtonStyle(open: boolean): string[] {
      switch (this.color) {
        case 'default':
          return ['hover:bg-gray-50', open ? 'bg-gray-100' : '']
        case 'green':
          return ['hover:bg-green-200', open ? 'bg-green-200' : '']
        case 'blue':
          return ['hover:bg-blue-700', open ? 'bg-blue-700' : '']
      }
    },
  },
})
</script>
