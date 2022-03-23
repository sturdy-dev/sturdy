<template>
  <div>
    <h1 v-if="isSuggesting" class="text-2xl font-bold text-gray-900">
      Suggesting to {{ workspace.suggestion.for.name }}
    </h1>
    <span class="w-full text-2xl font-bold text-gray-900 border-0 p-0 border-0 outline-none">
      {{ name }}
    </span>
    <p class="mt-2 text-sm text-gray-500">
      By
      {{ ' ' }}
      <span class="font-medium text-gray-900">
        {{ workspace.author.name }}
      </span>
      {{ ' ' }}
      in
      {{ ' ' }}
      <router-link
        :to="{
          name: 'codebaseHome',
          params: { codebaseSlug: codebaseSlug },
        }"
        class="font-medium text-gray-900"
      >
        {{ workspace.codebase.name }}
      </router-link>
    </p>
  </div>
</template>

<script lang="ts">
import { defineComponent, toRefs } from 'vue'
import type { PropType } from 'vue'
import { gql } from '@urql/vue'
import type { WorkspaceName_WorkspaceFragment } from './__generated__/WorkspaceName'
import { Slug } from '../slug'
import { useUpdateWorkspace } from '../mutations/useUpdateWorkspace'

export const WORKSPACE_FRAGMENT = gql`
  fragment WorkspaceName_Workspace on Workspace {
    id
    name
    suggestion {
      id
      for {
        id
        name
      }
    }
    author {
      id
      name
    }
    codebase {
      id
      shortID
      name
    }
  }
`

export default defineComponent({
  props: {
    workspace: {
      type: Object as PropType<WorkspaceName_WorkspaceFragment>,
      required: true,
    },
    disabled: {
      type: Boolean,
    },
  },
  setup(props) {
    const { workspace } = toRefs(props)
    const { updateWorkspace, mutating: updating } = useUpdateWorkspace()
    return {
      updating,
      updateTitle(title: string) {
        updateWorkspace({
          id: workspace.value.id,
          name: title,
        })
      },
    }
  },
  data() {
    return {
      name: this.workspace.name,
      updateTitleTimeout: null as null | ReturnType<typeof setTimeout>,
    }
  },
  computed: {
    isSuggesting() {
      return !!this.workspace.suggestion
    },
    codebaseSlug() {
      return Slug(this.workspace.codebase.name, this.workspace.codebase.shortID)
    },
  },
  watch: {
    'workspace.name': function (newName) {
      this.name = newName
    },
  },
  methods: {
    onKeyDown(e: KeyboardEvent) {
      const target = e.target as HTMLInputElement
      if (target.value.length === 0) return

      if (this.updateTitleTimeout) clearTimeout(this.updateTitleTimeout)
      this.updateTitleTimeout = setTimeout(() => {
        this.updateTitle(target.value)
      }, 300)
    },
  },
})
</script>
