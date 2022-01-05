<template>
  <div class="gap-2 max-w-[42rem] sm:rounded-lg block border-2 border-gray-200 p-3">
    <article>
      <Banner
        v-if="show_fail_message"
        class="mb-4"
        status="error"
        message="Could not submit your comment right now. Please try again later."
      />

      <div class="flex gap-2 px-4 py-2 w-full">
        <Avatar :author="user" size="8" />

        <form class="flex-1" @submit.stop.prevent="submit">
          <div>
            <label for="comment" class="sr-only">New comment</label>
            <TextareaAutosize
              id="comment"
              ref="comment"
              v-model="message"
              name="comment"
              rows="3"
              class="shadow-sm block w-full focus:ring-blue-500 focus:border-blue-500 sm:text-sm border-gray-300 rounded-md"
              placeholder="Leave a comment"
              :user="user"
              :members="members"
              @keydown="onkey"
            />
          </div>
          <div class="mt-3 flex items-center justify-between">
            <div class="flex-1" />
            <Button @click.stop.prevent="$emit('cancel')"> Cancel</Button>
            <Button button-type="submit" color="blue" class="ml-2" :disabled="message === ''">
              Comment
            </Button>
          </div>
        </form>
      </div>
    </article>
  </div>
</template>

<script lang="ts">
import { defineComponent, PropType } from 'vue'
import Button from '../shared/Button.vue'
import Banner from '../shared/Banner.vue'
import TextareaAutosize from '../shared/TextareaAutosize.vue'
import { useCreateComment } from '../../mutations/useCreateComment'
import { UserFragment, MemberFragment } from '../shared/__generated__/TextareaMentions'
import Avatar from '../shared/Avatar.vue'
import { gql } from '@urql/vue'
import {
  ReviewNewCommentViewFragment,
  ReviewNewCommentWorkspaceFragment,
  ReviewNewCommentChangeFragment,
} from './__generated__/ReviewNewComment'
import { ConvertEmojiToColons } from '../emoji/emoji'
import {
  CommentState,
  SetCommentComposingReply,
  temporaryNewCommentID,
} from '../comments/CommentState'

export const WORKSPACE = gql`
  fragment ReviewNewCommentWorkspace on Workspace {
    id
  }
`

export const VIEW = gql`
  fragment ReviewNewCommentView on View {
    id
  }
`

export const CHANGE = gql`
  fragment ReviewNewCommentChange on Change {
    id
  }
`

export default defineComponent({
  components: {
    Banner,
    Button,
    TextareaAutosize,
    Avatar,
  },
  props: {
    path: {
      type: String,
      required: true,
    },
    lineStart: {
      type: Number,
      required: true,
    },
    // lineEnd: Number, // TODO: Support multiline comments
    lineIsNew: {
      type: Boolean,
      required: true,
    },
    user: {
      type: Object as PropType<UserFragment>,
    },
    members: {
      type: Array as PropType<MemberFragment[]>,
      required: true,
    },

    workspace: {
      type: Object as PropType<ReviewNewCommentWorkspaceFragment>,
    },
    view: {
      type: Object as PropType<ReviewNewCommentViewFragment>,
    },
    change: {
      type: Object as PropType<ReviewNewCommentChangeFragment>,
    },

    commentsState: {
      type: Object as PropType<Map<string, CommentState>>,
      required: true,
    },
  },
  emits: ['cancel', 'submitted', 'set-comment-composing-reply'],
  setup() {
    const createCommentResult = useCreateComment()

    return {
      async createComment(
        message: string,
        path: string | undefined,
        lineStart: number | undefined,
        lineEnd: number | undefined,
        lineIsNew: boolean,
        workspaceID: string | null,
        viewID: string | null,
        changeID: string | null
      ) {
        await createCommentResult({
          message: ConvertEmojiToColons(message),
          path,
          lineStart,
          lineEnd,
          lineIsNew,
          workspaceID,
          viewID,
          changeID,
        })
      },
    }
  },
  data() {
    return {
      message: '',
      show_fail_message: false,
    }
  },
  mounted() {
    this.$nextTick(() => {
      this.$refs.comment.$el.focus()
    })
  },
  beforeMount() {
    // Restore message if the component has been unmounted and re-mounted
    let id = temporaryNewCommentID(this.$props.path, this.$props.lineStart, this.$props.lineIsNew)
    let state = this.$props.commentsState.get(id)
    if (state && state.composingReply) {
      this.message = state.composingReply
    }
  },
  methods: {
    onkey(e) {
      // Escape cancels if there is no message
      if (e.keyCode === 27 && !this.message) {
        this.$emit('cancel')
        e.stopPropagation()
        e.preventDefault()
        return
      }

      // Cmd + Enter submits
      if ((e.metaKey || e.ctrlKey) && e.keyCode === 13) {
        this.submit()
        e.stopPropagation()
        e.preventDefault()
        return
      }

      let evt: SetCommentComposingReply = {
        commentId: temporaryNewCommentID(this.path, this.lineStart, this.lineIsNew),
        composingReply: this.message,
      }
      this.$emit('set-comment-composing-reply', evt)

      // Stop bubbling (Cmd + A) should select all text, not allow to pick diffs, etc.
      e.stopPropagation()
    },
    submit() {
      if (!this.message) {
        return
      }

      this.show_fail_message = false

      this.createComment(
        this.message,
        this.path,
        this.lineStart,
        this.lineStart, // TODO: Support multiline comments
        this.lineIsNew,
        this.workspace?.id,
        this.view?.id,
        this.change?.id
      )
        .then(() => {
          this.$emit('submitted')
          this.emitter.emit('local-new-comment')
        })
        .catch((err) => {
          console.error(err)
          this.show_fail_message = true
        })
    },
  },
})
</script>
