<template>
  <DocumentationWithTableOfContents>
    <template #sidebar>
      <DocsSidebar />
    </template>
    <template #default>
      <div class="prose p-4 max-w-[800px]">
        <h1 id="setup-with-github">How to setup Sturdy with GitHub</h1>

        <p>
          It's easy to add existing GitHub repositories to Sturdy. With a few clicks, you can
          install Sturdy on as many repositories as you want. <strong>Sturdy for GitHub</strong> is
          a fully-featured version of Sturdy that runs on top of your existing Git repository.
        </p>

        <ol>
          <li>
            In the Sturdy App: Click on the three dots icon in the top left corner and select
            "Settings"
          </li>
          <li>Go to the GitHub-tab</li>
          <li>
            Follow the instructions to authenticate with GitHub, and install the Sturdy
            GitHub-application on the repositories that you want to use Sturdy with
          </li>
          <li>You will be redirected back to Sturdy</li>
          <li>
            On Sturdy, under the GitHub-tab, you'll now see the list of repositories that you can
            install Sturdy on
          </li>
          <li>Find the repository that you want to use, and click "Setup"</li>
          <li>
            After a few seconds, the codebase will be ready. Click "Open" to open the codebase in
            Sturdy
          </li>
          <li>
            To edit your code, follow the
            <router-link :to="{ name: 'v2DocsHowToEditCode' }"
              >"How to: edit code" guide
            </router-link>
          </li>
        </ol>

        <p>
          Once you have the repository in a Sturdy codebase, all of the latest code as well as the
          history will be available within Sturdy. The README.md file is nicely rendered and you can
          browse files and folders, just as you are used to.
        </p>

        <p>
          Sturdy runs on top of the repository, and will automatically keep itself updated. Whenever
          new code is pushed to the default branch (main/master/etc) on GitHub, by your or somebody
          else, Sturdy will keep itself up to date automatically.
        </p>

        <DocsInfoBox>
          The newly created Sturdy codebase is also aware of your existing pull requests, which will
          be imported and show up as workspaces.
        </DocsInfoBox>

        <p>
          <em>Coming soon:</em> Pull requests created by others on your team will also show up as
          workspaces so that you can get the code with one click. This will work even if others on
          the team are not using Sturdy.
        </p>

        <h1 id="source-of-truth">Source of truth</h1>
        <p>
          When GitHub repositories are connected to Sturdy, this is <strong>more</strong> than just
          an import at a point in time. Connected repositories can be configured in multiple ways:
        </p>
        <ol>
          <li>
            <strong>GitHub as the source of truth</strong> &mdash; In this mode, Sturdy will follow
            changes happening on GitHub. Changes made in Sturdy are merged through GitHub pull
            requests. When pull requests are merged on GitHub, they will automatically be synced to
            Sturdy. This works even when changes are made outside of Sturdy. In this mode, some team
            members can develop through Sturdy without affecting others who use GitHub directly.
          </li>
          <li>
            <strong>Sturdy as the source of truth</strong> &mdash; In this mode, the state on GitHub
            follows Sturdy. After a change is shared on Sturdy, it will be pushed to the HEAD branch
            on GitHub. Changes made directly on GitHub will not be synced to Sturdy. In this mode
            all development happens through Sturdy.
          </li>
          <li>
            <strong>Disabled</strong> &mdash; The connection can be disabled, in which case the copy
            of the code on Sturdy and GitHub will diverge.
          </li>
        </ol>
        <p>By default, connected repositories use GitHub as the source of truth.</p>
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

export default defineComponent({
  components: { DocsSidebar, DocumentationWithTableOfContents, DocsInfoBox },
  setup() {
    useHead({
      meta: [
        {
          name: 'description',
          content: 'Learn how to setup Sturdy for Github',
        },
        {
          name: 'keywords',
          content: 'study learn documentation setup sturdy on github git hub install',
        },
      ],
      title: 'How to setup Sturdy on GitHub | Sturdy',
    })
  },
})
</script>
