<template>
  <PublicLeftSidebar>
    <template #sidebar>
      <DocsSidebar />
    </template>
    <template #default>
      <div class="prose p-4 max-w-[800px]">
        <h1 id="quick-start">Quick Start</h1>

        <p>There are two ways you could quickly try Sturdy out:</p>
        <ol>
          <li>
            <strong>Sturdy in the cloud</strong> &mdash; zero setup [<a href="#sturdy-cloud"
              >jump to ⤵️</a
            >]
          </li>
          <li>
            <strong>Self-hosted Sturdy instance</strong> &mdash; with a Docker one-liner [<a
              href="#self-hosted-sturdy"
              >jump to ⤵️️</a
            >️]
          </li>
        </ol>

        <h2 id="sturdy-cloud">Sturdy in the cloud</h2>
        <p>Trying Sturdy in the cloud required zero setup:</p>
        <ol>
          <li>
            Download and install
            <router-link :to="{ name: 'download' }">Sturdy</router-link>
          </li>
          <li>After opening Sturdy, click "Sign up now"</li>
          <li>Enter your name and email address</li>
          <li>Paste the code that you received in your email</li>
        </ol>

        <p>
          Next steps: <router-link :to="{ name: 'v2DocsUsingSturdy' }">Using Sturdy</router-link>
        </p>

        <h2 id="self-hosted-sturdy">Self-hosted Sturdy instance</h2>

        <p>With Docker installed and running on your computer:</p>

        <ol>
          <li>
            In your terminal run the following docker command:
            <pre>{{ dockerOneliner }}</pre>
          </li>
          <li>
            Download and install the Sturdy App
            <router-link :to="{ name: 'download' }">Sturdy</router-link>
          </li>

          <li>
            <p>
              From the Sturdy tray-icon, connect to the self hosted instance: Sturdy > Server >
              Self-hosted.

              <img src="./images/ServerPicker.png" height="424" width="720" class="rounded-md" />

              <br />
              <span class="test-sm text-gray-700"
                ><em>Don't see your server?</em> Restart the Sturdy app and try again.</span
              >
            </p>
          </li>
          <li>Follow the instructions in the app to register an account and setup Study</li>
        </ol>

        <p>
          Next steps: <router-link :to="{ name: 'v2DocsUsingSturdy' }">Using Sturdy</router-link>
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

export default defineComponent({
  components: { DocsSidebar, PublicLeftSidebar },
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
          content:
            'Sturdy is an open-source version control platform focused on expresiveness and real-time feedback',
        },
        {
          name: 'keywords',
          content:
            'study learn documentation open-source version control platform expressiveness real-time feedback why',
        },
      ],
      title: 'Quick Start | Sturdy',
    })

    return {
      dockerOneliner: `docker run --interactive --tty \\
    --publish 30080:80 --publish 30022:22 \\
    --volume "$HOME/.sturdydata:/var/data" \\
    getsturdy/server`,
    }
  },
})
</script>
