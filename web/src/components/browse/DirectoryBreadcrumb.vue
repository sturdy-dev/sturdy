<template v-if="path">
  <h3 v-if="pathBreadcrumbs" class="text-left text-md font-medium leading-6 flex items-center">
    <router-link class="text-gray-600 hover:text-gray-500" :to="{ name: 'codebaseHome' }">
      {{ codebase.name }}
    </router-link>
    <Slash />

    <template v-for="(pathPart, idx) in pathBreadcrumbs" :key="idx">
      <router-link
        v-if="!pathPart.isCurrent"
        :to="{ name: 'browseFile', params: { path: pathPart.fullPath.split('/') } }"
        class="text-gray-600 hover:text-gray-500"
      >
        {{ pathPart.name }}
      </router-link>
      <span v-else class="text-gray-400">
        {{ pathPart.name }}
      </span>

      <Slash v-if="!pathPart.isCurrent" />
    </template>
  </h3>
  <h3 v-else class="text-left text-md font-medium leading-6">Code</h3>
</template>

<script lang="ts">
import { defineComponent, PropType } from 'vue'
import { Breadcrumbs } from './Breadcrumbs'
import Slash from './Slash.vue'
import { gql } from '@urql/vue'
import { DirectoryBreadcrumbFragment } from './__generated__/DirectoryBreadcrumb'

export const DIRECTORY_BREADCRUMB = gql`
  fragment DirectoryBreadcrumb on Codebase {
    id
    shortID
    name
  }
`

export default defineComponent({
  name: 'DirectoryBreadcrumb',
  components: { Slash },
  props: {
    path: {
      required: true,
      type: String,
    },
    codebase: {
      required: true,
      type: Object as PropType<DirectoryBreadcrumbFragment>,
    },
  },
  computed: {
    pathBreadcrumbs() {
      if (this.path) {
        return Breadcrumbs.paths(this.path)
      }
      return null
    },
  },
})
</script>
