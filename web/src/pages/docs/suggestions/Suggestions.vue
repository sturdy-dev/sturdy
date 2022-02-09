<template>
  <StaticPage
    title="Suggestions"
    subtitle="Learn how to use Suggestions in Sturdy"
    metadescription="Learn how to use Suggestions in Sturdy"
    image="https://images.unsplash.com/photo-1516321318423-f06f85e504b3?ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&ixlib=rb-1.2.1&auto=format&fit=crop&w=2070&q=80"
  >
    <div class="prose prose-yellow">
      <p>
        Sturdy's workspaces gives you real-time superpowers. All of the code that you create in a
        workspace is instantly available for your team to review and feedback on.
      </p>

      <p>
        Workspaces are powerful and allow you to comment and collaborate in real time. In this
        documentation we'll go through how to use
        <strong>Suggestions</strong>, to give feedback and suggest changes to a workspace &mdash;
        directly from your editor or IDE.
      </p>

      <ul>
        <li>Share feedback directly from your editor</li>
        <li>Explore any workspace from your IDE, with powerful code navigation</li>
        <li>Give and take feedback, without leaving your editor or browser</li>
        <li>No copy-pasting, and no commands to run on the command line</li>
      </ul>

      <p>
        Let's setup a project, invite my friend Thomas, write some code, and take a look at how to
        use suggestions!
      </p>

      <a :href="thomasStartImage" target="_blank">
        <img :src="thomasStartImage" height="2096" width="3562" />
      </a>

      <p>
        Here we are. Gustav has created a <code>fib()</code> function, but it has a slight bug (the
        result is never printed! 😱).
      </p>

      <p>
        Thomas, reviews this code, and spots the bug. Instead of leaving a comment, that would
        instruct Gustav how to fix the bug. Thomas decides to leave a <em>suggestion</em> with the
        fix instead. He does this by pressing the green
        <strong>Open for suggesting on ...</strong> button.
      </p>

      <p>
        With this single click, all files from Gustav's workspace are transferred to Thomas
        computer. If Thomas had any unsaved changes, they would be saved and will be automatically
        restored when he goes back to his own workspace.
      </p>

      <a :href="thomasSuggestingImage" target="_blank">
        <img :src="thomasSuggestingImage" height="2096" width="3562" />
      </a>

      <p>
        Any changes that Thomas makes in his editor, are automatically sent as suggestions to
        Gustav. No action required!
      </p>

      <a :href="gustavTakeSuggestionImage">
        <img :src="gustavTakeSuggestionImage" height="2096" width="3562" />
      </a>
      <p>
        Gustav now gets a notification that Thomas has made suggestions to his workspace, and can
        decide to either take the suggestion, or dismiss it.
      </p>

      <p>
        By clicking <strong>Take suggestion</strong>, the suggestions from Thomas can be instantly
        applied to Gustav's workspace, and are downloaded to his computer.
      </p>

      <h3>Read more</h3>

      <ul>
        <li>
          <router-link :to="{ name: 'blogVision' }">
            Unbreaking code collaboration &mdash; Our vision for the future
          </router-link>
        </li>
        <li>
          <router-link :to="{ name: 'resourcesDocs' }">Documentation</router-link>
        </li>
      </ul>
    </div>
  </StaticPage>
</template>

<script>
import StaticPage from '../../../layouts/StaticPage.vue'

import thomasStartImage from './thomas_start.png'
import thomasSuggestingImage from './thomas_suggesting.png'
import gustavTakeSuggestionImage from './gustav_take_suggestion.png'

export default {
  name: 'Suggestions',
  components: { StaticPage },
  setup() {
    return {
      thomasStartImage,
      thomasSuggestingImage,
      gustavTakeSuggestionImage,
    }
  },
}
</script>
