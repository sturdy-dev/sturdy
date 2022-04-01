<template>
  <div class="flex flex-col gap-4 divide-y divide-gray-200">
    <h2 id="activity-title" class="text-lg font-medium text-gray-900">Activity</h2>
    <div class="pt-4">
      <NewComment
        v-if="isAuthorized"
        :user="user"
        :members="change.codebase.members"
        :change-id="change.id"
      />
      <Activity
        :activity="change.activity"
        :codebase-slug="codebaseSlug"
        :user="user"
        :members="change.codebase.members"
      />
    </div>
  </div>
</template>

<script lang="ts">
import { gql } from '@urql/vue'
import type { PropType } from 'vue'
import { defineComponent } from 'vue'

import NewComment, { CODEBASE_FRAGMENT } from '../molecules/NewComment.vue'
import Activity, { WORKSPACE_ACTIVITY_FRAGMENT } from '../molecules/activity/Activity.vue'
import { MEMBER_FRAGMENT } from '../atoms/TextareaAutosize.vue'
import type { ChangeActivity_ChangeFragment } from './__generated__/ChangeActivitySidebar'

type Member = ChangeActivity_ChangeFragment['codebase']['members'][number]

export const CHANGE_FRAGMENT = gql`
  fragment ChangeActivity_Change on Change {
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
  components: { NewComment, Activity },
  props: {
    change: {
      type: Object as PropType<ChangeActivity_ChangeFragment>,
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
    isAuthenticated() {
      return !!this.user
    },
    isAuthorized() {
      const isMember = this.change.codebase.members.some(({ id }) => id === this.user?.id)
      return this.isAuthenticated && isMember
    },
  },
})
</script>
