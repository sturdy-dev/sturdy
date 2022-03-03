<template>
  <div
    class="border border-slate-300/30 relative overflow-hidden shadow-xl flex bg-slate-700/50 rounded-xl"
  >
    <div class="relative w-full flex flex-col">
      <div class="flex-none border-b border-slate-300/30">
        <div class="flex items-center h-8 space-x-1.5 px-3">
          <div class="w-2.5 h-2.5 bg-red-700 rounded-full"></div>
          <div class="w-2.5 h-2.5 bg-yellow-600 rounded-full"></div>
          <div class="w-2.5 h-2.5 bg-green-600 rounded-full"></div>
        </div>

        <div class="-mt-6 mb-1 text-sm font-semibold text-slate-400 text-center select-none">
          {{ windowTitle }}
        </div>
      </div>

      <div class="relative min-h-0 flex-auto flex flex-col">
        <div class="w-full flex-auto flex min-h-0">
          <div class="w-full flex-auto flex min-h-0 overflow-auto">
            <div class="w-full relative flex-auto">
              <div class="flex flex-col min-h-full text-sm tracking-tight leading-2">
                <code class="flex justify-center relative block text-slate-50 overflow-auto p-4">
                  <slot></slot>
                </code>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'

export default defineComponent({
  props: {
    windowTitle: {
      type: String,
      required: true,
    },
  },
})
</script>
