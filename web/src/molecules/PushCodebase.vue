<template>
  <Button v-if="isPushing" :disabled="true" :grouped="true" :last="true">
    <Spinner class="mr-2" />
    <span>Pushing</span></Button
  >
  <Button v-else :grouped="true" :last="true" @click="triggerPush">
    <ArrowSmUpIcon class="-ml-1 mr-2 h-5 w-5 text-gray-800" />
    <span>Push</span>
  </Button>
</template>

<script lang="ts">
import { gql } from '@urql/vue'
import { defineComponent, inject, PropType, toRef } from 'vue'
import { PushCodebaseRemoteFragment } from './__generated__/PushCodebase'
import { usePushCodebase } from '../mutations/usepushCodebase'
import Spinner from '../components/shared/Spinner.vue'
import Button from '../components/shared/Button.vue'
import { ArrowSmUpIcon } from '@heroicons/vue/solid'
import { Emitter } from 'mitt/src'

export const PUSH_CODEBASE_REMOTE_FRAGMENT = gql`
  fragment PushCodebaseRemote on Remote {
    id
    name
  }
`

export default defineComponent({
  components: { Spinner, Button, ArrowSmUpIcon },
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
    }
  },
})
</script>
