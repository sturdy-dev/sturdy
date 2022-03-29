<template>
  <div class="flex flex-col gap-4">
    <div class="grid grid-flow-row-dense grid-cols-2 gap-4 xl:grid-cols-1">
      <Presence :workspace="workspace" :user="user" class="hidden xl:flex" />
      <Comments :workspace="workspace" />
      <UpdatedAt :workspace="workspace" />
      <BasedOn :workspace="workspace" :codebase-slug="codebaseSlug" />
      <Watching
        v-if="isAuthorized && !isSelfOwnedWorkspace"
        :user="user"
        :watchers="workspace.watchers"
        :workspace-id="workspace.id"
      />
      <StatusDetails :statuses="workspace.statuses" />
      <GitHubPullRequest
        :git-hub-integration="workspace.codebase?.gitHubIntegration"
        :git-hub-pull-request="workspace.gitHubPullRequest"
      />
    </div>

    <div class="flex gap-6 border-t border-gray-200 py-6 justify-evenly xl:flex-col">
      <div class="flex flex-col flex-1">
        <h2 class="text-sm font-medium text-gray-500">Author</h2>
        <ul role="list" class="mt-3 space-y-3">
          <li class="flex justify-start">
            <a href="#" class="flex items-center space-x-3">
              <div class="flex-shrink-0">
                <Avatar :author="workspace.author" size="5" />
              </div>
              <div class="text-sm font-medium text-gray-900">
                {{ workspace.author.name }}
              </div>
            </a>
          </li>
        </ul>
      </div>
      <WorkspaceApproval
        v-if="showApproval"
        class="flex-1"
        :reviews="workspace.reviews"
        :workspace="workspace"
        :codebase-id="workspace.codebase.id"
        :user="user"
        :members="workspace.codebase.members"
      />
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import type { PropType } from 'vue'
import { gql } from '@urql/vue'

import Presence, {
  WORKSPACE_FRAGMENT as WORKSPACE_PRESENSE_FRAGMENT,
} from '../components/workspace/Presence.vue'
import Comments, {
  WORKSPACE_FRAGMENT as COMMENTS_WORKSPACE_FRAGMENT,
} from '../components/workspace/details/Comments.vue'
import UpdatedAt, {
  WORKSPACE_FRAGMENT as UPDATED_WORKSPACE_FRAGMENT,
} from '../components/workspace/details/UpdatedAt.vue'
import BasedOn from '../components/workspace/details/BasedOn.vue'
import Watching, { WORKSPACE_WATCHER_FRAGMENT } from '../components/workspace/details/Watching.vue'
import StatusDetails from '../components/statuses/StatusDetails.vue'
import GitHubPullRequest, {
  CODEBASE_GITHUB_INTEGRATION_FRAGMENT,
  GITHUB_PULL_REQUEST_FRAGMENT,
} from '../components/workspace/details/GitHubPullRequest.vue'
import Avatar from '../atoms/Avatar.vue'
import { AUTHOR } from '../atoms/AvatarHelper'
import WorkspaceApproval from '../components/workspace/WorkspaceApproval.vue'

import { Slug } from '../slug'
import type { WorkspaceDetails_WorkspaceFragment } from './__generated__/WorkspaceDetails'

export const WORKSPACE_FRAGMENT = gql`
  fragment WorkspaceDetails_Workspace on Workspace {
    id
    author {
      ...Author
    }
    watchers {
      ...WorkspaceWatcher
    }
    gitHubPullRequest @include(if: $isGitHubEnabled) {
      ...GitHubPullRequest
    }
    headChange {
      id
      title
      author {
        ...Author
      }
    }
    codebase {
      id
      name
      shortID
      gitHubIntegration @include(if: $isGitHubEnabled) {
        ...CodebaseGitHubIntegration
      }
      members {
        ...Author
      }
    }
    reviews {
      id
      grade
      createdAt
      isReplaced
      dismissedAt
      author {
        id
        name
        avatarUrl
      }
    }
    suggestion {
      id
    }
    ...Presence_Workspace
    ...CommentsCount_Workspace
    ...Updated_Workspace
  }
  ${CODEBASE_GITHUB_INTEGRATION_FRAGMENT}
  ${GITHUB_PULL_REQUEST_FRAGMENT}
  ${WORKSPACE_WATCHER_FRAGMENT}
  ${WORKSPACE_PRESENSE_FRAGMENT}
  ${COMMENTS_WORKSPACE_FRAGMENT}
  ${UPDATED_WORKSPACE_FRAGMENT}
  ${AUTHOR}
`

type Author = WorkspaceDetails_WorkspaceFragment['author']

export default defineComponent({
  components: {
    Presence,
    Comments,
    UpdatedAt,
    BasedOn,
    Watching,
    StatusDetails,
    GitHubPullRequest,
    Avatar,
    WorkspaceApproval,
  },
  props: {
    workspace: {
      type: Object as PropType<WorkspaceDetails_WorkspaceFragment>,
      required: true,
    },
    user: {
      type: Object as PropType<Author>,
    },
  },
  computed: {
    codebaseSlug() {
      return Slug(this.workspace.codebase.name, this.workspace.codebase.shortID)
    },
    isAuthenticated() {
      return !!this.user
    },
    isAuthorized() {
      const isMember = this.workspace.codebase.members.some(({ id }) => id === this.user?.id)
      return this.isAuthenticated && isMember
    },
    isSelfOwnedWorkspace() {
      return this.workspace && this.workspace.author.id === this.user?.id
    },

    isSuggesting() {
      return !!this.workspace.suggestion
    },
    showApproval() {
      return !this.isSuggesting
    },
  },
})
</script>
