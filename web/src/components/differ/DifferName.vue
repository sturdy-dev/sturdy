<template>
  <span
    v-for="(n, index) in fileNameParts"
    :key="index"
    class="text-sm font-medium text-gray-500 hover:text-gray-700 cursor-pointer"
    @click="$emit('addWithPrefix', n.prefix, !added)"
    ><span v-if="index > 0">/</span>{{ n.name }}</span
  >
</template>

<script lang="js">
export default {
  name: "DifferName",
  props: ['name', 'added'],
  emits: ['addWithPrefix'],
  computed: {
    fileNameParts() {
      let parts = this.name.split("/")
      let res = [];
      for (let i = 0; i < parts.length; i++) {
        res.push({
          name: parts[i],
          prefix: parts.slice(0, i + 1).join("/") + (i === parts.length - 1 ? '' : '/')
        })
      }
      return res;
    },
  }
}
</script>
