<template>
  <DocumentationWithTableOfContents>
    <template #sidebar>
      <DocsSidebar />
    </template>
    <template #default>
      <div class="prose p-4 max-w-[800px]">
        <h1 id="sturdy-documentation">Sturdy Documentation</h1>

        <p>
          Sturdy is an open-source version control platform that allows you to interact with your
          code at a higher abstraction level. This provides an ergonomic, real-time workflow for
          individuals and teams.
        </p>

        <p>
          In Sturdy, developers perform
          <strong>actions that are directly mapped to their intent</strong>, rather than thinking
          about the underlying data structures. For example: "Ship this code to production" instead
          of "Create a branch, stage files, make a commit, push to remote, merge, etc.".
        </p>

        <p>
          The dynamics of teams at tech companies are different from those of most open source
          projects. The biggest difference is in how closely people collaborate with each other
          &mdash; at a software shop you already synchronize, share knowledge and plan together.
          Sturdy is built specifically for this workflow.
        </p>
      </div>

      <div class="my-8 bg-gray-200 rounded-md px-4 py-2">
        <div class="prose max-w-full">
          <h2 id="get-started">Get started</h2>

          There are three versions of Sturdy:

          <ul>
            <li>
              <a href="https://getsturdy.com/"><strong>Sturdy in the Cloud</strong></a> lets you
              code and collaborate in real-time with your team.
            </li>
            <li>
              <strong>Sturdy Enterprise</strong> lets you run your own Sturdy instance in your own
              environment.
            </li>
            <li>
              <strong>Sturdy OSS</strong> is an open-source version of Sturdy that gives you all of
              the core functionality of Sturdy, but with some advanced features missing (OAuth,
              incremental migrations to Sturdy, etc).
            </li>
          </ul>

          <h3>Want to take Sturdy for a spin?</h3>

          Get started on <a href="https://getsturdy.com/"><strong>Sturdy Cloud</strong></a
          >, or run in Docker:

          <pre>{{ dockerOneliner }}</pre>
        </div>
      </div>

      <div class="prose max-w-[800px] p-4">
        <p>
          Sturdy is built on top of low-level Git data structures, so your code is always stored in
          a well-compatible format, and allows you to easily migrate your team to Sturdy, even
          partially if you want.
        </p>

        <p>
          You can use Sturdy <strong>by yourself</strong> and benefit from an optimized workflow
          even without anybody else on your team using it. Sturdy has a bridge integration with
          GitHub that allows you to interact with your existing projects.
        </p>

        <p>
          On the other hand, as more people from <strong>your team</strong> use Sturdy, it also
          becomes possible to benefit from a frictionless exchange of ideas &mdash; give code
          suggestions, early feedback.
        </p>

        <h2 id="collaboration">Collaboration</h2>

        <p>Collaboration in Sturdy is built around two core principles:</p>
        <ul>
          <li>Continuous feedback over traditional code review</li>
          <li>Small incremental changes over big-bang releases</li>
        </ul>

        <p>
          <router-link :to="{ name: 'v2DocsWorkingInTheOpen' }"
            ><strong>Continuous feedback</strong></router-link
          >
          reflects the reality of working together &mdash; getting help and suggestions on an early
          version of the code is more effective than review of a "completed" work. This means
          creating an environment which fosters experimentation and working in the open.
        </p>

        <p>
          <strong>Shipping small incremental changes</strong> is
          <a href="https://cloud.google.com/devops/state-of-devops">proven</a> to increase speed of
          development as well and reduce the risk of deployments. Sturdy makes shipping of small and
          incremental changes the intuitive default. It embraces
          <a href="https://cloud.google.com/architecture/devops/devops-tech-continuous-delivery"
            >continuous delivery</a
          >
          and
          <a href="https://cloud.google.com/architecture/devops/devops-tech-trunk-based-development"
            >trunk based development</a
          >.
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
  components: { DocumentationWithTableOfContents, DocsSidebar },
  setup() {
    useHead({
      meta: [
        {
          name: 'description',
          content:
            'Learn what Sturdy is, and how you can use Sturdy in the Cloud or self-host Sturdy on your own premises.',
        },
        {
          name: 'keywords',
          content:
            'study learn documentation cloud self host self-host docker oneliner on prem on-prem premises',
        },
      ],
      title: 'Documentation | Sturdy',
    })

    return {
      dockerOneliner: DOCKER_ONELINER,
    }
  },
})
</script>
