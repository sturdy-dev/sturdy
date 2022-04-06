<template>
  <div class="space-y-2">
    <div class="space-y-1">
      <label for="add-team-members" class="block text-sm font-medium text-gray-700">
        Members
      </label>

      <Banner v-if="showInvitedBanner" status="success">Invited!</Banner>
      <Banner v-if="showFailedBanner" status="error">
        User not found or could not be invited.
      </Banner>

      <template v-if="isMultiTenancyEnabled">
        <p class="text-sm text-gray-600">
          Invite a collaborator to <strong>{{ organization.name }}</strong> by entering the email
          address they used to sign up for Sturdy.
        </p>

        <div v-if="organization.writeable" class="flex">
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
      </template>
      <div v-else>
        <p class="text-sm text-gray-600">
          Users that sign up to this self-hosted instance of Sturdy will automatically join this
          organization.
        </p>
      </div>
    </div>

    <div v-if="!organization.writeable">
      <p class="text-sm text-gray-500">
        You don't have permissions to invite users to this organization, ask an admin for help if
        you want to invite someone.
      </p>
    </div>

    <div>
      <ul role="list" class="divide-y divide-gray-200">
        <li v-for="member in organization.members" :key="member.id" class="py-4 flex">
          <Avatar :author="member" size="10" />
          <div class="ml-3 flex flex-col flex-1">
            <span class="text-sm font-medium text-gray-900">
              {{ member.name }}
              <Pill v-if="!isActive(member.status)" color="gray">Pending invitation</Pill>
            </span>
            <span class="text-sm text-gray-500">{{ member.email }}</span>
          </div>
          <template v-if="organization.writeable">
            <Button v-if="member.id === user.id" @click="removeUser(member)">
              <UserRemoveIcon class="-ml-1 mr-1 h-5 w-5 text-gray-400" aria-hidden="true" />
              <span>Leave</span>
            </Button>
            <Button v-else @click="removeUser(member)">
              <UserRemoveIcon class="mr-1 h-5 w-5 text-gray-400" aria-hidden="true" />
              <span>Remove</span>
            </Button>
          </template>
        </li>
      </ul>
    </div>
  </div>
</template>

<script lang="ts">
import { computed, defineComponent, inject, ref } from 'vue'
import type { Ref, PropType } from 'vue'
import { gql } from '@urql/vue'
import type {
  OrganizationMembersOrganizationFragment,
  OrganizationMembersUserFragment,
} from './__generated__/OrganizationMembers'
import { PlusIcon, UserRemoveIcon } from '@heroicons/vue/solid'
import Avatar from '../../atoms/Avatar.vue'
import Button from '../../atoms/Button.vue'
import { Banner } from '../../atoms'
import { Feature, UserStatus } from '../../__generated__/types'
import { useAddUserToOrganization } from '../../mutations/useAddUserToOrganization'
import { useRemoveUserFromOrganization } from '../../mutations/useRemoveUserFromOrganization'
import Pill from '../../atoms/Pill.vue'

export const ORGANIZATION_FRAGMENT = gql`
  fragment OrganizationMembersOrganization on Organization {
    id
    name
    members {
      id
      name
      email
      avatarUrl
      status
    }
    writeable
  }
`

export const USER_FRAGMENT = gql`
  fragment OrganizationMembersUser on User {
    id
  }
`

export default defineComponent({
  components: { PlusIcon, Avatar, Button, Banner, UserRemoveIcon, Pill },
  props: {
    organization: {
      type: Object as PropType<OrganizationMembersOrganizationFragment>,
      required: true,
    },
    user: {
      type: Object as PropType<OrganizationMembersUserFragment>,
      required: true,
    },
  },
  setup() {
    const features = inject<Ref<Array<Feature>>>('features', ref([]))
    const isMultiTenancyEnabled = computed(() => features.value.includes(Feature.MultiTenancy))

    let addUserToOrganization = useAddUserToOrganization()
    let removeUserFromOrganization = useRemoveUserFromOrganization()

    return {
      isMultiTenancyEnabled,

      addUserToOrganization,
      removeUserFromOrganization,
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

      const variables = { email: this.addUserEmail, organizationID: this.organization.id }

      this.addUserToOrganization(variables)
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
    removeUser(member) {
      const variables = { userID: member.id, organizationID: this.organization.id }
      this.removeUserFromOrganization(variables).then(() => {
        this.emitter.emit('notification', {
          title: 'Removed ' + member.name,
          message: member.name + ' was removed from ' + this.organization.name,
        })
      })
    },
    isActive(status: UserStatus) {
      return status == UserStatus.Active
    },
  },
})
</script>
