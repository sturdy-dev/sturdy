<template>
  <div class="min-h-16">
    <h1 v-if="isSuggesting" class="text-2xl font-bold text-gray-900">
      Suggesting to {{ workspace.suggestion.for.name }}
    </h1>
    <h1 v-else class="text-2xl font-bold text-gray-900">
      {{ workspace.name }}
    </h1>
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
import { defineComponent, PropType } from 'vue'
import { gql } from '@urql/vue'
import { WorkspaceName_WorkspaceFragment } from './__generated__/WorkspaceName'
import { Slug } from '../slug'

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
  },
  computed: {
    isSuggesting() {
      return !!this.workspace.suggestion
    },
    codebaseSlug() {
      return Slug(this.workspace.codebase.name, this.workspace.codebase.shortID)
    },
  },
})
</script>
