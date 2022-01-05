<template>
  <StaticPage
    title="Integrations"
    category="solutions"
    metadescription="Integrate Sturdy with other tools and services"
    image="https://images.unsplash.com/photo-1513828742140-ccaa28f3eda0?ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&ixlib=rb-1.2.1&auto=format&fit=crop&w=1350&q=80"
  >
    <div class="text-base max-w-prose mx-auto lg:max-w-none">
      <p class="text-lg text-gray-500">
        A modern software team is using many tools and services, Sturdy comes with a powerful (and
        free) API, and natively integrates with the most popular services.
      </p>
    </div>
    <div
      class="mt-5 prose prose-yellow text-gray-500 mx-auto lg:max-w-none lg:row-start-1 lg:col-start-1"
    >
      <div class="flex space-x-4">
        <img src="../../components/ci/logos/GitHubLogo.png" class="h-8 w-8" />
        <div>
          <h3>GitHub</h3>
          <p>Bring your existing GitHub repositories to Sturdy without downtime</p>
          <p>
            <router-link
              :to="{ name: 'resourcesMigrateFromGitHub' }"
              class="font-medium text-red-600 hover:text-red-500"
            >
              Read more <span aria-hidden="true">&rarr;</span></router-link
            >
          </p>
        </div>
      </div>

      <div class="flex space-x-4">
        <img src="../../components/ci/logos/BuildkiteLogo.svg" class="h-8 w-8" />
        <div>
          <h3>Buildkite</h3>
          <p>Connect your codebase with pipelines on Buildkite.</p>
        </div>
      </div>

      <div class="flex space-x-4">
        <img src="../../components/ci/logos/GraphqlLogo.svg" class="h-8 w-8" />
        <div>
          <h3>Sturdy GraphQL API</h3>
          <p>Connect your applications directly to Sturdy</p>
          <p>
            <router-link
              :to="{ name: 'resourcesApi' }"
              class="font-medium text-red-600 hover:text-red-500"
            >
              Get started <span aria-hidden="true">&rarr;</span></router-link
            >
          </p>
        </div>
      </div>

      <div class="flex space-x-4">
        <img src="../../components/ci/logos/GitLogo.svg" class="h-8 w-8" />
        <div>
          <h3>Git</h3>
          <p>Import your existing git repository to Sturdy with <code>sturdy import</code>.</p>
        </div>
      </div>

      <h3>More coming soon&hellip;</h3>
      <p>
        We want to integrate with the complete suite of productivity tools that you and your team
        are using today. Send your wishlist to
        <a href="mailto:support@getsturdy.com">support@getsturdy.com</a> and we'll make it happen!
      </p>
    </div>
  </StaticPage>
</template>

<script>
import StaticPage from './StaticPage.vue'
export default {
  components: { StaticPage },
}
</script>
