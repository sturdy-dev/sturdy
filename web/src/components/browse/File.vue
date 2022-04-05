<template>
  <div class="border sm:rounded-lg flex flex-col overflow-hidden ph-no-capture">
    <div class="px-4 py-2 bg-gray-50 space-y-1 border-b border-gray-200">
      <span v-if="hideBreadcrumb">{{ file.path }}</span>
      <DirectoryBreadcrumb v-else :path="file.path" :codebase="codebase" />
      <div v-if="!isMarkdown" class="text-sm text-gray-500">{{ lines.length }} lines</div>
    </div>

    <div v-if="isImage && rawURL" class="m-auto p-4">
      <img :src="rawURL" class="" />
    </div>

    <div
      v-else-if="isMarkdown"
      class="p-4 prose prose-yellow break-word max-w-[60rem] readme-prose"
      v-html="render"
    />

    <div v-else-if="!showFile">
      <p class="p-4">
        Large files are hidden by default.
        <a href="#" class="text-blue-600" @click.stop.prevent="setForceShow">Show now.</a>
      </p>
    </div>

    <table v-else class="p-4 leading-4 text-sm px-4 font-mono">
      <tbody v-if="hl">
        <tr v-for="(line, idx) in hl" :key="idx">
          <td
            class="bg-white sticky left-0 text-gray-600 flex justify-end pr-2 pl-12 font-mono select-none w-2"
          >
            <span>{{ idx + 1 }}</span>
          </td>
          <td
            class="px-4 font-mono break-all border-l border-blue-500 w-full whitespace-pre"
            v-html="line"
          ></td>
        </tr>
      </tbody>
      <tbody v-else>
        <tr v-for="(line, idx) in lines" :key="idx">
          <td
            class="bg-white sticky left-0 text-gray-600 flex justify-end pr-2 pl-12 font-mono select-none w-2"
          >
            <span>{{ idx + 1 }}</span>
          </td>
          <td
            class="px-4 font-mono break-all border-l border-blue-500 bg-gray-50 w-full whitespace-pre"
          >
            {{ line }}
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<script lang="ts">
import highlight from '../../highlight/highlight_file'
import '../../highlight/highlight_common_languages'
import { gql } from '@urql/vue'
import type { PropType } from 'vue'
import { defineComponent, ref } from 'vue'
import type { OpenFileFragment } from './__generated__/File'
import type { DirectoryBreadcrumbFragment } from './__generated__/DirectoryBreadcrumb'
import DirectoryBreadcrumb from './DirectoryBreadcrumb.vue'

import { Marked } from '@ts-stack/markdown'
import { SturdyMarkdownRenderer } from './SturdyMarkdownRenderer'
import { FileType } from '../../__generated__/types'
import http from '../../http'

Marked.setOptions({
  renderer: new SturdyMarkdownRenderer(),
  gfm: true,
  tables: true,
  breaks: false,
  pedantic: false,
  sanitize: false,
  smartLists: true,
  smartypants: false,
})

export const OPEN_FILE = gql`
  fragment OpenFile on File {
    id
    path
    contents
    mimeType
    info {
      id
      fileType
      rawURL
    }
  }
`

export default defineComponent({
  props: {
    file: {
      required: true,
      type: Object as PropType<OpenFileFragment>,
    },
    codebase: {
      required: true,
      type: Object as PropType<DirectoryBreadcrumbFragment>,
    },
    hideBreadcrumb: {
      type: Boolean,
      required: false,
    },
  },
  components: { DirectoryBreadcrumb },
  setup() {
    let forceShowDiff = ref(false)

    return {
      forceShowDiff,
      setForceShow: function () {
        forceShowDiff.value = true
      },
    }
  },
  computed: {
    isMarkdown() {
      return this.file.mimeType === 'text/markdown'
    },
    render() {
      return Marked.parse(this.file.contents)
    },
    lines() {
      return this.file.contents.split('\n')
    },
    showFile() {
      return this.file.contents.length < 40_000 || this.forceShowDiff
    },
    enableHighlight() {
      return this.file.contents.length < 40_000
    },
    hl() {
      // TODO: A limitation in this highlighting is that blocks can't span more than one row. Multiline comments or string literals are not correctly highlighted.

      if (!this.enableHighlight) {
        return false
      }

      const ext = this.file.path.split('.').pop()
      if (!ext) {
        return false
      }

      const high = highlight(this.file.contents, ext)
      if (high) {
        return high.split('\n')
      }
      return false
    },
    isImage() {
      return this.file.info?.fileType === FileType.Image
    },
    rawURL() {
      const url = this.file.info?.rawURL
      if (url) {
        return this.apiURL(url)
      }
      return undefined
    },
  },
  methods: {
    apiURL(path: string): string {
      const base = http.url(path)
      // using the current browser location as the base, used if url() returns a relative url
      return new URL(base, new URL(window.location.href)).href
    },
  },
})
</script>

<style>
.readme-prose img {
  margin-top: 0 !important;
  margin-bottom: 0 !important;
  display: initial !important;
}
</style>
