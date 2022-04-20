<template>
  <div v-if="!newWorkspace">
    When the codebase is set up on your local machine, a new draft change will be created for you to
    code in.
  </div>
  <div v-else>
    <RouterLinkButton
      :to="{
        name: 'workspaceHome',
        params: {
          codebaseSlug: slug,
          id: newWorkspace.id,
        },
      }"
    >
      Go to {{ newWorkspace.name }}
    </RouterLinkButton>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import type { PropType } from 'vue'
import { gql } from 'graphql-tag'
import type {
  GoToWorkspaceStepCodebaseFragment,
  SetupUserViewsFragment,
} from './__generated__/SetupSturdyGoToWorkspaceStep'
import RouterLinkButton from '../../atoms/RouterLinkButton.vue'
import { Slug } from '../../slug'

export const SETUP_USER_VIEWS = gql`
  fragment SetupUserViews on User {
    id
    views {
      workspace {
        id
        name
        codebase {
          id
          name
          shortID
        }
      }
    }
  }
`

export const GO_TO_WORKSPACE_STEP_CODEBASE_FRAGMENT = gql`
  fragment GoToWorkspaceStepCodebase on Codebase {
    id
    name
    shortID
  }
`

export default defineComponent({
  components: { RouterLinkButton },
  props: {
    codebase: {
      required: true,
      type: Object as PropType<GoToWorkspaceStepCodebaseFragment>,
    },
    user: {
      type: Object as PropType<SetupUserViewsFragment | undefined>,
      default: undefined,
    },
  },
  computed: {
    newWorkspace(): SetupUserViewsFragment['views'][number]['workspace'] {
      return this.user?.views.filter((view) => view.workspace?.codebase.id === this.codebase.id)[0]
        ?.workspace
    },
    slug(): string {
      return Slug(this.codebase.name, this.codebase.shortID)
    },
  },
})
</script>
