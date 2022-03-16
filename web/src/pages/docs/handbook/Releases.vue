<template>
  <DocumentationWithTableOfContents>
    <template #sidebar>
      <HandbookSidebar />
    </template>

    <template #default>
      <div class="prose">
        <h1>Making Releases</h1>

        <p>This page documents how to make a new release of Sturdy.</p>

        <h2>The Docker "oneliner"</h2>

        <p>How to make a new release of the Sturdy Docker oneliner.</p>

        <ol>
          <li>In a new draft change, prepare release notes (CHANGELOG.md)</li>
          <li>
            Announce in <a href="https://getsturdy.slack.com/archives/C01MHP6SUKS">#general</a> that
            you're making a release
          </li>
          <li>
            Run <code>./scripts/build-upload-oneliner.sh --patch --push</code> to make and push the
            release to Docker Hub. Use --patch, --minor, or --major to bump the patch, minor, or
            major version number accordingly. Remove the --push argument to build the release
            without pushing it.
          </li>
          <li>Share/merge your updates to CHANGELOG.md</li>
          <li>
            <a href="https://github.com/sturdy-dev/sturdy/releases/new">
              Tag the release on GitHub
            </a>
            using the format <code>api/$VERSION</code>
          </li>
        </ol>

        <h2>App</h2>

        <p>How to make a new release of the Sturdy Electron app.</p>

        <ol>
          <li>In a new draft change, prepare release notes (CHANGELOG.md)</li>
          <li>
            Announce in <a href="https://getsturdy.slack.com/archives/C01MHP6SUKS">#general</a> that
            you're making a release
          </li>
          <li>Bump the version number in <code>app/package.json</code></li>
          <li>
            Run <code>cd app && ./build-electron-builder.sh --upload --stable</code> to make and
            push the release. Remove the --stable argument to build a new release to the
            Beta-channel. Remove the --upload argument to dry-run the build process.
          </li>
          <li>Share/merge your updates to CHANGELOG.md</li>
          <li>
            <a href="https://github.com/sturdy-dev/sturdy/releases/new">
              Tag the release on GitHub
            </a>
            using the format <code>app/$VERSION</code>
          </li>
        </ol>
      </div>
    </template>
  </DocumentationWithTableOfContents>
</template>

<script lang="ts" setup>
import DocumentationWithTableOfContents from '../../../layouts/DocumentationWithTableOfContents.vue'
import HandbookSidebar from '../../../organisms/docs/HandbookSidebar.vue'
</script>
