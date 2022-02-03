<template>
  <DocumentationWithTableOfContents>
    <template #sidebar>
      <DocsSidebar />
    </template>
    <template #default>
      <div class="prose p-4 max-w-[800px]">
        <h1 id="how-sturdy-augments-git">How to ship software to production</h1>
        <span class="">&hellip;or just get the code merged.</span>

        <p>
          Sturdy encourages shipping small incremental changes, so it optimizes for making the "code
          &mdash; describe changes &mdash; ship" loop as tight as possible.
        </p>

        <p><a href="#on-github">Jump to: Sturdy for GitHub ⤵️</a></p>

        <h2 id="on-sturdy">Sturdy</h2>

        <DocsInfoBox>
          These instructions are for codebases on Sturdy that are
          <strong>not connected to GitHub</strong>, or where
          <strong>Sturdy is the source of truth</strong>. Go to
          <a href="#on-github">Sturdy for GitHub</a> for instructions when using Sturdy on GitHub
          with GitHub as the source of truth.
        </DocsInfoBox>

        <p>
          With Sturdy setup and
          <router-link :to="{ name: 'v2DocsHowToEditCode' }">connected to a directory</router-link>
          on your computer:
        </p>

        <ol>
          <li>Create a new workspace, and give it a good name</li>
          <li>
            Open the directory with the code in any text editor
            <SubtleText>(tip: use the “Open in Visual Studio Code”-button)</SubtleText>
          </li>
          <li>Make the changes to your code</li>
          <li>Go back to Sturdy, and navigate to the workspace from the sidebar</li>
          <li>
            <SubtleText><em>Optional:</em> Undo changes you would like to dismiss</SubtleText>
          </li>
          <li>
            On the top of the page, use the description box to describe the changes that you have
            made
          </li>
          <li>
            <SubtleText
              ><em>Optional:</em> Resolving conflicts [<a href="#optional-resolving-conflicts"
                >go to</a
              >]</SubtleText
            >
          </li>
          <li>
            <SubtleText
              ><em>Optional:</em> Asking for feedback [<a href="#optional-asking-for-feedback"
                >go to</a
              >]</SubtleText
            >
          </li>
          <li>
            <SubtleText
              ><em>Optional:</em> Syncing on top of changelog [<a href="#optional-syncing">go to</a
              >]</SubtleText
            >
          </li>
          <li id="writing-a-description">
            After writing a description, the “Share” button will appear
          </li>
          <li>
            Click “Share!” to land all of the changes on top of the changelog, and to allow other
            workspaces to be based on top of your code
          </li>
          <li>
            After the change is shared, the workspace will be rebased to the latest version of the
            codebase automatically
          </li>
        </ol>

        <h2 id="on-github">Sturdy for GitHub</h2>

        <DocsInfoBox>
          These instructions are for codebases on Sturdy that are
          <strong>connected to GitHub</strong> and where GitHub is the
          <strong>source of truth</strong>.
        </DocsInfoBox>

        <p>
          Using Sturdy on with the GitHub bridge is similar to Sturdy without the git-bridge
          (standalone Sturdy), with some differences:
        </p>

        <ol>
          <li>
            Follow step 1 to 10 of the instructions for Sturdy without GitHub [<a href="#on-sturdy"
              >go to</a
            >]
          </li>
          <li>After writing a description, the “Create pull request” button will appear</li>
          <li>
            Click “Create pull request” to use Sturdy to a PR to the connected GitHub repository. A
            single commit will be created containing all of the changes in the workspace. The
            Workspace title and description will be used as the description of the commit, and the
            title/description of the pull request.
          </li>
          <li>
            <SubtleText
              ><em>Optional:</em> If the above fails, debug GitHub errors [<a
                href="#debugging-github-errors"
                >go to</a
              >]</SubtleText
            >
          </li>
          <li>
            New buttons will appear. “Go to pull request”, “Update pull request”, and “Merge”.
          </li>
          <li>“Go to pull request”: Open the PR in the browser</li>
          <li>“Update pull request”: update the PR title, description, and code from Sturdy.</li>
          <li>
            “Merge”: Merges the PR, it’s the same action as hitting “Merge” on GitHub, but without
            leaving Sturdy.
          </li>
          <li>
            Click “Merge” to merge the PR. After a few seconds the PR will be merged, and the
            workspace will be rebased to the latest version of the codebase automatically.
          </li>
        </ol>

        <h1>Help and debug</h1>

        <h2 id="optional-undo-changes" title="Undo changes">
          Optional: Undo changes you would like to dismiss
        </h2>

        <p>
          Using the checkbox on the left of file changes mark one or multiple changes you would like
          to dismiss
        </p>
        <p>
          In the blue banner at the top of the screen press "Undo" (NB: This action cannot be
          reversed)
        </p>
        <p>The changes that you had selected are removed from your workspace</p>

        <h2 id="optional-resolving-conflicts" title="Resolving conflicts">
          Optional: Resolving conflicts
        </h2>

        <p>
          If the changes in the workspace are conflicting with recent changes in the codebase, you
          have a conflict. Conflicts happen when two workspaces make changes to the same piece of
          code.
        </p>

        <p>A banner will be shown in the workspace if it’s conflicting.</p>

        <ol>
          <li>To start resolving a conflict, click “Sync” in the upper right corner</li>
          <li>For each file with conflicts, you now have three options.</li>
          <ol>
            <li>Pick the version of the entire file as it is on trunk</li>
            <li>Pick the version of the entire file as it is in the workspace</li>

            <li>Create a mix of the two (“Custom” resolution)</li>
          </ol>
          <li>Custom resolution:</li>
          <ol>
            <li>
              If you want to create a custom resolution, pick “Custom”. And open the file in your
              editor. It will now contain both versions of the file at the same time, using Git
              conflict markers (<code>&lt;&lt;&lt;&lt;&lt;&lt;&lt;</code>, <code>=======</code>,
              <code>&gt;&gt;&gt;&gt;&gt;&gt;&gt;</code>).
            </li>
            <li>Edit the file to your liking.</li>
          </ol>
          <li>
            When a resolution (trunk, workspace, or custom) has been picked for each file. Mark the
            conflicts as completed by clicking “Done”.
          </li>
          <li>
            The workspace should no longer have any conflicts, and the changes are ready to be
            shared.
          </li>
        </ol>

        <h2 id="optional-asking-for-feedback" title="Ask for feedback">
          Optional: Asking for feedback
        </h2>

        <DocsInfoBox
          >Sturdy makes it super easy to give feedback and suggest code to one another. It embraces
          working in the open within the team and favours early feedback over "reviewing" after the
          fact.</DocsInfoBox
        >

        <p>
          It might be a good idea to ask for some feedback from a colleague before sharing the
          changes. When browsing a workspace in Sturdy, use the “Ask for feedback” button on the
          right sidebar to ask someone to take a look at your code. They can either give you a
          thumbs up, or give you a comment with feedback.
        </p>

        <h2 id="optional-syncing" title="Syncing">Optional: Syncing</h2>

        <p>
          After creating the workspace, if new changes land on the changelog (either from others on
          the team or from you via another workspace), these changes are not available in your
          workspace by default.
        </p>

        <p>
          If there are new changes available, the "Sync changes" button will indicate this with a
          blue flash.
        </p>

        <p>To rebase your workspace on the latest changes:</p>

        <ol>
          <li>Press the “Sync changes” button</li>
          <li>Optional: Resolve conflicts [<a href="#optional-resolving-conflicts">go to</a>]</li>
        </ol>

        <p>
          It is not required to manually sync changes before sharing. This will be performed
          automatically.
        </p>

        <h2 id="debugging-github-errors" title="Debug GitHub">Debugging GitHub errors</h2>

        <p>
          The GitHub bridge from Sturdy ("Sturdy for GitHub") requires your user to have
          <strong>PUSH/WRITE</strong> permissions to the main repository (PUSH permissions to a fork
          is not enough).
        </p>

        <p>
          Sturdy will push to branches on the pattern <code>sturdy-pr-${UUID}</code>, make sure that
          there are no branch protection rules blocking pushes to branches with this pattern.
        </p>
      </div>
    </template>
  </DocumentationWithTableOfContents>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import DocumentationWithTableOfContents from '../../layouts/DocumentationWithTableOfContents.vue'
import { useHead } from '@vueuse/head'
import DocsSidebar from '../../organisms/docs/DocsSidebar.vue'
import DocsInfoBox from '../../molecules/DocsInfoBox.vue'
import SubtleText from '../../molecules/docs/SubtleText.vue'

export default defineComponent({
  components: { DocsSidebar, DocumentationWithTableOfContents, DocsInfoBox, SubtleText },
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
          content: 'Learn how to ship software to production with Sturdy',
        },
        {
          name: 'keywords',
          content: 'study learn documentation ship deploy software code to prod production',
        },
      ],
      title: 'How to ship software to production | Sturdy',
    })
  },
})
</script>
