<template>
  <div class="space-y-8 mt-8">
    <CodebaseInviteMembers :codebase-i-d="codebaseId" :show-header="true" />
    <CodebaseInviteCode :codebase-i-d="codebaseId" :show-header="true" />
  </div>
</template>

<script>
import CodebaseInviteMembers from '../CodebaseInviteMembers.vue'
import CodebaseInviteCode from '../CodebaseInviteCode.vue'

export default {
  name: 'SettingsCollaborators',
  components: {
    CodebaseInviteMembers,
    CodebaseInviteCode,
  },
  props: {
    codebaseId: String,
  },
}
</script>
