<template>
  <div v-if="!fetching" class="my-2 flex-col space-y-2">
    <Banner v-if="error.length > 0" status="error">{{ error }}</Banner>
    <Banner v-if="showSuccess">Saved!</Banner>
    <Banner v-if="showTriggered">The integration was successfully triggered!</Banner>

    <div class="flex items-center space-x-2">
      <Button v-if="editingIntegrationId" color="green" @click="updateIntegration"
        >Update Integration</Button
      >
      <Button v-else color="green" @click="updateIntegration">Create Integration</Button>
      <Button v-if="editingIntegrationId" @click="triggerIntegration">
        Trigger test build for the latest change:&nbsp;<em>{{ change.title }}</em>
      </Button>
    </div>
  </div>
</template>

<script lang="ts">
import { gql, useQuery } from '@urql/vue'
import { defineComponent, PropType, toRefs } from 'vue'
import { useCreateOrUpdateBuildkiteIntegration } from '../../mutations/useCreateOrUpdateBuildkiteIntegration'
import { useTriggerInstantIntegration } from '../../mutations/useTriggerInstantIntegration'
import {
  CreateOrUpdateBuildkiteIntegrationInput,
  IntegrationProvider,
  TriggerInstantIntegrationInput,
} from '../../__generated__/types'
import {
  BuildkiteHeadChangeFragment,
  BuildkiteIntegrationTestQuery,
  BuildkiteIntegrationTestQueryVariables,
} from './__generated__/TestIntegration'
import Button from '../shared/Button.vue'
import Banner from '../shared/Banner.vue'

export const BUILDKITE_HEAD_CHANGE = gql`
  fragment BuildkiteHeadChange on Change {
    id
    title
  }
`

export default defineComponent({
  components: { Button, Banner },
  props: {
    shortCodebaseId: {
      type: String,
      required: true,
    },
    apiToken: {
      type: String,
      required: true,
    },
    organizationName: {
      type: String,
      required: true,
    },
    pipelineName: {
      type: String,
      required: true,
    },
    webhookSecret: {
      type: String,
      required: true,
    },
    change: {
      type: Object as PropType<BuildkiteHeadChangeFragment>,
      required: true,
    },
    editingIntegrationId: {
      type: String,
      required: false,
    },
  },
  setup(props) {
    const { shortCodebaseId } = toRefs(props)
    const { data, fetching } = useQuery<
      BuildkiteIntegrationTestQuery,
      BuildkiteIntegrationTestQueryVariables
    >({
      query: gql`
        query BuildkiteIntegrationTest($shortID: ID!) {
          codebase(shortID: $shortID) {
            id
          }
        }
      `,
      variables: {
        shortID: shortCodebaseId.value,
      },
    })

    const triggerInstantIntegrationResult = useTriggerInstantIntegration()
    const createOrUpdateBuildkiteIntegrationResult = useCreateOrUpdateBuildkiteIntegration()
    return {
      data,
      fetching,

      async createOrUpdateBuildkiteIntegration(input: CreateOrUpdateBuildkiteIntegrationInput) {
        return await createOrUpdateBuildkiteIntegrationResult(input)
      },

      async triggerInstantIntegration(input: TriggerInstantIntegrationInput) {
        return await triggerInstantIntegrationResult(input)
      },
    }
  },
  data() {
    return {
      error: '',
      showSuccess: false,
      showTriggered: false,
    }
  },
  watch: {
    fetching: function (isFetching: boolean) {
      if (!isFetching) this.updateIntegration()
    },
  },
  methods: {
    updateIntegration() {
      this.showSuccess = false
      this.createOrUpdateBuildkiteIntegration({
        integrationID: this.editingIntegrationId,
        codebaseID: this.data!.codebase.id,
        apiToken: this.apiToken,
        organizationName: this.organizationName,
        pipelineName: this.pipelineName,
        webhookSecret: this.webhookSecret,
      })
        .then((res) => {
          this.error = ''
          this.showSuccess = true
          console.log(res)
          this.$router.replace({
            name: 'codebaseSettingsEditBuildkite',
            params: { integrationId: res.createOrUpdateBuildkiteIntegration.id },
          })
        })
        .catch((err: Error) => {
          this.error = err.message
        })
    },

    triggerIntegration() {
      this.triggerInstantIntegration({
        changeID: this.change.id,
        providers: [IntegrationProvider.Buildkite],
      })
        .then((statuses) => {
          this.error = ''
          console.log(statuses)
        })
        .catch((err: Error) => {
          this.error = err.message
        })
    },
  },
})
</script>
