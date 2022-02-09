<template>
  <StaticPage
    title="GraphQL API"
    subtitle="Integrate anything with Sturdy"
    category="documentation"
    metadescription="Use Sturdys GraphQL API to integrate anything with Sturdy"
    image="https://images.unsplash.com/photo-1517373116369-9bdb8cdc9f62?ixlib=rb-1.2.1&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1950&q=80"
  >
    <div class="mt-5 text-gray-500 mx-auto lg:max-w-none lg:row-start-1 lg:col-start-1">
      <div class="prose prose-yellow">
        <p>Use Sturdys GraphQL API to integrate <em>anything</em> with Sturdy!</p>
      </div>

      <div class="space-x-2 py-4">
        <LinkButton href="https://schema.getsturdy.com/">Schema Documentation</LinkButton>
        <LinkButton href="https://schema.getsturdy.com/sturdy.graphql">
          Download the full schema
        </LinkButton>
      </div>

      <div class="prose prose-yellow">
        <h2>The GraphQL Endpoint</h2>
        <p>In GraphQL, all queries are made to a single endpoint:</p>
        <pre>https://api.getsturdy.com/graphql</pre>

        <h2>Authentication</h2>
        <p>
          Authenticate yourself to the API by providing the
          <code>Authorization</code> header.
        </p>
        <pre><code>Authorization: bearer <em>token</em></code></pre>
        <p>Note that all queries must be authenticated. Get your personal access token below.</p>

        <ClientOnly>
          <h2>cURL example</h2>
          <pre>
curl -v -H "Authorization: bearer {{ token }}" \
  -d '{"query":"query { user { name email avatarUrl }}"}' \
  https://api.getsturdy.com/graphql
</pre
          >
        </ClientOnly>

        <h2>Recommendations</h2>
        <p>
          We're recommending to use a GraphQL app to explore the API and create your queries. Sturdy
          supports Schema Inspection, these clients will automatically download the API schema and
          documentation when configured to connect to the Sturdy API.
        </p>
        <ul>
          <li>
            <a href="https://altair.sirmuel.design/">Altair GraphQL Client</a>
          </li>
          <li>
            <a href="https://github.com/skevy/graphiql-app">GraphiQL App</a>
          </li>
        </ul>
      </div>

      <ClientOnly>
        <div>
          <div class="overflow-hidden sm:rounded-lg mt-4">
            <div class="py-5 prose">
              <h2 class="text-lg leading-6 text-gray-900">Your API Token</h2>
              <p v-if="fetchedToken" class="mt-1 max-w-2xl">
                Use this personal token in your requests to the Sturdy API. Note that the token is
                self-expiring in 30 days.
              </p>
              <p v-else class="mt-1 max-w-2xl">Login to your Sturdy account to get your token.</p>
            </div>
            <div v-if="fetchedToken">
              <div class="mt-1 flex rounded-md shadow-sm">
                <div class="relative flex items-stretch flex-grow focus-within:z-10">
                  <input
                    id="token"
                    v-model="token"
                    type="text"
                    readonly
                    class="focus:ring-blue-500 focus:border-blue-500 block w-full rounded-none rounded-l-md sm:text-sm border-gray-300"
                  />
                </div>
                <button
                  class="-ml-px relative inline-flex items-center space-x-2 px-4 py-2 border border-gray-300 text-sm font-medium rounded-r-md text-gray-700 bg-gray-50 hover:bg-gray-100 focus:outline-none focus:ring-1 focus:ring-blue-500 focus:border-blue-500"
                  @click="copyToClipboard"
                >
                  <svg
                    class="h-5 w-5 text-gray-400"
                    xmlns="http://www.w3.org/2000/svg"
                    fill="none"
                    viewBox="0 0 24 24"
                    stroke="currentColor"
                  >
                    <path
                      stroke-linecap="round"
                      stroke-linejoin="round"
                      stroke-width="2"
                      d="M8 5H6a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2v-1M8 5a2 2 0 002 2h2a2 2 0 002-2M8 5a2 2 0 012-2h2a2 2 0 012 2m0 0h2a2 2 0 012 2v3m2 4H10m0 0l3-3m-3 3l3 3"
                    />
                  </svg>
                  <span>Copy</span>
                </button>
              </div>
            </div>

            <div>
              <p class="mt-1 max-w-2xl text-sm text-gray-500">
                {{ message }}
              </p>
            </div>
          </div>
        </div>
      </ClientOnly>
    </div>
  </StaticPage>
</template>

<script lang="ts" setup>
import StaticPage from '../../layouts/StaticPage.vue'
import http from '../../http'
import LinkButton from '../../components/shared/LinkButton.vue'
import { ClientOnly } from 'vite-ssr/vue'
import { onMounted, ref } from 'vue'

let token = ref('LOGIN_TO_GET_YOUR_TOKEN')
let fetchedToken = ref(false)
let message = ref(null)

onMounted(async () => {
  try {
    const response = await fetch(http.url('v3/auth/client-token'), {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include',
    })

    http.checkStatus(response)

    const data = await response.json()

    token.value = data.token
    fetchedToken.value = true
  } catch (e) {
    console.error(e)
  }
})

let copyToClipboard = () => {
  let copyText = document.getElementById('token')
  if (copyText) {
    copyText.select()
    copyText.setSelectionRange(0, 99999)
    document.execCommand('copy')
    message.value = 'Copied!'
  }
}
</script>
