<template>
  <StaticPage
    title="Continuous Integration/Delivery (CI/CD)"
    metadescription="Learn how Sturdy integrates with your CI/CD workflows"
    image="https://images.unsplash.com/photo-1629904853716-f0bc54eea481?ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&ixlib=rb-1.2.1&auto=format&fit=crop&w=3540&q=80"
  >
    <div class="prose prose-yellow">
      <p>
        <router-link :to="{ name: 'resourcesMigrateFromGitHub' }"
          >Sturdy integrates with GitHub
        </router-link>
        . Any Continuous Integration/Continuous Delivery providers you use that report to GitHub
        will also have its status be reported in the Sturdy UI.
      </p>

      <h2>Running CI/CD Pipelines</h2>
      <p>
        If you have your codebase
        <router-link :to="{ name: 'resourcesMigrateFromGitHub' }"
          >hooked up with GitHub
        </router-link>
        , Sturdy will occasionally export changes to the corresponding repository on GitHub,
        depending on the integration/migration phase you're using.
      </p>

      <ul>
        <li>
          If you're using
          <router-link
            :to="{ name: 'resourcesMigrateFromGitHub', hash: '#github-as-source-of-truth' }"
            >GitHub as the source of truth
          </router-link>
          , changes you're making in Sturdy will be upstreamed to GitHub via pull requests. Every
          time you create or update a PR via Sturdy, the checks that are configured on the GitHub
          side will be run.
        </li>

        <li>
          If you're using
          <router-link
            :to="{ name: 'resourcesMigrateFromGitHub', hash: '#sturdy-as-source-of-truth' }"
            >Sturdy as the source of truth
          </router-link>
          , changes that are shared from a draft change and added to the changelog will be committed
          to the default branch on the GitHub repository. This, in turn, will run any workflows
          configured in GitHub.
        </li>
      </ul>

      <h2>Feedback in Sturdy</h2>
      <p>
        Whichever way you do it, Sturdy will report the status of the CI/CD runs directly in the UI.
        On all pages where entries from the changelog are shown, the status of the pipeline for each
        change will be added.
      </p>

      <figure>
        <img
          alt="Sturdy adds status icons to all the latest changes displayed on the top of the codebase page"
          src="./status-on-codebase.png"
          width="1622"
          height="652"
        />
        <figcaption>
          On the top of the codebase page, the list of latest changes will be complemented with a
          status icon indicating the result of any CI/CD workflow.
        </figcaption>
      </figure>
    </div>
  </StaticPage>
</template>

<script lang="ts" setup>
import StaticPage from '../../../layouts/StaticPage.vue'
</script>
