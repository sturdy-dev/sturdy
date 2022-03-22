<template>
  <Button v-if="isPulling" :disabled="true">
    <Spinner class="mr-2" />
    <span>Pulling</span></Button
  >
  <Button v-else @click="triggerPull">Pull from {{ remote.name }}</Button>
</template>

<script lang="ts">
import { gql } from '@urql/vue'
import { defineComponent, PropType, toRef } from 'vue'
import { PullCodebaseRemoteFragment } from './__generated__/PullCodebase'
import { usePullCodebase } from '../mutations/usePullCodebase'
import Spinner from '../components/shared/Spinner.vue'
import Button from '../components/shared/Button.vue'

export const PULL_CODEBASE_REMOTE_FRAGMENT = gql`
  fragment PullCodebaseRemote on Remote {
    id
    name
  }
`

export default defineComponent({
  components: { Spinner, Button },
  props: {
    remote: {
      type: Object as PropType<PullCodebaseRemoteFragment>,
      required: true,
    },
    codebaseId: {
      type: String,
      required: true,
    },
  },
  setup(props) {
    let { mutating: isPulling, pullCodebase } = usePullCodebase()

    const codebaseId = toRef(props, 'codebaseId')

    const triggerPull = async function () {
      const input = { codebaseID: codebaseId.value }

      await pullCodebase(input).catch((e) => {
        let title = 'Failed!'
        let message = 'Failed to pull'

        console.error(e)

        this.emitter.emit('notification', {
          title: title,
          message,
          style: 'error',
        })
      })
    }

    return {
      isPulling,
      triggerPull,
    }
  },
})
</script>
