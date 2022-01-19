<template>
  <PublicLeftSidebar>
    <template #sidebar>
      <DocsSidebar />
    </template>
    <template #default>
      <div class="prose p-4 max-w-[800px]">
        <h1 id="how-sturdy-augments-git">How to collaborate with others</h1>
        <p>
          Nowadays it is so easy to work together on text documents (think Google Docs) or on design (think Figma).
          We think exchanging code ideas should be just as easy, if not easier.
          Code is, of course, different from "regular text".
          When we talk about collaborative coding we don't mean editing the same copy concurrently,
          but rather making the action of trying someone else's code and giving suggestions on it extremely fast and easy.
        </p>

        <p><a href="#receiving-feedback">Jump to: Receiving feedback or help ⤵️</a></p>
        <p><a href="#giving-feedback">Jump to: Giving feedback or help to others ⤵️</a></p>

        <h2 id="receiving-feedback">Receiving feedback or help</h2>

        <p>With Sturdy setup and connected to a directory on your computer:</p>
        <!-- TODO: Add a link to setup directory instructions -->
        <ol>
          <li>Create a new workspace and give it a good name</li>
          <li>Open the directory with the code in any text editor (tip: use the “Open in Visual Studio Code”-button)</li>
          <li>Make the changes to your code</li>
          <li>Go back to Sturdy, and navigate to the workspace from the sidebar</li>
          <li>Use the "Ask for feedback" dropdown menu to select a team member to ask for feedback.</li>
          <ol>
            <li>You can request feedback only from team members who are members of the codebase. Here is how you invite
              somebody to a codebase
            </li>
          </ol>
          <!-- TODO: Add a link to inviting somebody to your codebase -->
          <li>The user that you tagged will receive a notification that you asked for feedback</li>
          <ol>
            <li>When somebody from the team is checking out a workspace, their avatar will appear in the top right
              corner
            </li>
          </ol>
          <li><em>Optional:</em> Comment on code [<a href="#comment-on-code">go to</a>]</li>
          <li><em>Optional:</em> Taking / Dismissing code suggestions [link]</li>
        </ol>

        <h3 id="comment-on-code">Optional: Comment on code</h3>
        <p>You can write comments on specific piece of code within the workspace</p>
        <ol>
          <li>Hover over the specific line that you wish to comment on - a plus icon will appear right next to the line
            number</li>
          <li>Click on the plus button</li>
          <li>Write your message</li>
          <ol>
            <li>You can tag a team member by using @ and typing their name</li>
          </ol>
          <li>Click the blue "Comment" button or press Command + Return on Mac to create the comment.</li>
        </ol>
        <p>Comments on Sturdy are not attached to line numbers. This means that if a piece of code with a comment is moved to another place in the file, the comments will persist.</p>

        <h3 id="taking-dismissing-suggestions">Optional: Taking / Dismissing code suggestions</h3>
        <p>
          Other members of the codebase can get the code and create suggestions [link to giving feedback]. When that happens you will see a green banner indicating that at the top of the workspace.
        </p>
        <ol>
          <li>When you have new suggestions you will get a notification - you can see your notifications via the bell
            icon at the bottom left. If configured, you will also receive a notification via email.</li>
          <li>To see the suggestions, click the "Show suggestions" button. This will display the code diffs between the
            state of the code in your workspace and the changes made by the person suggesting.</li>
          <li>Suggestions are indicated with the avatar of the team member making them.</li>
          <li>You can hide the suggestions by clicking the "Hide suggestions" button in the green banner.</li>
          <li>Suggestions are not in your workspace until you accept them.</li>
          <ol>
            <li>To take a suggestion, press the green checkmark button on the left side. By taking the suggestion, the
              code in your workspace will be updated to reflect that - you can immediately see the update in your IDE as
              well.</li>
            <li>To dismiss a suggestion, press the red cross button on the left.</li>
          </ol>
          <li>All suggestions can be dismissed at once by pressing the cross button in the right of the banner</li>
        </ol>

        <h2 id="giving-feedback">Giving feedback or help to others</h2>
        <p>
          You can give help and feedback to others spontaneously or when somebody asks for feedback. In the left sidebar of Sturdy you can see your own workspaces as well as those of others on the team.
        </p>
        <ol>
          <li>Navigate to a workspace by somebody else on the team by clicking on it in the sidebar. The avatars to the
            left of workspace names indicate who the author is.</li>
          <li>If the author of the workspace is currently working on it, you will see their avatar in the top right
            corner.</li>
          <li>If the author of the workspace makes code changes, you will see them in real time, the same way you would
            see changes in your own workspaces.</li>
          <li><em>Optional:</em> Comment on code [<a href="#comment-on-code">go to</a>]</li>
          <li><em>Optional:</em> Quickly indicating feedback [<a href="##quick-feedback">go to</a>]</li>
          <li><em>Optional:</em> Trying the code [link]</li>
        </ol>

        <h3 id="#quick-feedback">Optional: Quickly indicating feedback</h3>
        <p>You can give quick feedback without spelling it out in text.</p>
        <ol>
          <li>To quickly indicate that code looks good, use the button with the thumbs up icon in the right sidebar.</li>
          <li>To quickly indicate that you have some important feedback, use the button with the (i) icon in the right
            sidebar. You would typically do that to indicate if you have spotted an error but don't have time right now
            to elaborate.</li>
        </ol>

        <h3 id="#trying-code">Optional: Trying the code</h3>
        <ol>
          <li>While viewing the workspaces of another team member, press the "Connect
            <em>folder</em> for suggesting" button.
          </li>
          <ol>
            <li>You will see a new nested workspace named "Suggestion 1" below the workspace that you are trying.</li>
            <li>This temporary workspace has the code of your team member's workspace as it looked at the time you
              pressed the connect button.</li>
          </ol>
          <li>The folder on your computer now has the code of your team member's workspace.</li>
          <li>Explore or run the code as you normally would using your IDE / text editor / terminal</li>
          <li><em>Optional:</em> Give code suggestions [<a href="##give-suggestions">go to</a>]</li>
          <li>When you are done you can navigate back to your own workspace and press the "Connect <em>folder</em> for
            editing" button to continue where you left off.</li>
        </ol>
        
        <p>
          If the original author makes changes to their workspace after you have pressed the "Connect <em>folder</em> for suggesting" you will not see those new changes.
          This is a known limitation. To work around this:
        </p>

        <ol>
          <li>Navigate to the workspace of your teammate</li>
          <li>Press the "Connect <em>folder</em> for suggesting" button</li>
          <ol>
            <li>Observe that now there is a second temporary workspace named "Suggestion 2" to which you are
              connected</li>
          </ol>
        </ol>

        <h3 id="#give-suggestions">Give code suggestions</h3>

        <p>
          When you have the code from someone else's workspace on your computer, you can give suggestions by simply making changes in your IDE.
        </p>
        <p>
          While you are connected to the "Suggestion <em>n</em>" workspace:
        </p>

        <ol>
          <li>Type any changes as you normally would, in your IDE or text editor.</li>
          <li>Within the temporary suggestion workspace you can see the changes you have created</li>
          <li>The author of the original workspace is notified of the suggestions [<a href="#taking-dismissing-suggestions">go to</a>]</li>
        </ol>

        <p>You can suggest changes even to files that have not been modified in the original workspace. </p>

      </div>
    </template>
  </PublicLeftSidebar>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import PublicLeftSidebar from '../../layouts/PublicLeftSidebar.vue'
import { useHead } from '@vueuse/head'
import DocsSidebar from '../../organisms/docs/DocsSidebar.vue'

export default defineComponent({
  components: { DocsSidebar, PublicLeftSidebar },
  setup() {
    // TODO: Remove when we're launching!
    useHead({
      meta: [
        {
          name: 'robots',
          content: 'noindex',
        },
      ],
    })
  },
})
</script>
