<template>
  <p v-for="change in changes" :key="change.id">
    {{ change.title }}
  </p>
</template>

<script lang="ts">
import { gql } from '@urql/vue'
import { PropType } from 'vue'

import { STATUS_FRAGMENT } from '../../components/statuses/StatusBadge.vue'
import { AUTHOR } from '../../components/shared/AvatarHelper'

import { Changelog_ChangeFragment } from './__generated__/ChangeList'

export const CHANGELOG_CHANGE_FRAGMENT = gql`
  fragment Changelog_Change on Change {
    id
    title
    author {
      ...Author
    }
    description
    createdAt
    statuses {
      ...Status
    }
  }
  ${STATUS_FRAGMENT}
  ${AUTHOR}
`

export default {
  props: {
    changes: {
      type: Array as PropType<Changelog_ChangeFragment[]>,
      required: true,
    },
  },
}
</script>
