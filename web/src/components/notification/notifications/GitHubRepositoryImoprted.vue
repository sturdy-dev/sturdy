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
            params: { codebaseSlug: codebase_slug },
          }"
          @click="$emit('close')"
        >
          <strong>{{ data.repository.name }}</strong>
        </router-link>
        is ready
        {{ friendly_ago }}
      </p>
    </div>
  </div>
</template>

<script>
import time from '../../../time'
import { Slug } from '../../../slug'
import { DownloadIcon } from '@heroicons/vue/solid'
import { gql } from '@urql/vue'

export const GITHUB_REPOSITORY_IMPORTED_NOTIFICATION_FRAGMENT = gql`
  fragment GitHubRepositoryImported on GitHubRepositoryImported {
    id
    createdAt
    archivedAt
    type
    codebase {
      id
      name
      shortID

      members {
        id
        name
      }
    }
    repository {
      id
      name
      codebase {
        id
        name
        shortID
      }
    }
  }
`

export default {
  components: {
    DownloadIcon,
  },
  props: {
    data: { type: Object, required: true },
    now: { type: Object, required: true },
  },
  emits: ['close'],
  computed: {
    friendly_ago() {
      return time.getRelativeTime(new Date(this.data.createdAt * 1000), this.now)
    },
    codebase_slug() {
      return Slug(this.data.codebase.name, this.data.codebase.shortID)
    },
  },
}
</script>
