<template>
  <div class="relative flex items-start">
    <span class="relative inline-flex shadow-sm rounded-md z-10">
      <Button :first="true" :grouped="true" size="small" @click="$emit('add')">
        <span v-if="isAdded">Deselect file</span>
        <span v-else>Select file</span>
      </Button>
      <Button
        :last="true"
        :grouped="true"
        size="none"
        class="px-2"
        @click="toggleDropDown"
        @blur="fileDropDownBlur"
      >
        <span class="sr-only">Open options</span>
        <ChevronDownIcon class="h-5 w-5" />
      </Button>

      <transition
        enter-active-class="transition ease-out duration-100"
        enter-from-class="transform opacity-0 scale-95"
        enter-to-class="transform opacity-100 scale-100"
        leave-active-class="transition ease-in duration-75"
        leave-from-class="transform opacity-100 scale-100"
        leave-to-class="transform opacity-0 scale-95"
      >
        <div
          v-if="showDropdown"
          class="file-dropdown-overlay z-40 absolute right-0 mt-8 -mr-1 w-56 rounded-md shadow-lg bg-white ring-1 ring-black ring-opacity-5 focus:outline-none"
          role="menu"
          aria-orientation="vertical"
          aria-labelledby="option-menu"
        >
          <div class="py-1" role="none">
            <!--- If more options are added here, remember to adjust the min-height above to fit the menu on small files -->
            <!-- Only new files can be ignored -->
            <button
              :disabled="canIgnoreFile"
              class="block w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100 hover:text-gray-900 disabled:opacity-50"
              title="Ignore future changes to this file"
              role="menuitem"
              @click="$emit('ignore')"
            >
              Ignore file
            </button>

            <button
              class="block w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100 hover:text-gray-900 disabled:opacity-50"
              role="menuitem"
              @click="hide"
            >
              {{ isHidden ? 'Show' : 'Hide' }}
            </button>

            <button
              class="block w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100 hover:text-gray-900 disabled:opacity-50"
              role="menuitem"
              @click="$emit('undo')"
            >
              Undo the changes to this file
            </button>
          </div>
        </div>
      </transition>
    </span>
  </div>
</template>

<script lang="js">
import {ChevronDownIcon} from '@heroicons/vue/solid';
import {IsFocusChildOfElementWithClass} from "../../focus";
import Button from "../shared/Button.vue";

export default {
  name: "DifferAddButton",
  components: {ChevronDownIcon, Button},
  props: {
    isAdded: Boolean,
    isHidden: Boolean,
    canIgnoreFile: Boolean,
  },
  emits: ['add', 'hide', 'undo', 'ignore', 'showdropdown', 'hidedropdown', 'unhide'],
  data() {
    return {
      showDropdown: false,
    }
  },
  methods: {
    toggleDropDown() {
      this.showDropdown = !this.showDropdown;
      this.$emit('showdropdown')
      this.$emit('unhide')
    },
    hide() {
      this.showDropdown = false;
      this.$emit('hidedropdown')
      this.$emit('hide')
    },
    fileDropDownBlur(e) {
      // Lost blur to outside of dropdown
      if (!IsFocusChildOfElementWithClass(e, "file-dropdown-overlay")) {
        this.showDropdown = false;
        this.$emit('hidedropdown')
      }
    },
  }
}
</script>
