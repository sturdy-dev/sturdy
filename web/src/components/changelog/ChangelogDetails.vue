<template>
  <div
    class="xl:space-y-5 flex xl:block justify-between md:justify-start md:space-x-12 xl:space-x-0"
  >
    <div class="space-y-5">
      <div>
        <ul role="list" class="xl:mt-3 space-y-3">
          <li class="flex justify-start">
            <a class="flex items-center space-x-3">
              <div class="flex-shrink-0">
                <Avatar v-if="changeData?.change" :author="changeData.change.author" size="5" />
                <div v-else class="h-5 w-5 bg-gray-300 rounded-full animate-pulse"></div>
              </div>
              <div
                v-if="changeData?.change?.author?.name"
                class="text-sm font-medium text-gray-900"
              >
                {{ changeData.change.author.name }}
              </div>
              <div v-else class="h-4 rounded-md bg-gray-300 w-3/4 animate-pulse"></div>
            </a>
          </li>
        </ul>
      </div>
      <div v-if="changeData?.change" class="flex items-center space-x-2">
        <ChatAltIcon class="h-5 w-5 text-gray-400" aria-hidden="true" />
        <span
          v-if="changeData?.change?.comments.length === 0"
          class="text-gray-900 text-sm font-medium"
        >
          No comments
        </span>
        <span v-else class="text-gray-900 text-sm font-medium">
          {{ changeData?.change?.comments.length }}
          {{ changeData?.change?.comments.length === 1 ? 'comment' : 'comments' }}
        </span>
      </div>
      <div
        v-else
        class="flex items-center space-x-2 h-4 rounded-md w-3/4 bg-gray-300 animate-pulse"
      ></div>
    </div>

    <div class="space-y-5">
      <div v-if="changeData?.change?.createdAt > 0" class="flex items-center space-x-2">
        <CalendarIcon class="h-5 w-5 text-gray-400" aria-hidden="true" />
        <span class="text-gray-900 text-sm font-medium">
          {{ friendly_ago(changeData?.change?.createdAt) }}
        </span>
      </div>
      <a v-if="github_link" :href="github_link" class="flex items-center space-x-2">
        <CheckCircleIcon class="h-5 w-5 text-green-300" aria-hidden="true" />
        <span
          v-if="githubIntegration.gitHubIsSourceOfTruth"
          class="text-gray-900 text-sm font-medium"
        >
          Synced from GitHub
        </span>
        <span v-else class="text-gray-900 text-sm font-medium"> Synced to GitHub </span>
      </a>
    </div>

    <div class="space-y-5">
      <div class="flex items-center space-x-2">
        <StatusDetails :statuses="changeData?.change?.statuses" />
      </div>
    </div>

    <div v-if="changeData?.change">
      <div class="flex items-center space-x-2">
        <Button v-if="fetchingZipDownload" size="wider" disabled>
          <Spinner class="-ml-1 mr-2 h-5 w-5 text-gray-400" aria-hidden="true" />
          <span>Preparing</span>
        </Button>
        <Button
          v-else-if="!fetchingZipDownload && generateZipData && didTriggerDownload"
          size="wider"
          disabled
        >
          <DownloadIcon class="-ml-1 mr-2 h-5 w-5 text-gray-400" aria-hidden="true" />
          <span>Opened download in new tab</span>
        </Button>
        <Button v-else size="wider" @click="zipDownload">
          <DownloadIcon class="-ml-1 mr-2 h-5 w-5 text-gray-400" aria-hidden="true" />
          <span>Download as Zip</span>
        </Button>
      </div>
    </div>
  </div>
</template>
<script>
import Avatar from '../shared/Avatar.vue'
import time from '../../time'
import { CalendarIcon, ChatAltIcon, CheckCircleIcon } from '@heroicons/vue/solid'
import StatusDetails from '../statuses/StatusDetails.vue'
import Spinner from '../shared/Spinner.vue'
import DownloadIcon from '@heroicons/vue/outline/DownloadIcon'
import Button from '../shared/Button.vue'
import { ref, toRef, watch } from 'vue'
import { gql, useQuery } from '@urql/vue'

export default {
  name: 'CodebaseChangelogDetails',
  components: {
    Avatar,
    CalendarIcon,
    ChatAltIcon,
    CheckCircleIcon,
    StatusDetails,
    Spinner,
    DownloadIcon,
    Button,
  },
  props: {
    changeData: {},
    githubIntegration: {},
  },
  setup(props) {
    let changeRef = toRef(props, 'changeData')
    let loadedChangeID = ref(changeRef.value?.change?.id)

    let pauseGenerateZip = ref(true)
    let { data: generateZipData, fetching: fetchingZipDownload } = useQuery({
      query: gql`
        query ChangeDetailsDownloadZip($changeID: ID!) {
          change(id: $changeID) {
            id
            downloadZip {
              id
              url
            }
          }
        }
      `,
      pause: pauseGenerateZip,
      variables: {
        changeID: loadedChangeID,
      },
    })

    let didTriggerDownload = ref(false)
    watch(generateZipData, () => {
      let url = generateZipData.value?.change?.downloadZip?.url
      if (url && !didTriggerDownload.value) {
        didTriggerDownload.value = true // only once
        window.open(url, '_blank')?.focus()
      }
    })

    watch(changeRef, () => {
      // New change, reset data
      if (changeRef?.value?.change?.id && changeRef.value.change.id !== loadedChangeID.value) {
        loadedChangeID.value = changeRef.value.change?.id
        didTriggerDownload.value = false
        pauseGenerateZip.value = true
      }
    })

    return {
      generateZipData,
      fetchingZipDownload,
      didTriggerDownload,

      zipDownload() {
        console.log('change-id', changeRef.value?.change?.id, loadedChangeID)
        pauseGenerateZip.value = false
      },
    }
  },
  computed: {
    github_link() {
      if (!this.githubIntegration || !this.changeData?.change || !this.githubIntegration.enabled) {
        return false
      }
      return (
        'https://github.com/' +
        this.githubIntegration.owner +
        '/' +
        this.githubIntegration.name +
        '/commit/' +
        this.changeData.change.trunkCommitID
      )
    },
  },
  methods: {
    friendly_ago(ts) {
      return time.getRelativeTime(new Date(ts * 1000))
    },
  },
}
</script>
