<template>
  <button
    type="button"
    class="rounded-full w-3 h-3"
    :class="color"
    @click.prevent="() => $emit('click')"
  ></button>
</template>

<script>
export default {
  props: {
    server: {
      type: Object,
      required: true,
    },
  },
  emits: ['click'],
  computed: {
    color() {
      if (this.server.isUp == true) {
        return 'bg-green-500'
      } else if (this.server.isUp == false) {
        return 'bg-red-500'
      } else {
        return 'bg-gray-500'
      }
    },
  },
}
</script>
