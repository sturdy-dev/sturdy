<template>
  <SwitchGroup as="div" class="flex items-center">
    <Switch
      v-model="modelValue"
      class="flex-shrink-0 group relative rounded-full inline-flex items-center justify-center h-5 w-10 cursor-pointer focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-yellow-500"
    >
      <span class="sr-only">Use setting</span>
      <span
        aria-hidden="true"
        class="pointer-events-none absolute bg-white w-full h-full rounded-md"
      />
      <span
        aria-hidden="true"
        :class="[
          modelValue ? 'bg-yellow-600' : 'bg-gray-200',
          'pointer-events-none absolute h-4 w-9 mx-auto rounded-full transition-colors ease-in-out duration-200',
        ]"
      />
      <span
        aria-hidden="true"
        :class="[
          modelValue ? 'translate-x-5' : 'translate-x-0',
          'pointer-events-none absolute left-0 inline-block h-5 w-5 border border-gray-200 rounded-full bg-white shadow transform ring-0 transition-transform ease-in-out duration-200',
        ]"
      />
    </Switch>
    <SwitchLabel v-if="label" as="span" class="ml-3">
      <span class="text-sm font-medium text-gray-900">{{ label }}</span>
    </SwitchLabel>
  </SwitchGroup>
</template>

<script>
import { Switch, SwitchGroup, SwitchLabel } from '@headlessui/vue'

export default {
  props: {
    modelValue: { type: Boolean, default: false },
    label: {
      type: String,
      required: false,
    },
  },
  watch: {
    modelValue: {
      immediate: true,
      handler(value) {
        this.$emit('update:modelValue', value)
      },
    },
  },
  emits: ['update:modelValue'],
  components: {
    Switch,
    SwitchGroup,
    SwitchLabel,
  },
}
</script>
