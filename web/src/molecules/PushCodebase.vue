<template>
  <Button
    :disabled="isPushing"
    :grouped="true"
    :last="true"
    :icon="arrowSmUpIcon"
    :spinner="isPushing"
    color="white"
    @click="triggerPush"
  >
    <span v-if="isPushing">Pushing</span>
    <span v-else>Push</span>
  </Button>
</template>

<script lang="ts">
import { gql } from '@urql/vue'
import { defineComponent, inject, toRef } from 'vue'
import type { PropType } from 'vue'
import type { PushCodebaseRemoteFragment } from './__generated__/PushCodebase'
import { usePushCodebase } from '../mutations/usePushCodebase'
import Button from '../components/shared/Button.vue'
import { ArrowSmUpIcon } from '@heroicons/vue/solid'
import type { Emitter } from 'mitt/src'

export const PUSH_CODEBASE_REMOTE_FRAGMENT = gql`
  fragment PushCodebaseRemote on Remote {
    id
    name
  }
`

export default defineComponent({
  components: { Button },
  props: {
    remote: {
      type: Object as PropType<PushCodebaseRemoteFragment>,
      required: true,
    },
    codebaseId: {
      type: String,
      required: true,
    },
  },
  setup(props) {
    let { mutating: isPushing, pushCodebase } = usePushCodebase()

    const codebaseId = toRef(props, 'codebaseId')
    const remote = toRef(props, 'remote')

    let emitter = inject<Emitter>('emitter')

    const triggerPush = async function () {
      const input = { codebaseID: codebaseId.value }

      await pushCodebase(input)
        .catch((e) => {
          let title = 'Failed!'
          let message = 'Failed to push'

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
              title: 'Pushed!',
              message: 'Pushed changelog to ' + remote.value.name,
              style: 'success',
            })
          }
        })
    }

    return {
      isPushing,
      triggerPush,
      arrowSmUpIcon: ArrowSmUpIcon,
    }
  },
})
</script>
