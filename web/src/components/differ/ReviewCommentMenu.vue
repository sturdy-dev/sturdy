<template>
  <Menu as="div" class="relative inline-block text-left">
    <div>
      <MenuButton
        class="bg-white rounded-full flex items-center text-gray-400 hover:text-gray-600 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-offset-gray-100 focus:ring-indigo-500"
      >
        <span class="sr-only">Open options</span>
        <DotsVerticalIcon class="h-5 w-5" aria-hidden="true" />
      </MenuButton>
    </div>

    <transition
      enter-active-class="transition ease-out duration-100"
      enter-from-class="transform opacity-0 scale-95"
      enter-to-class="transform opacity-100 scale-100"
      leave-active-class="transition ease-in duration-75"
      leave-from-class="transform opacity-100 scale-100"
      leave-to-class="transform opacity-0 scale-95"
    >
      <MenuItems
        class="origin-top-right absolute right-0 mt-2 w-56 rounded-md shadow-lg bg-white ring-1 ring-black ring-opacity-5 focus:outline-none z-10"
      >
        <div class="py-1">
          <MenuItem v-if="canEdit" v-slot="{ active }">
            <a
              :class="[
                active ? 'bg-gray-100 text-gray-900' : 'text-gray-700',
                'block px-4 py-2 text-sm cursor-pointer',
              ]"
              @click.stop="$emit('startEdit')"
            >
              Edit
            </a>
          </MenuItem>
          <MenuItem v-if="canEdit" v-slot="{ active }">
            <a
              :class="[
                active ? 'bg-gray-100 text-gray-900' : 'text-gray-700',
                'block px-4 py-2 text-sm cursor-pointer',
              ]"
              @click.stop="$emit('delete')"
            >
              Delete
            </a>
          </MenuItem>
          <MenuItem v-slot="{ active }">
            <a
              :class="[
                active ? 'bg-gray-100 text-gray-900' : 'text-gray-700',
                'block px-4 py-2 text-sm cursor-pointer',
              ]"
              @click.stop="onCopyLink"
            >
              Copy link
            </a>
          </MenuItem>
        </div>
      </MenuItems>
    </transition>
  </Menu>
</template>

<script>
import { Menu, MenuButton, MenuItem, MenuItems } from '@headlessui/vue'
import { DotsVerticalIcon } from '@heroicons/vue/solid'

export default {
  components: {
    Menu,
    MenuButton,
    MenuItem,
    MenuItems,
    DotsVerticalIcon,
  },
  props: {
    canEdit: {
      type: Boolean,
      default: false,
    },
    comment: {
      type: Object,
      required: true,
    },
  },
  emits: ['startEdit', 'delete'],
  methods: {
    onCopyLink() {
      const currentUrl = new URL(window.location.href)
      currentUrl.hash = `#${this.comment.id}`

      navigator.clipboard.writeText(currentUrl.href)

      this.emitter.emit('notification', {
        title: 'Copied comment link',
        message: 'The comment link has been copied to your clipboard.',
      })
    },
  },
}
</script>
