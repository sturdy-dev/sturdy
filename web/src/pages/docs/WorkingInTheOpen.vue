<template>
  <DocumentationWithTableOfContents>
    <template #sidebar>
      <DocsSidebar />
    </template>
    <template #default>
      <div class="prose p-4 max-w-[800px]">
        <h1 id="working-in-the-open">Working in the open (within your team)</h1>

        <p>
          As a developer in a software team, you likely already do stand-ups, and planning together
          of what to work on next. Some teams also do regular demos and pairing sessions to boost
          knowledge sharing which really pays off in the long run. In this context, working in the
          open simply means continuously communicating within the team.
        </p>

        <p>
          Regardless of which processes a team has chosen, the main point is that the team succeeds
          or fails together. This is why most highly effective teams
          <a href="https://rework.withgoogle.com/print/guides/5721312655835136/"
            >excel at collaboration</a
          >.
        </p>

        <p>
          Sturdy is not a "process" tool. It does, however, follow the following two principles:
        </p>

        <ul>
          <li>Reducing the distance between developers</li>
          <li>Normalizing "trying things" and experimentation</li>
        </ul>

        <p>
          When it comes to building software as a team, early feedback is a good thing. A second
          pair of eyes gives an extra perspective which may spot a design flaw or an opportunity
          before a lot of effort was invested in a particular solution. Similarly, even just giving
          each other encouragement on the progress is motivating and helps keep us going. The later
          this feedback comes, the more it is perceived as judgement, and the less effective it is
          at producing better software.
        </p>

        <p>
          All of this is, of course, easier said than done. Traditionally, work starts off on a
          local branch and more often than not we want to "finish" it before sharing with others. As
          creators we have the natural tendency to want to perfect things. But perfectionism is
          holding us back from collaborating on ideas while they are small and moldable. It also
          prevents the compounding of ideas.
        </p>

        <p>
          <strong
            >Sturdy aims to establish a more relaxed vibe to working on code as a team.</strong
          >
        </p>

        <p>
          Instead of feeling pressure to perfect things before sharing with others for "review",
          teams using Sturdy celebrate trying things, creating experiments, and are building on each
          other's ideas.
        </p>

        <h2 id="reducing-distance-between-developers">Reducing distance between developers</h2>

        <p>
          Whenever somebody in a team works on something in isolation for a while, this creates
          distance between them and the team. The divergence comes not only from the new,
          unintegrated code but also from changes introduced by others on the team.
        </p>

        <p>
          To work effectively, it is necessary that a team is aligned (pulling in the same
          direction) and has a level of trust to be able to provide meaningful feedback and help to
          one another.
        </p>

        <p>
          While creating alignment and trust is something teams develop themselves, there is a case
          to be made that the tools that we use shape our behaviour.
        </p>

        <p>
          For example practices like
          <a href="https://cloud.google.com/architecture/devops/devops-tech-continuous-delivery"
            >Continuous Delivery</a
          >
          and
          <a href="https://cloud.google.com/architecture/devops/devops-tech-trunk-based-development"
            >Trunk based development</a
          >
          are known to yield not just faster development but a higher quality of the software.
        </p>

        <p>
          Sturdy is built around Trunk Based Development and Continuous Delivery while also making
          giving meaningful feedback early an easy and lightweight operation. It should be as easy
          and intuitive to exchange code ideas as it is to write messages in a Slack channel.
        </p>

        <p>
          Sturdy does all of this through it's <strong>real-time design</strong>. Here are a few
          examples:
        </p>

        <ul>
          <li>The team can easily discover code as it is being written by others</li>
          <li>You can get someone's code to your IDE with just one click</li>
          <li>
            Giving code suggestions to a team member is as simple as typing in your IDE while you
            have their code
          </li>
        </ul>

        <h2 id="normalizing-experimentation">Normalizing experimentation</h2>

        <p>
          Experimentation is very much core to software development, yet it is not prominently
          reflected in the process of many teams. The "failed" attempts and false starts before a
          feature is completed are rarely discussed. This is a pity because much of the learning
          happens right there.
        </p>

        <p>
          Sturdy normalizes experimentation and just "trying things". It's our core belief that this
          is an important part for not just knowledge sharing but also building trust.
        </p>

        <p>
          The way experimentation is supported in Sturdy is in the design of it's lightweight draft
          changes. If a developer has an idea or wants to quickly test something, they can quickly
          jump to a new draft. This is similar to git branches, but there is no need for manually
          stashing of current work. But also critically, because drafts are available to the team in
          real time, experiments like this can be part of the larger team conversation.
        </p>
      </div>
    </template>
  </DocumentationWithTableOfContents>
</template>

<script lang="ts" setup>
import DocumentationWithTableOfContents from '../../layouts/DocumentationWithTableOfContents.vue'
import { useHead } from '@vueuse/head'
import DocsSidebar from '../../organisms/docs/DocsSidebar.vue'

useHead({
  meta: [
    {
      name: 'description',
      content: 'Working in the open (within your team)',
    },
    {
      name: 'keywords',
      content: 'study learn documentation working code in the open publicly public',
    },
  ],
  title: 'Working in the open | Sturdy',
})
</script>
