<template>
  <div class="relative">
    <DownloadIcon class="rounded-full h-10 w-10 text text-gray-400 bg-white" />
  </div>
  <div class="min-w-0 flex-1 break-words">
    <div>
      <p class="mt-0.5 text-sm text-gray-500">
        Codebase
        <router-link
          class="underline"
          :to="{
            name: 'codebaseHome',
            params: { codebaseSlug: codebaseSlug },
          }"
          @click="$emit('close')"
        >
          <strong>{{ data.repository.name }}</strong>
        </router-link>
        is ready
        <RelativeTime :date="createdAt" />
      </p>
    </div>
  </div>
</template>

<script lang="ts">
import { Slug } from '../../../slug'
import { DownloadIcon } from '@heroicons/vue/solid'
import RelativeTime from '../../../atoms/RelativeTime.vue'
import { gql } from '@urql/vue'
import { defineComponent, type PropType } from 'vue'
import type { GitHubRepositoryImportedFragment } from './__generated__/GitHubRepositoryImoprted'

export const GITHUB_REPOSITORY_IMPORTED_NOTIFICATION_FRAGMENT = gql`
  fragment GitHubRepositoryImported on GitHubRepositoryImported {
    id
    createdAt
    archivedAt
    type
    repository {
      id
      name
      codebase {
        id
        name
        shortID
        members {
          id
          name
        }
      }
    }
  }
`

export default defineComponent({
  components: {
    DownloadIcon,
    RelativeTime,
  },
  props: {
    data: {
      type: Object as PropType<GitHubRepositoryImportedFragment>,
      required: true,
    },
  },
  emits: ['close'],
  computed: {
    createdAt() {
      return new Date(this.data.createdAt * 1000)
    },
    codebaseSlug() {
      return Slug(this.data.repository.codebase.name, this.data.repository.codebase.shortID)
    },
  },
})
</script>
