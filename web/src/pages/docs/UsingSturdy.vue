<template>
  <PublicLeftSidebar>
    <template #sidebar>
      <DocsSidebar />
    </template>
    <template #default>
      <div class="prose p-4 max-w-[800px]">
        <h1 id="using-sturdy">Using Sturdy</h1>
        <p>
          This tutorial takes you from setting up Sturdy with a repository from GitHub to merging a
          pull request from Sturdy. Make sure you have completed the Sturdy
          <router-link :to="{ name: 'v2DocsQuickStart' }">quick start</router-link> before
          continuing.
        </p>

        <h2 id="ui-overview">UI Overview</h2>
        <ol>
          <li>
            The left sidebar lists <strong>codebases</strong> and <strong>workspaces</strong> from
            your organization &mdash; click on a codebase or workspace to navigate to it
          </li>
          <li>
            The top left corner has a icon with tree dots &mdash; click the icon to open a menu with
            additional actions, such as <strong>creating a new codebase</strong>
          </li>
          <li>
            The bottom left corner has a bell icon &mdash; click the icon to open
            <strong>the notifications panel</strong>
          </li>
          <li>
            The bottom left corned also shows your avatar and a gearbox icon &mdash; click on either
            of them to go to your <strong>profile settings</strong>
          </li>
        </ol>

        <h2 id="import-a-repository">Setup Sturdy on a GitHub repository</h2>
        <ol>
          <li>Click on the three dots icon in the top left corner and select "Settings"</li>
          <li>Go to the GitHub-tab</li>
          <li>Authenticate yourself with GitHub</li>
          <li>Follow the GitHub prompt for authentication</li>
          <li>On GitHub, select the organization which contains the codebase you wish to setup</li>
          <li>You will be redirected back to Sturdy</li>
          <li>On Sturdy, click "Setup" on the repositories that you want to use from Sturdy</li>
        </ol>

        <h2 id="exploring-the-codebase">Exploring the codebase</h2>
        <p>
          Navigate to the codebase by selecting it's name in the sidebar on the left. This page
          shows you an overview of the codebase that shows you the following:
        </p>
        <ol>
          <li>A getting started panel</li>
          <li>The most recent changes</li>
          <li>A file explorer with the contents of the codebase</li>
          <li>The README.md file</li>
          <li>A list of collaborators</li>
        </ol>

        <h2 id="connecting-a-directory">Connecting a directory</h2>
        <p>With the codebase selected:</p>
        <ol>
          <li>
            Connect a directory where you want to have the code by clicking the "Connect directory"
            button
          </li>
          <li>
            Select a directory where you want to have the code. It is a good idea to create a new
            folder
          </li>
          <ol>
            <li>The folder you select must be empty</li>
            <li>The folder that you select should not be inside a git repository</li>
          </ol>
          <li>
            Observe that a new workspace was created with the name "Your Name's Workspace". In
            Sturdy, work always happens inside a workspace
          </li>
          <li>Workspaces are listed in the sidebar on the left, below the codebase</li>
          <li>Press the "Edit" button at the top of the page and rename it to "My first change"</li>
          <li>Any changes made files in the connected directory will show up here</li>
          <li>
            The folder you chose is connected to this workspace. This is indicated at the top of the
            workspace
          </li>
        </ol>

        <h2 id="creating-a-change">Creating your first change</h2>
        <p>With the workspace "My first change" selected:</p>
        <ol>
          <li>
            Open the connected folder on your computer with an IDE or text editor of your choice
          </li>
          <ol>
            <li>The directory connected to this workspace is indicated at the top</li>
            <li>
              <em>Optional:</em> Open the codebase in Visual Studio Code with the "Open in" button
            </li>
          </ol>
          <li>Using your IDE / text editor, make changes to any file in the codebase</li>
          <li>
            Within Sturdy &mdash; observe that changes you make are displayed as diffs in the "My
            first change" workspace
          </li>
        </ol>

        <h2 id="describing-a-change">Describing your changes</h2>
        <p>With the workspace "My first change" selected:</p>
        <ol>
          <li>Write a description of your changes in the text box at the top</li>
          <ol>
            <li>The description is automatically saved</li>
          </ol>
          <li>
            Observe that after a description has been provided a "Create pull request" button
            appears
          </li>
        </ol>

        <h2 id="creating-and-merging-a-pull-request">Creating and merging a pull request</h2>
        <p>With the workspace "My first change" selected and a description provided:</p>
        <ol>
          <li>Click the "Create pull request" button at the top</li>
          <li>Two new buttons will appear ("Update pull request" and "Merge")</li>
          <li>
            Click the "Merge" button. This is the same action as hitting "Merge" on GitHub, but
            without leaving Sturdy
          </li>
          <li>
            After a few seconds the PR will be merged, and the workspace will be rebased to the
            latest version of the codebase automatically
          </li>
        </ol>
        <DocsInfoBox>
          If your repository on GitHub requires approvals for pull requests, use the "Go to pull
          request" link in order to open the pull request in the browser.
        </DocsInfoBox>

        <h2 id="exploring-the-changelog">Exploring the codebase changelog</h2>
        <p>After the pull request from the "My first change" workspace has been merged:</p>
        <ol>
          <li>
            Navigate to the codebase overview by clicking on the codebase in the sidebar on the left
          </li>
          <li>The panel with recent changes now contains the change that you just created</li>
          <li>Click on the latest change</li>
          <li>You can see the diffs and description of the change as you created it</li>
          <li>
            Observe that on the right side you have a list of older changes &mdash; click on any of
            them to show you the details of that change
          </li>
        </ol>

        <h2 id="next-steps">Next steps</h2>
        <p>You can read more about:</p>
        <ul>
          <li>
            <router-link :to="{ name: 'v2DocsHowToCollaborateWithOthers' }">
              How to collaborate with others
            </router-link>
          </li>
          <li>
            <router-link :to="{ name: 'v2DocsHowSturdyAugmentsGit' }">
              How Sturdy augments Git
            </router-link>
          </li>
          <li>
            <router-link :to="{ name: 'v2DocsWorkingInTheOpen' }">
              Working in the open (within your team)
            </router-link>
          </li>
        </ul>
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
          content: 'How to navigate and use Sturdy',
        },
        {
          name: 'keywords',
          content: 'study learn documentation tutorial',
        },
      ],
      title: 'Using Sturdy',
    })
  },
})
</script>
