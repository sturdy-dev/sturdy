<template>
  <PaddedApp>
    <CreateCodebase>
      <template #header>
        <div class="py-8 px-4">
          <div class="">
            <h2
              class="text-4xl font-extrabold text-gray-900 sm:text-4xl sm:tracking-tight lg:text-4xl"
            >
              Create a new codebase
            </h2>
            <p class="mt-5 text-xl text-gray-500">
              Create your first Sturdy codebase. You'll soon be ready to code!
            </p>
          </div>
        </div>
      </template>
    </CreateCodebase>
  </PaddedApp>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import PaddedApp from '../layouts/PaddedApp.vue'
import CreateCodebase from '../organisms/CreateCodebase.vue'

export default defineComponent({
  components: { PaddedApp, CreateCodebase },
})
</script>
