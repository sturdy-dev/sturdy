<template>
  <Spacer>
    <div class="flex flex-1 items-center justify-center">
      <h1>Preferences</h1>
    </div>
  </Spacer>
  <div class="flex px-2 space-x-4">
    <Navigation />
    <section id="servers" class="flex flex-1">
      <ServersList />
    </section>
  </div>
</template>

<script lang="ts">
import Spacer from './components/Spacer.vue'
import Navigation from './components/Navigation.vue'
import ServersList from './components/ServersList.vue'

export default {
  components: {
    Spacer,
    Navigation,
    ServersList,
  },
}
</script>
