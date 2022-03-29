<template>
  <textarea ref="textarea" v-model="val" class="sans-serif"></textarea>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import type { PropType } from 'vue'
// only import type definitions here to support ssr
import type Tribute from 'tributejs/tributejs'
import type { TributeCollection, TributeItem } from 'tributejs/tributejs'
import { gql } from '@urql/vue'
import type { MemberFragment, UserFragment } from './__generated__/TextareaMentions'
import { initials, initialsBgColor } from './AvatarHelper'
import { emojis } from '../components/emoji/list/emojis'
import type { Emoji } from '../components/emoji/list/emojis'
import { EmojiConvertor } from 'emoji-js'

export const USER_FRAGMENT = gql`
  fragment User on User {
    id
  }
`

export const MEMBER_FRAGMENT = gql`
  fragment Member on Author {
    id
    name
    avatarUrl
  }
`

const withoutUser = function (user: UserFragment): (_: MemberFragment) => boolean {
  return function (m: MemberFragment): boolean {
    return m.id !== user.id
  }
}

const defaultEmojis = [
  {
    name: 'thumbsup',
  },
  {
    name: 'smiley',
  },
  {
    name: 'clap',
  },
  {
    name: 'raised_hands',
  },
  {
    name: 'star',
  },
]

export default defineComponent({
  props: {
    members: {
      type: Array as PropType<MemberFragment[]>,
      required: true,
    },
    user: {
      type: Object as PropType<UserFragment>,
    },
    modelValue: {
      type: [String, Number],
      default: '',
    },
  },
  emits: ['update:modelValue'],
  data() {
    let toNativeEmojiConverter = new EmojiConvertor()
    toNativeEmojiConverter.replace_mode = 'unified'
    toNativeEmojiConverter.allow_native = true

    return {
      tribute: {} as Tribute<Record<string | number | symbol, unknown>>,
      // data property for v-model binding with real text area tag
      val: null,
      toNativeEmojiConverter,
    }
  },
  computed: {
    people() {
      return this.user ? this.members.filter(withoutUser(this.user)) : this.members
    },

    tributeMentions(): TributeCollection<MemberFragment> {
      return {
        trigger: '@',
        values: this.people,
        selectTemplate: (item: TributeItem<MemberFragment>) => `@${item.original.name}`,
        containerClass:
          'mt-1 bg-white shadow-lg max-h-60 rounded-md py-1 text-base ring-1 ring-black ring-opacity-5 overflow-auto focus:outline-none sm:text-sm sturdy-no-click-outside',
        itemClass: 'text-gray-900 cursor-pointer select-none relative py-2 pl-3 pr-9',
        selectClass: 'item-selected',
        lookup: (i: MemberFragment) => i.name,
        noMatchTemplate: () => '',
        // template for displaying item in menu
        menuItemTemplate: function (i: TributeItem<MemberFragment>): string {
          if (i.original.avatarUrl) {
            return `<span class="inline-flex items-center gap-1"><img src="${i.original.avatarUrl}" class="w-4 h-4 rounded-full" /> <span>${i.original.name}</span></span>`
          }

          let bg = initialsBgColor(i.original)
          let n = initials(i.original)

          return `<span class="inline-flex items-center gap-1"><span class="w-4 h-4 rounded-full text-xs font-medium text-gray-600 inline-flex items-center justify-center ${bg}" style="font-size: 0.5rem">${n}</span> <span>${i.original.name}</span></span>`
        },
      }
    },

    tributeEmoji(): TributeCollection<Emoji> {
      return {
        trigger: ':',
        menuItemTemplate: (item: TributeItem<Emoji>): string => {
          let str = ':' + item.original.name + ':'
          let e = this.toNativeEmojiConverter.replace_colons(str)

          return `${e} ${item.original.name}`
        },
        selectTemplate: (item: TributeItem<Emoji>) => {
          let str = ':' + item.original.name + ':'
          return this.toNativeEmojiConverter.replace_colons(str)
        },
        lookup: (i: Emoji) => i.name,
        values: function (text, cb) {
          if (text) {
            cb(emojis)
          } else {
            cb(defaultEmojis)
          }
        },
        allowSpaces: false,
        autocompleteMode: true,
        menuItemLimit: 10,
        containerClass:
          'mt-1 bg-white shadow-lg max-h-60 rounded-md py-1 text-base ring-1 ring-black ring-opacity-5 overflow-auto focus:outline-none sm:text-sm sturdy-no-click-outside',
        itemClass: 'text-gray-900 cursor-pointer select-none relative py-2 pl-3 pr-9',
        selectClass: 'item-selected',
      }
    },
  },
  watch: {
    modelValue(val) {
      this.val = val
    },
    val(val) {
      this.$emit('update:modelValue', val)
    },
  },
  async mounted() {
    const textarea = this.$refs['textarea'] as Element
    // dynamically load implementation of TributeJS only if mounted
    const Tribute = await import('tributejs')
    this.tribute = new Tribute.default({ collection: [this.tributeMentions, this.tributeEmoji] })
    this.tribute.attach(textarea)
    textarea.addEventListener('tribute-replaced', (e) => {
      e.target?.dispatchEvent(new Event('input', { bubbles: true }))
    })
  },
  beforeUnmount() {
    const textarea = this.$refs['textarea'] as Element
    try {
      this.tribute?.detach(textarea)
    } catch {
      // do nothing
    }
  },
})
</script>
<style>
.item-selected {
  /* text-white */
  --tw-text-opacity: 1;
  color: rgba(255, 255, 255, var(--tw-text-opacity));
  /* bg-blue-600 */
  --tw-bg-opacity: 1;
  background-color: rgba(37, 99, 235, var(--tw-bg-opacity));
}
</style>
