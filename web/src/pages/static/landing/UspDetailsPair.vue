<template>
  <div class="flex flex-row justify-center space-x-16">
    <div class="w-1/2">
      <div class="flex flex-row space-x-4 items-center">
        <div class="w-8 h-8">
          <slot name="leftIcon"></slot>
        </div>
        <h1 class="text-slate-50 font-semibold tracking-tight highlight">
          {{ leftTitle }}
        </h1>
      </div>

      <div class="mt-2">
        <p>{{ leftBody }}</p>
        <div
          class="mt-2 text-sm cursor-pointer text-slate-400/80 hover:text-slate-50/80 font-semibold tracking-tight inline-flex group inline-flex items-center h-9"
        >
          <slot name="leftLink"></slot>
          <ChevronRightIcon class="h-4 w-4 ml-1" />
        </div>
      </div>
    </div>

    <div class="w-1/2">
      <div class="flex flex-row space-x-4 items-center">
        <div class="w-8 h-8">
          <slot name="rightIcon"></slot>
        </div>
        <h1 class="text-slate-50 font-semibold tracking-tight highlight">
          {{ rightTitle }}
        </h1>
      </div>
      <div class="mt-2">
        <p>{{ rightBody }}</p>

        <div
          class="mt-2 text-sm cursor-pointer text-slate-400/80 hover:text-slate-50/80 font-semibold tracking-tight inline-flex group inline-flex items-center h-9"
        >
          <slot name="rightLink"></slot>
          <ChevronRightIcon class="h-4 w-4 ml-1" />
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import { ChevronRightIcon } from '@heroicons/vue/outline'

export default defineComponent({
  components: {
    ChevronRightIcon,
  },
  props: {
    leftTitle: {
      type: String,
      required: true,
    },
    leftBody: {
      type: String,
      required: true,
    },
    rightTitle: {
      type: String,
      required: true,
    },
    rightBody: {
      type: String,
      required: true,
    },
  },
})
</script>

<style scoped></style>
