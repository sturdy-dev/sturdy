<template>
  <BlogPost
    title="Restore to any point in time, with 1-click"
    subtitle="Feel free to make mistakes, Sturdy got your back"
    date="April 1, 2021"
    :author="author"
  >
    <template #introduction>
      Hey everyone! This week we’ve been focusing on how Sturdy can help you if you make a mistake,
      or change your mind about what you’re building.
    </template>
    <template #default>
      <p>
        <strong>We’re introducing Snapshots</strong>, which allows you to time-travel to any
        previous point in history, similar to the history of a Google Docs document, or Apple’s Time
        Machine!
      </p>

      <p>
        With Sturdy it’s now safer than ever to do experiments or crazy refactoring, without having
        to worry about how to get back to a working state if your experiment doesn’t go according to
        plan — and according to me, this is the coolest feature of Sturdy so far!
      </p>

      <h2>Snapshots</h2>
      <p>
        It’s OK to make mistakes, for me it’s a core part of the development process. I write some
        code to experiment with a library, or an idea. The first attempt might not work out, so I
        back out, remove some of my code and try again! This is what the development process looks
        like for me many-many times per day.
      </p>

      <p>
        Sturdy now fully embraces this workflow, and allows you to switch back and forth between
        stages of your code, with what we’re calling <strong>Snapshots</strong>. Snapshots are fully
        automatic, and allows you to
        <strong>restore the workspace to what it looked like at any point in the past</strong>.
      </p>

      <p>I was previously hinting this on Twitter, and it’s now available for everyone!</p>

      <blockquote class="twitter-tweet">
        <p lang="en" dir="ltr">
          I&#39;ve been experimenting with automatic snapshots in Sturdy, with 1-click restore to
          any point in time. 👀 <a href="https://t.co/8Z3FU0hLLj">pic.twitter.com/8Z3FU0hLLj</a>
        </p>
        &mdash; Gustav Westling (@zegl)
        <a href="https://twitter.com/zegl/status/1376952939445690376?ref_src=twsrc%5Etfw"
          >March 30, 2021</a
        >
      </blockquote>

      <p>This is what it looks like (real screenshot, of me building the feature itself)!</p>

      <figure class="text-center">
        <img
          class="w-max-[100%] rounded-lg block m-auto"
          src="./001.png"
          alt=""
          width="895"
          height="672"
        />
      </figure>

      <p>I’m excited to see how Snapshots will improve everyones workflows!</p>

      <h2>Undo</h2>

      <h2>Sharing a change with the team</h2>
      <p>
        On the topic of making mistakes, after you’ve created a change (the equivalent of a
        “commit”), you can now easily undo it, if you for example change your mind or if you added
        to many or too few changes to it. Easy! Undoing takes you back to exactly where you where
        before the change was made.
      </p>

      <figure class="text-center">
        <img
          class="w-max-[100%] rounded-lg block m-auto"
          src="./002.png"
          alt=""
          width="978"
          height="213"
        />
      </figure>

      <h2>Smaller news</h2>
      <p>Codebases can now be archived!</p>

      <figure class="text-center">
        <img
          class="w-max-[100%] rounded-lg block m-auto"
          src="./003.png"
          alt=""
          width="978"
          height="254"
        />
      </figure>

      <p>
        Changes that where introduced by reverting an earlier change, now have a link to the change
        that they are reverting.
      </p>

      <figure class="text-center">
        <img
          class="w-max-[100%] rounded-lg block m-auto"
          src="./004.png"
          alt=""
          width="978"
          height="268"
        />
      </figure>

      <p>
        If you want to use Sturdy, sign up to the waiting list at
        <a href="https://getsturry.com/">getsturdy.com</a>, we’re accepting more and more users
        every week.
      </p>

      <p>
        Thanks everyone for reading, see you next week!<br />
        / Gustav
      </p>
    </template>
  </BlogPost>
</template>

<script lang="ts" setup>
import BlogPost from '../BlogPost.vue'
import { useHead } from '@vueuse/head'
import avatar from '../gustav.jpeg'

useHead({
  script: [
    {
      async: true,
      src: 'https://platform.twitter.com/widgets.js',
      charset: 'utf-8',
    },
  ],
})

const author = {
  name: 'Gustav Westling',
  avatar: avatar,
  link: 'https://twitter.com/zegl',
}
</script>
