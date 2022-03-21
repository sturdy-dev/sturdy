<template>
  <h2 class="sr-only">Details</h2>
  <div class="space-y-5">
    <Presence
      :presence="data.workspace.presence"
      :workspace="data.workspace"
      :user="user"
      class="hidden lg:flex"
    />
    <Comments :data="data" />
    <UpdatedAt :data="data" />
    <BasedOn :data="data" :codebase-slug="codebaseSlug" />
    <Watching
      v-if="isAuthorized && !isSelfOwnedWorkspace"
      :user="user"
      :watchers="data.workspace.watchers"
      :workspace-id="data.workspace.id"
    />
    <StatusDetails :statuses="data.workspace.statuses" />
    <GitHubPullRequest
      :git-hub-integration="data?.workspace?.codebase?.gitHubIntegration"
      :git-hub-pull-request="data?.workspace?.gitHubPullRequest"
    />
  </div>
  <div class="mt-6 border-t border-gray-200 py-6 space-y-8">
    <div>
      <h2 class="text-sm font-medium text-gray-500">Author</h2>
      <ul role="list" class="mt-3 space-y-3">
        <li class="flex justify-start">
          <a href="#" class="flex items-center space-x-3">
            <div class="flex-shrink-0">
              <Avatar :author="data.workspace.author" size="5" />
            </div>
            <div class="text-sm font-medium text-gray-900">
              {{ data.workspace.author.name }}
            </div>
          </a>
        </li>
      </ul>
    </div>
    <WorkspaceApproval
      v-if="showApproval"
      :reviews="data.workspace.reviews"
      :workspace="data.workspace"
      :codebase-id="data.workspace.codebase.id"
      :user="user"
      :members="data.workspace.codebase.members"
    />
  </div>

  <WorkspaceActivitySidebar
    v-if="showActivity"
    class="mt-6"
    :workspace="data.workspace"
    :codebase-slug="codebaseSlug"
    :user="user"
  />
</template>

<script lang="ts">
import { defineComponent } from 'vue'
export default defineComponent({})
</script>
