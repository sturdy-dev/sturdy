<template>
  <DocumentationWithTableOfContents>
    <template #sidebar>
      <DocsSidebar />
    </template>
    <template #default>
      <div class="prose p-4 max-w-[800px]">
        <h1 id="working-in-the-open">Connect with Git: Sturdy and Azure DevOps</h1>

        <p>
          <router-link :to="{ name: 'docsIntegrationsGit' }">
            &laquo; Back to "Integrate with Git"
          </router-link>
        </p>

        <p>
          Start by opening Sturdy, and go to the codebase that you want to setup the integration
          with. If you don't have a codebase yet, create the codebase on Sturdy first.
        </p>

        <p>On Sturdy: Go to the codebase > Settings > Integrations > Git > Add</p>

        <p>
          Under <strong>Remote URL (for push and pull)</strong>, enter a HTTP URL for push/pull. On
          Azure DevOps, you can find it under the <strong>Clone</strong> button on repository (see
          screenshot below).
        </p>

        <img src="./azure-overview.png" height="462" width="2298" class="border-2 rounded-md" />

        <img src="./azure-url.png" height="370" width="952" class="w-[30rem] border-2 rounded-md" />

        <p>
          Get your Azure credentials by clicking <strong>Generate git credentials</strong> on Azure.
          Enter the username and password on Sturdy.
        </p>

        <p>
          Fill in all of the other fields as instructed, when you're done. The configuration should
          look something like this:
        </p>
        <img src="./sturdy-azure.png" height="1944" width="1912" class="border-2 rounded-md" />
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
      content: 'Connect Sturdy and Azure DevOps',
    },
  ],
  title: 'Setup Guide: Sturdy and Azure DevOps',
})
</script>
