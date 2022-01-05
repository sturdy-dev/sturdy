<template>
  <div v-if="codebase.changes.length > 0" class="sm:flex sm:flex-row sm:gap-2">
    <router-link
      class="border rounded-md shadow-sm flex-1 px-4 py-3 items-stretch flex flex-col mb-2 sm:mb-0 hover:bg-gray-50"
      :to="{
        name: 'codebaseChangelog',
        params: { codebaseSlug: codebase.shortID, selectedChangeID: latestChange.id },
      }"
    >
      <div class="flex flex-row gap-2 mb-3 flex-none">
        <Avatar :author="latestChange.author" size="10" />
        <div class="text-gray-500 text-sm">
          <div class="font-semibold text-black">{{ latestChange.author.name }}</div>
          <div class="flex">
            <div class="mr-1">
              <StatusBadge :statuses="latestChange.statuses" />
            </div>
            <div>shared {{ timeAgo(latestChange.createdAt) }}</div>
          </div>
        </div>
      </div>
      <div class="text-sm flex-1 prose line-clamp-4" v-html="latestChange.description" />
      <div class="flex flex-row justify-end items-center flex-none">
        <div v-if="latestChange.comments.length > 0" class="flex-none mr-2">
          <ChangeCommentsIndicator :change="latestChange" />
        </div>
      </div>
    </router-link>

    <div v-if="codebase.changes.length > 1" class="flex flex-col gap-2 flex-1">
      <router-link
        v-for="change in codebase.changes.slice(1)"
        :key="change.id"
        class="border rounded-md shadow-sm px-4 py-2 hover:bg-gray-50"
        :to="{
          name: 'codebaseChangelog',
          params: { codebaseSlug: codebase.shortID, selectedChangeID: change.id },
        }"
      >
        <div class="flex flex-row gap-2 items-center justify-between">
          <div class="flex-none">
            <Avatar :author="change.author" size="6" />
          </div>
          <div class="text-sm flex-1">{{ change.title }}</div>
          <div v-if="change.comments.length > 0" class="flex-none">
            <ChangeCommentsIndicator :change="change" />
          </div>
          <div class="flex text-sm text-gray-500">
            <div class="mr-1">
              <StatusBadge :statuses="change.statuses" />
            </div>
            {{ timeAgo(change.createdAt) }}
          </div>
        </div>
      </router-link>

      <router-link
        :to="{ name: 'codebaseChangelog', params: { codebaseSlug: codebase.shortID } }"
        class="text-sm flex flex-row items-center gap-1 text-gray-500 hover:bg-gray-100 self-end px-2 py-1 rounded-md"
      >
        Older Changes
        <div class="h-4 w-4">
          <ArrowRightIcon />
        </div>
      </router-link>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent, PropType, toRefs } from 'vue'
import gql from 'graphql-tag'
import { TopOfChangelogFragment } from './__generated__/TopOfChangelogWidget'
import { ArrowRightIcon } from '@heroicons/vue/outline'
import Avatar from '../shared/Avatar.vue'
import time from '../../time'
import ChangeCommentsIndicator, { CHANGE_COMMENTS } from './ChangeCommentsIndicator.vue'
import StatusBadge, { STATUS_FRAGMENT } from '../statuses/StatusBadge.vue'
import { useUpdatedChangesStatuses } from '../../subscriptions/useUpdatedChangesStatuses'

export const TOP_OF_CHANGELOG = gql`
  fragment TopOfChangelog on Codebase {
    id
    shortID
    changes(input: { limit: 4 }) {
      id
      title
      description
      createdAt
      author {
        id
        name
        avatarUrl
      }
      ...ChangeComments
      statuses {
        ...Status
      }
    }
  }
  ${CHANGE_COMMENTS}
  ${STATUS_FRAGMENT}
`

type Change = TopOfChangelogFragment['changes'][number]

export default defineComponent({
  name: 'TopOfChangelogWidget',
  components: { Avatar, ArrowRightIcon, ChangeCommentsIndicator, StatusBadge },
  props: {
    codebase: {
      required: true,
      type: Object as PropType<TopOfChangelogFragment>,
    },
  },
  setup(props) {
    const { codebase } = toRefs(props)
    const changeIDs = codebase.value.changes.map((c) => c.id)
    useUpdatedChangesStatuses(changeIDs)
  },
  computed: {
    latestChange(): Change {
      return this.codebase.changes[0]
    },
  },
  methods: {
    timeAgo(timestamp: number) {
      return time.getRelativeTime(new Date(timestamp * 1000))
    },
  },
})
</script>
