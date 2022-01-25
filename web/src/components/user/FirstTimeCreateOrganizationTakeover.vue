<template>
  <div class="max-w-3xl mx-auto mt-16 flex flex-col gap-2">
    <OrganizationCreate
      class="max-w-3xl"
      :with-tier-picker="false"
      :user="user"
      :is-first="true"
      :is-multi-tennant="false"
    />
  </div>
</template>

<script lang="ts">
import { defineComponent, PropType } from 'vue'
import { gql } from '@urql/vue'
import OrganizationCreate from '../../organisms/organization/OrganizationCreate.vue'
import { FirstTimeCreateOrganizationTakeoverUserFragment } from './__generated__/FirstTimeCreateOrganizationTakeover'

export const CREATE_ORGANIZATION_TAKEOVER_USER = gql`
  fragment FirstTimeCreateOrganizationTakeoverUser on User {
    id
    name
  }
`

export default defineComponent({
  components: { OrganizationCreate },
  props: {
    user: {
      type: Object as PropType<FirstTimeCreateOrganizationTakeoverUserFragment>,
      required: true,
    },
  },
})
</script>
