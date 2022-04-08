<template>
  <div class="flex">
    <div
      v-if="oldFileInfo?.rawURL"
      class="bg-red-100 p-4 flex flex-col space-y-2 items-center flex-1"
    >
      <span class="text-center text-red-800">Old</span>
      <div>
        <img :src="apiURL(oldFileInfo.rawURL)" />
      </div>
    </div>
    <div
      v-if="newFileInfo?.rawURL"
      class="bg-green-100 p-4 flex flex-col space-y-2 items-center flex-1"
    >
      <span class="text-center text-green-800">New</span>
      <div>
        <img :src="apiURL(newFileInfo.rawURL)" />
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { gql } from '@urql/vue'
import { defineComponent, type PropType } from 'vue'
import type { RichDiffImage_FileInfoFragment } from './__generated__/RichDiffImage'
import http from '../../http'

export const RICH_DIFF_IMAGE_FILE_INFO = gql`
  fragment RichDiffImage_FileInfo on FileInfo {
    id
    fileType
    rawURL
  }
`

export default defineComponent({
  props: {
    oldFileInfo: {
      type: Object as PropType<RichDiffImage_FileInfoFragment>,
      required: false,
      default: null,
    },
    newFileInfo: {
      type: Object as PropType<RichDiffImage_FileInfoFragment>,
      required: false,
      default: null,
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
