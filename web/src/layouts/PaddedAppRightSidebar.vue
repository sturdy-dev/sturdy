<template>
  <div>
    <div
      class="fixed z-50 md:pr-64 w-full"
      :style="ipc ? 'top: calc(env(titlebar-area-height, 2rem) + 1px)' : 'top: 0'"
    >
      <slot name="toolbar" />
    </div>

    <div class="p-4 sm:p-8 grid grid-cols-1 xl:grid-cols-4">
      <div class="xl:col-span-3 xl:pr-8 xl:border-r xl:border-gray-200">
        <slot />
      </div>
      <div class="xl:pl-8 hidden xl:block">
        <slot name="sidebar" />
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import { gql, useQuery } from '@urql/vue'

export default defineComponent({
  components: {},
  setup() {
    const result = useQuery({
      query: gql`
        query PaddedApp {
          user {
            id
            name
          }
        }
      `,
    })

    return {
      data: result.data,
      ipc: window.ipc,
    }
  },
})
</script>
