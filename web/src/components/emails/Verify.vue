<template>
  <div class="flex items-start justify-center min-h-screen pt-4 px-4 pb-20 text-center sm:p-0">
    <div
      class="inline-block bg-white rounded-lg px-4 pt-5 pb-4 text-left overflow-hidden shadow-xl transform transition-all sm:my-8 sm:align-middle sm:max-w-sm sm:w-full sm:p-6"
    >
      <div>
        <div class="mx-auto flex items-center justify-center h-12 w-12 rounded-full bg-yellow-100">
          <Spinner class="h-6 w-6 text-yellow-600" aria-hidden="true" />
        </div>
        <div class="mt-3 text-center sm:mt-5">
          <h3 class="text-lg leading-6 font-medium text-gray-900">Please wait...</h3>
          <div class="mt-2">
            <p class="text-sm text-gray-500">We are verifying your email</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import { useVerifyEmail } from '../../mutations/useVerifyEmail'
import Spinner from '../../atoms/Spinner.vue'
import { useRoute, useRouter } from 'vue-router'

export default defineComponent({
  components: { Spinner },
  setup() {
    const verifyEmailResult = useVerifyEmail()
    const router = useRouter()
    const route = useRoute()
    const token = route.query['token']
    verifyEmailResult({ token })
      .then(() => router.push({ name: 'user', query: { email_verified: 'true' } }))
      .catch(() => router.push({ name: 'user', query: { email_not_verified: 'true' } }))

    return {}
  },
})
</script>
