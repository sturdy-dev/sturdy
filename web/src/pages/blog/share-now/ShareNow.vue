<template>
  <BlogPost
    title="Share now!"
    subtitle="Sharing changes with your team on Sturdy"
    date="March 16, 2021"
    :author="author"
  >
    <template #introduction>
      Hey coders! As you’ve noticed, we’re still exploring the fundamentals to the core workflow of
      Sturdy. This week we’ve made some larger changes to when changes are copied from a
      workspace<sup id="back-1"><a href="#note-1">1</a></sup> to a codebase<sup id="back-2"
        ><a href="#note-2">2</a></sup
      >. Let’s get into it!
    </template>
    <template #default>
      <h2>Sharing a change with the team</h2>
      <p>
        Earlier in Sturdy, making a change and sharing it with your team was a two step process.
        First you had to “Create change” where a change would be staged in the workspaces own
        changelog, but not yet available in other workspaces. To move a change to the codebase, the
        change first had to be found in the changelog and “Landed”.
      </p>

      <p>
        We didn’t like this workflow, as the actions that you had to take did not match the users
        intention of “I’m done coding this now, and want to share it with my team”. And honestly, it
        was just confusing. What does even “Land” mean? 🤷‍♂️
      </p>

      <p>
        We want to build something that’s better than sharing and reviewing changes late and in
        large batches, which is common in many Git based workflows. Git in itself does not have any
        limitations on how large or small a commit must be, but it’s often not ergonomic or
        effective to make small and independent changes, as creating even a single PR often requires
        many commands and switching between the terminal (or GUI) and the web.
      </p>

      <p>
        In Sturdy, sharing a change with the codebase and your team, is now one click (or two
        keystrokes) away! “Share now” takes your changes, and moves it to the codebase right away.
      </p>

      <figure class="text-center">
        <img class="w-max-[100%] rounded-lg block m-auto" src="./001.png" alt="" width="300" />
      </figure>

      <p>
        “Add to share later” prepares the change and places it in the workspaces own changelog,
        where you can change your mind about it, edit it, before deciding to share it.
      </p>

      <figure class="text-center">
        <img class="w-max-[100%] rounded-lg block m-auto" src="./002.png" alt="" width="300" />
      </figure>

      <p>
        I’ve been really happy with this improvement to the ease of use! Sturdy really is the
        easiest way for developers to collaborate on code.
      </p>

      <p>Now, next topic:</p>

      <h2>Keyboard shortcuts</h2>

      <p>
        Graphical UIs are great, they allow you to discover the product and features that you didn’t
        knew existed, but when it comes to productivity, nothing can beat keyboard shortcuts.
      </p>

      <p>
        This week we’ve started to add shortcuts to Sturdy, and have added three of them so far, as
        we’ve hinted in the screenshots above.
      </p>

      <p>
        When creating a change, the change can be shared immediately with <code>⌘⏎</code>. To share
        the change later use <code>⌘⇧⏎</code>. If you make a change in error, the change can be
        undone with <code>⌘Z</code>.
      </p>

      <p>
        This is a good start, but we’re not done adding shortcuts just yet, there will be more
        shortcuts coming soon, especially to the workflow of selecting diffs to include in a change.
      </p>

      <h2>Change the description of a change</h2>
      <p>
        A changes description is now editable! 🤯 If you make a typo, or just want to add some more
        context to a change that you or a team mate has made, you can now edit it, even if it’s been
        shared with the codebase! I bet you Git can’t do this. ;-)
      </p>

      <h2>Ignoring changes to a file</h2>
      <p>
        One of the first features that we built in Sturdy, before we even bootstrapped<sup
          id="back-3"
          ><a href="#note-3">3</a></sup
        >
        the project, was compatibility with .gitignore. It was a great feature, but was a hidden
        gem.
      </p>

      <p>
        To make it easier to ignore files, and to make it easier to discover this piece of git
        compatibility, we’ve now added a “Ignore file” button in the file dropdown when making a
        change. Ignoring a file adds this file to the codebases root .gitignore file in one click.
      </p>

      <p>
        It’s currently not possible to mark an entire directory, such as node_modules, as ignored
        through the web UI — but manually editing the .gitignore file works great now when you know
        that it’s possible to do so. :-)
      </p>

      <figure class="text-center">
        <img
          class="w-max-[100%] rounded-lg block m-auto"
          src="./003.png"
          alt=""
          width="554"
          height="398"
        />
      </figure>

      <p>
        Those are the highlights of this week. Have a nice weekend everyone, and we’ll see you again
        in next weeks post!
      </p>

      <hr />

      <ol>
        <li id="note-1">
          [Git: The equivalent of a unique branch and checkout]

          <a href="#back-1" class="ml-2">↩</a>
        </li>
        <li id="note-2">
          [Git: A repository and it’s HEAD branch, that all workspaces are tracking]
          <a href="#back-2" class="ml-2">↩</a>
        </li>
        <li id="note-3">
          Our first goal with Sturdy was to start developing Sturdy in Sturdy
          <a href="#back-3" class="ml-2">↩</a>
        </li>
      </ol>
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
