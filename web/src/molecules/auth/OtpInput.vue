<template>
  <div class="flex justify-center">
    <input autocomplete="off" name="hidden" type="text" class="hidden" />
    <SingleOtpInput
      v-for="(item, i) in numInputs"
      :key="i"
      :input-class="inputClasses"
      :focus="activeInput === i"
      :value="otp[i]"
      :is-last-child="i === numInputs - 1"
      :should-auto-focus="shouldAutoFocus"
      @change="handleOnChange"
      @keydown="handleOnKeyDown"
      @paste="handleOnPaste"
      @focus="handleOnFocus(i)"
      @blur="handleOnBlur"
    />
  </div>
</template>

<script lang="ts">
import { SingleOtpInput } from '../../atoms/auth'

const BACKSPACE = 8
const LEFT_ARROW = 37
const RIGHT_ARROW = 39
const DELETE = 46

export default {
  components: {
    SingleOtpInput,
  },
  props: {
    numInputs: {
      default: 6,
    },
    shouldAutoFocus: {
      default: true,
    },
    inputClasses: {
      type: String,
    },
  },
  emits: ['complete', 'change'],
  data() {
    return {
      activeInput: 0,
      otp: Array(this.numInputs).fill(''),
      oldOtp: Array(this.numInputs).fill(''),
    }
  },
  methods: {
    handleOnFocus(index: number) {
      this.activeInput = index
    },
    handleOnBlur() {
      this.activeInput = -1
    },
    checkFilledAllInputs() {
      if (this.otp.join('').length === this.numInputs) {
        return this.$emit('complete', this.otp.join(''))
      }
      return 'Wait until the user enters the required number of characters'
    },
    focusInput(index: number) {
      this.activeInput = Math.max(Math.min(this.numInputs - 1, index), 0)
    },
    focusNextInput() {
      this.focusInput(this.activeInput + 1)
    },
    focusPrevInput() {
      this.focusInput(this.activeInput - 1)
    },
    changeCodeAtFocus(value: number | string) {
      this.oldOtp = Object.assign([], this.otp)
      this.otp[this.activeInput] = value
      if (this.oldOtp.join('') !== this.otp.join('')) {
        this.$emit('change', this.otp.join(''))
        this.checkFilledAllInputs()
      }
    },
    handleOnPaste(event: any) {
      event.preventDefault()
      const pastedData = event.clipboardData
        .getData('text/plain')
        .replace('-', '')
        .slice(0, this.numInputs - this.activeInput)
        .split('')
      // Paste data from focused input onwards
      const currentCharsInOtp = this.otp.slice(0, this.activeInput)
      const combinedWithPastedData = currentCharsInOtp.concat(pastedData)
      combinedWithPastedData.slice(0, this.numInputs).forEach((value, i) => {
        this.otp[i] = value
      })
      this.focusInput(combinedWithPastedData.slice(0, this.numInputs).length)
      return this.checkFilledAllInputs()
    },
    handleOnChange(index: number) {
      this.changeCodeAtFocus(index)
      this.focusNextInput()
    },
    clearInput() {
      if (this.otp.length > 0) {
        this.$emit('change', '')
      }
      this.otp = []
      this.activeInput = 0
    },
    handleOnKeyDown(event: KeyboardEvent) {
      switch (event.keyCode) {
        case BACKSPACE:
          event.preventDefault()
          this.changeCodeAtFocus('')
          this.focusPrevInput()
          break
        case DELETE:
          event.preventDefault()
          this.changeCodeAtFocus('')
          break
        case LEFT_ARROW:
          event.preventDefault()
          this.focusPrevInput()
          break
        case RIGHT_ARROW:
          event.preventDefault()
          this.focusNextInput()
          break
        default:
          break
      }
    },
  },
}
</script>
