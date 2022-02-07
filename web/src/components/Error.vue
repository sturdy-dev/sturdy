<template>
  <div class="px-4 py-16 sm:px-6 sm:py-24 lg:px-8 flex justify-around">
    <div class="flex gap-12 flex-col">
      <main class="sm:flex">
        <p class="text-4xl font-extrabold text-yellow-400 sm:text-5xl" :class="[accentColor]">
          {{ code }}
        </p>
        <div class="sm:ml-6">
          <div class="sm:border-l sm:border-gray-200 sm:pl-6">
            <h1 class="text-4xl font-extrabold text-gray-900 tracking-tight sm:text-5xl">
              {{ title }}
            </h1>
            <p class="mt-1 text-base text-gray-500">
              {{ subtitle }}
            </p>
          </div>
        </div>
      </main>

      <div>
        <Button color="green" @click="goBack">Take me back to safety</Button>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import Button from './shared/Button.vue'
import { gql, useQuery } from '@urql/vue'
import { defineComponent } from 'vue'

export default defineComponent({
  name: 'ErrorPage',
  components: { Button },
  props: {
    error: { type: Error, required: true },
  },
  emits: ['reset-error'],
  setup() {
    let { data } = useQuery({
      query: gql`
        query ErrorPage {
          user {
            id
          }
        }
      `,
    })

    return {
      data,
      isApp: !!window.ipc,
    }
  },
  computed: {
    isNotFound() {
      if (this.error?.message === 'SturdyCodebaseNotFoundError') return true
      if (!this.error.graphQLErrors) return false
      const notFoundError = (e) => e.message === 'NotFoundError'
      return this.error.graphQLErrors.filter(notFoundError).length > 0
    },
    code() {
      if (this.isNotFound) {
        return '404'
      }
      return '500'
    },
    accentColor() {
      if (this.isNotFound) {
        return 'text-yellow-400'
      }
      return 'text-red-500'
    },
    title() {
      if (this.isNotFound) {
        return 'Page not found'
      }
      return 'Oops! Something went wrong'
    },
    subtitle() {
      if (this.isNotFound && this.isApp) {
        return 'This page could not be found.'
      }
      if (this.isNotFound) {
        return 'Please check the URL in the address bar and try again.'
      }
      return 'Please try again later.'
    },
  },
  mounted() {
    console.log(this.error)
  },
  methods: {
    async goBack() {
      if (this.data?.user?.id) {
        await this.$router.push({ name: 'home' })
      } else {
        await this.$router.push({ name: 'login' })
      }
      this.$emit('reset-error')
    },
  },
})
</script>
