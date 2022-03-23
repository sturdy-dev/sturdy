<template>
  <form @submit.stop.prevent="completeReply">
    <TextareaAutosize
      ref="replyComment"
      v-model="replyMessage"
      name="comment"
      :rows="expanded ? 3 : 1"
      class="shadow-sm block w-full focus:ring-blue-500 focus:border-blue-500 sm:text-sm border-gray-300 rounded-md"
      placeholder="Reply ..."
      :user="user"
      :members="members"
      @keydown="onkey"
      @click="expanded = true"
    />
    <div v-if="expanded" class="mt-3 flex items-center justify-between">
      <div class="flex-1" />
      <Button
        button-type="submit"
        color="blue"
        class="ml-2 inline-flex gap-2"
        :disabled="replyMessage === '' || isSending"
      >
        <Spinner v-if="isSending" />
        <span v-if="isSending">Replying</span>
        <span v-else>Reply</span>
      </Button>
    </div>
  </form>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import type { PropType } from 'vue'
import Button from '../shared/Button.vue'
import TextareaAutosize from '../shared/TextareaAutosize.vue'
import type { Comment } from '../differ/event'
import type { MemberFragment, UserFragment } from '../shared/__generated__/TextareaMentions'
import { useCreateComment } from '../../mutations/useCreateComment'
import Spinner from '../shared/Spinner.vue'
import type { ConvertEmojiToColons } from '../emoji/emoji'
import type { CommentState, SetCommentComposingReply } from './CommentState'

export default defineComponent({
  components: { Button, TextareaAutosize, Spinner },
  setup() {
    const createCommentResult = useCreateComment()
    return {
      async createComment(inReplyTo: string, message: string) {
        await createCommentResult({ inReplyTo, message })
      },
    }
  },
  data() {
    return {
      isSending: false,
      expanded: this.$props.startExpanded,
      replyMessage: '',
    }
  },
  beforeMount() {
    if (this.$props.parentCommentState.composingReply) {
      this.replyMessage = this.$props.parentCommentState.composingReply
      this.expanded = true
    }
  },
  props: {
    replyTo: {
      type: Object as PropType<Comment>,
      required: true,
    },
    // The logged in user
    user: {
      type: Object as PropType<UserFragment>,
    },
    // members of the selected codebase
    members: {
      type: Array as PropType<MemberFragment[]>,
      required: true,
    },
    startExpanded: {
      type: Boolean,
      required: true,
      default: () => {
        return false
      },
    },
    parentCommentState: {
      type: Object as PropType<CommentState>,
      default: () => {
        return { isExpanded: true, composingReply: undefined }
      },
    },
  },
  emits: {
    replied(payload: any) {
      return true // validation
    },
    setCommentComposingReply(payload: SetCommentComposingReply) {
      return true // validation
    },
  },
  methods: {
    onkey(e) {
      // Cmd + Enter submits reply
      if ((e.metaKey || e.ctrlKey) && e.keyCode === 13) {
        this.completeReply()
        e.stopPropagation()
        e.preventDefault()
        return
      }

      // Save current input
      this.$emit('setCommentComposingReply', {
        commentId: this.replyTo.id,
        composingReply: this.replyMessage,
      })

      // Stop bubbling (Cmd + A) should select all text, not allow to pick diffs, etc.
      e.stopPropagation()
    },
    async completeReply() {
      this.isSending = true
      await this.createComment(this.replyTo.id, ConvertEmojiToColons(this.replyMessage))
      this.isSending = false
      this.replyMessage = ''

      // reset
      this.$emit('setCommentComposingReply', {
        commentId: this.replyTo.id,
        composingReply: undefined,
      })

      this.$emit('replied')
    },
  },
})
</script>
