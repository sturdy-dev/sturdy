<template>
  <DocumentationWithTableOfContents>
    <template #sidebar>
      <DocsSidebar />
    </template>
    <template #default>
      <div class="prose p-4 max-w-[800px]">
        <h1 id="working-in-the-open">Integrate with Git</h1>

        <p>
          With the <strong>Git</strong> integration in Sturdy, you can connect Sturdy to any git
          hosting provider, such as GitLab, BitBucket, or Azure Dev Ops using the GIT protocol.
        </p>

        <ul>
          <li>Push a Sturdy Draft as a branch</li>
          <li>Sync (pull/push) the Sturdy changelog to a tracked branch (eg: <code>main</code>)</li>
          <li>Merge changes on Sturdy, and push to the tracked branch</li>
        </ul>

        <h2>Setup</h2>

        <p>
          On Sturdy, go the integrations page, and setup a new <strong>Git</strong> integration. The
          page will show you step-by-step instruction for how to setup the integration.
        </p>

        <img src="./sturdy-integrations.png" height="584" width="1618" />

        <p>Here are guides for how to setup Sturdy on some of the most popular Git providers.</p>

        <ul>
          <li>
            <router-link :to="{ name: 'docsIntegrationsGitAzureDevops' }">
              Setup guide: Azure DevOps
            </router-link>
          </li>
          <li>
            <router-link :to="{ name: 'docsIntegrationsGitGitLab' }">
              Setup guide: GitLab
            </router-link>
          </li>
        </ul>
      </div>
    </template>
  </DocumentationWithTableOfContents>
</template>

<script lang="ts" setup>
import DocumentationWithTableOfContents from '../../../layouts/DocumentationWithTableOfContents.vue'
import { useHead } from '@vueuse/head'
import DocsSidebar from '../../../organisms/docs/DocsSidebar.vue'

useHead({
  meta: [
    {
      name: 'description',
      content: 'Connect with Git',
    },
  ],
  title: 'Connect with Git | Sturdy',
})
</script>
