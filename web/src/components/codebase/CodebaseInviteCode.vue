<template>
  <div v-if="data">
    <template v-if="showHeader">
      <p class="block text-sm font-medium text-gray-700">Invite link</p>
      <p class="mt-1 text-sm text-gray-500">
        Invite someone to join the codebase by sending them a link
      </p>
    </template>

    <template v-if="data.codebase.inviteCode">
      <div class="mt-1 flex rounded-md shadow-sm">
        <div class="relative flex items-stretch flex-grow focus-within:z-10">
          <input
            id="invite_code"
            :value="inviteURL"
            type="text"
            readonly
            class="focus:ring-blue-500 focus:border-blue-500 block w-full rounded-none rounded-l-md sm:text-sm border-gray-300"
          />
        </div>
        <button
          class="-ml-px relative inline-flex items-center space-x-2 py-2 border border-gray-300 text-sm font-medium rounded-r-md text-gray-700 bg-gray-50 hover:bg-gray-100 focus:outline-none focus:ring-1 focus:ring-blue-500 focus:border-blue-500"
          :class="[small ? 'px-2' : 'px-4']"
          @click="copyInviteCode"
        >
          <svg
            class="h-5 w-5 text-gray-400"
            xmlns="http://www.w3.org/2000/svg"
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M8 5H6a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2v-1M8 5a2 2 0 002 2h2a2 2 0 002-2M8 5a2 2 0 012-2h2a2 2 0 012 2m0 0h2a2 2 0 012 2v3m2 4H10m0 0l3-3m-3 3l3 3"
            />
          </svg>
          <span v-if="!small">Copy</span>
        </button>
      </div>
      <div class="mt-1 z-0">
        <Button :grouped="true" :first="true" @click="updateInviteCode"> Refresh link </Button>
        <Button :grouped="true" :last="true" @click="disableInviteCode"> Disable link </Button>
      </div>
    </template>
    <div v-else class="mt-1">
      <Button @click="updateInviteCode"> Create invite link </Button>
    </div>
  </div>
</template>

<script>
import { gql, useMutation, useQuery } from '@urql/vue'
import Button from '../shared/Button.vue'

export default {
  name: 'CodebaseInviteCode',
  components: { Button },
  props: ['codebaseID', 'showHeader', 'small'],
  setup(props) {
    let { data, executeQuery } = useQuery({
      query: gql`
        query CodebaseInviteCode($id: ID, $shortID: ID) {
          codebase(id: $id, shortID: $shortID) {
            id
            inviteCode
          }
        }
      `,
      variables: {
        id: props.codebaseID,
      },
      requestPolicy: 'cache-and-network',
    })

    const { executeMutation: generateInviteCodeResult } = useMutation(gql`
      mutation CodebaseInviteCodeGenerate($id: ID!) {
        updateCodebase(input: { id: $id, generateInviteCode: true }) {
          id
          inviteCode
        }
      }
    `)

    const { executeMutation: disableInviteCodeResult } = useMutation(gql`
      mutation CodebaseInviteCodeDisable($id: ID!) {
        updateCodebase(input: { id: $id, disableInviteCode: true }) {
          id
          inviteCode
        }
      }
    `)

    return {
      data,
      refresh: () => {
        executeQuery({
          requestPolicy: 'network-only',
        })
      },

      async doGenerateInviteCode(id) {
        const variables = { id }
        await generateInviteCodeResult(variables).then((result) => {
          if (result.error) {
            throw new Error(result.error)
          }
        })
      },

      async doDisableInviteCode(id) {
        const variables = { id }
        await disableInviteCodeResult(variables).then((result) => {
          if (result.error) {
            throw new Error(result.error)
          }
        })
      },
    }
  },
  computed: {
    inviteURL() {
      return `${location.origin}/join/${this.data.codebase.inviteCode}`
    },
  },
  methods: {
    updateInviteCode() {
      this.doGenerateInviteCode(this.data.codebase.id).catch(() => {
        this.updateStatus = 'Something went wrong.'
      })
    },
    disableInviteCode() {
      this.doDisableInviteCode(this.data.codebase.id).catch(() => {
        this.updateStatus = 'Something went wrong.'
      })
    },
    copyInviteCode() {
      let copyText = document.getElementById('invite_code')
      copyText.select()
      copyText.setSelectionRange(0, 99999)
      document.execCommand('copy')
    },
  },
}
</script>
