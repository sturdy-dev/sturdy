<template>
  <div class="divide-y divide-gray-200">
    <div class="pb-4">
      <h2 id="activity-title" class="text-lg font-medium text-gray-900">Activity</h2>
    </div>
    <div class="pt-6">
      <NewComment
        v-if="isAuthorized"
        :user="user"
        :members="workspace.codebase.members"
        :workspace-id="workspace.id"
        :change-id="changeId"
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
import { PropType } from 'vue'

import NewComment, { CODEBASE_FRAGMENT } from '../molecules/NewComment.vue'
import WorkspaceActivity, { WORKSPACE_ACTIVITY_FRAGMENT } from '../molecules/activity/Activity.vue'
import { MEMBER_FRAGMENT } from '../components/shared/TextareaAutosize.vue'
import { WorkspaceActivity_WorkspaceFragment } from './__generated__/WorkspaceActivitySidebar'

type Member = WorkspaceActivity_WorkspaceFragment['codebase']['members'][number]

export const WORKSPACE_FRAGMENT = gql`
  fragment WorkspaceActivity_Workspace on Workspace {
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

export default {
  components: { NewComment, WorkspaceActivity },
  props: {
    workspace: {
      type: Object as PropType<WorkspaceActivity_WorkspaceFragment>,
      required: true,
    },
    changeId: { type: String },
    codebaseSlug: { type: String, required: true },
    user: { type: Object as PropType<Member> },
  },
  computed: {
    isAuthenticated() {
      return !!this.user
    },
    isAuthorized() {
      const isMember = this.workspace.codebase.members.some(({ id }) => id === this.user?.id)
      return this.isAuthenticated && isMember
    },
  },
}
</script>
