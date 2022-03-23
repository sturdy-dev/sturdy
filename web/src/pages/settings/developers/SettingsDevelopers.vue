<template>
  <PaddedAppLeftSidebar v-if="data" class="bg-white">
    <template #navigation>
      <SettingsVerticalNavigation />
    </template>

    <template #default>
      <div class="max-w-7xl">
        <Header>Developer settings</Header>
        <SettingsDeveloperCodebaseID :codebase-id="data.codebase.id" />
      </div>
    </template>
  </PaddedAppLeftSidebar>
</template>

<script lang="ts" setup>
import { gql, useQuery } from '@urql/vue'
import { useRoute } from 'vue-router'
import PaddedAppLeftSidebar from '../../../layouts/PaddedAppLeftSidebar.vue'
import SettingsVerticalNavigation from '../../../components/codebase/settings/SettingsVerticalNavigation.vue'
import Header from '../../../molecules/Header.vue'
import SettingsDeveloperCodebaseID from '../../../components/codebase/settings/SettingsDeveloperCodebaseID.vue'
import type {
  SettingsDevelopersQuery,
  SettingsDevelopersQueryVariables,
} from './__generated__/SettingsDevelopers'

let route = useRoute()

let { data } = useQuery<SettingsDevelopersQuery, SettingsDevelopersQueryVariables>({
  query: gql`
    query SettingsDevelopers($shortID: ID) {
      codebase(shortID: $shortID) {
        id
      }
    }
  `,
  variables: {
    shortID: route.params.codebaseSlug,
  },
})
</script>
