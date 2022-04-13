<template>
  <div class="gap-2 max-w-[42rem] sm:rounded-lg block border-2 border-gray-200">
    <div>
      <ReviewCommentChild
        :comment="comment"
        :members="members"
        :user="user"
        :show-collapse-button="isExpanded"
        :show-resolve-button="true"
        class="rounded-t-lg"
        @collapse="onCollapse"
      />
      <template v-if="isExpanded">
        <ReviewCommentChild
          v-for="reply in comment.replies"
          :key="reply.id"
          :comment="reply"
          :members="members"
          :user="user"
        />
      </template>
    </div>

    <div v-if="comment.replies.length > 0 && !isExpanded" class="p-2">
      <div
        class="px-4 py-2 inline-flex items-center gap-2 hover:bg-gray-100 cursor-pointer transition-all rounded-md w-full group text-sm font-medium"
        @click="onShowReplies"
      >
        <AvatarGroup :authors="replyingAuthors" />
        <p class="text-blue-300">
          <span class="group-hover:hidden">
            {{ comment.replies.length }}
            {{ comment.replies.length === 1 ? 'reply' : 'replies' }}
          </span>
          <span class="hidden group-hover:inline-block">View thread &raquo;</span>
        </p>
      </div>
    </div>

    <div v-if="isAuthenticated && (comment.replies.length === 0 || isExpanded)" class="px-4 py-2">
      <div class="flex gap-2 transition-all rounded-md w-full items-start text-sm">
        <Avatar :author="user" size="8" />
        <CommentReply
          ref="commentReply"
          class="flex-1"
          :reply-to="comment"
          :user="user"
          :members="members"
          :start-expanded="false"
          :parent-comment-state="commentState"
          @replied="onShowReplies"
          @set-comment-composing-reply="$emit('set-comment-composing-reply', $event)"
        />
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import type { PropType } from 'vue'
import type { Comment } from './event'
import ReviewCommentChild from './ReviewCommentChild.vue'
import { useUpdateComment } from '../../mutations/useUpdateComment'
import { useDeleteComment } from '../../mutations/useDeleteComment'
import type { MemberFragment, UserFragment } from '../../atoms/__generated__/TextareaMentions'
import CommentReply from '../comments/CommentReply.vue'
import Avatar from '../../atoms/Avatar.vue'
import AvatarGroup from '../../atoms/AvatarGroup.vue'
import type { CommentState, SetCommentExpandedEvent } from '../comments/CommentState'

export default defineComponent({
  components: {
    AvatarGroup,
    ReviewCommentChild,
    CommentReply,
    Avatar,
  },
  props: {
    comment: {
      type: Object as PropType<Comment>,
      required: true,
    },
    // The logged in user
    user: {
      type: Object as PropType<UserFragment>,
      required: false,
      default: null,
    },
    // members of the selected codebase
    members: {
      type: Array as PropType<MemberFragment[]>,
      required: true,
    },

    commentState: {
      type: Object as PropType<CommentState>,
      required: true,
    },
  },
  emits: ['replied', 'showReply', 'set-comment-expanded', 'set-comment-composing-reply'],
  setup() {
    const updateCommentResult = useUpdateComment()
    const deleteCommentResult = useDeleteComment()

    return {
      deleteComment(id: string) {
        return deleteCommentResult(id)
      },
      updateComment(id: string, message: string) {
        return updateCommentResult({ id, message })
      },
    }
  },
  data() {
    return {
      showReplies: false,
    }
  },
  computed: {
    isAuthenticated(): boolean {
      return !!this.user
    },
    isAuthorized(): boolean {
      const isMember = this.members.some(({ id }) => id === this.user?.id)
      return this.isAuthenticated && isMember
    },
    replyingAuthors(): Array<MemberFragment> {
      let authorByID = new Map<string, MemberFragment>()
      for (const reply of this.comment.replies) {
        authorByID.set(reply.author.id, reply.author)
      }
      return Array.from(authorByID.values())
    },
    isExpanded(): boolean {
      if (this.commentState && this.commentState.isExpanded) {
        return true
      }
      return false
    },
  },
  watch: {
    $route: function () {
      this.checkHighlighted()
    },
  },
  mounted() {
    this.checkHighlighted()
  },
  methods: {
    checkHighlighted() {
      const highlightedId = this.$route.hash.replace('#', '')
      if (this.comment.replies.some((reply) => reply.id === highlightedId)) {
        this.showReplies = true
      }
    },
    onCollapse() {
      let evt: SetCommentExpandedEvent = {
        commentId: this.comment.id,
        isExpanded: false,
      }
      this.$emit('set-comment-expanded', evt)
    },
    onShowReplies() {
      let evt: SetCommentExpandedEvent = {
        commentId: this.comment.id,
        isExpanded: true,
      }
      this.$emit('set-comment-expanded', evt)
    },
  },
  directives: {
    mousedownOutside: {
      beforeMount: (el: Element, binding: Record<string, unknown>): void => {
        el.mousedownOutsideEvent = (event) => {
          // If parent has .sturdy-no-click-outside, ignore
          let e = event.target
          while (e) {
            if (e.classList.contains('sturdy-no-click-outside')) {
              return
            }
            if (e.parentElement) {
              e = e.parentElement
            } else {
              break
            }
          }

          if (!(el === event.target || el.contains(event.target))) {
            binding.value()
          }
        }
        document.addEventListener('mousedown', el.mousedownOutsideEvent)
      },
      unmounted: (el: Record<string, unknown>): void => {
        document.removeEventListener('mousedown', el.mousedownOutsideEvent)
      },
    },
  },
})
</script>
