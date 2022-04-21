<template>
  <div class="mb-6">
    <div class="flex space-x-3">
      <div class="flex-shrink-0">
        <div class="relative">
          <Avatar :author="user" size="10" />
          <span class="absolute -bottom-0.5 -right-1 bg-white rounded-tl px-0.5 py-px">
            <ChatAltIcon class="h-5 w-5 text-gray-400" aria-hidden="true" />
          </span>
        </div>
      </div>
      <div class="min-w-0 flex-1">
        <Banner
          v-if="failing"
          class="mb-4"
          status="error"
          message="Could not submit your comment right now. Please try again later."
        />
        <form @submit.stop.prevent="submit">
          <div>
            <label for="comment" class="sr-only">Comment</label>
            <TextareaAutosize
              ref="comment"
              :key="counter"
              v-model="message"
              name="comment"
              :user="user"
              :members="members"
              rows="3"
              class="shadow-sm block w-full focus:ring-blue-500 focus:border-blue-500 sm:text-sm border-gray-300 rounded-md"
              placeholder="Leave a comment"
              @keydown="onkey"
            />
          </div>
          <div class="mt-4 flex justify-end">
            <Button color="black" size="taller" @click="submit">Comment</Button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import type { PropType } from 'vue'
import { gql } from '@urql/vue'
import { ChatAltIcon } from '@heroicons/vue/solid'
import Avatar from '../atoms/Avatar.vue'
import TextareaAutosize, { MEMBER_FRAGMENT } from '../components/../atoms/TextareaAutosize.vue'
import { Banner } from '../atoms'
import { ConvertEmojiToColons } from '../components/emoji/emoji'
import type { AuthorFragment } from '../atoms/__generated__/AvatarHelper'
import { useCreateComment } from '../mutations/useCreateComment'
import { defineComponent } from 'vue'
import Button from '../atoms/Button.vue'

export const CODEBASE_FRAGMENT = gql`
  fragment NewComment on Codebase {
    members {
      ...Member
    }
  }
  ${MEMBER_FRAGMENT}
`

export default defineComponent({
  components: {
    ChatAltIcon,
    Avatar,
    TextareaAutosize,
    Banner,
    Button,
  },
  props: {
    user: {
      type: Object as PropType<AuthorFragment>,
      required: false,
      default: null,
    },
    members: {
      type: Array as PropType<AuthorFragment[]>,
      required: true,
    },
    workspaceId: {
      type: String,
      required: false,
      default: null,
    },
    changeId: {
      type: String,
      required: false,
      default: null,
    },
  },
  setup() {
    const createCommentResult = useCreateComment()
    return {
      async createComment(
        message: string,
        workspaceID: string | undefined,
        changeID: string | undefined
      ) {
        return createCommentResult({
          message: ConvertEmojiToColons(message),
          workspaceID,
          changeID,
        })
      },
    }
  },
  data() {
    return {
      message: '',
      failing: false,

      // <TextareaAutosize> doesn't respond well to message getting reset from outside of the component.
      // Bump counter to re-create the component from scratch when message is reset.
      counter: 0,
    }
  },
  methods: {
    onkey(e: KeyboardEvent) {
      // Escape cancels if there is no message
      if (e.keyCode === 27 && !this.message) {
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

      // Stop bubbling (Cmd + A) should select all text, not allow to pick diffs, etc.
      e.stopPropagation()
    },
    async submit() {
      if (!this.message) {
        return
      }

      this.failing = false
      await this.createComment(this.message, this.workspaceId, this.changeId)
        .then(() => {
          this.emitter.emit('local-new-comment')
          this.message = ''
          this.counter++
        })
        .catch((err) => {
          console.error(err)
          this.failing = true
        })
    },
  },
})
</script>
