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

<script>
import Avatar from '../shared/Avatar.vue'
import CodebaseInviteCode from './CodebaseInviteCode.vue'

export default {
  components: { Avatar, CodebaseInviteCode },
  props: {
    codebaseId: {
      type: String,
      required: true,
    },
    members: {
      type: Array,
      required: true,
    },
    user: {
      type: Object,
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
