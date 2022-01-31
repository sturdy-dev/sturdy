<template>
  <div class="flex flex-col md:flex-row mt-4">
    <div v-if="editingName" class="flex flex-1 rounded-md shadow-sm mr-4" tabindex="0">
      <div class="relative flex items-stretch flex-grow focus-within:z-10">
        <input
          ref="codebase-name"
          v-model="newName"
          type="text"
          name="codebase_name"
          placeholder="Name the codebase ..."
          autocomplete="off"
          class="focus:ring-blue-500 focus:border-blue-500 block w-full rounded-none rounded-l-md sm:text-sm border-gray-300"
          @keydown="keyDown"
        />
      </div>
      <Button :is-grouped="true" :grouped="true" :last="true" class="-ml-1" @click="save">
        Save
      </Button>
    </div>

    <h1 v-else class="flex-1 text-lg font-medium">
      {{ codebaseName }}
    </h1>

    <div class="relative">
      <span class="relative z-0 inline-flex shadow-sm rounded-md">
        <Button @click="startStopEditName"> Edit Name </Button>
      </span>
    </div>
  </div>

  <Banner
    v-if="showRenameFailed"
    status="error"
    message="Could not update the name right now. Please try again!"
  />
</template>

<script>
import Button from '../../shared/Button.vue'
import { Banner } from '../../../atoms'
import { toRef, ref, watch } from 'vue'
import { gql, useMutation } from '@urql/vue'

export default {
  name: 'SettingsTitle',
  components: {
    Button,
    Banner,
  },
  props: {
    codebaseName: String,
    codebaseId: String,
  },
  setup(props) {
    let newName = ref('')
    let codebaseName = toRef(props, 'codebaseName')
    newName.value = codebaseName.value

    watch(codebaseName, (updatedName) => {
      if (updatedName) {
        newName.value = updatedName
      }
    })

    const { executeMutation: updateCodebaseResult } = useMutation(gql`
      mutation SettingsTitleUpdateName($id: ID!, $name: String) {
        updateCodebase(input: { id: $id, name: $name }) {
          id
          name
        }
      }
    `)

    return {
      newName,
      async updateCodebase(id, name) {
        const variables = { id, name }
        await updateCodebaseResult(variables).then((result) => {
          if (result.error) {
            throw new Error(result.error)
          }
        })
      },
    }
  },
  data() {
    return {
      editingName: false,
      showRenameFailed: false,
    }
  },
  methods: {
    save() {
      this.editingName = false
      this.showRenameFailed = false
      this.updateCodebase(this.codebaseId, this.newName)
        .then(() => {
          this.emitter.emit('notification', {
            title: 'Renamed codebase',
            message: 'Codebase has been renamed to ' + this.newName,
          })
        })
        .catch(() => {
          this.showRenameFailed = true
          setTimeout(() => {
            this.showRenameFailed = false
          }, 5000)
        })
    },

    startStopEditName() {
      if (this.editingName) {
        this.editingName = false
        return
      }

      this.editingName = true

      this.$nextTick(() => {
        this.$refs['codebase-name'].focus()
      })
    },

    keyDown(e) {
      // Enter
      if (e.keyCode === 13) {
        this.save()
      }
    },
  },
}
</script>
