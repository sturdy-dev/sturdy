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
                <Avatar :author="change.author" size="5" />
              </div>
              <div v-if="change.author?.name" class="text-sm font-medium text-gray-900">
                {{ change.author.name }}
              </div>
              <div v-else class="h-4 rounded-md bg-gray-300 w-3/4 animate-pulse"></div>
            </a>
          </li>
        </ul>
      </div>
      <div class="flex items-center space-x-2">
        <ChatAltIcon class="h-5 w-5 text-gray-400" aria-hidden="true" />
        <span v-if="change.comments.length === 0" class="text-gray-900 text-sm font-medium">
          No comments
        </span>
        <span v-else class="text-gray-900 text-sm font-medium">
          {{ change.comments.length }}
          {{ change.comments.length === 1 ? 'comment' : 'comments' }}
        </span>
      </div>
    </div>

    <div class="space-y-5">
      <div v-if="change.createdAt > 0" class="flex items-center space-x-2">
        <CalendarIcon class="h-5 w-5 text-gray-400" aria-hidden="true" />
        <RelativeTime class="text-gray-900 text-sm font-medium" :date="createdAt" />
      </div>
      <a v-if="github_link" :href="github_link" class="flex items-center space-x-2">
        <CheckCircleIcon class="h-5 w-5 text-green-300" aria-hidden="true" />
        <span
          v-if="gitHubIntegration?.gitHubIsSourceOfTruth"
          class="text-gray-900 text-sm font-medium"
        >
          Synced from GitHub
        </span>
        <span v-else class="text-gray-900 text-sm font-medium">Synced to GitHub</span>
      </a>
    </div>

    <div class="space-y-5">
      <div class="flex items-center space-x-2">
        <StatusDetails :statuses="change.statuses" />
      </div>
    </div>

    <div v-if="isDownloadAvailable">
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

<script lang="ts">
import { ref, toRefs, watch, inject, computed, Ref, defineComponent, PropType } from 'vue'

import RelativeTime from '../../atoms/RelativeTime.vue'
import Avatar from '../shared/Avatar.vue'
import { CalendarIcon, ChatAltIcon, CheckCircleIcon } from '@heroicons/vue/solid'
import StatusDetails, { STATUS_FRAGMENT } from '../statuses/StatusDetails.vue'
import Spinner from '../shared/Spinner.vue'
import DownloadIcon from '@heroicons/vue/outline/DownloadIcon'
import Button from '../shared/Button.vue'
import { gql, useQuery } from '@urql/vue'
import { Feature } from '../../__generated__/types'

import { AUTHOR } from '../shared/AvatarHelper'
import { ChangelogDetails_ChangeFragment } from './__generated__/ChangelogDetails'

export const CHANGE_FRAGMENT = gql`
  fragment ChangelogDetails_Change on Change {
    id
    trunkCommitID
    author {
      ...Author
    }
    comments {
      id
    }
    createdAt
    codebase {
      id
      gitHubIntegration {
        id
        enabled
        owner
        name
        gitHubIsSourceOfTruth
      }
    }
    statuses {
      ...Status
    }
  }
  ${AUTHOR}
  ${STATUS_FRAGMENT}
`

export default defineComponent({
  components: {
    RelativeTime,
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
    change: {
      type: Object as PropType<ChangelogDetails_ChangeFragment>,
      required: true,
    },
  },
  setup(props) {
    const features = inject<Ref<Array<Feature>>>('features', ref([]))
    const isDownloadAvailable = computed(() => features?.value?.includes(Feature.DownloadChanges))

    const { change } = toRefs(props)

    const pauseGenerateZip = ref(true)
    const { data: generateZipData, fetching: fetchingZipDownload } = useQuery({
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
        changeID: change.value.id,
      },
    })

    const didTriggerDownload = ref(false)
    watch(generateZipData, () => {
      let url = generateZipData.value?.change?.downloadZip?.url
      if (url && !didTriggerDownload.value) {
        didTriggerDownload.value = true // only once
        window.open(url, '_blank')?.focus()
      }
    })

    return {
      generateZipData,
      fetchingZipDownload,
      didTriggerDownload,

      isDownloadAvailable,

      zipDownload() {
        pauseGenerateZip.value = false
      },
    }
  },
  computed: {
    gitHubIntegration() {
      return this.change.codebase.gitHubIntegration
    },
    github_link() {
      const gitHubIntegration = this.change.codebase.gitHubIntegration
      return gitHubIntegration?.enabled
        ? `https://github.com/${gitHubIntegration.owner}/${gitHubIntegration.name}/commit/${this.change.trunkCommitID}`
        : undefined
    },
    createdAt() {
      return new Date(this.change.createdAt * 1000)
    },
  },
})
</script>
