<template>
  <form class="space-y-4">
    <Header>
      <span v-if="isFirst">Create your first organization 🎉</span>
      <span v-else>Create a new organization 🎉</span>
    </Header>

    <form class="space-y-4" @submit.stop.prevent="create">
      <TextInputWithLabel
        v-model="organizationName"
        placeholder="What's the name of your team or project?"
        label="Team name"
        name="org-name"
      />
      <OrganizationLicenseTierPicker v-if="withTierPicker" />
      <Button color="green" @click="create">Create team</Button>
    </form>

    <p class="text-gray-700 text-sm">
      Create a organization to manage your projects, codebases, members, and billing.
    </p>
    <p class="text-gray-700 text-sm">
      If you're creating a organization for work, use the company name as the name.
    </p>
  </form>
</template>

<script lang="ts">
import { defineComponent, PropType } from 'vue'
import Header from '../../molecules/Header.vue'
import OrganizationLicenseTierPicker from '../../organisms/organization/OrganizationLicenseTierPicker.vue'
import { gql, useMutation } from '@urql/vue'
import Button from '../../components/shared/Button.vue'
import { useRouter } from 'vue-router'
import { OrganizationCreateUserFragment } from './__generated__/OrganizationCreate'
import TextInputWithLabel from '../../molecules/TextInputWithLabel.vue'
import { useCreateOrganization } from '../../mutations/useCreateOrganization'

export const ORGANIZATION_CREATE_USER = gql`
  fragment OrganizationCreateUser on User {
    id
    name
  }
`

export default defineComponent({
  components: { TextInputWithLabel, Header, OrganizationLicenseTierPicker, Button },
  props: {
    withTierPicker: {
      type: Boolean,
      required: true,
    },
    user: {
      type: Object as PropType<OrganizationCreateUserFragment>,
      required: true,
    },
    isFirst: {
      type: Boolean,
      required: true,
    },
  },
  setup() {
    let executeCreateOrganization = useCreateOrganization()

    let router = useRouter()
    return {
      async createMutation(name: string) {
        const variables = { name }
        await executeCreateOrganization(variables).then((result) => {
          router.push({
            name: 'organizationListCodebases',
            params: { organizationSlug: result.createOrganization.shortID },
          })
        })
      },
    }
  },
  data() {
    return {
      organizationName: this.proposedTeamName(),
    }
  },
  computed: {},
  methods: {
    create() {
      this.createMutation(this.organizationName)
    },
    proposedTeamName(): string {
      let name = this.user.name.split(' ')
      if (name.length === 0) {
        return 'My startup'
      }

      let fname = name[0]
      let apos = "'s"
      if (fname.endsWith('s')) {
        apos = "'"
      }

      return `${fname}${apos} startup`
    },
  },
})
</script>
