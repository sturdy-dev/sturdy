<template>
  <Documentation>
    <div class="flex mt-8 px-4 flex-col md:flex-row text-sm tracking-tight">
      <div class="w-full flex-1 md:min-w-[240px]">
        <div class="hidden md:block md:sticky md:top-4">
          <slot name="sidebar"></slot>
        </div>
      </div>

      <div
        ref="content"
        class="docs w-full 2xl:max-w-[800px] xl:max-w-[650px] max-w-full pb-14 relative md:pl-16 xl:px-16 2xl:px-32 box-content overflow-auto"
      >
        <!-- Small screen table of contents -->
        <div class="w-full flex-1 pl-4 self-stretch block xl:hidden mb-10">
          <div class="top-4">
            <!-- Small Screen Link all documentation pages -->
            <span class="leading-8 text-slate-800 font-semibold">On this page</span>
            <ol class="leading-loose mt-4">
              <li v-for="(toc, idx) in tableOfContents" :key="idx" class="group">
                <a
                  :href="'#' + toc.id"
                  class="text-slate-700 hover:text-slate-900 space-x-2"
                  :class="toc.classes"
                >
                  <span>{{ toc.title }}</span>
                  <span class="text-gray-400 group-hover:text-gray-800">⤵️</span>
                </a>
              </li>
            </ol>

            <div class="leading-8 text-slate-800 mt-4 md:hidden">
              <router-link :to="{ name: 'v2DocsIndex' }">
                <span class="font-mono text-red-800">GOTO</span>: Sturdy Documentation Overview
              </router-link>
            </div>
          </div>
        </div>

        <slot name="default"></slot>
      </div>

      <!-- Large screen table of contents -->
      <div class="w-full flex-1 pl-4 self-stretch hidden xl:block">
        <div class="sticky top-4">
          <span class="leading-8 text-slate-800 font-semibold">On this page</span>
          <ol class="leading-loose mt-4">
            <li v-for="(toc, idx) in tableOfContents" :key="idx">
              <a
                :href="'#' + toc.id"
                class="text-slate-700 hover:text-slate-900 py-2"
                :class="toc.classes"
              >
                {{ toc.title }}
              </a>
            </li>
          </ol>
        </div>
      </div>
    </div>
  </Documentation>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import Documentation from './Documentation.vue'

type Item = {
  id: string
  title: string
  classes: Array<string>
}

export default defineComponent({
  components: {
    Documentation,
  },
  data() {
    return {
      tableOfContents: Array<Item>(),
    }
  },
  mounted() {
    this.$nextTick(function () {
      this.tableOfContents = this.buildTable()
    })
  },
  methods: {
    buildTable(): Array<Item> {
      let res = Array<Item>()

      const content = this.$refs.content
      if (!content) {
        return res
      }

      const el = content as HTMLElement
      const headings = el.querySelectorAll<HTMLElement>('h1, h2, h3')

      for (const heading of headings) {
        let classes = []
        if (heading.nodeName === 'H2') {
          classes.push('pl-4')
        }
        if (heading.nodeName === 'H3') {
          classes.push('pl-8')
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
