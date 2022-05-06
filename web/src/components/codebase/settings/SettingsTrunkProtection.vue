<template>
  <HorizontalDivider class="mt-4" bg="bg-white">Trunk Protection</HorizontalDivider>
  <Checkbox
    id="trunk-protection"
    v-model="required"
    title="Require tests to pass before merge?"
    :description="
      required
        ? 'All statuses must be healthy and up-to-date before merging.'
        : 'Tests are not required to pass before merging.'
    "
  ></Checkbox>
</template>

<script lang="ts" setup>
import HorizontalDivider from '../../../atoms/HorizontalDivider.vue'
import { useMutation } from '@urql/vue'
import { defineProps, ref, watch, withDefaults } from 'vue'
import type { CodebaseSettingsTrunkProtectionFragment } from './__generated__/SettingsTrunkProtection'
import Checkbox from '../../../atoms/Checkbox.vue'

interface Props {
  codebase: CodebaseSettingsTrunkProtectionFragment
}

const props = withDefaults(defineProps<Props>(), {})

const { executeMutation: updateCodebaseResult } = useMutation(gql`
  mutation SettingsTrunkProtection($id: ID!, $requireHealthyStatus: Boolean!) {
    updateCodebase(input: { id: $id, requireHealthyStatus: $requireHealthyStatus }) {
      id
      requireHealthyStatus
    }
  }
`)

const required = ref(false)

watch(
  props.codebase,
  () => {
    required.value = props.codebase.requireHealthyStatus
  },
  { immediate: true }
)

watch(required, () => {
  updateCodebaseResult({
    id: props.codebase.id,
    requireHealthyStatus: required.value,
  })
})
</script>

<script lang="ts">
import { gql } from '@urql/vue'

export const CODEBASE_SETTINGS_TRUNK_PROTECTION = gql`
  fragment CodebaseSettingsTrunkProtection on Codebase {
    id
    name
    requireHealthyStatus
    writeable
  }
`
</script>
