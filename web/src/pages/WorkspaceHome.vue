<template>
  <main v-if="data" class="bg-white">
    <div class="md:pr-64 fixed z-50 w-full">
      <SelectedHunksToolbar />
      <SearchToolbar />
    </div>

    <div class="py-8">
      <div class="mx-auto px-6 grid grid-cols-1 xl:grid-cols-4">
        <div class="xl:col-span-3 xl:pr-8 xl:border-r xl:border-gray-200">
          <div class="flex flex-col gap-2">
            <div class="flex justify-between gap-4">
              <WorkspaceName
                class="grow text-ellipsis overflow-hidden"
                :workspace="data.workspace"
                :disabled="!isAuthorized"
              />
              <div class="flex items-start gap-2">
                <ArchiveButton v-if="isAuthorized" :workspace-id="data.workspace.id" />
                <!-- sync button -->
                <div v-if="isAuthorized && showSync">
                  <OnboardingStep
                    id="SyncChanges"
                    :enabled="!data.workspace.upToDateWithTrunk && !!mutableView"
                  >
                    <template #title>Get up to date</template>
                    <template #description>
                      The codebase have new changes since this draft was started, and it's fallen
                      behind. Sync this draft to download all of the new changes.
                    </template>

                    <Tooltip :disabled="isSyncing" x-direction="left">
                      <template #tooltip>
                        <div v-if="data.workspace.upToDateWithTrunk">
                          This draft change is already up-to-date with the changelog.
                        </div>
                        <div v-else-if="viewConnectionState !== 'editing'">
                          You need to connect a local directory to this draft change before syncing.
                        </div>
                        <div v-else>
                          Get all the latest changes from the changelog into this draft change.
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
                          :icon="lightningBoltIcon"
                          :spinner="isSyncing"
                          @click="initSyncWithTrunk"
                        >
                          <span v-if="isSyncing">Syncing</span>
                          <span v-else>Sync</span>
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
                <!-- end sync button -->
              </div>
            </div>
            <!-- View Connection Status -->
            <div
              v-if="displayView && viewConnectionState === 'editing'"
              class="flex items-center gap-2"
            >
              <ViewStatusIndicator :view="displayView" />
              <OpenInEditor :view="displayView" />
            </div>
            <!-- end Connection Status -->

            <!-- Connect Button -->
            <template
              v-if="
                (viewConnectionState === 'own' || viewConnectionState === 'others') &&
                mostRecentSelfUserView
              "
            >
              <ButtonWithDropdown
                v-if="viewConnectionState === 'own'"
                :disabled="loadingNewWorkspace"
                class="z-20"
                @click="openWorkspaceOnView(data.workspace.id, mostRecentSelfUserView.id)"
              >
                <div class="inline-flex items-center gap-1">
                  <DesktopComputerIcon
                    class="hidden sm:block -ml-1 h-5 w-5 text-gray-400"
                    aria-hidden="true"
                  />
                  <span>Open in {{ mostRecentSelfUserView.shortMountPath }} for editing</span>
                </div>
                <template v-if="connectedViews.length > 1 || mutagenAvailable" #dropdown>
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

                  <MenuItem v-if="mutagenAvailable">
                    <button
                      class="text-sm text-left py-2 px-4 flex hover:bg-gray-50"
                      @click="createViewInDirectory"
                    >
                      <FolderAddIcon class="h-5 w-5 pr-0.5 mx-0.5" />
                      <span class="font-medium">Add another directory</span>
                    </button>
                  </MenuItem>
                </template>
              </ButtonWithDropdown>

              <OnboardingStep v-else-if="diffs.length > 0" id="OpenForSuggesting">
                <template #title>Help {{ data.workspace.author.name }} out</template>

                <template #description>
                  You're looking at someone else's draft change. If you want to give feedback in
                  code, you can temporarily connect your local directory to this draft and make
                  changes locally. These changes will appear as suggestions to
                  {{ data.workspace.author.name }}.
                </template>

                <ButtonWithDropdown
                  color="green"
                  :disabled="loadingNewWorkspace"
                  class="z-20"
                  @click="createSuggestion(data.workspace.id, mostRecentSelfUserView.id)"
                >
                  <div class="inline-flex gap-1">
                    <AnnotationIcon
                      class="hidden sm:block -ml-1 h-5 w-5 text-green-700"
                      aria-hidden="true"
                    />
                    <span>Open in {{ mostRecentSelfUserView.shortMountPath }} for suggesting</span>
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
            </template>
            <!-- end connect button -->

            <div class="flex items-center">
              <ConnectNewDirectory
                v-if="
                  viewConnectionState === 'own' && mutagenAvailable && connectedViews.length === 0
                "
                :codebase="data.workspace.codebase"
              />
            </div>

            <WorkspaceDescription
              class="max-w-prose"
              :workspace="data.workspace"
              :user="user"
              :diff-ids="diffs.flatMap((diff) => diff.hunks.map((hunk) => hunk.id))"
              :selected-hunk-ids="selectedHunkIDs"
            />

            <WorkspaceDetails class="mt-8 xl:hidden" :workspace="data.workspace" :user="user" />
          </div>

          <section>
            <div class="flex-grow relative min-w-0">
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

                <div v-if="rebaseStatus?.isRebasing && rebaseStatus.conflictingFiles">
                  <ResolveConflict
                    :rebasing="rebaseStatus"
                    :conflict-diffs="diffs"
                    @resolve-conflict="resolveConflict"
                  />
                  <div class="my-8">
                    <p
                      v-if="
                        rebaseStatus.conflictingFiles.length !== rebasing_conflict_resolutions.size
                      "
                      class="text-sm text-gray-500 pb-4 text-center"
                    >
                      You have
                      {{
                        rebaseStatus.conflictingFiles.length - rebasing_conflict_resolutions.size
                      }}
                      unresolved conflict{{
                        rebaseStatus.conflictingFiles.length - rebasing_conflict_resolutions.size >
                        1
                          ? 's'
                          : ''
                      }}
                    </p>
                    <div class="flex justify-center">
                      <Button
                        :disabled="
                          rebaseStatus.conflictingFiles.length !==
                          rebasing_conflict_resolutions.size
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
                  :is-stale="diffsStale"
                  :is-fetching="diffsFetching"
                  @codebase-updated="refresh"
                />
              </div>
            </div>
          </section>
        </div>
        <aside class="hidden xl:block xl:pl-8">
          <WorkspaceDetails :workspace="data.workspace" :user="user" />

          <WorkspaceActivitySidebar
            v-if="showActivity"
            class="mt-6"
            :workspace="data.workspace"
            :codebase-slug="codebaseSlug"
            :user="user"
          />
        </aside>
      </div>
    </div>
  </main>
</template>

<script lang="ts">
import { AnnotationIcon, DesktopComputerIcon, LightningBoltIcon } from '@heroicons/vue/solid'
import { FolderAddIcon } from '@heroicons/vue/outline'
import LiveDetails, {
  LIVE_DETAILS_DIFFS,
  LIVE_DETAILS_WORKSPACE,
} from '../components/workspace/LiveDetails.vue'
import http from '../http'
import { Banner } from '../atoms'
import ResolveConflict, { RESOLVE_CONFLICT_DIFF } from '../components/workspace/ResolveConflict.vue'
import Button from '../atoms/Button.vue'
import { gql, useQuery } from '@urql/vue'
import { useRoute, useRouter } from 'vue-router'
import { computed, defineComponent, inject, onUnmounted, ref, watch } from 'vue'
import type { Ref } from 'vue'
import { useHead } from '@vueuse/head'
import WorkspaceActivitySidebar, {
  WORKSPACE_FRAGMENT as WORKSPACE_ACTIVITY_WORKSPACE_FRAGMENT,
} from '../organisms/WorkspaceActivitySidebar.vue'
import ArchiveButton from '../organisms/workspace/ArchiveButton.vue'
import { useUpdatedComment } from '../subscriptions/useUpdatedComment'
import { useUpdatedWorkspace } from '../subscriptions/useUpdatedWorkspace'
import OnboardingStep from '../components/onboarding/OnboardingStep.vue'
import ViewStatusIndicator, { VIEW_STATUS_INDICATOR } from '../components/ViewStatusIndicator.vue'
import ButtonWithDropdown from '../molecules/ButtonWithDropdown.vue'
import { MenuItem } from '@headlessui/vue'
import { Feature, ViewStatusState } from '../__generated__/types'
import { useOpenWorkspaceOnView } from '../mutations/useOpenWorkspaceOnView'
import Tooltip from '../atoms/Tooltip.vue'
import { useUpdatedWorkspaceWatchers } from '../subscriptions/useUpdatedWorkspaceWathcers'
import { useCreateSuggestion } from '../mutations/useCreateSuggestion'
import SelectedHunksToolbar from '../components/workspace/SelectedHunksToolbar.vue'
import SearchToolbar from '../components/workspace/SearchToolbar.vue'
import OpenInEditor from '../components/workspace/OpenInEditor.vue'
import type {
  WorkspaceHomeDiffsQuery,
  WorkspaceHomeDiffsQueryVariables,
  WorkspaceHomeQuery,
  WorkspaceHomeQueryVariables,
} from './__generated__/WorkspaceHome'
import { useUpdatedWorkspaceDiffs } from '../subscriptions/useUpdatedWorkspaceDiffs'
import type { DeepMaybeRef } from '@vueuse/core'
import WorkspaceName, {
  WORKSPACE_FRAGMENT as WORKSPACE_NAME_FRAGMENT,
} from '../organisms/WorkspaceName.vue'
import WorkspaceDetails, {
  WORKSPACE_FRAGMENT as WORKSPACE_DETAILS_FRAGMENT,
} from '../organisms/WorkspaceDetails.vue'
import WorkspaceDescription, {
  WORKSPACE_FRAGMENT as WORKSPACE_DESCRIPTION_FRAGMENT,
} from '../organisms/WorkspaceDescription.vue'
import { useArchiveWorkspace } from '../mutations/useArchiveWorkspace'
import ConnectNewDirectory, {
  CODEBASE_FRAGMENT as CONNECT_NEW_DIRECTORY_CODEBASE_FRAGMENT,
} from '../organisms/electron/ConnectNewDirectory.vue'

type CodebaseView = WorkspaceHomeQuery['workspace']['codebase']['views'][number]

export default defineComponent({
  components: {
    SelectedHunksToolbar,
    Tooltip,
    FolderAddIcon,
    ButtonWithDropdown,
    ViewStatusIndicator,
    OnboardingStep,
    WorkspaceActivitySidebar,
    Button,
    ResolveConflict,
    LiveDetails,
    AnnotationIcon,
    Banner,
    DesktopComputerIcon,
    MenuItem,
    SearchToolbar,
    OpenInEditor,
    ArchiveButton,
    WorkspaceName,
    WorkspaceDetails,
    WorkspaceDescription,
    ConnectNewDirectory,
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
    const isRemoteEnabled = computed(() => features?.value?.includes(Feature.Remote))

    let route = useRoute()
    const router = useRouter()

    let workspaceID = computed(() => route.params.id as string)
    let shortCodebaseID = computed(() => route.params.codebaseSlug as string)

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

    let { data, fetching, error, executeQuery } = useQuery<
      WorkspaceHomeQuery,
      DeepMaybeRef<WorkspaceHomeQueryVariables>
    >({
      query: gql`
        query WorkspaceHome(
          $workspaceID: ID!
          $isGitHubEnabled: Boolean!
          $isRemoteEnabled: Boolean!
        ) {
          workspace(id: $workspaceID, allowArchived: true) {
            id
            lastLandedAt
            upToDateWithTrunk
            lastActivityAt
            author {
              id
              name
              avatarUrl
            }
            change {
              id
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
            codebase {
              id
              name

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
              ...ConnectNewDirectory_Codebase
            }
            rebaseStatus {
              id
              isRebasing
              conflictingFiles {
                id
                path
                workspaceDiff {
                  ...ResolveConflictDiff
                }
                trunkDiff {
                  ...ResolveConflictDiff
                }
              }
            }
            ...LiveDetailsWorkspace
            ...WorkspaceActivity_Workspace
            ...WorkspaceName_Workspace
            ...WorkspaceDetails_Workspace
            ...WorkspaceDescription_Workspace
          }
        }

        ${WORKSPACE_NAME_FRAGMENT}
        ${ViewFragment}
        ${WORKSPACE_ACTIVITY_WORKSPACE_FRAGMENT}
        ${VIEW_STATUS_INDICATOR}
        ${LIVE_DETAILS_WORKSPACE}
        ${RESOLVE_CONFLICT_DIFF}
        ${WORKSPACE_DETAILS_FRAGMENT}
        ${WORKSPACE_DESCRIPTION_FRAGMENT}
        ${CONNECT_NEW_DIRECTORY_CODEBASE_FRAGMENT}
      `,
      variables: {
        workspaceID: workspaceID,
        isGitHubEnabled: isGitHubEnabled,
        isRemoteEnabled: isRemoteEnabled,
      },
    })

    let {
      data: diffsData,
      fetching: diffsFetching,
      stale: diffsStale,
    } = useQuery<WorkspaceHomeDiffsQuery, DeepMaybeRef<WorkspaceHomeDiffsQueryVariables>>({
      query: gql`
        query WorkspaceHomeDiffs($workspaceID: ID!) {
          workspace(id: $workspaceID) {
            id
            diffs {
                ...LiveDetailsDiffs
            }
          }

          ${LIVE_DETAILS_DIFFS}
        }
      `,
      variables: { workspaceID: workspaceID },
      requestPolicy: 'cache-and-network',
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

    useUpdatedWorkspaceDiffs(
      workspaceID,
      computed(() => !workspaceID.value)
    )

    const openWorkspaceOnViewResult = useOpenWorkspaceOnView()

    let displayView = ref(null)
    let displayViewId = ref(null)

    watch(data, () => {
      if (data.value?.workspace?.view) {
        displayView.value = data.value?.workspace?.view
        displayViewId.value = displayView.value?.id
      } else {
        displayView.value = null
        displayViewId.value = null
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
    const archiveWorkspaceResult = useArchiveWorkspace()
    return {
      fetching,
      data,
      error,
      displayView,
      loadingNewWorkspace,

      diffsData,
      diffsFetching,
      diffsStale,

      openWorkspaceOnViewResult,

      selectedHunkIDs,

      async refresh() {
        await executeQuery({
          requestPolicy: 'network-only',
        })
      },
      async archiveWorkspace(id: string) {
        return archiveWorkspaceResult({ id }).then(({ archiveWorkspace }) => archiveWorkspace)
      },

      async createSuggestion(workspaceID: string, viewID: string) {
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

      lightningBoltIcon: LightningBoltIcon,
    }
  },
  data() {
    const { mutagenIpc, ipc } = window
    return {
      ...this.initialState(),
      ipc,
      mutagenIpc,
    }
  },
  computed: {
    isApp() {
      return !!this.ipc
    },
    mutagenAvailable() {
      return !!this.mutagenIpc
    },
    diffs() {
      let diffws = this.diffsData?.workspace?.id
      let wsID = this.data?.workspace?.id
      if (wsID && wsID !== diffws) {
        return []
      }

      let d = this.diffsData?.workspace?.diffs
      if (d) {
        return d
      }
      return []
    },
    showApproval() {
      return !this.isSuggesting
    },
    showActivity() {
      return !this.isSuggesting
    },
    showSync() {
      return (
        (!this.data.workspace || this.data.workspace.author.id === this.user?.id) &&
        !this.isSuggesting
      )
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
    mostRecentSelfUserView(): CodebaseView | undefined {
      return this.views.length > 0 ? this.views[0] : undefined
    },
    views(): CodebaseView[] {
      return (
        this.data?.workspace.codebase.views?.slice().sort((a: CodebaseView, b: CodebaseView) => {
          return Math.round(b.lastUsedAt / 100) - Math.round(a.lastUsedAt / 100)
        }) ?? []
      )
    },
    connectedViews(): CodebaseView[] {
      return this.views
        .filter(
          (v: CodebaseView) => v.status != null && v.status.state !== ViewStatusState.Disconnected
        )
        .sort((a: CodebaseView, b: CodebaseView) => {
          let aa = a.mountHostname + ' ' + a.shortMountPath
          let bb = b.mountHostname + ' ' + b.shortMountPath
          return aa.localeCompare(bb)
        })
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
    rebaseStatus: function () {
      return this.data.workspace.rebaseStatus
    },
  },
  watch: {
    '$route.params.id': function () {
      if (this.$route.params.id) {
        this.reset()
      }
    },
    'data.workspace.change.id': function (changeId) {
      if (!changeId) return
      this.$router.push({
        name: 'codebaseChange',
        params: {
          id: changeId,
          codebaseSlug: this.$route.params.codebaseSlug,
        },
        query: {
          new: true,
        },
      })
    },
    'data.workspace.codebase.id': function (n) {
      if (n) this.emitter.emit('codebase', n)
    },
  },
  unmounted() {
    this.emitter.off('differ-selected-hunk-ids', this.onSelectedHunkIDs)
    this.emitter.emit('search-toolbar-button-visible', false)
  },
  mounted() {
    this.emitter.on('differ-selected-hunk-ids', this.onSelectedHunkIDs)
    this.emitter.emit('search-toolbar-button-visible', true)
  },
  methods: {
    initialState() {
      return {
        workspaceID: this.$route.params.id,
        codebaseSlug: this.$route.params.codebaseSlug,

        rebasing_complete_no_conflicts: false,
        rebasing_complete_had_conflicts: false,
        rebasing_failed: false,
        rebasing_working: false,
        rebasing_conflict_resolutions: new Map(),

        loadingNewWorkspace: false,
        isSyncing: false,
        archiveWorkspaceActive: false,
      }
    },
    openWorkspaceOnView(workspaceID: string, viewID: string) {
      const variables = { workspaceID, viewID }
      this.loadingNewWorkspace = true
      this.openWorkspaceOnViewResult(variables)
        .catch((e) => {
          console.log(e)
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
      Object.assign(this.$data, this.initialState())
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
        .then(async (data) => {
          await this.refresh()
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
        .then(async (data) => {
          await this.refresh()
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

    async createViewInDirectory() {
      if (this.data?.workspace?.id == null) {
        return
      }

      const oldIsReady = this.mutagenIpc?.isReady && (await this.mutagenIpc.isReady())
      const newIsReady = this.ipc?.state && (await this.ipc.state()) === 'online'

      const mutagenReady = oldIsReady || newIsReady

      if (!mutagenReady) {
        this.emitter.emit('notification', {
          title: 'Sturdy is not running',
          message: 'Sturdy is still starting, please wait.',
          style: 'error',
        })
        return
      }

      await this.mutagenIpc
        .createNewViewWithDialog(this.data.workspace.id, this.data.workspace.codebase.name)
        .catch(async (e: any) => {
          if (e.message.includes('non-empty')) {
            this.emitter.emit('notification', {
              title: 'Directory is not empty',
              message: 'Please select an empty directory.',
              style: 'error',
            })
          } else if (e.message.includes('Cancelled')) {
            await this.archiveWorkspace(this.data.workspace.id)
          } else {
            throw e
          }
        })
    },
  },
})
</script>
