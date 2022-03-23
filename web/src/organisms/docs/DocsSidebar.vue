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
    name: 'Getting Started',
    links: [
      { route: 'v2DocsRoot', title: 'Getting Started' },
      { route: 'v2DocsProductIntro', title: 'Product Intro' },
      { route: 'v2DocsQuickStart', title: 'Quick start' },
      { route: 'v2DocsUsingSturdy', title: 'Using Sturdy' },
    ],
  },

  {
    name: 'Git',
    links: [
      { route: 'v2DocsHowSturdyInteractsWithGit', title: 'How Sturdy builds on Git' },

      { route: 'docsIntegrationsGit', title: 'Integrate with Git (GitLab, Azure, etc...)' },
      { route: 'v2DocsUsingSturdy', hash: 'import-a-repository', title: 'Sturdy for GitHub' },
    ],
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
      { route: 'v2DocsHowToSetupSturdyWithGitHub', title: 'Setup Sturdy with GitHub' },
      { route: 'v2DocsHowToEditCode', title: 'Edit code' },
      { route: 'v2DocsHowToSwitchBetweenTasks', title: 'Switch between tasks' },
    ],
  },

  {
    name: 'Enterprise',
    links: [
      { route: 'v2DocsPricing', title: 'Pricing' },
      { route: 'v2DocsSelfHosted', title: 'Self-hosted' },
    ],
  },

  {
    name: 'Features',
    links: [
      { route: 'featuresWorkspaceNavigation', title: 'Workspace Navigation' },
      { route: 'featuresIntegrations', title: 'Integrations' },
      { route: 'featuresLiveFeedback', title: 'Live Feedback' },
      // { route: 'featuresWorkflow', title: 'Workflow' },
      { route: 'featuresConflicts', title: 'Conflicts' },
      { route: 'featuresLargeFiles', title: 'Large Files' },
      { route: 'docsAccessControl', title: 'Access Control' },
      { route: 'docsSuggestions', title: 'Suggestions' },
      { route: 'docsCICD', title: 'CI / CD' },
      { route: 'docsIntegrationsGit', title: 'Integrate with Git (GitLab, Azure, etc ...)' },
    ],
  },

  {
    name: 'Learn',
    links: [
      { route: 'resourcesMigrateFromGitHub', title: 'Migrate from GitHub' },
      { route: 'resourcesHelpSturdyCommandLine', title: 'Sturdy for Servers (CLI)' },
      { route: 'resourcesApi', title: 'GraphQL API' },
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
