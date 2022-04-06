<template>
  <OnboardingStep id="LandingAChange" :dependencies="['MakingAChange', 'WorkspaceChanges']">
    <template #title>Publishing a Change</template>
    <template #description>
      When you're ready, use this button to save the changes you've made so far.
    </template>
    <Button
      color="blue"
      :disabled="merging || disabled"
      :class="[merging || disabled ? 'cursor-default' : '']"
      :show-tooltip="disabled"
      :tooltip-right="true"
      :spinner="merging"
      @click="shareChange"
    >
      <template #default>
        <template v-if="merging">Merging</template>
        <template v-else>Merge</template>
      </template>
      <template v-if="disabled" #tooltip>{{ disabledTooltipMessage }}</template>
    </Button>
  </OnboardingStep>
</template>

<script lang="ts">
import { defineComponent, inject, type PropType } from 'vue'

import OnboardingStep from '../components/onboarding/OnboardingStep.vue'
import Button from '../atoms/Button.vue'

import { useLandWorkspaceChange } from '../mutations/useLandWorkspaceChange'

export default defineComponent({
  components: {
    OnboardingStep,
    Button,
  },
  props: {
    workspaceId: {
      type: String,
      required: true,
    },
    disabled: {
      type: Boolean,
      required: false,
    },
    disabledTooltipMessage: {
      type: String,
      required: true,
    },
  },
  setup() {
    const { mutating: merging, landWorkspaceChange } = useLandWorkspaceChange()

    return {
      merging,
      landWorkspaceChange,
    }
  },
  methods: {
    shareChange() {
      return this.landWorkspaceChange({
        workspaceID: this.workspaceId,
      })
    },
  },
})
</script>
