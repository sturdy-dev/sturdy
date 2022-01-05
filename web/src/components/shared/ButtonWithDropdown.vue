<template>
  <Menu v-slot="{ open }" class="inline-block">
    <div class="relative">
      <div
        ref="buttonGroup"
        class="rounded-md border leading-4 relative inline-flex items-center font-medium flex-shrink-0 group leading-5 items-stretch divide-x"
        :class="[wrapperStyle, disabled ? 'opacity-80' : '']"
      >
        <button
          class="px-4 py-2 focus:z-10 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
          :class="[$slots.dropdown ? 'rounded-l-md' : 'rounded-md', buttonStyle]"
          :disabled="disabled"
          @click="$emit('click')"
        >
          <span class="text-sm font-medium contents">
            <slot></slot>
          </span>
        </button>
        <template v-if="$slots.dropdown">
          <MenuButton
            class="px-1.5 py-2 rounded-r-md focus:z-10 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
            :class="[menuButtonStyle(open), disabled ? 'opacity-80' : '']"
            :disabled="disabled"
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
            <slot name="dropdown"></slot>
          </MenuItems>
        </transition>
      </div>
    </div>
  </Menu>
</template>
<script lang="ts">
import { defineComponent, PropType } from 'vue'
import { ChevronDownIcon } from '@heroicons/vue/solid'
import { Menu, MenuItems, MenuButton } from '@headlessui/vue'

export default defineComponent({
  name: 'ButtonWithDropdown',
  components: { ChevronDownIcon, Menu, MenuItems, MenuButton },
  props: {
    color: {
      type: String as PropType<'default' | 'green'>,
      default: 'default',
    },
    disabled: {
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
      }
    },
    // eslint-disable-next-line vue/return-in-computed-property
    buttonStyle(): string {
      switch (this.color) {
        case 'default':
          return 'hover:bg-gray-50'
        case 'green':
          return 'hover:bg-green-200'
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
      }
    },
  },
})
</script>
