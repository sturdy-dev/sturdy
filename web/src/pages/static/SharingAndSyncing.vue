<template>
  <StaticPage
    title="Sharing and Syncing changes"
    subtitle="Staying up to date with changes done by others"
    category="documentation"
    metadescription="Staying up to date with changes done by others"
  >
    <div
      class="mt-5 prose prose-yellow text-gray-500 mx-auto lg:max-w-none lg:row-start-1 lg:col-start-1"
    >
      <p>
        In Sturdy code changes made by you and your teammates happen inside workspaces. Similar to
        Git branches, Sturdy workspaces allow for isolating work efforts from one another.
      </p>
      <p>
        It is inside of a workspace that you see and manage your in-progress code changes - you can
        describe what you are working on, select and undo updates. You can read more about
        workspaces in the article on
        <router-link :to="{ name: 'featuresLiveFeedback' }">real-time feedback</router-link>
      </p>

      <p>
        Let's say you are ready with your changes - you have gotten
        <router-link :to="{ name: 'featuresLiveFeedback' }">feedback and comments</router-link>
        from your team, perhaps even some
        <router-link :to="{ name: 'blogHumaneCodeReview' }">code suggestions</router-link>.
      </p>

      <div class="relative text-base mx-auto max-w-prose lg:max-w-none">
        <figure>
          <div class="aspect-w-12 aspect-h-7 lg:aspect-none">
            <img
              class="rounded-lg object-cover object-center"
              alt="Describing a code change in Sturdy"
              src="../../assets/Web/ChangeDescription.png"
            />
          </div>
          <figcaption class="mt-3 flex text-sm text-gray-500">
            <span class="ml-2">Describing the changes in a workspace</span>
          </figcaption>
        </figure>
      </div>

      <h3>What happens when you press share?</h3>
      <p>
        In short, after sharing, the code is considered 'landed' and available for deployment as
        well as for other workspaces. Sturdy embraces
        <a href="https://trunkbaseddevelopment.com/">trunk based development</a> - in a nutshell,
        encouraging frequent integration of the team's work. Frequent sharing of small incremental
        pieces of work avoids creating distance between developers as well helps in avoiding
        conflicts.
      </p>

      <div class="relative text-base mx-auto max-w-prose lg:max-w-none">
        <figure>
          <div class="aspect-w-12 aspect-h-7 lg:aspect-none">
            <img
              class="rounded-lg object-cover object-center"
              alt="Codebase changelog"
              src="../../assets/Web/CodebaseChangelog.png"
            />
          </div>
          <figcaption class="mt-3 flex text-sm text-gray-500">
            <span class="ml-2"
              >After a change has been shared from a workspace it will be visible in the codebase's
              changelog above the workspaces.</span
            >
          </figcaption>
        </figure>
      </div>

      <h3>Syncing</h3>
      <p>
        After a new change has been shared from a given workspace, all other workspaces are then no
        longer up to date with the main copy (trunk) of the code. This is similar to when a Git
        branch is behind the 'main' branch. Because Sturdy is aware of all workspaces, it is able to
        indicate when you are behind and prompt for syncing. In some cases it could also sync
        automatically.
      </p>

      <div class="relative text-base mx-auto max-w-prose lg:max-w-none">
        <figure>
          <div class="aspect-w-12 aspect-h-7 lg:aspect-none">
            <img
              class="rounded-lg object-cover object-center"
              alt="Codebase changelog"
              src="../../assets/Web/SyncPrompt.png"
            />
          </div>
          <figcaption class="mt-3 flex text-sm text-gray-500">
            <span class="ml-2">Indication that the workspace is behind</span>
          </figcaption>
        </figure>
      </div>
      <p>
        If there happen to be any conflicts between the code from the trunk and that in the
        workspace,
        <router-link :to="{ name: 'featuresConflicts' }"
          >Sturdy allows you to resolve them</router-link
        >
        visually in the browser or manually in your text editor of choice.
      </p>
    </div>
  </StaticPage>
</template>

<script>
import StaticPage from '../../layouts/StaticPage.vue'
export default {
  name: 'SharingAndSyncing',
  components: { StaticPage },
}
</script>

<style scoped></style>
