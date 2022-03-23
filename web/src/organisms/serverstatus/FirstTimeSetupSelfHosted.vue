<template>
  <div class="max-w-3xl mx-auto mt-16 flex flex-col gap-2">
    <OrganizationCreate
      :user="user"
      :with-tier-picker="false"
      :is-first="true"
      :is-multi-tennant="false"
    />
  </div>
</template>

<script lang="ts">
import type { PropType } from 'vue'
import { defineComponent } from 'vue'
import OrganizationCreate from '../organization/OrganizationCreate.vue'
import { gql } from '@urql/vue'
import type { FirstTimeSetupSelfHostedUserFragment } from './__generated__/FirstTimeSetupSelfHosted'

export const FIRST_TIME_USER = gql`
  fragment FirstTimeSetupSelfHostedUser on User {
    id
    name
  }
`

export default defineComponent({
  components: { OrganizationCreate },
  props: {
    user: {
      type: Object as PropType<FirstTimeSetupSelfHostedUserFragment>,
      required: true,
    },
  },
})
</script>
