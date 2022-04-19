<template>
  <div class="relative">
    <Avatar
      class="rounded-full bg-gray-400 flex items-center justify-center ring-8 ring-white"
      size="10"
      :author="data.comment.author"
    />

    <span class="absolute -bottom-0.5 -right-1 bg-white rounded-tl px-0.5 py-px">
      <ChatAltIcon class="h-5 w-5 text-gray-400" aria-hidden="true" />
    </span>
  </div>
  <div class="min-w-0 flex-1 break-words">
    <div>
      <div class="text-sm">
        <a href="#" class="font-medium text-gray-900">{{ data.comment.author.name }}</a>
      </div>
      <p v-if="data.comment.parent?.workspace" class="mt-0.5 text-sm text-gray-500">
        <router-link
          class="underline"
          :to="{
            name: 'workspaceHome',
            params: {
              codebaseSlug: codebaseSlug,
              id: data.comment.parent.workspace.id,
            },
            hash: `#${data.comment.id}`,
          }"
          @click="$emit('close')"
        >
          <template v-if="data.comment.parent.author.id === user.id">
            Replied to your comment on
          </template>
          <template v-else>
            Replied to {{ data.comment.parent.author.name }}'s comment on
          </template>

          <strong>{{ data.comment.parent.workspace.name }}</strong>
        </router-link>
        <RelativeTime :date="createdAt" />
      </p>
      <p v-else-if="data.comment.workspace" class="mt-0.5 text-sm text-gray-500">
        <router-link
          class="underline"
          :to="{
            name: 'workspaceHome',
            params: {
              codebaseSlug: codebaseSlug,
              id: data.comment.workspace.id,
            },
            hash: `#${data.comment.id}`,
          }"
          @click="$emit('close')"
        >
          Commented on
          <strong>{{ data.comment.workspace.name }}</strong>
        </router-link>
        <RelativeTime :date="createdAt" />
      </p>

      <p v-else-if="data.comment.parent?.change?.id" class="mt-0.5 text-sm text-gray-500">
        <router-link
          class="underline"
          :to="{
            name: 'codebaseChange',
            params: {
              codebaseSlug: codebaseSlug,
              id: data.comment.parent.change.id,
            },
            hash: `#${data.comment.id}`,
          }"
          @click="$emit('close')"
        >
          Replied to {{ data.comment.parent.author.name }}'s comment on
          <strong>{{ data.comment.parent.change.title }}</strong>
        </router-link>
        <RelativeTime :date="createdAt" />
      </p>
      <p v-else-if="data.comment.change?.id" class="mt-0.5 text-sm text-gray-500">
        <router-link
          class="underline"
          :to="{
            name: 'codebaseChange',
            params: {
              codebaseSlug: codebaseSlug,
              id: data.comment.change.id,
            },
            hash: `#${data.comment.id}`,
          }"
          @click="$emit('close')"
        >
          Commented on
          <strong>{{ data.comment.change.title }}</strong>
        </router-link>
        <RelativeTime :date="createdAt" />
      </p>
      <p v-else class="mt-0.5 text-sm text-gray-500">
        Commented
        <RelativeTime :date="createdAt" />
      </p>
    </div>
    <CommentMessage :message="data.comment.message" :user="user" :members="members" />
  </div>
</template>

<script lang="ts">
import { ChatAltIcon } from '@heroicons/vue/solid'
import RelativeTime from '../../../atoms/RelativeTime.vue'
import Avatar from '../../../atoms/Avatar.vue'
import { Slug } from '../../../slug'
import CommentMessage from '../../../atoms/CommentMessage.vue'
import type { User } from '../../../atoms/CommentMessage.vue'
import { gql } from '@urql/vue'
import type { NotificationCommentFragment } from './__generated__/Comment'
import { type PropType, defineComponent } from 'vue'

export const NOTIFICATION_COMMENT_FRAGMENT = gql`
  fragment NotificationComment on CommentNotification {
    id
    createdAt
    archivedAt
    type
    comment {
      id
      message
      createdAt

      codebase {
        id
        name
        shortID
        members {
          id
          name
        }
      }

      author {
        id
        name
        avatarUrl
      }
      ... on TopComment {
        workspace {
          id
          name
        }
        change {
          id
          title
          trunkCommitID
        }
      }
      ... on ReplyComment {
        parent {
          id
          author {
            id
            name
          }
          workspace {
            id
            name
          }
          change {
            id
            title
            trunkCommitID
          }
        }
      }
    }
  }
`

export default defineComponent({
  components: {
    ChatAltIcon,
    Avatar,
    CommentMessage,
    RelativeTime,
  },
  props: {
    data: {
      type: Object as PropType<NotificationCommentFragment>,
      required: true,
    },
    user: {
      type: Object as PropType<User>,
      required: true,
    },
  },
  emits: ['close'],
  computed: {
    createdAt() {
      return new Date(this.data.createdAt * 1000)
    },
    codebaseSlug() {
      return Slug(this.codebase.name, this.codebase.shortID)
    },
    codebase() {
      return this.data.comment.codebase
    },
    members() {
      return this.codebase.members
    },
  },
})
</script>
