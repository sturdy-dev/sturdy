<template>
  <banner v-if="!gitHubAppValidation.ok" status="error" >
    <label>The GitHub app is not correctly configured. Check the GitHub app permissions & events.
      <ul v-if="gitHubAppValidation.missingPermissions.length > 0" class="list-disc pl-4">
        <li>The app is missing the following permissions:
          <ul class="list-disc pl-4">
            <li v-for="permission in gitHubAppValidation.missingPermissions" :key="permission">{{permission}}</li>
          </ul>
        </li>
      </ul>
      <ul v-if="gitHubAppValidation.missingEvents.length > 0" class="list-disc pl-4">
        <li>The app is missing the following events:
          <ul class="list-disc pl-4">
            <li v-for="event in gitHubAppValidation.missingEvents" :key="event">{{event}}</li>
          </ul>
        </li>
      </ul>
    </label>
  </banner>
</template>

<script lang="ts">
import {defineComponent, PropType} from "vue";
import Banner from "../atoms/Banner.vue";
import {gql} from "@urql/vue";
import {GitHubAppErrorsBanner_GithubValidationAppFragment} from "./__generated__/GitHubAppErrorsBanner";

export const GITHUB_APP_ERRORS_BANNER_GITHUB_VALIDATION_APP_FRAGMENT = gql`
  fragment GitHubAppErrorsBanner_GithubValidationApp on GithubValidationApp {
      ok
      missingPermissions
      missingEvents
  }
`

export default defineComponent({
  components: {Banner},
  props: {
    gitHubAppValidation: {
      type: Object as PropType <GitHubAppErrorsBanner_GithubValidationAppFragment>,
      required: true,
    },
  },
})
</script>

<style scoped>

</style>