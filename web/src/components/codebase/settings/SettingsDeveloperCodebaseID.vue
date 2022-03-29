<template>
  <div>
    <HorizontalDivider class="mt-4" bg="bg-white"
      >Sturdy CLI (Sturdy for Servers)</HorizontalDivider
    >
    <div class="mx-4 divide-gray-200">
      <p class="mt-1 text-sm text-gray-500">
        When using
        <a class="text-blue-500 underline" target="_blank" href="https://getsturdy.com/docs/cli"
          >the Sturdy CLI</a
        >, use this ID to identify this codebase:
        <span class="font-mono underline">{{ codebaseId }}</span>
      </p>
    </div>
  </div>
</template>

<script lang="ts" setup>
import HorizontalDivider from '../../../atoms/HorizontalDivider.vue'
import { defineProps } from 'vue'

interface Props {
  codebaseId: string
}

const props = defineProps<Props>()
</script>
