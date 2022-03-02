<template>
  <PaddedAppRightSidebar v-if="data" class="bg-white">
    <ChangeListEmpty v-if="data.codebase.changes.length == 0" :codebase-slug="codebaseSlug" />
    <ChangeList
      v-else
      :changes="changes"
      :codebase-slug="codebaseSlug"
      :hasNextPage="hasNextPage"
      @next-page="onNextPage"
    />

    <template #sidebar>
      <AssembleTheTeam
        :user="user"
        :members="data.codebase.members"
        :codebase-id="data.codebase.id"
        :changes-count="data.codebase.changes.length"
      />
    </template>
  </PaddedAppRightSidebar>
</template>

<script lang="ts">
import { gql, useQuery } from '@urql/vue'
import { useRoute } from 'vue-router'
import { PropType, ref, watch, computed } from 'vue'
import { useHead } from '@vueuse/head'
import { DeepMaybeRef } from '@vueuse/core'

import { IdFromSlug } from '../../slug'

import ChangeList, { CHANGELOG_CHANGE_FRAGMENT } from '../../organisms/changelog/ChangeList.vue'
import ChangeListEmpty from '../../organisms/changelog/ChangeList.empty.vue'
import AssembleTheTeam, { CODEBASE_MEMBER_FRAGMENT } from '../../organisms/AssembleTheTeam.vue'
import PaddedAppRightSidebar from '../../layouts/PaddedAppRightSidebar.vue'

import { ChangelogV2Query, ChangelogV2QueryVariables } from './__generated__/Index'
import { User } from '../../__generated__/types'

const PAGE_QUERY = gql`
  query ChangelogV2($codebaseShortId: ID!, $before: ID, $limit: Int!) {
    codebase(shortID: $codebaseShortId) {
      id
      name
      changes(input: { limit: $limit, before: $before }) {
        ...ChangelogChange
      }
      members {
        ...Author
      }
    }
  }
  ${CHANGELOG_CHANGE_FRAGMENT}
  ${CODEBASE_MEMBER_FRAGMENT}
`

export default {
  components: { ChangeList, ChangeListEmpty, PaddedAppRightSidebar, AssembleTheTeam },
  props: {
    user: {
      type: Object as PropType<User>,
    },
  },
  setup() {
    const route = useRoute()

    const codebaseSlug = route.params.codebaseSlug as string
    const codebaseShortId = IdFromSlug(codebaseSlug)
    // fetch one more element than we actually show to know if there are more pages
    const limit = 11

    // store before state in the query
    const before = ref(route.query.before as string)
    watch(route, (newRoute) => {
      before.value = newRoute.query.before as string
    })

    const result = useQuery<ChangelogV2Query, DeepMaybeRef<ChangelogV2QueryVariables>>({
      query: PAGE_QUERY,
      variables: {
        codebaseShortId: codebaseShortId,
        before,
        limit,
      },
    })

    useHead({
      title: computed(() => {
        const name = result.data.value?.codebase?.name
        return name ? `${name} - changes` : 'Sturdy'
      }),
    })

    return {
      result,
      limit,

      data: result.data,
      fetching: result.fetching,
      error: result.error,

      codebaseSlug,
    }
  },
  computed: {
    hasNextPage() {
      return this.data?.codebase?.changes.length === this.limit
    },
    lastChangeId() {
      return this.data?.codebase?.changes.slice(-2)[0]?.id
    },
    changes() {
      const changes = this.data?.codebase?.changes || []
      return changes.length < this.limit ? changes : changes.slice(0, -1)
    },
  },
  methods: {
    onNextPage() {
      this.$router.push({
        query: {
          before: this.lastChangeId,
        },
      })
    },
  },
}
</script>
