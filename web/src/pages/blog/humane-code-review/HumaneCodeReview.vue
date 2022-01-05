<template>
  <BlogPost
    title="Humane Code Review"
    subtitle="Suggestions on Sturdy"
    date="July 10, 2021"
    description="We’ve been thinking a lot about what it means to work on a team lately. Most often, developers
      are working on the same codebase, but are not really working together. More often
      than not the first time a teammate of your sees the code that you have written is when you’re
      considering the feature to be complete, and you’ve sent it for code review."
    :author="author"
  >
    <template #introduction>
      We’ve been thinking a lot about what it means to work on a team lately. Most often, developers
      are working on the same codebase, but are not really <em>working together</em>. More often
      than not the first time a teammate of your sees the code that you have written is when you’re
      considering the feature to be complete, and you’ve sent it for code review.
    </template>
    <template #default>
      <p>
        When doing other types of work such as working on a blog post or a specification, on let’s
        say Google Docs, it’s natural for the author to write an initial draft and then share it
        with a colleague. By default, your colleagues have edit access but can also opt for leaving
        suggestions, when to use an edit or a suggestion is left for their own judgement.
      </p>
      <p>
        In the world of writing software, there isn’t anything that is near as flexible and
        straightforward as a Google Doc.
      </p>

      <p>
        As we’re all painfully aware, software has semantic meaning in a way that natural language
        doesn’t (making a typo doesn’t crash your document), multi-cursor editing of the same
        document in real-time is only practical if you can synchronise on who’s working on what. If
        I try to compile the app while someone else is in the middle of writing something, the odds
        of a compilation error are very high. So we can’t code in a doc.
      </p>

      <h2>What are software engineers doing instead?</h2>
      <p>
        If we’re not sitting behind the same computer at the same time or using a multi-cursor pair
        programming tool, almost all of us are used to giving and receiving feedback using the slow
        and cumbersome process that is <em>Code Review</em>.
      </p>
      <p>
        Code Review is interesting in that the more feedback you give, the more work you transfer to
        the author. And the result of the Code Review is often a laundry list of code to write and
        commands to execute, and you as the reviewer can’t practically help the author to do it.
        This leads to a slow cycle where one engineer tells the other engineer what to type, almost
        word for word.
      </p>

      <p>
        In <a href="https://getsturdy.com">Sturdy</a>, we want to change how we collaborate on code
        as a team, and enable you to meaningfully <em>help</em> each other.
      </p>

      <h2>Today we’re releasing a new way to help your colleagues!</h2>

      <p>
        Where you can give feedback in way that’s more like <strong>working together</strong>, and
        where receiving feedback isn’t a burden.
      </p>

      <p>It’s called Suggestions, and allows you to literally give feedback with one click.</p>

      <p class="pl-4 border-l-4 border-solid border-yellow-500">
        Some context: when working with Sturdy, you connect a directory on your computer to a Sturdy
        codebase. This directory becomes fully managed by Sturdy, and you control it entirely from
        the web.
      </p>

      <hr />

      <p>
        Let’s see how it works! In this example, my co-founder Kiril is taking a look at one of my
        workspaces, where I’ve been working on a Fizz Buzz program.
      </p>
      <figure>
        <img class="w-full rounded-lg" src="./001.png" alt="" width="903" height="667" />
        <figcaption>Kiril’s point of view when he’s looking at my workspace</figcaption>
      </figure>

      <p>
        Kiril sees that the program only does the “fizz-buzzing” for the number 50<sup id="back-1"
          ><a href="#note-1">1</a></sup
        >, instead of for all numbers from 1 to 50 as we have planned for it to do.
      </p>

      <p>
        Kiril knows that I’m busy working on something else right now, and decides that he want’s to
        help out to complete the program.
      </p>

      <p>
        By clicking “Open on src/my-backend on …”, a complete copy of my working directory is made
        available to Kiril on his computer. Any unsaved changes that he might have had are
        automatically saved, and will be restored when he comes back to his own workspace.
      </p>

      <p>
        He’s now in “Suggesting” mode, and to make suggestions — all he has to do is to type them in
        his favourite editor (rumours say that it’s GoLand).
      </p>

      <figure class="text-center">
        <img
          class="w-max-[100%] rounded-lg block m-auto"
          src="./002.png"
          alt=""
          width="434"
          height="103"
        />
        <figcaption>Kiril seeing that he’s now suggesting</figcaption>
      </figure>

      <p>
        From my point of view, suggestions are appearing in real-time, and I can decide to take the
        suggested improvements if and when I want to.
      </p>

      <figure>
        <img class="w-full rounded-lg" src="./003.png" alt="" width="906" height="577" />
        <figcaption>My point of view after Kiril has implemented his suggestions</figcaption>
      </figure>

      <p>
        Taking the feedback is easy, I’ll click on Kiril’s avatar to see his suggestions, and from
        there, it’s just one click to take the feedback.
      </p>

      <figure>
        <img class="w-full rounded-lg" src="./004.png" alt="" width="898" height="365" />
        <figcaption>My point of view when previewing a suggestion</figcaption>
      </figure>

      <p>When I take the suggestion, the code updates instantly on my computer! 🤯</p>

      <figure>
        <img class="w-max-[100%] rounded-lg" src="./005.png" alt="" width="903" height="280" />
        <figcaption>After taking the suggestions</figcaption>
      </figure>

      <p>
        <strong>This works incredibly well!</strong> We can both take and give feedback at our own
        pace, while still making giving and receiving feedback lightweight and fast.
      </p>

      <p>
        After Kiril is done giving feedback, he’ll switch back to his own workspace, and his working
        directory will automatically be reset to how it was before.
      </p>

      <p>
        Giving feedback only takes <strong>one click</strong>, and it’s the same to take it — and
        neither one of us had to run any commands on the command line, or figure out how to pull a
        branch that’s non-fast-forwardable!
      </p>

      <p>
        Suggestions are <strong>available now</strong> to all teams and users on the Sturdy
        invite-only Beta.
      </p>

      <p>
        Suggestions also works great for large files, and many files at once, and you can pick and
        chose only the suggestions (or part of a suggestion) that you want to have, and discard the
        rest.
      </p>

      <p>
        For me, the ability to help each other out earlier in the development lifecycle is a natural
        continuation of the trend to Shift Left<sup id="back-2"><a href="#note-2">2</a></sup
        >, and we’re looking forward to see how this workflow will be embraced in development teams
        of all sizes.
      </p>

      <p>
        In the coming weeks we’ll ship more features related to early-stage collaboration on
        feedback, stay tuned!
      </p>

      <p>
        Thanks for reading,<br />
        Gustav
      </p>

      <hr />

      <p>
        Want to work at Sturdy (funded by Y Combinator), and change how developers work on code?
        <strong>We’re hiring in Stockholm!</strong>
      </p>
      <ol>
        <li>
          <a
            href="https://www.ycombinator.com/companies/sturdy/jobs/0ra7yLR-backend-software-engineer"
          >
            Backend Software Engineer
          </a>
        </li>
        <li>
          <a
            href="https://www.ycombinator.com/companies/sturdy/jobs/vc25j1l-frontend-software-engineer"
          >
            Frontend Software Engineer
          </a>
        </li>
      </ol>

      <hr />

      <ol>
        <li id="note-1">
          To save you the trouble: the output is “buzz”. <a href="#back-1" class="ml-2">↩</a>
        </li>
        <li id="note-2">
          <a href="https://en.wikipedia.org/wiki/Shift-left_testing">
            https://en.wikipedia.org/wiki/Shift-left_testing
          </a>
          <a href="#back-2" class="ml-2">↩</a>
        </li>
      </ol>
    </template>
  </BlogPost>
</template>
<script>
import BlogPost from '../BlogPost.vue'
import avatar from '../gustav.jpeg'

export default {
  components: { BlogPost },
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
