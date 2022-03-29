<template>
  <BlogPost
    title="ACLs, and a fresh hot look"
    subtitle="This week at Sturdy"
    date="September 29, 2021"
    description="Sturdy changelog. Sturdy now ships with Access Control Lists and a brand new interface, that's faster and more lightweight than ever."
    reading-time="3 min"
    :author="author"
  >
    <template #introduction>
      Hi there! 👋 This issue of the changelog includes many highlights that we're super happy to
      announce. I hope that you'll like it!
    </template>
    <template #default>
      <h2>🔥 Hot design!</h2>
      <p>
        Many major areas of Sturdy have been redesigned with a lighter theme, that highlights the
        most important areas of Sturdy when you need them. Let's take a look!
      </p>
      <h3>1️⃣ Workspace</h3>
      <img
        src="./newsletter-workspace.png"
        height="2620"
        width="4000"
        alt="Screenshot of a Sturdy Workspace"
      />
      <p>
        The workspace now puts the code front and center. Distracting components like the changelog
        are now gone, and have made space for <em>feedback</em>, and the <em>new activity feed</em>.
      </p>
      <p>
        Use <strong>feedback</strong> to ask for feedback on your changes from a peer before you
        change them. Asking for feedback will send a notification to the person that you want to get
        feedback from, letting them know that you want their help.
      </p>
      <p>
        When they are done, they can either leave a thumbs up, or let you know that they have some
        feedback that they would like you to incorporate before you go ahead and share your changes.
      </p>
      <h3>2️⃣ Codebase</h3>
      <img
        src="./newsletter-codebase-readme.png"
        height="2620"
        width="4000"
        alt="Screenshot of a Sturdy Codebase with a README file"
      />
      <p>
        The home of your codebase has also been redesigned to match, with a lighter theme that's
        more focused on showing you what you need. The page is now featuring the codebases
        <code>README.md</code>, and the first-time setup instructions have been greatly improved.
      </p>
      <h3>3️⃣ Changelog</h3>
      <img
        src="./newsletter-codebase-changelog.png"
        height="2620"
        width="4000"
        alt="Screenshot of the changelog of a codebase on Sturdy"
      />
      <p>
        The changelog is now much faster, loading new data in just a few milliseconds. The page has
        also been redesigned to fit with the new colours and layout. From the changelog, you can
        find the complete history of your project, just like you're used to.
      </p>
      <p>
        From the changelog, you can now also <strong>view the entire file</strong> as it currently
        looks like on trunk.
      </p>
      <p>And last but, definitely not least. <strong>The new sidebar!</strong></p>
      <h3>4️⃣ Sidebar</h3>
      <p>
        I'm really stoked about how the new sidebar looks and feels. We've been taking inspiration
        from some of our favourite real-time collaboration tools, and are brining the same
        experience to coding.
      </p>
      <p>
        The sidebar organizes all of your codebases and workspaces, and makes it easy to get an
        overview of what's going on. And the sidebar will intelligently inform you where your
        attention is needed and where there has been new activity.
      </p>
      <h2 class="items-center flex">
        🔐 Access Control
        <Pill class="ml-4">BETA</Pill>
      </h2>
      <img
        src="../../docs/access-control/access-control.png"
        width="1552"
        height="1024"
        alt="Screenshot of the Access Control List permission editor on Sturdy"
      />
      <p>
        We're now launching the beta of
        <router-link :to="{ name: 'docsAccessControl' }">
          <strong>Access Control Lists</strong>
        </router-link>
        in Sturdy. ACLs in Sturdy are really powerful, and enables you to give fine grained access
        to codebases, <em>including to files</em>.
      </p>
      <p>
        The file ACL allows you to control who on your team can see what files. For example if you
        want to invite consultants to collaborate on only a part of your project. Or if you want to
        store sensitive files in your project, that only some team members should have access to.
      </p>
      <p>
        The permissions for <code>files</code> resources affects both which files that are synced to
        your workstation, and which files that you can see on the web.
      </p>
      <p>
        Together with Sturdy's integration with GitHub, Sturdy can now also operate as a access
        control-proxy for your existing projects on git.
      </p>
      <p>
        The
        <router-link :to="{ name: 'docsAccessControl' }">
          <strong>ACL documentation</strong>
        </router-link>
        contains more examples for how you can use ACLs in your project.
      </p>
      <p><em>To use ACLs, you need to install Sturdy v0.7.0 or newer.</em></p>

      <h2>🔔 Notifications</h2>
      <p>
        Notification in Sturdy has received a huge upgrade. Workspaces that need your attention are
        now <strong>highlighted in the sidebar</strong>, so that you get an easy overview of what
        has changed since the last time you where online.
      </p>
      <p>
        Sturdy now also supports <strong>browser notifications</strong>, to notify you immediately
        when something has changed. When you click the bell icon (lower left of the screen) the next
        time, Sturdy will ask you for your permission to send you notifications.
      </p>
      <h2>⤴️ GitHub integration: PRs on your behalf</h2>
      <p>
        When using Sturdy and GitHub, Sturdy is now creating Pull Requests under <em>your</em> name
        instead of our bot account. This means that you now get the credit of your work that you
        deserve!
      </p>
      <p>
        If you're curious, here's what a real world PR created from Sturdy looks like:
        <a href="https://github.com/zegl/kube-score/pull/393">kube-score#393</a>.
      </p>
      <h2>🕸 API</h2>
      <p>
        All of the changes above, are available Sturdy's
        <a href="https://getsturdy.com/api">GraphQL API</a>.
      </p>
      <h2>🖥 Sturdy Client</h2>
      <p>The Homebrew installation now supports darwin/arm64 and linux/arm64 out of the box! :-)</p>
      <hr />

      <p>Thanks for reading, see you in cyberspace!<br />&mdash; Gustav</p>

      <ul>
        <li><a href="https://discord.gg/5HnSdzMqtA">Join our Discord</a></li>
        <li><a href="https://twitter.com/getsturdy">Follow us on Twitter</a></li>
      </ul>
    </template>
  </BlogPost>
</template>

<script lang="ts" setup>
import BlogPost from '../BlogPost.vue'
import avatar from '../gustav.jpeg'
import Pill from '../../../atoms/Pill.vue'

const author = {
  name: 'Gustav Westling',
  avatar: avatar,
  link: 'https://twitter.com/zegl',
}
</script>
