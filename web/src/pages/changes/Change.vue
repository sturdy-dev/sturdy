<template>
  <PaddedAppRightSidebar class="bg-white">
    <template #toolbar>
      <SearchToolbar />
    </template>

    <template #sidebar>
      <ChangelogDetailsFetching v-if="fetching" />
      <ChangelogDetails v-else :change="data.change" />

      <div v-if="data?.change" class="mt-10 py-6 space-y-8">
        <div class="divide-y divide-gray-200">
          <div class="pb-4">
            <h2 id="activity-title" class="text-lg font-medium text-gray-900">More changes</h2>
          </div>
          <div class="pt-6">
            <ChangelogSidebar
              :codebase-id="data.change.codebase.id"
              :changes="data.change.codebase.changes"
              :selected-change-id="selectedChangeID"
              @selectCodebaseChange="onSelectCodebaseChange"
            />
          </div>
        </div>
      </div>
    </template>

    <div>
      <div class="md:flex md:items-center md:justify-between md:space-x-4">
        <div class="min-h-16 flex-1">
          <h1 class="text-2xl font-bold text-gray-900">
            <template v-if="fetching">
              <div class="h-6 w-1/2 bg-gray-300 animate-pulse rounded-md" />
            </template>
            <template v-else>
              {{ data.change.title }}
            </template>
          </h1>

          <p class="mt-2 text-sm text-gray-500">
            <template v-if="fetching">
              <div class="h-3 w-1/4 bg-gray-300 animate-pulse rounded-md" />
            </template>

            <template v-else>
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
                {{ data.change.codebase.name }}
              </router-link>
            </template>
          </p>
        </div>
      </div>

      <aside v-if="data" class="mt-8 xl:hidden">
        <ChangelogDetails :change="data.change" />
      </aside>
    </div>

    <section>
      <div class="flex-grow pt-4 z-10 relative min-w-0">
        <ChangeDetails
          v-if="data"
          :codebase-id="data.change.codebase.id"
          :change-id="selectedChangeID"
          :codebase-slug="codebaseSlug"
          :change="data.change"
          :user="user"
          :members="data.change.codebase.members"
        />
      </div>
    </section>
  </PaddedAppRightSidebar>
</template>

<script lang="ts">
import { ref, watch } from 'vue'
import { DeepMaybeRef } from '@vueuse/core'
import ChangeDetails, {
  CHANGE_DETAILS_CHANGE_FRAGMENT,
} from '../../components/changelog/ChangeDetails.vue'
import { gql, useQuery } from '@urql/vue'
import { useRoute } from 'vue-router'
import ChangelogDetails, { CHANGE_FRAGMENT } from '../../components/changelog/ChangelogDetails.vue'
import ChangelogDetailsFetching from '../../components/changelog/ChangelogDetails.fetching.vue'
import { STATUS_FRAGMENT } from '../../components/statuses/StatusBadge.vue'
import { MEMBER_FRAGMENT } from '../../components/shared/TextareaMentions.vue'
import SearchToolbar from '../../components/workspace/SearchToolbar.vue'
import PaddedAppRightSidebar from '../../layouts/PaddedAppRightSidebar.vue'
import ChangelogSidebar from '../../components/changelog/ChangelogSidebar.vue'

import { ChangePageQuery, ChangePageQueryVariables } from './__generated__/Change'

const PAGE_QUERY = gql`
  query ChangePage($changeID: ID!) {
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

      codebase {
        id
        shortID
        name

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
      ...ChangelogDetails_Change
      ...ChangeDetails_Change
    }
  }

  ${CHANGE_FRAGMENT}
  ${STATUS_FRAGMENT}
  ${MEMBER_FRAGMENT}
  ${CHANGE_DETAILS_CHANGE_FRAGMENT}
`

export default {
  components: {
    PaddedAppRightSidebar,
    ChangelogDetails,
    ChangelogDetailsFetching,
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
    const codebaseSlug = route.params.codebaseSlug as string
    const selectedChangeID = ref(route.params.selectedChangeID as string)
    watch(route, (newRoute) => {
      selectedChangeID.value = newRoute.params.selectedChangeID as string
    })

    const { data, fetching, error } = useQuery<
      ChangePageQuery,
      DeepMaybeRef<ChangePageQueryVariables>
    >({
      query: PAGE_QUERY,
      variables: {
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
        params: {
          selectedChangeID: event.commit_id,
        },
      })
    },
  },
}
</script>
