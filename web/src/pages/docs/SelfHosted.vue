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

        <h2 id="setup-github-integration">Optional: Setup GitHub Integration</h2>

        <h3 id="setup-proxy">Setup tunnel to localhost</h3>
        <p>
          GitHub integration heavily unilizes webhooks. To make that work, we have to ensure that
          self-hosted Sturdy instance can receive HTTP requests over the Internet.
        </p>

        <DocsInfoBox>
          If you have your own way of setting up a tunnel, you can skip this step.
        </DocsInfoBox>

        <p>
          Internally, at Sturdy, we are using <a href="https://ngrok.com/">ngrok</a> to setup a
          tunnel to localhost. It's a free and easy to use service that allows you to expose your
          local Sturdy instance to the Internet in no time. Here is how to set it up:
        </p>

        <ol>
          <li>
            Install <em>ngrok</em> client following instructions on their
            <a href="https://ngrok.com/download">Download page</a>
          </li>
          <li>
            Run it like so:
            <pre>ngrok http 30080</pre>
          </li>
        </ol>

        <p>
          Now your local port 30080 is exposed to the Internet, and in console, you will see a URL
          to access it. It would look something line this:
        </p>
        <pre>https://09f9-213-114-130-156.ngrok.io</pre>

        <h3 id="create-a-github-app">Create a GitHub App</h3>
        <p></p>

        <ol>
          <li>
            Go to
            <a href="https://github.com/settings/apps/new">https://github.com/settings/apps/new</a>
          </li>
          <li>Set the app name, for example <em>sturdy-self-hosted</em></li>
          <li>Set the homepage to <em>https://localhost:30080</em></li>
          <li>
            <u>Callback URLs</u>
            <ul>
              <li>sturdy:///setup-github</li>
            </ul>
          </li>
          <li><b>Untick</b> <em>Expire user authorization tokens</em></li>
          <li><b>Tick</b> <em>Request user authorization (OAuth) during installation</em></li>
          <li><b>Tick</b> <em>Redirect on update</em></li>
          <li>Make sure that webhooks are active</li>
          <li>
            Set the webhook URL to <em>${YOUR_LOCAL_TUNNEL_HOSTNAME}/api/v3/github/webhook</em>

            <DocsInfoBox>
              If you use ngrok, keep in mind that every time it will generate a new URL. So don't
              forget to update it later.
            </DocsInfoBox>
          </li>
          <li>
            <u>Permissions</u>
            <ul>
              <li>Checks - Read & Write</li>
              <li>Contents - Read & Write</li>
              <li>Metadata - Read only</li>
              <li>Pull Requests - Read & Write</li>
              <li>Workflows - Read & Write</li>
            </ul>
          </li>

          <li>
            <u>Subscribe to events</u>
            <ul>
              <li>Pull Request</li>
              <li>Pull request review</li>
              <li>Push</li>
            </ul>
          </li>

          <li>Click <em>Create GitHub App</em></li>
          <li>Click <em>Generate a new Client secret</em> and remember it</li>
          <li>In the bottom of the page, click <em>Generate a private key</em> and download it</li>
        </ol>

        <p>
          That's a lot of configuration! But now we are almost ready to start using Sturdy with
          GitHub integration.
        </p>

        <h3 id="add-github-app-configuration-to-the-app">
          Restart Sturdy with GitHub configuration
        </h3>

        <p>
          To finish the configuration, we need to restart Sturdy with some of the configuration from
          the app we've just created:
        </p>

        <pre>{{ dockerOnelinerWithGithub }}</pre>

        <p>
          Note that your github private key must be inside the
          <em>$HOME/.sturdydata</em> directory.
        </p>

        <p>Congratulations! You are now ready to use Sturdy with GitHub integration.</p>

        <DocsInfoBox
          >To learn more about how to use GitHub integration, see
          <router-link :to="{ name: 'v2DocsHowToSetupSturdyOnGitHub' }"
            >How to setup Sturdy on GitHub</router-link
          >
        </DocsInfoBox>

        <h2 id="license">License</h2>

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
import DocsInfoBox from '../../molecules/DocsInfoBox.vue'

export default defineComponent({
  components: { DocsSidebar, DocumentationWithTableOfContents, DocsInfoBox },
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
          content: 'Learn how to setup self-hosted Sturdy',
        },
        {
          name: 'keywords',
          content: 'study learn documentation self-hosted github enterprise local',
        },
      ],
      title: 'Self-hosted | Sturdy',
    })

    return {
      dockerOneliner: DOCKER_ONELINER,
      dockerOnelinerWithGithub: DOCKER_ONELINER.replace(
        'getsturdy/server',
        `--env STURDY_GITHUB_APP_ID=&lt;id> \\
    --env STURDY_GITHUB_APP_NAME=<name> \\
    --env STURDY_GITHUB_APP_CLIENT_ID=<client_id> \\
    --env STURDY_GITHUB_APP_SECRET=<secret> \\
    --env STURDY_GITHUB_APP_PRIVATE_KEY_PATH=/var/data/<filename> \\
        getsturdy/server`
      ),
    }
  },
})
</script>
