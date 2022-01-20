<template>
  <ol class="space-y-2">
    <template v-for="(group, idx) in groups" :key="idx">
      <li class="font-medium text-gray-400">{{ group.name }}</li>

      <li v-for="(link, linkIdx) in group.links" :key="linkIdx" class="ml-4">
        <router-link
          :to="{ name: link.route }"
          class="hover:text-gray-800 font-medium"
          :class="[isCurrent(link) ? 'text-gray-800' : 'text-gray-500']"
        >
          {{ link.title }}
        </router-link>
      </li>
    </template>
  </ol>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import { useRoute } from 'vue-router'

const groups = [
  {
    name: 'Getting Started',
    links: [
      { route: 'v2DocsProductIntro', title: 'Getting Started' },
      { route: 'v2DocsQuickStart', title: 'Quick start' },
      { route: 'v2DocsUsingSturdy', title: 'Using Sturdy' },
    ],
  },

  {
    name: 'Git',
    links: [{ route: 'v2DocsHowSturdyAugmentsGit', title: 'Augmenting Git' }],
  },

  {
    name: 'Vision',
    links: [{ route: 'v2DocsWorkingInTheOpen', title: 'Working in the open' }],
  },

  {
    name: 'How to',
    links: [
      { route: 'v2DocsHotToShipSoftwareToProduction', title: 'Ship to production' },
      { route: 'v2DocsHowToCollaborateWithOthers', title: 'Collaborate with others' },
      { route: 'v2DocsHowToSetupSturdyOnGitHub', title: 'Setup Sturdy on GitHub' },
      { route: 'v2DocsHowToEditCode', title: 'Edit code' },
      { route: 'v2DocsHowToSwitchBetweenTasks', title: 'Switch between tasks' },
    ],
  },
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
