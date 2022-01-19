<template>
  <div>
    <label :for="name" class="block text-sm font-medium text-gray-700">
      {{  label }}
    </label>
    <div class="mt-1">
      <TextInput v-model="value" :placeholder="placeholder" :name="name" :id="name" />
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import TextInput from "./TextInput.vue";
export default defineComponent({
  props: {
    modelValue: {
      type: String,
      required: true,
    },
    label: {
      type: String,
      required: true,
    },
    placeholder: {
      type: String,
      required: true,
    },
    name: {
      type: String,
      required: true,
    }
  },
  emits: ['update:modelValue'],
  components: {TextInput},
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
