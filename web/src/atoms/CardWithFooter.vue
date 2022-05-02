<template>
  <div class="bg-white overflow-hidden drop-shadow-lg rounded-lg">
    <div class="px-4 py-5 sm:p-6">
      <slot name="default"></slot>
    </div>
    <div class="bg-gray-50 px-4 py-4 sm:px-6">
      <slot name="footer"></slot>
    </div>
  </div>
</template>

<script lang="ts" setup></script>
