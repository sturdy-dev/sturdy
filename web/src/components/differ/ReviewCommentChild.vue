<template>
  <div
    v-if="show"
    :id="comment.id"
    ref="comment"
    :class="[highlighted ? 'ring mb-2 rounded-md' : '', ' p-4 group hover:bg-gray-50']"
  >
    <article :aria-labelledby="'comment-' + comment.id">
      <div>
        <div class="flex space-x-3 items-center">
          <div class="flex-shrink-0">
            <Avatar size="10" :author="comment.author" />
          </div>
          <div class="min-w-0 flex-1">
            <p class="text-sm font-medium text-gray-900">
              {{ comment.author.name }}
            </p>
            <p class="text-sm text-gray-500">
              <a :title="local_date(comment.created_at || comment.createdAt)">
                {{ friendly_ago(comment.created_at || comment.createdAt) }}
              </a>
            </p>
          </div>

          <ReviewCommentMenu
            class="group-hover:block hidden"
            :comment="comment"
            :can-edit="canEdit"
            @start-edit="startEdit"
            @delete="archive"
          />

          <div
            v-if="showCollapseButton"
            class="rounded-md text-blue-300 p-2 hover:bg-gray-100 cursor-pointer transition-all"
            @click="$emit('collapse')"
          >
            <ChevronDoubleUpIcon class="h-6 w-6" />
          </div>

          <Tooltip v-if="showResolveButton" x-direction="left">
            <template #default>
              <div
                class="rounded-md text-green-300 p-2 hover:bg-gray-100 cursor-pointer transition-all"
                @click="resolve"
              >
                <CheckIcon class="h-6 w-6" />
              </div>
            </template>
            <template #tooltip> Resolve comment </template>
          </Tooltip>
        </div>
      </div>
      <Banner v-if="show_delete_failed" status="error" class="my-2">
        Failed to archive the comment, try again later!
      </Banner>
      <form v-if="editing" @submit.stop.prevent="completeEdit">
        <TextareaAutosize
          ref="updatedComment"
          v-model="editingMessage"
          :value="editingMessage"
          name="comment"
          rows="3"
          class="shadow-sm block w-full focus:ring-blue-500 focus:border-blue-500 sm:text-sm border-gray-300 rounded-md mt-2"
          placeholder="Leave a comment"
          :user="user"
          :members="members"
          @click.stop
          @keydown="onkey"
        />
        <div class="mt-3 flex items-center justify-between">
          <div class="flex-1" />
          <Button @click.stop.prevent="editing = false"> Cancel</Button>
          <Button button-type="submit" color="blue" class="ml-2" :disabled="editingMessage === ''">
            Save
          </Button>
        </div>
      </form>
      <CommentMessage
        v-else
        :message="comment.message"
        :user="user"
        :members="members"
        :class="[clamped ? 'line-clamp-2 cursor-pointer' : '']"
        @click="clamped = false"
      />
    </article>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import type { PropType } from 'vue'
import type { Comment } from './event'
import time from '../../time'
import Avatar from '../../atoms/Avatar.vue'
import { Banner } from '../../atoms'
import ReviewCommentMenu from './ReviewCommentMenu.vue'
import Button from '../../atoms/Button.vue'
import TextareaAutosize from '../../atoms/TextareaAutosize.vue'
import { useUpdateComment } from '../../mutations/useUpdateComment'
import { useDeleteComment } from '../../mutations/useDeleteComment'
import type { UserFragment, MemberFragment } from '../../atoms/__generated__/TextareaMentions'
import CommentMessage from '../../atoms/CommentMessage.vue'
import mentionify from '../../atoms/mentionify'
import { ChevronDoubleUpIcon, CheckIcon } from '@heroicons/vue/solid'
import { useResolveComment } from '../../mutations/useResolveComment'
import Tooltip from '../../atoms/Tooltip.vue'

export default defineComponent({
  components: {
    Avatar,
    CommentMessage,
    Banner,
    ReviewCommentMenu,
    Button,
    TextareaAutosize,
    ChevronDoubleUpIcon,
    CheckIcon,
    Tooltip,
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
    showCollapseButton: {
      type: Boolean,
      required: false,
      default: () => {
        return false
      },
    },
    showResolveButton: {
      type: Boolean,
      required: false,
      default: () => {
        return false
      },
    },
  },
  emits: ['archived', 'prearchived', 'collapse'],
  setup() {
    const updateCommentResult = useUpdateComment()
    const deleteCommentResult = useDeleteComment()
    const resolveCommentResult = useResolveComment()

    return {
      deleteComment(id: string) {
        return deleteCommentResult(id)
      },
      resolveComment(id: string) {
        return resolveCommentResult(id)
      },
      updateComment(id: string, message: string) {
        return updateCommentResult({ id, message })
      },
    }
  },
  data() {
    return {
      now: new Date(),
      updateNowInterval: 0,
      show_delete_failed: false,
      show: true,
      editing: false,
      editingMessage: '',

      clamped: true,
    }
  },
  computed: {
    highlighted() {
      return `#${this.comment.id}` === this.$route.hash
    },
    canEdit() {
      return this.comment.author.id === this.user?.id
    },
  },
  mounted() {
    if (this.highlighted) this.scrollIntoView()
    this.updateNowInterval = window.setInterval(() => {
      this.now = new Date()
    }, 1000)
  },
  unmounted() {
    clearInterval(this.updateNowInterval)
  },
  methods: {
    scrollIntoView() {
      const comment = this.$refs.comment as HTMLElement
      comment.scrollIntoView({
        behavior: 'smooth',
        block: 'center',
        inline: 'center',
      })
    },
    friendly_ago(ts: number): string {
      return time.getRelativeTime(new Date(ts * 1000), this.now)
    },
    local_date(ts: number): string {
      return new Date(ts * 1000).toLocaleString()
    },
    archive() {
      this.show_delete_failed = false
      this.$emit('prearchived')

      this.deleteComment(this.comment.id)
        .then(() => {
          this.show = false
          this.$emit('archived')
        })
        .catch((err) => {
          console.error(err)
          this.show_delete_failed = true
        })
    },
    resolve() {
      this.resolveComment(this.comment.id)
        .then(() => {
          this.show = false
          this.$emit('archived')
        })
        .catch((err) => {
          console.error(err)
          this.show_delete_failed = true
        })
    },
    startEdit() {
      this.editing = true
      this.editingMessage = mentionify(this.comment.message, '@', this.members)

      this.$nextTick(() => {
        this.$nextTick(() => {
          this.$nextTick(() => {
            this.$refs.updatedComment.$el.focus()
          })
        })
      })
    },
    completeEdit() {
      this.updateComment(this.comment.id, this.editingMessage).then(() => {
        this.editing = false
      })
    },
    onkey(e: KeyboardEvent) {
      // Cmd + Enter submits
      if ((e.metaKey || e.ctrlKey) && e.keyCode === 13) {
        this.completeEdit()
        e.stopPropagation()
        e.preventDefault()
        return
      }

      // Stop bubbling (Cmd + A) should select all text, not allow to pick diffs, etc.
      e.stopPropagation()
    },
  },
})
</script>
