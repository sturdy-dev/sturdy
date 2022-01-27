<template>
  <div class="flex flex-row justify-center space-x-16">
    <div class="w-1/2">
      <slot name="leftIcon"></slot>
      <h1
        class="-ml-4 pl-4 border-l-2 border-amber-500/60 text-slate-50 font-semibold tracking-tight highlight"
      >
        {{ leftTitle }}
      </h1>
      <div class="mt-2">
        <p>{{ leftBody }}</p>
        <div
          class="cursor-pointer text-amber-500 hover:text-slate-50 font-semibold tracking-tight inline-flex group inline-flex items-center h-9"
        >
          <slot name="leftLink"></slot>
          <ChevronRightIcon class="h-4 w-4 ml-1" />
        </div>
      </div>
    </div>

    <div class="w-1/2">
      <slot name="rightIcon"></slot>
      <h1
        class="-ml-4 pl-4 border-l-2 border-amber-500/60 text-slate-50 font-semibold tracking-tight highlight"
      >
        {{ rightTitle }}
      </h1>
      <div class="mt-2">
        <p>{{ rightBody }}</p>

        <div
          class="cursor-pointer text-amber-500 hover:text-slate-50 font-semibold tracking-tight inline-flex group inline-flex items-center h-9"
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
