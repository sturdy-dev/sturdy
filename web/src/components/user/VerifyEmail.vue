<template>
  <Banner v-if="!emailVerified" status="warning">
    Your email is not verified.
    <a
      href="#"
      class="font-medium underline text-yellow-700 hover:text-yellow-600"
      @click="sendEmailVerification"
      >Resend verification email</a
    >
  </Banner>
  <Banner v-else-if="paramEmailVerified" status="success">Your email is verified.</Banner>
  <Banner v-else-if="paramEmailNotVerified" status="error">
    Email verification failed.
    <a
      href="#"
      class="font-medium underline text-yellow-700 hover:text-yellow-600"
      @click="sendEmailVerification"
      >Try again</a
    >
  </Banner>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import Banner from '../shared/Banner.vue'
import http from '../../http'
import { useRoute } from 'vue-router'

export default defineComponent({
  components: {
    Banner,
  },
  props: {
    emailVerified: {
      type: Boolean,
      required: true,
    },
  },
  emits: ['notification'],
  setup() {
    const route = useRoute()
    const paramEmailVerified = route.query['email_verified']
    const paramEmailNotVerified = route.query['email_not_verified']
    return {
      paramEmailVerified,
      paramEmailNotVerified,
    }
  },
  methods: {
    async sendEmailVerification(e: MouseEvent) {
      e.preventDefault()

      const resp = await fetch(http.url('v3/users/verify-email'), {
        method: 'POST',
        credentials: 'include',
      })

      if (resp.status === 200) {
        this.emitter.emit('notification', {
          title: 'Verification email sent',
          message: 'Please check your email',
        })
      } else {
        throw new Error(`Failed to send an email: status ${resp.status}`)
      }
    },
  },
})
</script>
