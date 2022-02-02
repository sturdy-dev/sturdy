<template>
  <div>
    <video
      id="video"
      autoplay
      loop
      muted
      playsinline
      alt="Sturdy demo"
      class="border border-slate-700/70 rounded-lg object-cover object-center"
    >
      <source v-if="id === 'trying-code'" src="./try-code.mp4" type="video/mp4" />
      <source v-if="id === 'suggesting-code'" src="./give-suggestion.mp4" type="video/mp4" />
      <source v-if="id === 'taking-code'" src="./accept-suggestion.mp4" type="video/mp4" />
      <source v-if="id === 'creating-pr'" src="./create-pr.mp4" type="video/mp4" />
    </video>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'

export default defineComponent({
  props: {
    id: {
      type: String,
      required: true,
    },
  },
  watch: {
    id: function () {
      let video = document.getElementById('video')
      video?.load()
    },
  },
})
</script>

<style scoped></style>
