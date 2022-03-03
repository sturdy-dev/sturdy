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
        <p v-if="item.comment.codeContext" class="mt-0.5 text-sm text-gray-500">
          <router-link :to="selfRoute" class="underline">
            Commented on {{ item.comment.codeContext.path }}
          </router-link>
          {{ friendly_ago(item.createdAt) }}
        </p>
        <p v-else-if="item.comment.parent" class="mt-0.5 text-sm text-gray-500">
          <router-link :to="selfRoute" class="underline">
            Replied to {{ item.comment.parent.author.name }}
          </router-link>
          {{ friendly_ago(item.createdAt) }}
        </p>
        <p v-else class="mt-0.5 text-sm text-gray-500">
          Commented {{ friendly_ago(item.createdAt) }}
        </p>
      </div>
      <div class="mt-2 text-sm text-gray-700">
        <CommentCodeContext v-if="item.comment.codeContext" :context="item.comment.codeContext" />
        <div v-if="item.comment.parent" class="border-l-4 border-gray-400 text-gray-600 px-2">
          <CommentMessage :message="item.comment.parent.message" :user="user" :members="members" />
        </div>
        <CommentMessage :message="item.comment.message" :user="user" :members="members" />
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
import Avatar from '../../components/shared/Avatar.vue'
import { ChatAltIcon, ReplyIcon } from '@heroicons/vue/solid'
import time from '../../time'
import CommentCodeContext from '../../components/workspace/CommentCodeContext.vue'
import CommentMessage, { User } from '../../components/shared/CommentMessage.vue'
import { gql } from '@urql/vue'
import { PropType, defineComponent } from 'vue'
import { WorkspaceCommentActivityFragment } from './__generated__/ActivityComment'
import Button from '../../components/shared/Button.vue'
import CommentReply from '../../components/comments/CommentReply.vue'

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
      }
      ... on ReplyComment {
        parent {
          id
          message
          author {
            id
            name
          }
        }
      }
    }
  }
`

export default defineComponent({
  components: {
    CommentCodeContext,
    CommentMessage,
    Avatar,
    ChatAltIcon,
    Button,
    ReplyIcon,
    CommentReply,
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
            params: { codebaseSlug: this.codebaseSlug, selectedChangeId: this.item.change.id },
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
