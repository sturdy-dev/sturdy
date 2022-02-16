<template>
  <div
    style="height: calc(env(titlebar-area-height, 2rem) + 1px); -webkit-app-region: drag"
    class="flex flex-1 items-center"
  >
    <slot></slot>
  </div>
</template>

<script>
export default {}
</script>
