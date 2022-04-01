<template>
  <div class="flex -space-x-1">
    <Avatar
      v-for="author in inRangeAuthors"
      :key="author.id || author.user_id"
      :author="author"
      size="8"
      class="ring-2 ring-white block"
    />
  </div>
</template>

<script lang="ts">
import type { PropType } from 'vue'
import { defineComponent } from 'vue'
import Avatar from './Avatar.vue'
import type { AuthorFragment } from './__generated__/AvatarHelper'

export default defineComponent({
  components: { Avatar },
  props: {
    authors: {
      type: Object as PropType<Array<AuthorFragment>>,
      required: true,
    },
    max: {
      type: Number,
      default: 5,
    },
  },
  data() {
    return {}
  },
  computed: {
    inRangeAuthors() {
      if (this.authors.length <= this.max) {
        return this.authors
      }
      return this.authors.slice(0, this.max)
    },
  },
})
</script>
