<template>
  <div class="flex flex-col gap-4 divide-y divide-gray-200">
    <h2 id="activity-title" class="text-lg font-medium text-gray-900">Activity</h2>
    <div class="pt-4">
      <NewComment
        v-if="isAuthorized"
        :user="user"
        :members="workspace.codebase.members"
        :workspace-id="workspace.id"
      />
      <WorkspaceActivity
        :activity="workspace.activity"
        :codebase-slug="codebaseSlug"
        :user="user"
        :members="workspace.codebase.members"
      />
    </div>
  </div>
</template>
<script lang="ts">
import { gql } from '@urql/vue'
import { type PropType, defineComponent } from 'vue'

import NewComment, { CODEBASE_FRAGMENT } from '../molecules/NewComment.vue'
import WorkspaceActivity, { WORKSPACE_ACTIVITY_FRAGMENT } from '../molecules/activity/Activity.vue'
import { MEMBER_FRAGMENT } from '../atoms/TextareaAutosize.vue'
import type { WorkspaceActivity_WorkspaceFragment } from './__generated__/WorkspaceActivitySidebar'

type Member = WorkspaceActivity_WorkspaceFragment['codebase']['members'][number]

export const WORKSPACE_FRAGMENT = gql`
  fragment WorkspaceActivity_Workspace on Workspace {
    id
    activity {
      ...WorkspaceActivity
    }
    codebase {
      id
      ...NewComment
      members {
        ...Member
      }
    }
  }
  ${MEMBER_FRAGMENT}
  ${CODEBASE_FRAGMENT}
  ${WORKSPACE_ACTIVITY_FRAGMENT}
`

export default defineComponent({
  components: { NewComment, WorkspaceActivity },
  props: {
    workspace: {
      type: Object as PropType<WorkspaceActivity_WorkspaceFragment>,
      required: true,
    },
    codebaseSlug: {
      type: String,
      required: true,
    },
    user: {
      type: Object as PropType<Member>,
      required: false,
      default: null,
    },
  },
  computed: {
    isAuthenticated(): boolean {
      return !!this.user
    },
    isAuthorized(): boolean {
      const isMember = this.workspace.codebase.members.some(({ id }) => id === this.user?.id)
      return this.isAuthenticated && isMember
    },
  },
})
</script>
