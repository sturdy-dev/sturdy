<template>
  <StaticPage
    title="Live Feedback"
    subtitle="Sturdy is about real-time collaboration"
    category="solutions"
    metadescription="Give and take feedback, in real time with Sturdy"
  >
    <div
      class="mt-5 prose prose-yellow text-gray-500 mx-auto lg:max-w-none lg:row-start-1 lg:col-start-1"
    >
      There are two key points to how Sturdy works that are very different from other source control
      tools:
      <ul>
        <li>
          Code management and collaboration happens in the browser - creating and undoing changes,
          sharing and discussing with colleagues.
        </li>
        <li>
          Editing and running code happens on your workstation using your favorite text editor or
          IDE.
        </li>
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
        When a directory on your computer is connected to Sturdy, it gains real-time superpowers. As
        you code in this directory, the changes that you make are instantly shared to Sturdy and
        made available for your team to comment and give feedback on.
      </p>

      <p>
        The way this works, you can use any of your existing tools and editors for writing and
        testing code "locally". At the same time you no longer need to worry about manually sharing
        code. It also works the other way around - if you undo a change in the browser or go back to
        an earlier state of your work, the sturdy directory will be instantly updated.
      </p>

      <h3>Workspaces</h3>
      <p>
        In Sturdy you're coding in isolation from each other, in what we're calling a Workspace. You
        can have as many workspaces as you want, and easily switch between which one that you want
        to open on your computer from Sturdy. The changes in a workspace are continuously saved, and
        when you're coming back to an older workspace, everything will be restored to just how you
        left it.
      </p>
      <div class="relative text-base mx-auto max-w-prose lg:max-w-none">
        <figure>
          <div class="aspect-w-12 aspect-h-7 lg:aspect-none">
            <img
              class="rounded-lg object-cover object-center"
              alt="Jumping between Workspaces in Sturdy"
              src="../../assets/Web/NavBanner.png"
            />
          </div>
          <figcaption class="mt-3 flex text-sm text-gray-500">
            <span class="ml-2">Jumping between Workspaces in Sturdy</span>
          </figcaption>
        </figure>
      </div>

      <p>
        Because Sturdy saves work automatically you can easily and safely
        <router-link :to="{ name: 'featuresWorkspaceNavigation' }">
          jump to another workspace
        </router-link>
        and try out the code on your computer and then go back to to what you are working on before.
      </p>

      <p>
        All workspaces in a codebase are available to all of it's collaborators. When you go to
        someone elses workspace in the browser, you'll see what they are up to
        <i>right now</i>. From there, you can <strong>comment on their code in real-time</strong>,
        and if you want to give suggestions for improvements, checking out their workspace on your
        computer is just one click away.
      </p>
      <div class="relative text-base mx-auto max-w-prose lg:max-w-none">
        <figure>
          <div class="aspect-w-12 aspect-h-7 lg:aspect-none">
            <img
              class="rounded-lg object-cover object-center"
              alt="Commenting on live code in Sturdy"
              src="../../assets/Web/FizzbuzzComment.png"
            />
          </div>
          <figcaption class="mt-3 flex text-sm text-gray-500">
            <span class="ml-2">Commenting on live code in Sturdy</span>
          </figcaption>
        </figure>
      </div>
      <p>
        Sturdy is <strong>shifting left</strong>, and you can give and take feedback on code as soon
        as it's been written! In Git-based workflows it's not uncommon for code to be shared only
        after it's considered to be &ldquo;complete&rdquo;, and you might discover way to late that
        you've been working on the wrong thing all along.
      </p>
      <p>
        Your team deserve live-feedback,
        <router-link :to="{ name: 'signup' }">sign up to Sturdy now!</router-link>
      </p>
    </div>
  </StaticPage>
</template>

<script>
import StaticPage from './StaticPage.vue'
export default {
  components: { StaticPage },
}
</script>
