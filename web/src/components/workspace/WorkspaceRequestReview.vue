<template>
  <Listbox v-slot="{ open }" as="div">
    <div class="mt-1 relative">
      <ListboxButton
        class="relative w-full bg-white border border-gray-300 rounded-md shadow-sm pl-3 pr-10 py-2 text-left cursor-pointer focus:outline-none focus:ring-1 focus:ring-blue-500 focus:border-blue-500 sm:text-sm group"
        @click="loadPeople"
      >
        <span class="flex items-center">
          <span class="ml-3 block truncate text-gray-500 group-hover:text-gray-900">
            Ask for feedback
          </span>
        </span>
        <span class="ml-3 absolute inset-y-0 right-0 flex items-center pr-2 pointer-events-none">
          <SelectorIcon
            class="h-5 w-5 text-gray-400 group-hover:text-gray-800"
            aria-hidden="true"
          />
        </span>
      </ListboxButton>

      <transition
        leave-active-class="transition ease-in duration-100"
        leave-from-class="opacity-100"
        leave-to-class="opacity-0"
      >
        <ListboxOptions
          class="absolute z-10 mt-1 w-full bg-white shadow-lg max-h-56 rounded-md py-1 text-base ring-1 ring-black ring-opacity-5 overflow-auto focus:outline-none sm:text-sm"
        >
          <ListboxOption v-if="(!data || !people) && open" as="template">
            <li class="text-gray-500 cursor-default select-none relative py-2 pl-3 pr-9">
              Loading...
            </li>
          </ListboxOption>
          <ListboxOption v-else-if="people.length === 0 && open" as="template">
            <li class="text-gray-500 cursor-default select-none relative py-2 pl-3 pr-9">
              You're the only coder here. Invite someone?
            </li>
          </ListboxOption>
          <ListboxOption
            v-for="person in people"
            v-else
            :key="person.id"
            v-slot="{ active }"
            as="template"
            :value="person"
            @click="requestReview(person.id)"
          >
            <li
              :class="[
                active ? 'text-white bg-blue-600' : 'text-gray-900',
                'cursor-pointer select-none relative py-2 pl-3 pr-9',
              ]"
            >
              <div class="flex items-center">
                <Avatar :author="person" size="6" class="flex-shrink-0" />
                <span :class="['font-normal ml-3 block truncate']">
                  {{ person.name }}
                </span>
              </div>
            </li>
          </ListboxOption>
        </ListboxOptions>
      </transition>
    </div>
  </Listbox>
</template>

<script lang="ts">
import { defineComponent, ref, toRef } from 'vue'
import type { Ref } from 'vue'
import { Listbox, ListboxButton, ListboxOption, ListboxOptions } from '@headlessui/vue'
import { SelectorIcon } from '@heroicons/vue/solid'
import { gql, useQuery } from '@urql/vue'
import Avatar from '../../atoms/Avatar.vue'
import { useRequestReview } from '../../mutations/useRequestReview'
import type { WorkspaceRequestReviewCodebaseQuery } from './__generated__/WorkspaceRequestReview'

export default defineComponent({
  components: {
    Listbox,
    ListboxButton,
    ListboxOption,
    ListboxOptions,
    SelectorIcon,
    Avatar,
  },
  props: {
    codebaseId: {
      type: String,
      required: true,
    },
    workspaceId: {
      type: String,
      required: true,
    },
  },
  setup(props) {
    const pauseLoadPeople = ref(true)
    let codebaseId = toRef(props, 'codebaseId')

    let { data, fetching, error } = useQuery({
      query: gql`
        query WorkspaceRequestReviewCodebase($id: ID) {
          user {
            id
          }

          codebase(id: $id) {
            id
            members {
              id
              name
              avatarUrl
            }
          }
        }
      `,
      variables: {
        id: codebaseId,
      },
      requestPolicy: 'cache-and-network',
      pause: pauseLoadPeople,
    })

    const requestReviewResult = useRequestReview()

    return {
      data: data as Ref<WorkspaceRequestReviewCodebaseQuery>, //
      fetching,

      async loadPeople() {
        pauseLoadPeople.value = false
      },

      async requestReview(userID: string) {
        await requestReviewResult({ workspaceID: props.workspaceId, userID })
      },
    }
  },
  computed: {
    people: function () {
      if (this.data?.codebase?.members && this.data?.user) {
        return this.data.codebase?.members.filter((m) => m.id !== this.data.user.id)
      }
      return null
    },
  },
})
</script>
