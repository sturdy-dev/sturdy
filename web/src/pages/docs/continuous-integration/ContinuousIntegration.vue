<template>
  <StaticPage
    title="Continuous Integration/Delivery (CI/CD)"
    metadescription="Learn how Sturdy integrates with your CI/CD workflows"
    image="https://images.unsplash.com/photo-1629904853716-f0bc54eea481?ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&ixlib=rb-1.2.1&auto=format&fit=crop&w=3540&q=80"
  >
    <div class="prose prose-yellow">
      <p>
        Sturdy has first-party integrations with popular CI/CD tools like Buildkite, and also
        supports bringing in existing CI/CD integrations via the
        <router-link :to="{ name: 'resourcesMigrateFromGitHub' }"
          >Sturdy for GitHub bridge</router-link
        >.
      </p>

      <h2>Supported CI/CD providers</h2>

      <ul>
        <li>Buildkite</li>
        <li>GitHub Actions</li>
        <li>Circle CI<sup>*</sup></li>
        <li>Jenkins<sup>*</sup></li>
        <li>TeamCity<sup>*</sup></li>
        <li>Bamboo<sup>*</sup></li>
        <li>Travis CI<sup>*</sup></li>
        <li>... and many more! <sup>*</sup></li>
      </ul>

      <p>
        <em><sup>*</sup> via Sturdy for GitHub</em>
      </p>

      <h2>Pre-merge testing</h2>

      <p>
        Pre-merge testing are the tests that you execute <em>before</em> merging a Draft on Sturdy.
        Usually this is a list of faster tests and linters that you never want to break.
      </p>

      <p>
        To run the tests, open up the Draft in Sturdy, and click on <em>Trigger CI</em> in the
        sidebar. This will schedule a build on all installed CI/CD providers, and their status will
        be reported back to Sturdy.
      </p>

      <div class="flex flex-col gap-4 md:gap-8">
        <div class="flex md:flex-row flex-col align-bottom gap-4 space-between w-full">
          <CardWithFooter>
            <template #default>
              <img
                src="./trigger-ci.png"
                alt="Trigger CI/CD on a Sturdy Draft"
                height="612"
                width="728"
              />
            </template>
            <template #footer>
              <em>Click <strong>Trigger CI</strong> to run tests</em>
            </template>
          </CardWithFooter>

          <CardWithFooter>
            <template #default>
              <img height="606" width="728" src="./tests-pending.png" alt="Tests are pending" />
            </template>
            <template #footer>
              <em>The tests are now running</em>
            </template>
          </CardWithFooter>
        </div>

        <CardWithFooter>
          <template #default>
            <img
              src="./detailed-tests-status.png"
              alt="Detailed test status"
              height="1314"
              width="1284"
            />
          </template>
          <template #footer>
            <em>Detailed test results in Sturdy</em>
          </template>
        </CardWithFooter>
      </div>

      <p>
        When a codebase is connected to GitHub, Sturdy will push a snapshot of the draft change as a
        commit on a branch called <code>sturdy-ci-${ID}</code>. If the tests does not start running,
        make sure to configure GitHub Actions (or other CI/CD integration) to trigger tests on
        pushes to branches with this name pattern.
      </p>

      <p>Any test results and statuses set via GitHub will automatically be forwarded to Sturdy.</p>

      <h2>Post-merge testing</h2>

      <p>
        Post-merge is the testing (and other CI/CD jobs) that happen after a change has been merged.
        This is where you might run tests that are slower and more expensive to run, and automated
        releases to production environments.
      </p>

      <p>
        Sturdy will trigger tests on <em>trunk</em> the same way that it triggers the pre-merge
        tests. All native configurations will be triggered automatically. When connecting via
        GitHub, the push to <em>main</em> (the default branch) will trigger the tests on connected
        providers on GitHub, and any test results will be forwarded to Sturdy.
      </p>

      <CardWithFooter>
        <template #default>
          <img
            alt="Sturdy adds status icons to all the latest changes displayed on the top of the codebase page"
            src="./status-on-codebase.png"
            width="1622"
            height="652"
          />
        </template>
        <template #footer>
          <em
            >On the top of the codebase page, the list of latest changes will be complemented with a
            status icon indicating the result of any CI/CD workflow.</em
          >
        </template>
      </CardWithFooter>

      <p>
        On the codebase overview page the test results for the latest 4 changes will be visible. On
        the changelog the status for each change will be visible.
      </p>

      <h2>Faster CI/CD with native integrations</h2>
      <p>
        To speed up CI/CD for large codebases, Sturdy exports snapshots of the codebase that are
        fast to download (as compared to doing a full git-clone, which for a large repository can be
        slower than the tests that you want to run). This is built on Sturdys
        <a href="https://schema.getsturdy.com/#definition-ContentsDownloadURL"
          >ContentsDownloadURL</a
        >
        API (for Changes and Drafts/Workspaces). And will soon support fast incremental downloads
        for larger files and repositories out-of-the-box.
      </p>

      <p>
        For compatability with providers that expects Git, Sturdy exposes a "fake" repository with a
        script to download the full snapshot. The repository only contains two files, and all you
        need to do is to run <code>./download</code> and Sturdy will take care of authenticating and
        downloading all files to a directory called <code>tmp_output</code>.
      </p>

      <ul>
        <li><code>download</code> – Script with built-in authentication</li>
        <li><code>sturdy.json</code> – Metadata file about what to download</li>
      </ul>

      <p>
        Example contents of the sturdy.json file. The contents will vary slightly depending on what
        should be downloaded (and is not guaranteed to be stable).
      </p>

      <pre>
{
  "codebase_id": "303cb44e-7127-42ce-a5f9-97fcc023c8ef",
  "workspace_id": "a76b7fae-7e9d-43c3-8b13-6c19a63d0241",
  "snapshot_id": "d8482aa2-2cb2-4192-ba6c-96c758637026"
}</pre
      >
    </div>
  </StaticPage>
</template>

<script lang="ts" setup>
import StaticPage from '../../../layouts/StaticPage.vue'
import CardWithFooter from '../../../atoms/CardWithFooter.vue'
</script>
