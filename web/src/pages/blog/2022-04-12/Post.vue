<template>
  <BlogPost
    title="This week at Sturdy #16"
    subtitle="This week at Sturdy"
    date="April 12, 2022"
    :author="author"
    :image="ogImageFull"
  >
    <template #introduction>
      Hello world! In this newsletter we're sharing some of the biggest updates and new features in
      Sturdy v1.7, let's get started!
    </template>
    <template #default>
      <p>
        Since
        <a href="https://getsturdy.com/blog/2022-03-10-introducing-draft-changes">the last issue</a
        >, where we introduced a whole new workflow with Draft Changes, we've release three minor
        versions of Sturdy (v1.5.0, v1.6.0, and v1.7.0).
        <i>
          (Minor releases when speaking in Semver, but major when it comes to what they enable, and
          how proud we are of them!)
        </i>
      </p>

      <h2>Highlights in v1.5, v1.6, and v1.7</h2>

      <img src="./IntegrateWithOtherServices.png" height="297" width="912" />

      <p>
        Sturdy <code>v1.5.0</code> introduced a new workflow in Sturdy, with Sturdy running on top
        of <strong>any</strong> Git repository and provider. Sturdy can now
        <router-link :to="{ name: 'docsIntegrationsGit' }">connect natively</router-link>
        over the git-protocol to any hosting provider, and brings all of the power and ease-of-use
        of using Sturdy to providers such as GitLab, Bitbucket, and Azure DevOps.
      </p>

      <p>
        The same release also added a much-needed improvement to the search in Sturdy, and support
        for running Sturdy on <code>linux/arm64</code> and automatic TLS support via Let's Encrypt.
      </p>

      <p><code>v1.6.0</code> introduced SSH support for the Git integration.</p>

      <img src="./CreateAndMergePullRequest.png" height="159" width="912" />

      <p>
        <code>v1.7.0</code> added a one-click-ship 🚢 feature to GitHub, where Sturdy can create a
        PR and merge it for you right away. Shipping your code with Sturdy is now
        <i>literally</i> 1-click away. We also focused on code-review in v1.7.0, bringing
        code-review of code written outside of Sturdy to Sturdy. Sturdy now makes all open pull
        requests to your repository available as <i>draft changes</i> in Sturdy.
      </p>

      <img src="./TriggerCI.png" height="108" width="912" />

      <p>
        This release also contained support for triggering and running your CI/CD pipelines based on
        any draft, at any time. So that you can utilize the powers of your CI runners to run your
        tests of your non-commited changes, and keeping my computer fans nice and idle. ;-)
      </p>

      <p>
        See
        <a href="https://github.com/sturdy-dev/sturdy/blob/main/CHANGELOG.md">the CHANGELOG</a> for
        a full list of changes.
      </p>

      <h2>Upgrading our team!</h2>

      <p>
        We now have three open positions to come and join our team! You'll get to work on some
        really cool open-source tech, and have fun while doing it!
      </p>

      <ul>
        <li>
          <a href="https://getsturdy.com/careers/senior-backend-software-engineer">
            Senior Backend Software Engineer
          </a>
          — (Go, PostgreSQL, GraphQL)
        </li>
        <li>
          <a href="https://getsturdy.com/careers/senior-frontend-software-engineer">
            Senior Frontend Software Engineer
          </a>
          — (Vue, TypeScript, urlq)
        </li>
        <li>
          <a href="https://getsturdy.com/careers/full-stack-engineer">Full Stack Engineer</a>
        </li>
      </ul>

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
import ogImage from './og_image_oss.png'

import { computed } from 'vue'

let ogImageFull = computed(() => new URL(ogImage, import.meta.env.VITE_WEB_HOST).href)

const author = {
  name: 'Gustav Westling',
  avatar: avatar,
  link: 'https://twitter.com/zegl',
}
</script>
