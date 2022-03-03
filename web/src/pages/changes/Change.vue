<template>
  <PaddedAppRightSidebar class="bg-white">
    <template #toolbar>
      <SearchToolbar />
    </template>

    <template #sidebar>
      <ChangelogDetailsFetching v-if="fetching" />
      <ChangelogDetails v-else :change="data.change" />

      <ChangeActivitySidebar
        v-if="!fetching"
        class="mt-6"
        :change="data.change"
        :codebase-slug="codebaseSlug"
        :user="user"
      />
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

    <section class="flex-grow pt-4 z-10 relative min-w-0">
      <ChangeDetails v-if="data" :change="data.change" :user="user" />
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
import ChangeActivitySidebar, {
  CHANGE_FRAGMENT as CHANGE_ACTIVITY_CHANGE_FRAGMENT,
} from '../../organisms/ChangeActivitySidebar.vue'

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
        name
      }

      ...ChangeActivity_Change
      ...ChangelogDetails_Change
      ...ChangeDetails_Change
    }
  }

  ${CHANGE_FRAGMENT}
  ${STATUS_FRAGMENT}
  ${MEMBER_FRAGMENT}
  ${CHANGE_DETAILS_CHANGE_FRAGMENT}
  ${CHANGE_ACTIVITY_CHANGE_FRAGMENT}
`

export default {
  components: {
    PaddedAppRightSidebar,
    ChangelogDetails,
    ChangelogDetailsFetching,
    ChangeDetails,
    SearchToolbar,
    ChangeActivitySidebar,
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
