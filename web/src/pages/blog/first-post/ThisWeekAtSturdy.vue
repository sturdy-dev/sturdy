<template>
  <BlogPost
    title="#000 - March 18th, 2021"
    description="Hey everyone! We've been working hard on Sturdy for the past few months, and as you might have
      seen, we've been making quite a lot of changes. I'm super excited about where we're going, and
      to be able to show of Sturdy and some of it's internals here in this post."
    subtitle="This week at Sturdy"
    date="March 18, 2021"
    :author="author"
  >
    <template #introduction>
      Hey everyone! We've been working hard on Sturdy for the past few months, and as you might have
      seen, we've been making quite a lot of changes. I'm super excited about where we're going, and
      to be able to show of Sturdy and some of it's internals here in this post.
    </template>
    <template #default>
      <h2>Sturdy is replacing Git with real-time collaboration in the cloud</h2>
      <p>Let me tell you how it works!</p>
      <p>
        Sturdy mainly consists of two components, the "Sturdy Agent" that runs on your computer
        (more on that later!) and the Sturdy website.
        <strong
          >The website is used for everything that you need to do in your daily workflow</strong
        >, such as describing a change, reviewing a change, and sharing a change with the team. No
        command like skills or third party GUIs required!
      </p>
      <p>
        Unlike Git, Sturdy has no concept of remotes (or checkouts), instead all developers on a
        project are working closer together. Developers are still working separately in what we're
        calling a workspace (there are no "shared cursors", which you might be used to in Google
        Docs or Replit). Your changes are synced with the changes of other developers only when
        you’re ready.
      </p>
      <p>
        The Sturdy Agent automatically syncs a directory on your computer with Sturdy, when you're
        making a change it will be automatically synced within a few milliseconds.
      </p>
      <div class="captioned-image-container">
        <figure>
          <a
            class="image-link image2 image2-710-1456"
            target="_blank"
            href="https://cdn.substack.com/image/fetch/f_auto,q_auto:good,fl_progressive:steep/https%3A%2F%2Fbucketeer-e05bbc84-baa3-437e-9518-adb32be77984.s3.amazonaws.com%2Fpublic%2Fimages%2Ff2972ac1-0f1e-46d2-8085-140c4cca944e_1970x960.png"
            ><img
              src="https://bucketeer-e05bbc84-baa3-437e-9518-adb32be77984.s3.amazonaws.com/public/images/f2972ac1-0f1e-46d2-8085-140c4cca944e_1970x960.png"
              data-attrs='{"src":"https://bucketeer-e05bbc84-baa3-437e-9518-adb32be77984.s3.amazonaws.com/public/images/f2972ac1-0f1e-46d2-8085-140c4cca944e_1970x960.png","height":710,"width":1456,"resizeWidth":null,"bytes":187966,"alt":null,"title":null,"type":"image/png","href":null}'
              alt=""
          /></a>
        </figure>
      </div>
      <blockquote>
        <p>Pick the changes that you want to share with the team - in just a few seconds!</p>
      </blockquote>
      <h2>The Sturdy agent (the part on your computer)</h2>
      <p>
        Sturdy runs a lightweight process on your computer that watches the file system for changes,
        and syncs it with Sturdy over a SSH-based protocol. If you make any changes to Sturdy via
        the web, such as switching to a different workspace, Sturdy automatically swaps out the
        files on your filesystem within seconds.
      </p>
      <p>Sturdy does not require any permissions on your computer, and is easy to setup.</p>
      <h2>Changes since last time</h2>
      <h4>1. Conflict resolver</h4>
      <p>
        Conflicts can happen to the best of us, and can happen even when using Sturdy. When a
        conflict can’t automatically be resolved, a human needs to step in and make a decision. The
        conflict resolver is now available for everyone, and is looking great!
      </p>
      <div class="captioned-image-container">
        <figure>
          <a
            class="image-link image2 image2-592-1456"
            target="_blank"
            href="https://cdn.substack.com/image/fetch/f_auto,q_auto:good,fl_progressive:steep/https%3A%2F%2Fbucketeer-e05bbc84-baa3-437e-9518-adb32be77984.s3.amazonaws.com%2Fpublic%2Fimages%2F9420b9ed-ca70-49cc-a9b3-32059a92e49f_1998x812.png"
            ><img
              src="https://bucketeer-e05bbc84-baa3-437e-9518-adb32be77984.s3.amazonaws.com/public/images/9420b9ed-ca70-49cc-a9b3-32059a92e49f_1998x812.png"
              data-attrs='{"src":"https://bucketeer-e05bbc84-baa3-437e-9518-adb32be77984.s3.amazonaws.com/public/images/9420b9ed-ca70-49cc-a9b3-32059a92e49f_1998x812.png","height":592,"width":1456,"resizeWidth":null,"bytes":153923,"alt":null,"title":null,"type":"image/png","href":null}'
              alt=""
          /></a>
        </figure>
      </div>
      <h4>2. Create a workspace in the past</h4>
      <p>
        From the changelog of your codebase, use the new “New workspace” button to create a
        workspace based on this revision. This is useful if you want to go back in time and see what
        a codebase used to look like.
      </p>
      <div class="captioned-image-container">
        <figure>
          <a
            class="image-link image2 image2-167-1456"
            target="_blank"
            href="https://cdn.substack.com/image/fetch/f_auto,q_auto:good,fl_progressive:steep/https%3A%2F%2Fbucketeer-e05bbc84-baa3-437e-9518-adb32be77984.s3.amazonaws.com%2Fpublic%2Fimages%2F1f9aa3c8-6876-49d5-a992-3a89b6c6bcef_2000x230.png"
            ><img
              src="https://bucketeer-e05bbc84-baa3-437e-9518-adb32be77984.s3.amazonaws.com/public/images/1f9aa3c8-6876-49d5-a992-3a89b6c6bcef_2000x230.png"
              data-attrs='{"src":"https://bucketeer-e05bbc84-baa3-437e-9518-adb32be77984.s3.amazonaws.com/public/images/1f9aa3c8-6876-49d5-a992-3a89b6c6bcef_2000x230.png","height":167,"width":1456,"resizeWidth":null,"bytes":44736,"alt":null,"title":null,"type":"image/png","href":null}'
              alt=""
          /></a>
        </figure>
      </div>
      <h4>3. Easy open of a workspace on an existing connection</h4>
      <p>
        It’s now easier than ever to switch between workspaces. If you’re creating a new workspace
        (or checking someone elses workspace out), use the “Use on…” button to open the workspace on
        a connected computer. And if you’re opening a workspace that isn’t used anywhere, a nice and
        helpful guide will let you open the workspace in no-time.
      </p>
      <div class="captioned-image-container">
        <figure>
          <a
            class="image-link image2 image2-365-1456"
            target="_blank"
            href="https://cdn.substack.com/image/fetch/f_auto,q_auto:good,fl_progressive:steep/https%3A%2F%2Fbucketeer-e05bbc84-baa3-437e-9518-adb32be77984.s3.amazonaws.com%2Fpublic%2Fimages%2Fa3aa352c-678d-4d7e-9c67-8fb4cda2dfe6_1978x496.png"
            ><img
              src="https://bucketeer-e05bbc84-baa3-437e-9518-adb32be77984.s3.amazonaws.com/public/images/a3aa352c-678d-4d7e-9c67-8fb4cda2dfe6_1978x496.png"
              data-attrs='{"src":"https://bucketeer-e05bbc84-baa3-437e-9518-adb32be77984.s3.amazonaws.com/public/images/a3aa352c-678d-4d7e-9c67-8fb4cda2dfe6_1978x496.png","height":365,"width":1456,"resizeWidth":null,"bytes":72392,"alt":null,"title":null,"type":"image/png","href":null}'
              alt=""
          /></a>
        </figure>
      </div>
      <p>
        That’s it for today! We hope that you like what you’re seeing. You can sign up to join the
        waiting list at <a href="https://getsturdy.com">getsturdy.com</a>, and I’ll send you an
        invite soon.
      </p>
      <p>Thanks,<br />Gustav</p>
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
