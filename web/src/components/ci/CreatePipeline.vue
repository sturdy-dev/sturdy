<template>
  <p class="my-2 max-w-2xl text-sm text-gray-500">
    Create a new pipeline on Buildkite for Sturdy to trigger. Give it a cool name, and enter it
    below.
  </p>
  <LinkButton :href="buildkitePipelinePage" target="_blank" class="my-2">
    <span>Create new pipeline</span>
    <ExternalLinkIcon class="h-5 w-5 ml-1" />
  </LinkButton>
  <Instructions
    description="When asked, enter the following settings:"
    :instructions="instructions"
  />
  <TextInput v-model="value" placeholder="The name of the Buildkite Pipeline" />
</template>

<script lang="ts">
import LinkButton from '../shared/LinkButton.vue'
import { ExternalLinkIcon } from '@heroicons/vue/solid'
import Instructions, { Instruction } from './Instructions.vue'
import TextInput from './TextInput.vue'
import { useCreateServiceToken } from '../../mutations/useCreateServiceToken'

export default {
  components: { LinkButton, ExternalLinkIcon, Instructions, TextInput },
  props: {
    modelValue: {
      type: String,
      required: true,
    },
    buildkiteOrganizationSlug: {
      type: String,
      required: true,
    },
    shortCodebaseId: {
      type: String,
      required: true,
    },
  },
  emits: ['update:modelValue'],
  data() {
    return {
      serviceTokenID: '',
      serviceTokenPassword: '',
      serviceTokenLoading: true,
    }
  },
  computed: {
    value: {
      get() {
        return this.modelValue
      },
      set(value: string) {
        this.$emit('update:modelValue', value)
      },
    },
    instructions(): Instruction[] {
      return [
        {
          name: 'Git Repository URL',
          value: this.gitRepositoryURL,
          pre: true,
        },
        { name: 'Command to run', value: './download', pre: true },
      ]
    },
    buildkitePipelinePage(): string {
      return `https://buildkite.com/organizations/${this.buildkiteOrganizationSlug}/pipelines/new`
    },
    gitRepositoryURL(): string {
      return `https://${this.serviceTokenID}:${this.serviceTokenPassword}@git.getsturdy.com`
    },
  },
  mounted() {
    this.fetchServiceToken()
  },
  methods: {
    fetchServiceToken() {
      useCreateServiceToken()({
        shortCodebaseID: this.shortCodebaseId,
        name: 'Buildkite CI token',
      }).then((result) => {
        this.serviceTokenID = result.id
        this.serviceTokenPassword = result.token!
        this.serviceTokenLoading = false
      })
    },
  },
}
</script>
