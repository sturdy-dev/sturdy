<template>
  <!-- Line -->
  <div class="flex">
    <div
      v-if="!isLast && status === 'current'"
      class="-ml-px absolute mt-0.5 top-4 left-4 w-0.5 h-full bg-gray-300"
      aria-hidden="true"
    />
    <div
      v-else-if="!isLast && status === 'completed'"
      class="-ml-px absolute mt-0.5 top-4 left-4 w-0.5 h-full bg-blue-600"
      aria-hidden="true"
    />
    <div
      v-else-if="!isLast"
      class="-ml-px absolute mt-0.5 top-4 left-4 w-0.5 h-full bg-gray-300"
      aria-hidden="true"
    />

    <!-- Circle -->
    <div
      class="relative flex items-start group space-x-4 flex-col flex-1 max-w-full"
      aria-current="step"
    >
      <div class="flex items-center space-x-4">
        <span class="h-9 block" aria-hidden="true">
          <span
            v-if="status === 'current'"
            class="relative z-10 w-8 h-8 flex items-center justify-center bg-white border-2 border-blue-600 rounded-full"
          >
            <span class="h-2.5 w-2.5 bg-blue-600 rounded-full" />
          </span>

          <span
            v-else-if="status === 'completed'"
            class="relative z-10 w-8 h-8 flex items-center justify-center bg-blue-600 rounded-full group-hover:bg-blue-800"
          >
            <CheckIcon class="w-5 h-5 text-white" aria-hidden="true" />
          </span>

          <span
            v-else
            class="relative z-10 w-8 h-8 flex items-center justify-center bg-white border-2 border-gray-300 rounded-full group-hover:border-gray-400"
          >
            <span class="h-2.5 w-2.5 bg-transparent rounded-full group-hover:bg-gray-300" />
          </span>
        </span>

        <div>
          <h3 class="text-lg leading-6 font-medium text-gray-900">{{ name }}</h3>
        </div>
      </div>

      <div class="pl-8 pr-8 flex-1 w-full">
        <slot></slot>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import type { PropType } from 'vue'
import { CheckIcon } from '@heroicons/vue/solid'

export type Status = 'pending' | 'current' | 'completed'

export default defineComponent({
  components: {
    CheckIcon,
  },
  props: {
    isLast: {
      type: Boolean,
      default: false,
    },
    status: {
      type: String as PropType<Status>,
      required: true,
    },
    name: {
      type: String,
      required: true,
    },
  },
})
</script>
