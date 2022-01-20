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
          <li>Sturdy in the cloud - zero setup [<a href="#sturdy-cloud">jump to ⤵️</a>]</li>
          <li>
            Self hosted Sturdy instance - with a Docker one-liner [<a href="#self-hosted-sturdy"
              >jump to ⤵</a
            >️]
          </li>
        </ol>

        <h2 id="sturdy-cloud">Sturdy in the cloud</h2>
        <p>Trying Sturdy in the cloud required zero setup:</p>
        <ol>
          <li>
            Download and install
            <router-link :to="{ name: 'download' }">the Sturdy app</router-link>
          </li>
          <li>After opening Sturdy, press the "Sign up now" button</li>
          <li>Enter your name and email address</li>
          <li>Paste the code that you received in your email</li>
        </ol>

        <DocsInfoBox>Sturdy in the cloud uses a passwordless login</DocsInfoBox>

        <p>Next steps: Using Sturdy end-to-end tutorial</p>
        <!-- TODO: Add link -->

        <h2 id="self-hosted-sturdy">Self hosted Sturdy instance</h2>

        <p>With Docker installed and running on your computer:</p>

        <ol>
          <li>
            In your terminal run the following docker command:
            <pre>
docker run --detach --publish 30080:80 --publish 30022:22 \
           --volume "$(pwd)/sturdydata:/var/data" \
           getsturdy/server</pre
            >
          </li>
          <!-- TODO: update -->
          <li>In your browser, go <a href="http://localhost:9001">http://localhost:9001</a></li>
          <li>Follow the instructions to setup your organization</li>
          <li>
            Download and install
            <router-link :to="{ name: 'download' }">the Sturdy app</router-link>
          </li>
          <li>After opening Sturdy, press the "Configure with self-hosted instance" button</li>
          <li>In the Sturdy host field enter "localhost" and press "Save"</li>
          <li>Press the "Sign up now" button</li>
          <li>Enter your name, email address and choose a password</li>
        </ol>

        <p>Next steps: Using Sturdy end-to-end tutorial</p>
        <!-- TODO: Add link -->
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
