<template>
  <div class="p-4 sm:p-8 flex flex-col md:flex-row md:space-x-4 space-y-8 md:space-y-0">
    <div class="md:w-60">
      <slot name="navigation" />
    </div>
    <div class="md:flex-1 md:max-w-7xl">
      <slot />
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
    }
  },
})
</script>
