<template>
  <div ref="target" class="contents">
    <slot></slot>
  </div>
</template>

<script lang="ts">
import { defineComponent, PropType } from 'vue'
import { useOnboarding } from './Onboarding.vue'
import { Step } from './Step'

export default defineComponent({
  props: {
    id: {
      type: String,
      required: true,
    },
    dependencies: {
      type: Array as PropType<string[]>,
      default: () => [],
    },
    enabled: {
      type: Boolean,
      default: true,
    },
  },

  setup() {
    const onboarding = useOnboarding()
    return {
      onboarding,
    }
  },

  watch: {
    enabled() {
      this.attemptRegisterStep()
    },
  },

  mounted() {
    this.attemptRegisterStep()
  },

  unmounted() {
    if (this.step) {
      this.onboarding.unregisterStep(this.step)
    }
  },

  methods: {
    attemptRegisterStep() {
      if (!this.$props.enabled) {
        return
      }

      if (this.$refs.target == null) {
        return
      }

      const step: Step = {
        highlightedElement: this.$refs.target,
        title: this.$slots.title,
        description: this.$slots.description,
        id: this.$props.id,
        dependencies: new Set(this.$props.dependencies),
      }
      this.step = step

      this.$nextTick(() => {
        if (step.highlightedElement.childElementCount > 0) {
          if (this.onboarding.registerStep(step)) {
            return
          }
        }
        const mo = new MutationObserver(() => {
          if (step.highlightedElement.childElementCount > 0) {
            if (this.onboarding.registerStep(step)) {
              mo.disconnect()
            }
          }
        })
        mo.observe(step.highlightedElement, { childList: true })
      })
    },
  },
})
</script>
