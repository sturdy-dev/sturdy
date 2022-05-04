<template>
  <Listbox v-slot="{ open }" as="div">
    <div class="mt-1 relative">
      <ListboxButton
        class="relative w-full bg-white border border-gray-300 rounded-md shadow-sm pl-3 pr-10 py-2 text-left cursor-pointer focus:outline-none focus:ring-1 focus:ring-blue-500 focus:border-blue-500 sm:text-sm group"
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
          <ListboxOption v-if="people.length === 0 && open" as="template">
            <li class="text-gray-500 cursor-default select-none relative py-2 pl-3 pr-9">
              <template v-if="workspace.codebase.members.length === 1">
                You're the only coder here.
              </template>
              Invite someone else?
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
import { defineComponent, toRefs, type PropType } from 'vue'
import { Listbox, ListboxButton, ListboxOption, ListboxOptions } from '@headlessui/vue'
import { SelectorIcon } from '@heroicons/vue/solid'
import { gql } from '@urql/vue'
import Avatar from '../../atoms/Avatar.vue'
import { AUTHOR } from '../../atoms/AvatarHelper'
import { useRequestReview } from '../../mutations/useRequestReview'
import type { WorkspaceRequestReview_WorkspaceFragment } from './__generated__/WorkspaceRequestReview'

export const WORKSPACE_FRAGMENT = gql`
  fragment WorkspaceRequestReview_Workspace on Workspace {
    id
    author {
      id
    }
    codebase {
      id
      members {
        id
        ...Author
      }
    }
  }
  ${AUTHOR}
`

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
    user: {
      type: Object as PropType<{ id: string }>,
      required: false,
    },
    workspace: {
      type: Object as PropType<WorkspaceRequestReview_WorkspaceFragment>,
      required: true,
    },
  },
  setup(props) {
    const requestReviewResult = useRequestReview()
    const { workspace } = toRefs(props)

    return {
      async requestReview(userID: string) {
        await requestReviewResult({ workspaceID: workspace.value.id, userID })
      },
    }
  },
  computed: {
    people() {
      return this.workspace.codebase.members
        .filter(({ id }) => id !== this.user?.id)
        .filter(({ id }) => id !== this.workspace.author.id)
    },
  },
})
</script>
