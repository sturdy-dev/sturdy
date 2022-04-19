<template>
  <div class="relative">
    <UserAddIcon class="rounded-full h-10 w-10 text text-gray-400 bg-white" />
  </div>
  <div class="min-w-0 flex-1 break-words">
    <div>
      <p class="mt-0.5 text-sm text-gray-500">
        You've been invited to
        <router-link
          class="underline"
          :to="{
            name: 'organizationListCodebases',
            params: { organizationSlug: organizationSlug },
          }"
          @click="$emit('close')"
        >
          <strong>{{ data.organization.name }}</strong>
        </router-link>
        <RelativeTime :date="createdAt" />
      </p>
    </div>
  </div>
</template>

<script lang="ts">
import { Slug } from '../../../slug'
import { UserAddIcon } from '@heroicons/vue/solid'
import RelativeTime from '../../../atoms/RelativeTime.vue'
import { gql } from '@urql/vue'
import { defineComponent, type PropType } from 'vue'
import type { InvitedToOrganizationNotificationFragment } from './__generated__/InvitedToOrganization'

export const INVITED_TO_ORGANIZATION_NOTIFICATION_FRAGMENT = gql`
  fragment InvitedToOrganizationNotification on InvitedToOrganizationNotification {
    id
    createdAt
    archivedAt
    type
    organization {
      id
      name
      shortID
      members {
        id
        name
      }
    }
  }
`

export default defineComponent({
  components: {
    UserAddIcon,
    RelativeTime,
  },
  props: {
    data: {
      type: Object as PropType<InvitedToOrganizationNotificationFragment>,
      required: true,
    },
  },
  emits: ['close'],
  computed: {
    createdAt() {
      return new Date(this.data.createdAt * 1000)
    },
    organizationSlug() {
      return Slug(this.data.organization.name, this.data.organization.shortID)
    },
  },
})
</script>
