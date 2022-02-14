<template>
  <button
    type="button"
    class="rounded-full w-3 h-3"
    :class="color"
    @click.prevent="handleTest"
  ></button>
</template>

<script>
import ipc from '../ipc'
import { ref, toRefs } from 'vue'

export default {
  props: {
    server: {
      type: Object,
      required: true,
    },
  },
  setup(props) {
    const { server } = toRefs(props)
    const status = ref(undefined)
    ipc.isHostUp(server.value).then((isUp) => {
      status.value = isUp
    })
    return {
      status,
    }
  },
  computed: {
    color() {
      if (this.status == true) {
        return 'bg-green-500'
      } else if (this.status == false) {
        return 'bg-red-500'
      } else {
        return 'bg-gray-500'
      }
    },
  },
  methods: {
    handleTest() {
      ipc.isHostUp(this.server).then((status) => {
        this.status = status
      })
    },
  },
}
</script>
