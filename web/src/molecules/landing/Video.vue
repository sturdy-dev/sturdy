<template>
  <div>
    <video
      id="video"
      ref="video"
      autoplay
      loop
      muted
      playsinline
      class="border border-slate-700/70 rounded-lg object-cover object-center"
    >
      <source :src="src" type="video/mp4" />
    </video>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'

export default defineComponent({
  props: {
    src: {
      type: String,
      required: true,
    },
  },
  watch: {
    src: function () {
      this.$refs.video.load()
    },
  },
})
</script>
