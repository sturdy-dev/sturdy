<template>
  <main v-if="data" class="bg-white">
    <div
      class="md:pr-64 fixed z-50 w-full"
      :style="ipc ? 'top: calc(env(titlebar-area-height, 2rem) + 1px)' : 'top: 0'"
    >
      <SelectedHunksToolbar />
      <SearchToolbar />
    </div>

    <ArchiveWorkspaceModal
      :is-active="archiveWorkspaceActive"
      :workspace-i-d="data.workspace.id"
      @deletedWorkspace="onWorkspaceArchived"
      @closeDeleteWorkspace="hideArchiveModal"
    />

    <div class="py-8">
      <div class="mx-auto px-6 grid grid-cols-1 xl:grid-cols-4">
        <div class="xl:col-span-3 xl:pr-8 xl:border-r xl:border-gray-200">
          <div>
            <div>
              <div class="md:flex md:items-center md:justify-between md:space-x-4">
                <!-- Workspace name -->
                <div v-if="editingName" class="h-16 inline-flex flex-row items-center">
                  <div
                    class="inline-flex rounded-md shadow-sm mr-4"
                    tabindex="0"
                    @focusout="saveName"
                  >
                    <div class="relative flex items-stretch focus-within:z-10">
                      <input
                        ref="workspaceName"
                        v-model="userEditingName"
                        type="text"
                        placeholder="Name your workspace, so that you know what you're working on"
                        style="min-width: 400px"
                        class="focus:ring-blue-500 focus:border-blue-500 block w-full rounded-none rounded-l-md sm:text-sm border-gray-300"
                        @keydown="editingNameKeyDown"
                      />
                    </div>
                    <button
                      class="-ml-px relative inline-flex items-center space-x-2 px-4 py-2 border border-gray-300 text-sm font-medium rounded-r-md text-gray-700 bg-gray-50 hover:bg-gray-100 focus:outline-none focus:ring-1 focus:ring-blue-500 focus:border-blue-500"
                      @click="saveName"
                    >
                      <span>Save</span>
                    </button>
                  </div>
                </div>
                <div v-else class="min-h-16">
                  <h1 v-if="isSuggesting" class="text-2xl font-bold text-gray-900">
                    Suggesting to {{ data.workspace.suggestion.for.name }}
                  </h1>
                  <h1 v-else class="text-2xl font-bold text-gray-900">
                    {{ data.workspace.name }}
                  </h1>
                  <p class="mt-2 text-sm text-gray-500">
                    By
                    {{ ' ' }}
                    <span class="font-medium text-gray-900">
                      {{ data.workspace.author.name }}
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
                      {{ data.workspace.codebase.name }}
                    </router-link>
                  </p>
                </div>

                <div v-if="isAuthorized" class="mt-4 flex space-x-3 md:mt-0 items-center">
                  <Button
                    v-if="showEdit"
                    size="wider"
                    :disabled="editingName"
                    @click="startEditingName"
                  >
                    <PencilIcon class="-ml-1 mr-2 h-5 w-5 text-gray-400" aria-hidden="true" />
                    <span>Edit</span>
                  </Button>

                  <Button size="wider" @click="showArchiveModal">
                    <ArchiveIcon class="-ml-1 mr-2 h-5 w-5 text-gray-400" aria-hidden="true" />
                    <span>Archive</span>
                  </Button>

                  <div v-if="showSync">
                    <OnboardingStep
                      id="SyncChanges"
                      :enabled="!data.workspace.upToDateWithTrunk && !!mutableView"
                    >
                      <template #title>Get up to date</template>
                      <template #description>
                        Other workspaces have published changes, and this workspace has fallen
                        behind. Before any changes in this workspace can be published, it needs to
                        be synchronized with the changelog. Try it out!
                      </template>

                      <Tooltip :disabled="isSyncing" x-direction="left">
                        <template #tooltip>
                          <div v-if="data.workspace.upToDateWithTrunk">
                            This workspace is already up-to-date with the changelog.
                          </div>
                          <div v-else-if="viewConnectionState !== 'editing'">
                            You need to connect a local directory to this workspace before syncing.
                          </div>
                          <div v-else>
                            Get all the latest changes from the changelog into this workspace.
                          </div>
                        </template>

                        <div class="relative inline-flex rounded-md shadow-sm">
                          <Button
                            :disabled="
                              !isOnAuthoritativeView ||
                              isSyncing ||
                              data.workspace.upToDateWithTrunk ||
                              viewConnectionState !== 'editing'
                            "
                            size="wider"
                            @click="initSyncWithTrunk"
                          >
                            <LightningBoltIcon
                              class="-ml-1 mr-2 h-5 w-5 text-gray-400"
                              aria-hidden="true"
                            />
                            <div v-if="isSyncing" class="inline-flex space-x-2 items-center">
                              <span>Syncing</span>
                              <Spinner />
                            </div>
                            <template v-else>
                              <span class="hidden md:block">Sync changes</span>
                              <span class="block md:hidden">Sync</span>
                            </template>
                          </Button>
                          <span
                            v-if="!data.workspace.upToDateWithTrunk && !isSyncing"
                            class="flex absolute h-3 w-3 top-0 right-0 -mt-1 -mr-1"
                          >
                            <span
                              class="animate-ping absolute inline-flex h-full w-full rounded-full bg-blue-400 opacity-75"
                            />
                            <span class="relative inline-flex rounded-full h-3 w-3 bg-blue-500" />
                          </span>
                        </div>
                      </Tooltip>
                    </OnboardingStep>
                  </div>
                </div>
              </div>

              <!-- Connect Button -->
              <div
                v-if="
                  (viewConnectionState === 'own' || viewConnectionState === 'others') &&
                  mostRecentSelfUserView
                "
                class="mt-5"
              >
                <ButtonWithDropdown
                  v-if="viewConnectionState === 'own'"
                  :disabled="loadingNewWorkspace"
                  class="z-20"
                  @click="openWorkspaceOnView(data.workspace.id, mostRecentSelfUserView.id)"
                >
                  <div class="flex">
                    <DesktopComputerIcon
                      class="-ml-1 mr-2 h-5 w-5 text-gray-400"
                      aria-hidden="true"
                    />
                    Connect {{ mostRecentSelfUserView.shortMountPath }} for editing
                  </div>
                  <template v-if="connectedViews.length > 1 || mutagenIpc != null" #dropdown>
                    <MenuItem v-for="view in connectedViews" :key="view.id">
                      <button
                        class="text-sm text-left py-2 px-4 flex hover:bg-gray-50"
                        :disabled="loadingNewWorkspace"
                        @click="openWorkspaceOnView(data.workspace.id, view.id)"
                      >
                        <ViewStatusIndicator v-if="view" class="pr-1" :view="view" compact />
                        <span class="font-medium pr-1">{{ view.shortMountPath }}</span>
                        <span class="text-gray-500"> on {{ view.mountHostname }}</span>
                      </button>
                    </MenuItem>

                    <MenuItem v-if="mutagenIpc != null">
                      <button
                        class="text-sm text-left py-2 px-4 flex hover:bg-gray-50"
                        @click="createViewInDirectory"
                      >
                        <FolderAddIcon class="h-5 w-5 pr-0.5 mx-0.5" />
                        <span class="font-medium">Connect another directory</span>
                      </button>
                    </MenuItem>
                  </template>
                </ButtonWithDropdown>

                <OnboardingStep v-else-if="diffs.length > 0" id="OpenForSuggesting">
                  <template #title>Help {{ data.workspace.author.name }} out</template>

                  <template #description>
                    You've made it to someone else's workspace. If you want to give feedback in
                    code, you can temporarily connect your local directory to this workspace and
                    make changes locally. These changes will appear as suggestions to
                    {{ data.workspace.author.name }}.
                  </template>

                  <ButtonWithDropdown
                    color="green"
                    :disabled="loadingNewWorkspace"
                    class="z-20"
                    @click="createSuggestion(data.workspace.id, mostRecentSelfUserView.id)"
                  >
                    <div class="flex">
                      <AnnotationIcon
                        class="-ml-1 mr-2 h-5 w-5 text-green-700"
                        aria-hidden="true"
                      />
                      Connect {{ mostRecentSelfUserView.shortMountPath }} for suggesting
                    </div>

                    <template v-if="connectedViews.length > 1" #dropdown>
                      <MenuItem v-for="view in connectedViews" :key="view.id">
                        <button
                          class="text-sm text-left py-2 px-4 flex hover:bg-gray-50"
                          :disabled="loadingNewWorkspace"
                          @click="createSuggestion(data.workspace.id, view.id)"
                        >
                          <ViewStatusIndicator v-if="view" class="pr-1" :view="view" compact />
                          <span class="font-medium pr-1">{{ view.shortMountPath }}</span>
                          <span class="text-gray-500"> on {{ view.mountHostname }}</span>
                        </button>
                      </MenuItem>
                    </template>
                  </ButtonWithDropdown>
                </OnboardingStep>
              </div>

              <div
                v-if="
                  viewConnectionState === 'own' && mutagenIpc != null && connectedViews.length === 0
                "
                class="mt-5"
              >
                <Button @click="createViewInDirectory">
                  <div class="flex items-center px-1">
                    <DesktopComputerIcon
                      class="-ml-1 mr-2 h-5 w-5 text-gray-400"
                      aria-hidden="true"
                    />
                    Connect directory
                  </div>
                </Button>
              </div>

              <!-- View Connection Status -->
              <div
                v-if="displayView && viewConnectionState === 'editing'"
                class="mt-4 flex space-x-8"
              >
                <ViewStatusIndicator :view="displayView" />
                <OpenInEditor :view="displayView" />
              </div>

              <aside class="mt-8 xl:hidden">
                <h2 class="sr-only">Details</h2>
                <div class="space-y-5 md:flex md:space-y-0">
                  <div class="space-y-5 md:flex-1">
                    <Comments :data="data" />
                    <UpdatedAt :data="data" />
                  </div>
                  <div class="space-y-5 md:flex-1">
                    <BasedOn :data="data" :codebase-slug="codebaseSlug" />
                    <Watching
                      v-if="user && !isSelfOwnedWorkspace"
                      :user="user"
                      :watchers="data.workspace.watchers"
                      :workspace-id="data.workspace.id"
                    />
                    <GitHubPullRequest
                      :git-hub-integration="data?.workspace?.codebase?.gitHubIntegration"
                      :git-hub-pull-request="data?.workspace?.gitHubPullRequest"
                    />
                  </div>
                </div>
                <div
                  class="mt-6 border-t border-b border-gray-200 py-6 space-y-8 md:space-y-0 md:grid md:grid-cols-2"
                >
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
              </aside>
              <div v-if="showDescription" class="pt-3 relative max-w-prose">
                <h2 class="sr-only">Description</h2>
                <OnboardingStep
                  id="MakingAChange"
                  :dependencies="['FindingYourWorkspace', 'WorkspaceChanges']"
                >
                  <template #title>Publishing a Change</template>
                  <template #description>
                    When you've made edits to your files and feel ready to make a checkpoint, write
                    a description of your change(s) here.
                  </template>

                  <Editor
                    :model-value="workspace_draft_description"
                    :editable="isAuthorized"
                    placeholder="Describe the changes in this workspace&hellip;"
                    @updated="onUpdatedDescription"
                  >
                    <transition
                      enter-active-class="transition ease-out duration-75"
                      enter-from-class="opacity-0 scale-75"
                      enter-to-class="opacity-100 scale-100"
                      leave-active-class="transition ease-in duration-75"
                      leave-from-class="opacity-100 scale-100"
                      leave-to-class="opacity-0 scale-75"
                    >
                      <ShareButton
                        v-if="data.workspace && canSubmitChange"
                        :workspace="data.workspace"
                        :all-hunk-ids="diffs.flatMap((diff) => diff.hunks.map((hunk) => hunk.id))"
                        @pre-create-change="preCreateChange"
                      />
                    </transition>
                  </Editor>
                </OnboardingStep>

                <transition
                  enter-active-class="transition ease-out duration-50"
                  enter-from-class="opacity-0 scale-75"
                  enter-to-class="opacity-100 scale-100"
                  leave-active-class="transition ease-in duration-25"
                  leave-from-class="opacity-100 scale-100"
                  leave-to-class="opacity-0 scale-75"
                >
                  <div
                    v-if="justSaved"
                    class="hidden xl:block text-gray-400 text-sm absolute bottom-full translate-y-1 right-0 origin-bottom-right"
                  >
                    Saved
                  </div>
                </transition>
              </div>
            </div>
          </div>

          <section>
            <div class="flex-grow z-10 relative min-w-0">
              <div class="pt-4">
                <!-- Rebasing -->
                <Banner
                  v-if="rebasing_complete_had_conflicts"
                  status="success"
                  class="mb-2"
                  message="Sync completed! Good job solving all of those conflicts!"
                />
                <Banner
                  v-if="rebasing_complete_no_conflicts"
                  status="success"
                  class="mb-2"
                  message="Sync completed! You're now up to date again!"
                />
                <Banner
                  v-if="rebasing_failed"
                  status="error"
                  class="mb-2"
                  message="Syncing failed. Please try again."
                />
                <Banner v-if="rebasing_working" status="info" message="Syncing..." class="mb-2" />

                <div v-if="rebasing?.is_rebasing && rebasing.conflicting_files">
                  <ResolveConflict
                    :rebasing="rebasing"
                    :conflict-diffs="conflictDiffs"
                    @resolveConflict="resolveConflict"
                  />
                  <div class="my-8">
                    <p
                      v-if="
                        rebasing.conflicting_files.length !== rebasing_conflict_resolutions.size
                      "
                      class="text-sm text-gray-500 pb-4 text-center"
                    >
                      You have
                      {{ rebasing.conflicting_files.length - rebasing_conflict_resolutions.size }}
                      unresolved conflict{{
                        rebasing.conflicting_files.length - rebasing_conflict_resolutions.size > 1
                          ? 's'
                          : ''
                      }}
                    </p>
                    <div class="flex justify-center">
                      <Button
                        :disabled="
                          rebasing.conflicting_files.length !== rebasing_conflict_resolutions.size
                        "
                        @click="sendConflictResolution"
                      >
                        Done
                      </Button>
                    </div>
                  </div>
                </div>

                <LiveDetails
                  v-else
                  :view="displayView"
                  :workspace="data.workspace"
                  :mutable="!!mutableView"
                  :is-suggesting="isSuggesting"
                  :diffs="diffs"
                  :comments="nonArchivedComments"
                  :user="user"
                  :members="data.workspace.codebase.members"
                  :is-on-authoritative-view="isOnAuthoritativeView"
                  :loaded-diffs="loadedDiffs"
                  @codebase-updated="refresh"
                  @pre-create-change="preCreateChange"
                />
              </div>
            </div>
          </section>
        </div>
        <aside class="hidden xl:block xl:pl-8">
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

          <div v-if="showActivity" class="mt-6 py-6 space-y-8">
            <div>
              <div class="divide-y divide-gray-200">
                <div class="pb-4">
                  <h2 id="activity-title" class="text-lg font-medium text-gray-900">Activity</h2>
                </div>
                <div class="pt-6">
                  <WorkspaceNewComment
                    v-if="isAuthorized"
                    :user="user"
                    :members="data.workspace.codebase.members"
                    :workspace-id="data.workspace.id"
                  />
                  <WorkspaceActivity
                    :activity="data.workspace.activity"
                    :codebase-slug="codebaseSlug"
                    :user="user"
                    :members="data.workspace.codebase.members"
                  />
                </div>
              </div>
            </div>
          </div>
        </aside>
      </div>
    </div>
  </main>
</template>

<script lang="ts">
import {
  AnnotationIcon,
  ArchiveIcon,
  DesktopComputerIcon,
  LightningBoltIcon,
  PencilIcon,
} from '@heroicons/vue/solid'
import { FolderAddIcon } from '@heroicons/vue/outline'
import LiveDetails, { LIVE_DETAILS_WORKSPACE } from '../components/workspace/LiveDetails.vue'
import { MEMBER_FRAGMENT } from '../components/shared/TextareaMentions.vue'
import http from '../http'
import { Banner } from '../atoms'
import ResolveConflict from '../components/workspace/ResolveConflict.vue'
import Button from '../components/shared/Button.vue'
import debounce from '../debounce'
import { gql, useMutation, useQuery } from '@urql/vue'
import { useRoute, useRouter } from 'vue-router'
import {
  computed,
  defineAsyncComponent,
  onUnmounted,
  ref,
  watch,
  toRefs,
  defineComponent,
  inject,
  Ref,
} from 'vue'
import { useHead } from '@vueuse/head'
import Spinner from '../components/shared/Spinner.vue'
import Avatar from '../components/shared/Avatar.vue'
import WorkspaceActivity, {
  WORKSPACE_ACTIVITY_FRAGMENT,
} from '../components/workspace/activity/Activity.vue'
import WorkspaceNewComment from '../components/workspace/WorkspaceNewComment.vue'
import WorkspaceApproval from '../components/workspace/WorkspaceApproval.vue'
import Watching, { WORKSPACE_WATCHER_FRAGMENT } from '../components/workspace/details/Watching.vue'
import Presence, { PRESENCE_FRAGMENT_QUERY } from '../components/workspace/Presence.vue'
import ArchiveWorkspaceModal from '../components/codebase/ArchiveWorkspaceModal.vue'
import GitHubPullRequest, {
  CODEBASE_GITHUB_INTEGRATION_FRAGMENT,
  GITHUB_PULL_REQUEST_FRAGMENT,
} from '../components/workspace/details/GitHubPullRequest.vue'
import Comments from '../components/workspace/details/Comments.vue'
import UpdatedAt from '../components/workspace/details/UpdatedAt.vue'
import BasedOn from '../components/workspace/details/BasedOn.vue'
import StatusDetails from '../components/statuses/StatusDetails.vue'
import { useUpdatedComment } from '../subscriptions/useUpdatedComment'
import { useUpdatedWorkspace } from '../subscriptions/useUpdatedWorkspace'
import OnboardingStep from '../components/onboarding/OnboardingStep.vue'
import ViewStatusIndicator, { VIEW_STATUS_INDICATOR } from '../components/ViewStatusIndicator.vue'
import ButtonWithDropdown from '../components/shared/ButtonWithDropdown.vue'
import { MenuItem } from '@headlessui/vue'
import { ViewStatusState } from '../__generated__/types'
import { useOpenWorkspaceOnView } from '../mutations/useOpenWorkspaceOnView'
import Tooltip from '../components/shared/Tooltip.vue'
import { useUpdatedWorkspaceWatchers } from '../subscriptions/useUpdatedWorkspaceWathcers'
import { useCreateSuggestion } from '../mutations/useCreateSuggestion'
import SelectedHunksToolbar from '../components/workspace/SelectedHunksToolbar.vue'
import SearchToolbar from '../components/workspace/SearchToolbar.vue'
import ShareButton, { SHARE_BUTTON } from '../components/workspace/ShareButton.vue'
import OpenInEditor from '../components/workspace/OpenInEditor.vue'
import { Feature } from '../__generated__/types'

export default defineComponent({
  components: {
    ShareButton,
    SelectedHunksToolbar,
    Tooltip,
    FolderAddIcon,
    ButtonWithDropdown,
    ViewStatusIndicator,
    OnboardingStep,
    WorkspaceApproval,
    Watching,
    WorkspaceNewComment,
    WorkspaceActivity,
    Avatar,
    Spinner,
    Button,
    ResolveConflict,
    LiveDetails,
    AnnotationIcon,
    Banner,
    Editor: defineAsyncComponent(() => import('../components/workspace/Editor.vue')),
    PencilIcon,
    ArchiveIcon,
    LightningBoltIcon,
    ArchiveWorkspaceModal,
    GitHubPullRequest,
    Comments,
    UpdatedAt,
    BasedOn,
    Presence,
    StatusDetails,
    DesktopComputerIcon,
    MenuItem,
    SearchToolbar,
    OpenInEditor,
  },
  props: {
    user: {
      type: Object,
    },
  },
  emits: ['workspaceUpdated', 'codebase-updated'],
  setup() {
    const features = inject<Ref<Array<Feature>>>('features', ref([]))
    const isGitHubEnabled = computed(() => features?.value?.includes(Feature.GitHub))

    let route = useRoute()
    const router = useRouter()
    let workspaceID = ref(route.params.id)
    let shortCodebaseID = ref(route.params.codebaseSlug)
    watch(
      () => route.params.id,
      (slug) => {
        workspaceID.value = slug
      }
    )
    watch(
      () => route.params.codebaseSlug,
      (slug) => {
        shortCodebaseID.value = slug
      }
    )

    const ViewFragment = gql`
      fragment ViewParts on View {
        id
        shortMountPath
        mountHostname
        lastUsedAt
        author {
          id
          name
          avatarUrl
        }
        workspace {
          id
        }
      }
    `

    let { data, fetching, error, executeQuery } = useQuery({
      query: gql`
        query WorkspaceHome($workspaceID: ID!, $isGitHubEnabled: Boolean!) {
          workspace(id: $workspaceID) {
            id
            name
            createdAt
            lastLandedAt
            updatedAt
            upToDateWithTrunk
            lastActivityAt
            draftDescription
            author {
              id
              name
              avatarUrl
            }
            watchers {
              ...WorkspaceWatcher
            }
            suggestion {
              id
              for {
                id
                name
              }
            }
            headChange {
              id
              title
              trunkCommitID
              createdAt
              author {
                id
                name
                avatarUrl
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
            activity {
              ...WorkspaceActivity
            }
            view {
              id
              mountPath
              ...ViewParts
              ...ViewStatusIndicator
            }
            comments {
              id
              message
              codeContext {
                id
                lineStart
                lineEnd
                lineIsNew
                context
                contextStartsAtLine
                path
              }
              createdAt
              deletedAt
              author {
                id
                name
                avatarUrl
              }
              replies {
                id
                message
                createdAt
                author {
                  id
                  name
                  avatarUrl
                }
              }
            }
            gitHubPullRequest @include(if: $isGitHubEnabled) {
              ...GitHubPullRequest
            }
            codebase {
              id
              name

              gitHubIntegration @include(if: $isGitHubEnabled) {
                ...CodebaseGitHubIntegration
              }
              members {
                ...Member
              }

              workspaces {
                id
                name
                lastActivityAt
                author {
                  id
                  name
                  avatarUrl
                }
              }
              views {
                id
                shortMountPath
                mountHostname
                lastUsedAt
                author {
                  id
                  name
                  avatarUrl
                }
                workspace {
                  id
                  name
                }
                status {
                  id
                  state
                }
                ...ViewStatusIndicator
              }
            }
            presence {
              ...PresenceParts
            }
            ...LiveDetailsWorkspace
            ...ShareButton
          }
        }

        ${ViewFragment}
        ${PRESENCE_FRAGMENT_QUERY}
        ${WORKSPACE_ACTIVITY_FRAGMENT}
        ${CODEBASE_GITHUB_INTEGRATION_FRAGMENT}
        ${GITHUB_PULL_REQUEST_FRAGMENT}
        ${VIEW_STATUS_INDICATOR}
        ${LIVE_DETAILS_WORKSPACE}
        ${MEMBER_FRAGMENT}
        ${WORKSPACE_WATCHER_FRAGMENT}
        ${SHARE_BUTTON}
      `,
      variables: { workspaceID: workspaceID, isGitHubEnabled },
    })

    useHead({
      title: computed(() => {
        let n = data.value?.workspace?.name
        if (n) {
          return n + ' | Sturdy'
        }
        return 'Sturdy'
      }),
    })

    useUpdatedWorkspace(workspaceID, {
      pause: computed(() => !shortCodebaseID.value || !workspaceID.value),
    })

    const openWorkspaceOnViewResult = useOpenWorkspaceOnView()

    const { executeMutation: updateWorkspaceResult } = useMutation(gql`
      mutation WorkspaceHomeUpdate($workspaceID: ID!, $name: String, $draftDescription: String) {
        updateWorkspace(
          input: { id: $workspaceID, name: $name, draftDescription: $draftDescription }
        ) {
          id
          name
          draftDescription
        }
      }
    `)

    let displayView = ref(null)
    let displayViewId = ref(null)

    watch(data, () => {
      if (data.value?.workspace?.view) {
        displayView.value = data.value?.workspace?.view
        displayViewId.value = displayView.value?.id
      }
    })

    useUpdatedWorkspaceWatchers(workspaceID)
    useUpdatedComment(workspaceID, displayViewId)

    // Re-run the main query every 15s
    let refreshInterval = setInterval(() => {
      executeQuery({
        requestPolicy: 'network-only',
      })
    }, 15000)
    onUnmounted(() => {
      clearInterval(refreshInterval)
    })

    // Workspace not found
    watch(error, (err) => {
      if (
        err?.graphQLErrors &&
        err?.graphQLErrors.length > 0 &&
        err?.graphQLErrors[0]?.message === 'NotFoundError'
      ) {
        router.push({
          name: 'codebaseHome',
          params: { codebaseSlug: route.params.codebaseSlug },
        })
      }
    })

    const selectedHunkIDs = ref(new Set())

    const createSuggestionResult = useCreateSuggestion()
    const loadingNewWorkspace = ref(false)
    return {
      fetching,
      data,
      error,
      displayView,
      loadingNewWorkspace,

      openWorkspaceOnViewResult,

      mutagenIpc: window.mutagenIpc,
      ipc: window.ipc,

      selectedHunkIDs,

      async refresh() {
        await executeQuery({
          requestPolicy: 'network-only',
        })
      },

      async createSuggestion(workspaceID, viewID) {
        const suggestion = await createSuggestionResult({
          workspaceID: workspaceID,
        })
        loadingNewWorkspace.value = true
        router.push({
          name: 'workspaceHome',
          params: {
            id: suggestion.workspace.id,
            codebaseSlug: route.params.codebaseSlug,
          },
        })
        await openWorkspaceOnViewResult({
          workspaceID: suggestion.workspace.id,
          viewID,
        })
        loadingNewWorkspace.value = false
      },

      async updateWorkspace(workspaceID, name = null, draftDescription = null) {
        const variables = { workspaceID, name, draftDescription }
        await updateWorkspaceResult(variables).then((result) => {
          console.log('update workspace', result)
        })
      },
    }
  },
  data() {
    return this.initialState()
  },
  computed: {
    showApproval() {
      return !this.isSuggesting
    },
    showActivity() {
      return !this.isSuggesting
    },
    showEdit() {
      return !this.isSuggesting
    },
    showSync() {
      return (
        (!this.data.workspace || this.data.workspace.author.id === this.user?.id) &&
        !this.isSuggesting
      )
    },
    showDescription() {
      return this.diffs.length > 0 && !this.isSuggesting
    },
    isAuthenticated() {
      return !!this.user
    },
    isAuthorized() {
      const isMember = this.data.workspace.codebase.members.some(({ id }) => id === this.user?.id)
      return this.isAuthenticated && isMember
    },
    mutableView() {
      return (
        this.data.workspace && this.data.workspace.author.id === this.user?.id && this.displayView
      )
    },
    isOnAuthoritativeView() {
      return this.data?.workspace?.view?.id === this.displayView?.id
    },
    isSelfOwnedWorkspace() {
      return this.data?.workspace && this.data?.workspace.author.id === this.user?.id
    },
    mostRecentSelfUserView() {
      return this.views[0]
    },
    views() {
      return (
        this.data?.workspace.codebase.views?.slice().sort((a, b) => {
          return b.lastUsedAt - a.lastUsedAt
        }) ?? []
      )
    },
    connectedViews() {
      return this.views.filter(
        (v) => v.status != null && v.status.state !== ViewStatusState.Disconnected
      )
    },
    nonArchivedComments() {
      return this.data?.workspace.comments.filter((c) => !c.deletedAt)
    },
    isSuggesting() {
      return !!this.data?.workspace.suggestion
    },
    viewConnectionState() {
      if (this.isSelfOwnedWorkspace) {
        if (this.views.some((v) => v.workspace?.id === this.data.workspace.id)) {
          return 'editing'
        } else {
          return 'own'
        }
      } else {
        return 'others'
      }
    },
    canSubmitChange() {
      if (this.data?.workspace == null) {
        return false
      }
      if (this.diffs.length === 0) {
        return false
      }
      // Have to have a change description before sharing
      if (this.data.workspace.draftDescription.length === 0) {
        return false
      }
      // Disallow users from sharing when they have selected hunks
      // (since it might lead them to think they're doing a partial share)
      if (this.selectedHunkIDs.size > 0) {
        return false
      }
      // If the workspace is up to date, we know for a fact that it doesn't conflict (it's a cheaper check, so mark as shareable right away)
      if (this.data.workspace.upToDateWithTrunk) {
        return true
      }
      return true
    },
  },
  watch: {
    '$route.params.id': function () {
      if (this.$route.params.id) {
        this.reset()
      }
    },
    'data.workspace.codebase.id': function (n) {
      if (n) this.emitter.emit('codebase', n)
    },
    'data.workspace.id': function (n) {
      if (n) {
        this.subscribe()
      }
    },
    'data.workspace.draftDescription': function () {
      this.setDraftDescription()
    },
    'displayView.id': function () {
      this.subscribe()
      if (this.displayView && this.displayView.id) {
        this.fetchRebasingStatus(this.displayView.id)
      }
    },
  },
  unmounted() {
    this.unsubscribe()
    this.emitter.off('differ-selected-hunk-ids', this.onSelectedHunkIDs)
  },
  mounted() {
    this.subscribe()
    this.setDraftDescription()
    this.emitter.on('differ-selected-hunk-ids', this.onSelectedHunkIDs)
  },
  methods: {
    initialState() {
      return {
        workspaceID: this.$route.params.id,
        codebaseSlug: this.$route.params.codebaseSlug,

        workspace_draft_description: null, // v-model (populated from data.workspace.draftDescription)
        workspace_draft_description_last_saved_val: null,
        updatedWorkspaceDescriptionDebounce: debounce(this.saveDraftDescription, 800),
        justSaved: false,
        unsetJustSavedFunc: debounce(() => {
          this.justSaved = false
        }, 4000),

        rebasing: {}, // Rebase status
        rebasing_complete_no_conflicts: false,
        rebasing_complete_had_conflicts: false,
        rebasing_failed: false,
        rebasing_working: false,
        rebasing_conflict_resolutions: new Map(),

        eventStream: null,
        eventPingTimeout: null,
        diffs: [],
        loadedDiffs: false,

        conflictDiffs: [],

        changes: [],
        reactive_changes_cancel_func: null,

        editingName: false, // if the name is being edited
        userEditingName: '', // model for the new name

        loadingNewWorkspace: false,
        isSyncing: false,

        archiveWorkspaceActive: false,
      }
    },
    openWorkspaceOnView(workspaceID, viewID) {
      const variables = { workspaceID, viewID }
      this.loadingNewWorkspace = true
      this.openWorkspaceOnViewResult(variables)
        .catch((e) => {
          const badRequestErrors = e.graphQLErrors.filter(
            (err) => err.message === 'BadRequestError'
          )
          const isBadRequest = badRequestErrors.length > 0
          if (isBadRequest) {
            const errorMessage = badRequestErrors[0].extensions.message
            this.emitter.emit('notification', {
              title: 'Can’t open workspace',
              message: `${errorMessage}`,
              style: 'error',
            })
          } else {
            throw e
          }
        })
        .finally(() => {
          this.loadingNewWorkspace = false
        })
    },
    onSelectedHunkIDs(hunkIDs) {
      this.selectedHunkIDs = hunkIDs
    },
    reset() {
      this.unsubscribe()
      Object.assign(this.$data, this.initialState())
      this.subscribe()

      this.setDraftDescription()
    },

    setDraftDescription() {
      this.workspace_draft_description = this.data?.workspace?.draftDescription
    },

    subscribe() {
      // Disconnect before starting a new connection
      if (this.eventStream) {
        this.unsubscribe()
      }

      if (!this.data?.workspace?.id) {
        return
      }

      const params = new URLSearchParams({
        workspace_id: this.data.workspace.id,
      })
      if (this.displayView && this.displayView.id) {
        params.append('view_id', this.displayView.id)
      }

      // Including the workspaceId in the URL to make sure that this view is using this workspace
      // If not, don't return any diff.
      let es = new EventSource(http.url('v3/stream?' + params.toString()), {
        withCredentials: true,
      })

      // Register timeout
      this.eventPingTimeout = setTimeout(this.eventTimeoutHandler, 25000)

      es.addEventListener(
        'error',
        (event) => {
          if (event.readyState === EventSource.CLOSED) {
            console.log('Event was closed')
            console.log(EventSource)
          }
        },
        false
      )

      // Ping
      // Diffs
      // CodebaseUpdated
      // ConflictDiffs
      // WorkspaceComments
      // WorkspaceUpdated

      es.addEventListener('Ping', (event) => {
        // reset timeout handler
        if (this.eventPingTimeout) {
          clearTimeout(this.eventPingTimeout)
        }
        // set a new timeout
        this.eventPingTimeout = setTimeout(this.eventTimeoutHandler, 25000) // The server sends a ping every 10s. Reconnect if one hasn't been received in 25s.
      })

      es.addEventListener('CodebaseUpdated', this.$emit('codebase-updated', {}))
      es.addEventListener('Diffs', (event) => this.handleDiffsEvent(JSON.parse(event.data)))
      es.addEventListener('ConflictDiffs', (event) =>
        this.handleConflictDiffsEvent(JSON.parse(event.data))
      )

      this.eventStream = es
    },
    unsubscribe() {
      clearTimeout(this.eventPingTimeout)
      if (this.eventStream) {
        this.eventStream.close()
        this.eventStream = null
      }
      // reset data
      this.diffs = []
      this.loadedDiffs = false
      this.conflictDiffs = []
    },
    eventTimeoutHandler() {
      this.unsubscribe()
      this.subscribe()
    },

    async saveName() {
      await this.updateWorkspace(this.data.workspace.id, this.userEditingName)
      this.editingName = false
    },
    async saveDraftDescription() {
      // Deduplication: avoid making requests if the description has not changed.
      // This is used to make sure that the client does not send a request to update the description after or during a
      // a change is being created from the workspace.
      if (this.workspace_draft_description_last_saved_val === this.workspace_draft_description) {
        return
      }

      this.workspace_draft_description_last_saved_val = this.workspace_draft_description
      this.justSaved = true
      this.unsetJustSavedFunc()
      await this.updateWorkspace(this.data.workspace.id, null, this.workspace_draft_description)
    },

    initSyncWithTrunk() {
      this.syncWithTrunk()
    },

    syncWithTrunk() {
      this.isSyncing = true

      let req = {
        workspace_id: this.data.workspace.id,
      }

      fetch(http.url('v3/rebase/' + this.displayView.id + '/start'), {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(req),
        credentials: 'include',
      })
        .then(http.checkStatus)
        .then((response) => response.json())
        .then((data) => {
          this.rebasing = data
          this.rebasing_conflict_resolutions = new Map()

          // Rebase is complete
          if (!data.is_rebasing) {
            this.rebasing_complete_no_conflicts = true
            setTimeout(() => {
              this.rebasing_complete_no_conflicts = false
            }, 3000)

            this.completedSync()
          }
        })
        .catch((e) => {
          console.log(e)
          this.rebasing_failed = true
          setTimeout(() => {
            this.rebasing_failed = false
          }, 3000)
        })
        .finally(() => {
          this.isSyncing = false
        })
    },

    fetchRebasingStatus(viewID) {
      fetch(http.url('v3/rebase/' + viewID), {
        credentials: 'include',
      })
        .then(http.checkStatus)
        .then((response) => response.json())
        .then((data) => {
          this.rebasing = data
        })
        .catch((e) => {
          console.log(e)
        })
    },

    resolveConflict(event) {
      let conflict = event.conflictingFile
      let version = event.version

      if (version !== 'todo') {
        this.rebasing_conflict_resolutions.set(conflict.path, {
          file_path: conflict.path,
          version: version,
        })
      } else {
        this.rebasing_conflict_resolutions.delete(conflict.path)
      }
    },

    sendConflictResolution() {
      let f = []
      this.rebasing_conflict_resolutions.forEach((v) => {
        f.push(v)
      })

      this.rebasing_working = true
      this.rebasing = null

      fetch(http.url('v3/rebase/' + this.displayView.id + '/resolve'), {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          files: f,
          workspace_id: this.data.workspace.id, // TODO: verify this value on the backend
        }),
        credentials: 'include',
      })
        .then(http.checkStatus)
        .then((response) => response.json())
        .then((data) => {
          this.rebasing = data
          this.rebasing_conflict_resolutions = new Map()

          // Request is done
          this.rebasing_working = false

          // Rebase is complete
          if (!data.is_rebasing) {
            this.rebasing_complete_had_conflicts = true
            setTimeout(() => {
              this.rebasing_complete_had_conflicts = false
            }, 3000)
            this.completedSync()
          }
        })
        .catch((e) => {
          console.log(e)
          alert('resolve failed')
        })
    },

    completedSync() {
      this.refresh()
    },

    startEditingName() {
      if (this.editingName) {
        this.editingName = false
        return
      }

      this.editingName = true
      this.userEditingName = this.data.workspace.name

      // Focus the input field
      this.$nextTick(() => {
        this.$refs.workspaceName.focus()
      })
    },
    editingNameKeyDown(e) {
      // Stop bubbling, to not trigger the key handler in LiveDetails (which captures events like Cmd+Enter, and Cmd+A)
      e.stopPropagation()

      // Enter
      if (e.keyCode === 13) {
        e.preventDefault()
        this.saveName()
      }
    },

    onUpdatedDescription(ev) {
      this.workspace_draft_description = ev.content
      this.justSaved = false
      this.unsetJustSavedFunc()

      if (ev.shouldSaveImmediately) {
        this.saveDraftDescription()
      } else {
        this.updatedWorkspaceDescriptionDebounce(ev.isInteractiveUpdate)
      }
    },
    handleDiffsEvent(ev) {
      this.loadedDiffs = true

      if (ev.diffs == null) {
        this.diffs = []
        return
      }

      this.diffs = ev.diffs
    },
    handleConflictDiffsEvent(ev) {
      if (ev.diffs == null) {
        this.conflictDiffs = []
        return
      }
      this.conflictDiffs = ev.diffs.filter((diff) => !diff.is_hidden)
    },

    showArchiveModal() {
      this.archiveWorkspaceActive = true
    },

    hideArchiveModal() {
      this.archiveWorkspaceActive = false
    },

    onWorkspaceArchived() {
      this.$router.push({
        name: 'codebaseHome',
        params: { codebaseSlug: this.$route.params.codebaseSlug },
      })
    },

    async preCreateChange() {
      // Save the description right away
      await this.saveDraftDescription()
    },

    async createViewInDirectory() {
      if (this.data?.workspace?.id == null) {
        return
      }

      let oldIsReady = this.mutagenIpc?.isReady && (await this.mutagenIpc.isReady())
      let newIsReady = this.ipc?.state && (await this.ipc.state()) === 'online'

      let mutagenReady = oldIsReady || newIsReady

      if (!mutagenReady) {
        this.emitter.emit('notification', {
          title: 'Sturdy is not running',
          message: 'Sturdy is still starting, please wait.',
          style: 'error',
        })
        return
      }

      try {
        await this.mutagenIpc.createNewViewWithDialog(this.data.workspace.id)
      } catch (e) {
        if (e.message.includes('non-empty')) {
          this.emitter.emit('notification', {
            title: 'Directory is not empty',
            message: 'Please select an empty directory.',
            style: 'error',
          })
        } else if (e.message.includes('Cancelled')) {
          // User cancelled the dialog. Do nothing
        } else {
          throw e
        }
      }
    },
  },
})
</script>
