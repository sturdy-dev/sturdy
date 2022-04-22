<template>
  <TextareaMentions
    ref="textarea"
    v-model="val"
    :style="styles"
    :user="user"
    :members="members"
    @focus="resize"
  ></TextareaMentions>
</template>

<script lang="ts">
import TextareaMentions, { MEMBER_FRAGMENT } from './TextareaMentions.vue'
import { computed, defineComponent, nextTick, type PropType, ref } from 'vue'
import type { UserFragment } from './__generated__/TextareaMentions'
import type { MemberFragment } from './__generated__/TextareaMentions'

export { MEMBER_FRAGMENT }

export default defineComponent({
  components: { TextareaMentions },
  props: {
    members: {
      type: Array as PropType<Array<MemberFragment>>,
      required: true,
    },
    user: {
      type: Object as PropType<UserFragment>,
      required: false,
      default: null,
    },
    modelValue: {
      type: [String, Number],
      default: '',
    },
    autosize: {
      type: Boolean,
      default: true,
    },
    minHeight: {
      type: [Number],
      default: null,
    },
    maxHeight: {
      type: [Number],
      default: null,
    },
    /*
     * Force !important for style properties
     */
    important: {
      type: [Boolean, Array],
      default: false,
    },
  },
  emits: ['update:modelValue'],
  setup(props) {
    const textarea = ref<InstanceType<typeof TextareaMentions>>()

    const height = ref<string>()
    const maxHeightScroll = ref<boolean>(false)
    const val = ref<string | number | undefined>(undefined)

    const isResizeImportant = computed(() => {
      const imp = props.important
      return imp === true || (Array.isArray(imp) && imp.includes('resize'))
    })

    const isHeightImportant = computed(() => {
      const imp = props.important
      return imp === true || (Array.isArray(imp) && imp.includes('height'))
    })

    const isOverflowImportant = computed(() => {
      const imp = props.important
      return imp === true || (Array.isArray(imp) && imp.includes('overflow'))
    })

    const resize = () => {
      nextTick(() => {
        if (!textarea.value) {
          return
        }

        let contentHeight = textarea.value.$el.scrollHeight + 1
        if (props.minHeight) {
          contentHeight = contentHeight < props.minHeight ? props.minHeight : contentHeight
        }
        if (props.maxHeight) {
          if (contentHeight > props.maxHeight) {
            contentHeight = props.maxHeight
            maxHeightScroll.value = true
          } else {
            maxHeightScroll.value = false
          }
        }

        const heightVal = contentHeight + 'px'
        height.value = `${heightVal}${isHeightImportant.value ? ' !important' : ''}`
      })
    }

    const styles = computed(() => {
      if (!props.autosize) return {}
      return {
        resize: !isResizeImportant.value ? 'none' : 'none !important',
        height: height.value,
        overflow: maxHeightScroll.value
          ? 'auto'
          : !isOverflowImportant.value
          ? 'hidden'
          : 'hidden !important',
      }
    })

    return {
      textarea,
      height,
      maxHeightScroll,
      val,
      resize,
      styles,
    }
  },
  watch: {
    modelValue(val) {
      // <- copy this
      this.val = val
    },
    val(val) {
      this.$nextTick(this.resize)
      this.$emit('update:modelValue', val)
    },
    minHeight() {
      this.$nextTick(this.resize)
    },
    maxHeight() {
      this.$nextTick(this.resize)
    },
    autosize(val) {
      if (val) this.resize()
    },
  },
  mounted() {
    this.val = this.modelValue
    this.resize()
  },
})
</script>
