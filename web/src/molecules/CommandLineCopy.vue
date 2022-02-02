<template>
  <div
    class="h-60 border border-slate-300/30 relative overflow-hidden shadow-xl flex bg-slate-700/50 rounded-xl"
  >
    <div class="relative w-full flex flex-col">
      <div class="flex-none border-b border-slate-300/30">
        <div class="flex items-center h-8 space-x-1.5 px-3">
          <div class="w-2.5 h-2.5 bg-red-700 rounded-full"></div>
          <div class="w-2.5 h-2.5 bg-yellow-600 rounded-full"></div>
          <div class="w-2.5 h-2.5 bg-green-600 rounded-full"></div>
        </div>

        <div class="-mt-6 mb-1 text-sm font-semibold text-slate-400 text-center select-none">
          Bash
        </div>
      </div>

      <div class="relative min-h-0 flex-auto flex flex-col">
        <div class="w-full flex-auto flex min-h-0">
          <div class="w-full flex-auto flex min-h-0 overflow-auto">
            <div class="w-full relative flex-auto">
              <div class="flex flex-col min-h-full text-sm tracking-tight leading-2">
                <code class="text-left flex relative block text-slate-50 overflow-auto p-4">
                  <span class="text-slate-400 select-none mr-3">$</span>
                  <span class="whitespace-pre">{{ content }}</span>
                </code>
                <div class="flex flex-row-reverse">
                  <DuplicateIcon
                    class="cursor-pointer m-4 w-8 h-8 hover:text-purple-400"
                    @click="copy"
                  />
                </div>
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
import { DuplicateIcon } from '@heroicons/vue/outline'

export default defineComponent({
  components: {
    DuplicateIcon,
  },
  props: {
    content: {
      type: String,
      required: true,
    },
  },
  methods: {
    copy() {
      navigator.clipboard.writeText(this.content)
    },
  },
})
</script>
