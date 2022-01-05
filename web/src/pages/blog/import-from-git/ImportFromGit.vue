<template>
  <BlogPost
    title="Importing from Git!"
    subtitle="Migrating to Sturdy"
    date="May 6, 2021"
    :author="author"
  >
    <template #introduction>
      Hey everyone! I hope that you’re doing great, and that May is treating you well!
    </template>
    <template #default>
      <p>
        Things are moving fast at <a href="https://getsturdy.com/">Sturdy</a>, and before diving in
        to the product updates, I’d like to make two announcements!
      </p>
      <p>
        First, we’ve now launched a
        <a href="https://discord.com/invite/5HnSdzMqtA">Discord Community</a>! Please join the
        Discord to chat about Sturdy, give us feedback, and talk about the future of code
        collaboration.
      </p>

      <h2>We're hiring!</h2>

      <p>
        We’re expanding the team here in Stockholm, and are hiring Software Engineers, and an
        engineer with a focus on Developer Relations!
      </p>

      <p>
        This is a great opportunity to be the first employee at a fast moving startup,
        <em>funded by Y Combinator</em>! We’re building tech for techies, and naturally really take
        pride in both the software that we’re writing and the product that we’re building, and hope
        that you do too!
      </p>

      <p>
        Our tech stack is built with Go, AWS, and PostgreSQL on the backend, and TypeScript and Vue
        on the frontend.
      </p>

      <p>
        If this sounds interesting to you or someone you know, please reach out to me on
        gustav@getsturdy.com and we’d love to chat with you!
      </p>

      <hr />

      <p>
        Now, on to the product updates! Here are the biggest updates to Sturdy since issue
        <a href="https://sturdy.substack.com/p/003-share-now">#003</a>, when we launched the new
        “Share” workflows.
      </p>

      <h2>1. Comments!</h2>

      <p>
        First up is comments! You can now write comments on code, both code that you have written as
        well as comments from others. We’re starting small, and intend to build more interesting
        workflows for commenting and collaborating on code, so this is just the first step, and it
        feels great to have it out! 😃
      </p>

      <figure>
        <img class="w-full rounded-lg" src="./comments.gif" alt="" width="786" height="566" />
        <figcaption>Commenting on code in Sturdy</figcaption>
      </figure>

      <h2>2. Import from Git</h2>

      <p>
        This is a huge one! We’ve known from the beginning that it’s important to keep compatibility
        with existing Version Control Systems, and have built Sturdy with this in mind, especially
        with compatibility with Git.
      </p>

      <p>
        You can now import an existing Git repository to Sturdy, without losing any of your data.
      </p>

      <figure>
        <img class="w-full rounded-lg" src="./import.png" alt="" width="786" height="566" />
        <figcaption>Importing a git repository to Sturdy</figcaption>
      </figure>

      <p>
        After creating a new Codebase the following banner will appear on the page. The import is
        done by running <code>“sturdy import”</code> from within a Git repository checkout on your
        computer.
      </p>

      <p>
        Internally, we already support exporting to Git, and we’re soon going to launch it in the
        application as well. Until then, if you want an export of your codebase to Git, please reach
        out to us on Discord!
      </p>

      <h2>3. Windows support</h2>
      <p>
        While we’re on the topic of the Sturdy CLI, since v0.5.10 released just yesterday, we’re
        supporting Sturdy on Windows 10!
      </p>
      <p>
        Sturdy works best in the new <a href="https://aka.ms/terminal">Windows Terminal</a> (because
        it supports emojis 🤩), but also works in PowerShell, the classic Command Prompt, and other
        shells.
      </p>
      <figure>
        <img class="w-full rounded-lg" src="./windows.png" alt="" width="1380" height="712" />
        <figcaption>Sturdy running on Windows</figcaption>
      </figure>

      <h2>4. Shortcuts</h2>

      <p>
        In the last post, we announced the first set of keyboard shortcuts available in Sturdy, and
        we promised to release more keyboard shortcuts. We kept that promise, and here’s the
        shortcuts that have been released since last time.
      </p>

      <ul>
        <li>⌘ + A selects all changes in a workspace</li>
        <li>⌘ + Shift + A de-selects all changes in a workspace</li>
        <li>⌘ + Enter submits a comment</li>
      </ul>

      <p>
        When listing and navigating changes (on the codebase or workspace changelogs), changes can
        now be browsed with browser native navigation
      </p>

      <ul>
        <li>Tab moves the cursor to the next change, and Shift + Tab moves to the previous</li>
        <li>Enter selects the highlighted change</li>
      </ul>
      <p>A full list of shortcuts is available from the “Help” page in Sturdy.</p>

      <h2>5. Other changes</h2>
      <p>
        We’ve made many other smaller changes to improve performance for large workspaces, and
        interactiveness on the site to make it more live, as well as launched fixes for many edge
        cases in the Sturdy CLI.
      </p>

      <p>I’m Gustav and I’ve been your captain this evening, thanks for flying with Sturdy!</p>
    </template>
  </BlogPost>
</template>
<script>
import BlogPost from '../BlogPost.vue'
import avatar from '../gustav.jpeg'

export default {
  components: { BlogPost },
  setup() {
    const author = {
      name: 'Gustav Westling',
      avatar: avatar,
      link: 'https://twitter.com/zegl',
    }
    return { author }
  },
}
</script>
