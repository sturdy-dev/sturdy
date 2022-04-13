<template>
  <div :id="item.comment.id" class="relative flex items-start space-x-3">
    <div class="relative">
      <Avatar :author="item.author" size="10" />
      <span class="absolute -bottom-2 -right-2 bg-white rounded-tl px-0.5 py-px">
        <ChatAltIcon class="h-5 w-5 text-gray-400" aria-hidden="true" />
      </span>
    </div>
    <div class="min-w-0 flex-1">
      <div>
        <div class="text-sm inline-flex justify-between w-full items-start">
          <a href="#" class="font-medium text-gray-900">{{ item.author.name }}</a>
          <Button size="small" @click="newReply">
            <ReplyIcon class="h-3 w-3 text-gray-500 hover:text-gray-900" />
          </Button>
        </div>

        <div class="mt-0.5 text-sm text-gray-500 inline-flex items-center">
          <span v-if="item.comment.codeContext" class="text-ellipsis overflow-hidden">
            <router-link :to="selfRoute" class="underline">
              Commented on {{ item.comment.codeContext.path }}
            </router-link>
            {{ friendly_ago(item.createdAt) }}
          </span>
          <span v-else-if="item.comment.parent">
            <router-link :to="selfRoute" class="underline">
              Replied to {{ item.comment.parent.author.name }}
            </router-link>
            {{ friendly_ago(item.createdAt) }}
          </span>
          <span v-else> Commented {{ friendly_ago(item.createdAt) }} </span>

          <Tooltip v-if="item.comment.resolved" x-direction="left">
            <template #default>
              <CheckIcon class="w-4 h-4 text-green-500 flex-shrink-0" />
            </template>
            <template #tooltip> Resolved </template>
          </Tooltip>
        </div>
      </div>
      <div class="mt-2 text-sm text-gray-700">
        <CommentCodeContext v-if="item.comment.codeContext" :context="item.comment.codeContext" />
        <div
          v-if="item.comment.parent"
          class="border-l-4 px-2"
          :class="[item.comment.parent.resolved ? 'border-green-400' : 'border-gray-400']"
        >
          <CommentMessage
            :message="item.comment.parent.message"
            :user="user"
            :members="members"
            :resolved="item.comment.parent.resolved"
          />
        </div>

        <CommentMessage
          :message="item.comment.message"
          :user="user"
          :members="members"
          :resolved="item.comment.resolved"
        />
      </div>

      <div v-if="isReplying" class="mt-2">
        <CommentReply
          ref="commentReply"
          :reply-to="item.comment.parent ?? item.comment"
          :user="user"
          :members="members"
          :start-expanded="true"
          @replied="isReplying = false"
        />
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import Avatar from '../../atoms/Avatar.vue'
import { ChatAltIcon, CheckIcon, ReplyIcon } from '@heroicons/vue/solid'
import time from '../../time'
import CommentCodeContext from '../../components/workspace/CommentCodeContext.vue'
import type { User } from '../../atoms/CommentMessage.vue'
import CommentMessage from '../../atoms/CommentMessage.vue'
import { gql } from '@urql/vue'
import type { PropType } from 'vue'
import { defineComponent } from 'vue'
import type { WorkspaceCommentActivityFragment } from './__generated__/ActivityComment'
import Button from '../../atoms/Button.vue'
import CommentReply from '../../components/comments/CommentReply.vue'
import Tooltip from '../../atoms/Tooltip.vue'

export const WORKSPACE_ACTIVITY_COMMENT_FRAGMENT = gql`
  fragment WorkspaceCommentActivity on WorkspaceCommentActivity {
    author {
      id
      name
      avatarUrl
    }
    createdAt
    change {
      id
      codebase {
        id
        members {
          id
          name
        }
      }
    }
    workspace {
      id
      codebase {
        id
        members {
          id
          name
        }
      }
    }
    comment {
      id
      message
      ... on TopComment {
        codeContext {
          id
          lineStart
          lineEnd
          lineIsNew
          context
          contextStartsAtLine
          path
        }

        resolved
      }
      ... on ReplyComment {
        parent {
          id
          message
          author {
            id
            name
          }
          resolved
        }
      }
    }
  }
`

export default defineComponent({
  components: {
    Tooltip,
    CommentCodeContext,
    CommentMessage,
    Avatar,
    ChatAltIcon,
    Button,
    ReplyIcon,
    CommentReply,
    CheckIcon,
  },
  props: {
    item: {
      type: Object as PropType<WorkspaceCommentActivityFragment>,
      required: true,
    },
    codebaseSlug: {
      type: String,
      required: true,
    },
    user: {
      type: Object as PropType<User>,
      required: false,
      default: null,
    },
  },
  data() {
    return {
      isReplying: false,
    }
  },
  computed: {
    members() {
      return this.item.workspace
        ? this.item.workspace.codebase.members
        : this.item.change
        ? this.item.change.codebase.members
        : []
    },
    selfRoute() {
      return this.item.workspace
        ? {
            name: 'workspaceHome',
            params: { codebaseSlug: this.codebaseSlug, id: this.item.workspace.id },
            hash: `#${this.item.comment.id}`,
          }
        : this.item.change
        ? {
            name: 'codebaseChange',
            params: { codebaseSlug: this.codebaseSlug, id: this.item.change.id },
            hash: `#${this.item.comment.id}`,
          }
        : {
            name: 'codebase',
            params: { codebaseSlug: this.codebaseSlug },
          }
    },
  },
  methods: {
    friendly_ago(ts: number) {
      return time.getRelativeTime(new Date(ts * 1000))
    },
    newReply() {
      this.isReplying = true
      this.$nextTick(() => {
        this.$nextTick(() => {
          this.$refs.commentReply.$refs.replyComment.$el.focus()
        })
      })
    },
  },
})
</script>
