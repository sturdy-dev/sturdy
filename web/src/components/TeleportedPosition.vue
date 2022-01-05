<template>
  <div ref="positioner"></div>
  <teleport to="#teleported-position">
    <div ref="container" class="absolute left-0 top-0">
      <slot></slot>
    </div>
  </teleport>
</template>

<script lang="ts">
import { defineComponent, PropType } from 'vue'

export default defineComponent({
  name: 'TeleportedPosition',
  props: {
    traceRef: {
      type: undefined as PropType<unknown>,
      default: null,
    },
  },
  watch: {
    traceRef() {
      this.timeout = setTimeout(this.updatePosition, 400)
    },
  },
  mounted() {
    let el: HTMLElement | null = this.$refs.positioner.parentElement
    const parentElements: HTMLElement[] = []

    while (el != null) {
      parentElements.push(el)
      el.addEventListener('scroll', this.updatePosition)
      el = el.parentElement
    }
    window.addEventListener('scroll', this.updatePosition)

    this.parentElements = parentElements
    this.updatePosition()
  },
  beforeUnmount() {
    clearTimeout(this.timeout)
    for (const el of this.parentElements) {
      el.removeEventListener('scroll', this.updatePosition)
    }
    window.removeEventListener('scroll', this.updatePosition)
  },
  methods: {
    updatePosition() {
      let { x, y } = this.$refs.positioner.getBoundingClientRect()
      if (y > window.innerHeight) {
        y = -1000
      }
      Object.assign(this.$refs.container.style, {
        left: `${x}px`,
        top: `${y + window.pageYOffset}px`,
      })
    },
  },
})
</script>
