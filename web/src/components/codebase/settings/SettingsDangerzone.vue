<template>
  <HorizontalDivider class="mt-4" color="red" bg="bg-white">Dangerzone</HorizontalDivider>

  <ConfirmModal
    :show="showArchiveModal"
    title="Archive Codebase"
    subtitle="Are you sure that you want to archive this codebase?"
    @confirmed="archive"
    @close="showArchiveModal = false"
  />

  <div class="mx-4 divide-y divide-gray-200 px-4 sm:px-6 gap-2 flex">
    <Button v-if="codebase.isPublic" color="red" :icon="lockClosedIcon" @click="doMakePrivate">
      Make this codebase <strong>private</strong>
    </Button>
    <Button v-else color="red" :icon="lockOpenIcon" @click="doMakePublic">
      Make this codebase <strong>public</strong>
    </Button>

    <Button color="red" :icon="archiveIcon" @click="showArchiveModal = true">
      Archive this codebase
    </Button>
  </div>
</template>

<script lang="ts">
import HorizontalDivider from '../../shared/HorizontalDivider.vue'
import { ArchiveIcon, LockOpenIcon, LockClosedIcon } from '@heroicons/vue/solid'
import { gql, useMutation } from '@urql/vue'
import { defineComponent } from 'vue'
import type { PropType } from 'vue'
import Button from '../../shared/Button.vue'
import type { SettingsDangerzoneFragment } from './__generated__/SettingsDangerzone'
import ConfirmModal from '../../../molecules/ConfirmModal.vue'

export const SETTINGS_DANGERZONE = gql`
  fragment SettingsDangerzone on Codebase {
    id
    name
    isPublic
  }
`

export default defineComponent({
  components: {
    ConfirmModal,
    HorizontalDivider,
    Button,
  },
  props: {
    codebase: {
      type: Object as PropType<SettingsDangerzoneFragment>,
    },
  },
  setup() {
    const { executeMutation: archiveCodebaseResult } = useMutation(gql`
      mutation SettingsDangerzoneArchive($id: ID!) {
        updateCodebase(input: { id: $id, archive: true }) {
          id
          name
          archivedAt
        }
      }
    `)

    const { executeMutation: makeCodebasePublicResult } = useMutation(gql`
      mutation SettingsMakePublic($id: ID!) {
        updateCodebase(input: { id: $id, isPublic: true }) {
          id
          isPublic
        }
      }
    `)

    const { executeMutation: makeCodebasePrivateResult } = useMutation(gql`
      mutation SettingsMakePrivate($id: ID!) {
        updateCodebase(input: { id: $id, isPublic: false }) {
          id
          isPublic
        }
      }
    `)

    return {
      async archiveCodebase(id) {
        const variables = { id }
        await archiveCodebaseResult(variables).then((result) => {
          if (result.error) {
            throw new Error(result.error)
          }
        })
      },

      async makePublic(id) {
        const variables = { id }
        await makeCodebasePublicResult(variables).then((result) => {
          if (result.error) {
            throw new Error(result.error)
          }
        })
      },

      async makePrivate(id) {
        const variables = { id }
        await makeCodebasePrivateResult(variables).then((result) => {
          if (result.error) {
            throw new Error(result.error)
          }
        })
      },

      lockClosedIcon: LockClosedIcon,
      lockOpenIcon: LockOpenIcon,
      archiveIcon: ArchiveIcon,
    }
  },
  data() {
    return {
      updateStatus: '',
      showArchiveModal: false,
    }
  },
  methods: {
    archive() {
      this.archiveCodebase(this.codebase.id)
        .then(() => {
          this.emitter.emit('notification', {
            title: 'Archived codebase',
            message: this.codebase.name + ' has been archived',
          })
          this.$router.push({ name: 'home' })
        })
        .catch(() => {
          this.updateStatus = 'Something went wrong.'
        })
    },
    doMakePublic() {
      this.makePublic(this.codebase.id)
        .then(() => {
          this.emitter.emit('notification', {
            title: 'Public!',
            message: this.codebase.name + ' is now a public codebase',
          })
        })
        .catch(() => {
          this.updateStatus = 'Something went wrong.'
        })
    },
    doMakePrivate() {
      this.makePrivate(this.codebase.id)
        .then(() => {
          this.emitter.emit('notification', {
            title: 'Private!',
            message: this.codebase.name + ' is now a private codebase',
          })
        })
        .catch(() => {
          this.updateStatus = 'Something went wrong.'
        })
    },
  },
})
</script>
