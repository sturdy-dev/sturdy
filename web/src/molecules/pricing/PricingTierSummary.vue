<template>
  <div class="bg-slate-800/90 rounded-md border-slate-500/30 border-2 flex flex-col min-w-[25rem]">
    <div class="p-6 h-[8rem]">
      <h1 class="text-4xl font-extrabold">{{ name }}</h1>
      <p v-if="price" class="font-bold mt-2 text-slate-400">
        <span class="text-3xl font-extrabold text-slate-200">${{ price }}</span>
        {{ ' ' }}
        <span class="text-base font-medium text-slate-400">/user/mo</span>
      </p>
    </div>

    <div class="p-6 flex flex-col flex-1 space-y-8">
      <ul class="flex flex-col flex-1 space-y-1">
        <slot></slot>
      </ul>

      <PricingButton class="flex-grow-0" :to="{ name: 'v2download' }"> Get started </PricingButton>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import PricingButton from './PricingButton.vue'

export default defineComponent({
  components: { PricingButton },
  props: {
    name: {
      type: String,
      required: true,
    },
    price: {
      type: String,
      required: false,
    },
  },
})
</script>
