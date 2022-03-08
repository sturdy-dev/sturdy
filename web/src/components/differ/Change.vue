<template>
  <div>
    <Banner
      v-if="hasHiddenChanges"
      class="my-2"
      status="warning"
      message="This change includes more files, but you don't have access to see them."
    />

    <Differ
      v-if="diffs"
      :diffs="diffs"
      :can-comment="authorized"
      :change-i-d="change.change_id || change.id"
      :comments="comments"
      :user="user"
      :members="members"
      :change="change"
      :show-full-file-button="showFullFileButton"
      :show-add-button="false"
      @submittedNewComment="$emit('submittedNewComment')"
    />
  </div>
</template>

<script>
import Differ from './Differ.vue'
import { Banner } from '../../atoms'

export default {
  name: 'Change',
  components: { Differ, Banner },
  props: ['change', 'comments', 'user', 'showFullFileButton', 'members'],
  emits: ['submittedNewComment'],
  data() {
    return {}
  },
  computed: {
    diffs() {
      return this.change.diffs.filter((d) => !d.isHidden)
    },
    hasHiddenChanges() {
      return this.change.diffs.length > this.diffs.length
    },
    authenticated() {
      return !!this.user
    },
    authorized() {
      const isMember = this.members.some(({ id }) => id == this.user?.id)
      return this.authenticated && isMember
    },
  },
}
</script>

<style scoped>
.change-description {
  width: 100%;
  word-break: break-word;
  white-space: pre-line;
}

.change-byline {
  display: inline-flex;
  align-items: center;
}

.change-byline .sturdy-card-avatar {
  margin: 0 4px;
  padding: 0;
  height: 24px;
  width: 24px;
}
</style>
