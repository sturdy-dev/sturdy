<template>
  <BlogPost
    title="This week at Sturdy #18"
    subtitle="This week at Sturdy"
    date="May 3, 2022"
    description="Hello world! In this newsletter we're sharing some of the biggest updates and new features in
      Sturdy v1.8, let's get started!"
    :author="author"
    :image="ogImageFull"
  >
    <template #introduction>
      Hello world! In this newsletter we're sharing some of the biggest updates and new features in
      Sturdy v1.8, let's get started!
    </template>
    <template #default>
      <p>
        Hey ducklings (that's you, all Sturdy fans!) Sturdy <code>v1.8.0</code> is just
        <em>fresh of the compiler</em> and we'd like to share all of the goodies in this new
        release.
      </p>

      <p>
        The
        <a href="https://github.com/sturdy-dev/sturdy/blob/api/v1.8.0/CHANGELOG.md">CHANGELOG</a>
        contains a condensed list of everything that's new, but in this post we're digging deeper
        into all of the new features and improvements.
      </p>

      <h2>Undo and redo</h2>

      <p>
        Behind the scenes, Sturdy has been been operating on "snapshots" for a while now. It's
        what's powering all of our drafts, and for example powers the automatic store and restore of
        changes when you're jumping between drafts.
      </p>

      <p>
        Snapshots are a stored in a linked list, where each draft has a "head" snapshot, that in
        turn can have a "parent" snapshot (this becomes important later).
      </p>

      <p>
        Sturdy automatically creates snapshots in the background when you're coding (or really, when
        we detect changes to files on your filesystem), and when you're doing manual actions to a
        draft through the UI (such as undoing a file).
      </p>

      <p>
        The new undo/redo functionality operates entirely on the snapshots, and exposes them to the
        user.
      </p>

      <CardWithFooter>
        <template #default>
          <img src="./undo-redo.png" height="364" width="1706" alt="Undo and redo on Sturdy" />
        </template>
        <template #footer> Undo and redo in Sturdy (in the upper right corner). </template>
      </CardWithFooter>

      <h2>CI/CD on drafts</h2>

      <p>
        <a href="https://getsturdy.com/blog/2022-04-12-this-week-at-sturdy">In v1.7.0</a> we
        introduced support for running your automated testing via first-party integrations. In
        <code>v1.8.0</code> we're expanding the support for CI/CD over the "Sturdy for GitHub"
        bridge and are adding support for all CI/CD providers that integrate with GitHub, including
        GitHub Actions, CircleCI, and more! This works by pushing a the draft to a branch named
        <code>sturdy-ci-${NAME}</code> to GitHub, and using the push event as the trigger to run the
        tests.
      </p>

      <CardWithFooter>
        <template #default>
          <img
            src="./ci-cd-via-github.png"
            height="522"
            width="1392"
            alt="Sturdy CI/CD via GitHub"
          />
        </template>
        <template #footer>
          How you can integrate any CI/CD provider with Sturdy via <em>Sturdy for GitHub</em>.
        </template>
      </CardWithFooter>

      <p>
        The <a href="https://getsturdy.com/docs/continuous-integration">CI/CD documentation</a> has
        been updated and goes further in-depth about how this works.
      </p>

      <h2>Highlights</h2>

      <ul>
        <li>
          <strong>[Improvement]</strong> Improved caching of codebase contents, making operations
          like "undo" and "merge" significantly faster
        </li>
        <li>
          <strong>[Improvement]</strong> Fixed a data-race where sometimes a change could be
          imported (from GitHub or other) multiple times (leading to a confusing changelog)
        </li>
        <li>
          <strong>[Improvement]</strong> Improved reliability when importing extremely large pull
          requests (+100k lines changed)
        </li>
        <li>
          <strong>[Improvement]</strong> Better performance when GitHub webhook delivery is slow
          (added internal handling that does not rely on webhooks)
        </li>
        <li>
          <strong>[Improvement]</strong> Register the Sturdy app as a handler for the
          <code>sturdy://</code> protocol scheme on Linux (App Images, deb, and rpm)
        </li>
        <li>
          <strong>[Fix]</strong> Improved first time boot performance of the server, and fixed a
          race condition where sometimes the server did not successfully start the first time.
        </li>
        <li>
          <strong>[Fix]</strong> Fixed a bug where navigation between drafts could take you to the
          wrong page
        </li>
        <li>
          <strong>[Fix]</strong> Fixed a bug where the callback from GitHub after updating
          permissions for the Sturdy app could take you to an unexpected page
        </li>
        <li>
          <strong>[Fix]</strong> Fixed a bug where comments on "live" code could sometimes "jump"
          around
        </li>
      </ul>

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
import CardWithFooter from '../../../atoms/CardWithFooter.vue'

let ogImageFull = computed(
  () => new URL(ogImage, import.meta.env.VITE_WEB_HOST as string | undefined).href
)

const author = {
  name: 'Gustav Westling',
  avatar: avatar,
  link: 'https://twitter.com/zegl',
}
</script>
