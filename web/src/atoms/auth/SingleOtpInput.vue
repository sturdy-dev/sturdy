<template>
  <div class="flex items-center">
    <input
      ref="input"
      v-model="model"
      type="text"
      :class="inputClass"
      maxlength="1"
      pattern="[a-zA-Z0-9]"
      @input="handleOnChange"
      @keydown="handleOnKeyDown"
      @paste="handleOnPaste"
      @focus="handleOnFocus"
      @blur="handleOnBlur"
    />
  </div>
</template>

<script lang="ts">
import { defineComponent, ref } from 'vue'
import type { Ref } from 'vue'

const isCodeNumeric = (charCode: number) =>
  (charCode >= 48 && charCode <= 57) || (charCode >= 96 && charCode <= 105)

const isCodeAlpha = (charCode: number) =>
  (charCode >= 65 && charCode <= 90) || (charCode >= 97 && charCode <= 122)

export default defineComponent({
  props: {
    value: {
      type: String,
      required: false,
      default: null,
    },
    inputClass: {
      type: String,
      required: false,
      default: '',
    },
    focus: {
      type: Boolean,
    },
    shouldAutoFocus: {
      type: Boolean,
      default: true,
    },
    isLastChild: {
      type: Boolean,
    },
  },
  emits: ['change', 'keydown', 'paste', 'focus', 'blur'],
  setup(props) {
    const model = ref(props.value || '')
    const input = ref<HTMLInputElement | null>(null) as Ref<HTMLInputElement>

    return {
      input,
      model,
    }
  },
  watch: {
    value: function (newValue: string) {
      this.model = newValue
    },
    focus: function (newFocusValue, oldFocusValue) {
      // Check if focusedInput changed
      // Prevent calling function if input already in focus
      if (oldFocusValue !== newFocusValue && this.input && this.focus) {
        this.input.focus()
        this.input.select()
      }
    },
  },
  mounted() {
    if (this.input.value && this.focus && this.shouldAutoFocus) {
      this.input.focus()
      this.input.select()
    }
  },
  methods: {
    handleOnChange() {
      return this.$emit('change', this.model)
    },
    handleOnKeyDown(event: KeyboardEvent) {
      // Only allow characters a-zA-Z0-9, DEL, Backspace, Enter, Right and Left Arrows, and Pasting
      const keyEvent = event || window.event
      const charCode = keyEvent.which ? keyEvent.which : keyEvent.keyCode
      if (
        isCodeAlpha(charCode) ||
        isCodeNumeric(charCode) ||
        [8, 9, 13, 37, 39, 46, 86].includes(charCode)
      ) {
        this.$emit('keydown', event)
      } else {
        keyEvent.preventDefault()
      }
    },
    handleOnPaste(event: KeyboardEvent) {
      this.$emit('paste', event)
    },
    handleOnFocus() {
      this.input.select()
      return this.$emit('focus')
    },
    handleOnBlur() {
      this.$emit('blur')
    },
  },
})
</script>
