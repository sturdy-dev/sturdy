<template>
  <PaddedApp class="bg-white">
    <main v-if="data" class="flex-1 relative focus:outline-none">
      <div class="block space-y-4 xl:border-gray-200">
        <div>
          <div class="flex justify-center flex-col space-y-2">
            <div class="flex justify-end">
              <div class="mt-4 flex space-x-3 md:mt-0"></div>
            </div>

            <div v-if="fetching || (stale && data.codebase.file == null)">
              <Spinner />
            </div>

            <FileOrDirectory
              v-else-if="data.codebase?.file"
              :file-or-directory="data.codebase.file"
              :codebase="data.codebase"
            />

            <div v-else>
              <h1 class="text-3xl mb-3">
                <span class="text-gray-400">404</span>
                Not Found
              </h1>

              <p class="max-w-prose">
                This file could've been removed, or it was never here in the first place...
              </p>

              <router-link
                :to="{ name: 'codebaseHome' }"
                class="px-3 py-1 bg-blue-600 text-white rounded-md inline-block mt-3"
              >
                Go back to {{ data.codebase.name }}
              </router-link>
            </div>
          </div>
        </div>
      </div>
    </main>
  </PaddedApp>
</template>

<script lang="ts">
import { gql, useQuery } from '@urql/vue'
import { useRoute } from 'vue-router'
import FileOrDirectory, { OPEN_FILE_OR_DIRECTORY } from '../components/browse/FileOrDirectory.vue'
import { ref, watch } from 'vue'
import Spinner from '../components/shared/Spinner.vue'
import { DIRECTORY_BREADCRUMB } from '../components/browse/DirectoryBreadcrumb.vue'
import PaddedApp from '../layouts/PaddedApp.vue'

function packPath(path: string | string[]): string {
  if (Array.isArray(path)) {
    return path.join('/')
  }
  return path
}

export default {
  components: { PaddedApp, Spinner, FileOrDirectory },
  setup() {
    const route = useRoute()
    const path = ref(packPath(route.params.path))
    watch(route, () => {
      path.value = packPath(route.params.path)
    })

    const { data, fetching, stale } = useQuery({
      query: gql`
        query BrowseFile($id: ID!, $path: String!) {
          codebase(shortID: $id) {
            id
            name
            file(path: $path) {
              ...OpenFileOrDirectory
            }
            ...DirectoryBreadcrumb
          }
        }
        ${OPEN_FILE_OR_DIRECTORY}
        ${DIRECTORY_BREADCRUMB}
      `,
      variables: {
        id: route.params.codebaseSlug as string,
        path: path,
      },
      requestPolicy: 'cache-and-network',
    })
    return {
      data,
      fetching,
      stale,
    }
  },
}
</script>
