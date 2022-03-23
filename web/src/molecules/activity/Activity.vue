<template>
  <div class="flow-root">
    <ul role="list" class="-mb-8">
      <li v-for="(item, itemIdx) in activity" :key="item.id">
        <div class="relative pb-8">
          <span
            v-if="itemIdx !== activity.length - 1"
            class="absolute top-5 left-5 -ml-px h-full w-0.5 bg-gray-200"
            aria-hidden="true"
          />

          <WorkspaceActivityComment
            v-if="item.__typename === 'WorkspaceCommentActivity'"
            :item="item"
            :codebase-slug="codebaseSlug"
            :user="user"
          />
          <WorkspaceActivityCreatedChange
            v-else-if="item.__typename === 'WorkspaceCreatedChangeActivity'"
            :item="item"
            :codebase-slug="codebaseSlug"
          />
          <WorkspaceActivityRequestedReview
            v-else-if="item.__typename === 'WorkspaceRequestedReviewActivity'"
            :item="item"
            :codebase-slug="codebaseSlug"
          />
          <WorkspaceActivityReviewed
            v-else-if="item.__typename === 'WorkspaceReviewedActivity'"
            :item="item"
            :codebase-slug="codebaseSlug"
          />
        </div>
      </li>
    </ul>
  </div>
</template>

<script lang="ts">
import WorkspaceActivityCreatedChange, {
  WORKSPACE_ACTIVITY_CREATED_CHANGE_FRAGMENT,
} from './ActivityCreatedChange.vue'
import WorkspaceActivityRequestedReview, {
  WORKSPACE_ACTIVITY_REQUESTED_REVIEW_FRAGMENT,
} from './ActivityRequestedReview.vue'
import WorkspaceActivityComment, {
  WORKSPACE_ACTIVITY_COMMENT_FRAGMENT,
} from './ActivityComment.vue'
import WorkspaceActivityReviewed, {
  WORKSPACE_ACTIVITY_REVIEWED_FRAGMENT,
} from './ActivityReviewed.vue'
import { defineComponent, onUnmounted, toRefs, watch } from 'vue'
import type { PropType } from 'vue'
import { gql, useMutation } from '@urql/vue'
import type {
  WorkspaceActivityCodebaseMemberFragment,
  WorkspaceActivityFragment,
} from './__generated__/Activity'

export const WORKSPACE_ACTIVITY_FRAGMENT = gql`
  fragment WorkspaceActivity on WorkspaceActivity {
    __typename
    id
    createdAt
    author {
      id
      name
      avatarUrl
    }
    ... on WorkspaceRequestedReviewActivity {
      ...WorkspaceActivityRequestedReview
    }
    ... on WorkspaceCreatedChangeActivity {
      ...WorkspaceCreatedChangeActivity
    }
    ... on WorkspaceReviewedActivity {
      ...WorkspaceReviewedActivity
    }
    ... on WorkspaceCommentActivity {
      ...WorkspaceCommentActivity
    }
  }
  ${WORKSPACE_ACTIVITY_REQUESTED_REVIEW_FRAGMENT}
  ${WORKSPACE_ACTIVITY_CREATED_CHANGE_FRAGMENT}
  ${WORKSPACE_ACTIVITY_REVIEWED_FRAGMENT}
  ${WORKSPACE_ACTIVITY_COMMENT_FRAGMENT}
`

export const WORKSPACE_ACTIVITY_CODEBASE_MEMBER_FRAGMENT = gql`
  fragment WorkspaceActivityCodebaseMember on User {
    id
    name
    avatarUrl
  }
`

export default defineComponent({
  components: {
    WorkspaceActivityComment,
    WorkspaceActivityCreatedChange,
    WorkspaceActivityRequestedReview,
    WorkspaceActivityReviewed,
  },
  props: {
    user: {
      type: Object,
    },
    activity: {
      type: Object as PropType<Array<WorkspaceActivityFragment>>,
      required: true,
    },
    codebaseSlug: {
      type: String,
      required: true,
    },
    members: {
      type: Array as PropType<Array<WorkspaceActivityCodebaseMemberFragment>>,
      required: true,
    },
  },
  setup(props) {
    const { activity, user, members } = toRefs(props)
    const isAuthenticated = !!user.value
    const isMember = members.value.some(({ id }) => id === user.value?.id)
    const isAuthorized = isAuthenticated && isMember

    const { executeMutation: readWorkspaceActivity } = useMutation(gql`
      mutation WorkspaceActivityRead($id: ID!) {
        readWorkspaceActivity(input: { id: $id }) {
          id
          isRead
        }
      }
    `)

    let isVisible = document.visibilityState === 'visible'
    let hasFocus = document.hasFocus()

    let markAsRead = () => {
      if (!isAuthorized) return

      if (!isVisible || !hasFocus) {
        return
      }

      if (activity.value && activity.value.length > 0) {
        readWorkspaceActivity({ id: activity.value[0].id }).then((result) => {
          if (result.error) {
            throw new Error(result.error.toString())
          }
        })
      }
    }

    let visibilityListener = () => {
      isVisible = document.visibilityState === 'visible'
      hasFocus = document.hasFocus()

      // Mark as unread in 1s (if still visible)
      if (isVisible && hasFocus) {
        setTimeout(markAsRead, 1000)
      }
    }

    document.addEventListener('visibilitychange', visibilityListener)
    window.addEventListener('focus', visibilityListener)
    window.addEventListener('blur', visibilityListener)
    onUnmounted(() => {
      document.removeEventListener('visibilitychange', visibilityListener)
      window.removeEventListener('focus', visibilityListener)
      window.removeEventListener('blur', visibilityListener)
    })

    watch(activity, (n, o) => {
      if (n && o) {
        if (!n.length || !o.length) {
          markAsRead()
        } else if (n.length > 0 && o.length > 0 && n[0].id !== o[0].id) {
          markAsRead()
        }
      } else if (n && !o && n.length > 0) {
        markAsRead()
      }
    })

    // On load, after 1s, mark the newest activity as read
    setTimeout(markAsRead, 1000)

    return {}
  },
})
</script>
