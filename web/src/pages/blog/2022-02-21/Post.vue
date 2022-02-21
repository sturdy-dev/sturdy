<template>
  <BlogPost
    title="Sturdy is now open-source"
    subtitle="This week at Sturdy"
    date="February 21, 2022"
    :author="author"
  >
    <template #introduction>
      We have some really exciting news to share &mdash; Sturdy is now open-source! As developers, we love
      OSS, so to say that we are thrilled about this announcement would be an understatement!
    </template>
    <template #default>
      <h2>Open-source</h2>
      <p>
        Sturdy is now officially an open-source project. You can check out the
        <a href="https://github.com/sturdy-dev/sturdy">GitHub repository</a>
        and, while you are there, giving us a ⭐️ would be awesome!
      </p>
      <p>
        The project is open to contributions from everyone! We love to receive your feedback, bug
        reports, issues, and even code! If you would like to get involved, check out the
        <a href="https://github.com/sturdy-dev/sturdy/blob/main/CONTRIBUTING.md"
          >contribution guidelines</a
        >.
      </p>
      <p>
        Within the Sturdy team we use Sturdy for it's own development and have
        <a href="https://getsturdy.com/sturdy-zyTDsnY">the codebase</a>
        and <a href="https://github.com/sturdy-dev/sturdy">the repo</a>
        connected using the
        <router-link :to="{ name: 'v2DocsHowToSetupSturdyWithGitHub' }"
          >GitHub integration</router-link
        >. This allows us to develop on Sturdy and also take community contributions.
      </p>

      <h2>Running your own Sturdy instance</h2>
      <p>
        We know that for some teams, trying a new SaaS at work involves a lot of red tape. This is
        why we created an easy way to try Sturdy locally using Docker - the setup takes only a few
        minutes and you can still try most of the features.
      </p>
      <p>To get started, run this one-liner:</p>
      <pre>{{ dockerOneliner }}</pre>
      <p>
        You can find detailed description of the setup over in our
        <router-link :to="{ name: 'v2DocsSelfHosted' }">docs</router-link>.
      </p>

      <p>
        Of course, Sturdy in the Cloud is still the easiest way to use Sturdy with zero setup. At
        the same time, adding the ability to self-host Sturdy, makes it accessible to even more
        teams.
      </p>

      <h2>Organizations support</h2>

      <h2>Updated Documentation</h2>
      <p>
        Over the past few weeks we have been working on beefing up our docs. As developers ourselves
        we always check the docs of tools in order to understand how they <em>really</em> work.
      </p>

      <p>
        We have added articles with detailed run-downs of how different aspects of Sturdy work. For
        example:
      </p>
      <ul>
        <li>
          <router-link :to="{ name: 'v2DocsUsingSturdy' }">
            A getting started tutorial
          </router-link>
        </li>
        <li>
          <router-link :to="{ name: 'v2DocsHowToCollaborateWithOthers' }">
            Collaborating with others
          </router-link>
        </li>
        <li>
          <router-link :to="{ name: 'v2DocsHowToSwitchBetweenTasks' }">
            Switching between tasks
          </router-link>
        </li>
      </ul>
      <p>...and many others. As we move forward we will be investing even more in documentation.</p>

      <h2>New home page</h2>
      <p>
        Our <router-link :to="{ name: 'v2Index' }">home page</router-link> has received a redesign
        to better the goals and vision for Sturdy. It has a few GIFs showing a number of workflows -
        here is a teaser:
      </p>

      <video autoplay loop muted playsinline>
        <source src="../../landing/assets/collaborate.mp4" type="video/mp4" />
      </video>

      <p>
        Thanks for being with us on this journey for better code collaboration!
        <br />&mdash; Kiril
      </p>
    </template>
  </BlogPost>
</template>

<script lang="ts" setup>
import BlogPost from '../BlogPost.vue'
import avatar from '../kiril.jpeg'
import { DOCKER_ONELINER as dockerOneliner } from '../../../docker'

const author = {
  name: 'Kiril Videlov',
  avatar: avatar,
  link: 'https://twitter.com/krlvi',
}
</script>
