<template>
  <DocumentationWithTableOfContents>
    <template #sidebar>
      <DocsSidebar />
    </template>
    <template #default>
      <div class="prose p-4 max-w-[800px]">
        <h1 id="how-to-collaborate">How to collaborate with others</h1>
        <p>
          Nowadays it is so easy to work together on text documents
          <SubtleText>(think Google Docs)</SubtleText>
          or on design
          <SubtleText>(think Figma)</SubtleText>
          . We think exchanging code ideas should be
          <strong>just as easy</strong>, if not easier!
        </p>
        <p>
          Code is, of course, different from "regular text". When we talk about collaborative coding
          we don't mean editing the same copy concurrently, but rather making the action of trying
          someone else's code and giving suggestions on it extremely fast and easy.
        </p>

        <p><a href="#receiving-feedback">Jump to: Receiving feedback or help ⤵️</a></p>
        <p><a href="#giving-feedback">Jump to: Giving feedback or help to others ⤵️</a></p>

        <h2 id="receiving-feedback">Receiving feedback or help</h2>

        <p>
          With Sturdy setup and
          <router-link :to="{ name: 'v2DocsHowToEditCode' }">connected to a directory</router-link>
          on your computer:
        </p>
        <ol>
          <li>Create a new draft change and give it a good name</li>
          <li>
            Open the directory with the code in any text editor
            <SubtleText>(tip: use the “Open in Visual Studio Code”-button)</SubtleText>
          </li>
          <li>Make the changes to your code</li>
          <li>Go back to Sturdy, and navigate to the draft from the sidebar</li>
          <li>
            Use the <strong>"Ask for feedback"</strong> dropdown menu to the right to select a team
            member to ask for feedback
          </li>
          <ol>
            <li>
              You can request feedback only from team members who are members of the organization or
              codebase
            </li>
          </ol>
          <!-- TODO: Add a link to inviting somebody to your organization or codebase -->
          <li>The user that you tagged will receive a notification that you asked for feedback</li>
          <ol>
            <li>
              When somebody from the team is browsing a draft, their avatar will appear in the top
              right corner
            </li>
          </ol>
          <li>
            <SubtleText
              ><em>Optional:</em> Comment on code [<a href="#comment-on-code">go to</a>]
            </SubtleText>
          </li>
          <li>
            <SubtleText
              ><em>Optional:</em> Taking / Dismissing code suggestions [<a
                href="#taking-dismissing-suggestions"
                >go to</a
              >]
            </SubtleText>
          </li>
        </ol>

        <h3 id="comment-on-code" title="Comment on code">Optional: Comment on code</h3>
        <p>You can write comments on specific piece of code within the draft change</p>
        <ol>
          <li>
            Point the mouse over the specific line that you wish to comment on &mdash; a plus icon
            will appear to the right next to the line number
          </li>
          <li>
            Click on the plus button
            <PlusIcon class="h-4 w-4 inline" />
          </li>
          <li>Write your message</li>
          <ol>
            <li>
              <SubtleText>
                Tip: You can tag a team member by typing <code>@</code> followed by their name
              </SubtleText>
            </li>
            <li>
              <SubtleText>
                Tip: Write an emoji by typing <code>:</code> followed by the emoji name!
              </SubtleText>
            </li>
          </ol>
          <li>
            Click the blue "Comment" button or press <code>⌘ + enter</code> or
            <code>ctrl + enter</code> to send the comment.
          </li>
        </ol>

        <DocsInfoBox>
          Comments on Sturdy are not attached to line numbers. A comment is attached to the piece of
          code that the comment was made on, if the code moves, so will the comment!
        </DocsInfoBox>

        <h3 id="taking-dismissing-suggestions" title="Taking / Dismissing code suggestions">
          Optional: Taking / Dismissing code suggestions
        </h3>
        <p>
          Other members of the codebase can get the code from your draft and
          <a href="#give-suggestions">create suggestions for you</a>. When that happens you will see
          a green banner indicating it at the top of the draft.
        </p>
        <ol>
          <li>
            When you have new suggestions you will get a notification &mdash; you can see your
            notifications via the bell icon
            <BellIcon class="h-4 w-4 inline" />
            at the bottom left of Sturdy. If configured, you will also receive a notification via
            email
            <SubtleText>
              (email notifications are only available for Sturdy in the Cloud)
            </SubtleText>
            .
          </li>
          <li>
            To see the suggestions, click the "Show suggestions" button. This will display the code
            diffs between the state of the code in your draft and the changes made by the person
            suggesting.
          </li>
          <li>Suggestions are indicated with the avatar of the team member making them.</li>
          <li>
            You can hide the suggestions by clicking the "Hide suggestions" button in the green
            banner.
          </li>
          <li>Suggestions are not in your draft change until you accept them.</li>
          <ol>
            <li>
              To take a suggestion, press the green checkmark button on the left side. By taking the
              suggestion, the code in your draft will be updated to reflect that – you can
              immediately see the update in your IDE as well.
            </li>
            <li>To dismiss a suggestion, press the red cross button on the left.</li>
          </ol>
          <li>
            All suggestions can be dismissed at once by pressing the cross button in the right of
            the banner
          </li>
        </ol>

        <h2 id="giving-feedback">Giving feedback or help to others</h2>
        <p>
          You can give help and feedback to others spontaneously or when somebody asks for feedback.
          In the left sidebar of Sturdy you can see your own drafts as well as those of others on
          the team.
        </p>
        <ol>
          <li>
            Navigate to a draft by somebody else on the team by clicking on it in the sidebar. The
            avatars to the left of draft names indicate who the author is.
          </li>
          <li>
            If the author of the draft is currently working on it, you will see their avatar in the
            top right corner.
          </li>
          <li>
            If the author of the draft is making code changes, you will see them in real-time, the
            same way you would see changes in your own draft.
          </li>
          <li>
            <SubtleText
              ><em>Optional:</em> Comment on code [<a href="#comment-on-code">go to</a>]
            </SubtleText>
          </li>
          <li>
            <SubtleText
              ><em>Optional:</em> Quickly indicating feedback [<a href="#quick-feedback">go to</a>]
            </SubtleText>
          </li>
          <li>
            <SubtleText
              ><em>Optional:</em> Trying the code [<a href="#trying-code">go to</a>]
            </SubtleText>
          </li>
        </ol>

        <h3 id="quick-feedback" title="Quickly indicating feedback">
          Optional: Quickly indicating feedback
        </h3>
        <p>You can give quick feedback without spelling it out in text.</p>
        <ol>
          <li>
            To quickly indicate that code looks good, use the button with the thumbs up icon in the
            right sidebar.
          </li>
          <li>
            To quickly indicate that you have some important feedback, use the button with the (i)
            icon in the right sidebar. You would typically do that to indicate if you have spotted
            an error but don't have time right now to elaborate.
          </li>
        </ol>

        <h3 id="trying-code" title="Trying the code">Optional: Trying the code</h3>
        <ol>
          <li>
            While viewing a draft of another team member, press the "Connect
            <em>folder</em> for suggesting" button.
          </li>
          <ol>
            <li>
              You will see a new nested draft change named "Suggestion 1" below the draft that you
              are trying.
            </li>
            <li>
              This temporary draft has the code of your team member's draft as it looked at the time
              you pressed the connect button.
            </li>
          </ol>
          <li>The folder on your computer now has the code of your team member's draft.</li>
          <li>
            Explore or run the code as you normally would using your IDE / text editor / terminal
          </li>
          <li>
            <SubtleText
              ><em>Optional:</em> Give code suggestions [<a href="#give-suggestions">go to</a>]
            </SubtleText>
          </li>
          <li>
            When you are done you can navigate back to your own draft change and press the "Connect
            <em>folder</em> for editing" button to continue where you left off.
          </li>
        </ol>

        <p>
          If the original author makes changes to their draft change after you have pressed the
          "Connect <em>folder</em> for suggesting" you will not see those new changes. This is a
          known limitation. To work around this:
        </p>

        <ol>
          <li>Navigate to the draft of your teammate</li>
          <li>Press the "Connect <em>folder</em> for suggesting" button</li>
          <ol>
            <li>
              Observe that now there is a second temporary draft named "Suggestion 2" to which you
              are connected
            </li>
          </ol>
        </ol>

        <h3 id="give-suggestions">Give code suggestions</h3>

        <p>
          When you have the code from someone else's draft change on your computer, you can give
          suggestions by simply making changes in your IDE.
        </p>
        <p>While you are connected to the "Suggestion <em>n</em>" draft change:</p>

        <ol>
          <li>Type any changes as you normally would, in your IDE or text editor.</li>
          <li>
            Within the temporary suggestion draft change you can see the changes you have created
          </li>
          <li>
            The author of the original draft change is notified of the suggestions [<a
              href="#taking-dismissing-suggestions"
              >go to</a
            >]
          </li>
        </ol>

        <p>
          You can suggest changes even to files that have not been modified in the original draft
          change.
        </p>
      </div>
    </template>
  </DocumentationWithTableOfContents>
</template>

<script lang="ts" setup>
import DocumentationWithTableOfContents from '../../layouts/DocumentationWithTableOfContents.vue'
import { useHead } from '@vueuse/head'
import DocsSidebar from '../../organisms/docs/DocsSidebar.vue'
import SubtleText from '../../molecules/docs/SubtleText.vue'
import DocsInfoBox from '../../molecules/DocsInfoBox.vue'
import { BellIcon, PlusIcon } from '@heroicons/vue/solid'

useHead({
  meta: [
    {
      name: 'description',
      content: 'Learn how to collaborate in real-time with Sturdy',
    },
    {
      name: 'keywords',
      content: 'study learn documentation collaboration team-work team work others',
    },
  ],
  title: 'How to collaborate with others | Sturdy',
})
</script>
