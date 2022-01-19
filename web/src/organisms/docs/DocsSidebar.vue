<template>
  <ol class="space-y-1">
    <li v-for="(link, idx) in links" :key="idx">
      <router-link
        :to="{ name: link.name }"
        class="hover:text-gray-800 font-medium"
        :class="[isCurrent(link) ? 'text-gray-800' : 'text-gray-500']"
      >
        {{ link.title }}
      </router-link>
    </li>
  </ol>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import { useRoute } from 'vue-router'

const links = [
  { name: 'v2DocsProductIntro', title: 'Getting Started' },
  { name: 'v2DocsHowSturdyAugmentsGit', title: 'Augmenting Git' },
  { name: 'v2DocsWorkingInTheOpen', title: 'Working in the open' },
  { name: 'v2DocsHotToShipSoftwareToProduction', title: 'How to: ship to production' },
  { name: 'v2DocsHowToCollaborateWithOthers', title: 'How to: collaborate with others' },
  { name: 'v2DocsHowToImportCodeFromGitHub', title: 'How to: import code from GitHub'}
  //  { name: 'index', title: 'Self-host' },
  //  { name: 'index', title: 'Cloud' },
  //  { name: 'index', title: 'User Guides' },
  //  { name: 'index', title: 'API' },
  //  { name: 'index', title: 'API' },
  //  { name: 'index', title: 'API' },
  //  { name: 'index', title: 'API' },
]

export default defineComponent({
  setup() {
    let route = useRoute()

    return {
      links,
      route,
    }
  },
  methods: {
    isCurrent(link) {
      return link.name === this.route.name
    },
  },
})
</script>
