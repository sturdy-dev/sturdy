<template>
  <div>
    <div
      class="flex flex-col w-[29rem] h-56 text-slate-300/80 bg-[#0c1322] border border-slate-700/30 px-8 py-4 mb-2"
    >
      <div class="w-28 h-6 brightness-75">
        <slot name="logo"> </slot>
      </div>
      <div class="grow relative tracking-tight text-sm mt-5">
        <svg
          class="absolute text-amber-500/50 top-0 left-0 transform -translate-x-7 -translate-y-4 h-7 w-7 text-indigo-600"
          fill="currentColor"
          viewBox="0 0 32 32"
          aria-hidden="true"
        >
          <path
            d="M9.352 4C4.456 7.456 1 13.12 1 19.36c0 5.088 3.072 8.064 6.624 8.064 3.36 0 5.856-2.688 5.856-5.856 0-3.168-2.208-5.472-5.088-5.472-.576 0-1.344.096-1.536.192.48-3.264 3.552-7.104 6.624-9.024L9.352 4zm16.512 0c-4.8 3.456-8.256 9.12-8.256 15.36 0 5.088 3.072 8.064 6.624 8.064 3.264 0 5.856-2.688 5.856-5.856 0-3.168-2.304-5.472-5.184-5.472-.576 0-1.248.096-1.44.192.48-3.264 3.456-7.104 6.528-9.024L25.864 4z"
          />
        </svg>
        <p class="relative">
          <slot name="quote"></slot>
        </p>
      </div>
      <div class="flex text-sm items-center gap-4 justify-end mt-2">
        <slot class="" name="face"></slot>
        <div class="flex flex-col">
          <span class="font-semibold">{{ name }}</span>
          <span class="text-slate-300/70">{{ title }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'

export default defineComponent({
  props: {
    name: {
      type: String,
      required: true,
    },
    title: {
      type: String,
      required: false,
    },
  },
})
</script>

<style scoped></style>
