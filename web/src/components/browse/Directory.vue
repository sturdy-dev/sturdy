<template>
  <div class="bg-white border overflow-hidden sm:rounded-lg flex flex-col ph-no-capture">
    <div class="px-4 py-2 bg-gray-50 space-y-1">
      <DirectoryBreadcrumb :path="directory?.path" :codebase="codebase" />
      <div class="text-sm text-gray-500">{{ directory.children.length }} files</div>
    </div>

    <router-link
      v-if="haveParent"
      class="flex flex-row items-center px-4 py-1 border-t text-sm hover:bg-gray-50"
      :to="{ name: 'browseFile', params: { path: parentPath } }"
    >
      <div class="flex-none mr-2">
        <FolderIcon class="text-gray-300 h-5 w-5" />
      </div>
      <div class="flex-1">..</div>
    </router-link>

    <router-link
      v-for="child in directory.children"
      :key="child.id"
      class="flex flex-row items-center px-4 py-1 border-t text-sm hover:bg-gray-50"
      :to="{ name: 'browseFile', params: { path: child.path.split('/') } }"
    >
      <div class="flex-none mr-2">
        <DocumentTextIcon v-if="child.__typename === 'File'" class="text-gray-400 h-4 w-5" />
        <FolderIcon v-else-if="child.__typename === 'Directory'" class="text-yellow-300 h-5 w-5" />
      </div>
      <div class="flex-1">
        {{ child.path.replace(directory.path + '/', '') }}
      </div>
    </router-link>
  </div>

  <File
    v-if="directory.readme"
    :file="directory.readme"
    :codebase="codebase"
    :hide-breadcrumb="true"
  />
</template>

<script lang="ts">
import { defineComponent, PropType } from 'vue'
import { gql } from '@urql/vue'
import { OpenDirectoryFragment } from './__generated__/Directory'
import { DocumentTextIcon } from '@heroicons/vue/outline'
import { FolderIcon } from '@heroicons/vue/solid'
import File, { OPEN_FILE } from './File.vue'
import DirectoryBreadcrumb from './DirectoryBreadcrumb.vue'
import { DirectoryBreadcrumbFragment } from './__generated__/DirectoryBreadcrumb'

export const OPEN_DIRECTORY = gql`
  fragment OpenDirectory on Directory {
    id
    path

    children {
      ... on File {
        id
        path
        mimeType
      }

      ... on Directory {
        id
        path
      }
    }

    readme {
      id
      ...OpenFile
    }
  }
  ${OPEN_FILE}
`

export default defineComponent({
  name: 'Directory',
  components: {
    DirectoryBreadcrumb,
    DocumentTextIcon,
    FolderIcon,
    File,
  },
  props: {
    directory: {
      required: true,
      type: Object as PropType<OpenDirectoryFragment>,
    },
    codebase: {
      required: true,
      type: Object as PropType<DirectoryBreadcrumbFragment>,
    },
  },
  computed: {
    haveParent() {
      return this.directory?.path !== '/' && this.directory?.path !== ''
    },
    parentPath() {
      let p = this.directory.path.split('/')
      p.pop()
      return p
    },
  },
})
</script>
