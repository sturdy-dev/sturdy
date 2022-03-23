<template>
  <div :class="[!isLast ? 'pb-10' : '', 'relative']">
    <StepIndicator :status="status" :is-last="isLast" :name="name">
      <slot v-if="status !== 'pending'" />
    </StepIndicator>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import type { PropType } from 'vue'
import StepIndicator from './StepIndicator.vue'
import type { Status } from './StepIndicator.vue'

export default defineComponent({
  components: { StepIndicator },
  props: {
    name: { type: String, required: true },
    status: { type: String as PropType<Status>, required: true },
    isLast: { type: Boolean, default: false },
  },
})
</script>
