<template>
  <BlogPost
    title="YC W21 Demo Day!"
    subtitle="Behind the scenes"
    date="March 24, 2021"
    :author="author"
  >
    <template #introduction>
      The past week has been the culmination of our YC W21 batch, with Demo Day1 on Tuesday! A huge
      congratulations and thanks to all other companies in the batch!
    </template>
    <template #default>
      <p>
        The past week has been the culmination of our YC W21 batch, with Demo Day<sup
          ><a id="footnote-anchor-1" class="footnote-anchor" href="#footnote-1">1</a></sup
        >
        on Tuesday! A huge congratulations and thanks to all other companies in the batch!
      </p>
      <p>
        W21 is the second batch in a row to be completely virtual, including Demo Day. This is what
        it looked like when Kiril went on stage:
      </p>
      <div class="captioned-image-container">
        <figure>
          <a
            class="image-link image2 image2-428-1456"
            target="_blank"
            href="https://cdn.substack.com/image/fetch/f_auto,q_auto:good,fl_progressive:steep/https%3A%2F%2Fbucketeer-e05bbc84-baa3-437e-9518-adb32be77984.s3.amazonaws.com%2Fpublic%2Fimages%2F9214f484-995f-4435-943c-9857c442874d_1756x516.png"
            ><img
              src="https://bucketeer-e05bbc84-baa3-437e-9518-adb32be77984.s3.amazonaws.com/public/images/9214f484-995f-4435-943c-9857c442874d_1756x516.png"
              data-attrs='{"src":"https://bucketeer-e05bbc84-baa3-437e-9518-adb32be77984.s3.amazonaws.com/public/images/9214f484-995f-4435-943c-9857c442874d_1756x516.png","height":428,"width":1456,"resizeWidth":null,"bytes":433132,"alt":null,"title":null,"type":"image/png","href":null}'
              alt=""
          /></a>
        </figure>
      </div>
      <p>
        Here’s the other point of view, with Kirils makeshift teleprompter, which is one of the best
        hacks that I’ve ever seen:
      </p>
      <div class="captioned-image-container">
        <figure>
          <a
            class="image-link image2 image2-535-1024"
            target="_blank"
            href="https://cdn.substack.com/image/fetch/f_auto,q_auto:good,fl_progressive:steep/https%3A%2F%2Fbucketeer-e05bbc84-baa3-437e-9518-adb32be77984.s3.amazonaws.com%2Fpublic%2Fimages%2Fe9b05dd7-980b-442f-b06f-f4fd6836d361_1024x535.jpeg"
            ><img
              src="https://bucketeer-e05bbc84-baa3-437e-9518-adb32be77984.s3.amazonaws.com/public/images/e9b05dd7-980b-442f-b06f-f4fd6836d361_1024x535.jpeg"
              data-attrs='{"src":"https://bucketeer-e05bbc84-baa3-437e-9518-adb32be77984.s3.amazonaws.com/public/images/e9b05dd7-980b-442f-b06f-f4fd6836d361_1024x535.jpeg","height":535,"width":1024,"resizeWidth":null,"bytes":215795,"alt":null,"title":null,"type":"image/jpeg","href":null}'
              alt=""
          /></a>
        </figure>
      </div>

      <h3>Changes since last week</h3>
      <p>
        Demo Day didn’t stop us from continuing to ship improvements to Sturdy. This week we’ve
        focused on the collaborative aspects, and ease of use in the core workflows.
      </p>
      <p>Here are some of the highlights!</p>
      <h4>Sync changes</h4>
      <p>
        Our ambition for Sturdy is that you should never have to worry about making sure that the
        code on your computer is up to date before starting a new piece of work.
      </p>
      <p>
        We have not yet nailed that experience, and are currently falling back to making it obvious,
        and extremely easy to stay up to date. If someone else has made changes that are not yet in
        your workspace, the “Sync changes” button will start flashing. From there, syncing is
        instantaneous and requires no extra steps.
      </p>
      <p>
        “Sync changes” was previously known as “Sync on Trunk”, as the workflow in Sturdy is largely
        inspired by Trunk Based Development<sup
          ><a id="footnote-anchor-2" class="footnote-anchor" href="#footnote-2">2</a></sup
        >. This reference is not very clear to newcomers, so we’ve simply renamed it to Sync
        changes. :-)
      </p>
      <div class="captioned-image-container">
        <figure>
          <a
            class="image-link image2 image2-62-316"
            target="_blank"
            href="https://cdn.substack.com/image/fetch/f_auto,q_auto:good,fl_progressive:steep/https%3A%2F%2Fbucketeer-e05bbc84-baa3-437e-9518-adb32be77984.s3.amazonaws.com%2Fpublic%2Fimages%2Fce7ed6b9-9ac4-43a6-8247-931970148f31_600x116.gif"
            ><img
              src="https://bucketeer-e05bbc84-baa3-437e-9518-adb32be77984.s3.amazonaws.com/public/images/ce7ed6b9-9ac4-43a6-8247-931970148f31_600x116.gif"
              data-attrs='{"src":"https://bucketeer-e05bbc84-baa3-437e-9518-adb32be77984.s3.amazonaws.com/public/images/ce7ed6b9-9ac4-43a6-8247-931970148f31_600x116.gif","height":116,"width":600,"resizeWidth":316,"bytes":96926,"alt":null,"title":null,"type":"image/gif","href":null}'
              alt=""
          /></a>
        </figure>
      </div>
      <h4><strong>Picking multiple changes got easier!</strong></h4>
      <p>
        You can now pick a directory in the file picker, and all files under that directory will be
        added or unadded in a single click — easy peasy!
      </p>
      <div class="captioned-image-container">
        <figure>
          <a
            class="image-link image2 image2-500-800"
            target="_blank"
            href="https://cdn.substack.com/image/fetch/f_auto,q_auto:good,fl_progressive:steep/https%3A%2F%2Fbucketeer-e05bbc84-baa3-437e-9518-adb32be77984.s3.amazonaws.com%2Fpublic%2Fimages%2F2f748721-a1a6-40e9-9f59-1b1b284c4062_800x500.gif"
            ><img
              src="https://bucketeer-e05bbc84-baa3-437e-9518-adb32be77984.s3.amazonaws.com/public/images/2f748721-a1a6-40e9-9f59-1b1b284c4062_800x500.gif"
              data-attrs='{"src":"https://bucketeer-e05bbc84-baa3-437e-9518-adb32be77984.s3.amazonaws.com/public/images/2f748721-a1a6-40e9-9f59-1b1b284c4062_800x500.gif","height":500,"width":800,"resizeWidth":null,"bytes":374951,"alt":null,"title":null,"type":"image/gif","href":null}'
              alt=""
          /></a>
        </figure>
      </div>
      <h4>Seeing a teammates work</h4>
      <p>
        When looking at someone else’s workspace, you’ll now be able to see what’s going on in
        there! Sturdy now enables you to give your teammates very early feedback, without having to
        ask them to push or screenshare. Going forward, we’re going to add collaborative tools right
        here in the workspace page, so that you can help out creating changes, or even fixing bugs
        before landing a change!
      </p>
      <div class="captioned-image-container">
        <figure>
          <a
            class="image-link image2 image2-658-863"
            target="_blank"
            href="https://cdn.substack.com/image/fetch/f_auto,q_auto:good,fl_progressive:steep/https%3A%2F%2Fbucketeer-e05bbc84-baa3-437e-9518-adb32be77984.s3.amazonaws.com%2Fpublic%2Fimages%2F69c3ab56-f50f-446c-a524-fc84093f7e50_863x658.png"
            ><img
              src="https://bucketeer-e05bbc84-baa3-437e-9518-adb32be77984.s3.amazonaws.com/public/images/69c3ab56-f50f-446c-a524-fc84093f7e50_863x658.png"
              data-attrs='{"src":"https://bucketeer-e05bbc84-baa3-437e-9518-adb32be77984.s3.amazonaws.com/public/images/69c3ab56-f50f-446c-a524-fc84093f7e50_863x658.png","height":658,"width":863,"resizeWidth":null,"bytes":59524,"alt":null,"title":null,"type":"image/png","href":null}'
              alt=""
          /></a>
        </figure>
      </div>
      <h2>Want to try Sturdy?</h2>
      <p>
        We’re up Sturdy to more and more users, if you want an invite,
        <a href="https://getsturdy.com/">request access from our landing page</a>, and we’ll send
        you an invite soon!<br />
      </p>
      <p>Thanks for reading, see you soon again!<br />/ Gustav</p>
      <hr />

      <ol>
        <li id="footnote-1" class="footnote">
          <a href="https://www.ycombinator.com/demoday">https://www.ycombinator.com/demoday</a>
          <a href="#footnote-anchor-1" class="ml-2">↩</a>
        </li>
        <li id="footnote-2" class="footnote">
          <a href="https://trunkbaseddevelopment.com/">https://trunkbaseddevelopment.com/</a>
          <a href="#footnote-anchor-2" class="ml-2">↩</a>
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
