<template>
  <div v-if="data && data.codebase" class="flex flex-col">
    <div>
      <template v-if="showHeader">
        <label for="add_team_members" class="block text-sm font-medium text-gray-700">
          Add collaborator
        </label>
        <p class="mt-1 text-sm text-gray-500">
          Invite a collaborator to <strong>{{ data.codebase.name }}</strong> using the email address
          of their Sturdy-account
        </p>
      </template>
      <p id="add_team_members_helper" class="sr-only">Search by email address</p>
      <div class="flex mt-1">
        <div class="flex-grow">
          <input
            id="add_team_members"
            v-model="inviteUserEmail"
            type="text"
            name="add_team_members"
            class="block w-full shadow-sm focus:ring-light-blue-500 focus:border-light-blue-500 sm:text-sm border-gray-300 rounded-md"
            placeholder="Email address"
            aria-describedby="add_team_members_helper"
            @keydown.enter="inviteMember"
          />
        </div>
        <span class="ml-3">
          <button
            type="button"
            class="bg-white inline-flex items-center px-4 py-2 border border-gray-300 shadow-sm text-sm font-medium rounded-md text-gray-700 hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-light-blue-500"
            @click.stop.prevent="inviteMember"
          >
            <PlusIcon class="-ml-2 mr-1 h-5 w-5 text-gray-400" />
            <span>Add</span>
          </button>
        </span>
      </div>
    </div>

    <p>{{ inviteMemberStatus }}</p>

    <div class="mt-2">
      <p
        v-if="
          data.codebase.organization &&
          data.codebase.indirectMembers &&
          data.codebase.indirectMembers.length > 0
        "
        class="my-1 text-sm text-gray-500"
      >
        The following users are direct members of this codebase
      </p>

      <ul class="divide-y divide-gray-200">
        <li v-for="member in data.codebase.directMembers" :key="member.id" class="py-4 flex">
          <Avatar :author="member" size="10" />
          <div class="ml-3 flex flex-col flex-1">
            <span class="text-sm font-medium text-gray-900">{{ member.name }}</span>
            <span class="text-sm text-gray-500">{{ member.email }}</span>
          </div>

          <template v-if="data.codebase.writeable">
            <Button v-if="member.id === data.user.id" @click="removeMember(member)">
              <UserRemoveIcon class="-ml-1 mr-1 h-5 w-5 text-gray-400" aria-hidden="true" />
              <span>Leave</span>
            </Button>
            <Button v-else @click="removeMember(member)">
              <UserRemoveIcon class="mr-1 h-5 w-5 text-gray-400" aria-hidden="true" />
              <span>Remove</span>
            </Button>
          </template>
        </li>
      </ul>
    </div>

    <template
      v-if="
        data.codebase.organization &&
        data.codebase.indirectMembers &&
        data.codebase.indirectMembers.length > 0
      "
    >
      <div class="mt-2">
        <p class="my-1 text-sm text-gray-500">
          The following users are members of {{ data.codebase.organization.name }}, and also have
          access to this codebase:
        </p>

        <ul class="divide-y divide-gray-200">
          <li v-for="member in data.codebase.indirectMembers" :key="member.id" class="py-4 flex">
            <Avatar :author="member" size="10" />
            <div class="ml-3 flex flex-col">
              <span class="text-sm font-medium text-gray-900">{{ member.name }}</span>
              <span class="text-sm text-gray-500">{{ member.email }}</span>
            </div>
          </li>
        </ul>
      </div>
    </template>
  </div>
</template>

<script lang="ts">
import { PlusIcon, UserRemoveIcon } from '@heroicons/vue/solid'
import Avatar from '../shared/Avatar.vue'
import { gql, useQuery } from '@urql/vue'
import { defineComponent } from 'vue'
import type {
  CodebaseInviteMembersQuery,
  CodebaseInviteMembersQueryVariables,
} from './__generated__/CodebaseInviteMembers'
import { useAddUserToCodebase } from '../../mutations/useAddUserToCodebase'
import { useRemoveUserFromCodebase } from '../../mutations/useRemoveUserFromCodebase'
import Button from '../shared/Button.vue'

export default defineComponent({
  components: { PlusIcon, Avatar, UserRemoveIcon, Button },
  props: ['codebaseID', 'showHeader'],
  setup(props) {
    let { data, executeQuery } = useQuery<
      CodebaseInviteMembersQuery,
      CodebaseInviteMembersQueryVariables
    >({
      query: gql`
        query CodebaseInviteMembers($id: ID, $shortID: ID) {
          codebase(id: $id, shortID: $shortID) {
            id
            name

            directMembers: members(filterDirectAccess: true) {
              id
              email
              name
              avatarUrl
            }

            indirectMembers: members(filterDirectAccess: false) {
              id
              email
              name
              avatarUrl
            }

            writeable

            organization {
              id
              name
            }
          }

          user {
            id
          }
        }
      `,
      variables: {
        id: props.codebaseID,
      },
      requestPolicy: 'cache-and-network',
    })

    let addUserToCodebase = useAddUserToCodebase()
    let removeUserFromCodebase = useRemoveUserFromCodebase()

    return {
      data,

      addUserToCodebase,
      removeUserFromCodebase,

      refresh() {
        executeQuery({
          requestPolicy: 'network-only',
        })
      },
    }
  },
  data() {
    return {
      inviteUserEmail: '', // Form model
      inviteMemberStatus: '', // Form status message
    }
  },
  methods: {
    inviteMember() {
      const variables = { codebaseID: this.codebaseID, email: this.inviteUserEmail }

      this.addUserToCodebase(variables)
        .then(() => {
          this.inviteMemberStatus = 'The user was added!'
          this.inviteUserEmail = ''
        })
        .catch(() => {
          this.inviteMemberStatus =
            'The user could not be added. Check that you entered the right email, and try again.'
        })
        .finally(this.refresh)
    },

    removeMember(member: any) {
      const variables = { codebaseID: this.codebaseID, userID: member.id }
      this.removeUserFromCodebase(variables)
        .then(() => {
          this.emitter.emit('notification', {
            title: 'Removed ' + member.name,
            message: member.name + ' was removed from ' + this.data.codebase.name,
          })
        })
        .finally(this.refresh)
    },
  },
})
</script>
