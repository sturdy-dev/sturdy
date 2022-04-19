<template>
  <div class="relative">
    <UserAddIcon class="rounded-full h-10 w-10 text text-gray-400 bg-white" />
  </div>
  <div class="min-w-0 flex-1 break-words">
    <div>
      <span class="mt-0.5 text-sm text-gray-500">
        You've been invited to
        <router-link
          class="underline"
          :to="{
            name: 'codebaseHome',
            params: { codebaseSlug: codebaseSlug },
          }"
          @click="$emit('close')"
        >
          <strong>{{ data.codebase.name }}</strong>
        </router-link>
        <RelativeTime :date="createdAt" />
      </span>
    </div>
  </div>
</template>

<script lang="ts">
import { Slug } from '../../../slug'
import { UserAddIcon } from '@heroicons/vue/solid'
import RelativeTime from '../../../atoms/RelativeTime.vue'
import { gql } from '@urql/vue'
import { defineComponent, type PropType } from 'vue'
import type { InvitedToCodebaseNotificationFragment } from './__generated__/InvitedToCodebase'

export const INVITED_TO_CODEBASE_NOTIFICATION_FRAGMENT = gql`
  fragment InvitedToCodebaseNotification on InvitedToCodebaseNotification {
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
  }
`

export default defineComponent({
  components: {
    UserAddIcon,
    RelativeTime,
  },
  props: {
    data: {
      type: Object as PropType<InvitedToCodebaseNotificationFragment>,
      required: true,
    },
  },
  emits: ['close'],
  computed: {
    createdAt() {
      return new Date(this.data.createdAt * 1000)
    },
    codebaseSlug() {
      return Slug(this.data.codebase.name, this.data.codebase.shortID)
    },
  },
})
</script>
