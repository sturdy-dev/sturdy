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
import { defineComponent, PropType } from 'vue'

import Avatar from '../shared/Avatar.vue'
import { AuthorFragment } from '../shared/__generated__/AvatarHelper'
import CodebaseInviteCode from './CodebaseInviteCode.vue'

import { AUTHOR } from '../shared/AvatarHelper'

export { AUTHOR as CODEBASE_MEMBER_FRAGMENT }

export default defineComponent({
  components: { Avatar, CodebaseInviteCode },
  props: {
    codebaseId: {
      type: String,
      required: true,
    },
    members: {
      type: Array as PropType<AuthorFragment[]>,
      required: true,
    },
    user: {
      type: Object as PropType<AuthorFragment>,
      required: false,
      default: null,
    },
  },
  computed: {
    authenticated(): boolean {
      if (this.user) {
        return true
      }
      return false
    },
    authorized() {
      const isMember = this.members.some(({ id }) => id === this.user?.id)
      return this.authenticated && isMember
    },
  },
})
</script>
