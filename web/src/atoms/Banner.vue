<template>
  <div
    class="p-4 w-full border-l-4"
    :class="{
      'bg-green-50 border-green-400': status === 'success',
      'bg-blue-50 border-blue-400': status === 'info',
      'bg-red-50 border-red-400': status === 'error',
      'bg-yellow-50 border-yellow-400': status === 'warning',
    }"
  >
    <div class="flex w-full">
      <div v-if="showIcon" class="flex-shrink-0 mr-3">
        <CheckCircleIcon
          v-if="status === 'success'"
          class="h-5 w-5 text-green-400"
          aria-hidden="true"
        />
        <InformationCircleIcon
          v-else-if="status === 'info'"
          class="h-5 w-5 text-blue-400"
          aria-hidden="true"
        />
        <XCircleIcon
          v-else-if="status === 'error'"
          class="h-5 w-5 text-red-400"
          aria-hidden="true"
        />
        <ExclamationIcon
          v-else-if="status === 'warning'"
          class="h-5 w-5 text-yellow-400"
          aria-hidden="true"
        />
      </div>
      <div
        class="w-full text-sm font-medium"
        :class="{
          'text-green-800': status === 'success',
          'text-blue-800': status === 'info',
          'text-red-800': status === 'error',
        }"
      >
        <p v-if="message" class="">
          <span>{{ message }}</span>
        </p>
        <slot v-else> This banner has no message! Oops! </slot>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import {
  ExclamationIcon,
  CheckCircleIcon,
  InformationCircleIcon,
  XCircleIcon,
} from '@heroicons/vue/solid'
import { defineComponent, type PropType } from 'vue'

export default defineComponent({
  name: 'Banner',
  components: { ExclamationIcon, CheckCircleIcon, InformationCircleIcon, XCircleIcon },
  props: {
    status: {
      type: String as PropType<'success' | 'warning' | 'info' | 'error'>,
      default: 'success',
      required: true,
    },
    message: {
      type: String,
      default: '',
    },
    showIcon: {
      type: Boolean,
      default: true,
    },
  },
})
</script>
