<template>
  <Spacer />
  <div class="flex flex-row">
    <Navigation />
    <section id="servers">
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
