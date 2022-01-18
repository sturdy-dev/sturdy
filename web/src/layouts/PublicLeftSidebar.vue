<template>
  <PublicOpenSource>
    <div class="flex mt-8 px-4 flex-col md:flex-row">
      <div class="w-full flex-1 md:min-w-[240px]">
        <div class="md:sticky md:top-4">
          <slot name="sidebar"></slot>
        </div>
      </div>

      <div
        ref="content"
        class="w-full 2xl:max-w-[800px] xl:max-w-[650px] max-w-full pb-14 relative md:pl-16 xl:px-16 2xl:px-32 box-content overflow-auto"
      >
        <slot name="default"></slot>
      </div>

      <div class="w-full flex-1 pl-4 self-stretch hidden xl:block">
        <div class="sticky top-4">
          <span class="leading-8 text-gray-500 font-medium">On this page</span>
          <ol class="leading-loose mt-4">
            <li v-for="(toc, idx) in tableOfContents" :key="idx">
              <a
                :href="'#' + toc.id"
                class="text-gray-600 hover:text-gray-800 font-medium"
                :class="toc.classes"
                >
                {{ toc.title }}
              </a>
            </li>
          </ol>
        </div>
      </div>
    </div>
  </PublicOpenSource>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import PublicOpenSource from "./PublicOpenSource.vue";

export default defineComponent({
  components: {
    PublicOpenSource,
  },
  setup() {
    return {
      user: null,
    }
  },
  data() {
    return {
      tableOfContents: [],
    }
  },
  mounted() {
    this.$nextTick(function () {
      this.tableOfContents = this.buildTable()
    })
  },
  methods: {
    buildTable() {
      let res = []

      const content = this.$refs.content
      if (!content) {
        return res
      }

      const headings = content.querySelectorAll('h1, h2')

      for (const heading of headings) {
        let classes = []
        if (heading.nodeName === 'H2') {
          classes.push('pl-4')
        }

        let title
        if (heading.title) {
          title = heading.title
        } else if (heading.innerText) {
          title = heading.innerText
        }
        if (!title) {
          continue
        }

        res.push({
          title,
          classes,
          id: heading.id,
        })
      }

      return res
    },
  },
})
</script>
