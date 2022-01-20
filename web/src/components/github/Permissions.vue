<template>
  <fieldset class="space-y-4">
    <legend class="sr-only">Notifications</legend>
    <div
      v-for="permission in permissions"
      :key="permission.title"
      class="relative flex items-start"
    >
      <div class="flex items-center h-5">
        <CheckIcon class="h-4 w-4 text-green-700" />
      </div>
      <div class="ml-3 text-sm space-y-1">
        <span class="font-medium text-gray-700">{{ permission.title }}</span> &mdash;
        <span>{{ permission.level }}</span>
        <p><span v-html="permission.reason"></span></p>
      </div>
    </div>
  </fieldset>
</template>
<script>
import CheckIcon from '@heroicons/vue/solid/CheckIcon'

const permissions = [
  {
    title: 'Contents',
    gitHubDescription: 'Repository contents, commits, branches, downloads, releases, and merges.',
    level: 'Read & write',
    reason:
      'Sturdy pulls and pushes code to allow you use Sturdy and GitHub at the same time. Code is pushed to branches prefixed with <code>sturdy-pr-</code>, and optionally to the default branch.',
  },
  {
    title: 'Metadata',
    level: 'Read-only',
    gitHubDescription: 'Search repositories, list collaborators, and access repository metadata.',
    reason:
      'Sturdy uses the name of the repository, and connects all collaborators with their Sturdy account.',
  },
  {
    title: 'Pull requests',
    level: 'Read & write',
    gitHubDescription:
      'Pull requests and related comments, assignees, labels, milestones, and merges.',
    reason:
      'Sturdy creates pull requests from Workspaces when GitHub is the source of truth of a codebase.',
  },
  {
    title: 'Workflows',
    level: 'Read & write',
    gitHubDescription: 'Update GitHub Action workflow files.',
    reason: 'Allows you to edit files under <code>.github/workflows</code> via Sturdy.',
  },
]

export default {
  components: { CheckIcon },
  setup() {
    return { permissions }
  },
}
</script>
