<template>
  <StaticPage
    title="Migrate from GitHub to Sturdy"
    category="knowledgebase"
    metadescription="How to migrate from GitHub to Sturdy, and gain real-time superpowers"
    image="https://images.unsplash.com/photo-1580894908361-967195033215?ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&ixlib=rb-1.2.1&auto=format&fit=crop&w=1650&q=80"
  >
    <div class="text-base max-w-prose mx-auto lg:max-w-none">
      <p class="text-lg text-gray-500">
        Learn how to migrate from GitHub to Sturdy without any downtime, or losing any work.
      </p>
    </div>
    <div
      class="mt-5 prose prose-yellow text-gray-500 mx-auto lg:max-w-none lg:row-start-1 lg:col-start-1"
    >
      <p>
        We know &mdash; migrating from one Version Control System to another without any team
        "downtime" sounds like a nightmare.
      </p>
      <p>
        This is why we've built a migration path where you and your team slowly can migrate from
        GitHub to Sturdy without losing any existing work, or risking breaking your CI/CD pipeline.
      </p>
      <p>
        <LinkButton href="https://getsturdy.com/user" class="!no-underline !text-black">
          Connect GitHub to Sturdy
        </LinkButton>
      </p>
      <p>
        A migration has three phases, that start with respecting GitHub as it's current
        <i>source of truth</i>.
      </p>
      <h3 id="github-as-source-of-truth">GitHub as the source of truth</h3>
      <p>
        When connecting a GitHub repository to Sturdy the first time, nothing changes on your GitHub
        repository. You can use Sturdy and GitHub side-by-side for as long as you need to, and you
        don't have to switch to using Sturdy as the source of truth until you and your team are
        ready.
      </p>

      <div class="relative text-base mx-auto max-w-prose lg:max-w-none">
        <figure>
          <div class="lg:aspect-w-12 aspect-h-7 lg:aspect-none">
            <img
              class="rounded-lg object-cover object-center"
              alt="A Pull Request created by Sturdy"
              src="../../assets/Web/GitHubPullRequest.png"
              height="498"
              width="937"
            />
          </div>
          <figcaption class="mt-3 flex text-sm text-gray-500">
            <span class="ml-2">Creating GitHub pull requests from Sturdy</span>
          </figcaption>
        </figure>
      </div>

      <ul>
        <li>Sturdy syncs new commits when they are pushed to the <i>default</i> GitHub branch</li>
        <li>Use Sturdy draft changes, with live feedback and suggestions</li>
        <li>
          When a draft change is ready to be shared, Sturdy opens a pull request towards your repo.
          Where you can use your existing code review and CI/CD tools.
        </li>
        <li>When the pull request is merged on GitHub, it will be synced back to Sturdy again</li>
      </ul>
      <h3 id="sturdy-as-source-of-truth">Sturdy as the source of truth</h3>
      <p>
        To enable the full workflow on Sturdy, you have to make Sturdy the source of truth. In this
        mode you're sharing changes between each other only using Sturdy, and don't have to work
        with GitHub pull requests anymore!
      </p>
      <ul>
        <li>Changes are shared directly on Sturdy</li>
        <li>Use Sturdy draft changes, with live feedback and suggestions</li>
        <li>
          After a change has been <i>shared</i> on Sturdy, a copy can be pushed to your GitHub
          repository, as a backup, or to trigger your release pipelines.
        </li>
        <li>
          Branches and commits merged on GitHub without Sturdy will <i>not</i> be synced to Sturdy.
        </li>
      </ul>
      <h3>Disconnect from GitHub</h3>
      <p>
        Last but not least, when you don't need GitHub anymore, disconnect the Codebase from GitHub.
        <strong>Congrats, you're now fully migrated to Sturdy!</strong>
      </p>
      <h3>Permissions required by the Sturdy GitHub-app</h3>
      <Permissions />
      <h3>Any questions or problems?</h3>
      <p>
        If anything goes wrong, or you have any questions about migrating from GitHub to Sturdy.
        Please email us at <a href="mailto:support@getsturdy.com">support@getsturdy.com</a>.
      </p>
    </div>
  </StaticPage>
</template>

<script lang="ts" setup>
import LinkButton from '../../atoms/LinkButton.vue'
import StaticPage from '../../layouts/StaticPage.vue'
import Permissions from '../../components/github/Permissions.vue'
</script>
