<template>
  <div class="space-y-2">
    <div class="space-y-1">
      <label class="block text-sm font-medium text-gray-700"> Organization name </label>
      <div class="flex justify-between gap-4">
        <text-input v-model="organizationName" />
        <Button :disabled="updateBtnEnable" @click="update"> Update </Button>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import { PropType } from 'vue/dist/vue'
import { OrganizationMembersOrganizationFragment } from './__generated__/OrganizationMembers'
import TextInput from '../../molecules/TextInput.vue'
import Button from '../../components/shared/Button.vue'
import { useUpdateOrganization } from '../../mutations/useUpdateOrganization'

export default defineComponent({
  name: 'OrganizationUpdate',
  components: { Button, TextInput },
  props: {
    organization: {
      type: Object as PropType<OrganizationMembersOrganizationFragment>,
      required: true,
    },
  },
  setup() {
    let updateOrganization = useUpdateOrganization()

    return {
      updateOrganization,
    }
  },
  data() {
    return {
      organizationName: this.organization.name,
    }
  },
  computed: {
    updateBtnEnable() {
      return this.organizationName.length === 0 || this.organizationName === this.organization.name
    },
  },
  methods: {
    async update() {
      const variables = { id: this.organization.id, name: this.organizationName }
      this.updateOrganization(variables)
        .then(() => {
          this.emitter.emit('notification', {
            title: 'Organization',
            message: 'Organization name updated!',
          })
        })
        .catch(() => {
          this.emitter.emit('notification', {
            title: 'Organization',
            message: 'Failed to update organization name! Please try later.',
            style: 'error',
          })
        })
    },
  },
})
</script>

<style scoped></style>
