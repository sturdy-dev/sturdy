<template>
  <PaddedAppRightSidebar v-if="data" class="bg-white">
    <main class="relative md:overflow-y-auto focus:outline-none">
      <div class="block space-y-4">
        <div>
          <div class="flex items-center">
            <!-- For Sturdy the App: Show "Connect Directory" that automatically creates a view + workspace -->
            <ConnectNewDirectory v-if="showAppConnectDirectory" :codebase="data.codebase" />

            <!-- Spacer to make layout render both in app and in browser -->
            <div class="flex-1"></div>

            <div class="relative">
              <span class="relative z-0 inline-flex space-x-4">
                <Button
                  v-if="showCliSetupToggleButton"
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

                <RouterLinkButton
                  :to="{ name: 'codebaseChanges' }"
                  size="wider"
                  color="white"
                  :icon="viewListIcon"
                >
                  Changelog
                </RouterLinkButton>

                <RouterLinkButton
                  v-if="isAuthorized"
                  :to="{ name: 'codebaseSettings' }"
                  :icon="cogIcon"
                  size="wider"
                  color="white"
                >
                  Settings
                </RouterLinkButton>
              </span>
            </div>
          </div>
        </div>

        <!-- setup new view instructions for both CLI and app -->
        <SetupNewView
          v-if="showSetupNewView"
          :codebase="data.codebase"
          :current-user-has-a-view="currentUserHasAView"
          :codebase-slug="codebaseSlug"
        />

        <!-- TODO(gustav): remove or fix -->
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
    </main>

    <template #sidebar>
      <div class="space-y-4">
        <PushPullCodebase
          v-if="data.codebase.remote"
          :remote="data.codebase.remote"
          :codebase-id="data.codebase.id"
        />

        <AssembleTheTeam
          :user="user"
          :members="data.codebase.members"
          :codebase-id="data.codebase.id"
          :changes-count="data.codebase.changes.length"
        />
      </div>
    </template>
  </PaddedAppRightSidebar>
</template>

<script lang="ts">
import { gql, useQuery } from '@urql/vue'
import { useRoute, useRouter } from 'vue-router'
import { computed, defineComponent, inject, onUnmounted, ref, watch } from 'vue'
import type { PropType, Ref } from 'vue'
import SetupNewView from '../components/codebase/SetupNewView.vue'
import Button from '../atoms/Button.vue'
import { useHead } from '@vueuse/head'
import { ChevronDownIcon, ChevronUpIcon, CogIcon, ViewListIcon } from '@heroicons/vue/solid'
import ImportFromGit from '../components/codebase/ImportFromGit.vue'
import { useUpdatedWorkspaceByCodebase } from '../subscriptions/useUpdatedWorkspace'
import Directory, { OPEN_DIRECTORY } from '../components/browse/Directory.vue'
import TopOfChangelogWidget, { TOP_OF_CHANGELOG } from '../organisms/TopOfChangelogWidget.vue'
import NoFilesCodebase from '../components/codebase/NoFilesCodebase.vue'
import ConnectNewDirectory, {
  CODEBASE_FRAGMENT as CONNECT_NEW_DIRECTORY_CODEBASE_FRAGMENT,
} from '../organisms/electron/ConnectNewDirectory.vue'
import WorkspaceList, { WORKSPACE_LIST } from '../components/codebase/WorkspaceList.vue'
import PaddedAppRightSidebar from '../layouts/PaddedAppRightSidebar.vue'
import type {
  CodebaseHomeCodebaseQuery,
  CodebaseHomeCodebaseQueryVariables,
} from './__generated__/CodebaseHome'
import { Feature } from '../__generated__/types'
import type { User } from '../__generated__/types'
import RouterLinkButton from '../atoms/RouterLinkButton.vue'
import AssembleTheTeam from '../organisms/AssembleTheTeam.vue'
import PushPullCodebase, {
  PUSH_PULL_CODEBASE_REMOTE_FRAGMENT,
} from '../molecules/PushPullCodebase.vue'
import type { DeepMaybeRef } from '@vueuse/core'

export default defineComponent({
  components: {
    PushPullCodebase,
    AssembleTheTeam,
    PaddedAppRightSidebar,
    Directory,
    SetupNewView,
    Button,
    ChevronUpIcon,
    ChevronDownIcon,
    ImportFromGit,
    TopOfChangelogWidget,
    NoFilesCodebase,
    ConnectNewDirectory,
    WorkspaceList,
    RouterLinkButton,
  },
  props: {
    user: {
      type: Object as PropType<User>,
    },
  },
  setup() {
    const route = useRoute()
    const codebaseSlug = computed(() =>
      route.params.codebaseSlug ? (route.params.codebaseSlug as string) : ''
    )

    const features = inject<Ref<Array<Feature>>>('features', ref([]))
    const isRemoteEnabled = computed(() => features?.value?.includes(Feature.Remote))

    let { data, fetching, error, executeQuery } = useQuery<
      CodebaseHomeCodebaseQuery,
      DeepMaybeRef<CodebaseHomeCodebaseQueryVariables>
    >({
      query: gql`
        query CodebaseHomeCodebase($shortCodebaseID: ID!, $isRemoteEnabled: Boolean!) {
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
            remote @include(if: $isRemoteEnabled) {
              ...PushPullCodebaseRemote
            }
            ...ConnectNewDirectory_Codebase
          }
        }
        ${OPEN_DIRECTORY}
        ${TOP_OF_CHANGELOG}
        ${WORKSPACE_LIST}
        ${PUSH_PULL_CODEBASE_REMOTE_FRAGMENT}
        ${CONNECT_NEW_DIRECTORY_CODEBASE_FRAGMENT}
      `,
      requestPolicy: 'cache-and-network',
      variables: {
        shortCodebaseID: codebaseSlug,
        isRemoteEnabled,
      },
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

    const isMultiTenancyEnabled = computed(() => features?.value?.includes(Feature.MultiTenancy))

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

      isMultiTenancyEnabled,

      cogIcon: CogIcon,
      viewListIcon: ViewListIcon,
    }
  },
  data() {
    const ipc = window.ipc
    return {
      showSetupInstructions: false,
      ipc,
    }
  },
  computed: {
    thisIsApp() {
      return !!this.ipc
    },

    showCliSetupToggleButton() {
      if (!this.isMultiTenancyEnabled) return false
      if (!this.isAuthorized) return false
      return false
    },

    showSetupNewView() {
      if (!this.isAuthorized) return false
      return (
        !this.currentUserHasAView ||
        !(this.data?.codebase?.rootDir?.children?.length > 0) ||
        this.showSetupInstructions
      )
    },

    showDownloadApp() {
      if (this.thisIsApp) return false
      return true
    },

    showAppConnectDirectory() {
      return this.isAuthorized && this.thisIsApp
    },

    currentUserHasAView() {
      if (this.data) {
        return this.data.codebase?.views.filter((vw) => vw.author.id === this.user?.id).length > 0
      }
      return false
    },
    isAuthenticated() {
      return !!this.user
    },

    isAuthorized() {
      if (this.data) {
        const isMember = this.data.codebase?.members.some(({ id }) => id === this.user?.id)
        return this.isAuthenticated && isMember
      }
      return false
    },
  },
  watch: {
    'data.codebase.id': function (id) {
      if (id) this.emitter.emit('codebase', id)
    },
    error: function (err) {
      if (err) throw err
    },
  },
})
</script>
