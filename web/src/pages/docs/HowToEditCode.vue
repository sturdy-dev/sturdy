<template>
  <PublicLeftSidebar>
    <template #sidebar>
      <DocsSidebar />
    </template>
    <template #default>
      <div class="prose p-4 max-w-[800px]">
        <h1 id="how-sturdy-augments-git">How to edit code</h1>
        <p>
          With Sturdy, you access and edit your code with your choice of local IDE or text editor – no different from what you are used to.
        </p>

        <p>After you have created, imported or joined an existing codebase:</p>

        <ol>
          <li>Navigate to the codebase by selecting it in the sidebar on the left.</li>
          <li>Connect a directory where you want to have the code by clicking the "Connect directory" button.</li>
          <li>Select a directory where you want to have the code. It is a good idea to create a new folder.</li>
          <ol>
            <li>The folder you select must be empty</li>
            <li>The folder that you select cannot be inside a git repository</li>
            <li>A current limitation - once setup, the folder cannot be moved</li>
          </ol>
          <li>You can use the folder for accessing and editing code with any IDE or text editor you are used to (Tip: Use
            the Open in VSCode button for quick access)</li>
          <li>Observe that a new workspace was created named "Your Name's Workspace"</li>
        </ol>

        <p>
          Next steps:
          <router-link :to="{ name: 'v2DocsHotToShipSoftwareToProduction' }">How do I ship software to production</router-link>
        </p>

        <DocsInfoBox>Any time you create code changes, they "happen" in a Sturdy workspace. Workspaces provide isolation, just
          like Git branches – it’s where you work. The main difference between Sturdy and others is the highly automated
          workflow that it enables.</DocsInfoBox>

        <p>
          As you write the code, you can see your changes in Sturdy, undo them,
          <router-link :to="{ name: 'v2DocsHowToCollaborateWithOthers' }">ask for feedback</router-link>
          or
          <router-link :to="{ name: 'v2DocsHotToShipSoftwareToProduction' }">ship them to production</router-link>
          – all of this is just a click away.
          You can read more about how Sturdy augments Git on
          <router-link :to="{ name: 'v2DocsHowSturdyAugmentsGit' }">this page</router-link>
        </p>
        <p>
          <em>Coming soon:</em> It will be possible to toggle an auto-sync option for a workspace. With that enabled, you will always get the latest changes by others.
        </p>
      </div>
    </template>
  </PublicLeftSidebar>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import PublicLeftSidebar from '../../layouts/PublicLeftSidebar.vue'
import { useHead } from '@vueuse/head'
import DocsSidebar from '../../organisms/docs/DocsSidebar.vue'
import DocsInfoBox from '../../molecules/DocsInfoBox.vue'

export default defineComponent({
  components: { DocsSidebar, PublicLeftSidebar, DocsInfoBox },
  setup() {
    // TODO: Remove when we're launching!
    useHead({
      meta: [
        {
          name: 'robots',
          content: 'noindex',
        },
      ],
    })
  },
})
</script>
