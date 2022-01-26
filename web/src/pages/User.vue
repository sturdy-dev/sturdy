<template>
  <PaddedApp v-if="data">
    <div v-if="data" class="max-w-screen-xl mx-auto w-full">
      <form
        class="divide-y divide-gray-200 lg:col-span-9"
        enctype="multipart/form-data"
        @submit.stop.prevent="save"
      >
        <div class="space-y-8">
          <div>
            <h2 class="text-lg leading-6 font-medium text-gray-900">Profile</h2>

            <div class="mt-6 flex flex-col lg:flex-row">
              <div class="flex-grow space-y-6">
                <div class="col-span-12 sm:col-span-6">
                  <label for="name" class="block text-sm font-medium text-gray-700">Name</label>
                  <input
                    id="name"
                    v-model="userName"
                    type="text"
                    name="name"
                    autocomplete="name"
                    class="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-light-blue-500 focus:border-light-blue-500 sm:text-sm"
                  />
                </div>

                <div v-if="passwordEnabled" class="col-span-12 sm:col-span-6">
                  <label for="password" class="block text-sm font-medium text-gray-700"
                    >Password</label
                  >
                  <input
                    id="password"
                    v-model="userPassword"
                    type="password"
                    name="name"
                    class="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-light-blue-500 focus:border-light-blue-500 sm:text-sm"
                  />
                </div>

                <div class="col-span-12 sm:col-span-6">
                  <label for="email" class="block text-sm font-medium text-gray-700">Email</label>
                  <input
                    id="email"
                    v-model="userEmail"
                    type="text"
                    name="email"
                    autocomplete="email"
                    class="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-light-blue-500 focus:border-light-blue-500 sm:text-sm"
                  />
                </div>
              </div>

              <div class="mt-6 flex-grow lg:mt-0 lg:ml-6 lg:flex-grow-0 lg:flex-shrink-0">
                <p class="text-sm font-medium text-gray-700" aria-hidden="true">Photo</p>

                <!-- This is the mobile version! -->
                <div class="mt-1 lg:hidden">
                  <div class="flex items-center">
                    <div
                      class="flex-shrink-0 inline-block rounded-full overflow-hidden h-12 w-12"
                      aria-hidden="true"
                    >
                      <img
                        v-if="data.user.avatarUrl"
                        class="relative rounded-full h-full w-full"
                        :src="data.user.avatarUrl"
                        alt=""
                      />
                      <div v-else class="relative rounded-full h-full w-full bg-gray-200" />
                    </div>
                    <div class="ml-5 rounded-md shadow-sm">
                      <div
                        class="group relative border border-gray-300 rounded-md py-2 px-3 flex items-center justify-center hover:bg-gray-50 focus-within:ring-2 focus-within:ring-offset-2 focus-within:ring-light-blue-500"
                      >
                        <label
                          for="user_photo"
                          class="relative text-sm leading-4 font-medium text-gray-700 pointer-events-none"
                        >
                          <span>Change</span>
                          <span class="sr-only"> user photo</span>
                        </label>
                        <input
                          id="user_photo"
                          accept="image/jpeg, image/png"
                          name="user_photo"
                          type="file"
                          class="absolute w-full h-full opacity-0 cursor-pointer border-gray-300 rounded-md"
                          @change="uploadAvatar"
                        />
                      </div>
                    </div>
                  </div>
                </div>

                <!-- Desktop version! -->
                <div class="hidden relative rounded-full overflow-hidden lg:block">
                  <img
                    v-if="data.user.avatarUrl"
                    class="relative rounded-full w-40 h-40"
                    :src="data.user.avatarUrl"
                    alt=""
                  />
                  <div v-else class="relative rounded-full w-40 h-40 bg-gray-200" />

                  <label
                    for="user-photo"
                    class="absolute inset-0 w-full h-full bg-black bg-opacity-75 flex items-center justify-center text-sm font-medium text-white opacity-0 hover:opacity-100 focus-within:opacity-100"
                  >
                    <span>Change</span>
                    <span class="sr-only"> user photo</span>
                    <input
                      id="user-photo"
                      accept="image/jpeg, image/png"
                      type="file"
                      name="user-photo"
                      class="absolute inset-0 w-full h-full opacity-0 cursor-pointer border-gray-300 rounded-md"
                      @change="uploadAvatar"
                    />
                  </label>
                </div>
              </div>
            </div>
          </div>

          <Banner v-if="status_success" status="success" message="Updated!" />
          <Banner
            v-if="status_failed"
            status="error"
            class="mb-2"
            message="Sorry, your profile could not be updated right now. Try again later."
          />
          <Banner
            v-if="status_avatar_failed"
            status="error"
            class="mb-2"
            message="Sorry, your avatar could not be updated. Please try again later."
          />
          <Banner
            v-if="status_avatar_success"
            status="success"
            class="mb-2"
            message="Your avatar has been updated!"
          />

          <Integrations v-if="isGitHubEnabled" :user="data.user" :git-hub-app="data.gitHubApp" />

          <div>
            <div>
              <h2 class="text-lg leading-6 font-medium text-gray-900">Notifications and emails</h2>
              <p class="mt-1 text-sm text-gray-500">If you want Sturdy to contact you</p>
            </div>

            <VerifyEmail :email-verified="data.user.emailVerified" />

            <NotificationPreferences
              :features="features"
              :preferences="data.user.notificationPreferences"
              :email-verified="data.user.emailVerified"
            />

            <ul class="divide-y divide-gray-200">
              <li class="flex items-center justify-between">
                <div class="flex flex-col">
                  <p class="text-sm font-medium text-gray-900">Newsletters</p>
                  <p v-if="data.user.notificationsReceiveNewsletter" class="text-sm text-gray-500">
                    You're subscribed to the Sturdy newsletter.
                  </p>
                  <p v-else class="text-sm text-gray-500">
                    You'll not receive future newsletters by email.
                  </p>
                </div>
                <div class="w-16 justify-around items-center flex">
                  <input
                    ref="notificationsReceiveNewsletter"
                    v-model="userNotificationsReceiveNewsletter"
                    :checked="userNotificationsReceiveNewsletter"
                    type="checkbox"
                    class="focus:ring-blue-500 h-4 w-4 text-blue-600 border-gray-300 rounded"
                  />
                </div>
              </li>
            </ul>
          </div>

          <div class="mt-16 flex justify-end">
            <Button type="button" @click="refresh">Cancel</Button>
            <Button type="submit" color="blue" class="ml-5"> Save </Button>
          </div>
        </div>
      </form>
    </div>
  </PaddedApp>
</template>

<script>
import http from '../http'
import Banner from '../components/shared/Banner.vue'
import { gql, useMutation, useQuery } from '@urql/vue'
import Button from '../components/shared/Button.vue'
import { ref, watch, toRefs } from 'vue'
import NotificationPreferences from '../components/user/NotificationPreferences.vue'
import VerifyEmail from '../components/user/VerifyEmail.vue'
import PaddedApp from '../layouts/PaddedApp.vue'
import Integrations, {
  INTEGRATIONS_GITHUB_APP_FRAGMENT,
  INTEGRATIONS_USER_FRAGMENT,
} from '../organisms/user/Integrations.vue'
import { Feature } from '../__generated__/types'

export default {
  components: {
    PaddedApp,
    Banner,
    Button,
    NotificationPreferences,
    VerifyEmail,
    Integrations,
  },
  props: {
    features: {
      type: Array,
      required: true,
    },
  },
  setup(props) {
    const { features } = toRefs(props)
    const isGitHubEnabled = features.value.includes(Feature.GitHub)

    let { data, fetching, error, executeQuery } = useQuery({
      query: gql`
        query UserPage($isGitHubEnabled: Boolean!) {
          gitHubApp @include(if: $isGitHubEnabled) {
            ...IntegrationsGitHubApp
          }
          user {
            id
            name
            email
            emailVerified
            avatarUrl
            notificationsReceiveNewsletter
            notificationPreferences {
              type
              channel
              enabled
            }
            ...IntegrationsUser @include(if: $isGitHubEnabled)
          }
        }
        ${INTEGRATIONS_GITHUB_APP_FRAGMENT}
        ${INTEGRATIONS_USER_FRAGMENT}
      `,
      requestPolicy: 'cache-and-network',
      variables: {
        isGitHubEnabled,
      },
    })

    const { executeMutation: updateUserResult } = useMutation(gql`
      mutation UpdateUser(
        $name: String
        $password: String
        $email: String
        $notificationsReceiveNewsletter: Boolean
      ) {
        updateUser(
          input: {
            name: $name
            email: $email
            password: $password
            notificationsReceiveNewsletter: $notificationsReceiveNewsletter
          }
        ) {
          id
          name
          email
          notificationsReceiveNewsletter
        }
      }
    `)

    // One way data bindings
    let userName = ref('')
    let userEmail = ref('')
    let userPassword = ref('')
    let userNotificationsReceiveNewsletter = ref(false)
    watch(data, () => {
      if (data && data.value && data.value.user) {
        userName.value = data.value.user.name
        userEmail.value = data.value.user.email
        userNotificationsReceiveNewsletter.value = data.value.user.notificationsReceiveNewsletter
      }
    })

    return {
      isGitHubEnabled,

      data,
      fetching,
      error,

      userName,
      userEmail,
      userPassword,
      userNotificationsReceiveNewsletter,

      refresh() {
        executeQuery({
          requestPolicy: 'network-only',
        })
      },

      async updateUser(name, password, email, notificationsReceiveNewsletter) {
        const variables = {
          name,
          email,
          password,
          notificationsReceiveNewsletter,
        }
        await updateUserResult(variables).then((result) => {
          if (result.error) {
            throw new Error(result.error)
          }
          console.log('update user', result)
        })
      },
    }
  },
  data: function () {
    return {
      status_success: false,
      status_failed: false,
      status_avatar_failed: false,
      status_avatar_success: false,
    }
  },
  watch: {
    error: function (err) {
      if (err) throw err
    },
  },
  methods: {
    save() {
      this.updateUser(
        this.userName,
        this.userPassword,
        this.userEmail,
        this.userNotificationsReceiveNewsletter
      )
        .then(() => {
          this.status_success = true
          this.status_failed = false
        })
        .catch(() => {
          this.status_success = false
          this.status_failed = true
        })
    },
    uploadAvatar(event) {
      this.status_avatar_success = false
      this.status_avatar_failed = false

      let files = event.target.files

      const formData = new FormData()
      formData.append('file', files[0])

      fetch(http.url('v3/user/update-avatar'), {
        method: 'POST',
        body: formData,
        credentials: 'include',
      })
        .then(http.checkStatus)
        .then((response) => response.json())
        .then(() => {
          this.status_avatar_success = true
          this.status_avatar_failed = false
          this.emitter.emit('reload-user', {})
        })
        .catch(() => {
          this.status_avatar_success = false
          this.status_avatar_failed = true
        })
        .finally(this.refresh)
    },
  },
  computed: {
    passwordEnabled() {
      return this.features.includes(Feature.PasswordAuth)
    },
  },
}
</script>
