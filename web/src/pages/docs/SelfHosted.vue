<template>
  <DocumentationWithTableOfContents>
    <template #sidebar>
      <DocsSidebar />
    </template>
    <template #default>
      <div class="prose p-4 max-w-[800px]">
        <h1 id="working-in-the-open">Self-hosted Sturdy</h1>

        <p>Run Sturdy anywhere, in your own datacenter or on your own workstation.</p>

        <h2>Self-hosted Sturdy in Docker</h2>

        <p>
          The easiest way to run Sturdy is to run the one-liner bundle.
          <code>getsturdy/server</code> contains the full Sturdy server and all of it's dependencies
          in a single container.
        </p>

        <p>
          The container supports both <code>linux/amd64</code> and <code>linux/arm64</code> out of
          the box.
        </p>

        <ol>
          <li>
            Start Sturdy in Docker
            <pre>{{ dockerOneliner }}</pre>
          </li>
          <li>
            <router-link :to="{ name: 'v2download' }">
              Download and install the Sturdy app for macOS and Windows
            </router-link>
          </li>
          <li>Start the app</li>
          <li>
            <p>
              <span
                >From the tray icon, switch to connect to the self-hosted instance of Sturdy running
                on your computer</span
              >
              <img class="rounded-md" src="./images/ServerPicker.png" height="424" width="720" />
            </p>
          </li>
          <li>
            The Sturdy app will auto-detect instances running locally on the default ports
            <code>30080</code> and
            <code>30022</code>
          </li>
          <li>
            Follow the instructions in the app to setup your account and configure your server!
          </li>
        </ol>

        <h2>License</h2>

        <p>
          The published Docker image <code>getsturdy/server</code> contains Sturdy Enterprise, and
          is licensed under the Sturdy Enterprise License.
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
import { DOCKER_ONELINER } from '../../docker'

export default defineComponent({
  components: { DocsSidebar, DocumentationWithTableOfContents },
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
          content: 'Working in the open (within your team)',
        },
        {
          name: 'keywords',
          content: 'study learn documentation working code in the open publicly public',
        },
      ],
      title: 'Self-hosted Sturdy',
    })

    return {
      dockerOneliner: DOCKER_ONELINER,
    }
  },
})
</script>
