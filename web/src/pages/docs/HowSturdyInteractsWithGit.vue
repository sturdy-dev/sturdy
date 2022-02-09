<template>
  <DocumentationWithTableOfContents>
    <template #sidebar>
      <DocsSidebar />
    </template>
    <template #default>
      <div class="prose p-4 max-w-[800px]">
        <h1 id="how-sturdy-interacts-with-git">Interacting with Git</h1>

        <p>
          Sturdy utilizes low-level Git data structures, but what does that mean, and why? In a
          nutshell, it is all about creating a more leveraged development environment &mdash; a
          productivity hack.
        </p>
        <p>
          Git is powerful as a <em>database</em>, but as a <em>protocol</em>, it places constraints
          on the development workflow. We believe that these constraints don't fit modern product
          development which is fast-paced and highly collaborative.
        </p>

        <h2 id="expressiveness">Expressiveness</h2>
        <p>
          The dictionary definition of the term "expressiveness" is "effectively conveying a
          thought". In the context of software development, our goal is to map developer intent to
          concrete actions in Sturdy. The command line is awesome because it is fast and efficient
          for getting things done, but even faster is
          <em>not doing those things in the first place!</em>
        </p>

        <p>Here is a typical cycle of code contribution with Git:</p>

        <ol>
          <li>Pull main or master</li>
          <li>Create a new branch</li>
          <li>Write the code</li>
          <li>Add the changes</li>
          <li>Commit</li>
          <li>Push the branch</li>
          <li>Create a pull request</li>
          <li>Receive feedback</li>
          <li>GOTO 3.</li>
        </ol>

        <p>
          Regardless of whether you use the terminal or a Git GUI, these steps are part of the
          <em>protocol</em>. If we however distill the developer intent, it was to solve a business
          problem, together with a team.
        </p>

        <p>
          Sturdy is built by <em>truly lazy</em> developers who thought &mdash; “can’t we just do
          the coding and collaboration, and automate the rest of the steps?”.
        </p>

        <h2 id="abstractions">Abstractions</h2>

        <p>
          Generally, computers are built on many layers of abstraction. Those abstractions exist
          because they provide leverage when building new things. For example, programming languages
          exist on a spectrum of control / expressiveness.
        </p>

        <p>
          Lower level languages allow for precise control &mdash; e.g. choosing specific CPU
          instructions (assembly) or manually allocating and freeing memory (C, C++). The tradeoff
          is that more code needs to be written to achieve the goal. For some, this is the
          appropriate abstraction level.
        </p>

        <p>
          Higher level languages (e.g. Java, Go, Python) abstract away things like memory management
          and allow the coder to focus at a higher level (e.g. a specific business problem). Under
          the hood there is of course still a stack and heap, but a Python programmer doesn’t need
          to think about it. For many, this is the appropriate abstraction level.
        </p>

        <p>
          In the same spirit, Sturdy is a higher level tool, an abstraction allowing coders to focus
          on building software and exchanging feedback and ideas with their team. There are still
          branches, commits, checkouts, trees and refs, but they are fully managed by Sturdy.
        </p>

        <p>
          This is why Sturdy does not have a CLI – the simplest interface is the absence of one. The
          only interface that Sturdy has revolves around collaboration and initiating higher level
          operations.
        </p>

        <p>
          As to the big question &mdash; why? Increasing the expressiveness of the version control
          and collaboration means getting more stuff done, or just finishing quicker, depending on
          teams priorities. Sturdy streamlines Continuous Delivery &mdash; shipping small
          incremental changes, frequently, which are easier to review and reason around.
        </p>

        <h2 id="there-is-no-local-or-remote" title="There is no local or remote">
          There is no <span class="line-through">spoon</span> local or remote
        </h2>

        <p>
          The big, glaring difference between Sturdy and Git is that there is no distinction between
          local and remote. Just like Git, code is read and modified through a folder on a computer,
          using an IDE or text editor. However, changes that are made with Sturdy are instantly
          discoverable by others on the team (through Sturdy Cloud or a self-hosted Sturdy).
        </p>

        <p>
          The implication is that a coder no longer has to run Git commands to manage state in a
          local directory. In Sturdy all the operations are at a high level (compared to using
          vanilla Git where several steps are needed to achieve the equivalent outcome).
        </p>

        <p>
          All the Sturdy tricks are enabled by the real-time synchronization of files and state, and
          the fact that code management, review‚ and feedback happen in the same place:
        </p>

        <ol>
          <li>Discover work-in-progress code within the team, get early feedback.</li>
          <li>No manual synchronizing of local and remote.</li>
          <li>
            Getting the code from someone else’s workspace (roughly equivalent to a branch or pull
            request) with a single click.
          </li>
          <li>Trivial to give or accept code suggestions.</li>
          <li>Automatically stay up to day with the main / default branch.</li>
        </ol>

        <h2 id="how-git-and-sturdy-concepts-relate">How Git and Sturdy concepts relate</h2>
        <p>
          Having mentioned that Sturdy operates at a higher level, it is worth stressing that the
          output is fully compatible with Git. The underlying primitives are the same and a few
          higher level concepts are introduced. Here is a TL;DR:
        </p>

        <table>
          <thead>
            <tr>
              <th>Sturdy</th>
              <th>Git</th>
            </tr>
          </thead>
          <tbody>
            <tr>
              <td>Codebase</td>
              <td>Repository</td>
            </tr>
            <tr>
              <td>Workspace</td>
              <td>Branch and pull request (<em>but live</em>)</td>
            </tr>
            <tr>
              <td>Connected directory</td>
              <td>Checkout / clone of a repository</td>
            </tr>
            <tr>
              <td><em>N/A</em></td>
              <td>Pull</td>
            </tr>
            <tr>
              <td><em>N/A</em></td>
              <td>Push</td>
            </tr>
          </tbody>
        </table>

        <p>Pull and Push are not user facing operations in Sturdy.</p>

        <h2 id="further-reading">Further reading</h2>
        <p>
          You can learn more about how Sturdy interoperates with GitHub on
          <router-link :to="{ name: 'v2DocsHowToSetupSturdyWithGitHub' }">this page</router-link>
          .
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
      content: 'How Sturdy interacts with git.',
    },
    {
      name: 'keywords',
      content: 'study learn documentation interacting git github',
    },
  ],
  title: 'Interacting with Git | Sturdy',
})
</script>
