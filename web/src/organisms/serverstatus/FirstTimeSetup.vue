<template>
  <div>
    <h1>Hey! Welcome to your Sturdy instance!</h1>
    <p>Let's create your server...</p>
    <OrganizationCreate />
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import OrganizationCreate from '../organization/OrganizationCreate.vue'

export default defineComponent({
  components: { OrganizationCreate },
})
</script>
