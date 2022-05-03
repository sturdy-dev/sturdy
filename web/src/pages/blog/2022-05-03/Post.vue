<template>
  <BlogPost
    title="What is stopping us from merging 30 pull requests per day?"
    subtitle="The hard parts about collaborating on code"
    date="May 03, 2022"
    :author="author"
    :dive-in-banner="false"
    description="Most software projects are done by teams rather than individuals, so how we integrate work is critical. Collaborating on code is surprisingly challenging because any change made in isolation creates distance between developers. Shipping small and often is probably the most effective strategy for mitigating this. But just how frequently can we ship?"
    reading-time="12 minutes"
  >
    <template #introduction>
      <p>
        Most software projects are done by teams rather than individuals, so how we integrate work
        is critical. Collaborating on code is surprisingly challenging because any change made in
        isolation creates distance between developers.
        <strong>Shipping small and often</strong> is probably the most effective strategy for
        mitigating this. But just how frequently can we ship?
      </p>
    </template>

    <template #default>
      <h2>Programming is as close as we get to magic</h2>
      <p>
        I have been interested in computers since I was a kid. I remember playing
        <a href="https://en.wikipedia.org/wiki/Prince_of_Persia_(1989_video_game)"
          >Prince of Persia</a
        >
        and wondering “What drives the pixels on the screen?” as well as “How is pressing →
        connected to forward movement in the game?”. At first, I was using my imagination and
        guessing how it worked. I was 11 and had a few wildly inaccurate ideas of what happens
        “inside” the computer.
      </p>

      <img src="./Img-1.png" alt="Prince of Persia classic" width="768" />

      <p>
        Later on, learning about how computers actually work was one of the most exciting
        experiences of my childhood. I would compare the feeling to that of when somebody shows you
        how a magic trick works. A huge “Aha!” moment, followed by a feeling of power in actually
        understanding the trick.
      </p>

      <blockquote>
        I wanted to learn how to program these silicon rocks that we have tricked into “thinking”
      </blockquote>

      <p>
        Probably my top “Aha” moment was reading about how a CPU works, described in simple terms,
        starting with how it adds numbers together. From that moment, I knew that I wanted to learn
        how to program these silicon rocks that we have tricked into “thinking”. I still believe,
        even today, that this is as close as one gets to performing magic. For completeness, here is
        <a href="https://www.youtube.com/watch?v=cNN_tTXABUA">a really cool video</a> explaining how
        a CPU works.
      </p>

      <p>
        Years later, I had the opportunity to study software engineering at uni. I learned about
        data structures, the network stacks and protocols and even graphics pipelines. While I was
        getting more and more familiar with how software was “made”, there was one
        <b>meta-problem</b> which sitting unanswered for me…
      </p>

      <h2>Collaborating on code mystery</h2>

      <blockquote>
        If large systems require teams of developers to create, how do developers effectively add
        their code to the same system?
      </blockquote>

      <p>
        I am not necessarily talking about version control, but rather the process of integrating
        logic and assumptions written by different programmers — making a system work as one whole.
      </p>

      <p>
        I never got a satisfactory explanation to this mystery. At my first job as a programmer in a
        team of developers, I was introduced to a number of processes and buzzwords. Being fresh out
        of uni, I was eager to learn how it was done for real. Before long, however, I was convinced
        this only skirted around the real issue of collaborating on code. Let me illustrate with an
        example.
      </p>

      <p>
        The company was building a financial system which consisted of a “Base product” and
        “Customization” layers, developed by different teams. Think of the “Base product” code as an
        upstream which gets forked and modified by different customization teams. Perhaps many of
        you can already spot the challenge.
      </p>

      <p>
        A huge issue we experienced was that code in the “Base product” and the “Customizations”
        always diverged significantly over time, no matter how much people coordinated development.
        Keeping any of the customization codebases up to date with the “Base product” always meant
        dealing with not just merge conflicts syntactically, but also semantically. Some weeks, this
        overhead was more than 50% of the work effort.
      </p>

      <img src="./Img-2.png" alt="Feature customization via forking" width="768" />

      <p>
        It is obvious, in retrospect, that the separation between “Base product” and “Customization”
        was an inappropriate one. We mitigated this challenge by increasing the frequency of
        integration. This was an extreme example of how code written by any given developer diverges
        over time from that of other developers. Furthermore, an example of the costs associated
        with it.
      </p>

      <p>
        Writing code independently and in parallel creates distance between programmers — this was
        precisely the core of my unanswered question about collaborating on code.
      </p>

      <h2>Minimizing the distance between developers</h2>

      <blockquote>Ship small and often. PS: Push your branch before lunch.</blockquote>

      <p>
        The best strategy for effective collaboration on software I have seen is one of minimizing
        the code distance between developers. There are two popular techniques for achieving this —
        trunk-based Development and Continuous Deployment. This is usually what people are referring
        to when they speak of “shipping small and often”.
      </p>

      <h3>Trunk-based development</h3>
      <p>
        <a href="https://trunkbaseddevelopment.com/">Trunk-based development</a> simply means
        avoiding prolonged development of features on branches in favor of merging small code
        contributions more frequently into the code “trunk”, often referred to as the “main”.
      </p>

      <p>
        Instead of completing a feature before merging it a week (or several) later, it is hidden
        behind a feature flag but merged incrementally multiple times a day. This way, the code
        context is available to everybody on the team.
      </p>

      <img src="./Img-4.png" alt="Trunk-based development" width="768" />

      <div class="text-sm text-right">
        Source: <a href="https://trunkbaseddevelopment.com">trunkbaseddevelopment.com</a>
      </div>

      <p>
        This is quite different from strategies like “<a
          href="https://www.atlassian.com/git/tutorials/comparing-workflows/gitflow-workflow"
          >Git Flow</a
        >” which promotes working in relative isolation from the team and dealing with the
        complexity it entails later. I bet you know that feeling when you are asked to review a pull
        request with 500+ lines of changed code. Trunk-based development is practiced by small and
        large teams alike precisely because it encourages sharing small chunks of code that are
        easier to reason around.
      </p>

      <h3>Continuous Deployment</h3>
      <p>
        While the purpose of trunk-based development is to minimize the distance between developers,
        Continuous Deployment is about minimizing the distance between developers and
        <strong>their users</strong>.
      </p>

      <p>
        Continuous Deployment is the practice of always having the latest version of the “trunk”
        code (main / master branch) deployed in the production environment. Knowing that the latest
        code is what the users experience dramatically simplifies the mental model of code and its
        operations.
      </p>

      <blockquote>(Safely) testing in production.</blockquote>

      <p>
        The reality of software development is that you don't really know anything until you test
        your code in a real environment. Embracing this fact means that you can
        <strong>plan for it</strong>. Deploying code continuously may sound scary, but it is known
        to
        <a
          href="https://cloud.google.com/blog/products/devops-sre/announcing-dora-2021-accelerate-state-of-devops-report"
          >minimize risks</a
        >. This is because shipping smaller changes spreads out the large risk of large Big Bang
        deploys over time.
      </p>

      <p>
        With smaller and more predictable deploys, teams get to invest in observability and
        monitoring of the software and really understand how changes they make in code affect users.
        By avoiding long-lived feature branches and working in isolation we are prompted to validate
        each other's assumptions early rather than after a week of development.
      </p>

      <blockquote>
        Here's something I've come to believe: Creators need an immediate connection to what they're
        creating. That's my principle.
      </blockquote>

      <p>
        This is a quote from a great talk by
        <a href="https://en.wikipedia.org/wiki/Bret_Victor">Bret Victor</a> called
        <a href="https://www.youtube.com/watch?v=PUv66718DII">"Inventing on Principle (video)"</a>.
        His point is largely about minimizing feedback loops, and this is also exactly what
        trunk-based development and Continuous Delivery are about.
      </p>

      <p>
        Okay, if we know that increasing the integration frequency of our code is a good thing, how
        far can we take it?
      </p>

      <h2>How often can we ship?</h2>

      <p>
        Nowadays, both trunk-based development and Continuous Deployment are very popular with
        modern software teams. Most of us strive to “Keep pull requests small” and “ship small &
        often”. Over the past 4 years, the number of teams deploying multiple times per day has
        <a
          href="https://cloud.google.com/blog/products/devops-sre/announcing-dora-2021-accelerate-state-of-devops-report"
          >almost quadrupled</a
        >.
      </p>

      <blockquote>What is stopping us from merging 30 pull request per day?</blockquote>

      <p>
        I have been wondering — given how successful the strategy of integrating code more
        frequently has been, can we take this even further? I mean, what is stopping us from merging
        30 pull requests per day?
      </p>

      <p>
        There are two parts to this. Firstly, there is a <strong>fixed overhead</strong> per code
        contribution. In addition to that, the <strong>traditional code review</strong> process is
        fairly heavy, as it was designed for a different workflow.
      </p>

      <h3>Fixed overhead per code contribution</h3>

      <p>
        Each code change, no matter how small, requires some logistics to be integrated. Let's
        consider what it takes to fix a typo in a documentation file:
      </p>

      <ul>
        <li>Fetching the latest code from a remote <code>git pull</code></li>
        <li>Creating a branch for the code change <code>git checkout -b feature/fix-typo</code></li>
        <li>Fixing the typo</li>
        <li>Staging the file with the relevant change <code>git add myfile.py</code></li>
        <li>Committing the code change <code>git commit -m "Fix typo in documentation"</code></li>
        <li>
          Pushing the code change to the remote
          <code>git push --set-upstream origin feature/fix-typo</code>
        </li>
        <li>Creating a pull (merge) request in your git provider</li>
      </ul>

      <p>
        Of course, there is a technical reason for the existence of all these steps. But consider
        the disconnect between the intent "I want to fix a typo in a documentation file" and the
        actual steps. Defaults matter, and this default discourages shipping small changes. This
        effort is negligible for features that span days or weeks, but it is quite significant if we
        strive to "ship small and often".
      </p>
      <blockquote>
        This friction means that the path of least resistance is to grow the scope of code
        contributions and delaying integration.
      </blockquote>

      <h3>Traditional code reviews</h3>

      <p>
        Fixing a typo in the documentation does not carry the same risks as replacing an
        authentication middleware. So, why do we apply the same review process in both cases?
        Maintaining code quality and knowledge sharing are the two main reasons the software
        industry has adopted strict code reviews as a “best practice”. But as everything in
        engineering, there is an associated cost.
      </p>

      <blockquote>
        If I have to wait for review anyway, I might as well batch my typo fix together with this
        other thing I am building
      </blockquote>

      <p>
        The way formal code reviews are implemented today represents a hard blocking step, it
        contradicts the “shipping small and often” principle. This further skews the default
        behavior towards expanding the scope of code contributions.
      </p>

      <p>
        It takes a deliberate effort to keep contributions small. Some of the best engineers I have
        met create stacks of pull requests but this certainly takes some skill and focus to pull
        off.
      </p>
      <p>
        The question "How do developers contribute effectively their code to the same system?" has
        been on my mind since my days at uni. At this point, I have come to believe that the unit of
        code contribution is flawed, but never gets challenged because it is so ingrained.
      </p>

      <h2>The unit of contribution is flawed</h2>

      <p>
        Today, the term "pull request" is almost synonymous with the unit of contribution which is
        being reviewed and integrated. It's worth highlighting that this workflow originates from
        open-source, where more often than not, contributors don't know one another and there is no
        implicit trust. Teams coding at work, though, are quite different.
      </p>
      <p>
        Consider how different the work dynamics are — the developers are colleagues who know each
        other, plan together and likely have a daily sync and common chat channels. Development at
        work also happens at a much higher pace than open-source projects — with the median lead
        time to merge a pull request being about half as long.
      </p>

      <img src="./Img-7.png" alt="Median time until pull request merge" width="768" />

      <div class="text-right text-sm">
        Source:
        <a href="https://octoverse.github.com/writing-code-faster/#merging-pull-requests">GitHub</a>
      </div>

      <p>
        A friend of mine was recently telling me about the steps he takes to maintain a friendly
        collaboration vibe in the team. It boiled down to always having a discussion on Slack
        <i>before</i> submitting a pull request or leaving code feedback on another developer's PR.
        Somehow, in the pull request itself, the conversation was always a bit more tense and
        formal.
      </p>

      <p>
        Too often submitting pull requests feels like submitting homework and reviewing them feels
        like being a judge in a district court. These are the defaults that we work with, and it is
        up to individual teams to work out a process that makes collaboration feel more friendly.
      </p>

      <p>
        Perhaps this is the reason why many developers are choosing to avoid this process altogether
        in favor of
        <a href="https://martinfowler.com/articles/on-pair-programming.html">pair programming</a>.
        It fulfills the same objectives of quality gating and knowledge sharing in a non-blocking
        way.
      </p>

      <img src="./Img-8.png" alt="Twitter discussions on code reviews" width="768" />

      <h3>Does pair programming answer my question about effective collaboration on code?</h3>
      <p>
        Well, no, because it feels like cheating the semantics of the question. It is a form of
        <a href="https://en.wikipedia.org/wiki/Scalability#Vertical_or_scale_up">vertical scaling</a
        >, where the pair of developers acts as an extra-smart-developer. If my code collaboration
        puzzle was a question about a race condition, this solution is the equivalent of making the
        application single threaded.
      </p>

      <p>
        I would like to conclude with some thoughts and a hypothesis I have about making
        multithreaded code collaboration efficient and intuitive.
      </p>

      <h2>Code collaboration in teams via working in the open</h2>

      <p>
        Building on the strategy of integrating code contributions in small, incremental chunks is
        the idea of working in the open. This idea extends past the scope of what current generation
        development tools do (think Git platforms). There are three main elements to working in the
        open:
      </p>

      <ul>
        <li>
          Work-in-progress code is discoverable — the team can see what everyone is working on.
        </li>
        <li>Work-in-progress code is interactive — developers can easily try each other's code.</li>
        <li>Feedback is early, proactive and asynchronous.</li>
      </ul>

      <p>
        Practically speaking, this allows feedback and review steps to happen early in the process,
        rather than later on, when a lot of time and effort has been put in. This approach is
        similar to
        <a href="https://en.wikipedia.org/wiki/Shift-left_testing">shift-left testing</a>. Because
        good code discussions so timing and context dependent, an environment that provides
        discoverability of work-in-progress can allow the compounding of ideas.
      </p>

      <blockquote>Think of this as asynchronous pair programming.</blockquote>

      <p>
        Consider for a moment how tools like Figma and Google Docs have changed collaboration around
        design assets and text documents, respectively. During the same time period, the
        fundamentals of code collaborations have changed very little. Even if we view software
        development tools from a purely utilitarian point of view, it is clear that there is a
        mismatch between what they were designed to do and how they are used by teams today.
      </p>

      <h3>Tools are changing</h3>
      <p>
        These days, there are a number of hypotheses for the future of developer tools. From
        environments in the cloud
        <a href="https://gitpod.io/">[1]</a><a href="https://github.com/features/codespaces">[2]</a
        ><a href="https://aws.amazon.com/cloud9/">[3]</a> to pair programming tools
        <a href="https://git.live/">[4]</a
        ><a href="https://www.jetbrains.com/code-with-me/">[5]</a> and browser IDEs
        <a href="https://replit.com/">[6]</a>.
      </p>

      <blockquote>New tools should augment the existing ecosystem</blockquote>

      <p>
        It is my belief that tools should be built and adapted after how people already work, and
        not vice versa.
        <br />With <strong>trunk-based development</strong> and
        <strong>Continuous deployment</strong> are at the core of modern software development, what
        would a tool specifically designed for this workflow look like? If you would like to read
        our take on this, check out the
        <router-link :to="{ name: 'v2DocsRoot' }">Sturdy Docs</router-link>.
      </p>

      <p>
        So, how does a team <em>effectively</em> contribute code to the same system? I have come to
        believe that it boils down to continuous, high-quality communication around the code, and
        tools play an important role in this.
      </p>

      <p>
        Thanks for reading!
        <br />
        - <a href="https://twitter.com/krlvi">Kiril</a>
      </p>
      <hr />
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
