<template>
  <File v-if="isFile" :file="fileOrDirectory" :codebase="codebase" />
  <Directory v-else-if="isDirectory" :directory="fileOrDirectory" :codebase="codebase" />
</template>

<script lang="ts">
import type { OpenFileOrDirectoryFragment } from './__generated__/FileOrDirectory'
import type { PropType } from 'vue'
import { gql } from '@urql/vue'
import Directory, { OPEN_DIRECTORY } from './Directory.vue'
import File, { OPEN_FILE } from './File.vue'
import type { DirectoryBreadcrumbFragment } from './__generated__/DirectoryBreadcrumb'

export const OPEN_FILE_OR_DIRECTORY = gql`
  fragment OpenFileOrDirectory on FileOrDirectory {
    __typename
    ... on File {
      ...OpenFile
    }
    ... on Directory {
      ...OpenDirectory
    }
  }
  ${OPEN_FILE}
  ${OPEN_DIRECTORY}
`

export default {
  components: {
    Directory,
    File,
  },
  props: {
    fileOrDirectory: {
      type: Object as PropType<OpenFileOrDirectoryFragment>,
      required: true,
    },
    codebase: {
      required: true,
      type: Object as PropType<DirectoryBreadcrumbFragment>,
    },
  },
  computed: {
    isFile() {
      return this.fileOrDirectory.__typename === 'File'
    },
    isDirectory() {
      return this.fileOrDirectory.__typename === 'Directory'
    },
  },
}
</script>
