<template>
  <ArchiveWorkspaceModal
    :is-active="modalOpened"
    :workspace-id="workspaceId"
    @archived="onWorkspaceArchived"
    @close="hideModal"
    @archiving="setArchiving"
  />
  <Button size="wider" @click="showModal">
    <Spinner v-if="archiving" class="-ml-1 mr-2 h-5 w-5 text-gray-400" aria-hidden="true" />
    <ArchiveIcon v-else class="-ml-1 mr-2 h-5 w-5 text-gray-400" aria-hidden="true" />
    <span>Archive</span>
  </Button>
</template>

<script lang="ts">
import Button from '../../components/shared/Button.vue'
import Spinner from '../../components/shared/Spinner.vue'
import { ArchiveIcon } from '@heroicons/vue/solid'
import ArchiveWorkspaceModal from '../../components/codebase/ArchiveWorkspaceModal.vue'

export default {
  components: {
    ArchiveIcon,
    ArchiveWorkspaceModal,
    Button,
    Spinner,
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
