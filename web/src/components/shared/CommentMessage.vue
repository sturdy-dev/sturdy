<template>
  <p
    class="mt-2 text-gray-700 space-y-4 break-words whitespace-pre-wrap"
    :class="[isOnlyEmoji ? 'text-2xl' : 'text-sm']"
    v-html="ify(message)"
  ></p>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import type { PropType } from 'vue'
import linkify from 'linkify-string'
import mentionify from './mentionify'
import { EmojiConvertor } from 'emoji-js'

export type User = {
  id: string
  name: string
}

export type Author = User

export default defineComponent({
  props: {
    message: {
      type: String,
      required: true,
    },
    members: {
      type: Array as PropType<Author[]>,
      required: true,
    },
    user: {
      type: Object as PropType<User>,
    },
  },
  data() {
    let emojiConvertor = new EmojiConvertor()
    emojiConvertor.replace_mode = 'unified'
    emojiConvertor.allow_native = true
    return { emojiConvertor }
  },
  computed: {
    isOnlyEmoji(): boolean {
      let msg = this.emojiConvertor.replace_colons(this.message)

      // TODO: This regex does not match all emojis (such as flags), look into extending it
      return /^[\p{Extended_Pictographic}\u{1F3FB}-\u{1F3FF}\u{1F9B0}-\u{1F9B3}\s]+$/u.test(msg)
    },
  },
  methods: {
    ify(txt: string) {
      // linkify the content
      txt = linkify(txt, { className: 'text-gray-600 underline' })
      if (this.user) {
        // mentionify current user
        txt = mentionify(txt, '@', [this.user], 'font-semibold bg-yellow-100')
      }
      // mentionify other users
      txt = mentionify(txt, '@', this.members, 'font-semibold')

      txt = this.emojiConvertor.replace_colons(txt)

      return txt
    },
  },
})
</script>
