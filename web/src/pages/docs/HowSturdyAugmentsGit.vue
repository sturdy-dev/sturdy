<template>
  <PublicLeftSidebar>
    <template #sidebar>
      <DocsSidebar />
    </template>
    <template #default>
      <div class="prose p-4 max-w-[800px]">
        <h1 id="how-sturdy-augments-git">How Sturdy augments Git</h1>

        <p>
          Sturdy can run on top of Git and GitHub, but what does that mean? In a nutshell, it is all
          about creating a more leveraged development environment &mdash; a productivity hack.
        </p>

        <h2 id="expressiveness">Expressiveness</h2>
        <p>
          The command line is awesome because it is fast and efficient for getting things done, but
          it’s not as fast as not to do anything!
        </p>

        <p>The typical cycle of code contribution with git is something like this:</p>

        <ol>
          <li>Pull main or master</li>
          <li>Create a new branch</li>
          <li>Write the code</li>
          <li>Add the changes</li>
          <li>Commit</li>
          <li>Push the branch</li>
          <li>Collaborate and exchange feedback with the team (through a PR or MR)</li>
          <li>Got to 3.</li>
        </ol>

        <p>
          This can be done with a Git GUI or the Git CLI, except for the feedback exchange, which
          happens in the browser. Sturdy is built by <em>truly lazy</em> developers that thought
          &mdash; “can’t we just do the coding and collaboration, and automate the rest of the
          steps?”.
        </p>

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
          teams priorities. Sturdys jam is Continuous Delivery &mdash; shipping small incremental
          changes, frequently, which are easier to review and reason around.
        </p>

        <h2 id="there-is-no-local-or-remote" title="There is no local or remote">
          There is no <span class="line-through">spoon</span> local or remote
        </h2>

        <p>
          The big, glaring difference between Sturdy and Git is that there is no distinction between
          local and remote. Just like git, code is read and modified through a folder on a computer,
          using an IDE or text editor. However, changes that are made with Sturdy are instantly
          available for others to review (through Sturdy Cloud or a self-hosted Sturdy).
        </p>

        <p>
          The implication is that a coder no longer has to perform git commands to manage state in a
          local directory. In Sturdy all the operations are at a high level (compared to using
          vanilla Git where several steps are needed to achieve the equivalent outcome).
        </p>

        <p>
          All the Sturdy tricks are enabled by the real-time synchronization of files and state and
          the fact that code management, review and feedback happen in the same place:
        </p>

        <ol>
          <li>Discover work in progress code within the team, get early feedback.</li>
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
          higher level concepts are introduced. Here is a TL:DR;
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

        <h2 id="implementation">Implementation</h2>

        <p>
          When Sturdy is configured with a git-bridge (such as "Sturdy for GitHub"), the codebase in
          Sturdy runs on top of the existing git repository.
        </p>

        <p>
          Sturdy will clone the repository to the Sturdy backend, and import all changes and the
          full history. Webhooks from GitHub keep the trunk in Sturdy up-to-date with the HEAD
          branch automatically, and continuously imports more changes to Sturdy.
        </p>

        <p>
          Workspaces in Sturdy enables collaboration, suggestions, and giving and taking feedback.
          When Sturdy runs on top of GitHub, to share code, a Pull Request will be opened towards
          the repository. This PR is subject to the branch protection rules configured on GitHub,
          and if everything is green, it can be merged.
        </p>

        <p>
          Sturdy keeps track of PRs created from Workspaces, and the PR can be updated with more
          changes, or merged directly from Sturdy.
        </p>

        <p>
          <em>Known limitations</em>: git-submodules and LFS files can currently not be imported to
          Sturdy (the metadata will be there, but not the "full" file).
        </p>
      </div>
    </template>
  </PublicLeftSidebar>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import PublicLeftSidebar from '../../layouts/PublicLeftSidebar.vue'
import { useHead } from '@vueuse/head'
import DocsSidebar from '../../organisms/docs/DocsSidebar.vue'

export default defineComponent({
  components: { DocsSidebar, PublicLeftSidebar },
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
          content: 'How Sturdy is augmenting and building on top of git.',
        },
        {
          name: 'keywords',
          content: 'study learn documentation augmenting git github',
        },
      ],
      title: 'How Sturdy Augments Git | Sturdy',
    })
  },
})
</script>
