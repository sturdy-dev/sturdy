<template>
  <div>
    <div class="text-sm text-gray-500 text-center">
      <div v-if="isMultiTenancyEnabled">Logging in to Sturdy Cloud</div>
      <div v-else>
        Logging in to {{ serverHost }}<br />
        <span v-if="serverInfo" class="capitalize">
          Sturdy {{ serverInfo.distributionType }} {{ serverInfo.version }}
        </span>
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { gql } from '@urql/vue'
import { computed, defineProps, inject, Ref, ref } from 'vue'
import { Feature } from '../../__generated__/types'
import { ServerInfoFragment } from './__generated__/ServerInfo'

const SERVER_INFO = gql`
  fragment ServerInfo on Installation {
    version
    distributionType
  }
`

const features = inject<Ref<Array<Feature>>>('features', ref([]))
const isMultiTenancyEnabled = computed(() => features?.value?.includes(Feature.MultiTenancy))

interface Props {
  serverInfo?: ServerInfoFragment
}

const props = defineProps<Props>()

let serverHost = window.location.host
</script>
