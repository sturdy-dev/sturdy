<template>
  <Public>
    <div class="flex mt-8 px-4 flex-col md:flex-row">
      <div class="w-full flex-1 md:min-w-[240px]">
        <div class="md:sticky md:top-4">
          <slot name="sidebar"></slot>
        </div>
      </div>

      <div class=" w-full 2xl:max-w-[800px] xl:max-w-[650px] max-w-full pb-14 relative md:pl-16 xl:px-16 2xl:px-32 box-content overflow-auto">
        <slot name="default"></slot>
      </div>

      <div class="w-full flex-1 pl-4 self-stretch hidden xl:block">
        <div class="sticky top-4">
          <span class="leading-8 text-gray-500 font-medium">On this page</span>
          <ol class="leading-loose">
            <li><a href="#" class="text-gray-600 hover:text-gray-800 font-medium">Foo</a></li>
            <li><a href="#" class="text-gray-600 hover:text-gray-800 font-medium pl-4">Bar</a></li>
            <li><a href="#" class="text-gray-600 hover:text-gray-800 font-medium pl-4">Bar baz foo baz</a></li>
            <li><a href="#" class="text-gray-600 hover:text-gray-800 font-medium">Foo</a></li>
          </ol>
        </div>
      </div>
    </div>
  </Public>
</template>

<script lang="ts">
import {defineComponent} from 'vue'
import Public from "./Public.vue";

export default defineComponent({
  components: {
    Public,
  },
  setup() {
    return {
      user: null,
    }
  },
})
</script>
