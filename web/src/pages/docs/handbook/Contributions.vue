<template>
  <DocumentationWithTableOfContents>
    <template #sidebar>
      <HandbookSidebar />
    </template>

    <template #default>
      <div class="prose">
        <h1>External Contributions</h1>

        <p>
          Sturdy accepts (and welcomes!) contributions to our code from external parties. We're
          currently requiring all contributors to sign a Contributor License Agreement (CLA) before
          they can contribute.
        </p>

        <ol>
          <li>
            The contributor (or their employer) needs to sign the
            <a href="https://forms.gle/u8tfb8LyNtLNitQ27">Contributor License Agreement (CLA)</a>
          </li>
          <li>
            <em>Sturdy Employee:</em>
            <a
              href="https://docs.google.com/spreadsheets/d/1YyYFxEYo8dovmvrVMUHx_gC0EyLga9UztrteVzRmwrY/edit?usp=sharing"
              >Verify</a
            >
            that the contributor has filled in the form
          </li>
          <li>
            <em>Sturdy Employee:</em> Add the contributors GitHub username to
            <a href="https://getsturdy.com/clabot-config-AySxptY">the .clabot configuration</a>
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
