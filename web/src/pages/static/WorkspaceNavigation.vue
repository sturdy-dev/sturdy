<template>
  <StaticPage
    title="Instant Workspace switching"
    subtitle="Navigate between coding workspaces with one click"
    category="solutions"
    metadescription="Navigate between coding workspaces with one click"
  >
    <div
      class="mt-5 prose prose-yellow text-gray-500 mx-auto lg:max-w-none lg:row-start-1 lg:col-start-1"
    >
      <p>
        Because Sturdy
        <router-link :to="{ name: 'featuresLiveFeedback' }"
          >manages the codebase folder on your workstation</router-link
        >, jumping between different work efforts is really easy. It takes a single click to move
        your development environment to another workspace, either yours or created by your team. All
        of this without losing any of your work in progress.
      </p>
      <p>
        We are building Sturdy in order to make experimentation and collaboration on code easier.
        Instant workspace switching is particularly useful when you:
      </p>
      <ul>
        <li>Suddenly spot something small that you wish to fix and deploy separately</li>
        <li>Would like to run somebody else's code on your machine</li>
        <li>Wish to make a code suggestion to your teammate</li>
      </ul>

      <div class="relative text-base mx-auto max-w-prose lg:max-w-none">
        <video
          autoplay
          loop
          muted
          playsinline
          alt="Navigating between Sturdy workspaces with one click"
          class="rounded-lg object-cover object-center"
        >
          <source src="../../assets/Landing/WorkspaceSwitching.webm" type="video/webm" />
          <source src="../../assets/Landing/WorkspaceSwitching.mp4" type="video/mp4" />
        </video>
      </div>
      <p>
        Compared to checking out different Git branches, switching workspaces is a very lightweight
        operation.
      </p>
      <p>
        There is no need to stash your work in progress because Sturdy manages keeps track of every
        change automatically in the cloud.
      </p>
      <p>
        There is no need to add and fetch remotes in order to get your colleague's latest code
        because changes are always synchronized in real time.
      </p>
      <p>
        <router-link :to="{ name: 'featuresLiveFeedback' }"
          >Read more: Instant feedback</router-link
        >
      </p>
    </div>
  </StaticPage>
</template>

<script>
import StaticPage from '../../layouts/StaticPage.vue'
export default {
  name: 'WorkspaceNavigation',
  components: { StaticPage },
}
</script>

<style scoped></style>
