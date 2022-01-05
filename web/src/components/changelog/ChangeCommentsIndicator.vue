<template>
  <div class="text-xs rounded-md bg-gray-100 flex flex-row items-center px-1 py-0.5 gap-0.5">
    <AnnotationIcon class="h-3 w-3" />
    <span class="font-semibold"> {{ change.comments.length }}</span>
  </div>
</template>

<script lang="ts">
import { defineComponent, PropType } from 'vue'
import gql from 'graphql-tag'
import { ChangeCommentsFragment } from './__generated__/ChangeCommentsIndicator'
import { AnnotationIcon } from '@heroicons/vue/solid'

export const CHANGE_COMMENTS = gql`
  fragment ChangeComments on Change {
    id
    comments {
      id
    }
  }
`

export default defineComponent({
  name: 'ChangeCommentsIndicator',
  components: { AnnotationIcon },
  props: {
    change: {
      required: true,
      type: Object as PropType<ChangeCommentsFragment>,
    },
  },
})
</script>
