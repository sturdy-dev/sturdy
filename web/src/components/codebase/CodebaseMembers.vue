<template>
  <div class="space-y-5">
    <h2 class="text-sm font-medium text-gray-500">Collaborators</h2>

    <div class="space-y-2 mt-3">
      <div v-for="member in members" :key="member.id" class="flex items-center space-x-2">
        <Avatar :author="member" size="5" class="flex-shrink-0" />
        <span class="text-gray-900 text-sm font-medium">{{ member.name }}</span>
      </div>
    </div>

    <div v-if="authorized">
      <p class="text-sm text-gray-500">Invite your team to Sturdy with this link</p>
      <CodebaseInviteCode :codebase-i-d="codebaseId" :small="true" />
    </div>
  </div>
</template>

<script lang="ts">
import { gql } from '@urql/vue'
import { PropType } from 'vue'

import Avatar from '../shared/Avatar.vue'
import { AuthorFragment } from '../shared/__generated__/Avatar'
import CodebaseInviteCode from './CodebaseInviteCode.vue'

import { CodebaseMemberFragment } from './__generated__/CodebaseMembers'

export const CODEBASE_MEMBER_FRAGMENT = gql`
  fragment CodebaseMember on Author {
    id
    name
  }
`

export default {
  components: { Avatar, CodebaseInviteCode },
  props: {
    codebaseId: {
      type: String,
      required: true,
    },
    members: {
      type: Array as PropType<CodebaseMemberFragment[]>,
      required: true,
    },
    user: {
      type: Object as PropType<AuthorFragment>,
    },
  },
  computed: {
    authenticated() {
      return !!this.user
    },
    authorized() {
      const isMember = this.members.some(({ id }) => id === this.user?.id)
      return this.authenticated && isMember
    },
  },
}
</script>
