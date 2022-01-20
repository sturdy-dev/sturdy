<template>
  <PublicLeftSidebar>
    <template #sidebar>
      <DocsSidebar />
    </template>
    <template #default>
      <div class="prose p-4 max-w-[800px]">
        <h1 id="how-sturdy-augments-git">How to import code from GitHub</h1>

        <p>
          If you have an existing repository on GitHub, adding it to Sturdy is easy. With a few
          clicks you can give access to one or multiple specific repositories. This integration is
          not just for getting the code but also for the bridge interaction between Sturdy and
          GitHub.
        </p>

        <ol>
          <li>Click on the three dots icon in the top left corner and select "New codebase"</li>
          <li>
            On the codebase creation screen click the "Import repositories" button within the
            "Coming from GitHub?" panel in the bottom.
          </li>
          <li>On GitHub, select the organization which contains the codebase you want to import</li>
          <li>Follow the GitHub prompt for authentication</li>
          <li>Select one or multiple repositories you wish to work with</li>
          <li>You will be redirected back to Sturdy</li>
          <li>
            When the repositories are ready within Sturdy, you will be notified - this is usually
            instant
          </li>
          <li>
            To edit your code, follow the
            <router-link :to="{ name: 'v2DocsHowToEditCode' }"
              >"How to: edit code" guide</router-link
            >
          </li>
        </ol>

        <p>
          Once you have the repository in a Sturdy codebase, all of the latest code as well as the
          history will be available within Sturdy. The README.md file is nicely rendered and you can
          browse files and folders, just as you are used to.
        </p>

        <p>
          It is important to note that whenever new code is pushed to the default (main/master)
          branch on GitHub, by your or somebody else, Sturdy will keep itself up to date
          automatically.
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
