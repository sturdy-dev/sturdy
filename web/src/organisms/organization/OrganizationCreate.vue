<template>
  <form class="space-y-4">
    <Header>
      <span>Create a new organization</span>
    </Header>

    <TextInput v-model="organizationName" placeholder="Ex: Apple Inc." />

    <OrganizationLicenseTierPicker />

    <Button color="green" @click="create">Create</Button>
  </form>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import TextInput from '../../molecules/TextInput.vue'
import Header from '../../molecules/Header.vue'
import OrganizationLicenseTierPicker from '../../organisms/organization/OrganizationLicenseTierPicker.vue'
import { gql, useMutation } from '@urql/vue'
import Button from '../../components/shared/Button.vue'
import { useRouter } from 'vue-router'
import {
  CreateOrganizationMutation,
  CreateOrganizationMutationVariables,
} from './__generated__/OrganizationCreate'

export default defineComponent({
  components: { TextInput, Header, OrganizationLicenseTierPicker, Button },
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
            name: 'organizationView',
            params: { id: result.data?.createOrganization.id },
          })
        })
      },
    }
  },
  data() {
    return {
      organizationName: '',
    }
  },
  methods: {
    create() {
      this.createMutation(this.organizationName)
    },
  },
})
</script>
