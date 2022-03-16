<template>
  <StaticPage
    title="Conflicts"
    category="solutions"
    metadescription="How Sturdy deals with and prevents conflicts"
  >
    <div class="text-base max-w-prose mx-auto lg:max-w-none">
      <p class="text-lg text-gray-500">How Sturdy deals with and prevents conflicts.</p>
    </div>
    <div
      class="mt-5 prose prose-yellow text-gray-500 mx-auto lg:max-w-none lg:row-start-1 lg:col-start-1"
    >
      <p>
        When two (or more) developers are touching the same lines of code, Sturdy won't be able to
        automatically sync the code between the two developers.
      </p>
      <p>
        As an example, imagine that we have a function <code>Add</code> that takes two ints and
        returns another int. Two colleagues, Frank and Sarah happen to both make changes to this
        function in parallel. Sarah changes the types to float64, and Frank changes the names of the
        arguments.
      </p>

      <pre>
- func Add(a, b int) int
+ func Add(a, b float64) float64</pre
      >

      <pre>
- func Add(a, b int) int
+ func Add(x, y int) int</pre
      >

      <p>
        It's no longer possible to automatically sync these two changes as they are touching the
        same lines of code.
      </p>
      <p>So what do you do?</p>
      <p>
        Let's say that Sarahs changes are approved first, and she shares them with the rest of the
        team. The next time that Frank syncs new changes into his draft change, the conflict will be
        detected, and he'll have to decide how to resolve the conflict.
      </p>

      <div class="relative text-base mx-auto max-w-prose lg:max-w-none">
        <figure>
          <div class="lg:aspect-w-12 aspect-h-7 lg:aspect-none">
            <img
              class="rounded-lg object-cover object-center"
              alt="Resolving conflicts in Sturdy"
              src="../../assets/Web/Conflict.png"
              width="2046"
              height="556"
            />
          </div>
          <figcaption class="mt-3 flex text-sm text-gray-500">
            <span class="ml-2">Resolving Conflicts in Sturdy</span>
          </figcaption>
        </figure>
      </div>

      <p>
        For each file where there are conflicts, Frank will be given three options. He can either
        decide to keep only one version of the changes, keep the file as it looked before syncing,
        and override Sarahs changes. Or doing the opposite, allowing Sarahs changes to replace the
        one he made. If neither of these alternatives are good enough, Frank can resolve the
        conflict manually in his editor, and submit a manual resolution.
      </p>

      <div class="relative text-base mx-auto max-w-prose lg:max-w-none">
        <figure>
          <div class="lg:aspect-w-12 aspect-h-7 lg:aspect-none">
            <img
              class="rounded-lg object-cover object-center"
              alt="Resolving conflicts in Sturdy"
              src="../../assets/Web/ConflictCustomResolutionWeb.png"
              height="750"
              width="2056"
            />
          </div>
          <figcaption class="mt-3 flex text-sm text-gray-500">
            <span class="ml-2">Initiating a manual conflict resolution on Sturdy</span>
          </figcaption>
        </figure>
      </div>

      <p>
        Solving a conflict manually can look a bit intimidating, but is easy to do if you're
        familiar with Git. Sturdy exposes conflicts in a git-compatible format (using
        &lt;&lt;&lt;&lt;&lt;&lt;&lt;, =======, and &gt;&gt;&gt;&gt;&gt;&gt;&gt; as separators).
        Resolve the conflict in your local text editor, and when you're done, mark the conflict as
        done on Sturdy.
      </p>

      <div class="relative text-base mx-auto max-w-prose lg:max-w-none">
        <figure>
          <div class="aspect-w-12 aspect-h-7 lg:aspect-none">
            <img
              class="rounded-lg object-cover object-center"
              alt="Resolving conflicts in Sturdy"
              src="../../assets/Web/ConflictCustomResolution.png"
            />
          </div>
          <figcaption class="mt-3 flex text-sm text-gray-500">
            <span class="ml-2">Resolving a conflict manually in Sturdy</span>
          </figcaption>
        </figure>
      </div>

      <p>And that's how easy it is to resolve conflicts in Sturdy!</p>

      <p v-if="false">
        Your team deserves an easy way to resolve conflicts,
        <router-link :to="{ name: 'signup' }">get started with Sturdy now!</router-link>
      </p>
    </div>
  </StaticPage>
</template>

<script lang="ts" setup>
import StaticPage from '../../layouts/StaticPage.vue'
</script>
