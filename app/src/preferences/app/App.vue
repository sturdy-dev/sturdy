<template>
  <TitleBar>
    <template #header>
      <h1>Preferences</h1>
    </template>

    <template #default>
      <div class="flex p-2 space-x-4">
        <Navigation />
        <section id="servers" class="flex flex-1">
          <ServersList />
        </section>
      </div>
    </template>
  </TitleBar>
</template>

<script lang="ts">
import TitleBar from './components/TitleBar.vue'
import Navigation from './components/Navigation.vue'
import ServersList from './components/ServersList.vue'

export default {
  components: {
    TitleBar,
    Navigation,
    ServersList,
  },
}
</script>
