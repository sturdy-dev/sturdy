<template>
  <div class="relative">
    <Avatar
      class="rounded-full bg-gray-400 flex items-center justify-center ring-8 ring-white"
      size="10"
      :author="data.suggestion.author"
    />

    <span class="absolute -bottom-0.5 -right-1 bg-white rounded-tl px-0.5 py-px">
      <CodeIcon class="h-5 w-5 text-gray-400" aria-hidden="true" />
    </span>
  </div>
  <div class="min-w-0 flex-1 break-words">
    <div>
      <div v-if="data.suggestion.author" class="text-sm">
        <a href="#" class="font-medium text-gray-900">{{ data.suggestion.author.name }}</a>
      </div>

      <p class="mt-0.5 text-sm text-gray-500">
        Suggested a new change on
        <router-link
          class="underline"
          :to="{
            name: 'workspaceHome',
            params: { codebaseSlug: codebase_slug, id: data.suggestion.for.id },
          }"
          @click="$emit('close')"
        >
          <strong>{{ data.suggestion.for.name }}</strong>
        </router-link>
        {{ friendly_ago }}
      </p>
    </div>
  </div>
</template>

<script lang="ts">
import time from '../../../time'
import { Slug } from '../../../slug'
import { CodeIcon } from '@heroicons/vue/solid'
import { gql } from '@urql/vue'
import Avatar from '../../../atoms/Avatar.vue'
import type { NewSuggestionNotificationFragment } from './__generated__/NewSuggestion'
import type { PropType } from 'vue'

export const NEW_SUGGESTION_NOTIFICATION_FRAGMENT = gql`
  fragment NewSuggestionNotification on NewSuggestionNotification {
    id
    type
    createdAt
    codebase {
      id
      shortID
      name
    }
    suggestion {
      id
      author {
        id
        name
        avatarUrl
      }
      for {
        id
        name
      }
    }
  }
`

export default {
  components: {
    Avatar,
    CodeIcon,
  },
  props: {
    data: {
      type: Object as PropType<NewSuggestionNotificationFragment>,
      required: true,
    },
    now: {
      type: Object as PropType<Date>,
      required: true,
    },
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
