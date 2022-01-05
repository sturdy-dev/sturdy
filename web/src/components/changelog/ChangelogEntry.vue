<template>
  <div
    class="relative flex items-start space-x-3 group cursor-pointer"
    @click.stop.prevent="select"
  >
    <div class="relative">
      <Avatar :author="change.author" size="10" />
      <span
        v-if="change.comments.length > 0"
        class="absolute -bottom-2 -right-2 bg-white rounded-tl px-0.5 py-px"
      >
        <ChatAltIcon
          class="h-5 w-5 group-hover:text-gray-600"
          :class="[isSelected ? 'text-gray-500' : 'text-gray-400']"
          aria-hidden="true"
        />
      </span>
    </div>
    <div class="min-w-0 flex-1">
      <div>
        <div>
          <a
            class="text-sm font-medium group-hover:text-gray-900 cursor-pointer"
            :class="[isSelected ? 'text-gray-900' : 'text-gray-500']"
          >
            {{ change.title }}
          </a>
        </div>
        <div
          class="text-sm font-normal group-hover:text-gray-900 cursor-pointer"
          :class="[isSelected ? 'text-gray-800' : 'text-gray-400']"
        ></div>
        <div
          v-if="change.createdAt > 0"
          class="flex text-sm font-normal group-hover:text-gray-900 cursor-pointer"
          :class="[isSelected ? 'text-gray-800' : 'text-gray-400']"
        >
          <div class="mr-1">
            <StatusBadge :statuses="change.statuses" />
          </div>
          {{ friendly_ago(change.createdAt) }}
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ChatAltIcon } from '@heroicons/vue/outline'
import Avatar from '../shared/Avatar.vue'
import time from '../../time'
import StatusBadge from '../statuses/StatusBadge.vue'

export default {
  name: 'ChangeLogEntry',
  components: { ChatAltIcon, Avatar, StatusBadge },
  props: ['isSelected', 'change'],
  data() {
    return {}
  },
  methods: {
    select(ev) {
      this.$emit('select')
    },
    friendly_ago(ts) {
      return time.getRelativeTime(new Date(ts * 1000))
    },
  },
}
</script>
