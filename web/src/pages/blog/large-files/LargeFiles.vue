<template>
  <BlogPost
    title="Enormous files, GitHub migration, and more!"
    subtitle="This week at Sturdy"
    date="September 9, 2021"
    :author="author"
  >
    <template #introduction>
      Hey Sturdy-fans! Here's an update with all of the major changes since the last time we hit
      your inbox.
    </template>
    <template #default>
      <h2>The doors are open!</h2>
      <p>Sturdy is now open for everyone &mdash; for you, your team, and your friends!</p>

      <p>
        <a href="https://getsturdy.com/signup"> &raquo; Sign up </a>
      </p>

      <h2>Large files</h2>
      <p>Sturdy now supports large files, and <em>huge</em> codebases.</p>
      <p>
        Thanks to some heavy work behind the scenes. Using large (gigabyte-sized) assets such as
        images or large binary files with Sturdy is now straight forward.
      </p>
      <p>
        Large files on Sturdy are downloaded on-demand, and the size of the projects history doesn't
        slow you down, or require you to use additional storage on connected computers.
      </p>
      <p>
        <a href="https://getsturdy.com/features/large-files">
          &raquo; Read more about large files
        </a>
      </p>

      <h2>GitHub integration</h2>
      <p><img src="../../../assets/Web/GitHubPullRequest.png" height="498" width="937" /></p>
      <p>
        Sturdy now plays nicely with Git and GitHub! It's always been important for us to build
        Sturdy with compatability in mind, and we're now reaping the benefits of that.
      </p>
      <p>
        Using the GitHub plugin for Sturdy, the two can now be used side-by-side, so that you can
        get the benefits of Sturdy without having to migrate everything of of GitHub.
      </p>
      <p>
        You can for example use Sturdy for all your collaboration needs, while still using your
        CI/CD integrated through GitHub.
      </p>
      <p>
        Changes on Sturdy are automatically exported as git commits, and can be pushed to your
        repository through Pull Requests or directly to your HEAD branch.
      </p>
      <p>
        <a href="https://getsturdy.com/features/migrate-from-github">
          &raquo; Read more in our knowledgebase
        </a>
      </p>

      <h2>Notifications</h2>
      <p>
        We get it &mdash; you're a busy person, and are juggling many tasks and workspaces in
        parallel. Sturdy now shows you when you have something important waiting. Here's an example
        of what it looks like:
      </p>
      <div class="flex items-center">
        <div>
          <img src="./bell.png" height="247" width="127" />
        </div>
        <div>
          <img src="./notification.png" height="376" width="880" />
        </div>
      </div>

      <div class="clear-both"></div>

      <h2>Windows installer</h2>
      <p>
        We finally have a installation wizard for Windows, making it really easy to get started with
        Sturdy. Download links are available from the downloads page.
      </p>
      <p>
        <img src="./installer.png" height="422" width="544" />
      </p>

      <p>Thanks to everyone that's using Sturdy, you're the best!</p>
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
