<template>
  <PublicLeftSidebar>
    <template #sidebar>
      <DocsSidebar />
    </template>
    <template #default>
      <div class="prose p-4 max-w-[800px]">
        <h1 id="how-to-switch-between-tasks">How to switch between tasks</h1>
        <p>
          Sturdy allows you to easily and quickly switch between tasks by using multiple workspaces.
          Workspaces are a unit of isolation in Sturdy, which means that code changes in one
          workspace are independent from those in another. It is quick to create and dispose of
          workspaces. This guide will go through how you can use multiple workspaces to seamlessly
          switch between tasks. Make sure you are familiar with
          <router-link :to="{ name: 'v2DocsHowToEditCode' }">editing code</router-link>
          in Sturdy before continuing.
        </p>

        <h2 id="switching-workspaces">Multitasking with workspaces</h2>
        <p>
          You can use workspaces to jump between multiple work in progress items. After you have
          created, imported or joined an existing codebase:
        </p>

        <ol>
          <li>
            Create a new workspace by clicking the "plus" icon next to the codebase name in the left
            sidebar
          </li>
          <li>Navigate to the new workspace named "<em>Username</em>'s Workspace"</li>
          <li>Rename the newly created workspace to "Feature Alpha" using the "Edit" button</li>
          <li>
            Click the "Connect <em>folder</em> for editing" button at the top of the workspace
          </li>
          <li>The chosen <em>folder</em> is now connected to the "Feature Alpha" workspace</li>
          <li>Using your IDE of choice, create some changes to the code</li>
          <li>The diffs of your changes appear in the "Feature Alpha" workspace</li>
          <li>
            Create a new workspace by clicking the "plus" icon next to the codebase name in the left
            sidebar
          </li>
          <li>Navigate to the new workspace named "<em>Username</em>'s Workspace"</li>
          <li>Rename the newly created workspace to "Feature Beta" using the "Edit" button</li>
          <li>
            Click the "Connect <em>folder</em> for editing" button at the top of the workspace
          </li>
          <li>The chosen <em>folder</em> is now connected to the "Feature Beta" workspace</li>
          <ol>
            <li>
              Observe that in the left sidebar, the computer icon has moved from "Feature Alpha" to
              "Feature Beta"
            </li>
            <li>
              Hover over the computer icon next to "Feature Beta" to reveal additional information
            </li>
          </ol>
          <li>
            In your IDE of choice, observe that the code changes that you made earlier are
            <em>not</em> present
          </li>
          <ol>
            <li>
              This is because the changes were made in workspace "Feature Alpha" and your directory
              is now connected to workspace "Feature Beta"
            </li>
          </ol>
          <li>Using your IDE:</li>
          <ol>
            <li>
              Create some new file changes, different from the ones you did earlier for the first
              workspace
            </li>
          </ol>
          <li>The diffs of your changes appear in the "Feature Beta" workspace</li>
          <li>Using the sidebar on the left, navigate back to the "Feature Alpha" workspace</li>
          <ol>
            <li>In Sturdy you can see the diffs that you made earlier</li>
          </ol>
          <li>
            Click the "Connect <em>folder</em> for editing" button at the top of the workspace page
            to connect your folder to "Feature Alpha"
          </li>
          <ol>
            <li>
              Observe that in the left sidebar, the computer icon has moved from "Feature Beta" to
              "Feature Alpha"
            </li>
            <li>
              In your IDE, observe that your code now has the changes from "Feature Alpha" and
              <em>no longer</em> has the "Feature Beta" changes
            </li>
          </ol>
        </ol>

        <p>
          Workspaces in Sturdy are lightweight and can be created at any time. This allows you to
          have multiple experiments at the same time, or the ability to quickly fix a bug without
          disrupting your primary task.
        </p>

        <DocsInfoBox
          >In Sturdy, there is no need to 'save' or 'stash' changes when switching between
          workspaces.</DocsInfoBox
        >

        <DocsInfoBox class="mt-8"
          >You can connect the directory on your computer to workspaces created by others on your
          team. Learn more about
          <router-link :to="{ name: 'v2DocsHowToCollaborateWithOthers' }"
            >collaborating with others</router-link
          >.
        </DocsInfoBox>

        <h2 id="workspaces-based-on-other-workspaces">Workspaces bases on other workspaces</h2>
        <p>
          It can be useful to create workspaces based on other workspaces or parts of. Consider the
          following two examples:
        </p>
        <ol>
          <li>
            You have written some code on which you wish to get feedback, and deploy independently.
            You wish to continue building on this code, as a separate unit of work
          </li>
          <li>
            While you are building a feature you spot a small bug or typo in some unrelated code,
            which you decide to fix. You wish to extract the bug/typo fix and have it deployed
            separately
          </li>
        </ol>
        <p>
          Sturdy allows you to create new workspaces based on the work in an existing workspace.
          Within an existing workspace which <em>already has</em> some work in progress code
          changes:
        </p>
        <ol>
          <li>
            Using the checkbox on the left of changes, select one or multiple diffs that you wish to
            copy to a new workspace
          </li>
          <ol>
            <li>
              <em>Optional</em>: Use the "Select All" checkbox at the top in order to quickly select
              all diffs
            </li>
          </ol>
          <li>
            Observe the banner in blue at the top of the app indicating "<em>n</em> change(s)
            selected"
          </li>
          <li>
            Click the "Copy to New Workspace" button to create a new workspace with the selected
            changes
          </li>
          <li>
            Observe that a new workspace named "Fork of <em>name of the original</em>" has been
            created
          </li>
          <li>
            Using the sidebar on the left, navigate to the workspace "Fork of
            <em>name of the original</em>"
          </li>
          <li>Observe that the new workspace has only the changes that you selected</li>
          <li>Rename the workspace to reflect the changes included, for example "Quick bugfix"</li>
        </ol>

        <p>
          Creating workspaces based on other workspaces is a powerful tool that allows you to
          quickly and easily ship small incremental changes.
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
    useHead({
      meta: [
        // TODO: Remove when we're launching!
        {
          name: 'robots',
          content: 'noindex',
        },
        {
          name: 'description',
          content: 'Learn how to quickly jump between tasks in Sturdy',
        },
        {
          name: 'keywords',
          content: 'study learn documentation switch change task workspace coding branches',
        },
      ],
      title: 'How to switch between tasks | Sturdy',
    })
  },
})
</script>
