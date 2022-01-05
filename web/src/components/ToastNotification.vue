<template>
  <div>
    <transition
      enter-active-class="transform ease-out duration-300 transition"
      enter-from-class="translate-y-2 opacity-0 sm:translate-y-0 sm:translate-x-2"
      enter-to-class="translate-y-0 opacity-100 sm:translate-x-0"
      leave-active-class="transition ease-in duration-100"
      leave-from-class="opacity-100"
      leave-to-class="opacity-0"
    >
      <div
        v-show="show"
        class="mb-4 max-w-sm w-full bg-white shadow-lg rounded-lg pointer-events-auto ring-1 ring-black ring-opacity-5 overflow-hidden"
      >
        <div class="p-4">
          <div class="flex items-start">
            <div class="flex-shrink-0">
              <CheckCircleIcon v-if="style === 'success'" class="h-6 w-6 text-green-400" />
              <ExclamationCircleIcon v-else class="h-6 w-6 text-red-400" />
            </div>
            <div class="ml-3 w-0 flex-1 pt-0.5">
              <p class="text-sm font-medium" :class="titleClasses">
                {{ title }}
              </p>
              <p class="mt-1 text-sm text-gray-500">
                {{ message }}
              </p>
            </div>
            <div class="ml-4 flex-shrink-0 flex">
              <button
                class="bg-white rounded-md inline-flex text-gray-400 hover:text-gray-500 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
                @click="show = false"
              >
                <span class="sr-only">Close</span>
                <XIcon class="h-5 w-5" />
              </button>
            </div>
          </div>
        </div>
      </div>
    </transition>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import { CheckCircleIcon, ExclamationCircleIcon, XIcon } from '@heroicons/vue/solid'

export default defineComponent({
  components: {
    CheckCircleIcon,
    ExclamationCircleIcon,
    XIcon,
  },
  props: {
    id: String,
    title: String,
    message: String,
    pos: Number,
    style: String,
  },
  data() {
    return {
      show: true,
    }
  },
  computed: {
    titleClasses() {
      if (this.style === 'success') {
        return ['text-green-400']
      } else {
        return ['text-red-400']
      }
    },
  },
  mounted() {
    setTimeout(() => {
      this.close()
    }, 10000) // Auto hide after 10s
  },
  methods: {
    close() {
      this.emitter.on('notification-close', { id: this.id })
      this.show = false
    },
  },
})
</script>
