<template>
  <BlogPost
    title="Launching Sturdy for Desktop!"
    subtitle="This week at Sturdy"
    date="December 7, 2021"
    :author="author"
  >
    <template #introduction>
      Hey! Here’s a summary of all of the latest and greatest features available in the latest
      version of Sturdy.
    </template>
    <template #default>
      <h2>Launching Sturdy for Desktop!</h2>

      <img src="./app.png" width="2878" height="2224" />

      <p>
        We’re now launching the Sturdy app for Mac and Windows! A big thanks to all beta testers and
        early adopters.
      </p>

      <p>
        The app makes it incredibly easy to manage your connected computers, and to connect new
        ones. From now on, you’ll never have to touch the terminal to use Sturdy!
      </p>

      <p><a href="https://getsturdy.com/download?ref=newsletter">Download now!</a></p>

      <h2>Suggestions</h2>

      <p>
        <strong>Code suggestions to workspaces</strong> have gotten a major improvement, mostly
        behind the scenes, but also improvements in the app.
      </p>

      <p>
        When you’re suggesting to another workspace, you’re now suggesting from a “lightweight”
        workspace, that appears in the sidebar underneath the workspace that you’re suggesting to.
      </p>

      <p>
        This allows you to both disconnect and reconnect to your suggestions, and resume the
        suggesting later. You can also have multiple suggestions from different workspaces.
      </p>

      <p>
        And much more! Give it a test by inviting a collaborator to your codebase on Sturdy, and go
        to an active workspace of theirs.
      </p>

      <h2>Selection</h2>

      <p>
        The “selection” of files and hunks on the Workspace page has received a huge facelift! Take
        a look:
      </p>

      <video autoplay loop muted playsinline>
        <source src="./selection.mp4" type="video/mp4" />
        <img src="./selection.gif" />
      </video>

      <h2>Search</h2>

      <p>
        Sturdy now comes with built-in search for code! It’s available as <code>⌘+F</code> in the
        app, and <code>⌘+K</code> on the web.
      </p>

      <video autoplay loop muted playsinline>
        <source src="./search.mp4" type="video/mp4" />
        <img src="./search.gif" />
      </video>

      <h2>Advent of Code</h2>

      <img src="./advent-of-code.jpeg" height="600" width="1200" class="rounded" />

      <p>
        Sturdy is happy to be <strong>sponsoring Advent of Code</strong> this year!
        <a href="https://adventofcode.com/" taget="_blank">Advent of Code</a>
        is a combination of programming puzzle solving and an advent calendar, launching a new
        puzzle every day from the first of December until Christmas.
      </p>

      <p>
        To celebrate, we’ve added support for <em>public codebases</em> on Sturdy, and have made
        <a href="https://getsturdy.com/advent-of-code-2021-uoeIDQk">
          our own Advent of Code codebase
        </a>
        public! We’re doing all of the problems every day (and sometimes even waking up at 06:00 to
        solve them!).
      </p>

      <h2>Sturdy on ARM</h2>

      <p>
        We’re now building and releasing Sturdy for three new ARM architectures on Linux - ARMv5,
        ARMv6, and ARMv7 - in addition to the already existing release for ARM64.
      </p>

      <p>
        This means that you can run Sturdy on your small and lightweight devices such as a Raspberry
        Pi.
      </p>

      <p>
        <a href="https://getsturdy.com/install?ref=newsletter">You can find the downloads here.</a>
      </p>

      <p>Thanks for using Sturdy, see you soon with another blog post!<br />&mdash; Gustav</p>
    </template>
  </BlogPost>
</template>

<script lang="ts" setup>
import BlogPost from '../BlogPost.vue'
import avatar from '../gustav.jpeg'

const author = {
  name: 'Gustav Westling',
  avatar: avatar,
  link: 'https://twitter.com/zegl',
}
</script>
