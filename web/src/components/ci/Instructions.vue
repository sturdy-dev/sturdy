<template>
  <div v-if="description.length > 0">
    <p class="mt-1 max-w-2xl text-sm text-gray-500">{{ description }}</p>
  </div>

  <div v-if="instructions.length > 0" class="mt-5 border-t border-gray-200">
    <dl
      v-for="instruction in instructions"
      :key="instruction.name"
      class="sm:divide-y sm:divide-gray-200"
    >
      <div class="py-4 sm:py-5 sm:grid sm:grid-cols-3 sm:gap-4">
        <dt class="text-sm font-medium text-gray-500">{{ instruction.name }}</dt>
        <dd
          class="mt-1 text-sm text-gray-900 sm:mt-0 sm:col-span-2"
          :class="[instruction.pre ? 'font-mono' : '']"
        >
          {{ instruction.value }}
        </dd>
      </div>
    </dl>
  </div>
</template>

<script lang="ts">
import type { PropType } from 'vue'

export type Instruction = {
  name: string
  value: string
  pre?: boolean
}

export default {
  props: {
    description: {
      type: String,
      default: '',
    },
    instructions: {
      type: Array as PropType<Instruction[]>,
      default: () => [],
    },
  },
}
</script>
