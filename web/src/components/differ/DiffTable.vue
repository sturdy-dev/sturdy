<template>
  <div class="d2h-code-wrapper">
    <table
      v-if="parsedDiff"
      class="d2h-diff-table leading-4 z-0"
      style="border-collapse: separate; border-spacing: 0"
    >
      <tbody
        v-for="(hunk, hunkIndex) in parsedDiff"
        :key="hunkIndex"
        class="d2h-diff-tbody d2h-file-diff z-0"
      >
        <template
          v-for="block in highlightedBlocks(hunk.blocks, hunk.language)"
          :key="block.header"
        >
          <tr class="h-full overflow-hidden z-0 leading-7">
            <td class="d2h-code-linenumber d2h-info h-full sticky left-0 z-20 bg-white">
              <div class="font-sans text-sm p-1">
                <slot name="blockIndexAction" />
              </div>
            </td>
            <td class="d2h-info h-full bg-blue-50 left-0 z-0">
              <div class="flex items-center sticky left-0">
                <div class="d2h-code-line d2h-info text-gray-500">
                  {{ block.header }}
                </div>
              </div>
            </td>
          </tr>

          <template v-for="(row, rowIndex) in block.lines" :key="rowIndex">
            <tr :data-row-index="rowIndex" class="z-0" :class="[grayedOut ? 'opacity-70' : '']">
              <td
                :class="[
                  'd2h-code-linenumber bg-white sticky left-0 z-20',
                  row.type === 'insert' ? 'border-r border-l border-green-500' : '',
                  row.type === 'delete' ? 'border-r border-l border-red-500' : '',
                ]"
              >
                <label class="select-none text-gray-600 flex">
                  <div class="line-num">{{ row.oldNumber }}</div>
                  <div class="line-num">{{ row.newNumber }}</div>
                </label>
              </td>

              <td
                :class="[
                  'code-row-wrapper relative z-10',
                  row.type === 'insert' ? 'bg-green-50' : '',
                  row.type === 'delete' ? 'bg-red-50' : '',
                ]"
              >
                <div class="d2h-code-line relative z-0 px-4">
                  <span v-if="row.type === 'context'" class="d2h-code-line-prefix">&nbsp;</span>
                  <span v-else class="d2h-code-line-prefix">{{ row.prefix }}</span>
                  <span
                    v-if="row.content"
                    class="d2h-code-line-ctn whitespace-pre"
                    v-html="row.content"
                  />
                  <span v-else class="d2h-code-line-ctn whitespace-pre">{{
                    row.originalContent
                  }}</span>
                </div>
              </td>
            </tr>
          </template>
        </template>
      </tbody>
    </table>
  </div>
</template>

<script lang="ts">
import { defineComponent, PropType } from 'vue'
import * as Diff2Html from 'diff2html'
import highlight from '../../highlight/highlight'
import { Block, HighlightedBlock } from './event'
import { gql } from '@urql/vue'
import { DiffTable_HunkFragment } from './__generated__/DiffTable'

// This component shares a lot of code with DifferFile, can they be combined and/or split in some nicer way?

export const HUNK_FRAGMENT = gql`
  fragment DiffTable_Hunk on Hunk {
    id
    patch
  }
`

export default defineComponent({
  props: {
    unparsedDiff: {
      type: Object as PropType<DiffTable_HunkFragment>,
      required: true,
    },
    grayedOut: {
      type: Boolean,
      required: true,
    },
  },
  computed: {
    parsedDiff() {
      return Diff2Html.parse(this.unparsedDiff.patch)
    },
  },
  methods: {
    highlightedBlocks(input: Array<Block>, lang: string): Array<HighlightedBlock> {
      return highlight(input, lang, true)
    },
  },
})
</script>
