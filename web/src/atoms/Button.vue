<template>
  <button
    :type="buttonType"
    :class="classes"
    :disabled="disabled"
    class="disabled:opacity-50 relative inline-flex items-center text-sm font-medium flex-shrink-0 hover:bg-gray-50 focus:z-10 focus:outline-none focus:ring-1 focus:ring-blue-500 focus:border-blue-500 group leading-5 gap-2"
  >
    <Spinner v-if="spinner" />
    <component :is="icon" v-else-if="icon" class="h-5 w-5" :class="iconColor" />

    <slot>Button</slot>

    <div
      v-if="showTooltip"
      class="absolute bottom-0 flex-col mb-8 hidden group-hover:flex z-50"
      :class="{
        'bottom-full -translate-y-1 group-hover:-translate-y-2': tooltipPosition === 'top',
        'top-full translate-y-1 group-hover:translate-y-2': tooltipPosition === 'bottom',
        'right-0': tooltipPosition === 'left',
        'left-0': tooltipPosition === 'right' || tooltipRight,
      }"
    >
      <span
        class="relative p-2 text-xs leading-none rounded text-white whitespace-nowrap bg-black shadow-lg"
      >
        <slot name="tooltip" />
      </span>
      <div
        class="w-3 h-3 transform rotate-45 bg-black absolute bottom-0 -mb-1"
        :class="{
          'top-full -translate-y-2': tooltipPosition === 'top',
          'bottom-full translate-y-2': tooltipPosition === 'bottom',
          'left-2': tooltipPosition === 'left',
          'right-2': tooltipPosition === 'right' || tooltipRight,
        }"
      />
    </div>
  </button>
</template>

<script lang="ts">
import { defineComponent, type FunctionalComponent, type PropType } from 'vue'
import Spinner from './Spinner.vue'

export default defineComponent({
  components: { Spinner },
  props: {
    disabled: {
      type: Boolean,
      default: false,
    },
    buttonType: {
      type: String as PropType<'button' | 'submit' | 'reset'>,
      default: 'button',
    },
    color: {
      type: String,
      default: 'white',
    },
    size: {
      type: String,
      default: 'normal',
    },
    grouped: {
      type: Boolean,
      default: false,
    },
    first: {
      type: Boolean,
      default: false,
    },
    last: {
      type: Boolean,
      default: false,
    },
    showTooltip: {
      type: Boolean,
      default: false,
    },
    // deprecated, use tooltipPosition instead
    tooltipRight: {
      type: Boolean,
      default: false,
    },
    tooltipPosition: {
      type: String as PropType<'left' | 'right' | 'top' | 'bottom'>,
      default: 'top',
    },
    autoFocus: {
      type: Boolean,
      default: false,
    },
    icon: {
      type: Function as PropType<FunctionalComponent>,
      required: false,
      default: undefined,
    },
    spinner: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    classes() {
      return [this.rounded + ' ' + this.colors + ' ' + this.border + ' ' + this.margins]
    },
    rounded() {
      if (this.grouped && this.first) {
        return 'rounded-l-md'
      } else if (this.grouped && this.last) {
        return 'rounded-r-md'
      } else if (!this.grouped) {
        return 'rounded-md'
      }
      return ''
    },
    border() {
      if (this.grouped && !this.last) {
        return 'border border-r-0'
      }
      return 'border'
    },
    iconColor() {
      if (this.color === 'red') {
        return 'text-red-400'
      }
      if (this.color === 'blue') {
        return 'text-white'
      }
      return 'text-gray-400'
    },
    colors() {
      if (this.color === 'white') {
        return 'text-gray-700 bg-white hover:bg-gray-50 border-gray-300'
      }
      if (this.color === 'blue') {
        return 'text-white bg-blue-600 hover:bg-blue-700 border-blue-700'
      }
      if (this.color === 'lightblue') {
        return 'text-gray-800 bg-blue-200 hover:bg-blue-300 border-blue-300'
      }
      if (this.color === 'green') {
        return 'text-green-800 bg-green-100 hover:bg-green-200 border-green-200'
      }
      if (this.color === 'red') {
        return 'text-red-800 bg-red-100 hover:bg-red-200 border-red-200'
      }
      if (this.color === 'slate') {
        return 'bg-slate-800 border-slate-200 font-semibold hover:bg-gray-900'
      }
      return ''
    },
    margins() {
      if (this.size === 'wider') {
        return 'px-4 py-2 leading-4'
      }
      if (this.size === 'normal') {
        return 'px-3 py-2 leading-4'
      }
      if (this.size === 'small') {
        return 'px-4 py-1'
      }

      if (this.size === 'tiny') {
        return 'px-3 py-0.5'
      }

      return ''
    },
  },
})
</script>
