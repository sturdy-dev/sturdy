<template>
  <div class="block relative rounded-full select-none leading-none" :class="[sizeClasses]">
    <img v-if="calcUrl" class="rounded-full" :src="calcUrl" alt="" />
    <span
      v-else
      class="inline-flex items-center justify-center rounded-full w-full h-full"
      :class="[initialsBgColor]"
    >
      <span
        class="font-medium text-gray-600"
        :class="[initialsFontSize]"
        :style="[size <= 5 ? 'font-size: 0.5rem' : '']"
      >
        {{ initials }}
      </span>
    </span>

    <span
      v-if="showOnline"
      class="absolute bottom-0 right-0 block rounded-full ring-2 ring-white"
      data-testid="avatar-is-online"
      :class="[online ? 'bg-green-400' : 'bg-red-400', onlineSize]"
    />
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import type { PropType } from 'vue'
import { initials, initialsBgColor } from './AvatarHelper'
import type { AuthorFragment } from './__generated__/AvatarHelper'

export default defineComponent({
  props: {
    url: {
      type: String,
      required: false,
      default: null,
    },
    author: {
      type: Object as PropType<AuthorFragment>,
      required: true,
    },
    size: {
      type: String,
      required: false,
      default: null,
    },
    showOnline: {
      type: Boolean,
      default: false,
    },
    online: Boolean,
    onlineSize: {
      type: String,
      default: 'h-1.5 w-1.5',
    },
  },
  data() {
    return {}
  },
  computed: {
    sizeClasses() {
      if (this.size) {
        return 'w-' + this.size + ' h-' + this.size
      }
      return ''
    },
    calcUrl() {
      // REST API
      if (this.author && this.author.avatar_url) {
        return this.author.avatar_url
      }
      // GraphQL
      if (this.author && this.author.avatarUrl) {
        return this.author.avatarUrl
      }
      return this.url
    },
    initials() {
      return initials(this.author)
    },
    initialsBgColor() {
      return initialsBgColor(this.author)
    },
    initialsFontSize() {
      if (this.size < 8) {
        return 'text-xs'
      }
      return 'text-base'
    },
  },
})
</script>
