<template>
  <div>
    <div class="flex justify-between items-center">
      <h2
        v-if="isAuthorized || nonDismissedReviews.length > 0"
        class="text-sm font-medium text-gray-500"
      >
        Feedback
      </h2>

      <div v-if="isAuthorized && !isOwnWorkspace">
        <Button
          size="wider"
          class="focus:ring-0 focus:border-gray-300"
          :grouped="true"
          :first="true"
          :show-tooltip="true"
          @click="createOrUpdateReview(workspace.id, 'Approve')"
        >
          <template #tooltip>Looks good!</template>
          <template #default>
            <ThumbUpIcon
              class="h-4 w-4"
              :class="[selfUserReview?.grade === 'Approve' ? 'text-green-400' : 'text-gray-300']"
            />
          </template>
        </Button>
        <Button
          class="focus:ring-0 focus:border-gray-300"
          size="wider"
          :grouped="true"
          :last="true"
          :show-tooltip="true"
          :tooltip-right="true"
          @click="createOrUpdateReview(workspace.id, 'Reject')"
        >
          <template #tooltip>I have some feedback</template>
          <template #default>
            <InformationCircleIcon
              class="h-4 w-4"
              :class="[selfUserReview?.grade === 'Reject' ? 'text-orange-400' : 'text-gray-300']"
            />
          </template>
        </Button>
      </div>
    </div>

    <ul role="list" class="mt-3 space-y-3">
      <li v-for="(review, idx) in nonDismissedReviews" :key="idx" class="flex justify-start">
        <span class="flex items-center space-x-3">
          <div class="flex-shrink-0">
            <Avatar :author="review.author" size="5" />
          </div>

          <Tooltip>
            <template #tooltip>
              <span v-if="review.grade === 'Approve'">Looks good to me!</span>
              <span v-else-if="review.grade === 'Reject'">I have feedback</span>
              <span v-else-if="review.grade === 'Requested'">Waiting for feedback</span>
            </template>
            <template #default>
              <span
                class="relative inline-flex items-center rounded-full border border-gray-300 px-3 py-0.5"
              >
                <ThumbUpIcon
                  v-if="review.grade === 'Approve'"
                  class="h-5 w-5 text-green-400"
                  title="Approved"
                />
                <InformationCircleIcon
                  v-else-if="review.grade === 'Reject'"
                  e
                  class="h-5 w-5 text-orange-400"
                  title="Rejected"
                />
                <ClockIcon
                  v-else-if="review.grade === 'Requested'"
                  class="h-5 w-5 text-gray-300"
                  title="Pending review"
                />
              </span>
            </template>
          </Tooltip>

          <div class="text-sm font-medium text-gray-900">
            {{ review.author.name }}
          </div>
          <a v-if="isAuthorized" title="Dismiss this review">
            <XIcon
              class="h-3 w-3 text-gray-300 hover:text-gray-500 cursor-pointer"
              @click="dismissReview(review.id)"
            />
          </a>
        </span>
      </li>
    </ul>

    <WorkspaceRequestReview v-if="isAuthorized" :workspace="workspace" :user="user" class="mt-4" />
  </div>
</template>

<script lang="ts">
import Avatar from '../../atoms/Avatar.vue'
import { AUTHOR } from '../../atoms/AvatarHelper'
import Tooltip from '../../atoms/Tooltip.vue'
import Button from '../../atoms/Button.vue'
import { gql, useMutation } from '@urql/vue'
import { ClockIcon, InformationCircleIcon, ThumbUpIcon, XIcon } from '@heroicons/vue/solid'
import WorkspaceRequestReview, {
  WORKSPACE_FRAGMENT as WORKSPACE_REQUEST_REVIEW_WORKSPACE_FRAGMENT,
} from './WorkspaceRequestReview.vue'
import { useCreateOrUpdateReview } from '../../mutations/useCreateOrUpdateReview'
import { ReviewGrade } from '../../__generated__/types'
import type { PropType } from 'vue'
import type { WorkspaceApproval_WorkspaceFragment } from './__generated__/WorkspaceApproval'

export const WORKSPACE_FRAGMENT = gql`
  fragment WorkspaceApproval_Workspace on Workspace {
    id
    ...WorkspaceRequestReview_Workspace

    author {
      id
    }

    reviews {
      id
      grade
      createdAt
      isReplaced
      dismissedAt
      author {
        id
        ...Author
      }
    }

    codebase {
      id
      members {
        id
        ...Author
      }
    }
  }
  ${WORKSPACE_REQUEST_REVIEW_WORKSPACE_FRAGMENT}
  ${AUTHOR}
`

export default {
  components: {
    WorkspaceRequestReview,
    Avatar,
    Button,
    ThumbUpIcon,
    XIcon,
    ClockIcon,
    InformationCircleIcon,
    Tooltip,
  },
  props: {
    user: {
      type: Object as PropType<{ id: string }>,
      default: null,
    },
    workspace: {
      type: Object as PropType<WorkspaceApproval_WorkspaceFragment>,
      required: true,
    },
  },
  setup() {
    const createOrUpdateReviewResult = useCreateOrUpdateReview()

    const { executeMutation: dismissReviewResult } = useMutation(gql`
      mutation WorkspaceApprovalDismiss($id: ID!) {
        dismissReview(input: { id: $id }) {
          id
          dismissedAt
        }
      }
    `)

    return {
      async createOrUpdateReview(workspaceID: string, grade: ReviewGrade) {
        const variables = { workspaceID, grade }
        await createOrUpdateReviewResult(variables)
      },

      async dismissReview(id: string) {
        const variables = { id }
        await dismissReviewResult(variables).then((result) => {
          console.log('dismissReviewResult', result)
        })
      },
    }
  },
  data() {
    return {
      showRequestReview: false,
    }
  },
  computed: {
    isOwnWorkspace() {
      return this.user?.id === this.workspace.author.id
    },
    isAuthenticated() {
      return !!this.user
    },
    isAuthorized() {
      const isMember = this.workspace.codebase.members.some(({ id }) => id === this.user?.id)
      return this.isAuthenticated && isMember
    },
    selfUserReview() {
      if (!this.isAuthenticated) return null
      const r = this.workspace.reviews.filter(
        (r) => r.author.id === this.user?.id && !r.dismissedAt && !r.isReplaced
      )
      if (r && r.length > 0) {
        return r[0]
      }
      return null
    },
    nonDismissedReviews() {
      return this.workspace.reviews
        ?.filter((r) => !r.dismissedAt && !r.isReplaced)
        .sort((a, b) => a.author.name.localeCompare(b.author.name))
    },
  },
}
</script>
