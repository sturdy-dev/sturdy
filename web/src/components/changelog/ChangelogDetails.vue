<template>
  <div
    class="grid auto-rows-min auto-cols-min grid-flow-row-dense grid-cols-3 xl:grid-cols-1 gap-4"
  >
    <div class="flex items-center space-x-3">
      <div class="flex-shrink-0">
        <Avatar :author="change.author" size="5" />
      </div>
      <div v-if="change.author?.name" class="text-sm font-medium text-gray-900">
        {{ change.author.name }}
      </div>
      <div v-else class="h-4 rounded-md bg-gray-300 w-3/4 animate-pulse"></div>
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

    <div class="flex items-center space-x-2">
      <StatusDetails :statuses="change.statuses" />
    </div>

    <DownloadChangeButton v-if="isDownloadAvailable" :change-id="change.id" />

    <div class="flex gap-2 cols-span-2 text-sm">
      <router-link v-if="!!change.parent" :to="{ params: { id: change.parent.id } }">
        <Button class="flex gap-2" size="small">
          <ArrowLeftIcon class="h-4 w-4" />
          Previous
        </Button>
      </router-link>
      <Button v-else :disabled="true" class="flex gap-2" size="small">
        <ArrowLeftIcon class="h-4 w-4" />
        Previous
      </Button>

      <router-link v-if="!!change.child" :to="{ params: { id: change.child.id } }">
        <Button class="flex gap-2" size="small">
          Next
          <ArrowRightIcon class="h-4 w-4" />
        </Button>
      </router-link>
      <Button v-else :disabled="true" class="flex gap-2" size="small">
        Next
        <ArrowRightIcon class="h-4 w-4" />
      </Button>
    </div>
  </div>
</template>

<script lang="ts">
import { ref, inject, computed, Ref, defineComponent, PropType } from 'vue'

import RelativeTime from '../../atoms/RelativeTime.vue'
import Avatar from '../shared/Avatar.vue'
import {
  CalendarIcon,
  ChatAltIcon,
  CheckCircleIcon,
  ArrowSmRightIcon as ArrowRightIcon,
  ArrowSmLeftIcon as ArrowLeftIcon,
} from '@heroicons/vue/solid'
import DownloadChangeButton from '../../molecules/DownloadChangeButton.vue'
import StatusDetails, { STATUS_FRAGMENT } from '../statuses/StatusDetails.vue'
import Button from '../shared/Button.vue'
import { gql } from '@urql/vue'
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

    parent {
      id
    }
    child {
      id
    }
  }
  ${AUTHOR}
  ${STATUS_FRAGMENT}
`

export default defineComponent({
  components: {
    RelativeTime,
    DownloadChangeButton,
    Avatar,
    CalendarIcon,
    ChatAltIcon,
    CheckCircleIcon,
    StatusDetails,
    Button,
    ArrowLeftIcon,
    ArrowRightIcon,
  },
  props: {
    change: {
      type: Object as PropType<ChangelogDetails_ChangeFragment>,
      required: true,
    },
  },
  setup() {
    const features = inject<Ref<Array<Feature>>>('features', ref([]))
    const isDownloadAvailable = computed(() => features?.value?.includes(Feature.DownloadChanges))
    return {
      isDownloadAvailable,
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
