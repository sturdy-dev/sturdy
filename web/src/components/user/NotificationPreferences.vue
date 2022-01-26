<template>
  <fieldset>
    <div class="relative flex items-start py-4">
      <div class="min-w-0 flex-1 text-sm"></div>
      <div class="w-16 justify-around flex items-center h-5">
        <p class="font-medium text-gray-700">Web</p>
      </div>
      <div class="w-16 flex justify-around items-center h-5">
        <p class="font-medium text-gray-700">Email</p>
      </div>
    </div>

    <div v-for="preference in grouped" :key="preference" class="relative flex items-start py-4">
      <div class="min-w-0 flex-1 text-sm">
        <label class="font-medium text-gray-700">{{ preference.title }}</label>
        <p id="candidates-description" class="text-gray-500">
          {{ preference.description }}
        </p>
      </div>
      <div class="w-16 flex items-center h-5 justify-around">
        <input
          v-model="preference.webEnabled"
          type="checkbox"
          class="focus:ring-indigo-500 h-4 w-4 text-indigo-600 border-gray-300 rounded"
          @click="preference.toggleWeb"
        />
      </div>
      <div class="w-16 flex items-center h-5 justify-around">
        <input
          v-model="preference.emailEnabled"
          :disabled="!emailVerified"
          :class="{ 'disabled:opacity-50': !emailVerified }"
          type="checkbox"
          class="focus:ring-indigo-500 h-4 w-4 text-indigo-600 border-gray-300 rounded"
          @click="preference.toggleEmail"
        />
      </div>
    </div>
  </fieldset>
</template>

<script lang="ts">
import { defineComponent, PropType, inject, Ref, ref, computed } from 'vue'
import {
  NotificationPreference,
  NotificationType,
  NotificationChannel,
  Feature,
} from '../../__generated__/types'
import { useUpdateNotificationPreference } from '../../mutations/useUpdateNotificationPreference'

type Preference = {
  title: string
  description: string
  emailEnabled: boolean
  webEnabled: boolean
  toggleWeb: () => void
  toggleEmail: () => void
}

export default defineComponent({
  name: 'NotificationPreferences',
  props: {
    preferences: {
      type: Array as PropType<NotificationPreference[]>,
      required: true,
    },
    emailVerified: {
      type: Boolean,
      required: true,
    },
  },
  setup() {
    const features = inject<Ref<Array<Feature>>>('features', ref([]))
    const isGitHubEnabled = computed(() => features?.value?.includes(Feature.GitHub))

    const updateNotificationPreferenceResult = useUpdateNotificationPreference()
    return {
      isGitHubEnabled,

      async updateNotificationPreference(
        type: NotificationType,
        channel: NotificationChannel,
        enabled: boolean
      ) {
        await updateNotificationPreferenceResult({
          type,
          channel,
          enabled,
        })
      },
    }
  },
  computed: {
    grouped(): Preference[] {
      const enabledByTypeByChannel = new Map()
      this.preferences.forEach((p) => {
        if (p.type === NotificationType.GitHubRepositoryImported && !this.isGitHubEnabled) {
          return
        }
        if (!enabledByTypeByChannel.get(p.type)) {
          const m = new Map()
          m.set(p.channel, p.enabled)
          enabledByTypeByChannel.set(p.type, m)
        } else {
          enabledByTypeByChannel.get(p.type).set(p.channel, p.enabled)
        }
      })

      const result = [] as Preference[]
      enabledByTypeByChannel.forEach((enabledByChannel, typ) => {
        const emailEnabled = enabledByChannel.get(NotificationChannel.Email)
        const webEnabled = enabledByChannel.get(NotificationChannel.Web)
        result.push({
          title: this.title(typ as NotificationType),
          description: this.description(typ as NotificationType),
          emailEnabled: emailEnabled,
          webEnabled: webEnabled,
          toggleWeb: () =>
            this.updateNotificationPreference(typ, NotificationChannel.Web, !webEnabled),
          toggleEmail: () =>
            this.updateNotificationPreference(typ, NotificationChannel.Email, !emailEnabled),
        })
      })

      return result.sort((p1, p2) => p1.title.localeCompare(p2.title))
    },
  },
  methods: {
    description(typ: NotificationType): string {
      switch (typ) {
        case NotificationType.Comment:
          return 'Get notified when someone writes a new comment'
        case NotificationType.NewSuggestion:
          return 'Get notified when someone sends you a new suggestion'
        case NotificationType.RequestedReview:
          return 'Get notified when someone requests your review'
        case NotificationType.Review:
          return 'Get notified when someone sends you a review'
        case NotificationType.GitHubRepositoryImported:
          return 'Get notified when a new repository is imported'
        default:
          throw Error(`unsupported type ${typ}`)
      }
    },

    title(typ: NotificationType): string {
      switch (typ) {
        case NotificationType.Comment:
          return 'New comments'
        case NotificationType.NewSuggestion:
          return 'New suggestions'
        case NotificationType.RequestedReview:
          return 'Review requested'
        case NotificationType.Review:
          return 'Review received'
        case NotificationType.GitHubRepositoryImported:
          return 'GitHub repository imported'
        default:
          throw Error(`unsupported type ${typ}`)
      }
    },
  },
})
</script>
