<template>
  <Listbox as="div" v-model="selected" class="w-32">
    <ListboxLabel class="block text-sm font-medium text-gray-700"> Channel </ListboxLabel>
    <div class="mt-1 relative">
      <ListboxButton
        class="bg-white relative w-full border border-gray-300 rounded-md shadow-sm pl-3 pr-10 py-2 text-left cursor-default focus:outline-none focus:ring-1 focus:ring-yellow-500 focus:border-indigo-500 sm:text-sm"
      >
        <span class="block truncate">{{ selected.name }}</span>
        <span class="absolute inset-y-0 right-0 flex items-center pr-2 pointer-events-none">
          <SelectorIcon class="h-5 w-5 text-gray-400" aria-hidden="true" />
        </span>
      </ListboxButton>

      <transition
        leave-active-class="transition ease-in duration-100"
        leave-from-class="opacity-100"
        leave-to-class="opacity-0"
      >
        <ListboxOptions
          class="absolute z-10 mt-1 w-full bg-white shadow-lg max-h-60 rounded-md py-1 text-base ring-1 ring-black ring-opacity-5 overflow-auto focus:outline-none sm:text-sm"
        >
          <ListboxOption
            as="template"
            v-for="channel in channels"
            :key="channel.id"
            :value="channel"
            v-slot="{ active, selected }"
          >
            <li
              :class="[
                active ? 'text-white bg-yellow-600' : 'text-gray-900',
                'cursor-default select-none relative py-2 pl-3 pr-9',
              ]"
            >
              <span :class="[selected ? 'font-semibold' : 'font-normal', 'block truncate']">
                {{ channel.name }}
              </span>

              <span
                v-if="selected"
                :class="[
                  active ? 'text-white' : 'text-yellow-600',
                  'absolute inset-y-0 right-0 flex items-center pr-4',
                ]"
              >
                <CheckIcon class="h-5 w-5" aria-hidden="true" />
              </span>
            </li>
          </ListboxOption>
        </ListboxOptions>
      </transition>
    </div>
  </Listbox>
</template>

<script>
import { ref } from 'vue'
import {
  Listbox,
  ListboxButton,
  ListboxLabel,
  ListboxOption,
  ListboxOptions,
} from '@headlessui/vue'
import { CheckIcon, SelectorIcon } from '@heroicons/vue/solid'
import ipc from '../ipc'

const channels = [
  { id: 'latest', name: 'Stable' },
  { id: 'beta', name: 'Beta' },
  { id: 'alpha', name: 'Alpha' },
]

export default {
  components: {
    Listbox,
    ListboxButton,
    ListboxLabel,
    ListboxOption,
    ListboxOptions,
    CheckIcon,
    SelectorIcon,
  },
  setup() {
    const selected = ref(channels[1])
    ipc.getChannel().then((channel) => {
      selected.value = channels.find((c) => c.id === channel)
    })
    return {
      channels,
      selected,
    }
  },
  watch: {
    selected: function (channel) {
      ipc.setChannel(channel.id)
    },
  },
}
</script>
