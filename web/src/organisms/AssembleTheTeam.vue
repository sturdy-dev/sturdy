<template>
  <OnboardingStep
    id="InvitingTeamMembers"
    :dependencies="['FindingYourCodebase']"
    :enabled="changesCount > 0"
  >
    <template #title> Assemble the Team</template>
    <template #description>
      In this section over here you can invite other people to collaborate on this codebase with
      you.
    </template>

    <CodebaseMembers :user="user" :members="members" :codebase-id="codebaseId" />
  </OnboardingStep>
</template>

<script lang="ts">
import { PropType } from 'vue'

import { AuthorFragment } from '../components/shared/__generated__/Avatar'
import { CodebaseMemberFragment } from '../components/codebase/__generated__/CodebaseMembers'

import CodebaseMembers from '../components/codebase/CodebaseMembers.vue'

export default {
  components: {
    CodebaseMembers,
  },
  props: {
    user: {
      type: Object as PropType<AuthorFragment>,
    },
    codebaseId: {
      type: String,
      required: true,
    },
    members: {
      type: Array as PropType<CodebaseMemberFragment[]>,
      required: true,
    },
    changesCount: {
      type: Number,
      required: true,
    },
  },
}
</script>
