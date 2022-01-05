<template>
  <StaticPage
    title="Sturdy for Git users"
    subtitle="Learn Sturdy"
    category="solutions"
    metadescription="Learn Sturdy by comparing it to Git"
  >
    <div
      class="mt-5 prose prose-yellow text-gray-500 mx-auto lg:max-w-none lg:row-start-1 lg:col-start-1"
    >
      <p>
        While Sturdy and Git are fundamentally different, there are still some similarities between
        the two. On this page we'll go trough the major features of both Sturdy and Git (with GitHub
        or GitLab) and see how they compare.
      </p>

      <table>
        <thead>
          <tr>
            <th>Sturdy</th>
            <th>Git</th>
          </tr>
        </thead>
        <tbody>
          <tr>
            <td>Codebase</td>
            <td>Repository</td>
          </tr>
          <tr>
            <td>
              <router-link :to="{ name: 'featuresWorkspaceNavigation' }">Workspace</router-link>
            </td>
            <td>Branch, Pull Request</td>
          </tr>
          <tr>
            <td><router-link :to="{ name: 'syncing' }">Share</router-link></td>
            <td>Merge a Pull Request</td>
          </tr>
          <tr>
            <td>Change</td>
            <td>Commit</td>
          </tr>
          <tr>
            <td><router-link :to="{ name: 'syncing' }">Sync</router-link></td>
            <td>Rebase</td>
          </tr>
          <tr>
            <td class="text-gray-400"><em>N/A</em></td>
            <td>Push</td>
          </tr>
          <tr>
            <td class="text-gray-400"><em>N/A</em></td>
            <td>Pull</td>
          </tr>
          <tr>
            <td class="text-gray-400"><em>N/A</em></td>
            <td>Command Line</td>
          </tr>
        </tbody>
      </table>

      <p>
        <router-link :to="{ name: 'resourcesMigrateFromGitHub' }">
          Read More: Migrate from GitHub to Sturdy
        </router-link>
      </p>
    </div>
  </StaticPage>

  <DiveIn />
</template>

<script>
import StaticPage from '../../StaticPage.vue'
import DiveIn from '../../../../components/index/DiveIn.vue'
export default {
  name: 'WorkspaceNavigation',
  components: { StaticPage, DiveIn },
}
</script>

<style scoped></style>
