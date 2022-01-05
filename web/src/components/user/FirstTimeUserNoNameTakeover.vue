<template>
  <div class="max-w-3xl mx-auto mt-16 flex flex-col gap-2">
    <h1 class="text-gray-800 text-4xl font-bold">Welcome to Sturdy!</h1>
    <p class="text-gray-500">Allow me to introduce myself, I'm Sturdy. Who are you?</p>

    <div class="mt-1 flex rounded-md shadow-sm">
      <div class="relative flex items-stretch flex-grow focus-within:z-10">
        <input
          id="email"
          v-model="name"
          type="email"
          name="email"
          class="focus:ring-blue-500 focus:border-blue-500 block w-full rounded-none rounded-l-md sm:text-sm border-gray-300"
          placeholder="Enter your name"
          @keydown.enter="setName"
        />
      </div>
      <button
        type="button"
        class="-ml-px relative inline-flex items-center space-x-2 px-4 py-2 border border-gray-300 text-sm font-medium rounded-r-md text-yellow-700 bg-yellow-50 hover:bg-yellow-100 focus:outline-none focus:ring-1 focus:ring-blue-500 focus:border-blue-500 disabled:bg-gray-100 disabled:text-gray-500"
        :disabled="!name"
        @click="setName"
      >
        <LightningBoltSolidIcon v-if="name" class="h-5 w-5 text-yellow-500" aria-hidden="true" />
        <LightningBoltOutlineIcon v-else class="h-5 w-5 text-gray-500" aria-hidden="true" />
        <span>Get started!</span>
      </button>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent, PropType } from 'vue'
import { gql, useMutation } from '@urql/vue'
import { NoNameTakeoverUserFragment } from './__generated__/FirstTimeUserNoNameTakeover'
import { LightningBoltIcon as LightningBoltOutlineIcon } from '@heroicons/vue/outline'
import { LightningBoltIcon as LightningBoltSolidIcon } from '@heroicons/vue/solid'

export const NO_NAME_TAKEOVER_USER = gql`
  fragment NoNameTakeoverUser on User {
    id
    name
  }
`

export default defineComponent({
  components: { LightningBoltOutlineIcon, LightningBoltSolidIcon },
  props: {
    user: {
      type: Object as PropType<NoNameTakeoverUserFragment>,
      required: true,
    },
  },
  setup() {
    const { executeMutation: updateUserResult } = useMutation(gql`
      mutation NoNameTakeoverUserUpdateUser($name: String) {
        updateUser(input: { name: $name }) {
          id
          name
        }
      }
    `)

    return {
      async updateUser(name: string) {
        const variables = {
          name,
        }
        await updateUserResult(variables).then((result) => {
          if (result.error) {
            throw new Error(result.error)
          }
          console.log('update user', result)
        })
      },
    }
  },
  data() {
    return {
      name: '',
    }
  },
  methods: {
    setName() {
      this.updateUser(this.name)
    },
  },
})
</script>
