<template>
  <div class="p-4 sm:p-8">
    <FirstTimeUserNoNameTakeover v-if="data && data.user && !data.user.name" :user="user" />
    <slot v-else></slot>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import { gql, useQuery } from '@urql/vue'
import FirstTimeUserNoNameTakeover from '../components/user/FirstTimeUserNoNameTakeover.vue'

export default defineComponent({
  components: { FirstTimeUserNoNameTakeover },
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
