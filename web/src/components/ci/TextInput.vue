<template>
  <input
    v-model="value"
    type="text"
    autocomplete="off"
    class="focus:ring-blue-500 focus:border-blue-500 flex-grow block w-full min-w-0 rounded-md sm:text-sm border-gray-300"
  />
</template>

<script lang="ts">
import { defineComponent } from 'vue'
export default defineComponent({
  props: {
    modelValue: {
      type: String,
      required: true,
    },
  },
  emits: ['update:modelValue'],
  computed: {
    value: {
      get() {
        return this.modelValue
      },
      set(value: string) {
        this.$emit('update:modelValue', value)
      },
    },
  },
})
</script>
