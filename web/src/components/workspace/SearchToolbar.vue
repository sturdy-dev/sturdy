<template>
  <transition
    enter-active-class="transition ease-out duration-100"
    enter-from-class="transform opacity-0 -translate-y-1"
    enter-to-class="transform opacity-100 translate-y-none"
    leave-active-class="transition ease-in duration-75"
    leave-from-class="transform opacity-100 translate-y-none"
    leave-to-class="transform opacity-0 -translate-y-1"
  >
    <div
      v-if="showSearch"
      class="bg-yellow-400 text-white shadow-lg flex flex-row items-center border-b border-yellow-700 px-5 py-1 gap-2"
    >
      <div class="text-gray-800 ml-3 flex-1">
        <label for="search" class="sr-only">Search</label>
        <div class="relative rounded-md shadow-sm">
          <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
            <SearchIcon class="h-5 w-5 text-gray-800" aria-hidden="true" />
          </div>
          <input
            id="search"
            ref="search"
            v-model="searchQuery"
            type="search"
            name="search"
            class="focus:ring-yellow-500 focus:border-yellow-500 block w-full pl-10 sm:text-sm border-gray-300 rounded-md"
            placeholder="Search"
            @keyup.enter="searchKeyUpEnter"
            @keyup.esc="searchStop"
            @keyup="emitSearch"
          />
        </div>
      </div>
      <div class="flex-1 select-none flex gap-2 items-center">
        <div v-if="matchesCount > 0" class="flex gap-1">
          <ArrowSmUpIcon
            class="h-8 w-8 text-white hover:bg-yellow-600 cursor-pointer p-1 rounded-md"
            @click.stop.prevent="searchPrev"
          />
          <ArrowSmDownIcon
            class="h-8 w-8 text-white hover:bg-yellow-600 cursor-pointer p-1 rounded-md"
            @click.stop.prevent="searchNext"
          />
        </div>

        <span v-if="matchesCount > 0">{{ searchCurrentIdx + 1 }} / {{ matchesCount }}</span>
        <span v-else-if="searchQuery">No matches</span>
      </div>
      <XIcon
        class="h-8 w-8 text-white hover:bg-yellow-600 cursor pointer p-1 rounded-md"
        @click="() => searchStop()"
      />
    </div>
  </transition>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import { ArrowSmDownIcon, ArrowSmUpIcon, SearchIcon, XIcon } from '@heroicons/vue/solid'

export default defineComponent({
  components: {
    SearchIcon,
    ArrowSmDownIcon,
    ArrowSmUpIcon,
    XIcon,
  },
  data() {
    const ipc = window.ipc
    return {
      ipc,
      searchQuery: '',
      showSearch: false,
      searchCurrentIdx: -1,
      matchesCount: 0,
      searchCurrentSelectedId: null,
    }
  },
  computed: {
    isApp() {
      return !!this.ipc
    },
  },
  watch: {
    searchQuery: function () {
      this.searchCurrentSelectedId = null
      this.emitSearch()
    },
  },
  mounted() {
    window.addEventListener('keydown', this.globalKeyDown)
    this.emitter.on('search-result', this.onSearchResult)
  },
  unmounted() {
    window.removeEventListener('keydown', this.globalKeyDown)
    this.emitter.off('search-result', this.onSearchResult)
  },
  methods: {
    globalKeyDown(event) {
      // Available as cmd+k on the web
      let keys = [
        75, // K
      ]

      // Also available as cmd+f in the app
      if (this.isApp) {
        keys.push(70) // F
      }

      if (keys.indexOf(event.keyCode) > -1 && (event.ctrlKey || event.metaKey)) {
        this.showSearch = true
        event.stopPropagation()
        event.preventDefault()
        this.$nextTick(() => {
          this.$refs.search.focus()
        })
        return false
      }

      return true
    },
    searchKeyUpEnter(event) {
      if (event.shiftKey) {
        this.searchPrev()
      } else {
        this.searchNext()
      }
    },
    searchNext() {
      this.searchCurrentIdx++
      this.searchScrollTo()
    },
    searchPrev() {
      this.searchCurrentIdx--
      this.searchScrollTo()
    },
    searchScrollTo() {
      let allMatches = document.getElementsByClassName('sturdy-searchmatch')
      if (!allMatches) {
        return
      }

      if (this.searchCurrentIdx >= allMatches.length) {
        this.searchCurrentIdx = 0
      }
      if (this.searchCurrentIdx < 0) {
        this.searchCurrentIdx = allMatches.length - 1
      }

      if (this.searchCurrentIdx < allMatches.length) {
        let el = allMatches[this.searchCurrentIdx]
        el.scrollIntoView({ block: 'center', inline: 'start' })
        this.searchCurrentSelectedId = el.id
      }
    },
    searchStop() {
      this.showSearch = false
      this.searchQuery = ''
      this.searchCurrentSelectedId = null
      this.emitSearch()
    },
    emitSearch() {
      this.emitter.emit('search', {
        searchQuery: this.searchQuery,
        searchCurrentSelectedId: this.searchCurrentSelectedId,
      })
    },
    onSearchResult(event) {
      this.matchesCount = event.matchesCount
    },
  },
})
</script>
