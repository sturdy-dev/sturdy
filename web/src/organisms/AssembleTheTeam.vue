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

import OnboardingStep from '../components/onboarding/OnboardingStep.vue'
import CodebaseMembers, {
  CODEBASE_MEMBER_FRAGMENT,
} from '../components/codebase/CodebaseMembers.vue'

import { AuthorFragment } from '../components/shared/__generated__/Avatar'
import { CodebaseMemberFragment } from '../components/codebase/__generated__/CodebaseMembers'

export { CODEBASE_MEMBER_FRAGMENT }

export default {
  components: {
    CodebaseMembers,
    OnboardingStep,
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
