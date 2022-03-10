<template>
  <BlogPost
    title="Introducing Draft changes"
    subtitle="This week at Sturdy"
    date="March 10, 2022"
    :author="author"
    :image="ogImageFull"
  >
    <template #introduction>
      Hey everyone! In this week's update I would like to share with you an update to the Sturdy
      workflow with Draft changes, as well as the release of Sturdy for Linux. Today is also an
      exciting day for us because we are launching Sturdy on
      <a class="text-amber-600" href="https://www.producthunt.com/posts/sturdy">Product Hunt</a>!
    </template>
    <template #default>
      <h2>📝 Draft changes</h2>
      <p>
        We felt that workspaces in Sturdy were sometimes unnecessarily complex. For example, after
        you merge the code in your workspace, should you continue working in the same workspace or
        create a new one?
      </p>
      <p>Drafts replace Workspaces. Here's how they work:</p>
      <ul>
        <li>All code changes start from a Draft</li>
        <li>
          If a Draft is empty, it's name will be <span class="font-semibold italic">italic</span>,
          and it will not be visible to others
        </li>
        <li>Drafts with code updates are discoverable in the app sidebar</li>
        <li>Whenever a Draft is merged, it simply becomes a Change entry in the Changelog</li>
      </ul>
      <div class="flex justify-center items-center">
        <img src="./drafts-italic.png" width="250" class="rounded shadow-md" />
      </div>

      <p>
        This change makes the 'develop &mdash; collaborate &mdash; integrate' loop even tighter. In
        addition, as a team you will no longer have to search which workspaces have relevant
        changes.
      </p>
      <p>
        We would love to hear what you think! Feel free to drop us a message on
        <a href="https://discord.gg/fQcH9QAVpX">Discord</a>.
      </p>

      <h2>🚀 Launching on Product Hunt</h2>
      <p>
        It is always very exciting to share Sturdy with the wider developer community. It also
        drives us to think how to better describe what we stand for.
      </p>
      <p>
        Today we are launching
        <a href="https://www.producthunt.com/posts/sturdy">Sturdy on Product Hunt</a>! Check out the
        workflow video and screenshots that we shared :). We would also be very thankful for your
        upvote there 🙏❤️!
      </p>

      <img src="./skip-the-git-overhead.png" height="600" width="1200" class="rounded shadow-md" />

      <h2>🐧 Sturdy on Linux</h2>

      <p>
        The Sturdy desktop client 🐥 now has a Linux 🐧 version! We had a number of users asking for
        it, and now we have finally added <strong>.deb</strong> and
        <strong>.rpm</strong> distributions on our
        <a href="https://getsturdy.com/download">download page</a>.
      </p>

      <img src="./sturdy-linux.jpeg" height="600" width="1200" class="rounded shadow-md" />

      <p>This means that Sturdy is now available for Mac, Windows, and Linux!</p>

      <h2>🎙 Kodsnack podcast (in Swedish)</h2>
      <p>
        I had the pleasure of talking about Sturdy, version control, and our vision for code
        collaboration on the Kodsnack podcast.
      </p>
      <p>
        Fredrik had some great questions and the conversation was a lot of fun! It is in Swedish,
        and you can check it out <a href="https://kodsnack.se/463/">here</a>.
      </p>

      <img src="./kodsnack-logo.png" width="300" class="rounded shadow-md" />

      <br />
      <p>
        Thanks for reading, and until the next post!
        <br />&mdash; Gustav and team Sturdy
      </p>
    </template>
  </BlogPost>
</template>

<script lang="ts" setup>
import BlogPost from '../BlogPost.vue'
import avatar from '../gustav.jpeg'
import ogImage from './skip-the-git-overhead.png'

import { computed } from 'vue'

let ogImageFull = computed(() => new URL(ogImage, import.meta.env.VITE_WEB_HOST).href)

const author = {
  name: 'Gustav Westling',
  avatar: avatar,
  link: 'https://twitter.com/zegl',
}
</script>
