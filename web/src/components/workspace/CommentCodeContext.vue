<template>
  <div
    class="d2h-file-wrapper bg-white border border-gray-200 rounded-md my-4 z-0 relative overflow-y-hidden overflow-x-auto"
  >
    <div>
      <div class="d2h-code-wrapper">
        <table
          class="d2h-diff-table leading-4 z-0"
          style="border-collapse: separate; border-spacing: 0"
        >
          <tbody :class="['d2h-diff-tbody d2h-file-diff z-0']">
            <tr v-if="false" class="h-full overflow-hidden z-0">
              <td
                class="d2h-code-linenumber d2h-info h-full sticky left-0 z-20 bg-white min-w-[80px]"
              ></td>
              <td class="d2h-info h-full bg-blue-50 left-0 z-0 w-full">
                <div class="flex items-center sticky left-0">
                  <div class="d2h-code-line d2h-info text-gray-500">
                    &nbsp;&nbsp;{{ block.header }}
                  </div>
                </div>
              </td>
            </tr>

            <template v-for="(row, rowIndex) in contextLines" :key="rowIndex">
              <tr class="z-0">
                <td
                  :class="[
                    'd2h-code-linenumber bg-white sticky left-0 z-20 border-r border-l border-blue-500',
                  ]"
                >
                  <div class="select-none text-gray-600 flex">
                    <div class="line-num">
                      {{ context.contextStartsAtLine + rowIndex }}
                    </div>
                  </div>
                </td>

                <td :class="['code-row-wrapper relative z-10']">
                  <div class="d2h-code-line relative z-0 px-4">
                    <span class="d2h-code-line-ctn whitespace-pre">
                      {{ row }}
                    </span>
                  </div>
                </td>
              </tr>
            </template>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'

export default defineComponent({
  props: {
    context: { type: Object, required: true },
  },
  computed: {
    contextLines() {
      let l = this.context.context.split('\n')
      l.splice(-1)
      return l
    },
  },
})
</script>
