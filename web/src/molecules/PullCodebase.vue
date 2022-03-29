<template>
  <Button
    :disabled="isPulling"
    :grouped="true"
    :first="true"
    :icon="arrowSmDownIcon"
    :spinner="isPulling"
    color="white"
    @click="triggerPull"
  >
    <span v-if="isPulling">Pulling</span>
    <span v-else>Pull</span>
  </Button>
</template>

<script lang="ts">
import { gql } from '@urql/vue'
import { defineComponent, inject, toRef } from 'vue'
import type { PropType } from 'vue'
import type { PullCodebaseRemoteFragment } from './__generated__/PullCodebase'
import { usePullCodebase } from '../mutations/usePullCodebase'
import Button from '../atoms/Button.vue'
import { ArrowSmDownIcon } from '@heroicons/vue/solid'
import type { Emitter } from 'mitt/src'

export const PULL_CODEBASE_REMOTE_FRAGMENT = gql`
  fragment PullCodebaseRemote on Remote {
    id
    name
  }
`

export default defineComponent({
  components: { Button },
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
  setup: function (props) {
    let { mutating: isPulling, pullCodebase } = usePullCodebase()

    const codebaseId = toRef(props, 'codebaseId')
    const remote = toRef(props, 'remote')

    let emitter = inject<Emitter>('emitter')

    const triggerPull = async function () {
      const input = { codebaseID: codebaseId.value }

      await pullCodebase(input)
        .catch((e) => {
          let title = 'Failed!'
          let message = 'Failed to pull'

          console.error(e)

          if (emitter) {
            emitter.emit('notification', {
              title: title,
              message,
              style: 'error',
            })
          }
        })
        .then(() => {
          if (emitter) {
            emitter.emit('notification', {
              title: 'Pulled!',
              message: 'Pulled from ' + remote.value.name + ', the changelog has been updated!',
              style: 'success',
            })
          }
        })
    }

    return {
      isPulling,
      triggerPull,
      arrowSmDownIcon: ArrowSmDownIcon,
    }
  },
})
</script>
