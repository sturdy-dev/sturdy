<template>
  <div class="bg-white border overflow-hidden sm:rounded-lg">
    <div class="px-4 py-5 sm:px-6 flex justify-between">
      <div>
        <h3 class="text-lg leading-6 font-medium text-gray-900">... or import from Git?</h3>
        <p class="mt-1 max-w-2xl text-sm text-gray-500">
          Bring an existing repository and all of its history to Sturdy
        </p>
      </div>
      <div>
        <Button v-if="!showImportCommand" @click="show"> Show instructions</Button>
      </div>
    </div>
    <div v-if="showImportCommand" class="border-t border-gray-200 px-4 py-5 sm:px-6 text-sm">
      <p class="text-gray-500">
        Import your existing history to Sturdy using the
        <router-link to="/install" class="font-medium text-black"> Sturdy CLI</router-link>
      </p>

      <div class="mt-2">
        <div class="mt-1 flex rounded-md shadow-sm">
          <div class="relative flex items-stretch flex-grow focus-within:z-10">
            <input
              id="import-command"
              :value="sturdyImportCommand"
              type="text"
              readonly
              name="codebase_name"
              class="focus:ring-blue-500 focus:border-blue-500 block w-full rounded-none rounded-l-md sm:text-sm border-gray-300"
            />
          </div>
          <button
            class="-ml-px relative inline-flex items-center space-x-2 px-4 py-2 border border-gray-300 text-sm font-medium rounded-r-md text-gray-700 bg-gray-50 hover:bg-gray-100 focus:outline-none focus:ring-1 focus:ring-blue-500 focus:border-blue-500"
            @click="copy('import-command')"
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
            <span>Copy</span>
          </button>
        </div>
      </div>

      <div v-if="data?.codebase.changes?.length > 0" class="mt-2">
        <Banner message="Imported successfully!" status="success" />
        <p class="mt-2 flex items-center space-x-2">
          <span>The latest change in the codebase is now</span>

          <Avatar :author="data.codebase.changes[0].author" size="5" />

          <router-link
            class="font-medium text-black"
            :to="{
              name: 'codebaseChange',
              params: { id: codebaseId, initialSelectedChangeID: data.codebase.changes[0].id },
            }"
          >
            {{ data.codebase.changes[0].title }}
          </router-link>
        </p>
        <p class="mt-2">
          Follow the "First time setup guide" above to connect to the new codebase.
        </p>
      </div>
      <div v-else class="mt-8 flex">
        <p class="text-gray-500">Waiting for import ...</p>

        <svg
          class="animate-spin ml-3 -mr-1 h-5 w-5 text-black"
          xmlns="http://www.w3.org/2000/svg"
          fill="none"
          viewBox="0 0 24 24"
        >
          <circle
            class="opacity-25"
            cx="12"
            cy="12"
            r="10"
            stroke="currentColor"
            stroke-width="4"
          />
          <path
            class="opacity-75"
            fill="currentColor"
            d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
          />
        </svg>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { computed, defineComponent, onUnmounted, ref, toRef } from 'vue'
import Button from '../shared/Button.vue'
import { Banner } from '../../atoms'
import Avatar from '../shared/Avatar.vue'
import { gql, useQuery } from '@urql/vue'

export default defineComponent({
  components: { Button, Banner, Avatar },
  props: {
    codebaseId: String,
  },
  setup(props) {
    let codebaseID = toRef(props, 'codebaseId')
    let showImportCommand = ref(false)

    let { data, executeQuery } = useQuery({
      query: gql`
        query ImportFromGit($codebaseID: ID!) {
          codebase(id: $codebaseID) {
            id
            changes(input: { limit: 1 }) {
              id
              title
              author {
                id
                name
                avatarUrl
              }
            }
          }
        }
      `,
      variables: {
        codebaseID: codebaseID,
      },
      requestPolicy: 'network-only',
      pause: computed(() => !showImportCommand.value),
    })

    const interval = setInterval(() => {
      executeQuery({
        requestPolicy: 'network-only',
      })
    }, 2000)
    onUnmounted(() => {
      clearInterval(interval)
    })

    return {
      data,
      showImportCommand,
      async show() {
        showImportCommand.value = true
      },
    }
  },
  computed: {
    sturdyImportCommand(): string {
      return 'sturdy import ' + this.codebaseId
    },
  },
  methods: {
    copy(id: string) {
      let copyText = document.getElementById(id) as HTMLInputElement
      copyText.select()
      copyText.setSelectionRange(0, 99999)
      document.execCommand('copy')
    },
  },
})
</script>
