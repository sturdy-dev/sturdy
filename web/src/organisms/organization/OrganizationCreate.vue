<template>
  <form class="space-y-4">
    <Header>
      <span>Create a new team 🎉</span>
    </Header>

    <form class="space-y-4" @submit.stop.prevent="create">
      <TextInputWithLabel
        v-model="organizationName"
        placeholder="What's the name of your team or project?"
        label="Organization name"
        name="org-name"
      />
      <OrganizationLicenseTierPicker v-if="withTierPicker" />
      <Button color="green" @click="create">Get started</Button>
    </form>
  </form>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import Header from '../../molecules/Header.vue'
import OrganizationLicenseTierPicker from '../../organisms/organization/OrganizationLicenseTierPicker.vue'
import { gql, useMutation } from '@urql/vue'
import Button from '../../components/shared/Button.vue'
import { useRouter } from 'vue-router'
import {
  CreateOrganizationMutation,
  CreateOrganizationMutationVariables,
} from './__generated__/OrganizationCreate'
import TextInputWithLabel from '../../molecules/TextInputWithLabel.vue'
export default defineComponent({
  components: { TextInputWithLabel, Header, OrganizationLicenseTierPicker, Button },
  props: {
    withTierPicker: {
      type: Boolean,
      required: true,
    },
  },
  setup() {
    let { executeMutation: executeCreateOrganization } = useMutation<
      CreateOrganizationMutation,
      CreateOrganizationMutationVariables
    >(gql`
      mutation createOrganization($name: String!) {
        createOrganization(input: { name: $name }) {
          id
          name
        }
      }
    `)
    let router = useRouter()
    return {
      async createMutation(name: string) {
        const variables = { name }
        await executeCreateOrganization(variables).then((result) => {
          if (result.error) {
            throw new Error(result.error)
          }
          router.push({
            name: 'codebaseOverview',
            params: { id: result.data?.createOrganization.id },
          })
        })
      },
    }
  },
  data() {
    return {
      organizationName: '',
      organizationLegalName: '',
    }
  },
  methods: {
    create() {
      this.createMutation(this.organizationName)
    },
  },
})
</script>
