<template>
  <PaddedAppRightSidebar class="bg-white">
    <template #toolbar>
      <SearchToolbar />
    </template>

    <template #sidebar>
      <ChangelogDetails
        v-if="!isChangeNotFound"
        :change-data="data"
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
                :codebase-id="data.codebase.id"
                :changes="data.codebase.changes"
                :selected-change-id="selectedChangeID"
                @selectCodebaseChange="onSelectCodebaseChange"
              />
            </div>
          </div>
        </div>
      </div>
    </template>

    <div>
      <div class="md:flex md:items-center md:justify-between md:space-x-4">
        <div v-if="data && data?.change" class="min-h-16 flex-1">
          <h1 class="text-2xl font-bold text-gray-900">
            {{ data.change.title }}
          </h1>
          <p class="mt-2 text-sm text-gray-500">
            By
            {{ ' ' }}
            <span class="font-medium text-gray-900">
              {{ data.change.author?.name }}
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
              {{ data.codebase.name }}
            </router-link>
          </p>
        </div>
        <div v-else-if="!isChangeNotFound" class="h-16 flex-1">
          <!-- Loading -->
          <div class="h-6 w-1/2 bg-gray-300 animate-pulse rounded-md"></div>
        </div>
      </div>

      <aside v-if="data" class="mt-8 xl:hidden">
        <ChangelogDetails
          :change-data="data"
          :github-integration="data?.codebase?.gitHubIntegration"
        />
      </aside>
    </div>

    <section>
      <div class="flex-grow pt-4 z-10 relative min-w-0">
        <div class="pl-1">
          <ChangeDetails
            v-if="data && data?.change"
            :codebase-id="data?.codebase.id"
            :change-id="selectedChangeID"
            :codebase-slug="codebaseSlug"
            :change="data?.change"
            :user="user"
            :members="data?.codebase.members"
          />
        </div>
      </div>
    </section>
  </PaddedAppRightSidebar>
</template>

<script lang="ts">
import ChangeDetails from '../../components/changelog/ChangeDetails.vue'
import { gql, useQuery } from '@urql/vue'
import { useRoute } from 'vue-router'
import ChangelogDetails from '../../components/changelog/ChangelogDetails.vue'
import { STATUS_FRAGMENT } from '../../components/statuses/StatusBadge.vue'
import { MEMBER_FRAGMENT } from '../../components/shared/TextareaMentions.vue'
import SearchToolbar from '../../components/workspace/SearchToolbar.vue'
import PaddedAppRightSidebar from '../../layouts/PaddedAppRightSidebar.vue'
import ChangelogSidebar from '../../components/changelog/ChangelogSidebar.vue'

import { ChangePageQuery, ChangePageQueryVariables } from './__generated__/Change'

const PAGE_QUERY = gql`
  query ChangePage($id: ID, $shortID: ID, $changeID: ID!) {
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
`

export default {
  components: {
    PaddedAppRightSidebar,
    ChangelogDetails,
    ChangeDetails,
    SearchToolbar,
    ChangelogSidebar,
  },
  props: {
    user: {
      type: Object,
      required: false,
    },
  },
  setup() {
    const route = useRoute()
    const selectedChangeID = route.params.selectedChangeID as string
    const codebaseSlug = route.params.codebaseSlug as string

    const { data, fetching, error } = useQuery<ChangePageQuery, ChangePageQueryVariables>({
      query: PAGE_QUERY,
      variables: {
        shortID: codebaseSlug,
        changeID: selectedChangeID,
      },
    })

    return {
      selectedChangeID,
      codebaseSlug,

      data,
      fetching,
      error,
    }
  },
  computed: {
    isChangeNotFound: function () {
      if (!this.error) return
      return (
        this.error.graphQLErrors?.length > 0 &&
        this.error.graphQLErrors[0].message === 'NotFoundError'
      )
    },
  },
  watch: {
    'data.codebase.id': function (id) {
      if (id) this.emitter.emit('codebase', id)
    },
    error: function (error) {
      throw error
    },
  },
  methods: {
    onSelectCodebaseChange(event: { commit_id: string }) {
      this.$router.push({
        name: 'codebaseChange',
        params: {
          codebaseSlug: this.codebaseSlug,
          selectedChangeID: event.commit_id,
        },
      })
    },
  },
}
</script>
