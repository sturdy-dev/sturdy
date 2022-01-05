<template>
  <a
    :class="classes"
    :disabled="disabled"
    class="disabled:opacity-50 relative inline-flex items-center text-sm font-medium flex-shrink-0 hover:bg-gray-50 leading-5 focus:outline-none focus:ring-1 focus:ring-blue-500 focus:border-blue-500"
  >
    <slot>Button</slot>
  </a>
</template>

<script>
export default {
  name: 'LinkButton',
  props: {
    disabled: {
      type: Boolean,
      default: false,
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
    colors() {
      if (this.color === 'white') {
        return 'text-gray-700 bg-white hover:bg-gray-50 border-gray-300'
      }
      if (this.color === 'blue') {
        return 'text-white bg-blue-600 hover:bg-blue-700 border-transparent'
      }
      if (this.color === 'green') {
        return 'text-green-800 bg-green-100 hover:bg-green-200 border-transparent'
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
}
</script>
