import { defineComponent, PropType, VNode } from 'vue'

export default defineComponent({
  props: {
    nodes: {
      type: Array as PropType<VNode[]>,
      required: true,
    },
  },
  render(_: unknown, __: unknown, { nodes }: { nodes: VNode[] }) {
    return nodes
  },
})
