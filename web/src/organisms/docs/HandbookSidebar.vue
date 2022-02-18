<template>
  <ol class="space-y-2 text-sm">
    <template v-for="(group, idx) in groups" :key="idx">
      <li class="mt-2 font-medium tracking-tight font-semibold text-slate-800">{{ group.name }}</li>

      <br />
      <li v-for="(link, linkIdx) in group.links" :key="linkIdx" class="">
        <router-link
          :to="{ name: link.route }"
          class="border-l pl-4 py-2"
          :class="[
            isCurrent(link)
              ? 'text-amber-600 border-amber-600 font-semibold'
              : 'text-slate-700 hover:text-slate-900 hover:border-slate-400',
          ]"
        >
          {{ link.title }}
        </router-link>
      </li>

      <br />
    </template>
  </ol>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import { useRoute } from 'vue-router'

const groups = [
  {
    name: 'Handbook',
    links: [
      { route: 'handbookCodeOfConduct', title: 'Code of Conduct' },
      { route: 'handbookReleases', title: 'Releasing' },
    ],
  },
]

export default defineComponent({
  setup() {
    let route = useRoute()

    return {
      groups,
      route,
    }
  },
  methods: {
    isCurrent(link) {
      return link.route === this.route.name
    },
  },
})
</script>
