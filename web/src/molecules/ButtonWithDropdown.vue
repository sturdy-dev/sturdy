<template>
  <Menu v-slot="{ open }" class="inline-block relative" as="div">
    <div
      ref="buttonGroup"
      class="leading-4 rounded-md relative inline-flex items-center font-medium flex-shrink-0 group leading-5 items-stretch"
      :class="[wrapperStyle]"
    >
      <Button
        :grouped="!!$slots.dropdown"
        :first="!!$slots.dropdown"
        :disabled="disabled"
        :color="color"
        :spinner="spinner"
        :show-tooltip="showTooltip"
        :size="size"
        :tooltip-right="tooltipRight"
        @click="$emit('click')"
      >
        <template #tooltip>
          <slot name="tooltip" />
        </template>
        <template #default>
          <slot />
        </template>
      </Button>

      <template v-if="$slots.dropdown">
        <MenuButton
          class="px-1.5 py-2 rounded-r-md border focus:z-10 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
          :class="[menuButtonStyle(open), disabled ? 'opacity-80' : '']"
        >
          <ChevronDownIcon class="w-5 h-5" />
        </MenuButton>
      </template>
    </div>

    <template v-if="$slots.dropdown">
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
    </template>
  </Menu>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import type { PropType } from 'vue'
import { ChevronDownIcon } from '@heroicons/vue/solid'
import { Menu, MenuItems, MenuButton } from '@headlessui/vue'
import Button from '../atoms/Button.vue'

export default defineComponent({
  components: { ChevronDownIcon, Menu, MenuItems, MenuButton, Button },
  props: {
    color: {
      type: String as PropType<'default' | 'green' | 'blue'>,
      default: 'default',
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    spinner: {
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
    size: {
      type: String,
      default: 'normal',
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
