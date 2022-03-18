<template>
  <DocumentationWithTableOfContents>
    <template #sidebar>
      <DocsSidebar />
    </template>
    <template #default>
      <div class="prose p-4 max-w-[800px]">
        <h1 id="using-sturdy">Using Sturdy</h1>

        <p>
          This tutorial takes you from setting up Sturdy with a repository from GitHub to merging a
          pull request from Sturdy. Make sure you have completed the Sturdy
          <router-link :to="{ name: 'v2DocsQuickStart' }">quick start</router-link>
          before continuing.
        </p>

        <h2 id="ui-overview">UI Overview</h2>

        <img src="./images/UI_overview.png" width="1551" height="937" alt="Sturdy UI overview" />
        <ol>
          <li>
            The left sidebar lists <strong>codebases</strong> and open
            <strong>draft changes</strong> from your organization &mdash; click on a codebase or
            draft change to navigate to it
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

        <h2 id="first-codebase">Your first codebase</h2>
        <p>In Sturdy you can create a codebase in two ways:</p>
        <ul>
          <li>Create new codebases [<a href="#create-a-codebase">go to</a>]</li>
          <li>
            Connect one or multiple existing repositories from GitHub [<a
              href="#import-a-repository"
              >go to</a
            >]
          </li>
        </ul>

        <DocsInfoBox>
          If you are running a self-hosted Sturdy instance and would like to connect a GitHub
          repository, you will need to first perform additional setup.
          <br />
          Follow
          <router-link :to="{ name: 'v2DocsSelfHosted', hash: '#setup-github-integration' }"
            >this guide
          </router-link>
          before continuing.
        </DocsInfoBox>

        <h3 id="create-a-codebase">Create a Sturdy standalone codebase</h3>
        <ol>
          <li>Click on the three dots icon in the top left corner and select "New codebase"</li>
          <li>Type in a codebase name and press "Create"</li>
        </ol>

        <h3 id="import-a-repository">Setup Sturdy on a GitHub repository</h3>
        <p>
          Repositories imported from GitHub will remain up-to-date with changes made outside of
          Sturdy (<router-link
            :to="{ name: 'v2DocsHowToSetupSturdyWithGitHub', hash: '#source-of-truth' }"
            >read more</router-link
          >).
        </p>
        <ol>
          <li>Click on the three dots icon in the top left corner and select "Settings"</li>
          <li>Go to the GitHub-tab</li>
          <li>Authenticate yourself with GitHub</li>
          <li>
            Update the GitHub installation, selecting one or multiple repositories
            <img
              src="./images/Add_repository_on_GitHub.png"
              width="1175"
              height="586"
              alt="Updating the selected repositories on GitHub"
            />
          </li>
          <li>On GitHub, select the organization which contains the codebase you wish to setup</li>
          <img
            src="./images/GitHub_select_repository.png"
            width="1157"
            height="640"
            alt="Selecting repository access on GitHub"
          />
          <li>You will be redirected back to Sturdy</li>
          <li>On Sturdy, click "Setup" on the repositories that you want to use from Sturdy</li>
          <img
            src="./images/Setup_repository.png"
            width="1184"
            height="682"
            alt="Creating a Sturdy codebase connected to a GitHub repository"
          />
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
          <img
            src="./images/Connect_directory.png"
            width="1248"
            height="834"
            alt="Connecting a local directory to a Sturdy draft change"
          />
          <li>
            Select a directory where you want to have the code in. Sturdy will create a new
            directory inside the directory that you select with the same name as the codebase.
          </li>
          <li>
            Observe that a new draft change was created with the name "Untitled Draft". In Sturdy,
            work always happens inside a draft change.
          </li>
          <li>Open draft changes are listed in the sidebar on the left, below the codebase</li>
          <li>Press the "Edit" button at the top of the page and rename it to "My first change"</li>
          <li>Any changes made files in the connected directory will show up here</li>
          <li>
            The folder you chose is connected to this draft. This is indicated at the top of the
            draft change page.
          </li>
        </ol>

        <h2 id="creating-a-change">Creating your first change</h2>
        <p>With the draft change "My first change" selected:</p>
        <ol>
          <li>
            Open the connected folder on your computer with an IDE or text editor of your choice
          </li>
          <ol>
            <li>The directory connected to this draft is indicated at the top</li>
            <li>
              <SubtleText>
                <em>Optional:</em> Open the codebase in Visual Studio Code with the "Open in" button
              </SubtleText>
            </li>
          </ol>
          <img
            src="./images/Creating_first_change.png"
            width="1248"
            height="1102"
            alt="Creating your first change on Sturdy"
          />
          <li>Using your editor, open a file in the codebase and make some changes to it</li>
          <li>
            Within Sturdy &mdash; observe that changes you make are displayed as diffs in the "My
            first change" draft
          </li>
        </ol>

        <h2 id="describing-a-change">Describing your changes</h2>
        <p>With the draft "My first change" selected:</p>
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
        <p>With the draft change "My first change" selected and a description provided:</p>
        <ol>
          <li>Click the "Create pull request" button at the top</li>
          <img
            src="./images/Before_creating_PR.png"
            height="222"
            width="1088"
            alt="Change description in a Sturdy draft change"
          />
          <li>Two new buttons will appear ("Update pull request" and "Merge")</li>
          <li>
            Click the "Merge" button. This is the same action as hitting "Merge" on GitHub, but
            without leaving Sturdy
          </li>
          <img
            src="./images/After_creating_PR.png"
            height="235"
            width="1118"
            alt="Managing a GitHub Pull Request from Sturdy"
          />
          <li>
            After a few seconds the PR will be merged. A new draft change will automatically be
            created for you, so that you're ready to code again right away!
          </li>
        </ol>
        <DocsInfoBox>
          If your repository on GitHub requires approvals for pull requests, use the "Go to pull
          request" link in order to open the pull request in the browser.
        </DocsInfoBox>

        <h2 id="exploring-the-changelog">Exploring the codebase changelog</h2>
        <p>After the pull request from the "My first change" draft has been merged:</p>
        <ol>
          <li>
            Navigate to the codebase overview by clicking on the codebase in the sidebar on the left
          </li>
          <li>The panel with recent changes now contains the change that you just created</li>
          <li>Click on the Changelog</li>
          <li>You can see the list of changes, including the one you've just created</li>
          <img
            src="./images/Previous_changes.png"
            height="837"
            width="1312"
            alt="Exploring previous changes"
          />
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
            <router-link :to="{ name: 'v2DocsHowSturdyInteractsWithGit' }">
              How Sturdy interacts with Git
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
      content: 'How to navigate and use Sturdy',
    },
    {
      name: 'keywords',
      content: 'study learn documentation tutorial',
    },
  ],
  title: 'Using Sturdy',
})
</script>
