<template>
  <PaddedAppLeftSidebar v-if="data" class="bg-white">
    <template #navigation>
      <SettingsVerticalNavigation />
    </template>

    <template #default>
      <SettingsTitle :codebase-id="data.codebase.id" :codebase-name="data.codebase.name" />
      <SettingsGitHubIntegration :git-hub-integration="data.codebase.gitHubIntegration" />
      <SettingsDangerzone :codebase="data.codebase" />
    </template>
  </PaddedAppLeftSidebar>
</template>

<script>
import SettingsTitle from '../../components/codebase/settings/SettingsTitle.vue'
import SettingsGitHubIntegration from '../../components/codebase/settings/SettingsGitHubIntegration.vue'
import SettingsDangerzone, {
  SETTINGS_DANGERZONE,
} from '../../components/codebase/settings/SettingsDangerzone.vue'
import { gql, useQuery } from '@urql/vue'
import { useRoute } from 'vue-router'
import PaddedAppLeftSidebar from '../../layouts/PaddedAppLeftSidebar.vue'
import SettingsVerticalNavigation from '../../components/codebase/settings/SettingsVerticalNavigation.vue'

export default {
  name: 'CodebaseSettings',
  components: {
    SettingsTitle,
    SettingsGitHubIntegration,
    SettingsDangerzone,
    PaddedAppLeftSidebar,
    SettingsVerticalNavigation,
  },
  setup() {
    let route = useRoute()

    let { data } = useQuery({
      query: gql`
        query Settings($id: ID, $shortID: ID) {
          codebase(id: $id, shortID: $shortID) {
            id
            name
            gitHubIntegration {
              id
              owner
              name
              enabled
              gitHubIsSourceOfTruth
              trackedBranch
              lastPushAt
              lastPushErrorMessage
            }
            ...SettingsDangerzone
          }
        }

        ${SETTINGS_DANGERZONE}
      `,
      variables: {
        shortID: route.params.codebaseSlug,
      },
    })

    return {
      data,
    }
  },
  watch: {
    'data.codebase.id': function (id) {
      if (id) this.emitter.emit('codebase', id)
    },
  },
  data() {
    return {
      updateStatus: '',
      showRenameFailed: false,
    }
  },
}
</script>
