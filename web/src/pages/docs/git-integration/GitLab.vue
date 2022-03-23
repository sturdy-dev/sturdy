<template>
  <DocumentationWithTableOfContents>
    <template #sidebar>
      <DocsSidebar />
    </template>
    <template #default>
      <div class="prose p-4 max-w-[800px]">
        <h1 id="working-in-the-open">Connect with Git: Sturdy and GitLab</h1>

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
          GitLab, you can find it under the <strong>Clone</strong> button on a repository (see
          screenshot below).
        </p>

        <img
          src="./gitlab-overview.png"
          height="554"
          width="1990"
          class="w-82 border-2 rounded-md"
        />

        <img
          src="./gitlab-url.png"
          height="558"
          width="708"
          class="w-[30rem] border-2 rounded-md"
        />

        <p>
          Get your GitLab credentials by generating a new Personal Access Token (here:

          <a target="_blank" href="https://gitlab.com/-/profile/personal_access_tokens"
            >https://gitlab.com/-/profile/personal_access_tokens</a
          >), and give the token the <code>write_repository</code> permission. Use your GitLab
          username as the username, and your Personal Access Token as the password of the
          integration.
        </p>

        <p>
          Fill in all of the other fields as instructed, when you're done. The configuration should
          look something like this:
        </p>

        <img src="./sturdy-gitlab.png" height="1920" width="1616" class="border-2 rounded-md" />
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
