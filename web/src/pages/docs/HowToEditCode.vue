<template>
  <DocumentationWithTableOfContents>
    <template #sidebar>
      <DocsSidebar/>
    </template>
    <template #default>
      <div class="prose p-4 max-w-[800px]">
        <h1 id="how-to-edit-code">How to edit code</h1>
        <p>
          With Sturdy, you access and edit your code with any local IDE or text editor &mdash; no
          different from what you are already used to.
        </p>

        <p>After you have created, imported or joined an existing codebase:</p>

        <ol>
          <li>Navigate to the codebase by selecting it in the sidebar on the left.</li>
          <li>
            Connect a directory where you want to have the code by clicking the "Connect directory"
            button.
          </li>
          <li>
            Select a directory where you want to have the code in. Sturdy will create a new
            directory inside the directory that you select with the same name as the codebase.
          </li>
          <li>
            You can use the folder for accessing and editing code with any IDE or text editor
            <SubtleText>
              (tip: Use the "Open in Visual Studio Code"-button for quick access)
            </SubtleText>
          </li>
          <li>Observe that a new draft change named "Untitled Draft" was created</li>
        </ol>

        <p>
          Next steps:
          <router-link :to="{ name: 'v2DocsHotToShipSoftwareToProduction' }">
            How do I ship software to production
          </router-link>
        </p>

        <DocsInfoBox>
          Any time you create code changes, they "happen" in a Sturdy draft change. Drafts provide
          isolation, just like Git branches &mdash; it’s where you work. The main difference between
          Sturdy and others is the highly automated workflow that it enables.
        </DocsInfoBox>

        <p>
          As you write the code, you can see your changes in Sturdy, undo them,
          <router-link :to="{ name: 'v2DocsHowToCollaborateWithOthers' }"
          >ask for feedback
          </router-link>
          or
          <router-link :to="{ name: 'v2DocsHotToShipSoftwareToProduction' }"
          >ship them to production
          </router-link>
          &mdash; all of this is just a click away. You can read more about how Sturdy interacts
          with Git on
          <router-link :to="{ name: 'v2DocsHowSturdyInteractsWithGit' }">this page</router-link>
          .
        </p>
        <p>
          <em>Coming soon:</em> It will be possible to toggle an auto-sync option on a draft change.
          With that enabled, you will always get the latest changes by others.
        </p>
      </div>
    </template>
  </DocumentationWithTableOfContents>
</template>

<script lang="ts" setup>
import DocumentationWithTableOfContents from '../../layouts/DocumentationWithTableOfContents.vue'
import {useHead} from '@vueuse/head'
import DocsSidebar from '../../organisms/docs/DocsSidebar.vue'
import DocsInfoBox from '../../molecules/DocsInfoBox.vue'
import SubtleText from '../../molecules/docs/SubtleText.vue'

useHead({
  meta: [
    {
      name: 'description',
      content: 'Learn how to edit code with Sturdy',
    },
    {
      name: 'keywords',
      content: 'study learn documentation coding code editing',
    },
  ],
  title: 'How to edit code | Sturdy',
})
</script>
