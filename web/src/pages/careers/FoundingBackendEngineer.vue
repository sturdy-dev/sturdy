<template>
  <Public>
    <JobDescription title="Founding Backend Engineer" />
  </Public>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import Public from '../../layouts/Public.vue'
import JobDescription from '../../organisms/careers/JobDescription.vue'

export default defineComponent({
  components: {
    Public,
    JobDescription,
  },
})
</script>
