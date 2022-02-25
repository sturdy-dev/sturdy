<template>
  <PaddedAppRightSidebar v-if="!fetching">
    <ChangeList :changes="data.codebase.changes" />

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
import { PropType } from 'vue'

import { IdFromSlug } from '../../slug'

import ChangeList, { CHANGELOG_CHANGE_FRAGMENT } from '../../organisms/changelog/ChangeList.vue'
import AssembleTheTeam, { CODEBASE_MEMBER_FRAGMENT } from '../../organisms/AssembleTheTeam.vue'
import PaddedAppRightSidebar from '../../layouts/PaddedAppRightSidebar.vue'

import { ChangelogV2Query, ChangelogV2QueryVariables } from './__generated__/Index'
import { User } from '../../__generated__/types'

export default {
  components: { ChangeList, PaddedAppRightSidebar, AssembleTheTeam },
  props: {
    user: {
      type: Object as PropType<User>,
    },
  },
  setup() {
    const route = useRoute()
    const codebaseSlug = route.params.codebaseSlug as string
    const codebaseShortId = IdFromSlug(codebaseSlug)
    const { data, fetching, error } = useQuery<ChangelogV2Query, ChangelogV2QueryVariables>({
      query: gql`
        query ChangelogV2($codebaseShortId: ID!) {
          codebase(shortID: $codebaseShortId) {
            id
            changes {
              ...Changelog_Change
            }
            members {
              ...CodebaseMember
            }
          }
        }
        ${CHANGELOG_CHANGE_FRAGMENT}
        ${CODEBASE_MEMBER_FRAGMENT}
      `,
      variables: {
        codebaseShortId: codebaseShortId,
      },
    })
    return {
      data,
      fetching,
      error,
    }
  },
}
</script>
