<template>
  <div v-if="item.change" class="relative flex items-start space-x-3">
    <div>
      <div class="relative px-1">
        <Avatar :author="item.author" size="8" />
      </div>
    </div>

    <div class="min-w-0 flex-1 py-1.5">
      <div class="text-sm text-gray-500">
        <a class="font-medium text-gray-900">
          {{ item.author.name }}
        </a>
        {{ ' ' }}
        created
        {{ ' ' }}
        <router-link
          :to="{
            name: 'codebaseChange',
            params: {
              codebaseSlug: codebaseSlug,
              selectedChangeID: item.change.id,
            },
          }"
          class="font-medium text-gray-900"
        >
          {{ item.change.title }}
        </router-link>
        {{ ' ' }}
        <span class="flex whitespace-nowrap mr-1">
          <StatusDetails :statuses="item.change.statuses" :show-text="false" />
          <RelativeTime :date="createdAt" />
        </span>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import Avatar from '../../components/shared/Avatar.vue'
import RelativeTime from '../../atoms/RelativeTime.vue'
import { PropType, toRef } from 'vue'
import { WorkspaceCreatedChangeActivityFragment } from './__generated__/ActivityCreatedChange'
import { STATUS_FRAGMENT } from '../../components/statuses/StatusBadge.vue'
import { gql } from '@urql/vue'
import { useUpdatedChangesStatuses } from '../../subscriptions/useUpdatedChangesStatuses'
import StatusDetails from '../../components/statuses/StatusDetails.vue'

export const WORKSPACE_ACTIVITY_CREATED_CHANGE_FRAGMENT = gql`
  fragment WorkspaceCreatedChangeActivity on WorkspaceCreatedChangeActivity {
    author {
      id
      name
      avatarUrl
    }
    createdAt
    change {
      id
      title
      trunkCommitID
      statuses {
        ...Status
      }
    }
  }
  ${STATUS_FRAGMENT}
`

export default {
  name: 'WorkspaceActivityCreatedChange',
  components: { Avatar, StatusDetails, RelativeTime },
  props: {
    item: {
      type: Object as PropType<WorkspaceCreatedChangeActivityFragment>,
      required: true,
    },
    codebaseSlug: {
      type: String,
      required: true,
    },
  },
  setup(props) {
    const item = toRef(props, 'item')
    if (item.value.change) useUpdatedChangesStatuses([item.value.change.id])
  },
  computed: {
    createdAt() {
      return new Date(this.item.createdAt * 1000)
    },
  },
}
</script>
