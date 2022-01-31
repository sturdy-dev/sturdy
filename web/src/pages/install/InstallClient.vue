<template>
  <PaddedApp>
    <div>
      <Banner status="info">
        These are installation instructions for the older Sturdy CLI. Looking to
        <router-link :to="{ name: 'download' }" class="text-black underline"
          >download the Sturdy app for Mac and Windows</router-link
        >?
      </Banner>

      <div class="bg-white shadow overflow-hidden sm:rounded-lg mt-4">
        <div class="px-4 py-5 sm:px-6 flex justify-between">
          <div>
            <h3 class="text-lg leading-6 font-medium text-gray-900">Install Sturdy</h3>
            <p class="mt-1 max-w-2xl text-sm text-gray-500">Installation guide</p>
          </div>
          <div>
            <Button
              v-if="data"
              color="green"
              @click="
                $router.push({
                  name: 'codebaseHome',
                  params: { codebaseSlug: $route.params.codebaseSlug },
                })
              "
            >
              <span>Continue the setup of {{ data.codebase.name }}</span>
              <ArrowCircleRightIcon class="-mr-1 ml-2 h-5 w-5 text-green-700" />
            </Button>
          </div>
        </div>
        <div class="border-t border-gray-200 px-4 py-5 sm:px-6">
          <div>
            <div class="flow-root">
              <InstallationInstructions class="-my-5 py-5" />
            </div>

            <div class="mt-6">
              <Button
                v-if="data"
                color="green"
                @click="
                  $router.push({
                    name: 'codebaseHome',
                    params: { codebaseSlug: $route.params.codebaseSlug },
                  })
                "
              >
                <span>Continue setup of {{ data.codebase.name }}</span>
                <ArrowCircleRightIcon class="-mr-1 ml-2 h-5 w-5 text-green-700" />
              </Button>
              <router-link
                v-else
                class="w-full flex justify-center items-center px-4 py-2 border border-gray-300 shadow-sm text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50"
                to="codebases"
              >
                <span>Create your first codebase</span>
                <ArrowCircleRightIcon class="-mr-1 ml-2 h-5 w-5 text-green-700" />
              </router-link>
            </div>
          </div>
        </div>
      </div>
    </div>
  </PaddedApp>
</template>

<script>
import { ArrowCircleRightIcon } from '@heroicons/vue/solid'
import InstallationInstructions from '../../components/install/InstallationInstructions.vue'
import { gql, useQuery } from '@urql/vue'
import { useRoute } from 'vue-router'
import Button from '../../components/shared/Button.vue'
import { computed, ref, watch } from 'vue'
import { Banner } from '../../atoms'
import PaddedApp from '../../layouts/PaddedApp.vue'

export default {
  name: 'InstallClient',
  components: {
    PaddedApp,
    Banner,
    InstallationInstructions,
    ArrowCircleRightIcon,
    Button,
  },
  setup() {
    let route = useRoute()
    let codebaseSlug = ref(route.params.codebaseSlug)
    watch(route, () => {
      codebaseSlug.value = route.params.codebaseSlug
    })

    let { data } = useQuery({
      query: gql`
        query InstallClient($shortID: ID) {
          codebase(shortID: $shortID) {
            id
            name
          }
        }
      `,
      pause: computed(() => !codebaseSlug.value),
      variables: {
        shortID: codebaseSlug.value,
      },
      requestPolicy: 'cache-and-network',
    })

    return {
      data,
    }
  },
  watch: {
    '$route.params.codebaseSlug': function () {
      this.setVisited()
    },
    error: function (err) {
      if (err) throw err
    },
  },
  mounted() {
    this.setVisited()
  },
  methods: {
    setVisited() {
      if (this.$route.params.codebaseSlug) {
        console.log('set visit', this.$route.params.codebaseSlug)
        localStorage.setItem('visitedInstallClient', true)
      } else {
        console.log('not set')
      }
    },
  },
}
</script>
