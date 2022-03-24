<template>
  <div class="w-max-full">
    <div class="mt-1 flex rounded-md shadow-sm w-max-full">
      <div class="relative flex items-stretch flex-grow focus-within:z-10">
        <input
          id="token"
          :value="val"
          type="text"
          readonly
          class="focus:ring-blue-500 focus:border-blue-500 block w-full rounded-none rounded-l-md sm:text-sm border-gray-300"
        />
      </div>
      <button
        class="-ml-px relative inline-flex items-center space-x-2 px-4 py-2 border border-gray-300 text-sm font-medium rounded-r-md text-gray-700 bg-gray-50 hover:bg-gray-100 focus:outline-none focus:ring-1 focus:ring-blue-500 focus:border-blue-500"
        @click="copyToClipboard"
      >
        <svg
          class="h-5 w-5 text-gray-400"
          xmlns="http://www.w3.org/2000/svg"
          fill="none"
          viewBox="0 0 24 24"
          stroke="currentColor"
        >
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M8 5H6a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2v-1M8 5a2 2 0 002 2h2a2 2 0 002-2M8 5a2 2 0 012-2h2a2 2 0 012 2m0 0h2a2 2 0 012 2v3m2 4H10m0 0l3-3m-3 3l3 3"
          />
        </svg>
        <span>Copy</span>
      </button>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { defineProps, toRef } from 'vue'

interface Props {
  value: string
}

const props = defineProps<Props>()
const val = toRef(props, 'value')

const copyToClipboard = function () {
  const copyText = document.getElementById('token')
  copyText.select()
  copyText.setSelectionRange(0, 99999)
  document.execCommand('copy')
}
</script>
