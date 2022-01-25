<template>
  <div class="space-y-2">
    <div class="space-y-1">
      <label for="add-team-members" class="block text-sm font-medium text-gray-700">
        Team Members
      </label>

      <Banner v-if="showInvitedBanner" status="success">Invited!</Banner>
      <Banner v-if="showFailedBanner" status="error">
        User not found or could not be invited.
      </Banner>

      <div class="flex">
        <div class="flex-grow">
          <input
            id="add-team-members"
            v-model="addUserEmail"
            type="text"
            name="add-team-members"
            class="block w-full shadow-sm focus:ring-sky-500 focus:border-sky-500 sm:text-sm border-gray-300 rounded-md"
            placeholder="Email address"
            aria-describedby="add-team-members-helper"
            @keydown.enter="invite"
          />
        </div>
        <span class="ml-3">
          <Button :disabled="!addUserEmail" @click="invite">
            <PlusIcon class="-ml-1 mr-1 h-5 w-5 text-gray-400" aria-hidden="true" />
            <span>Invite</span>
          </Button>
        </span>
      </div>
    </div>

    <div class="border-b border-gray-200">
      <ul role="list" class="divide-y divide-gray-200">
        <li v-for="member in members" :key="member.id" class="py-4 flex">
          <Avatar :author="member" size="10" />
          <div class="ml-3 flex flex-col">
            <span class="text-sm font-medium text-gray-900">{{ member.name }}</span>
            <span class="text-sm text-gray-500">{{ member.email }}</span>
          </div>
        </li>
      </ul>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent, PropType } from 'vue'
import { gql, useMutation } from '@urql/vue'
import {
  InviteUserToOrganizationMutation,
  InviteUserToOrganizationMutationVariables,
  OrganizationMembersFragment,
} from './__generated__/OrganizationMembers'
import { PlusIcon } from '@heroicons/vue/solid'
import Avatar from '../../components/shared/Avatar.vue'
import Button from '../../components/shared/Button.vue'
import Banner from '../../components/shared/Banner.vue'

export const ORGANIZATION_MEMBERS_FRAGMENT = gql`
  fragment OrganizationMembers on Author {
    id
    name
    email
    avatarUrl
  }
`

export default defineComponent({
  components: { PlusIcon, Avatar, Button, Banner },
  props: {
    organizationId: {
      type: String,
      required: true,
    },
    members: {
      type: Array as PropType<Array<OrganizationMembersFragment>>,
      required: true,
    },
  },
  setup() {
    let { executeMutation: execInviteUserToOrganization } = useMutation<
      InviteUserToOrganizationMutation,
      InviteUserToOrganizationMutationVariables
    >(gql`
      mutation InviteUserToOrganization($email: String!, $organizationID: ID!) {
        addUserToOrganization(input: { email: $email, organizationID: $organizationID }) {
          id
          members {
            id
          }
        }
      }
    `)

    return {
      async inviteUserToOrganization(email: string, organizationID: string) {
        const variables = { email, organizationID }
        return execInviteUserToOrganization(variables).then((result) => {
          if (result.error) {
            throw new Error(result.error)
          }
        })
      },
    }
  },
  data() {
    return {
      addUserEmail: '',
      showInvitedBanner: false,
      showFailedBanner: false,
    }
  },
  methods: {
    async invite() {
      this.showInvitedBanner = false
      this.showFailedBanner = false

      this.inviteUserToOrganization(this.addUserEmail, this.organizationId)
        .then(() => {
          this.showInvitedBanner = true
          this.showFailedBanner = false
          this.addUserEmail = ''
        })
        .catch(() => {
          this.showInvitedBanner = false
          this.showFailedBanner = true
        })
    },
  },
})
</script>
