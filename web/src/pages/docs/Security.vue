<template>
  <StaticPage
    title="Security"
    category="knowledgebase"
    metadescription="The safety and integrity of your code is our highest priority. Your code and your data is
        always treated with respect. It's always encrypted in transport and at rest, and is stored
        in the European Union."
    :wide="true"
    image="https://images.unsplash.com/photo-1609358905581-e5381612486e?ixlib=rb-1.2.1&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=708&q=80"
  >
    <div class="text-base max-w-prose mx-auto lg:max-w-none">
      <p class="text-lg text-gray-500">
        The safety and integrity of your code is our highest priority. Your code and your data is
        always treated with respect. It's always encrypted in transport and at rest, and is stored
        in the European Union.
      </p>
    </div>
    <div
      class="mt-5 prose prose-red text-gray-500 mx-auto lg:max-w-none lg:row-start-1 lg:col-start-1"
    >
      <p>Sturdy stores data in the European Union and is compliant with the GDPR.</p>
      <p>Your data is always encrypted in transit and at rest.</p>
    </div>

    <div
      class="mt-5 prose prose-red text-gray-500 mx-auto lg:max-w-none lg:row-start-1 lg:col-start-1"
    >
      <h2 class="text-center">Platform security</h2>
      <h3>Physical Access</h3>
      <p>Sturdy does not have any physical servers. All of our services run in the cloud.</p>

      <h2 class="text-center">Application Security</h2>
      <div class="grid grid-cols-2 gap-8">
        <div>
          <h3>Files on your computer</h3>
          <p>
            Sturdy stores metadata in the users home directory in the following directories and
            files
          </p>
          <p></p>
          <ul>
            <li><code>~/.sturdy-agent</code></li>
            <li><code>~/.sturdy-sync</code></li>
            <li><code>~/.sturdy</code></li>
          </ul>

          <p>
            Sturdy also needs read/write access to the user-provided directories that Sturdy should
            monitor for code changes.
          </p>

          <p>
            Sturdy will also automatically access the <code>~/.ssh/known_hosts</code> file to
            configure SSH trust to Sturdy's servers.
          </p>
        </div>

        <div>
          <h3>Testing and review</h3>
          <p>
            Sturdy applies code review and automated and manual testing of all changes to the Sturdy
            software.
          </p>

          <h3>File communication to Sturdy</h3>
          <p>
            The local application communicates with the Sturdy API over HTTP with TLS, using a
            temporary (30 day) JWT (JSON Web Token). This token is rotated by the client in the
            background.
          </p>

          <p>
            File transfers between a client and Sturdy are made over SSH, and are authenticated with
            4096-bit RSA keys. The private key(s) are generated on the client, and will never leave
            the client computer.
          </p>

          <h3>User logins</h3>
          <p>
            User passwords are cryptographically hashed and salted before being stored in our
            database.
          </p>
        </div>
      </div>

      <h2 class="text-center">Security Policies</h2>

      <div class="grid grid-cols-2 gap-8">
        <div>
          <h3>Employee access to data</h3>
          <p>
            Sturdy employees will only ever access customer data when it’s required for support
            related duties. When this happens, staff will do their best to respect customer privacy,
            and only access the code and settings required to debug the issue at hand.
          </p>
          <p>All sessions and actions are logged, and regularly reviewed.</p>
        </div>
      </div>

      <h2 class="text-center">Contact us</h2>
      <div class="flex justify-around">
        <div class="text-center">
          <p>Don't hesitate to reach out to us if you have any questions.</p>

          <LinkButton href="mailto:support@getsturdy.com" class="!no-underline !text-black !m-0">
            <span>support@getsturdy.com</span>
          </LinkButton>
        </div>
      </div>
    </div>
  </StaticPage>
</template>

<script lang="ts" setup>
import StaticPage from '../../layouts/StaticPage.vue'
import LinkButton from '../../components/shared/LinkButton.vue'
</script>
