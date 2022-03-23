<template>
  <PaddedAppLeftSidebar v-if="data" class="bg-white">
    <template #navigation>
      <VerticalNavigation />
    </template>

    <template #header>
      <OrganizationSettingsHeader :name="data.organization.name" />
    </template>

    <template #default>
      <Header>
        <span>Installation</span>
      </Header>

      <div class="space-y-4">
        <p class="text-sm">Manage this Sturdy installation.</p>

        <template v-if="haveLicense">
          <div>
            <SubHeader class="mb-1">Current license</SubHeader>
            <OrganizationListLicenses :licenses="data.organization.licenses" />
          </div>
        </template>

        <div>
          <SubHeader v-if="haveLicense" class="mb-1"> Update the license key </SubHeader>
          <SubHeader v-else class="mb-1">Enter your license key</SubHeader>

          <Banner v-if="failedMessage" status="error" class="my-1">
            Unable to update the license key: {{ failedMessage }}
          </Banner>

          <div class="flex items-center space-x-2">
            <div class="flex-1">
              <TextInput
                v-model="licenseKeyInput"
                placeholder="Enter a new license key"
                label="License Key"
                @keydown.enter.stop="doUpdateLicenseKey"
              />
            </div>
            <Button :disabled="!licenseKeyInput" @click="doUpdateLicenseKey">Save</Button>
          </div>
        </div>
      </div>
    </template>
  </PaddedAppLeftSidebar>
</template>

<script lang="ts">
import { computed, defineComponent } from 'vue'
import Header from '../../molecules/Header.vue'
import SubHeader from '../../molecules/SubHeader.vue'
import VerticalNavigation from '../../organisms/organization/VerticalNavigation.vue'
import PaddedAppLeftSidebar from '../../layouts/PaddedAppLeftSidebar.vue'
import { gql, useQuery } from '@urql/vue'
import { useRoute } from 'vue-router'
import OrganizationSettingsHeader from '../../organisms/organization/OrganizationSettingsHeader.vue'
import OrganizationListLicenses, {
  ORGANIZATION_LIST_SINGLE_LICENSE,
} from '../../organisms/organization/OrganizationListLicenses.vue'
import type {
  ManageInstallationPageQuery,
  ManageInstallationPageQueryVariables,
} from './__generated__/ManageInstallationPage'
import TextInput from '../../molecules/TextInput.vue'
import Button from '../../components/shared/Button.vue'
import { useUpdateInstallation } from '../../mutations/useUpdateInstallation'
import { Banner } from '../../atoms'

export default defineComponent({
  components: {
    TextInput,
    OrganizationListLicenses,
    Header,
    SubHeader,
    VerticalNavigation,
    PaddedAppLeftSidebar,
    OrganizationSettingsHeader,
    Button,
    Banner,
  },
  setup() {
    let route = useRoute()

    let { data, executeQuery } = useQuery<
      ManageInstallationPageQuery,
      ManageInstallationPageQueryVariables
    >({
      query: gql`
        query ManageInstallationPage($shortID: ID!) {
          organization(shortID: $shortID) {
            id
            name
            licenses {
              ...OrganizationListSingleLicense
            }
          }
        }
        ${ORGANIZATION_LIST_SINGLE_LICENSE}
      `,
      requestPolicy: 'cache-and-network',
      variables: {
        shortID: computed(() => route.params.organizationSlug as string),
      },
    })

    let updateInstallation = useUpdateInstallation()

    return {
      data,
      updateInstallation,

      reload() {
        executeQuery({
          requestPolicy: 'network-only',
        })
      },
    }
  },
  data() {
    return {
      licenseKeyInput: null,
      failedMessage: null,
    }
  },
  computed: {
    haveLicense() {
      if (this.data.organization.licenses && this.data.organization.licenses.length > 0) {
        return true
      }
      return false
    },
  },
  methods: {
    doUpdateLicenseKey() {
      if (!this.licenseKeyInput) {
        return
      }

      this.failedMessage = null

      this.updateInstallation({ licenseKey: this.licenseKeyInput })
        .catch((err) => {
          if (err?.graphQLErrors && err?.graphQLErrors[0]?.extensions?.message) {
            this.failedMessage = err?.graphQLErrors[0]?.extensions?.message
          } else {
            throw err
          }
        })
        .then(() => {
          this.licenseKeyInput = null
        })
        .finally(() => {
          this.reload()
        })
    },
  },
})
</script>
