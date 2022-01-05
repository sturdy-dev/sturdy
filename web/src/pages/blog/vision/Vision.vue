<template>
  <BlogPost
    title="Unbreaking code collaboration"
    surtitle="Sturdy"
    subtitle="Our vision for the future"
    :author="author"
    date="August 18, 2021"
    description="I feel like collaborating on code today is broken. It is too hard to get started and even after learning all the basic tools and commands there is just too much effort in making contributions as a team. Why is it so hard to give meaningful feedback on each other's code? How come the rest of the stack has evolved, when collaboration hasn't?"
    reading-time="10 minutes"
  >
    <template #default>
      <p>
        I feel like collaborating on code today is broken. It is too hard to get started and even
        after learning all the basic tools and commands there is just too much effort in making
        contributions as a team. Why is it so hard to give meaningful feedback on each other's code?
        How come the rest of the stack has evolved, when collaboration hasn't?
      </p>

      <p>
        Kiril and I could not let go of the idea that code collaboration doesn't have to be like
        this – and that there must be a better way. So we started sketching and building on a new
        version control system to make collaborating on code easier.
      </p>

      <p>
        We quit our jobs to work on Sturdy just over 9 months ago. In this post I want to cover what
        we've learned since then, where we are now, and what our vision for the future of
        collaborating on code looks like.
      </p>

      <p>
        When designing what would later become Sturdy, our first goal was to create a version
        control that is <strong>higher level and more leveraged</strong> than existing systems. This
        means creating a system with actions closer to what a human might want. For example the
        intention to "ship the code that I just have written to production" translates into a single
        action, instead of a chain of low-level commands. Other examples of this would be undoing
        code to an earlier version, or running a colleague's code on your computer, which should
        also be one-click actions. In a way this is similar to how high level programming languages
        like Python are more expressive and powerful than C, at the cost of giving up some control.
      </p>

      <p>
        We have implemented this highly expressive environment by exposing the codebase through a
        two-way-synced "magic directory" on your computer (à la Dropbox). The management of what's
        "checked out" in the directory happens through Sturdy's web-based application. I'll cover
        more of this later, and how we're making it fast and easy to use.
      </p>

      <p>
        Our second goal with Sturdy is to enable <strong>great collaboration</strong>.
        Traditionally, the VCS and the collaboration software are tightly coupled, but versioning
        and working together are not the same thing.
      </p>

      <p>
        Sharing, discussing, and trying each other's code should not be a big deal. The way that
        engineers collaborate today, mostly by sharing code in large batches (via Pull Requests,
        Merge Requests, etc) feels too formal and rigid. We took inspiration from how some teams
        have chosen to do pair programming, and can casually discuss ideas and explore code
        together, and wanted to create an asynchronous version of that.
      </p>

      <p>
        We knew that we wanted to build a VCS tailored for teams who work with each other on a daily
        basis. Teams have already aligned their goals, on standups, or through the organization's
        quarterly objectives. Engineers in these teams are also incentivized to collaborate and help
        each other out, to ship features faster, and to spread knowledge.
      </p>

      <p>
        When it comes to spotting bugs, just browsing changes in a Pull Request doesn't quite cut
        it. We wanted to make it super easy (with approximately one-click) to get a copy of someone
        else's work in progress on your computer so that you can run and explore it in your IDE.
        Taking this even further we made it so that any edits that you make while having someone
        else's code checked out, are instantly available to the author as suggestions! The proposed
        code changes can be small or entire refactorings, which the author can take or reject. I
        believe that this enables a workflow with early feedback, that feels relaxed, and is easy to
        use, and that still allows you to work asynchronously at your own pace. With Sturdy's
        two-way-syncing, any suggestions the author accepts are synced down to their computer.
      </p>

      <p>
        If a software project was a human, the VCS would be it's veins, connecting to many other
        tools and services. We knew that to even have a chance in creating something that was
        possible for a team to adapt, that we had to have <strong>great compatibility</strong>.
        Compatibility with existing version control, continuous integration, and provide a migration
        path that doesn't break the workflow for the team. We're launching Sturdy with import from
        Git, as well as a non-breaking continuous migration from path GitHub, so that you can use
        Sturdy and GitHub side-by-side during the migration. The code can always be exported to Git
        at any time.
      </p>

      <p>
        The need for compatibility is also the reason why we're <em>not</em> building a development
        environment in the browser, or an extension to an IDE or editor. Any program that can read
        and write files from disk can be used with Sturdy.
      </p>

      <p>
        We started building Sturdy (it was called "collabd" at the time) as a side-project, before
        we decided to quit our jobs and do this full time. We mostly worked on Sturdy for fun, and I
        accidentally sneaked in another design goal – support enormous monorepos. I was frustrated
        by having to clone hundreds of gigabytes of history to my computer, and keep a large
        checkout of files even if those files would never be read on my workstation.
      </p>

      <p>
        For our first prototype, we built a VCS where the history tree was turned upside down and
        purely based on patches. The server only had a full copy of each file as it looked in the
        current "HEAD", and stored "negative" and "positive" patches for each change. To calculate
        the diff between two revisions, you would need to find a common ancestor, and add the
        patches together. We built the entire storage on CockroachDB, and interfaced to it via a
        FUSE (Filesystem in Userspace) mount. It allowed us to retrieve files just-in-time to reduce
        the amount of disk space used, but required you to always be online.
      </p>

      <p>
        There were many problems with this approach, and the number of enormous monrepos in the
        world is surprisingly small – the fact that FUSE only works well on Linux didn't help. So
        when we decided to build a company around what later became Sturdy, we came to our senses,
        and simplified things quite a bit.
      </p>

      <p>
        We scratched the monorepo goal, and redesigned the core of the version control.
        <strong>Sturdy is now a snapshot based version control system</strong> (like Git,
        Subversion, and a few others) as compared to the earlier "changeset" based system, and we're
        doing versioning and file storage on disk, with Postgres as a metadata database. FUSE lived
        on for a month or two, before we scrapped it as well, in favour for a system where the
        clients have a full copy of the "checkout" written to their normal file system, and a
        background daemon that's responsible for syncing files to and from Sturdy to this directory.
      </p>

      <p>
        Sturdy now has three main components, the version control backend, the file syncer, and the
        web app. Coding happens in a local text editor or IDE of choice, while codebase management
        happens on the web. On the web new changes are recorded, reviewed, workspaces are updated to
        the latest trunk, and changes are triggered to be synced down to connected computers.
      </p>

      <figure class="text-center">
        <img
          class="w-max-[100%] rounded-lg block m-auto"
          src="./graph.png"
          alt=""
          width="740"
          height="764"
        />
        <figcaption>Overview of Sturdy</figcaption>
      </figure>

      <p>
        The way we allow programmers to work in isolation from each other is through workspaces. It
        is inside a workspace where discussion and suggestions happen, and when ready – changes land
        towards a common trunk. We're doing this so that a team can work on as many things in
        parallel as they need to, while still encouraging the authors to land their changes as soon
        as possible, to minimize the divergence between engineers.
      </p>

      <p>
        Two-way-syncing is hard. When a workspace lands on the trunk, syncing is easy, as it’s all
        happening on the backend (and we can use a mutex). We know ahead-of-time if there will be a
        “merge conflict”, and ask the user to resolve it before the landing can happen. In the other
        case however, when we’re syncing files between the computer and Sturdy it’s not as easy. If
        the internet connection is stable enough, and the Sturdy daemon is started before code
        changes, there won’t be any problems.
      </p>

      <p>
        On the other hand, if there are “simultaneous” changes to both the local directory, and
        pending changes (like undoing the recent changes to a file) to be downloaded from the web,
        things can (and do) go haywire. Currently Sturdy will prioritize changes from the user over
        changes from the server if there is a disagreement. This way no data is lost, but it’s still
        an annoyance as you can end up in an unexpected state. We’re working on making this setup
        better soon, but it’s good enough for us to dare to launch.
      </p>

      <p>
        Sturdy's cloud-first approach enables interactions and workflows that were impossible to do
        before. My favourite example of this is the capacity to store the
        <strong>history in a high resolution</strong> between the final "revisions", just like
        Google Docs does for text. When a file is modified, and Sturdy slurps those changes up to
        the backend, a new snapshot is created. The first feature that we built off of this was the
        ability to "time travel" between earlier versions of the code, for example with one click
        you can go back to what a workspace looked like exactly 15 minutes ago.
      </p>

      <p>
        This history is only available to the original author of the code, but we’re excited to
        explore ideas for how to integrate this data (in a safe way) in code review and code
        exploration. If the reviewer can see how the code unfolded over time, it might be easier to
        review it, rather than looking at diffs in alphabetical order.
      </p>

      <p>
        Sturdy is online-first, which means that
        <strong>code changes are synced in real-time for team members to review and run</strong>. In
        traditional code review, code is only ever looked at when it’s considered to be “complete”,
        or at least close to complete, by the author. At this stage, it’s hard to discuss the bigger
        picture of the change, as the author might have already spent a substantial amount of time
        on it, and might be more emotionally invested in the solution than they care to admit.
      </p>
      <p>
        In my experience, this leads to a code review process that either entirely focusses on
        nitpicks on typos, or highlighting that a method is slightly too long, instead of guiding,
        reviewing, and helping out to create a good piece of software. Well, it's that – or skimming
        through the changes and declaring <i>"LGTM"</i>.
      </p>

      <p>
        Each bug that slips through the code review, and reaches production, is a sign that code
        review as we know it today is broken. That it's both too hard to really review code, and
        that we don't like it when we're blocked by a reviewer, so we just pretend to review each
        other's code, to be able to get back to coding sooner.
      </p>

      <p>
        We haven't built formal code review into Sturdy yet, with the ability to block the
        contribution if it hasn't been approved, as we haven't completely explored what we want to
        do, but we know for sure that we don't want to copy the incumbents. Lately, I've been
        questioning why every change needs to be signed off by a peer, and if there is a way to
        allow for safe contribution without blocking the author when waiting for review. What if the
        version control is able to determine if a change is safe enough to skip additional feedback?
        Or if work connected to a JIRA ticket labeled "quick fix" only "recommends" a review, but
        doesn't require it? Or if it can be allowed to be deployed to a staging environment or a
        subset of production servers, while it's awaiting review, and only require sign-off before
        deploying to production or 100% of the traffic.
      </p>

      <p>
        The real-time sync in Sturdy enables us to give <strong>real-time insights</strong>. Waiting
        for tests to run is a pain, and I’ve previously had the habit to create and push throwaway
        commits to get the continuous integration to run the tests for me in the cloud, to avoid
        having to slow integration tests on my computer. I’m excited about the possibility of
        rethinking when and how CI should execute tests, with a vision to always provide up-to-date
        test results. We’re not supporting CI just yet, but it’s coming soon.
      </p>

      <p>
        Another cool thing is that Sturdy is aware of what the entire team is working on, in
        real-time. It can, for example, notify a developer when they are about to work on the same
        stuff as a teammate is already working on. In that case it would be better to encourage
        collaboration between them, or at the very least be aware early on of the conflict.
      </p>

      <p>
        Sturdy is available <em>right now</em>. There are many opportunities ahead of us, and a lot
        of innovation waiting to happen. We hope that you like our vision, and what we already have
        launched today in our first version!
      </p>

      <p>
        While Sturdy is free for now, we plan to offer Sturdy as a SaaS, with a free tier for
        smaller teams, and a paid tier for larger organisations and codebases. You won’t lose any
        data the day we introduce the paid tiers. Getting paid hasn’t been our first priority, but
        we just want to make it clear that Sturdy is going to be a paid service, and that we’re not
        going to make money from other sources (like ads, or selling your data).
      </p>

      <p>
        Working as a developer nowadays is a messy ordeal – teams try to move quickly, requirements
        change, bugs need hotfixing. We are passionate about building tech to make the life of a
        developer a little more enjoyable.
      </p>

      <p>
        Thanks for reading,<br />
        – Kiril and Gustav
      </p>

      <ul>
        <li><a href="https://getsturdy.com/blog">Read more about Sturdy on our blog</a></li>
        <li>
          <a href="https://getsturdy.com/quickstart">Get started &mdash; our quickstart guide</a>
        </li>
        <li>
          <a href="https://news.ycombinator.com/item?id=28221109">Discuss this on Hacker News</a>
        </li>
      </ul>

      <Waitinglist />
    </template>
  </BlogPost>
</template>

<script>
import BlogPost from '../BlogPost.vue'
import avatar from '../gustav.jpeg'
import Waitinglist from '../../../components/Waitinglist.vue'

export default {
  components: { Waitinglist, BlogPost },
  setup() {
    const author = {
      name: 'Gustav Westling',
      avatar: avatar,
      link: 'https://twitter.com/zegl',
    }
    return { author }
  },
}
</script>
