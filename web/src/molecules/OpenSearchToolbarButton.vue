<template v-if="show">
  <Tooltip y-direction="down" x-direction="right">
    <template #tooltip>Open search toolbar</template>
    <button
        v-if="show"
      class="h-7 w-7 p-1 flex disabled:opacity-30 cursor-auto text-gray-500 disabled:bg-transparent disabled:text-gray-500 hover:bg-gray-200 hover:text-gray-600 rounded-md items-center justify-center"
      @click="open"
    >
      <SearchIcon class="w-5" />
    </button>
  </Tooltip>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import { SearchIcon } from '@heroicons/vue/solid'
import Tooltip from '../components/shared/Tooltip.vue'

export default defineComponent({
  components: { Tooltip, SearchIcon },
  data(){
    return{
      show: false
    }
  },
  methods: {
    open() {
      this.emitter.emit('open-search-toolbar')
    },
    searchToolbarButtonVisible(visible: boolean) {
      this.show = visible
    },
  },
  mounted() {
    this.emitter.on('search-toolbar-button-visible', this.searchToolbarButtonVisible)
  },
  unmounted() {
    this.emitter.on('search-toolbar-button-visible', this.searchToolbarButtonVisible)
  }
})
</script>
