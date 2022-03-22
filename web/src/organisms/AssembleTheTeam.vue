<template>
  <div>
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
  </div>
</template>

<script lang="ts">
import { defineComponent, PropType } from 'vue'

import OnboardingStep from '../components/onboarding/OnboardingStep.vue'
import CodebaseMembers, {
  CODEBASE_MEMBER_FRAGMENT,
} from '../components/codebase/CodebaseMembers.vue'

import { AuthorFragment } from '../components/shared/__generated__/AvatarHelper'

export { CODEBASE_MEMBER_FRAGMENT }

export default defineComponent({
  components: {
    CodebaseMembers,
    OnboardingStep,
  },
  props: {
    user: {
      type: Object as PropType<AuthorFragment>,
      required: false,
      default: null,
    },
    codebaseId: {
      type: String,
      required: true,
    },
    members: {
      type: Array as PropType<AuthorFragment[]>,
      required: true,
    },
    changesCount: {
      type: Number,
      required: true,
    },
  },
})
</script>
