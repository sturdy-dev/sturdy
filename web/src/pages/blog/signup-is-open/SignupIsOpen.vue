<template>
  <BlogPost
    title="Sturdy is Here"
    subtitle="Sign ups for Sturdy are now open!"
    :author="author"
    date="August 12, 2021"
    description="
      Making a new version control system is definitely not for the faint-hearted!
      It has taken us a lot of work to lay the foundation for the code collaboration platform that we are building.
      Now that you can try out Sturdy for yourself, I would like to take this opportunity to also share some
      more on the vision Gustav and I are working towards."
  >
    <template #default>
      <p>
        Making a new version control system is definitely not for the faint-hearted! It has taken us
        a lot of work to lay the foundation for the code collaboration platform that we are
        building.
      </p>

      <p>
        Now that you can
        <router-link :to="{ name: 'signup' }"> try out Sturdy for yourself</router-link>
        , I would like to take this opportunity to also share some more on the vision we are working
        towards.
      </p>

      <h2>Okay, why make a new version control when everybody uses Git?</h2>
      <p>
        Git is great at versioning, so it is no surprise it has become so popular with open source
        projects. However, it turns out that in <em>software teams</em>, managing code changes, and
        collaborating with other developers are conflated.
      </p>

      <p>
        Like many other coders, Gustav and I have also used Git at our past jobs &mdash; making
        software products in team settings, trying to help one another.
      </p>

      <p>
        The thing is, Git was made specifically for the needs and workflow of the Linux kernel
        hackers, who are working in a hierarchical, distributed, and slower paced environment. This
        workflow fits open source just fine but we have seen
        <strong>how different the workflow at work happens to be</strong>.
      </p>

      <p>
        At a software company your code usually has a single source of truth. You collaborate
        closely with the rest of your team to solve problems for users, which often means you wish
        to have quick turnarounds in getting fixes and shipping code.
      </p>

      <h2>Observations we have made working as developers in a team setting</h2>
      <p>
        <strong
          >It is preferable to ship small and incremental changes as frequently as you can.</strong
        >
        This means that with a tighter feedback loop, preferably with CI/CD, the risk of having a
        breaking change is greatly reduced, and peer review becomes much easier. However, the
        overhead of doing this with Git is enormous, and often prompts team members to remind each
        other to “not write more than X lines of code” per change.
      </p>

      <p>
        <strong>Giving meaningful feedback on code is hard.</strong> Many teams are doing peer
        review as part of their development process, but how many bugs do we really catch and
        prevent with this process, vs. the hidden cost of waiting for review? Bugs still reach
        production after all, even if the change was declared <em>"LGTM"</em> by a team member.
        Unfortunately, Pull Requests often make it easy to focus on nitpicking and difficult to
        casually discuss the design.
      </p>

      <p>
        <strong>Fetching and running someone else's work-in-progress code takes effort.</strong>
        This is especially if you are in the middle of something else. Giving code suggestions
        asynchronously is even harder.
      </p>

      <p>
        <strong>You want to know what your team-members are up to.</strong>
        It is natural for engineers to want to help their team with knowledge that they have. In a
        dynamic work setting people also want to prevent collision of work, as well as keep track of
        progress. Sometimes teams would actively encourage each other to have a cadence of “pushing”
        code to their branch regularly, and before the end of day.
      </p>

      <p>
        <strong>Git is actually hard.</strong> A huge chunk of the questions on Stack Overflow are
        “How do I X in git”. It's a powerful tool, but it comes with a steep learning curve that you
        have to overcome very early in your career. However, an engineer may rather focus on their
        specific business problem, instead of having to think about about low-level VCS operations
        (even if they are fully proficient in it).
      </p>

      <figure>
        <img
          class="w-full rounded-lg"
          src="./git-stackoverflow.png"
          alt="Three out of the top 4 Stack Overflow questions are 'How do I X in Git?'"
          width="950"
        />
        <figcaption>
          Three out of the top 4 Stack Overflow questions are "How do I X in Git?"
        </figcaption>
      </figure>

      <h2>Our vision for the future of team code collaboration</h2>

      <p>
        We believe that developer tools should empower teams to experiment, share, and discuss
        software ideas as well as ultimately ship their work with as little friction as possible. In
        a nutshell, this means working at a <strong>higher abstraction level</strong> than Git and
        prior VCS systems.
      </p>
      <p>
        Greater abstraction means a more leveraged development environment &mdash; getting more done
        with less effort. This is similar to how high level programming languages, such as Python,
        are more expressive and powerful than for example C, at the cost of giving up some control.
      </p>

      <p>
        Making our code versioning and collaboration tools smoother, means that engineers get to
        focus on the business problems they are solving instead of googling how to undo something.
      </p>

      <h2>Sturdy is cloud first</h2>
      <p>
        We have designed Sturdy as a cloud first platform because it just fits the way software
        teams work &mdash; together. The
        <router-link :to="{ name: 'featuresLiveFeedback' }">online nature</router-link>
        of Sturdy means you don't have to spend any effort synchronizing code.
      </p>

      <p>
        With awareness of changes from the entire team in real-time, Sturdy can provide a lot of
        actionable insights as well enabling easy
        <router-link :to="{ name: 'blogRestoreToAnyPointInTime' }">one click</router-link>
        actions aiding the development process. This falls in the category of workflow augmentation
        and includes anything from "click to undo", to more frequent CI feedback, as well as
        notifying of potential conflicting work.
      </p>
      <p>
        What is even more exciting is enabling developers to help each other,
        <router-link :to="{ name: 'blogHumaneCodeReview' }"
          >discuss and try out each other's code
        </router-link>
        with little to no effort. Consider for a moment teamwork with Google Docs compared to
        manually emailing and merging files. Sturdy enables this level of collaboration while still
        using your local text editor, and without forcing you to use multi-cursor collaboration.
      </p>

      <div class="relative text-base mx-auto max-w-prose lg:max-w-none">
        <figure>
          <div class="aspect-w-12 aspect-h-7 lg:aspect-none">
            <img
              class="rounded-lg object-cover object-center"
              alt="Jumping between Workspaces in Sturdy"
              src="../../../assets/Web/NavBanner.png"
            />
          </div>
          <figcaption class="mt-3 flex text-sm text-gray-500">
            <span class="ml-2">Jumping between Workspaces in Sturdy</span>
          </figcaption>
        </figure>
      </div>

      <h2>Just The Beginning</h2>
      <p>There are have big plans for Sturdy and we are just getting started!</p>

      <p>
        If you have any questions, check out
        <router-link :to="{ name: 'resourcesDocs' }"
          >our documentation &mdash; we are continuously adding new articles to it
        </router-link>
        .
      </p>

      <p>
        I hope you enjoy Sturdy and would love to hear what you think &mdash; join our
        <a href="https://discord.gg/5HnSdzMqtA">Discord community</a> or message me on Twitter
        <a href="https://twitter.com/krlvi">@krlvi</a>. <br />&mdash; Kiril
      </p>
    </template>
  </BlogPost>
</template>

<script lang="ts" setup>
import BlogPost from '../BlogPost.vue'
import avatar from '../kiril.jpeg'

const author = {
  name: 'Kiril Videlov',
  avatar: avatar,
  link: 'https://twitter.com/krlvi',
}
</script>
