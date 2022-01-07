<template>
  <div class="border sm:rounded-lg flex flex-col overflow-hidden">
    <div class="px-4 py-2 bg-gray-50 space-y-1">
      <span v-if="hideBreadcrumb">{{ file.path }}</span>
      <DirectoryBreadcrumb v-else :path="file.path" :codebase="codebase" />
      <div v-if="!isMarkdown" class="text-sm text-gray-500">{{ lines.length }} lines</div>
    </div>

    <div
      v-if="isMarkdown"
      class="border-t border-gray-200 p-4 prose prose-yellow break-all"
      v-html="render"
    />

    <div v-else-if="!showFile">
      <p class="p-4">
        Large files are hidden by default.
        <a href="#" class="text-blue-600" @click.stop.prevent="setForceShow">Show now.</a>
      </p>
    </div>

    <table v-else class="border-t border-gray-200 p-4 leading-4 text-sm px-4 font-mono">
      <tbody v-if="hl">
        <tr v-for="(line, idx) in hl" :key="idx">
          <td
            class="bg-white sticky left-0 z-20 text-gray-600 flex justify-end pr-2 pl-12 font-mono select-none w-2"
          >
            <span>{{ idx + 1 }}</span>
          </td>
          <td class="px-4 font-mono break-all border-l border-blue-500 w-full" v-html="line"></td>
        </tr>
      </tbody>
      <tbody v-else>
        <tr v-for="(line, idx) in lines" :key="idx">
          <td
            class="bg-white sticky left-0 z-20 text-gray-600 flex justify-end pr-2 pl-12 font-mono select-none w-2"
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
import showdown from 'showdown'
import highlight from '../../highlight/highlight_file'
import '../../highlight/highlight_common_languages'
import { gql } from '@urql/vue'
import { PropType, ref } from 'vue'
import { OpenFileFragment } from './__generated__/File'
import { DirectoryBreadcrumbFragment } from './__generated__/DirectoryBreadcrumb'
import DirectoryBreadcrumb from './DirectoryBreadcrumb.vue'

export const OPEN_FILE = gql`
  fragment OpenFile on File {
    id
    path
    contents
    mimeType
  }
`

export default {
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
    let conv = new showdown.Converter()
    conv.setFlavor('github')

    let forceShowDiff = ref(false)

    return {
      conv,
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
      return this.conv.makeHtml(this.file.contents)
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
  },
}
</script>

<style>
.prose img {
  margin-top: 0 !important;
  margin-bottom: 0 !important;
  display: initial !important;
}
</style>
