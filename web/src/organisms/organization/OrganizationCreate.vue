<template>
  <form class="space-y-4">
    <template v-if="isFirst && isMultiTennant">
      <h1 class="text-gray-800 text-4xl font-bold">Let's get you setup 🎉</h1>
      <p class="text-gray-500">Create your first organization</p>
    </template>
    <template v-else-if="isMultiTennant">
      <h1 class="text-gray-800 text-4xl font-bold">Create a new organization 🎉</h1>
      <p class="text-gray-500">Create a organization to organize your work</p>
    </template>
    <template v-else-if="!isMultiTennant">
      <h1 class="text-gray-800 text-4xl font-bold">Let's get you setup 🎉</h1>
      <p class="text-gray-500">Create an organization for your Sturdy server</p>
    </template>

    <form class="space-y-4" @submit.stop.prevent="create">
      <TextInputWithLabel
        v-model="organizationName"
        placeholder="What's the name of your team or project?"
        label="Team name"
        name="org-name"
      />
      <OrganizationLicenseTierPicker v-if="withTierPicker" />
      <Button color="green" @click="create">Get started</Button>
    </form>

    <template v-if="isMultiTennant">
      <p class="text-gray-700 text-sm">
        Create a organization to manage your projects, codebases, members, and billing.
      </p>
      <p class="text-gray-700 text-sm">
        If you're creating a organization for work, use the company name as the name.
      </p>
    </template>
    <template v-else>
      <p class="text-gray-700 text-sm">Create an organization to setup your Sturdy server.</p>
      <p class="text-gray-700 text-sm">
        If you're creating a organization for work, use the company name as the name of the
        organization.
      </p>
    </template>
  </form>
</template>

<script lang="ts">
import { defineComponent, PropType } from 'vue'
import OrganizationLicenseTierPicker from '../../organisms/organization/OrganizationLicenseTierPicker.vue'
import { gql } from '@urql/vue'
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
  components: { TextInputWithLabel, OrganizationLicenseTierPicker, Button },
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
    isMultiTennant: {
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
