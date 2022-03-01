<template>
  <div class="divide-y divide-gray-200">
    <div class="pb-4">
      <h2 id="activity-title" class="text-lg font-medium text-gray-900">Activity</h2>
    </div>
    <div class="pt-6">
      <NewComment
        v-if="isAuthorized"
        :user="user"
        :members="members"
        :workspace-id="workspaceId"
        :change-id="changeId"
      />
      <WorkspaceActivity
        :activity="activity"
        :codebase-slug="codebaseSlug"
        :user="user"
        :members="members"
      />
    </div>
  </div>
</template>
<script lang="ts">
import { gql } from '@urql/vue'
import { PropType } from 'vue'

import NewComment, { CODEBASE_FRAGMENT } from '../molecules/NewComment.vue'
import WorkspaceActivity, { WORKSPACE_ACTIVITY_FRAGMENT } from '../molecules/activity/Activity.vue'
import { WorkspaceActivity_WorkspaceFragment } from './__generated__/WorkspaceActivitySidebar'

type Activity = WorkspaceActivity_WorkspaceFragment['activity'][number]
type Member = WorkspaceActivity_WorkspaceFragment['codebase']['members'][number]

export const WORKSPACE_FRAGMENT = gql`
  fragment WorkspaceActivity_Workspace on Workspace {
    activity {
      ...WorkspaceActivity
    }
    codebase {
      ...NewComment
    }
  }
  ${CODEBASE_FRAGMENT}
  ${WORKSPACE_ACTIVITY_FRAGMENT}
`

export default {
  components: { NewComment, WorkspaceActivity },
  props: {
    workspaceId: { type: String },
    changeId: { type: String },
    codebaseSlug: { type: String, required: true },
    activity: { type: Array as PropType<Activity[]>, required: true },
    user: { type: Object as PropType<Member> },
    members: { type: Array as PropType<Member[]>, default: () => [] },
  },
  computed: {
    isAuthenticated() {
      return !!this.user
    },
    isAuthorized() {
      const isMember = this.members.some(({ id }) => id === this.user?.id)
      return this.isAuthenticated && isMember
    },
  },
}
</script>
