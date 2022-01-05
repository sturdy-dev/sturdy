<template>
  <TransitionRoot as="template" :show="open">
    <Dialog
      as="div"
      static
      class="fixed inset-0 overflow-hidden z-50"
      :open="open"
      @close="$emit('close')"
    >
      <div class="absolute inset-0 overflow-hidden">
        <DialogOverlay class="absolute inset-0" />

        <div class="fixed inset-y-0 right-0 pl-10 max-w-full flex">
          <TransitionChild
            as="template"
            enter="transform transition ease-in-out duration-100 sm:duration-200"
            enter-from="translate-x-full"
            enter-to="translate-x-0"
            leave="transform transition ease-in-out duration-100 sm:duration-200"
            leave-from="translate-x-0"
            leave-to="translate-x-full"
          >
            <div class="w-screen max-w-md">
              <div class="h-full flex flex-col py-6 bg-white shadow-xl overflow-y-scroll">
                <div class="px-4 sm:px-6">
                  <div class="flex items-start justify-between">
                    <DialogTitle class="text-lg font-medium text-gray-900">
                      Notifications
                    </DialogTitle>
                    <div class="ml-3 h-7 flex items-center">
                      <button
                        class="bg-white rounded-md text-gray-400 hover:text-gray-500 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
                        @click="$emit('close')"
                      >
                        <span class="sr-only">Close panel</span>
                        <XIcon class="h-6 w-6" aria-hidden="true" />
                      </button>
                    </div>
                  </div>
                </div>
                <div class="mt-6 relative flex-1 px-4 sm:px-6">
                  <NotificationFeed
                    :notifications="notifications"
                    :user="user"
                    @close="$emit('close')"
                  />
                </div>
              </div>
            </div>
          </TransitionChild>
        </div>
      </div>
    </Dialog>
  </TransitionRoot>
</template>

<script>
import {
  Dialog,
  DialogOverlay,
  DialogTitle,
  TransitionChild,
  TransitionRoot,
} from '@headlessui/vue'
import { XIcon } from '@heroicons/vue/outline'
import NotificationFeed from './Feed.vue'

export default {
  components: {
    NotificationFeed,
    Dialog,
    DialogOverlay,
    DialogTitle,
    TransitionChild,
    TransitionRoot,
    XIcon,
  },
  props: {
    open: Boolean,
    notifications: Object,
    user: Object,
  },
  emits: ['close'],
}
</script>
