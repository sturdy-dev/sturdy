<template>
  <PaddedApp class="bg-white">
    <OrganizationCreate />
  </PaddedApp>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import PaddedApp from '../../layouts/PaddedApp.vue'
import OrganizationCreate from '../../organisms/organization/OrganizationCreate.vue'

export default defineComponent({
  components: { PaddedApp, OrganizationCreate },
})
</script>
