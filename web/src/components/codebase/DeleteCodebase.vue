<template>
  <div>
    <div :class="modal">
      <div class="modal-background" />
      <div class="modal-card">
        <div class="modal-card-head">
          <p class="modal-card-title">Are you sure you want to delete this codebase?</p>
          <button class="delete" aria-label="close" @click="close" />
        </div>
        <section class="modal-card-body">
          You are about to delete <b>{{ codebaseName }}</b
          >. This action can not be undone.
        </section>
        <footer class="modal-card-foot">
          <button :class="['button', 'is-link']" @click="deleteCodebase">Delete codebase</button>
          <button class="button" @click="close">Cancel</button>
        </footer>
      </div>
    </div>
  </div>
</template>

<script>
import http from '../../http'

export default {
  name: 'DeleteCodebase',
  props: ['isActive', 'codebaseId', 'codebaseName'],
  emits: ['closeDeleteCodebase', 'deletedCodebase'],
  computed: {
    modal() {
      return {
        'is-active': this.isActive,
        modal: true,
      }
    },
  },
  methods: {
    deleteCodebase() {
      fetch(http.url('v3/codebases/' + this.codebaseId), {
        method: 'DELETE',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include',
      }).then(() => {
        this.$emit('deletedCodebase')
        this.$emit('closeDeleteCodebase')
      })
    },
    close() {
      this.$emit('closeDeleteCodebase')
    },
  },
}
</script>
