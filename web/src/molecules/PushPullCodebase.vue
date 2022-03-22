<template>
  <div>
    <h2 class="text-sm font-medium text-gray-500">Connected to {{ remote.name }}</h2>

    <PullCodebase :codebase-id="codebaseId" :remote="remote" />
    <PushCodebase :codebase-id="codebaseId" :remote="remote" />
  </div>
</template>

<script lang="ts">
import { gql } from '@urql/vue'
import { defineComponent, PropType } from 'vue'
import PushCodebase from './PushCodebase.vue'
import PullCodebase from './PullCodebase.vue'
import { PushPullCodebaseRemoteFragment } from './__generated__/PushPullCodebase'

export const PUSH_PULL_CODEBASE_REMOTE_FRAGMENT = gql`
  fragment PushPullCodebaseRemote on Remote {
    id
    name
  }
`

export default defineComponent({
  components: { PushCodebase, PullCodebase },
  props: {
    remote: {
      type: Object as PropType<PushPullCodebaseRemoteFragment>,
      required: true,
    },
    codebaseId: {
      type: String,
      required: true,
    },
  },
})
</script>
