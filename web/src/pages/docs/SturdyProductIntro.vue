<template>
  <DocumentationWithTableOfContents>
    <template #sidebar>
      <DocsSidebar />
    </template>
    <template #default>
      <div class="prose p-4 max-w-[800px]">
        <h1 id="product-intro">Product Intro</h1>
        <h2 id="what-is-sturdy">What is Sturdy?</h2>
        <p>Sturdy is an open-source version control platform focused on:</p>
        <ul>
          <li>
            <em>Expressiveness</em> &mdash; actions performed by developers should be directly
            mapped to intent
          </li>
          <li>
            <em>Real-time feedback</em> &mdash; exchange of code ideas in a team should be
            frictionless
          </li>
        </ul>

        <h2 id="what-does-sturdy-do">What does Sturdy do?</h2>
        <h3 id="individual">As an individual</h3>
        <p>Sturdy gives you a highly automated workflow. Here are a few examples:</p>
        <ul>
          <li>
            Get straight to coding &mdash; no need for pushing or fetching changes [<router-link
              :to="{ name: 'v2DocsHowToEditCode' }"
              >docs</router-link
            >]
          </li>
          <li>
            Seamlessly switch between tasks &mdash; no need for stashing or creating temporary
            commits [<router-link :to="{ name: 'v2DocsHowToSwitchBetweenTasks' }">docs</router-link
            >]
          </li>
          <li>
            Develop faster with an tighter "code &mdash; describe changes &mdash; ship" loop
            [<router-link :to="{ name: 'v2DocsHotToShipSoftwareToProduction' }">docs</router-link>]
          </li>
        </ul>

        <h3 id="team">As a team</h3>
        <p>
          Using Sturdy makes
          <em>shipping of small and incremental changes the intuitive default</em>. It is built
          around real-time collaborative features:
        </p>

        <ul>
          <li>
            See each other's code live, as it's being written, opening the opportunity for early
            discussions and feedback [<router-link
              :to="{ name: 'v2DocsHowToCollaborateWithOthers', hash: '#receiving-feedback' }"
              >docs</router-link
            >]
          </li>
          <li>
            Try someone else's code on your computer with a single click [<router-link
              :to="{ name: 'v2DocsHowToCollaborateWithOthers', hash: '#trying-code' }"
              >docs</router-link
            >]
          </li>
          <li>
            Suggest code changes/ideas to a colleague by simply typing in your IDE [<router-link
              :to="{ name: 'v2DocsHowToCollaborateWithOthers', hash: '#give-suggestions' }"
              >docs</router-link
            >]
          </li>
        </ul>

        <h2 id="who-is-sturdy-for">Who is Sturdy for?</h2>
        <p>Sturdy is for:</p>
        <ul>
          <li>
            Developers and teams who are interested in adopting continuous delivery or taking it to
            the next level
          </li>
          <li>
            Engineers who just want to focus on coding and iterating quickly while spending less
            time with Git
          </li>
          <li>
            Experienced developers who review and give feedback to many pull requests in their team
          </li>
          <li>Software teams who want to boost collaboration and knowledge sharing</li>
        </ul>

        <DocsInfoBox class="mt-8">
          You can use Sturdy by yourself and benefit from an optimized workflow even without anybody
          else on your team using it.
        </DocsInfoBox>

        <h2>Can I use Sturdy together with Git?</h2>
        <p>
          <strong>Yes.</strong> Sturdy is built using the same low-level components as git (via
          <a href="https://libgit2.org/">libgit2</a>) and has a
          <router-link :to="{ name: 'v2DocsHowSturdyInteractsWithGit' }"
            >bridge integration
          </router-link>
          with GitHub. This means that you can use Sturdy <strong>on top</strong> of your existing
          GitHub setup, compatible with your existing branch protection and CI/CD and other
          integrations.
        </p>
        <p>
          Sturdy can also import your existing git history, as well as export the history back out
          to git.
        </p>

        <h2 id="how-do-i-get-started">How do I get started?</h2>
        <p>
          Try out
          <router-link :to="{ name: 'v2DocsQuickStart', hash: '#sturdy-cloud' }"
            >Sturdy in the Cloud
          </router-link>
          with zero setup (for free, no credit card required). Or run a
          <router-link :to="{ name: 'v2DocsQuickStart', hash: '#self-hosted-sturdy' }"
            >self-hosted instance
          </router-link>
          of Sturdy anywhere.
        </p>
      </div>
    </template>
  </DocumentationWithTableOfContents>
</template>

<script lang="ts" setup>
import DocumentationWithTableOfContents from '../../layouts/DocumentationWithTableOfContents.vue'
import { useHead } from '@vueuse/head'
import DocsSidebar from '../../organisms/docs/DocsSidebar.vue'
import DocsInfoBox from '../../molecules/DocsInfoBox.vue'

useHead({
  meta: [
    {
      name: 'description',
      content: 'Learn how to start using Sturdy!',
    },
    {
      name: 'keywords',
      content: 'study learn documentation get started tutorial setup quickstart',
    },
  ],
  title: 'Product Intro | Sturdy',
})
</script>
