<template>
  <DocumentationWithTableOfContents>
    <div class="prose p-4 max-w-[800px]">
      <h1>Sturdy Public Roadmap</h1>
      <sub>Updated May 2022</sub>
      <p>
        What's this? This is the Sturdy Public Roadmap, it's describing what we're working on right
        now, and what we're planning to work on in the future.
      </p>
      <p>
        New ideas, feedback, and discussions are welcome, and best happens on our
        <a href="https://discord.gg/fQcH9QAVpX">Discord</a>. :-)
      </p>
    </div>

    <div class="flex gap-4 flex-col md:flex-row">
      <div v-for="column in columns" :key="column.title" class="flex-1 flex-shrink-0">
        <h2 class="text-center uppercase font-bold text-gray-700 mb-2">{{ column.title }}</h2>
        <div class="flex flex-col gap-4">
          <div v-for="item in column.items" :key="item.title" class="bg-yellow-100 p-2 rounded-md">
            <h4 class="font-bold text-gray-700">{{ item.title }}</h4>
            <p class="text-sm text-gray-600">{{ item.description }}</p>
          </div>
        </div>
      </div>
    </div>
  </DocumentationWithTableOfContents>
</template>

<script lang="ts" setup>
import DocumentationWithTableOfContents from '../layouts/DocumentationWithTableOfContents.vue'

const later = {
  title: 'Later',
  items: [
    {
      title: 'Editor Integrations',
      description: 'Integrate with VSCode and JetBrains based IDEs for a deeper integration',
    },
    {
      title: 'Ticket integrations',
      description: 'Integrate with task management systems like JIRA and Trello',
    },
    {
      title: 'Jump-to-next-task',
      description:
        'Suggest drafts to review, and new drafts to start working (based on integration with ticketing systems)',
    },
    {
      title: 'CI/CD performance metrics',
      description: 'Track performance of CI/CD pipelines, and suggest how to make pipelines faster',
    },
    {
      title: 'Instant Integration (II)',
      description: 'Continiously run tests on spare CI/CD resources',
    },
    {
      title: 'Quantum-repos',
      description: 'Monorepo or polyrepo? Why not both?',
    },
  ],
}

const soon = {
  title: 'Soon',

  items: [
    {
      title: 'Webhooks',
      description: 'Register custom webhook receivers to connect Sturdy to any system in real-time',
    },
    {
      title: 'Stacked Drafts',
      description: 'Base one draft on top of another, and require them to be merged in-order',
    },
    {
      title: 'Enforce CI/CD',
      description: 'Option to enforce tests to pass before merging',
    },
    {
      title: 'Suggested reviewers',
      description: 'Suggest reviewers based on the current state of the code',
    },
    {
      title: 'Permission management',
      description: 'Manage permissions for users and teams in organisations and codebases',
    },
    {
      title: 'Service Accounts',
      description:
        'Bot accounts with limited scopes and access for easier scripting of your Sturdy resources',
    },
  ],
}

const shipped = {
  title: 'Shipped',
  items: [
    {
      title: 'Push/Pull to any git host',
      description: 'Integrate with GitLab, BitBucket, Azure DevOps, etc',
    },
    {
      title: 'Undo and redo',
      description: 'Undo and redo changes in a draft',
    },
  ],
}

const columns = [later, soon, shipped]
</script>
