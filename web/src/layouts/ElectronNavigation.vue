<template>
  <main class="flex flex-1 flex-col">
    <AppTitleBarSpacer
      fixed
      pad-right="1rem"
      class="bg-gray-50 border-b left-0 right-0"
      :class="{
        'border-black': appPlatform === 'win32',
      }"
    >
      <div class="flex-1 flex flex-row justify-between gap-8 items-center">
        <AppTitleBarSpacer v-slot="{ ipc }" pad-left="1rem">
          <AppHistoryNavigationButtons :ipc="ipc" class="flex" />
        </AppTitleBarSpacer>

        <AppMutagenStatus class="flex-1" />
        <AppShareButton class="flex-shrink-0" />
      </div>
    </AppTitleBarSpacer>

    <!-- Primary column -->
    <section
      class="flex-1 flex flex-col overflow-x-auto"
      :class="{ 'spacer-padding': !!appPlatform }"
    >
      <slot></slot>
    </section>
  </main>
</template>

<script lang="ts" setup>
import AppTitleBarSpacer from '../components/AppTitleBarSpacer.vue'
import AppHistoryNavigationButtons from '../components/AppHistoryNavigationButtons.vue'
import AppMutagenStatus from '../components/AppMutagenStatus.vue'
import AppShareButton from '../components/AppShareButton.vue'

const appPlatform = appEnvironment?.platform
</script>
