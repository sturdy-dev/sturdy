<template>
  <DocumentationWithTableOfContents>
    <template #sidebar>
      <DocsSidebar />
    </template>
    <template #default>
      <div class="prose p-4 max-w-[800px]">
        <h1 id="how-to-switch-between-tasks">How to switch between tasks</h1>
        <p>
          Sturdy allows you to easily and quickly switch between tasks by using multiple draft change.
          Draft changes are a unit of isolation in Sturdy, which means that code changes in one
          draft are independent from those in another. It is quick to create and dispose of
          draft. This guide will go through how you can use multiple drafts to seamlessly
          switch between tasks. Make sure you are familiar with
          <router-link :to="{ name: 'v2DocsHowToEditCode' }">editing code</router-link>
          in Sturdy before continuing.
        </p>

        <h2 id="switching-workspaces">Multitasking with draft changes</h2>
        <p>
          You can use draft change to jump between multiple work in progress items. After you have
          created, imported or joined an existing codebase:
        </p>

        <ol>
          <li>
            Create a new draft change by clicking the "plus" icon next to the codebase name in the left
            sidebar
          </li>
          <li>Navigate to the new draft named "<em>Untitled draft</em>"</li>
          <li>Rename the newly created draft to "Feature Alpha" using the "Edit" button</li>
          <li>
            Click the "Connect <em>folder</em> for editing" button at the top of the draft change page
          </li>
          <li>The chosen <em>folder</em> is now connected to the "Feature Alpha" draft</li>
          <li>Using your IDE of choice, create some changes to the code</li>
          <li>The diffs of your changes appear in the "Feature Alpha" d</li>
          <li>
            Create a new draft change by clicking the "plus" icon next to the codebase name in the left
            sidebar
          </li>
          <li>Navigate to the new draft named "<em>Untitled draft</em>"</li>
          <li>Rename the newly created draft change to "Feature Beta" using the "Edit" button</li>
          <li>
            Click the "Connect <em>folder</em> for editing" button at the top of the draft change page
          </li>
          <li>The chosen <em>folder</em> is now connected to the "Feature Beta" draft change</li>
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
              This is because the changes were made in draft "Feature Alpha" and your directory
              is now connected to the "Feature Beta" draft change
            </li>
          </ol>
          <li>Using your IDE:</li>
          <ol>
            <li>
              Create some new file changes, different from the ones you did earlier for the first
              draft
            </li>
          </ol>
          <li>The diffs of your changes appear in the "Feature Beta" draft</li>
          <li>Using the sidebar on the left, navigate back to the "Feature Alpha" draft</li>
          <ol>
            <li>In Sturdy you can see the diffs that you made earlier</li>
          </ol>
          <li>
            Click the "Connect <em>folder</em> for editing" button at the top of the draft change page
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
          Draft changes in Sturdy are lightweight and can be created at any time. This allows you to
          have multiple experiments at the same time, or the ability to quickly fix a bug without
          disrupting your primary task.
        </p>

        <DocsInfoBox>
          In Sturdy, there is no need to 'save' or 'stash' changes when switching between
          drafts.
        </DocsInfoBox>

        <DocsInfoBox class="mt-8">
          You can connect the directory on your computer to draft changes created by others on your
          team. Learn more about
          <router-link :to="{ name: 'v2DocsHowToCollaborateWithOthers' }"
            >collaborating with others
          </router-link>
          .
        </DocsInfoBox>

        <h2 id="workspaces-based-on-other-workspaces">Draft changes bases on other drafts</h2>
        <p>
          It can be useful to create drafts based on other drafts or on parts of one. Consider the
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
          Sturdy allows you to create new draft change based on the work in an existing draft.
          Within an existing draft which <em>already has</em> some work in progress code
          changes:
        </p>
        <ol>
          <li>
            Using the checkbox on the left of changes, select one or multiple diffs that you wish to
            copy to a new draft change
          </li>
          <ol>
            <li>
              <SubtleText>
                <em>Optional</em>: Use the "Select All" checkbox at the top in order to quickly
                select all diffs
              </SubtleText>
            </li>
          </ol>
          <li>
            Observe the banner in blue at the top of the app indicating "<em>n</em> change(s)
            selected"
          </li>
          <li>
            Click the "Copy to New Draft" button to create a new draft change with the selected
            changes
          </li>
          <li>
            Observe that a new draft change named "Fork of <em>name of the original</em>" has been
            created
          </li>
          <li>
            Using the sidebar on the left, navigate to the draft "Fork of
            <em>name of the original</em>"
          </li>
          <li>Observe that the new draft change has only the changes that you selected</li>
          <li>Rename the draft change to reflect the changes included, for example "Quick bugfix"</li>
        </ol>

        <p>
          Creating draft change based on other drafts is a powerful tool that allows you to
          quickly and easily ship small incremental changes.
        </p>
      </div>
    </template>
  </DocumentationWithTableOfContents>
</template>

<script lang="ts" setup>
import DocumentationWithTableOfContents from '../../layouts/DocumentationWithTableOfContents.vue'
import { useHead } from '@vueuse/head'
import DocsSidebar from '../../organisms/docs/DocsSidebar.vue'
import DocsInfoBox from '../../molecules/DocsInfoBox.vue'
import SubtleText from '../../molecules/docs/SubtleText.vue'

useHead({
  meta: [
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
</script>
