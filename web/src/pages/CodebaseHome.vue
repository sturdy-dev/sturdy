<template>
  <PaddedApp v-if="data" class="bg-white">
    <main class="relative md:overflow-y-auto focus:outline-none">
      <ArchiveWorkspaceModal
        :is-active="archiveWorkspaceActive"
        :workspace-i-d="workspaceIdToArchive"
        @deletedWorkspace="fetchWorkspaces"
        @closeDeleteWorkspace="hideArchiveModal"
      />

      <div class="">
        <div class="grid grid-cols-1 xl:grid-cols-4">
          <div
            class="block col-start-1 col-span-4 xl:col-span-3 space-y-4 xl:pr-8 xl:border-r xl:border-gray-200"
          >
            <div>
              <div class="flex items-center">
                <!-- For Sturdy the App: Show "Connect Directory" that automatically creates a view + workspace -->
                <CreateViewAndWorkspace
                  v-if="!showSetupNewView"
                  :codebase-id="data.codebase.id"
                  :codebase-slug="codebaseSlug"
                />

                <!-- Spacer to make layout render both in app and in browser -->
                <div class="flex-1"></div>

                <div class="relative">
                  <span class="relative z-0 inline-flex space-x-4">
                    <Button
                      v-if="currentUserHasAView && !ipc"
                      size="wider"
                      color="lightblue"
                      class="hidden lg:inline-flex"
                      @click="showSetupInstructions = !showSetupInstructions"
                    >
                      <ChevronUpIcon
                        v-if="showSetupInstructions"
                        class="-ml-1 mr-2 h-5 w-5 text-gray-800"
                        aria-hidden="true"
                      />
                      <ChevronDownIcon
                        v-else
                        class="-ml-1 mr-2 h-5 w-5 text-gray-800"
                        aria-hidden="true"
                      />
                      <span>Setup</span>
                    </Button>

                    <Button
                      size="wider"
                      @click="
                        $router.push({
                          name: 'codebaseChangelog',
                          params: { codebaseSlug: codebaseSlug },
                        })
                      "
                    >
                      <ViewListIcon class="-ml-1 mr-2 h-5 w-5 text-gray-400" aria-hidden="true" />
                      <span>Changelog</span>
                    </Button>

                    <Button
                      v-if="isAuthorized"
                      size="wider"
                      @click="
                        $router.push({
                          name: 'codebaseSettings',
                          params: { codebaseSlug: codebaseSlug },
                        })
                      "
                    >
                      <CogIcon class="-ml-1 mr-2 h-5 w-5 text-gray-400" aria-hidden="true" />
                      <span>Settings</span>
                    </Button>
                  </span>
                </div>
              </div>
            </div>

            <SetupNewView
              v-if="showSetupNewView"
              :codebase="data.codebase"
              :current-user-has-a-view="currentUserHasAView"
              :codebase-slug="codebaseSlug"
            />

            <ImportFromGit
              v-if="data.codebase.changes?.length === 0 && data.codebase.workspaces.length === 0"
              :codebase-id="data.codebase.id"
            />

            <TopOfChangelogWidget v-if="data.codebase" :codebase="data.codebase" />

            <Directory
              v-if="data.codebase?.rootDir?.children.length > 0"
              :directory="data.codebase.rootDir"
              :codebase="data.codebase"
            />
            <NoFilesCodebase v-else class="pt-24" />

            <!-- Workspace list for mobile -->
            <WorkspaceList :workspaces="data.codebase.workspaces" class="block md:hidden" />
          </div>

          <aside class="col-start-4 col-span-1 px-4 pl-8 hidden xl:block">
            <OnboardingStep
              id="InvitingTeamMembers"
              :dependencies="['FindingYourCodebase']"
              :enabled="data.codebase.changes.length > 0"
            >
              <template #title> Assemble the Team</template>
              <template #description>
                In this section over here you can invite other people to collaborate on this
                codebase with you.
              </template>

              <CodebaseMembers
                :user="user"
                :members="data.codebase.members"
                :codebase-id="data.codebase.id"
              />
            </OnboardingStep>
          </aside>
        </div>
      </div>
    </main>
  </PaddedApp>
</template>

<script>
import ArchiveWorkspaceModal from '../components/codebase/ArchiveWorkspaceModal.vue'
import time from '../time'
import { gql, useQuery } from '@urql/vue'
import { useRoute, useRouter } from 'vue-router'
import { computed, onUnmounted, ref, watch } from 'vue'
import SetupNewView from '../components/codebase/SetupNewView.vue'
import Button from '../components/shared/Button.vue'
import { useHead } from '@vueuse/head'
import CodebaseMembers from '../components/codebase/CodebaseMembers.vue'
import { ChevronDownIcon, ChevronUpIcon, CogIcon, ViewListIcon } from '@heroicons/vue/solid'
import ImportFromGit from '../components/codebase/ImportFromGit.vue'
import { useUpdatedWorkspaceByCodebase } from '../subscriptions/useUpdatedWorkspace'
import Directory, { OPEN_DIRECTORY } from '../components/browse/Directory.vue'
import TopOfChangelogWidget, {
  TOP_OF_CHANGELOG,
} from '../components/changelog/TopOfChangelogWidget.vue'
import OnboardingStep from '../components/onboarding/OnboardingStep.vue'
import NoFilesCodebase from '../components/codebase/NoFilesCodebase.vue'
import CreateViewAndWorkspace from '../components/codebase/CreateViewAndWorkspace.vue'
import WorkspaceList, { WORKSPACE_LIST } from '../components/codebase/WorkspaceList.vue'
import PaddedApp from '../layouts/PaddedApp.vue'

export default {
  name: 'CodebaseHome',
  components: {
    PaddedApp,
    OnboardingStep,
    Directory,
    SetupNewView,
    ArchiveWorkspaceModal,
    CodebaseMembers,
    Button,
    ViewListIcon,
    CogIcon,
    ChevronUpIcon,
    ChevronDownIcon,
    ImportFromGit,
    TopOfChangelogWidget,
    NoFilesCodebase,
    CreateViewAndWorkspace,
    WorkspaceList,
  },
  props: ['user'],
  setup() {
    let route = useRoute()
    let codebaseSlug = ref(route.params.codebaseSlug)
    watch(
      () => route.params.codebaseSlug,
      (slug) => {
        codebaseSlug.value = slug
      }
    )

    let { data, fetching, error, executeQuery } = useQuery({
      query: gql`
        query CodebaseHomeCodebase($shortCodebaseID: ID!) {
          codebase(shortID: $shortCodebaseID) {
            id
            shortID
            name
            description
            inviteCode
            createdAt
            archivedAt
            members {
              id
              name
              avatarUrl
            }
            workspaces {
              id
              ...WorkspaceList
            }
            rootDir: file(path: "/") {
              ... on Directory {
                ...OpenDirectory
              }
            }
            changes(input: { limit: 4 }) {
              id
            }
            views {
              id
              author {
                id
              }
            }
            ...TopOfChangelog
          }
        }
        ${OPEN_DIRECTORY}
        ${TOP_OF_CHANGELOG}
        ${WORKSPACE_LIST}
      `,
      requestPolicy: 'cache-and-network',
      variables: { shortCodebaseID: codebaseSlug },
    })

    let codebaseID = ref('')
    watch(data, (n) => {
      if (n && n.codebase === null) {
        throw new Error('SturdyCodebaseNotFoundError')
      }

      if (n && n.codebase) {
        codebaseID.value = n.codebase.id
      }
    })

    const router = useRouter()

    useHead({
      title: computed(() => {
        let n = data.value?.codebase?.name
        if (n) {
          return n + ' on Sturdy'
        }
        return 'Sturdy'
      }),
      base: {
        href: router.resolve({
          name: 'browseFile',
          params: {
            codebaseID: codebaseID,
            path: ['/'],
          },
        }).href,
      },
    })

    useUpdatedWorkspaceByCodebase(codebaseSlug)

    let now = ref(new Date())
    let nowInterval = setInterval(() => {
      now.value = new Date()
    }, 3000)
    onUnmounted(() => {
      clearInterval(nowInterval)
    })

    return {
      fetching: fetching,
      data: data,
      error,
      refresh() {
        executeQuery({
          requestPolicy: 'network-only',
        })
      },
      now,
      codebaseSlug,
      ipc: window.ipc,
    }
  },
  data() {
    return this.initialState()
  },
  computed: {
    thisIsApp() {
      return !!this.ipc
    },
    showSetupNewView() {
      if (!this.isAuthorized) return false
      return (
        !this.currentUserHasAView ||
        !(this.data.codebase?.rootDir?.children.length > 0) ||
        this.showSetupInstructions
      )
    },
    currentUserHasAView() {
      if (this.data) {
        return this.data.codebase.views.filter((vw) => vw.author.id === this.user?.id).length > 0
      }
      return false
    },
    isAuthenticated() {
      return !!this.user
    },
    isAuthorized() {
      if (this.data) {
        const isMember = this.data.codebase.members.some(({ id }) => id === this.user?.id)
        return this.isAuthenticated && isMember
      }
      return false
    },
  },
  watch: {
    error: function (err) {
      if (err) throw err
    },
    '$route.params.codebaseSlug': function () {
      if (this.$route.params.codebaseSlug) {
        this.refresh()
      }
    },
  },
  unmounted() {
    clearTimeout(this.fetchWorkspacesTimeout)
  },
  methods: {
    initialState() {
      return {
        archiveWorkspaceActive: false,
        workspaceIdToArchive: '',
        showSetupInstructions: false,
      }
    },
    reset() {
      clearTimeout(this.fetchWorkspacesTimeout)
      Object.assign(this.$data, this.initialState())
    },
    showArchiveModal(ws) {
      this.archiveWorkspaceActive = true
      this.workspaceIdToArchive = ws.id
    },
    hideArchiveModal() {
      this.archiveWorkspaceActive = false
      this.workspaceIdToArchive = ''
    },
    fetchWorkspaces() {
      this.refresh()
    },

    friendly_ago(ts) {
      return time.getRelativeTime(new Date(ts * 1000), this.now)
    },
  },
}
</script>
