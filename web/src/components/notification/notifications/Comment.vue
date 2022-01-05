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
              codebaseSlug: codebase_slug,
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
        {{ friendly_ago }}
      </p>
      <p v-else-if="data.comment.workspace" class="mt-0.5 text-sm text-gray-500">
        <router-link
          class="underline"
          :to="{
            name: 'workspaceHome',
            params: {
              codebaseSlug: codebase_slug,
              id: data.comment.workspace.id,
            },
            hash: `#${data.comment.id}`,
          }"
          @click="$emit('close')"
        >
          Commented on
          <strong>{{ data.comment.workspace.name }}</strong>
        </router-link>
        {{ friendly_ago }}
      </p>

      <p v-else-if="data.comment.parent?.change?.id" class="mt-0.5 text-sm text-gray-500">
        <router-link
          class="underline"
          :to="{
            name: 'codebaseChangelog',
            params: {
              codebaseSlug: codebase_slug,
              selectedChangeID: data.comment.parent.change.id,
            },
            hash: `#${data.comment.id}`,
          }"
          @click="$emit('close')"
        >
          Replied to {{ data.comment.parent.author.name }}'s comment on
          <strong>{{ data.comment.parent.change.title }}</strong>
        </router-link>
        {{ friendly_ago }}
      </p>
      <p v-else-if="data.comment.change?.id" class="mt-0.5 text-sm text-gray-500">
        <router-link
          class="underline"
          :to="{
            name: 'codebaseChangelog',
            params: {
              codebaseSlug: codebase_slug,
              selectedChangeID: data.comment.change.id,
            },
            hash: `#${data.comment.id}`,
          }"
          @click="$emit('close')"
        >
          Commented on
          <strong>{{ data.comment.change.title }}</strong>
        </router-link>
        {{ friendly_ago }}
      </p>
      <p v-else class="mt-0.5 text-sm text-gray-500">Commented {{ friendly_ago }}</p>
    </div>
    <CommentMessage :message="data.comment.message" :user="user" :members="data.codebase.members" />
  </div>
</template>

<script lang="ts">
import { ChatAltIcon } from '@heroicons/vue/solid'
import Avatar from '../../shared/Avatar.vue'
import time from '../../../time'
import { Slug } from '../../../slug'
import CommentMessage, { User } from '../../shared/CommentMessage.vue'
import { gql } from '@urql/vue'
import { NotificationCommentFragment } from './__generated__/CommentNotification'
import { PropType } from 'vue'

export const NOTIFICATION_COMMENT_FRAGMENT = gql`
  fragment NotificationComment on CommentNotification {
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
    comment {
      id
      message
      createdAt
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

export default {
  components: {
    ChatAltIcon,
    Avatar,
    CommentMessage,
  },
  props: {
    data: {
      type: Object as PropType<NotificationCommentFragment>,
      required: true,
    },
    now: {
      type: Object as PropType<Date>,
      required: true,
    },
    user: {
      type: Object as PropType<User>,
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
