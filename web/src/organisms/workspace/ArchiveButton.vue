<template>
  <ArchiveWorkspaceModal
    :is-active="modalOpened"
    :workspace-id="workspaceId"
    @archived="onWorkspaceArchived"
    @close="hideModal"
    @archiving="setArchiving"
  />
  <Button :icon="archiveIcon" :spinner="archiving" @click="showModal">
    <span class="hidden sm:block">Archive</span>
  </Button>
</template>

<script lang="ts">
import Button from '../../components/shared/Button.vue'
import { ArchiveIcon } from '@heroicons/vue/solid'
import ArchiveWorkspaceModal from '../../components/codebase/ArchiveWorkspaceModal.vue'

export default {
  components: {
    ArchiveWorkspaceModal,
    Button,
  },
  props: {
    workspaceId: {
      type: String,
      required: true,
    },
  },
  data() {
    return {
      modalOpened: false,
      archiving: false,
      archiveIcon: ArchiveIcon,
    }
  },
  methods: {
    setArchiving() {
      this.archiving = true
    },
    showModal() {
      this.modalOpened = true
    },
    hideModal() {
      this.modalOpened = false
    },
    onWorkspaceArchived() {
      this.archiving = false
      this.$router.push({
        name: 'codebaseHome',
        params: { codebaseSlug: this.$route.params.codebaseSlug },
      })
    },
  },
}
</script>
