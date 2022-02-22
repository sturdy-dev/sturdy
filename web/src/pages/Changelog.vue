<template>
  <PaddedAppRightSidebar class="bg-white">
    <template #toolbar>
      <SearchToolbar />
    </template>

    <template #default>
      <ChangelogEmpty
        v-if="codebaseEmpty"
        :codebase-slug="codebase_slug"
        class="px-4 sm:px-6 lg:px-8"
      />
      <div v-else>
        <div class="md:flex md:items-center md:justify-between md:space-x-4">
          <div v-if="changeData?.change && data" class="min-h-16 flex-1">
            <h1 class="text-2xl font-bold text-gray-900">
              {{ changeData?.change?.title }}
            </h1>
            <p class="mt-2 text-sm text-gray-500">
              By
              {{ ' ' }}
              <span class="font-medium text-gray-900">
                {{ changeData?.change?.author?.name }}
              </span>
              {{ ' ' }}
              in
              {{ ' ' }}
              <router-link
                :to="{
                  name: 'codebaseHome',
                  params: { codebaseSlug: codebase_slug },
                }"
                class="font-medium text-gray-900"
              >
                {{ data.codebase.name }}
              </router-link>
            </p>
          </div>
          <div v-else class="h-16 flex-1">
            <!-- Loading -->
            <div class="h-6 w-1/2 bg-gray-300 animate-pulse rounded-md"></div>
          </div>
        </div>

        <aside v-if="changeData" class="mt-8 xl:hidden">
          <ChangelogDetails
            :change-data="changeData"
            :github-integration="data?.codebase?.gitHubIntegration"
          />
        </aside>
      </div>

      <section>
        <div class="flex-grow pt-4 z-10 relative min-w-0">
          <div class="pl-1">
            <div v-if="showChangeError" class="text-center text-gray-500 flex flex-col space-y-4">
              <div>This change does not exist.</div>
              <div>
                <Button
                  @click="
                    $router.push({
                      to: 'codebaseChangelog',
                      params: {
                        codebaseSlug: codebase_slug,
                        selectedChangeID: null,
                      },
                    })
                  "
                >
                  Go to latest
                </Button>
              </div>
            </div>
            <ChangeDetails
              v-else-if="selectedChangeID && changeData && data"
              :codebase-id="data.codebase.id"
              :change-id="selectedChangeID"
              :codebase-slug="codebase_slug"
              :change="changeData.change"
              :user="user"
              :members="data.codebase.members"
              @commented="refreshChange"
            />
          </div>
        </div>
      </section>
    </template>

    <template #sidebar>
      <ChangelogDetails
        :change-data="changeData"
        :github-integration="data?.codebase?.gitHubIntegration"
      />

      <div class="mt-10 py-6 space-y-8">
        <div>
          <div class="divide-y divide-gray-200">
            <div class="pb-4">
              <h2 id="activity-title" class="text-lg font-medium text-gray-900">More changes</h2>
            </div>
            <div class="pt-6">
              <ChangelogSidebar
                v-if="data"
                :codebase-i-d="data.codebase.id"
                :changes="data.codebase.changes"
                :selected-change-id="selectedChangeID"
                @selectCodebaseChange="onSelectCodebaseChange"
              />
            </div>
          </div>
        </div>
      </div>
    </template>
  </PaddedAppRightSidebar>
</template>

<script>
import ChangeDetails from '../components/changelog/ChangeDetails.vue'
import { Slug } from '../slug'
import ChangelogSidebar from '../components/changelog/ChangelogSidebar.vue'
import ChangelogEmpty from '../components/changelog/ChangelogEmpty.vue'
import { gql, useQuery } from '@urql/vue'
import { useRoute } from 'vue-router'
import Button from '../components/shared/Button.vue'
import { ref, watch, computed } from 'vue'
import ChangelogDetails from '../components/changelog/ChangelogDetails.vue'
import time from '../time'
import { STATUS_FRAGMENT } from '../components/statuses/StatusBadge.vue'
import { MEMBER_FRAGMENT } from '../components/shared/TextareaMentions.vue'
import SearchToolbar from '../components/workspace/SearchToolbar.vue'
import PaddedAppRightSidebar from '../layouts/PaddedAppRightSidebar.vue'

export default {
  components: {
    PaddedAppRightSidebar,
    ChangelogDetails,
    ChangeDetails,
    ChangelogSidebar,
    Button,
    ChangelogEmpty,
    SearchToolbar,
  },
  props: {
    user: {
      type: Object,
    },
  },
  setup() {
    const route = useRoute()

    let { data, fetching, error } = useQuery({
      query: gql`
        query ChangelogCodebase($id: ID, $shortID: ID) {
          codebase(id: $id, shortID: $shortID) {
            id
            shortID
            name

            gitHubIntegration {
              id
              owner
              name
              enabled
              gitHubIsSourceOfTruth
            }

            members {
              ...Member
            }

            changes {
              id
              title
              createdAt
              author {
                id
                name
                avatarUrl
              }
              comments {
                id
              }
              statuses {
                ...Status
              }
            }
          }
        }
        ${STATUS_FRAGMENT}
        ${MEMBER_FRAGMENT}
      `,
      variables: {
        shortID: route.params.codebaseSlug,
      },
      requestPolicy: 'cache-and-network',
    })

    let selectedChangeID = ref(null)
    let codebaseEmpty = ref(false)

    const calcChangeToShow = () => {
      codebaseEmpty.value = data.value?.codebase?.changes?.length === 0

      if (route.params.selectedChangeID) {
        selectedChangeID.value = route.params.selectedChangeID
      } else if (data.value?.codebase?.changes?.length > 0) {
        selectedChangeID.value = data.value.codebase.changes[0].id
      } else {
        selectedChangeID.value = null
      }
    }

    calcChangeToShow()
    watch(data, () => {
      calcChangeToShow()
    })
    watch(route, () => {
      calcChangeToShow()
    })

    let {
      data: changeData,
      fetching: fetchingChange,
      refresh: refreshChange,
      error: changeError,
    } = useQuery({
      query: gql`
        query Changelog($changeID: ID!) {
          change(id: $changeID) {
            id

            title
            description
            createdAt
            trunkCommitID

            statuses {
              ...Status
            }

            author {
              id
              name
              avatarUrl
            }

            diffs {
              id
              origName
              newName
              preferredName

              isDeleted
              isNew
              isMoved

              isLarge
              largeFileInfo {
                id
                size
              }

              isHidden

              hunks {
                id
                patch
              }
            }

            comments {
              id
              message
              codeContext {
                id
                path
                lineEnd
                lineStart
                lineIsNew
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
          }
        }
        ${STATUS_FRAGMENT}
      `,
      requestPolicy: 'cache-and-network',
      pause: computed(() => !selectedChangeID.value),
      variables: { changeID: selectedChangeID },
    })

    let showChangeError = ref(false)
    watch(changeError, (newErr) => {
      if (
        newErr?.graphQLErrors?.length > 0 &&
        newErr.graphQLErrors[0].message === 'NotFoundError'
      ) {
        showChangeError.value = true
      } else {
        showChangeError.value = false
      }
    })

    return {
      data,
      fetching,
      error,
      selectedChangeID,
      changeData,
      fetchingChange,
      refreshChange,
      showChangeError,
      codebaseEmpty,

      ipc: window.ipc,
    }
  },
  data() {
    return {
      reactive_changes_cancel_func: null,
    }
  },
  computed: {
    codebase_slug() {
      if (this.data) {
        return Slug(this.data.codebase.name, this.data.codebase.shortID)
      }
      return null
    },
  },
  watch: {
    'data.codebase.id': function (id) {
      if (id) this.emitter.emit('codebase', id)
    },
  },
  created() {
    this.emitter.on('select-codebase-changelog-change', this.selectedCodebaseChange)
  },
  unmounted() {
    this.emitter.off('select-codebase-changelog-change', this.selectedCodebaseChange)

    if (this.reactive_changes_cancel_func) {
      this.reactive_changes_cancel_func()
    }
  },
  methods: {
    onSelectCodebaseChange(ev) {
      this.$router.push({
        name: 'codebaseChangelog',
        params: {
          codebaseSlug: this.codebase_slug,
          selectedChangeID: ev.commit_id,
        },
      })
    },
    friendly_ago(ts) {
      return time.getRelativeTime(new Date(ts * 1000))
    },
  },
}
</script>
