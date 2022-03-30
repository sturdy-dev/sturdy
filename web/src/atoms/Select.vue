<template>
  <Listbox
    v-slot="{ open }"
    v-model="selected"
    class="inline-block relative flex border rounded-md"
    as="div"
    :class="wrapperStyle"
  >
    <slot name="selected" :item="selected" />
    <ListboxButton
      class="hover:bg-gray-50 px-1.5 py-2 border-l rounded-r-md focus:z-10 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
      :class="buttonStyle(open)"
    >
      <ChevronDownIcon class="w-5 h-5" />
    </ListboxButton>
    <transition
      enter-active-class="transition ease-out duration-100"
      enter-from-class="transform opacity-0 scale-95"
      enter-to-class="transform opacity-100 scale-100"
      leave-active-class="transition ease-in duration-75"
      leave-from-class="transform opacity-100 scale-100"
      leave-to-class="transform opacity-0 scale-95"
    >
      <ListboxOptions
        class="origin-top absolute z-10 top-full left-0 bg-white shadow-lg rounded-md border mt-1 flex flex-col divide-y divide-gray-200 whitespace-nowrap"
      >
        <ListboxOption
          v-for="(option, i) in options"
          :key="i"
          :value="option"
          class="cursor-pointer"
          @click="onOptionSelected(i)"
        >
          <component :is="option" class="w-full pointer-events-none" />
        </ListboxOption>
      </ListboxOptions>
    </transition>
  </Listbox>
</template>

<script lang="ts">
import {
  Listbox,
  ListboxButton,
  ListboxOptions,
  ListboxOption,
  ListboxLabel,
} from '@headlessui/vue'
import { ChevronDownIcon } from '@heroicons/vue/solid'

import { defineComponent, toRefs, type PropType } from 'vue'

export default defineComponent({
  props: {
    color: {
      type: String as PropType<'white' | 'green' | 'blue'>,
      default: 'white',
    },
    id: {
      type: String,
      required: false,
    },
  },
  components: {
    ChevronDownIcon,
    Listbox,
    ListboxButton,
    ListboxOptions,
    ListboxOption,
    ListboxLabel,
  },
  setup(props, { slots }) {
    const { id } = toRefs(props)
    const selectedIdx = localStorage.getItem(id.value)
    return {
      options: slots.options(),
      selected: slots.options()[selectedIdx || 0],
    }
  },
  computed: {
    wrapperStyle() {
      switch (this.color) {
        case 'green':
          return 'divide-green-200 bg-green-100 border-green-200'
        case 'blue':
          return 'divide-blue-700 bg-blue-600 border-blue-700 text-white'
        default:
          return 'divide-gray-200 bg-white border-gray-200 text-gray-700'
      }
    },
  },
  methods: {
    buttonStyle(open: boolean) {
      switch (this.color) {
        case 'green':
          return ['hover:bg-green-200', open ? 'bg-green-200' : '']
        case 'blue':
          return ['hover:bg-blue-700', open ? 'bg-blue-700' : '']
        default:
          return ['hover:bg-gray-50', open ? 'bg-gray-100' : '']
      }
    },
    onOptionSelected(index: number) {
      if (this.id) {
        localStorage.setItem(this.id, index.toString())
      }
    },
  },
})
</script>
