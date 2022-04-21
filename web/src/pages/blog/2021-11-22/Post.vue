<template>
  <BlogPost
    title="Sturdy-the-app is coming!"
    subtitle="This week at Sturdy"
    date="November 22, 2021"
    description="Hey! Here’s a summary of all of the latest and greatest features available in the latest
      version of Sturdy."
    :author="author"
  >
    <template #introduction>
      Hey! Here’s a summary of all of the latest and greatest features available in the latest
      version of Sturdy.
    </template>
    <template #default>
      <h2>Native App (coming soon!)</h2>

      <img src="./codebaseHome.png" height="1934" width="2774" />

      <p>
        Sturdy has previously been a mixed experience between using Sturdy in the browser, and the
        Sturdy CLI (to manage the connection to your filesystem). This is finally about to change
        &mdash; <strong>Sturdy is going native!</strong>
      </p>

      <p>
        The Sturdy app (for Mac, Linux, and Windows) greatly improves how Sturdy interacts with your
        system. The app is bundled with everything you need to get started, and there’s no longer
        any need to use the cli (in fact, we’re deprecating it).
      </p>

      <p>
        No additional authentication step is required, all you need to do to get started is to
        download the app and log in, and you’re ready to code!
      </p>

      <p>
        The ability to use Sturdy in the browser is not going away, but some actions, like
        connecting a Sturdy to a directory on your computer, will only be available through the app.
      </p>

      <p>
        <strong>We’re looking for more beta testers</strong> (especially testers on Linux and
        Windows), if you’re interested, please reply to this email!
      </p>

      <h2>Codebase homepage</h2>

      <p>
        The home for each codebase in Sturdy has received a facelift! It’s now highlighting up to 4
        of the latest changes in the codebase, giving you quick access to all files in the codebase
        (<i
          ><a href="https://getsturdy.com/features/access-control" target="_blank"
            >that you have access to</a
          ></i
        >).
      </p>

      <h2>Importing PRs</h2>

      <p>
        When you’re setting up Sturdy from GitHub, any open pull requests that you might have will
        now be imported as Workspaces in Sturdy. Just like before, any updates that you make to the
        code in this workspace, will be pushed to your existing PR.
      </p>

      <h2>Watching workspaces, @-mentions, and notifications</h2>

      <img src="./sidebar.png" height="207" width="240" class="float-left mb-4 mr-4" />

      <p>
        To keep your busy projects from annoying you, you can now
        <strong>watch your favourite workspaces</strong>, and you’ll only be notified about changes
        in your watched workspaces.
      </p>

      <p>
        Interacting with a workspace, such as when giving some feedback on it, will cause you to
        watch it.
      </p>

      <p>
        Sturdy now also supports @-mentions in comments, mentioning someone will always notify them,
        and also highlight the workspace that they were mentioned in with a badge in the sidebar.
      </p>

      <h2>Native CI/CD with Buildkite (beta!)</h2>

      <p>
        Sturdy has previously supported
        <a href="https://getsturdy.com/docs/continuous-integration" target="_blank">
          CI/CD via your connected GitHub repository </a
        >, but we’re now happy to introduce a <strong>native CI/CD experience</strong> through
        <a href="https://buildkite.com/" target="_blank">Buildkite</a>!
      </p>

      <p>
        This is a completely Sturdy-native experience, and does not rely on git. The native CI/CD
        integration supports large files, file level permissions, and more. Oh, and it’s very fast!
      </p>

      <p>
        Reach out if you’re interested in getting access to the beta, and let us know which CI/CD
        provider that we should integrate with next!
      </p>

      <h2>More....</h2>

      <p>
        Since the last newsletter, we’ve shipped over 500 features and improvements, so there’s way
        too many to cover them all. Here are some of my favourites.
      </p>

      <ul>
        <li>Connecting to a workspace has been simplified. Click connect, and you’re ready!</li>
        <li>
          First time users now get a tutorial through Sturdy when they use Sturdy for the first time
        </li>
        <li>
          <a href="https://getsturdy.com/" target="_blank">Big updates to the landing page</a>
        </li>
        <li>
          Improved integration with pull requests, merged PRs are now appearing in the workspace
          activity feed, and more
        </li>
        <li>API to download an archive of your codebase</li>
        <li>Improved supports for files larger than 500 MB</li>
        <li>Windows releases of the sturdy-cli are now signed and trusted by default</li>
        <li>New codebases are seeded with a welcome workspace</li>
      </ul>

      <p>
        Thanks for using Sturdy, see you next week (yes, really!) with another blog post!<br />&mdash;
        Gustav and Team Sturdy
      </p>
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
