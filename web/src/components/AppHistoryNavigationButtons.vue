<template>
  <div class="flex gap-0.5">
    <button
      :disabled="!canGoBack"
      class="h-7 w-7 p-1 flex disabled:opacity-30 cursor-auto text-gray-500 disabled:bg-transparent disabled:text-gray-500 hover:bg-gray-300 hover:text-gray-600 rounded-md items-center justify-center"
      @click="back"
    >
      <ArrowLeftIcon />
    </button>
    <button
      :disabled="!canGoForward"
      class="h-7 w-7 p-1 flex disabled:opacity-30 cursor-auto text-gray-500 disabled:bg-transparent disabled:text-gray-500 hover:bg-gray-300 hover:text-gray-600 rounded-md items-center justify-center"
      @click="forward"
    >
      <ArrowRightIcon />
    </button>
  </div>
</template>

<script lang="ts">
import { defineComponent, PropType, ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import { ArrowRightIcon, ArrowLeftIcon } from '@heroicons/vue/solid'

export default defineComponent({
  components: { ArrowRightIcon, ArrowLeftIcon },
  props: {
    ipc: {
      type: Object,
      required: true,
    },
  },
  setup(props) {
    const canGoBack = ref(false)
    const canGoForward = ref(true)

    const route = useRoute()

    watch(route, () => {
      props.ipc.canGoBack().then((can) => {
        canGoBack.value = can
      })
      props.ipc.canGoForward().then((can) => {
        canGoForward.value = can
      })
    })

    props.ipc.canGoBack().then((can) => {
      canGoBack.value = can
    })
    props.ipc.canGoForward().then((can) => {
      canGoForward.value = can
    })

    return {
      canGoBack,
      canGoForward,
    }
  },
  methods: {
    async back() {
      await this.ipc.goBack()
      this.canGoForward = true
      this.canGoBack = await this.ipc.canGoBack()
    },
    async forward() {
      await this.ipc.goForward()
      this.canGoForward = await this.ipc.canGoForward()
      this.canGoBack = true
    },
  },
})
</script>
